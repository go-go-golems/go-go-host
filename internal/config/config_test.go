package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExpandsEnvironmentVariables(t *testing.T) {
	t.Setenv("GO_GO_HOST_TEST_DSN", "postgres://user:pass@example.test:5432/go_go_host?sslmode=disable")
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("controlDbDsn: ${GO_GO_HOST_TEST_DSN}\ndevAuth: false\noidcIssuer: http://issuer.example\noidcClientId: dashboard\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.ControlDBDSN != os.Getenv("GO_GO_HOST_TEST_DSN") {
		t.Fatalf("expected expanded dsn, got %q", cfg.ControlDBDSN)
	}
	if cfg.OIDCDeviceClientID != "go-go-host-cli" {
		t.Fatalf("expected default device client id, got %q", cfg.OIDCDeviceClientID)
	}
	if got := cfg.OIDCAcceptedClientIDs; len(got) != 2 || got[0] != "dashboard" || got[1] != "go-go-host-cli" {
		t.Fatalf("expected dashboard and default CLI accepted clients, got %#v", got)
	}
}

func TestLoadPreservesExplicitAcceptedOIDCClientIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("devAuth: false\noidcIssuer: http://issuer.example\noidcClientId: dashboard\noidcDeviceClientId: custom-cli\noidcAcceptedClientIds:\n  - dashboard\n  - custom-cli\n  - another-client\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got := cfg.OIDCAcceptedClientIDs; len(got) != 3 || got[1] != "custom-cli" || got[2] != "another-client" {
		t.Fatalf("expected explicit accepted clients to be preserved, got %#v", got)
	}
}
