package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
	"github.com/google/uuid"
)

func newIntegrationHandlerWithConfig(t *testing.T, configure func(*config.Config)) http.Handler {
	t.Helper()
	dsn := os.Getenv("GO_GO_HOST_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("GO_GO_HOST_TEST_DATABASE_URL is not set; skipping Postgres integration test")
	}
	ctx := context.Background()
	st, err := store.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(st.Close)
	if err := st.ApplyMigrations(ctx); err != nil {
		t.Fatalf("migrations: %v", err)
	}
	cfg := config.Default()
	cfg.BaseDomain = "example.test"
	cfg.DevAuth = true
	cfg.DataDir = t.TempDir()
	if configure != nil {
		configure(&cfg)
	}
	return NewHandler(control.NewCoreWithStore(cfg, st))
}

func TestAdminInventoryRequiresPlatformAdminAndListsTenants(t *testing.T) {
	suffix := uuid.NewString()[:8]
	adminUser := "platform-admin-" + suffix
	tenantUser := "tenant-" + suffix
	h := newIntegrationHandlerWithConfig(t, func(cfg *config.Config) { cfg.DevPlatformAdminSubjects = []string{adminUser} })
	org := createTestOrgViaAPI(t, h, tenantUser, "admin-inv-org-"+suffix)
	_ = createTestSiteViaAPI(t, h, tenantUser, org.ID, "admin-inv-site-"+suffix)

	for _, path := range []string{"/api/v1/admin/orgs", "/api/v1/admin/users", "/api/v1/admin/sites", "/api/v1/admin/deployments"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("X-Go-Go-Host-User", tenantUser)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("%s non-admin: expected 403, got %d body=%s", path, rec.Code, rec.Body.String())
		}
	}

	orgReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orgs", nil)
	orgReq.Header.Set("X-Go-Go-Host-User", adminUser)
	orgRec := httptest.NewRecorder()
	h.ServeHTTP(orgRec, orgReq)
	if orgRec.Code != http.StatusOK {
		t.Fatalf("admin orgs: expected 200, got %d body=%s", orgRec.Code, orgRec.Body.String())
	}
	var orgs []adminOrgDTO
	if err := json.Unmarshal(orgRec.Body.Bytes(), &orgs); err != nil {
		t.Fatal(err)
	}
	foundOrg := false
	for _, row := range orgs {
		if row.ID == org.ID && row.SiteCount >= 1 {
			foundOrg = true
		}
	}
	if !foundOrg {
		t.Fatalf("expected admin org inventory to include %#v with site count, got %#v", org, orgs)
	}

	siteReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/sites", nil)
	siteReq.Header.Set("X-Go-Go-Host-User", adminUser)
	siteRec := httptest.NewRecorder()
	h.ServeHTTP(siteRec, siteReq)
	if siteRec.Code != http.StatusOK {
		t.Fatalf("admin sites: expected 200, got %d body=%s", siteRec.Code, siteRec.Body.String())
	}
	var sites []adminSiteDTO
	if err := json.Unmarshal(siteRec.Body.Bytes(), &sites); err != nil {
		t.Fatal(err)
	}
	foundSite := false
	for _, row := range sites {
		if row.OrgID == org.ID {
			foundSite = true
		}
	}
	if !foundSite {
		t.Fatalf("expected admin site inventory to include org %s, got %#v", org.ID, sites)
	}
}
