package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestSiteSettingsDomainsCapabilitiesAndAudit(t *testing.T) {
	h := newIntegrationHandler(t)
	suffix := uuid.NewString()[:8]
	user := "phase11-" + suffix
	org := createTestOrgViaAPI(t, h, user, "phase11-org-"+suffix)
	site := createTestSiteViaAPI(t, h, user, org.ID, "phase11-site-"+suffix)

	putConfig := httptest.NewRequest(http.MethodPut, "/api/v1/sites/"+site.ID+"/config", strings.NewReader(`{"key":"theme.title","value":{"text":"Hello"}}`))
	putConfig.Header.Set("X-Go-Go-Host-User", user)
	putConfig.Header.Set("Content-Type", "application/json")
	putConfigRec := httptest.NewRecorder()
	h.ServeHTTP(putConfigRec, putConfig)
	if putConfigRec.Code != http.StatusOK {
		t.Fatalf("put config: %d %s", putConfigRec.Code, putConfigRec.Body.String())
	}

	listConfig := httptest.NewRequest(http.MethodGet, "/api/v1/sites/"+site.ID+"/config", nil)
	listConfig.Header.Set("X-Go-Go-Host-User", user)
	listConfigRec := httptest.NewRecorder()
	h.ServeHTTP(listConfigRec, listConfig)
	if listConfigRec.Code != http.StatusOK {
		t.Fatalf("list config: %d %s", listConfigRec.Code, listConfigRec.Body.String())
	}
	var configItems []siteConfigDTO
	if err := json.Unmarshal(listConfigRec.Body.Bytes(), &configItems); err != nil {
		t.Fatalf("decode config: %v", err)
	}
	if len(configItems) != 1 || configItems[0].Key != "theme.title" || !bytes.Contains(configItems[0].Value, []byte("Hello")) {
		t.Fatalf("unexpected config items: %+v", configItems)
	}

	putCap := httptest.NewRequest(http.MethodPut, "/api/v1/sites/"+site.ID+"/capabilities", strings.NewReader(`{"capability":"assets","enabled":false,"config":{}}`))
	putCap.Header.Set("X-Go-Go-Host-User", user)
	putCap.Header.Set("Content-Type", "application/json")
	putCapRec := httptest.NewRecorder()
	h.ServeHTTP(putCapRec, putCap)
	if putCapRec.Code != http.StatusOK {
		t.Fatalf("put capability: %d %s", putCapRec.Code, putCapRec.Body.String())
	}
	listCap := httptest.NewRequest(http.MethodGet, "/api/v1/sites/"+site.ID+"/capabilities", nil)
	listCap.Header.Set("X-Go-Go-Host-User", user)
	listCapRec := httptest.NewRecorder()
	h.ServeHTTP(listCapRec, listCap)
	if listCapRec.Code != http.StatusOK {
		t.Fatalf("list capabilities: %d %s", listCapRec.Code, listCapRec.Body.String())
	}
	var caps []siteCapabilityDTO
	if err := json.Unmarshal(listCapRec.Body.Bytes(), &caps); err != nil {
		t.Fatalf("decode caps: %v", err)
	}
	foundAssetsDisabled := false
	for _, cap := range caps {
		if cap.Capability == "assets" && !cap.Enabled {
			foundAssetsDisabled = true
		}
	}
	if !foundAssetsDisabled {
		t.Fatalf("assets capability was not disabled: %+v", caps)
	}

	addDomain := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+site.ID+"/domains", strings.NewReader(`{"hostname":"www.`+suffix+`.example.test"}`))
	addDomain.Header.Set("X-Go-Go-Host-User", user)
	addDomain.Header.Set("Content-Type", "application/json")
	addDomainRec := httptest.NewRecorder()
	h.ServeHTTP(addDomainRec, addDomain)
	if addDomainRec.Code != http.StatusCreated {
		t.Fatalf("add domain: %d %s", addDomainRec.Code, addDomainRec.Body.String())
	}
	var domain siteDomainDTO
	if err := json.Unmarshal(addDomainRec.Body.Bytes(), &domain); err != nil {
		t.Fatalf("decode domain: %v", err)
	}
	if domain.Status != "pending" || domain.VerificationToken == "" {
		t.Fatalf("unexpected new domain: %+v", domain)
	}

	verifyDomain := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+site.ID+"/domains/"+domain.ID+"/verify", nil)
	verifyDomain.Header.Set("X-Go-Go-Host-User", user)
	verifyDomainRec := httptest.NewRecorder()
	h.ServeHTTP(verifyDomainRec, verifyDomain)
	if verifyDomainRec.Code != http.StatusOK {
		t.Fatalf("verify domain: %d %s", verifyDomainRec.Code, verifyDomainRec.Body.String())
	}
	var verified siteDomainDTO
	if err := json.Unmarshal(verifyDomainRec.Body.Bytes(), &verified); err != nil {
		t.Fatalf("decode verified domain: %v", err)
	}
	if verified.Status != "verified" || verified.VerifiedAt == "" {
		t.Fatalf("unexpected verified domain: %+v", verified)
	}

	envReq := httptest.NewRequest(http.MethodGet, "/api/v1/sites/"+site.ID+"/environment", nil)
	envReq.Header.Set("X-Go-Go-Host-User", user)
	envRec := httptest.NewRecorder()
	h.ServeHTTP(envRec, envReq)
	if envRec.Code != http.StatusOK || !strings.Contains(envRec.Body.String(), "Secrets/environment variables are intentionally not implemented") {
		t.Fatalf("environment placeholder: %d %s", envRec.Code, envRec.Body.String())
	}

	auditReq := httptest.NewRequest(http.MethodGet, "/api/v1/orgs/"+org.ID+"/audit", nil)
	auditReq.Header.Set("X-Go-Go-Host-User", user)
	auditRec := httptest.NewRecorder()
	h.ServeHTTP(auditRec, auditReq)
	if auditRec.Code != http.StatusOK {
		t.Fatalf("audit: %d %s", auditRec.Code, auditRec.Body.String())
	}
	body := auditRec.Body.String()
	for _, action := range []string{"site.config.upsert", "site.capability.update", "site.domain.add", "site.domain.verify"} {
		if !strings.Contains(body, action) {
			t.Fatalf("expected audit action %s in %s", action, body)
		}
	}
}
