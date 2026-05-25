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

func TestMaintenanceMetadataPruneAndReadyz(t *testing.T) {
	h := newIntegrationHandler(t)
	suffix := uuid.NewString()[:8]
	user := "maint-" + suffix
	org := createTestOrgViaAPI(t, h, user, "maint-org-"+suffix)
	site := createTestSiteViaAPI(t, h, user, org.ID, "maint-site-"+suffix)

	readyReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	readyRec := httptest.NewRecorder()
	h.ServeHTTP(readyRec, readyReq)
	if readyRec.Code != http.StatusOK || !strings.Contains(readyRec.Body.String(), `"db":"ok"`) {
		t.Fatalf("readyz: %d %s", readyRec.Code, readyRec.Body.String())
	}

	metaReq := httptest.NewRequest(http.MethodGet, "/api/v1/sites/"+site.ID+"/export/metadata", nil)
	metaReq.Header.Set("X-Go-Go-Host-User", user)
	metaRec := httptest.NewRecorder()
	h.ServeHTTP(metaRec, metaReq)
	if metaRec.Code != http.StatusOK {
		t.Fatalf("metadata export: %d %s", metaRec.Code, metaRec.Body.String())
	}
	var metadata struct {
		Site struct {
			ID string `json:"id"`
		} `json:"site"`
		Capabilities []any  `json:"capabilities"`
		ExportedAt   string `json:"exportedAt"`
	}
	if err := json.Unmarshal(metaRec.Body.Bytes(), &metadata); err != nil {
		t.Fatalf("decode metadata: %v", err)
	}
	if metadata.Site.ID != site.ID || len(metadata.Capabilities) == 0 || metadata.ExportedAt == "" {
		t.Fatalf("unexpected metadata: %+v", metadata)
	}

	pruneReq := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+site.ID+"/deployments/prune", bytes.NewReader([]byte(`{"dryRun":true,"olderThan":"0d"}`)))
	pruneReq.Header.Set("X-Go-Go-Host-User", user)
	pruneReq.Header.Set("Content-Type", "application/json")
	pruneRec := httptest.NewRecorder()
	h.ServeHTTP(pruneRec, pruneReq)
	if pruneRec.Code != http.StatusOK || !strings.Contains(pruneRec.Body.String(), `"dryRun":true`) {
		t.Fatalf("prune dry-run: %d %s", pruneRec.Code, pruneRec.Body.String())
	}
}
