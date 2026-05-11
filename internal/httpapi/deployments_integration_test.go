package httpapi

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestDeploymentUploadActivateAndServeFlow(t *testing.T) {
	h := newIntegrationHandler(t)
	suffix := uuid.NewString()[:8]
	user := "deploy-" + suffix

	org := createTestOrgViaAPI(t, h, user, "deploy-org-"+suffix)
	site := createTestSiteViaAPI(t, h, user, org.ID, "deploy-site-"+suffix)
	bundle := writeHelloBundle(t)

	uploadRec := uploadBundleViaAPI(t, h, user, site.ID, bundle)
	if uploadRec.Code != http.StatusCreated {
		t.Fatalf("upload status: got %d body %s", uploadRec.Code, uploadRec.Body.String())
	}
	var upload struct {
		Deployment deploymentDTO `json:"deployment"`
	}
	if err := json.Unmarshal(uploadRec.Body.Bytes(), &upload); err != nil {
		t.Fatalf("decode upload: %v", err)
	}
	if upload.Deployment.Status != "validated" {
		t.Fatalf("expected validated deployment, got %#v", upload.Deployment)
	}

	activateReq := httptest.NewRequest(http.MethodPost, "/api/v1/deployments/"+upload.Deployment.ID+"/activate", nil)
	activateReq.Header.Set("X-Go-Go-Host-User", user)
	activateRec := httptest.NewRecorder()
	h.ServeHTTP(activateRec, activateReq)
	if activateRec.Code != http.StatusOK {
		t.Fatalf("activate status: got %d body %s", activateRec.Code, activateRec.Body.String())
	}

	publicReq := httptest.NewRequest(http.MethodGet, "http://"+site.PrimaryHost+"/", nil)
	publicReq.Host = site.PrimaryHost
	publicRec := httptest.NewRecorder()
	h.ServeHTTP(publicRec, publicReq)
	if publicRec.Code != http.StatusOK {
		t.Fatalf("public status: got %d body %s", publicRec.Code, publicRec.Body.String())
	}
	if !bytes.Contains(publicRec.Body.Bytes(), []byte("Hello")) {
		t.Fatalf("expected hello body, got %s", publicRec.Body.String())
	}
}

func createTestOrgViaAPI(t *testing.T, h http.Handler, user, slug string) orgDTO {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orgs", bytes.NewReader([]byte(`{"slug":"`+slug+`","name":"Deploy Org"}`)))
	req.Header.Set("X-Go-Go-Host-User", user)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("org: %d %s", rec.Code, rec.Body.String())
	}
	var org orgDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &org); err != nil {
		t.Fatal(err)
	}
	return org
}

func createTestSiteViaAPI(t *testing.T, h http.Handler, user, orgID, slug string) siteDTO {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+orgID+"/sites", bytes.NewReader([]byte(`{"slug":"`+slug+`","name":"Deploy Site"}`)))
	req.Header.Set("X-Go-Go-Host-User", user)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("site: %d %s", rec.Code, rec.Body.String())
	}
	var site siteDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &site); err != nil {
		t.Fatal(err)
	}
	return site
}

func uploadBundleViaAPI(t *testing.T, h http.Handler, user, siteID, bundlePath string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	part, err := mw.CreateFormFile("bundle", filepath.Base(bundlePath))
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(part, f); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	_ = mw.WriteField("message", "integration test")
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+siteID+"/deployments", &body)
	req.Header.Set("X-Go-Go-Host-User", user)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func writeHelloBundle(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "hello.tar.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	files := map[string]string{
		"go-go-host.json":  `{"scriptsDir":"scripts","assetsDir":"assets","smokePath":"/","capabilities":["time","timer"]}`,
		"scripts/app.js":   `const express = require('express'); const app = express.app(); app.get('/', (req, res) => '<h1>Hello from deployed bundle</h1>');`,
		"assets/style.css": "body { color: black; }",
	}
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
