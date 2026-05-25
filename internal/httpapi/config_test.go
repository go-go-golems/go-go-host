package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-go-golems/go-go-host/internal/config"
	"github.com/go-go-golems/go-go-host/internal/control"
)

func TestConfigResponseIncludesOIDCWhenDevAuthDisabled(t *testing.T) {
	cfg := config.Default()
	cfg.DevAuth = false
	cfg.OIDCIssuer = "http://issuer.example/realms/go-go-host"
	cfg.OIDCClientID = "dashboard"
	cfg.OIDCDeviceClientID = "cli"
	h := NewHandler(control.NewCoreWithStore(cfg, nil))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/config", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d body %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode config: %v", err)
	}
	oidc, ok := body["oidc"].(map[string]any)
	if !ok {
		t.Fatalf("expected oidc config in response: %#v", body)
	}
	if oidc["issuer"] != cfg.OIDCIssuer || oidc["clientId"] != cfg.OIDCClientID || oidc["deviceClientId"] != cfg.OIDCDeviceClientID {
		t.Fatalf("unexpected oidc config: %#v", oidc)
	}
}
