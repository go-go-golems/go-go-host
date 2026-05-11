package httpapi

import (
	"bytes"
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

func newIntegrationHandler(t *testing.T) http.Handler {
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
	return NewHandler(control.NewCoreWithStore(cfg, st))
}

func TestDevAuthMeOrgAndSiteFlow(t *testing.T) {
	h := newIntegrationHandler(t)
	suffix := uuid.NewString()[:8]

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	meReq.Header.Set("X-Go-Go-Host-User", "api-"+suffix)
	meRec := httptest.NewRecorder()
	h.ServeHTTP(meRec, meReq)
	if meRec.Code != http.StatusOK {
		t.Fatalf("me status: got %d body %s", meRec.Code, meRec.Body.String())
	}

	orgBody := []byte(`{"slug":"api-org-` + suffix + `","name":"API Org"}`)
	orgReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs", bytes.NewReader(orgBody))
	orgReq.Header.Set("X-Go-Go-Host-User", "api-"+suffix)
	orgReq.Header.Set("Content-Type", "application/json")
	orgRec := httptest.NewRecorder()
	h.ServeHTTP(orgRec, orgReq)
	if orgRec.Code != http.StatusCreated {
		t.Fatalf("org status: got %d body %s", orgRec.Code, orgRec.Body.String())
	}
	var org orgDTO
	if err := json.Unmarshal(orgRec.Body.Bytes(), &org); err != nil {
		t.Fatalf("decode org: %v", err)
	}

	siteBody := []byte(`{"slug":"api-site-` + suffix + `","name":"API Site"}`)
	siteReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+org.ID+"/sites", bytes.NewReader(siteBody))
	siteReq.Header.Set("X-Go-Go-Host-User", "api-"+suffix)
	siteReq.Header.Set("Content-Type", "application/json")
	siteRec := httptest.NewRecorder()
	h.ServeHTTP(siteRec, siteReq)
	if siteRec.Code != http.StatusCreated {
		t.Fatalf("site status: got %d body %s", siteRec.Code, siteRec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/orgs/"+org.ID+"/sites", nil)
	listReq.Header.Set("X-Go-Go-Host-User", "api-"+suffix)
	listRec := httptest.NewRecorder()
	h.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status: got %d body %s", listRec.Code, listRec.Body.String())
	}
	var sites []siteDTO
	if err := json.Unmarshal(listRec.Body.Bytes(), &sites); err != nil {
		t.Fatalf("decode sites: %v", err)
	}
	if len(sites) != 1 {
		t.Fatalf("expected 1 site, got %d", len(sites))
	}
}
