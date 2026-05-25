package httpapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/httpapi/docfs"
	"github.com/go-go-golems/go-go-host/internal/webadmin"
)

const Version = "0.1.0-dev"

func NewHandler(core *control.Core) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /readyz", handleReady(core))
	mux.HandleFunc("GET /api/v1/version", handleVersion)
	dashboard := webadmin.NewHandler()
	mux.Handle("/app", http.StripPrefix("/app", dashboard))
	mux.Handle("/app/", http.StripPrefix("/app", dashboard))
	mux.Handle("/admin", http.StripPrefix("/admin", webadmin.NewPlaceholderHandler()))
	mux.Handle("/admin/", http.StripPrefix("/admin", webadmin.NewPlaceholderHandler()))
	mux.HandleFunc("GET /api/v1/config", func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]any{
			"baseDomain":    core.Config.BaseDomain,
			"publicBaseUrl": core.Config.PublicBaseURL,
			"devAuth":       core.Config.DevAuth,
		}
		if !core.Config.DevAuth && core.Config.OIDCIssuer != "" && core.Config.OIDCClientID != "" {
			response["oidc"] = map[string]any{
				"issuer":             core.Config.OIDCIssuer,
				"clientId":           core.Config.OIDCClientID,
				"deviceClientId":     core.Config.OIDCDeviceClientID,
				"scopes":             core.Config.OIDCScopes,
				"redirectPath":       core.Config.OIDCRedirectPath,
				"logoutRedirectPath": core.Config.OIDCLogoutRedirectPath,
			}
		}
		writeJSON(w, http.StatusOK, response)
	})

	api := http.NewServeMux()
	api.HandleFunc("GET /api/v1/me", handleMe(core))
	api.HandleFunc("GET /api/v1/orgs", handleListOrgs(core))
	api.HandleFunc("POST /api/v1/orgs", handleCreateOrg(core))
	api.HandleFunc("GET /api/v1/orgs/{org_id}/sites", handleListSites(core))
	api.HandleFunc("POST /api/v1/orgs/{org_id}/sites", handleCreateSite(core))
	api.HandleFunc("GET /api/v1/orgs/{org_id}/agents", handleListAgents(core))
	api.HandleFunc("POST /api/v1/orgs/{org_id}/agents", handleCreateAgent(core))
	api.HandleFunc("POST /api/v1/orgs/{org_id}/agents/{agent_id}/revoke", handleRevokeAgent(core))
	api.HandleFunc("POST /api/v1/orgs/{org_id}/agents/{agent_id}/enrollment-token", handleCreateAgentEnrollmentToken(core))
	api.HandleFunc("GET /api/v1/orgs/{org_id}/agents/{agent_id}/keys", handleListAgentKeys(core))
	api.HandleFunc("POST /api/v1/orgs/{org_id}/agents/{agent_id}/keys/{key_id}/revoke", handleRevokeAgentKey(core))
	api.HandleFunc("POST /api/v1/orgs/{org_id}/agents/{agent_id}/grants", handleUpsertAgentGrant(core))
	api.HandleFunc("GET /api/v1/orgs/{org_id}/audit", handleListAudit(core))
	api.HandleFunc("POST /api/v1/agent/enroll", handleEnrollAgent(core))
	api.HandleFunc("POST /api/v1/agent/deploy-runs", handleCreateAgentDeployRun(core))
	api.HandleFunc("POST /api/v1/agent/deploy-runs/{run_id}/upload", handleAgentDeployRunUpload(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/runtime", handleRuntimeStatus(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/db/stats", handleDBStats(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/config", handleListSiteConfig(core))
	api.HandleFunc("PUT /api/v1/sites/{site_id}/config", handleUpsertSiteConfig(core))
	api.HandleFunc("DELETE /api/v1/sites/{site_id}/config", handleDeleteSiteConfig(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/capabilities", handleListSiteCapabilities(core))
	api.HandleFunc("PUT /api/v1/sites/{site_id}/capabilities", handleUpsertSiteCapability(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/domains", handleListSiteDomains(core))
	api.HandleFunc("POST /api/v1/sites/{site_id}/domains", handleAddSiteDomain(core))
	api.HandleFunc("POST /api/v1/sites/{site_id}/domains/{domain_id}/verify", handleVerifySiteDomain(core))
	api.HandleFunc("DELETE /api/v1/sites/{site_id}/domains/{domain_id}", handleDeleteSiteDomain(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/environment", handleSiteEnvironmentPlaceholder(core))
	api.HandleFunc("POST /api/v1/sites/{site_id}/deployments", handleUploadDeployment(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/deployments", handleListDeployments(core))
	api.HandleFunc("POST /api/v1/sites/{site_id}/rollback", handleRollbackDeployment(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/export/metadata", handleExportSiteMetadata(core))
	api.HandleFunc("GET /api/v1/sites/{site_id}/export/db", handleExportSiteDB(core))
	api.HandleFunc("POST /api/v1/sites/{site_id}/deployments/prune", handlePruneDeployments(core))
	api.HandleFunc("GET /api/v1/deployments/{deployment_id}", handleGetDeployment(core))
	api.HandleFunc("GET /api/v1/deployments/{deployment_id}/bundle", handleExportDeploymentBundle(core))
	api.HandleFunc("POST /api/v1/deployments/{deployment_id}/activate", handleActivateDeployment(core))
	api.HandleFunc("GET /api/v1/admin/runtimes/summary", handleAdminRuntimeSummary(core))
	api.HandleFunc("POST /api/v1/admin/runtimes/{site_id}/restart", handleAdminRuntimeRestart(core))
	api.HandleFunc("POST /api/v1/admin/runtimes/{site_id}/stop", handleAdminRuntimeStop(core))
	api.HandleFunc("GET /api/v1/admin/orgs", handleAdminListOrgs(core))
	api.HandleFunc("GET /api/v1/admin/users", handleAdminListUsers(core))
	api.HandleFunc("GET /api/v1/admin/sites", handleAdminListSites(core))
	api.HandleFunc("GET /api/v1/admin/deployments", handleAdminListDeployments(core))
	api.HandleFunc("GET /api/v1/admin/deployments/{deployment_id}", handleAdminGetDeployment(core))
	api.HandleFunc("GET /api/v1/admin/agents", handleAdminListAgents(core))
	api.HandleFunc("GET /api/v1/admin/audit", handleAdminListAudit(core))
	api.HandleFunc("GET /api/v1/admin/quotas", handleAdminListQuotas(core))
	api.HandleFunc("GET /api/v1/admin/capabilities", handleAdminListCapabilities(core))
	api.HandleFunc("GET /api/v1/admin/domains", handleAdminListDomains(core))
	api.HandleFunc("POST /api/v1/admin/audit/retention", handleAuditRetention(core))
	api.HandleFunc("GET /api/v1/docs", docfs.HandleListDocs)
	api.HandleFunc("GET /api/v1/docs/{slug}", docfs.HandleGetDoc)

	authn := &oidcAuthenticator{cfg: core.Config, st: core.Store}
	authedAPI := authMiddleware(api, authn, core.Config.DevAuth)
	mux.Handle("/api/v1/me", authedAPI)
	mux.Handle("/api/v1/orgs", authedAPI)
	mux.Handle("/api/v1/orgs/", authedAPI)
	mux.Handle("/api/v1/sites/", authedAPI)
	mux.Handle("/api/v1/deployments/", authedAPI)
	mux.Handle("/api/v1/admin/", authedAPI)
	mux.Handle("/api/v1/docs", authedAPI)
	mux.Handle("/api/v1/docs/", authedAPI)
	mux.Handle("/api/v1/agent/", api)

	return withRequestID(rootRedirectHandler{next: withFallback(mux, core.Supervisor), dashboardHost: hostFromPublicBaseURL(core.Config.PublicBaseURL)})
}

type rootRedirectHandler struct {
	next          http.Handler
	dashboardHost string
}

func (h rootRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && h.isDashboardHost(r.Host) {
		http.Redirect(w, r, "/app", http.StatusFound)
		return
	}
	h.next.ServeHTTP(w, r)
}

func (h rootRedirectHandler) isDashboardHost(host string) bool {
	if h.dashboardHost == "" {
		return true
	}
	return strings.EqualFold(stripPort(host), h.dashboardHost)
}

func hostFromPublicBaseURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return stripPort(u.Host)
}

func stripPort(host string) string {
	if strings.HasPrefix(host, "[") {
		if end := strings.LastIndex(host, "]"); end >= 0 {
			return strings.Trim(host[:end+1], "[]")
		}
	}
	if idx := strings.LastIndex(host, ":"); idx >= 0 {
		return host[:idx]
	}
	return host
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleReady(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checks := map[string]string{"db": "ok", "dataDir": "ok"}
		status := http.StatusOK
		if core.Store == nil || core.Store.Ping(r.Context()) != nil {
			checks["db"] = "failed"
			status = http.StatusServiceUnavailable
		}
		if err := os.MkdirAll(core.Config.DataDir, 0o755); err != nil {
			checks["dataDir"] = "failed"
			status = http.StatusServiceUnavailable
		} else if f, err := os.CreateTemp(core.Config.DataDir, ".readyz-*"); err != nil {
			checks["dataDir"] = "failed"
			status = http.StatusServiceUnavailable
		} else {
			name := f.Name()
			_ = f.Close()
			_ = os.Remove(name)
		}
		state := "ready"
		if status != http.StatusOK {
			state = "not_ready"
		}
		writeJSON(w, status, map[string]any{"status": state, "checks": checks})
	}
}

func handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": Version})
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = time.Now().UTC().Format("20060102T150405.000000000Z")
		}
		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": message})
}
