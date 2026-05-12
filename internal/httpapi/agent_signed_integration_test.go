package httpapi

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	createBody := []byte(`{"name":"ci","siteId":"` + site.ID + `","allowedChannels":["default"],"allowedPaths":["bundles/**"]}`)
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

	goodBody := []byte(`{"siteId":"` + site.ID + `","channel":"default","path":"bundles/app.tar.gz","action":"deploy"}`)
	goodReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-good")
	goodRec := httptest.NewRecorder()
	h.ServeHTTP(goodRec, goodReq)
	if goodRec.Code != http.StatusCreated {
		t.Fatalf("good signed run: %d %s", goodRec.Code, goodRec.Body.String())
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

	wrongPathBody := []byte(`{"siteId":"` + site.ID + `","channel":"default","path":"private/app.tar.gz","action":"deploy"}`)
	wrongPathReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", wrongPathBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-wrong-path")
	wrongPathRec := httptest.NewRecorder()
	h.ServeHTTP(wrongPathRec, wrongPathReq)
	if wrongPathRec.Code != http.StatusForbidden {
		t.Fatalf("expected wrong path denial, got %d %s", wrongPathRec.Code, wrongPathRec.Body.String())
	}

	wrongSiteBody := []byte(`{"siteId":"` + otherSite.ID + `","channel":"default","path":"bundles/app.tar.gz","action":"deploy"}`)
	wrongSiteReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", wrongSiteBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-wrong-site")
	wrongSiteRec := httptest.NewRecorder()
	h.ServeHTTP(wrongSiteRec, wrongSiteReq)
	if wrongSiteRec.Code != http.StatusForbidden {
		t.Fatalf("expected wrong site denial, got %d %s", wrongSiteRec.Code, wrongSiteRec.Body.String())
	}

	badSigReq := signedAgentRequest(t, http.MethodPost, "/api/v1/agent/deploy-runs", goodBody, created.Agent.ID, enrolled.KeyID, priv, time.Now().UTC(), "nonce-bad-sig")
	badSigReq.Header.Set("X-Go-Go-Agent-Signature", base64.StdEncoding.EncodeToString([]byte("not a real signature")))
	badSigRec := httptest.NewRecorder()
	h.ServeHTTP(badSigRec, badSigReq)
	if badSigRec.Code != http.StatusForbidden {
		t.Fatalf("expected bad signature denial, got %d %s", badSigRec.Code, badSigRec.Body.String())
	}
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
