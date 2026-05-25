package webadmin

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHandlerServesIndexForNestedRoutes(t *testing.T) {
	h := NewHandler()
	for _, path := range []string{"/", "/orgs/org_123/sites", "/orgs/org_123/sites/site_123/deployments"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "go-go-host") {
			t.Fatalf("%s: expected dashboard index", path)
		}
	}
}

func TestNewHandlerServesAssets(t *testing.T) {
	assetPath := ""
	if err := fs.WalkDir(dashboardFS, "dist/assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || assetPath != "" {
			return err
		}
		if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
			assetPath = strings.TrimPrefix(path, "dist")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if assetPath == "" {
		t.Fatal("expected embedded asset")
	}
	h := NewHandler()
	req := httptest.NewRequest(http.MethodGet, assetPath, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected asset 200, got %d", rec.Code)
	}
}
