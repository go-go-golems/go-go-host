package runtime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

func TestSiteRuntimeRendersFixture(t *testing.T) {
	ctx := context.Background()
	rt := newFixtureRuntime(t, ctx)
	defer rt.Close(ctx)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://hello.localhost/", nil)
	rt.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Hello from go-go-host") {
		t.Fatalf("expected rendered hello page, got %s", rec.Body.String())
	}
}

func TestDatabaseConfigureDisabled(t *testing.T) {
	ctx := context.Background()
	rt := newFixtureRuntime(t, ctx)
	defer rt.Close(ctx)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://hello.localhost/config-test", nil)
	rt.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"ok":true`) {
		t.Fatalf("expected configure failure to be reported as ok, got %s", rec.Body.String())
	}
}

func TestRuntimeHealthCheck(t *testing.T) {
	ctx := context.Background()
	rt := newFixtureRuntime(t, ctx)
	defer rt.Close(ctx)

	if err := rt.HealthCheck(ctx); err != nil {
		t.Fatalf("health check failed: %v", err)
	}
}

func TestExecAndFSUnavailableByDefault(t *testing.T) {
	ctx := context.Background()
	rt := newFixtureRuntime(t, ctx)
	defer rt.Close(ctx)

	for _, moduleName := range []string{"exec", "fs"} {
		moduleName := moduleName
		_, err := rt.runtime.Owner.Call(ctx, "require-"+moduleName, func(_ context.Context, vm *goja.Runtime) (any, error) {
			_, err := vm.RunString("require('" + moduleName + "')")
			return nil, err
		})
		if err == nil {
			t.Fatalf("expected require(%q) to fail", moduleName)
		}
	}
}

func newFixtureRuntime(t *testing.T, ctx context.Context) *SiteRuntime {
	t.Helper()
	rt, err := NewSiteRuntime(ctx, Spec{
		SiteID:       "site_test",
		OrgID:        "org_test",
		DeploymentID: "dep_test",
		Hosts:        []string{"hello.localhost"},
		ScriptsDir:   "testdata/sites/hello/scripts",
		AssetsDir:    "testdata/sites/hello/assets",
		DBPath:       t.TempDir() + "/app.sqlite",
		Dev:          true,
		Capabilities: DefaultCapabilities(),
	})
	if err != nil {
		t.Fatalf("create runtime: %v", err)
	}
	return rt
}
