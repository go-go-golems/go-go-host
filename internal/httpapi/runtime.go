package httpapi

import (
	"errors"
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/runtime"
	"github.com/go-go-golems/go-go-host/internal/store"
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

func handleAdminRuntimeRestart(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := requirePlatformAdmin(core, w, r)
		if !ok {
			return
		}
		siteID := r.PathValue("site_id")
		if err := core.Supervisor.Restart(r.Context(), siteID); err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, runtime.ErrRuntimeNotFound) {
				status = http.StatusNotFound
			}
			writeError(w, status, err.Error())
			return
		}
		status, _ := core.Supervisor.Status(siteID)
		_, _ = core.Store.InsertAuditEvent(r.Context(), store.AuditEvent{OrgID: status.OrgID, ActorType: "user", ActorID: p.User.ID, Action: "runtime.restart", ResourceType: "site", ResourceID: siteID})
		writeJSON(w, http.StatusOK, status)
	}
}

func handleAdminRuntimeStop(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := requirePlatformAdmin(core, w, r)
		if !ok {
			return
		}
		siteID := r.PathValue("site_id")
		previous, _ := core.Supervisor.Status(siteID)
		if err := core.Supervisor.Stop(r.Context(), siteID); err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, runtime.ErrRuntimeNotFound) {
				status = http.StatusNotFound
			}
			writeError(w, status, err.Error())
			return
		}
		status, _ := core.Supervisor.Status(siteID)
		orgID := status.OrgID
		if orgID == "" {
			orgID = previous.OrgID
		}
		_, _ = core.Store.InsertAuditEvent(r.Context(), store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: p.User.ID, Action: "runtime.stop", ResourceType: "site", ResourceID: siteID})
		writeJSON(w, http.StatusOK, status)
	}
}

func handleAdminRuntimeSummary(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		writeJSON(w, http.StatusOK, core.Supervisor.Summary())
	}
}
