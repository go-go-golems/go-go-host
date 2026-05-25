package deploy

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateAndStoreRejectsBadPaths(t *testing.T) {
	bundle := writeTestTarGz(t, map[string]string{"go-go-host.json": `{"scriptsDir":"scripts"}`, "../evil.js": "bad"})
	prepared, err := ValidateAndStore(context.Background(), bundle, filepath.Join(t.TempDir(), "bundle.tgz"), filepath.Join(t.TempDir(), "unpack"), Options{MaxBytes: 1024 * 1024})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if prepared.Report.Valid || !containsError(prepared.Report.Errors, "invalid path") {
		t.Fatalf("expected invalid path report, got %#v", prepared.Report)
	}
}

func TestValidateAndStoreAcceptsDotSlashManifest(t *testing.T) {
	bundle := writeTestTarGz(t, map[string]string{"./go-go-host.json": `{"scriptsDir":"scripts","entrypoint":"scripts/app.js","smokePath":"/"}`, "./scripts/app.js": "console.log('hi')"})
	prepared, err := ValidateAndStore(context.Background(), bundle, filepath.Join(t.TempDir(), "bundle.tgz"), filepath.Join(t.TempDir(), "unpack"), Options{MaxBytes: 1024 * 1024})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if !prepared.Report.Valid {
		t.Fatalf("expected valid report, got %#v", prepared.Report)
	}
}

func TestValidateAndStoreAllowsDoubleStarPolicy(t *testing.T) {
	bundle := writeTestTarGz(t, map[string]string{"go-go-host.json": `{"scriptsDir":"scripts","entrypoint":"scripts/app.js","smokePath":"/"}`, "scripts/app.js": "console.log('hi')", "assets/style.css": "body{}"})
	prepared, err := ValidateAndStore(context.Background(), bundle, filepath.Join(t.TempDir(), "bundle.tgz"), filepath.Join(t.TempDir(), "unpack"), Options{MaxBytes: 1024 * 1024, AllowedPaths: []string{"**"}})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if !prepared.Report.Valid {
		t.Fatalf("expected double-star policy to allow nested archive paths, got %#v", prepared.Report)
	}
}

func TestValidateAndStoreRejectsMissingManifest(t *testing.T) {
	bundle := writeTestTarGz(t, map[string]string{"scripts/app.js": "console.log('hi')"})
	prepared, err := ValidateAndStore(context.Background(), bundle, filepath.Join(t.TempDir(), "bundle.tgz"), filepath.Join(t.TempDir(), "unpack"), Options{MaxBytes: 1024 * 1024})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if prepared.Report.Valid || !containsError(prepared.Report.Errors, "missing go-go-host.json") {
		t.Fatalf("expected missing manifest report, got %#v", prepared.Report)
	}
}

func TestValidateAndStoreRejectsOversizedBundle(t *testing.T) {
	bundle := writeTestTarGz(t, map[string]string{"go-go-host.json": `{"scriptsDir":"scripts"}`, "scripts/app.js": strings.Repeat("x", 128)})
	prepared, err := ValidateAndStore(context.Background(), bundle, filepath.Join(t.TempDir(), "bundle.tgz"), filepath.Join(t.TempDir(), "unpack"), Options{MaxBytes: 32})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if prepared.Report.Valid || !containsError(prepared.Report.Errors, "exceed") {
		t.Fatalf("expected oversized report, got %#v", prepared.Report)
	}
}

func TestValidateAndStoreRejectsForbiddenCapability(t *testing.T) {
	bundle := writeTestTarGz(t, map[string]string{"go-go-host.json": `{"scriptsDir":"scripts","capabilities":["exec"]}`, "scripts/app.js": "console.log('hi')"})
	prepared, err := ValidateAndStore(context.Background(), bundle, filepath.Join(t.TempDir(), "bundle.tgz"), filepath.Join(t.TempDir(), "unpack"), Options{MaxBytes: 1024 * 1024})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if prepared.Report.Valid || !containsError(prepared.Report.Errors, "capability") {
		t.Fatalf("expected forbidden capability report, got %#v", prepared.Report)
	}
}

func writeTestTarGz(t *testing.T, files map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bundle.tar.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
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

func containsError(errors []string, needle string) bool {
	for _, err := range errors {
		if strings.Contains(err, needle) {
			return true
		}
	}
	return false
}
