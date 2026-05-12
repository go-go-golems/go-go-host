package httpapi

import (
	"net/http"
	"strconv"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type adminOrgDTO struct {
	ID              string `json:"id"`
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	CreatedAt       string `json:"createdAt"`
	MemberCount     int64  `json:"memberCount"`
	SiteCount       int64  `json:"siteCount"`
	DeploymentCount int64  `json:"deploymentCount"`
}

type adminUserDTO struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	DisplayName   string `json:"displayName"`
	CreatedAt     string `json:"createdAt"`
	LastLoginAt   string `json:"lastLoginAt,omitempty"`
	PlatformAdmin bool   `json:"platformAdmin"`
	OrgCount      int64  `json:"orgCount"`
}

type adminSiteDTO struct {
	ID                 string `json:"id"`
	OrgID              string `json:"orgId"`
	OrgSlug            string `json:"orgSlug"`
	OrgName            string `json:"orgName"`
	Slug               string `json:"slug"`
	Name               string `json:"name"`
	PrimaryHost        string `json:"primaryHost"`
	Status             string `json:"status"`
	ActiveDeploymentID string `json:"activeDeploymentId"`
	CreatedAt          string `json:"createdAt"`
	RuntimeStatus      string `json:"runtimeStatus"`
	RequestsTotal      int64  `json:"requestsTotal"`
	ErrorsTotal        int64  `json:"errorsTotal"`
	LastError          string `json:"lastError,omitempty"`
}

type adminDeploymentDTO struct {
	ID             string `json:"id"`
	SiteID         string `json:"siteId"`
	SiteSlug       string `json:"siteSlug"`
	PrimaryHost    string `json:"primaryHost"`
	OrgID          string `json:"orgId"`
	OrgSlug        string `json:"orgSlug"`
	OrgName        string `json:"orgName"`
	Version        int    `json:"version"`
	Status         string `json:"status"`
	BundleRef      string `json:"bundleRef"`
	UnpackedPath   string `json:"unpackedPath"`
	ManifestJSON   string `json:"manifestJson"`
	ValidationJSON string `json:"validationJson"`
	CreatedByType  string `json:"createdByType"`
	CreatedByID    string `json:"createdById"`
	CreatedAt      string `json:"createdAt"`
	ActivatedAt    string `json:"activatedAt,omitempty"`
}

func requirePlatformAdmin(core *control.Core, w http.ResponseWriter, r *http.Request) (principal, bool) {
	p, err := requirePrincipal(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return principal{}, false
	}
	admin, err := core.Store.IsPlatformAdmin(r.Context(), p.User.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return principal{}, false
	}
	if !admin {
		writeError(w, http.StatusForbidden, "platform admin required")
		return principal{}, false
	}
	return p, true
}

func handleAdminListOrgs(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		orgs, err := core.Store.ListAdminOrgs(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminOrgDTO, 0, len(orgs))
		for _, org := range orgs {
			out = append(out, adminOrgDTO{ID: org.ID, Slug: org.Slug, Name: org.Name, CreatedAt: org.CreatedAt, MemberCount: org.MemberCount, SiteCount: org.SiteCount, DeploymentCount: org.DeploymentCount})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListUsers(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		users, err := core.Store.ListAdminUsers(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminUserDTO, 0, len(users))
		for _, user := range users {
			out = append(out, adminUserDTO{ID: user.ID, Email: user.Email, DisplayName: user.DisplayName, CreatedAt: user.CreatedAt, LastLoginAt: user.LastLoginAt, PlatformAdmin: user.PlatformAdmin, OrgCount: user.OrgCount})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListSites(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		sites, err := core.Store.ListAdminSites(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminSiteDTO, 0, len(sites))
		for _, site := range sites {
			out = append(out, adminSiteDTO{ID: site.ID, OrgID: site.OrgID, OrgSlug: site.OrgSlug, OrgName: site.OrgName, Slug: site.Slug, Name: site.Name, PrimaryHost: site.PrimaryHost, Status: site.Status, ActiveDeploymentID: site.ActiveDeploymentID, CreatedAt: site.CreatedAt, RuntimeStatus: site.RuntimeStatus, RequestsTotal: site.RequestsTotal, ErrorsTotal: site.ErrorsTotal, LastError: site.LastError})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListDeployments(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		deployments, err := core.Store.ListAdminDeployments(r.Context(), store.AdminDeploymentFilter{OrgID: r.URL.Query().Get("orgId"), SiteID: r.URL.Query().Get("siteId"), Status: r.URL.Query().Get("status"), Limit: limit})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminDeploymentDTO, 0, len(deployments))
		for _, deployment := range deployments {
			out = append(out, adminDeploymentDTO{ID: deployment.ID, SiteID: deployment.SiteID, SiteSlug: deployment.SiteSlug, PrimaryHost: deployment.PrimaryHost, OrgID: deployment.OrgID, OrgSlug: deployment.OrgSlug, OrgName: deployment.OrgName, Version: deployment.Version, Status: deployment.Status, BundleRef: deployment.BundleRef, UnpackedPath: deployment.UnpackedPath, ManifestJSON: string(deployment.ManifestJSON), ValidationJSON: string(deployment.ValidationJSON), CreatedByType: deployment.CreatedByType, CreatedByID: deployment.CreatedByID, CreatedAt: deployment.CreatedAt, ActivatedAt: deployment.ActivatedAt})
		}
		writeJSON(w, http.StatusOK, out)
	}
}
