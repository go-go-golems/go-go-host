package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type siteConfigDTO struct {
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	UpdatedAt string          `json:"updatedAt"`
}

type upsertSiteConfigRequest struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type deleteSiteConfigRequest struct {
	Key string `json:"key"`
}

type siteCapabilityDTO struct {
	SiteID     string          `json:"siteId"`
	Capability string          `json:"capability"`
	Enabled    bool            `json:"enabled"`
	Config     json.RawMessage `json:"config"`
	UpdatedAt  string          `json:"updatedAt"`
}

type upsertSiteCapabilityRequest struct {
	Capability string          `json:"capability"`
	Enabled    bool            `json:"enabled"`
	Config     json.RawMessage `json:"config"`
}

type siteDomainDTO struct {
	ID                string `json:"id"`
	SiteID            string `json:"siteId"`
	Hostname          string `json:"hostname"`
	Status            string `json:"status"`
	VerificationToken string `json:"verificationToken"`
	VerifiedAt        string `json:"verifiedAt,omitempty"`
	CreatedAt         string `json:"createdAt"`
}

type addSiteDomainRequest struct {
	Hostname string `json:"hostname"`
}

type siteEnvironmentPlaceholderDTO struct {
	SiteID       string   `json:"siteId"`
	Status       string   `json:"status"`
	Supported    []string `json:"supported"`
	NotSupported []string `json:"notSupported"`
	Message      string   `json:"message"`
}

func handleListSiteConfig(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		items, err := core.Sites.ListConfig(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeControlError(w, err)
			return
		}
		out := make([]siteConfigDTO, 0, len(items))
		for _, item := range items {
			out = append(out, siteConfigDTO{Key: item.Key, Value: json.RawMessage(item.ValueJSON), UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleUpsertSiteConfig(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req upsertSiteConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if len(req.Value) == 0 {
			req.Value = []byte("null")
		}
		if err := core.Sites.UpsertConfig(r.Context(), p.User.ID, r.PathValue("site_id"), req.Key, req.Value); err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleDeleteSiteConfig(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req deleteSiteConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if err := core.Sites.DeleteConfig(r.Context(), p.User.ID, r.PathValue("site_id"), req.Key); err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleListSiteCapabilities(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		caps, err := core.Sites.ListCapabilities(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeControlError(w, err)
			return
		}
		out := make([]siteCapabilityDTO, 0, len(caps))
		for _, cap := range caps {
			out = append(out, siteCapabilityDTO{SiteID: cap.SiteID, Capability: cap.Capability, Enabled: cap.Enabled, Config: json.RawMessage(cap.ConfigJSON), UpdatedAt: cap.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")})
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleUpsertSiteCapability(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req upsertSiteCapabilityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if len(req.Config) == 0 {
			req.Config = []byte("{}")
		}
		if err := core.Sites.UpsertCapability(r.Context(), p.User.ID, r.PathValue("site_id"), req.Capability, req.Enabled, req.Config); err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleListSiteDomains(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		domains, err := core.Sites.ListDomains(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeControlError(w, err)
			return
		}
		out := make([]siteDomainDTO, 0, len(domains))
		for _, domain := range domains {
			out = append(out, siteDomainToDTO(domain))
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleAddSiteDomain(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req addSiteDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		domain, err := core.Sites.AddDomain(r.Context(), p.User.ID, r.PathValue("site_id"), req.Hostname)
		if err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, siteDomainToDTO(*domain))
	}
}

func handleVerifySiteDomain(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		domain, err := core.Sites.VerifyDomain(r.Context(), p.User.ID, r.PathValue("site_id"), r.PathValue("domain_id"))
		if err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, siteDomainToDTO(*domain))
	}
}

func handleDeleteSiteDomain(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if err := core.Sites.DeleteDomain(r.Context(), p.User.ID, r.PathValue("site_id"), r.PathValue("domain_id")); err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleSiteEnvironmentPlaceholder(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		siteID := r.PathValue("site_id")
		if _, err := core.Sites.ListConfig(r.Context(), p.User.ID, siteID); err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, siteEnvironmentPlaceholderDTO{SiteID: siteID, Status: "design-placeholder", Supported: []string{"non-secret site config via /config"}, NotSupported: []string{"process env passthrough", "plaintext secret values in API responses", "runtime require('fs') secret reads"}, Message: "Secrets/environment variables are intentionally not implemented in v1. Use non-secret site config only until encrypted secret storage and runtime injection are designed."})
	}
}

func siteDomainToDTO(domain store.SiteDomain) siteDomainDTO {
	verifiedAt := ""
	if !domain.VerifiedAt.IsZero() {
		verifiedAt = domain.VerifiedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return siteDomainDTO{ID: domain.ID, SiteID: domain.SiteID, Hostname: domain.Hostname, Status: domain.Status, VerificationToken: domain.VerificationToken, VerifiedAt: verifiedAt, CreatedAt: domain.CreatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

func writeControlError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, control.ErrPermissionDenied) {
		status = http.StatusForbidden
	}
	writeError(w, status, err.Error())
}
