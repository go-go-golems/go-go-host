package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/control"
)

type meResponse struct {
	User          userDTO         `json:"user"`
	Memberships   []membershipDTO `json:"memberships"`
	PlatformAdmin bool            `json:"platformAdmin"`
}

type userDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type membershipDTO struct {
	OrgID   string `json:"orgId"`
	OrgSlug string `json:"orgSlug"`
	OrgName string `json:"orgName"`
	Role    string `json:"role"`
}

type createOrgRequest struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type createSiteRequest struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func handleMe(core *control.Core) http.HandlerFunc {
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
		admin, err := core.Store.IsPlatformAdmin(r.Context(), p.User.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := meResponse{User: userDTO{ID: p.User.ID, Email: p.User.Email, DisplayName: p.User.DisplayName}, Memberships: []membershipDTO{}, PlatformAdmin: admin}
		for _, m := range memberships {
			out.Memberships = append(out.Memberships, membershipDTO{OrgID: m.OrgID, OrgSlug: m.OrgSlug, OrgName: m.OrgName, Role: m.Role})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleCreateOrg(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req createOrgRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		org, err := core.Orgs.CreateOrg(r.Context(), p.User.ID, req.Slug, req.Name)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, orgToDTO(org))
	}
}

func handleListSites(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		orgID := r.PathValue("org_id")
		sites, err := core.Sites.ListSites(r.Context(), p.User.ID, orgID)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, control.ErrPermissionDenied) {
				status = http.StatusForbidden
			}
			writeError(w, status, err.Error())
			return
		}
		out := make([]siteDTO, 0, len(sites))
		for _, site := range sites {
			out = append(out, siteToDTO(site))
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleCreateSite(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req createSiteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		site, err := core.Sites.CreateSite(r.Context(), p.User.ID, r.PathValue("org_id"), req.Slug, req.Name)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, control.ErrPermissionDenied) {
				status = http.StatusForbidden
			}
			writeError(w, status, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, siteToDTO(*site))
	}
}
