package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestAgentsAndAuditAPIFlow(t *testing.T) {
	h := newIntegrationHandler(t)
	suffix := uuid.NewString()[:8]
	user := "agent-" + suffix
	org := createTestOrgViaAPI(t, h, user, "agent-org-"+suffix)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/orgs/"+org.ID+"/agents", bytes.NewReader([]byte(`{"name":"ci-bot"}`)))
	createReq.Header.Set("X-Go-Go-Host-User", user)
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	h.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create agent: %d %s", createRec.Code, createRec.Body.String())
	}
	var agent agentDTO
	if err := json.Unmarshal(createRec.Body.Bytes(), &agent); err != nil {
		t.Fatal(err)
	}
	if agent.Status != "active" {
		t.Fatalf("expected active agent, got %#v", agent)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/orgs/"+org.ID+"/agents", nil)
	listReq.Header.Set("X-Go-Go-Host-User", user)
	listRec := httptest.NewRecorder()
	h.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list agents: %d %s", listRec.Code, listRec.Body.String())
	}
	var agents []agentDTO
	if err := json.Unmarshal(listRec.Body.Bytes(), &agents); err != nil {
		t.Fatal(err)
	}
	if len(agents) != 1 {
		t.Fatalf("expected one agent, got %#v", agents)
	}

	auditReq := httptest.NewRequest(http.MethodGet, "/api/v1/orgs/"+org.ID+"/audit?action=agent.create", nil)
	auditReq.Header.Set("X-Go-Go-Host-User", user)
	auditRec := httptest.NewRecorder()
	h.ServeHTTP(auditRec, auditReq)
	if auditRec.Code != http.StatusOK {
		t.Fatalf("audit: %d %s", auditRec.Code, auditRec.Body.String())
	}
	var events []auditDTO
	if err := json.Unmarshal(auditRec.Body.Bytes(), &events); err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 || events[0].Action != "agent.create" {
		t.Fatalf("expected agent.create event, got %#v", events)
	}
}
