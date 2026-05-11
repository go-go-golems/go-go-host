package httpapi

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type deploymentDTO struct {
	ID             string `json:"id"`
	SiteID         string `json:"siteId"`
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

type rollbackRequest struct {
	SiteID string `json:"siteId"`
}

func handleUploadDeployment(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "expected multipart form with bundle file")
			return
		}
		file, header, err := r.FormFile("bundle")
		if err != nil {
			writeError(w, http.StatusBadRequest, "missing multipart bundle file")
			return
		}
		defer file.Close()
		tmp, err := os.CreateTemp("", "go-go-host-bundle-*"+filepath.Ext(header.Filename))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer os.Remove(tmp.Name())
		if _, err := io.Copy(tmp, file); err != nil {
			_ = tmp.Close()
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = tmp.Close()
		result, err := core.Deployments.Upload(r.Context(), control.UploadDeploymentInput{ActorUserID: p.User.ID, SiteID: r.PathValue("site_id"), BundlePath: tmp.Name(), Message: r.FormValue("message"), Channel: r.FormValue("channel")})
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		status := http.StatusCreated
		if !result.Report.Valid {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]any{"deployment": deploymentToDTO(result.Deployment), "report": result.Report, "manifest": result.Manifest})
	}
}

func handleListDeployments(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		deps, err := core.Deployments.List(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		out := make([]deploymentDTO, 0, len(deps))
		for _, dep := range deps {
			out = append(out, deploymentToDTO(&dep))
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleGetDeployment(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		dep, err := core.Deployments.Get(r.Context(), p.User.ID, r.PathValue("deployment_id"))
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, deploymentToDTO(dep))
	}
}

func handleActivateDeployment(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		dep, err := core.Deployments.Activate(r.Context(), p.User.ID, r.PathValue("deployment_id"))
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, deploymentToDTO(dep))
	}
}

func handleRollbackDeployment(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		dep, err := core.Deployments.Rollback(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, deploymentToDTO(dep))
	}
}

func writeDeploymentError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, control.ErrPermissionDenied) {
		status = http.StatusForbidden
	}
	writeError(w, status, err.Error())
}

func deploymentToDTO(dep *store.Deployment) deploymentDTO {
	return deploymentDTO{ID: dep.ID, SiteID: dep.SiteID, Version: dep.Version, Status: dep.Status, BundleRef: dep.BundleRef, UnpackedPath: dep.UnpackedPath, ManifestJSON: string(dep.ManifestJSON), ValidationJSON: string(dep.ValidationJSON), CreatedByType: dep.CreatedByType, CreatedByID: dep.CreatedByID, CreatedAt: dep.CreatedAt, ActivatedAt: dep.ActivatedAt}
}
