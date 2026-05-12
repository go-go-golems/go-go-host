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
}
