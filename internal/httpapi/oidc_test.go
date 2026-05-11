package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
)

func TestOIDCModeRejectsMissingBearerToken(t *testing.T) {
	cfg := config.Default()
	cfg.DevAuth = false
	cfg.OIDCIssuer = "http://127.0.0.1:1/realms/test"
	cfg.OIDCClientID = "go-go-host-dashboard"
	h := NewHandler(control.NewCoreWithStore(cfg, nil))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized && rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected auth failure, got %d body %s", rec.Code, rec.Body.String())
	}
}
