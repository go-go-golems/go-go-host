package httpapi

import (
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/control"
)

func handleRuntimeStatus(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		siteID := r.PathValue("site_id")
		site, err := core.Store.GetSite(r.Context(), siteID)
		if err != nil {
			writeError(w, http.StatusNotFound, "site not found")
			return
		}
		if err := core.Orgs.EnsureRole(r.Context(), p.User.ID, site.OrgID, "org_owner", "org_developer", "org_viewer"); err != nil {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		status, ok := core.Supervisor.Status(siteID)
		if !ok {
			writeJSON(w, http.StatusOK, map[string]any{"siteId": siteID, "status": "stopped"})
			return
		}
		writeJSON(w, http.StatusOK, status)
	}
}

func handleAdminRuntimeSummary(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		admin, err := core.Store.IsPlatformAdmin(r.Context(), p.User.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !admin {
			writeError(w, http.StatusForbidden, "platform admin required")
			return
		}
		writeJSON(w, http.StatusOK, core.Supervisor.Summary())
	}
}
