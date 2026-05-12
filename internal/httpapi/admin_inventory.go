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

type adminAgentDTO struct {
	ID              string `json:"id"`
	OrgID           string `json:"orgId"`
	OrgSlug         string `json:"orgSlug"`
	OrgName         string `json:"orgName"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	CreatedByUserID string `json:"createdByUserId"`
	CreatedAt       string `json:"createdAt"`
	LastSeenAt      string `json:"lastSeenAt,omitempty"`
	GrantCount      int64  `json:"grantCount"`
}

type adminQuotaDTO struct {
	SiteID           string `json:"siteId"`
	SiteSlug         string `json:"siteSlug"`
	PrimaryHost      string `json:"primaryHost"`
	OrgID            string `json:"orgId"`
	OrgSlug          string `json:"orgSlug"`
	OrgName          string `json:"orgName"`
	BundleMaxBytes   int64  `json:"bundleMaxBytes"`
	DBSoftMaxBytes   int64  `json:"dbSoftMaxBytes"`
	DBHardMaxBytes   int64  `json:"dbHardMaxBytes"`
	RequestTimeoutMS int    `json:"requestTimeoutMs"`
	UpdatedAt        string `json:"updatedAt"`
	RequestsTotal    int64  `json:"requestsTotal"`
	ErrorsTotal      int64  `json:"errorsTotal"`
}

type adminCapabilityDTO struct {
	SiteID     string `json:"siteId"`
	SiteSlug   string `json:"siteSlug"`
	OrgID      string `json:"orgId"`
	OrgSlug    string `json:"orgSlug"`
	OrgName    string `json:"orgName"`
	Capability string `json:"capability"`
	Enabled    bool   `json:"enabled"`
	ConfigJSON string `json:"configJson"`
	UpdatedAt  string `json:"updatedAt"`
}

type adminDomainDTO struct {
	ID                string `json:"id"`
	SiteID            string `json:"siteId"`
	SiteSlug          string `json:"siteSlug"`
	OrgID             string `json:"orgId"`
	OrgSlug           string `json:"orgSlug"`
	OrgName           string `json:"orgName"`
	Hostname          string `json:"hostname"`
	Status            string `json:"status"`
	VerificationToken string `json:"verificationToken"`
	VerifiedAt        string `json:"verifiedAt,omitempty"`
	CreatedAt         string `json:"createdAt"`
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

func adminDeploymentToDTO(deployment *store.AdminDeployment) adminDeploymentDTO {
	return adminDeploymentDTO{ID: deployment.ID, SiteID: deployment.SiteID, SiteSlug: deployment.SiteSlug, PrimaryHost: deployment.PrimaryHost, OrgID: deployment.OrgID, OrgSlug: deployment.OrgSlug, OrgName: deployment.OrgName, Version: deployment.Version, Status: deployment.Status, BundleRef: deployment.BundleRef, UnpackedPath: deployment.UnpackedPath, ManifestJSON: string(deployment.ManifestJSON), ValidationJSON: string(deployment.ValidationJSON), CreatedByType: deployment.CreatedByType, CreatedByID: deployment.CreatedByID, CreatedAt: deployment.CreatedAt, ActivatedAt: deployment.ActivatedAt}
}

func handleAdminGetDeployment(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		deployment, err := core.Store.GetAdminDeployment(r.Context(), r.PathValue("deployment_id"))
		if err != nil {
			writeError(w, http.StatusNotFound, "deployment not found")
			return
		}
		writeJSON(w, http.StatusOK, adminDeploymentToDTO(deployment))
	}
}

func handleAdminListAgents(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		agents, err := core.Store.ListAdminAgents(r.Context(), store.AdminAgentFilter{OrgID: r.URL.Query().Get("orgId"), Status: r.URL.Query().Get("status")})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminAgentDTO, 0, len(agents))
		for _, agent := range agents {
			out = append(out, adminAgentDTO{ID: agent.ID, OrgID: agent.OrgID, OrgSlug: agent.OrgSlug, OrgName: agent.OrgName, Name: agent.Name, Status: agent.Status, CreatedByUserID: agent.CreatedByUserID, CreatedAt: agent.CreatedAt, LastSeenAt: agent.LastSeenAt, GrantCount: agent.GrantCount})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListAudit(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		events, err := core.Store.ListAdminAuditEvents(r.Context(), store.AuditFilter{OrgID: r.URL.Query().Get("orgId"), ResourceID: r.URL.Query().Get("resourceId"), ActorType: r.URL.Query().Get("actorType"), ActorID: r.URL.Query().Get("actorId"), Action: r.URL.Query().Get("action"), Limit: limit})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]auditDTO, 0, len(events))
		for _, event := range events {
			out = append(out, auditToDTO(event))
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListQuotas(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		rows, err := core.Store.ListAdminQuotas(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminQuotaDTO, 0, len(rows))
		for _, row := range rows {
			out = append(out, adminQuotaDTO{SiteID: row.SiteID, SiteSlug: row.SiteSlug, PrimaryHost: row.PrimaryHost, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, BundleMaxBytes: row.BundleMaxBytes, DBSoftMaxBytes: row.DBSoftMaxBytes, DBHardMaxBytes: row.DBHardMaxBytes, RequestTimeoutMS: row.RequestTimeoutMS, UpdatedAt: row.UpdatedAt, RequestsTotal: row.RequestsTotal, ErrorsTotal: row.ErrorsTotal})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListCapabilities(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		rows, err := core.Store.ListAdminCapabilities(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminCapabilityDTO, 0, len(rows))
		for _, row := range rows {
			out = append(out, adminCapabilityDTO{SiteID: row.SiteID, SiteSlug: row.SiteSlug, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Capability: row.Capability, Enabled: row.Enabled, ConfigJSON: string(row.ConfigJSON), UpdatedAt: row.UpdatedAt})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAdminListDomains(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePlatformAdmin(core, w, r); !ok {
			return
		}
		rows, err := core.Store.ListAdminDomains(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out := make([]adminDomainDTO, 0, len(rows))
		for _, row := range rows {
			out = append(out, adminDomainDTO{ID: row.ID, SiteID: row.SiteID, SiteSlug: row.SiteSlug, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Hostname: row.Hostname, Status: row.Status, VerificationToken: row.VerificationToken, VerifiedAt: row.VerifiedAt, CreatedAt: row.CreatedAt})
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
			out = append(out, adminDeploymentToDTO(&deployment))
		}
		writeJSON(w, http.StatusOK, out)
	}
}
