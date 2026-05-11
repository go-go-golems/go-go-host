package runtime

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSupervisorRoutesByHost(t *testing.T) {
	ctx := context.Background()
	sup := NewSupervisor()
	if err := sup.Activate(ctx, fixtureSpec(t, "site_a", "a.localhost")); err != nil {
		t.Fatalf("activate a: %v", err)
	}
	if err := sup.Activate(ctx, fixtureSpec(t, "site_b", "b.localhost")); err != nil {
		t.Fatalf("activate b: %v", err)
	}

	for _, host := range []string{"a.localhost", "b.localhost"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://"+host+"/", nil)
		req.Host = host
		sup.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d body %s", host, rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "Hello from go-go-host") {
			t.Fatalf("%s: expected hello body, got %s", host, rec.Body.String())
		}
	}
	if got := sup.Summary(); got.ActiveSites != 2 {
		t.Fatalf("expected 2 active sites, got %#v", got)
	}
}

func TestSupervisorUnknownHost404(t *testing.T) {
	ctx := context.Background()
	sup := NewSupervisor()
	if err := sup.Activate(ctx, fixtureSpec(t, "site_a", "a.localhost")); err != nil {
		t.Fatalf("activate: %v", err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://unknown.localhost/", nil)
	req.Host = "unknown.localhost"
	sup.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestSupervisorFailedActivationDoesNotReplaceCurrentRuntime(t *testing.T) {
	ctx := context.Background()
	sup := NewSupervisor()
	if err := sup.Activate(ctx, fixtureSpec(t, "site_a", "a.localhost")); err != nil {
		t.Fatalf("activate good: %v", err)
	}
	bad := fixtureSpec(t, "site_a", "a.localhost")
	bad.ScriptsDir = t.TempDir() + "/missing"
	if err := sup.Activate(ctx, bad); err == nil {
		t.Fatalf("expected bad activation to fail")
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://a.localhost/", nil)
	req.Host = "a.localhost"
	sup.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected old runtime to remain serving, got %d body %s", rec.Code, rec.Body.String())
	}
	st, ok := sup.Status("site_a")
	if !ok || st.Status != StatusFailed {
		t.Fatalf("expected failed status, got %#v ok=%v", st, ok)
	}
}

func TestSupervisorStop(t *testing.T) {
	ctx := context.Background()
	sup := NewSupervisor()
	if err := sup.Activate(ctx, fixtureSpec(t, "site_a", "a.localhost")); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if err := sup.Stop(ctx, "site_a"); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := sup.Stop(ctx, "site_a"); !errors.Is(err, ErrRuntimeNotFound) {
		t.Fatalf("expected runtime not found, got %v", err)
	}
}

func fixtureSpec(t *testing.T, siteID, host string) Spec {
	t.Helper()
	return Spec{
		SiteID:       siteID,
		OrgID:        "org_test",
		DeploymentID: "dep_" + siteID,
		Hosts:        []string{host},
		ScriptsDir:   "testdata/sites/hello/scripts",
		AssetsDir:    "testdata/sites/hello/assets",
		DBPath:       t.TempDir() + "/app.sqlite",
		Dev:          true,
		Capabilities: DefaultCapabilities(),
	}
}
