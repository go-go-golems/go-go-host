package httpapi

import (
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/control"
)

func handleDBStats(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		siteID := r.PathValue("site_id")
		site, err := core.Store.GetSite(r.Context(), siteID)
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		if _, err := core.Store.MembershipRole(r.Context(), site.OrgID, p.User.ID); err != nil {
			writeDeploymentError(w, control.ErrPermissionDenied)
			return
		}
		stats, ok, err := core.Supervisor.DBStats(siteID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "runtime not found")
			return
		}
		writeJSON(w, http.StatusOK, stats)
	}
}
