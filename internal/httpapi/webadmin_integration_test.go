package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
)

func TestDashboardRoutesServeEmbeddedSPAAndAPIRoutesStillWork(t *testing.T) {
	h := NewHandler(control.NewCore(config.Config{BaseDomain: "localhost", PublicBaseURL: "http://127.0.0.1:8080", DevAuth: true}))
	for _, path := range []string{"/app", "/app/", "/app/orgs/org_123/sites"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "go-go-host") {
			t.Fatalf("%s: expected dashboard index", path)
		}
	}
	apiRec := httptest.NewRecorder()
	h.ServeHTTP(apiRec, httptest.NewRequest(http.MethodGet, "/api/v1/version", nil))
	if apiRec.Code != http.StatusOK || !strings.Contains(apiRec.Body.String(), Version) {
		t.Fatalf("expected version JSON, got code=%d body=%s", apiRec.Code, apiRec.Body.String())
	}
}
