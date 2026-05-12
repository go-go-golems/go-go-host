package control

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-host/internal/store"
)

type ExportMetadata struct {
	Site         *store.Site            `json:"site"`
	Quota        *store.SiteQuota       `json:"quota,omitempty"`
	Config       []store.SiteConfigItem `json:"config"`
	Capabilities []store.SiteCapability `json:"capabilities"`
	Domains      []store.SiteDomain     `json:"domains"`
	Deployments  []store.Deployment     `json:"deployments"`
	ExportedAt   string                 `json:"exportedAt"`
}

type PruneInput struct {
	ActorUserID string
	SiteID      string
	Statuses    []string
	OlderThan   time.Time
	KeepLatest  int
	DryRun      bool
}

type PruneResult struct {
	DryRun     bool               `json:"dryRun"`
	Deleted    int                `json:"deleted"`
	Candidates []store.Deployment `json:"candidates"`
}

type MaintenanceService struct {
	store   *store.Store
	dataDir string
}

func (s *MaintenanceService) ExportSiteMetadata(ctx context.Context, actorUserID, siteID string) (*ExportMetadata, error) {
	site, err := s.authorizeSiteView(ctx, actorUserID, siteID)
	if err != nil {
		return nil, err
	}
	quota, _ := s.store.GetSiteQuota(ctx, siteID)
	config, err := s.store.ListSiteConfig(ctx, siteID)
	if err != nil {
		return nil, err
	}
	caps, err := s.store.ListSiteCapabilities(ctx, siteID)
	if err != nil {
		return nil, err
	}
	domains, err := s.store.ListSiteDomains(ctx, siteID)
	if err != nil {
		return nil, err
	}
	deployments, err := s.store.ListDeploymentsBySite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	return &ExportMetadata{Site: site, Quota: quota, Config: config, Capabilities: caps, Domains: domains, Deployments: deployments, ExportedAt: time.Now().UTC().Format(time.RFC3339)}, nil
}

func (s *MaintenanceService) SiteDBPath(ctx context.Context, actorUserID, siteID string) (string, error) {
	if _, err := s.authorizeSiteView(ctx, actorUserID, siteID); err != nil {
		return "", err
	}
	path := filepath.Join(s.dataDir, "sites", siteID, "db", "app.sqlite")
	if err := ensureInside(s.dataDir, path); err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}

func (s *MaintenanceService) DeploymentBundlePath(ctx context.Context, actorUserID, deploymentID string) (string, *store.Deployment, error) {
	dep, err := s.store.GetDeployment(ctx, deploymentID)
	if err != nil {
		return "", nil, err
	}
	if _, err := s.authorizeSiteView(ctx, actorUserID, dep.SiteID); err != nil {
		return "", nil, err
	}
	if dep.BundleRef == "" || dep.BundleRef == "pending" {
		return "", nil, os.ErrNotExist
	}
	if err := ensureInside(s.dataDir, dep.BundleRef); err != nil {
		return "", nil, err
	}
	if _, err := os.Stat(dep.BundleRef); err != nil {
		return "", nil, err
	}
	return dep.BundleRef, dep, nil
}

func (s *MaintenanceService) PruneDeployments(ctx context.Context, input PruneInput) (*PruneResult, error) {
	site, err := s.authorizeSiteDeploy(ctx, input.ActorUserID, input.SiteID)
	if err != nil {
		return nil, err
	}
	statuses := input.Statuses
	if len(statuses) == 0 {
		statuses = []string{store.DeploymentStatusRejected, store.DeploymentStatusSuperseded}
	}
	if input.OlderThan.IsZero() {
		input.OlderThan = time.Now().UTC()
	}
	candidates, err := s.store.ListPrunableDeployments(ctx, input.SiteID, site.ActiveDeploymentID, statuses, input.OlderThan)
	if err != nil {
		return nil, err
	}
	if input.KeepLatest > 0 && len(candidates) > input.KeepLatest {
		candidates = candidates[:len(candidates)-input.KeepLatest]
	} else if input.KeepLatest > 0 {
		candidates = nil
	}
	result := &PruneResult{DryRun: input.DryRun, Candidates: candidates}
	if input.DryRun {
		return result, nil
	}
	for _, dep := range candidates {
		_ = os.Remove(dep.BundleRef)
		if dep.UnpackedPath != "" {
			_ = os.RemoveAll(dep.UnpackedPath)
		}
		if err := s.store.DeleteDeployment(ctx, dep.ID); err != nil {
			return nil, err
		}
		result.Deleted++
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: input.ActorUserID, Action: "deployment.prune", ResourceType: "site", ResourceID: input.SiteID, MetadataJSON: fmt.Sprintf(`{"deleted":%d,"dryRun":%t}`, result.Deleted, input.DryRun)})
	return result, nil
}

func (s *MaintenanceService) RetainAudit(ctx context.Context, actorUserID string, olderThan time.Time) (int64, error) {
	if ok, err := s.store.IsPlatformAdmin(ctx, actorUserID); err != nil || !ok {
		if err != nil {
			return 0, err
		}
		return 0, ErrPermissionDenied
	}
	if olderThan.IsZero() {
		return 0, errors.New("olderThan is required")
	}
	count, err := s.store.DeleteAuditEventsBefore(ctx, olderThan)
	if err != nil {
		return 0, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{ActorType: "user", ActorID: actorUserID, Action: "audit.retention", ResourceType: "audit_log", ResourceID: "global", MetadataJSON: fmt.Sprintf(`{"deleted":%d,"olderThan":%q}`, count, olderThan.Format(time.RFC3339))})
	return count, nil
}

func (s *MaintenanceService) authorizeSiteView(ctx context.Context, actorUserID, siteID string) (*store.Site, error) {
	if s.store == nil {
		return nil, errors.New("store is not configured")
	}
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := ensureViewRole(ctx, s.store, actorUserID, site.OrgID); err != nil {
		return nil, err
	}
	return site, nil
}

func (s *MaintenanceService) authorizeSiteDeploy(ctx context.Context, actorUserID, siteID string) (*store.Site, error) {
	if s.store == nil {
		return nil, errors.New("store is not configured")
	}
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := ensureDeployRole(ctx, s.store, actorUserID, site.OrgID); err != nil {
		return nil, err
	}
	return site, nil
}

func ensureInside(root, path string) error {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(rootAbs, pathAbs)
	if err != nil {
		return err
	}
	if rel == "." || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return fmt.Errorf("path escapes data dir")
	}
	return nil
}
