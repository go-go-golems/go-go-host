package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
	hostruntime "github.com/go-go-golems/go-go-host/internal/runtime"
)

func TestHandlerFallsBackToSupervisorForPublicHostTraffic(t *testing.T) {
	ctx := context.Background()
	core := control.NewCore(config.Default())
	if err := core.Supervisor.Activate(ctx, hostruntime.Spec{
		SiteID:       "site_public",
		OrgID:        "org_public",
		DeploymentID: "dep_public",
		Hosts:        []string{"public.localhost"},
		ScriptsDir:   "../runtime/testdata/sites/hello/scripts",
		AssetsDir:    "../runtime/testdata/sites/hello/assets",
		DBPath:       t.TempDir() + "/app.sqlite",
		Dev:          true,
		Capabilities: hostruntime.DefaultCapabilities(),
	}); err != nil {
		t.Fatalf("activate fixture: %v", err)
	}
	h := NewHandler(core)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://public.localhost/", nil)
	req.Host = "public.localhost"
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected public host to be served, got %d body %s", rec.Code, rec.Body.String())
	}
}
