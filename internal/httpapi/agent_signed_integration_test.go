package httpapi

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/google/uuid"
)

func TestAgentSignedDeployRunSecurity(t *testing.T) {
	h := newIntegrationHandler(t)
	suffix := uuid.NewString()[:8]
	user := "signed-agent-" + suffix
	org := createTestOrgViaAPI(t, h, user, "signed-agent-org-"+suffix)
	site := createTestSiteViaAPI(t, h, user, org.ID, "signed-agent-site-"+suffix)
	otherSite := createTestSiteViaAPI(t, h, user, org.ID, "signed-agent-other-"+suffix)

	createBody := []byte(`{"name":"ci","siteId":"` + site.ID + `","allowedChannels":["default"],"allowedBundlePaths":["bundles/**"],"canActivate":true}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+org.ID+"/agents", bytes.NewReader(createBody))
	createReq.Header.Set("X-Go-Go-Host-User", user)
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	h.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create agent: %d %s", createRec.Code, createRec.Body.String())
	}
	var created createAgentResponse
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	enrollBody := []byte(`{"token":"` + created.EnrollmentToken + `","publicKey":"` + base64.StdEncoding.EncodeToString(pub) + `"}`)
	enrollReq := httptest.NewRequest(http.MethodPost, "/api/v1/agent/enroll", bytes.NewReader(enrollBody))
	enrollReq.Header.Set("Content-Type", "application/json")
	enrollRec := httptest.NewRecorder()
	h.ServeHTTP(enrollRec, enrollReq)
	if enrollRec.Code != http.StatusCreated {
		t.Fatalf("enroll: %d %s", enrollRec.Code, enrollRec.Body.String())
	}
	var enrolled enrollAgentResponse
	if err := json.Unmarshal(enrollRec.Body.Bytes(), &enrolled); err != nil {
		t.Fatal(err)
	}

	goodBody := []byte(`{"siteId":"` + site.ID + `","channel":"default","bundlePath":"bundles/app.tar.gz","action":"deploy","activate":true}`)
	goodReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-good")
	goodRec := httptest.NewRecorder()
	h.ServeHTTP(goodRec, goodReq)
	if goodRec.Code != http.StatusCreated {
		t.Fatalf("good signed run: %d %s", goodRec.Code, goodRec.Body.String())
	}
	var goodRun createDeployRunResponse
	if err := json.Unmarshal(goodRec.Body.Bytes(), &goodRun); err != nil {
		t.Fatal(err)
	}
	uploadRec := uploadAgentBundleViaAPI(t, h, goodRun.ID, goodRun.UploadToken, writeHelloBundle(t))
	if uploadRec.Code != http.StatusCreated {
		t.Fatalf("agent upload: %d %s", uploadRec.Code, uploadRec.Body.String())
	}
	var upload agentUploadResponse
	if err := json.Unmarshal(uploadRec.Body.Bytes(), &upload); err != nil {
		t.Fatal(err)
	}
	if !upload.Activated || upload.Deployment.Status != "active" {
		t.Fatalf("expected auto-activated upload, got %#v", upload)
	}
	secondUploadRec := uploadAgentBundleViaAPI(t, h, goodRun.ID, goodRun.UploadToken, writeHelloBundle(t))
	if secondUploadRec.Code != http.StatusForbidden {
		t.Fatalf("expected second upload denial, got %d %s", secondUploadRec.Code, secondUploadRec.Body.String())
	}
	keysReq := httptest.NewRequest(http.MethodGet, "/api/v1/orgs/"+org.ID+"/agents/"+created.Agent.ID+"/keys", nil)
	keysReq.Header.Set("X-Go-Go-Host-User", user)
	keysRec := httptest.NewRecorder()
	h.ServeHTTP(keysRec, keysReq)
	if keysRec.Code != http.StatusOK {
		t.Fatalf("list keys: %d %s", keysRec.Code, keysRec.Body.String())
	}
	var keys []agentKeyDTO
	if err := json.Unmarshal(keysRec.Body.Bytes(), &keys); err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0].LastUsedAt == "" || keys[0].Fingerprint == "" {
		t.Fatalf("expected one used key with fingerprint, got %#v", keys)
	}

	replayReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-good")
	replayRec := httptest.NewRecorder()
	h.ServeHTTP(replayRec, replayReq)
	if replayRec.Code != http.StatusBadRequest {
		t.Fatalf("expected replay denial, got %d %s", replayRec.Code, replayRec.Body.String())
	}

	oldReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC().Add(-10*time.Minute), "nonce-old")
	oldRec := httptest.NewRecorder()
	h.ServeHTTP(oldRec, oldReq)
	if oldRec.Code != http.StatusBadRequest {
		t.Fatalf("expected old timestamp denial, got %d %s", oldRec.Code, oldRec.Body.String())
	}

	wrongPathBody := []byte(`{"siteId":"` + site.ID + `","channel":"default","bundlePath":"private/app.tar.gz","action":"deploy"}`)
	wrongPathReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", wrongPathBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-wrong-path")
	wrongPathRec := httptest.NewRecorder()
	h.ServeHTTP(wrongPathRec, wrongPathReq)
	if wrongPathRec.Code != http.StatusForbidden {
		t.Fatalf("expected wrong path denial, got %d %s", wrongPathRec.Code, wrongPathRec.Body.String())
	}

	wrongSiteBody := []byte(`{"siteId":"` + otherSite.ID + `","channel":"default","bundlePath":"bundles/app.tar.gz","action":"deploy"}`)
	wrongSiteReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", wrongSiteBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-wrong-site")
	wrongSiteRec := httptest.NewRecorder()
	h.ServeHTTP(wrongSiteRec, wrongSiteReq)
	if wrongSiteRec.Code != http.StatusForbidden {
		t.Fatalf("expected wrong site denial, got %d %s", wrongSiteRec.Code, wrongSiteRec.Body.String())
	}

	revokeKeyReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+org.ID+"/agents/"+created.Agent.ID+"/keys/"+enrolled.KeyID+"/revoke", bytes.NewReader([]byte(`{"reason":"integration test"}`)))
	revokeKeyReq.Header.Set("X-Go-Go-Host-User", user)
	revokeKeyReq.Header.Set("Content-Type", "application/json")
	revokeKeyRec := httptest.NewRecorder()
	h.ServeHTTP(revokeKeyRec, revokeKeyReq)
	if revokeKeyRec.Code != http.StatusOK {
		t.Fatalf("revoke key: %d %s", revokeKeyRec.Code, revokeKeyRec.Body.String())
	}
	revokedKeyReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-revoked-key")
	revokedKeyRec := httptest.NewRecorder()
	h.ServeHTTP(revokedKeyRec, revokedKeyReq)
	if revokedKeyRec.Code != http.StatusForbidden {
		t.Fatalf("expected revoked key denial, got %d %s", revokedKeyRec.Code, revokedKeyRec.Body.String())
	}
	rotateReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+org.ID+"/agents/"+created.Agent.ID+"/enrollment-token", nil)
	rotateReq.Header.Set("X-Go-Go-Host-User", user)
	rotateRec := httptest.NewRecorder()
	h.ServeHTTP(rotateRec, rotateReq)
	if rotateRec.Code != http.StatusCreated {
		t.Fatalf("rotation token: %d %s", rotateRec.Code, rotateRec.Body.String())
	}
	var rotation createAgentEnrollmentTokenResponse
	if err := json.Unmarshal(rotateRec.Body.Bytes(), &rotation); err != nil {
		t.Fatal(err)
	}
	newPub, newPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	newEnrollBody := []byte(`{"token":"` + rotation.EnrollmentToken + `","publicKey":"` + base64.StdEncoding.EncodeToString(newPub) + `"}`)
	newEnrollReq := httptest.NewRequest(http.MethodPost, "/api/v1/agent/enroll", bytes.NewReader(newEnrollBody))
	newEnrollReq.Header.Set("Content-Type", "application/json")
	newEnrollRec := httptest.NewRecorder()
	h.ServeHTTP(newEnrollRec, newEnrollReq)
	if newEnrollRec.Code != http.StatusCreated {
		t.Fatalf("replacement enroll: %d %s", newEnrollRec.Code, newEnrollRec.Body.String())
	}
	var replacement enrollAgentResponse
	if err := json.Unmarshal(newEnrollRec.Body.Bytes(), &replacement); err != nil {
		t.Fatal(err)
	}
	replacementReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, replacement.KeyID, newPriv, time.Now().UTC(), "nonce-replacement")
	replacementRec := httptest.NewRecorder()
	h.ServeHTTP(replacementRec, replacementReq)
	if replacementRec.Code != http.StatusCreated {
		t.Fatalf("expected replacement key success, got %d %s", replacementRec.Code, replacementRec.Body.String())
	}

	badSigReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, replacement.KeyID, newPriv, time.Now().UTC(), "nonce-bad-sig")
	badSigReq.Header.Set("X-Go-Go-Agent-Signature", base64.StdEncoding.EncodeToString([]byte("not a real signature")))
	badSigRec := httptest.NewRecorder()
	h.ServeHTTP(badSigRec, badSigReq)
	if badSigRec.Code != http.StatusForbidden {
		t.Fatalf("expected bad signature denial, got %d %s", badSigRec.Code, badSigRec.Body.String())
	}

	restrictGrantBody := []byte(`{"siteId":"` + site.ID + `","canDeploy":true,"canActivate":false,"allowedChannels":["default"],"allowedBundlePaths":["go-go-host.json"]}`)
	restrictGrantReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+org.ID+"/agents/"+created.Agent.ID+"/grants", bytes.NewReader(restrictGrantBody))
	restrictGrantReq.Header.Set("X-Go-Go-Host-User", user)
	restrictGrantReq.Header.Set("Content-Type", "application/json")
	restrictGrantRec := httptest.NewRecorder()
	h.ServeHTTP(restrictGrantRec, restrictGrantReq)
	if restrictGrantRec.Code != http.StatusOK {
		t.Fatalf("restrict grant: %d %s", restrictGrantRec.Code, restrictGrantRec.Body.String())
	}

	actionActivateBody := []byte(`{"siteId":"` + site.ID + `","channel":"default","bundlePath":"go-go-host.json","action":"activate"}`)
	actionActivateReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", actionActivateBody, created.Agent.ID, replacement.KeyID, newPriv, time.Now().UTC(), "nonce-action-activate")
	actionActivateRec := httptest.NewRecorder()
	h.ServeHTTP(actionActivateRec, actionActivateReq)
	if actionActivateRec.Code != http.StatusForbidden {
		t.Fatalf("expected action-based activation denial, got %d %s", actionActivateRec.Code, actionActivateRec.Body.String())
	}

	restrictedRunBody := []byte(`{"siteId":"` + site.ID + `","channel":"default","bundlePath":"go-go-host.json","action":"deploy"}`)
	restrictedRunReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", restrictedRunBody, created.Agent.ID, replacement.KeyID, newPriv, time.Now().UTC(), "nonce-restricted-run")
	restrictedRunRec := httptest.NewRecorder()
	h.ServeHTTP(restrictedRunRec, restrictedRunReq)
	if restrictedRunRec.Code != http.StatusCreated {
		t.Fatalf("restricted signed run: %d %s", restrictedRunRec.Code, restrictedRunRec.Body.String())
	}
	var restrictedRun createDeployRunResponse
	if err := json.Unmarshal(restrictedRunRec.Body.Bytes(), &restrictedRun); err != nil {
		t.Fatal(err)
	}
	restrictedUploadRec := uploadAgentBundleViaAPI(t, h, restrictedRun.ID, restrictedRun.UploadToken, writeHelloBundle(t))
	if restrictedUploadRec.Code != http.StatusBadRequest {
		t.Fatalf("expected restricted upload rejection, got %d %s", restrictedUploadRec.Code, restrictedUploadRec.Body.String())
	}
}

type agentUploadResponse struct {
	Activated  bool          `json:"activated"`
	Deployment deploymentDTO `json:"deployment"`
}

func uploadAgentBundleViaAPI(t *testing.T, h http.Handler, runID, uploadToken, bundlePath string) *httptest.ResponseRecorder {
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
		_ = f.Close()
		t.Fatal(err)
	}
	_ = f.Close()
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/deploy-runs/"+runID+"/upload", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("X-Go-Go-Upload-Token", uploadToken)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func signedAgentRequest(t *testing.T, method, path string, body []byte, agentID, keyID string, priv ed25519.PrivateKey, ts time.Time, nonce string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	timestamp := ts.Format(time.RFC3339)
	canonical := control.AgentCanonicalString(method, path, control.HashBody(body), timestamp, nonce)
	sig := ed25519.Sign(priv, []byte(canonical))
	req.Header.Set("X-Go-Go-Agent-ID", agentID)
	req.Header.Set("X-Go-Go-Agent-Key-ID", keyID)
	req.Header.Set("X-Go-Go-Agent-Timestamp", timestamp)
	req.Header.Set("X-Go-Go-Agent-Nonce", nonce)
	req.Header.Set("X-Go-Go-Agent-Signature", base64.StdEncoding.EncodeToString(sig))
	return req
}
