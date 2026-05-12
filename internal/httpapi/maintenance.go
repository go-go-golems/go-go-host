package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-host/internal/control"
)

type pruneDeploymentsRequest struct {
	Statuses   []string `json:"statuses"`
	OlderThan  string   `json:"olderThan"`
	KeepLatest int      `json:"keepLatest"`
	DryRun     bool     `json:"dryRun"`
}

type auditRetentionRequest struct {
	OlderThan string `json:"olderThan"`
	DryRun    bool   `json:"dryRun"`
}

func handleExportSiteMetadata(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		meta, err := core.Maintenance.ExportSiteMetadata(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeControlError(w, err)
			return
		}
		filename := fmt.Sprintf("go-go-host-%s-metadata.json", r.PathValue("site_id"))
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		writeJSON(w, http.StatusOK, meta)
	}
}

func handleExportSiteDB(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		path, err := core.Maintenance.SiteDBPath(r.Context(), p.User.ID, r.PathValue("site_id"))
		if err != nil {
			writeDownloadError(w, err)
			return
		}
		w.Header().Set("Content-Disposition", `attachment; filename="`+r.PathValue("site_id")+`.sqlite"`)
		http.ServeFile(w, r, path)
	}
}

func handleExportDeploymentBundle(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		path, dep, err := core.Maintenance.DeploymentBundlePath(r.Context(), p.User.ID, r.PathValue("deployment_id"))
		if err != nil {
			writeDownloadError(w, err)
			return
		}
		name := dep.ID + filepath.Ext(path)
		w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
		http.ServeFile(w, r, path)
	}
}

func handlePruneDeployments(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req pruneDeploymentsRequest
		if r.Body != nil && r.ContentLength != 0 {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}
		}
		olderThan, err := parseMaintenanceTime(req.OlderThan, time.Now().UTC())
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := core.Maintenance.PruneDeployments(r.Context(), control.PruneInput{ActorUserID: p.User.ID, SiteID: r.PathValue("site_id"), Statuses: req.Statuses, OlderThan: olderThan, KeepLatest: req.KeepLatest, DryRun: req.DryRun})
		if err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func handleAuditRetention(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req auditRetentionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		olderThan, err := parseMaintenanceTime(req.OlderThan, time.Time{})
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if req.DryRun {
			writeJSON(w, http.StatusOK, map[string]any{"dryRun": true, "olderThan": olderThan.Format(time.RFC3339), "message": "audit retention dry-run does not delete rows; run with dryRun=false to delete"})
			return
		}
		deleted, err := core.Maintenance.RetainAudit(r.Context(), p.User.ID, olderThan)
		if err != nil {
			writeControlError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"dryRun": false, "deleted": deleted, "olderThan": olderThan.Format(time.RFC3339)})
	}
}

func parseMaintenanceTime(raw string, fallback time.Time) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if fallback.IsZero() {
			return time.Time{}, errors.New("olderThan is required")
		}
		return fallback, nil
	}
	if strings.HasSuffix(raw, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(raw, "d"))
		if err != nil || days < 0 {
			return time.Time{}, fmt.Errorf("invalid day duration %q", raw)
		}
		return time.Now().UTC().Add(-time.Duration(days) * 24 * time.Hour), nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RFC3339 time %q", raw)
	}
	return parsed.UTC(), nil
}

func writeDownloadError(w http.ResponseWriter, err error) {
	if errors.Is(err, control.ErrPermissionDenied) {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeError(w, http.StatusNotFound, err.Error())
}
