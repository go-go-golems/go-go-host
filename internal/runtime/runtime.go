package runtime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/go-go-golems/go-go-goja/engine"
	databasemod "github.com/go-go-golems/go-go-goja/modules/database"
	"github.com/go-go-golems/go-go-host/internal/sitejs/dbguard"
	"github.com/go-go-golems/go-go-host/internal/sitejs/uidsl"
	"github.com/go-go-golems/go-go-host/internal/sitejs/web"
	_ "github.com/mattn/go-sqlite3"
)

type Spec struct {
	SiteID       string
	OrgID        string
	DeploymentID string
	Hosts        []string
	ScriptsDir   string
	AssetsDir    string
	DBPath       string
	Dev          bool
	HealthPath   string
	Capabilities CapabilitySet
}

type CapabilitySet struct {
	Database bool
	Timers   bool
	Assets   bool
}

func DefaultCapabilities() CapabilitySet {
	return CapabilitySet{Database: true, Timers: true, Assets: true}
}

type SiteRuntime struct {
	spec    Spec
	db      *sql.DB
	guard   *dbguard.Guard
	runtime *engine.Runtime
	host    *web.Host
	started time.Time
}

func NewSiteRuntime(ctx context.Context, spec Spec) (*SiteRuntime, error) {
	if spec.ScriptsDir == "" {
		return nil, fmt.Errorf("scripts dir is required")
	}
	if spec.DBPath == "" {
		return nil, fmt.Errorf("db path is required")
	}
	if err := os.MkdirAll(filepath.Dir(spec.DBPath), 0o755); err != nil && filepath.Dir(spec.DBPath) != "." {
		return nil, fmt.Errorf("create site db directory: %w", err)
	}
	db, err := sql.Open("sqlite3", spec.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open site sqlite database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping site sqlite database: %w", err)
	}

	host := web.NewHost(web.HostOptions{Dev: spec.Dev, Renderer: uidsl.RenderAny})
	guard := dbguard.New(db, spec.DBPath)
	meteredDB := dbguard.NewMeteredDB(db, guard)
	databaseModule := databasemod.New(databasemod.WithPreconfiguredDB(meteredDB), databasemod.WithConfigureEnabled(false))
	dbAliasModule := databasemod.New(databasemod.WithName("db"), databasemod.WithPreconfiguredDB(meteredDB), databasemod.WithConfigureEnabled(false))

	builder := engine.NewBuilder().WithModules(
		engine.NativeModuleSpec{ModuleID: "database:app", ModuleName: databaseModule.Name(), Loader: databaseModule.Loader},
		engine.NativeModuleSpec{ModuleID: "database:db-alias", ModuleName: dbAliasModule.Name(), Loader: dbAliasModule.Loader},
	)
	middleware := []string{"path"}
	if spec.Capabilities.Timers {
		middleware = append(middleware, "time", "timer")
	}
	builder = builder.UseModuleMiddleware(engine.MiddlewareOnly(middleware...)).WithRuntimeModuleRegistrars(web.NewExpressRegistrar(host), uidsl.NewRegistrar(), dbguard.NewRegistrar(guard))
	factory, err := builder.Build()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("build site runtime factory: %w", err)
	}
	rt, err := factory.NewRuntime(ctx)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create site runtime: %w", err)
	}
	host.SetRuntime(rt.Owner)
	if spec.AssetsDir != "" && spec.Capabilities.Assets {
		host.RegisterStatic("/assets", spec.AssetsDir)
	}
	sr := &SiteRuntime{spec: spec, db: db, guard: guard, runtime: rt, host: host, started: time.Now().UTC()}
	if err := sr.LoadScripts(ctx); err != nil {
		_ = sr.Close(ctx)
		return nil, err
	}
	return sr, nil
}

func (s *SiteRuntime) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.host.ServeHTTP(w, r) }
func (s *SiteRuntime) Handler() http.Handler                            { return s.host }
func (s *SiteRuntime) Spec() Spec                                       { return s.spec }
func (s *SiteRuntime) StartedAt() time.Time                             { return s.started }
func (s *SiteRuntime) DBGuard() *dbguard.Guard                          { return s.guard }

func (s *SiteRuntime) HealthCheck(ctx context.Context) error {
	path := s.spec.HealthPath
	if path == "" {
		path = "/"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://site-runtime.local"+path, nil)
	if err != nil {
		return err
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if rec.Code < 200 || rec.Code >= 400 {
		return fmt.Errorf("runtime health check %s returned status %d: %s", path, rec.Code, strings.TrimSpace(rec.Body.String()))
	}
	return nil
}

func (s *SiteRuntime) Close(ctx context.Context) error {
	var errs []error
	if s.runtime != nil {
		if err := s.runtime.Close(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (s *SiteRuntime) LoadScripts(ctx context.Context) error {
	files, err := scriptFiles(s.spec.ScriptsDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		file := file
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read script %s: %w", file, err)
		}
		_, err = s.runtime.Owner.Call(ctx, "load-script", func(_ context.Context, vm *goja.Runtime) (any, error) {
			_, err := vm.RunScript(file, string(data))
			return nil, err
		})
		if err != nil {
			return fmt.Errorf("execute script %s: %w", file, err)
		}
	}
	return nil
}

func scriptFiles(dir string) ([]string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("stat scripts directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scripts path %s is not a directory", dir)
	}
	var files []string
	if err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".js") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk scripts directory: %w", err)
	}
	sort.Strings(files)
	return files, nil
}
