package httpapi

import (
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/control"
)

func handleListOrgs(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		memberships, err := core.Store.ListMembershipsForUser(r.Context(), p.User.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]membershipDTO, 0, len(memberships))
		for _, m := range memberships {
			out = append(out, membershipDTO{OrgID: m.OrgID, OrgSlug: m.OrgSlug, OrgName: m.OrgName, Role: m.Role})
		}
		writeJSON(w, http.StatusOK, out)
	}
}
