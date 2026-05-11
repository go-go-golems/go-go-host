package cmds

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIHTTPHelpersSmoke(t *testing.T) {
	bundle := writeSmokeBundle(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Go-Go-Host-User"); got != "alice" {
			t.Fatalf("expected dev user header, got %q", got)
		}
		switch r.URL.Path {
		case "/api/v1/me":
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
		case "/api/v1/sites/site_123/deployments":
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				t.Fatalf("parse multipart: %v", err)
			}
			if _, _, err := r.FormFile("bundle"); err != nil {
				t.Fatalf("missing bundle: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"deployment": map[string]any{"id": "dep_123", "siteId": "site_123", "status": "validated"}, "report": map[string]any{"valid": true}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	var me map[string]any
	if err := getJSONWithAuth(server.URL, "/api/v1/me", "alice", "", &me); err != nil {
		t.Fatalf("get json: %v", err)
	}
	var deploy map[string]any
	if err := postMultipartBundleWithAuth(server.URL, "/api/v1/sites/site_123/deployments", "alice", "", bundle, map[string]string{"message": "smoke"}, &deploy); err != nil {
		t.Fatalf("multipart: %v", err)
	}
	if _, ok := deploy["deployment"]; !ok {
		t.Fatalf("expected deployment response, got %#v", deploy)
	}
}

func TestCLIHTTPErrorIncludesResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON := `{"error":"validation failed"}`
		http.Error(w, writeJSON, http.StatusBadRequest)
	}))
	defer server.Close()
	var out map[string]any
	err := getJSONWithAuth(server.URL, "/fail", "alice", "", &out)
	if err == nil || !strings.Contains(err.Error(), "validation failed") {
		t.Fatalf("expected response body in error, got %v", err)
	}
}

func writeSmokeBundle(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bundle.tar.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	files := map[string]string{"go-go-host.json": `{"scriptsDir":"scripts"}`, "scripts/app.js": ""}
	for name, body := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(body)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}
