package control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-go-golems/go-go-host/internal/deploy"
	hostruntime "github.com/go-go-golems/go-go-host/internal/runtime"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type DeploymentService struct {
	store      *store.Store
	supervisor *hostruntime.Supervisor
	dataDir    string
}

type UploadDeploymentInput struct {
	ActorUserID  string
	ActorType    string
	ActorID      string
	SiteID       string
	BundlePath   string
	Message      string
	Channel      string
	AllowedPaths []string
}

type DeploymentResult struct {
	Deployment *store.Deployment       `json:"deployment"`
	Report     deploy.ValidationReport `json:"report"`
	Manifest   deploy.Manifest         `json:"manifest"`
}

func (s *DeploymentService) Upload(ctx context.Context, input UploadDeploymentInput) (*DeploymentResult, error) {
	if s.store == nil {
		return nil, errors.New("store is not configured")
	}
	site, err := s.store.GetSite(ctx, input.SiteID)
	if err != nil {
		return nil, err
	}
	actorType := input.ActorType
	actorID := input.ActorID
	if actorType == "" {
		actorType = "user"
	}
	if actorID == "" {
		actorID = input.ActorUserID
	}
	if actorType == "user" {
		if err := ensureDeployRole(ctx, s.store, input.ActorUserID, site.OrgID); err != nil {
			return nil, err
		}
	} else if actorType != "agent" {
		return nil, ErrPermissionDenied
	}
	quota, err := s.store.GetSiteQuota(ctx, site.ID)
	if err != nil {
		return nil, err
	}
	versionDir := filepath.Join(s.dataDir, "sites", site.ID, "incoming")
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		return nil, err
	}
	// Reserve a deployment id/version before validation so immutable paths are stable.
	placeholderReport := deploy.ValidationReport{Valid: true}
	placeholderBytes, _ := json.Marshal(placeholderReport)
	dep, err := s.store.CreateDeployment(ctx, store.CreateDeploymentInput{
		SiteID: site.ID, Status: store.DeploymentStatusUploaded, BundleRef: "pending", ManifestJSON: []byte("{}"), ValidationJSON: placeholderBytes, CreatedByType: actorType, CreatedByID: actorID,
	})
	if err != nil {
		return nil, err
	}
	archiveDest := filepath.Join(s.dataDir, "bundles", site.ID, dep.ID+filepath.Ext(input.BundlePath))
	unpackDest := filepath.Join(s.dataDir, "sites", site.ID, "deployments", dep.ID)
	policyCaps, err := s.siteCapabilityPolicy(ctx, site.ID)
	if err != nil {
		return nil, err
	}
	prepared, err := deploy.ValidateAndStore(ctx, input.BundlePath, archiveDest, unpackDest, deploy.Options{MaxBytes: quota.BundleMaxBytes, AllowedPaths: input.AllowedPaths, Channel: input.Channel, PolicyCaps: policyCaps})
	if err != nil {
		return nil, err
	}
	validationJSON, _ := json.Marshal(prepared.Report)
	manifestJSON := prepared.ManifestJSON
	if len(manifestJSON) == 0 {
		manifestJSON = []byte("{}")
	}
	if prepared.Report.Valid {
		drySpec := hostruntime.Spec{SiteID: site.ID, OrgID: site.OrgID, DeploymentID: dep.ID, Hosts: []string{site.PrimaryHost}, ScriptsDir: filepath.Join(unpackDest, filepath.FromSlash(prepared.Manifest.ScriptsDir)), AssetsDir: filepath.Join(unpackDest, filepath.FromSlash(prepared.Manifest.AssetsDir)), DBPath: filepath.Join(s.dataDir, "sites", site.ID, "dry-run", dep.ID+".sqlite"), Dev: true, HealthPath: prepared.Manifest.SmokePath, Capabilities: hostruntime.DefaultCapabilities(), DBSoftMaxBytes: quota.DBSoftMaxBytes, DBHardMaxBytes: quota.DBHardMaxBytes, RequestTimeoutMS: quota.RequestTimeoutMS}
		dryRuntime, dryErr := hostruntime.NewSiteRuntime(ctx, drySpec)
		if dryErr != nil {
			prepared.Report.Valid = false
			prepared.Report.Errors = append(prepared.Report.Errors, "dry-run runtime load failed: "+dryErr.Error())
		} else {
			if err := dryRuntime.HealthCheck(ctx); err != nil {
				prepared.Report.Valid = false
				prepared.Report.Errors = append(prepared.Report.Errors, "dry-run smoke check failed: "+err.Error())
			}
			_ = dryRuntime.Close(ctx)
		}
		validationJSON, _ = json.Marshal(prepared.Report)
	}
	status := store.DeploymentStatusValidated
	if !prepared.Report.Valid {
		status = store.DeploymentStatusRejected
	}
	if err := s.store.UpdateDeploymentArtifacts(ctx, dep.ID, status, archiveDest, unpackDest, manifestJSON, validationJSON, prepared.BundleSHA256); err != nil {
		return nil, err
	}
	dep, err = s.store.GetDeployment(ctx, dep.ID)
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: actorType, ActorID: actorID, Action: "deployment.upload", ResourceType: "deployment", ResourceID: dep.ID})
	if status == store.DeploymentStatusRejected {
		_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: actorType, ActorID: actorID, Action: "deployment.validation_failed", ResourceType: "deployment", ResourceID: dep.ID})
	}
	return &DeploymentResult{Deployment: dep, Report: prepared.Report, Manifest: prepared.Manifest}, nil
}

func (s *DeploymentService) List(ctx context.Context, actorUserID, siteID string) ([]store.Deployment, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := ensureViewRole(ctx, s.store, actorUserID, site.OrgID); err != nil {
		return nil, err
	}
	return s.store.ListDeploymentsBySite(ctx, siteID)
}

func (s *DeploymentService) Get(ctx context.Context, actorUserID, deploymentID string) (*store.Deployment, error) {
	dep, err := s.store.GetDeployment(ctx, deploymentID)
	if err != nil {
		return nil, err
	}
	site, err := s.store.GetSite(ctx, dep.SiteID)
	if err != nil {
		return nil, err
	}
	if err := ensureViewRole(ctx, s.store, actorUserID, site.OrgID); err != nil {
		return nil, err
	}
	return dep, nil
}

func (s *DeploymentService) Activate(ctx context.Context, actorUserID, deploymentID string) (*store.Deployment, error) {
	dep, err := s.store.GetDeployment(ctx, deploymentID)
	if err != nil {
		return nil, err
	}
	site, err := s.store.GetSite(ctx, dep.SiteID)
	if err != nil {
		return nil, err
	}
	if err := ensureDeployRole(ctx, s.store, actorUserID, site.OrgID); err != nil {
		return nil, err
	}
	return s.activate(ctx, "user", actorUserID, dep, site)
}

func (s *DeploymentService) ActivateAsAgent(ctx context.Context, agentID, deploymentID string) (*store.Deployment, error) {
	dep, err := s.store.GetDeployment(ctx, deploymentID)
	if err != nil {
		return nil, err
	}
	site, err := s.store.GetSite(ctx, dep.SiteID)
	if err != nil {
		return nil, err
	}
	agent, err := s.store.GetAgent(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if agent.OrgID != site.OrgID || agent.Status != store.AgentStatusActive {
		return nil, ErrPermissionDenied
	}
	grants, err := s.store.ListAgentSiteGrants(ctx, agentID)
	if err != nil {
		return nil, err
	}
	allowed := false
	for _, grant := range grants {
		if grant.SiteID == site.ID && grant.CanActivate && (grant.ExpiresAt.IsZero() || time.Now().UTC().Before(grant.ExpiresAt)) {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, ErrPermissionDenied
	}
	return s.activate(ctx, "agent", agentID, dep, site)
}

func (s *DeploymentService) RestoreActiveRuntimes(ctx context.Context) error {
	deployments, err := s.store.ListAdminDeployments(ctx, store.AdminDeploymentFilter{Status: store.DeploymentStatusActive, Limit: 1000})
	if err != nil {
		return err
	}
	for _, row := range deployments {
		dep := &store.Deployment{ID: row.ID, SiteID: row.SiteID, Version: row.Version, Status: row.Status, BundleRef: row.BundleRef, UnpackedPath: row.UnpackedPath, ManifestJSON: row.ManifestJSON, ValidationJSON: row.ValidationJSON, CreatedByType: row.CreatedByType, CreatedByID: row.CreatedByID, CreatedAt: row.CreatedAt, ActivatedAt: row.ActivatedAt, BundleSHA256: row.BundleSHA256}
		site := &store.Site{ID: row.SiteID, OrgID: row.OrgID, Slug: row.SiteSlug, Name: row.SiteSlug, PrimaryHost: row.PrimaryHost, Status: store.SiteStatusActive, ActiveDeploymentID: row.ID}
		if _, err := s.activate(ctx, "system", "startup-restore", dep, site); err != nil {
			return fmt.Errorf("restore active runtime for site %s deployment %s: %w", row.SiteID, row.ID, err)
		}
	}
	return nil
}

func (s *DeploymentService) activate(ctx context.Context, actorType, actorID string, dep *store.Deployment, site *store.Site) (*store.Deployment, error) {
	if dep.Status != store.DeploymentStatusValidated && dep.Status != store.DeploymentStatusSuperseded && dep.Status != store.DeploymentStatusActive {
		return nil, fmt.Errorf("deployment %s is not activatable from status %q", dep.ID, dep.Status)
	}
	manifest, err := manifestFromDeployment(dep)
	if err != nil {
		return nil, err
	}
	unpackedPath := dep.UnpackedPath
	if unpackedPath == "" {
		unpackedPath = filepath.Join(s.dataDir, "sites", site.ID, "deployments", dep.ID)
	}
	quota, _ := s.store.GetSiteQuota(ctx, site.ID)
	hosts := []string{site.PrimaryHost}
	if domains, err := s.store.ListVerifiedSiteDomains(ctx, site.ID); err == nil {
		for _, domain := range domains {
			hosts = append(hosts, domain.Hostname)
		}
	}
	spec := hostruntime.Spec{SiteID: site.ID, OrgID: site.OrgID, DeploymentID: dep.ID, Hosts: hosts, ScriptsDir: filepath.Join(unpackedPath, filepath.FromSlash(manifest.ScriptsDir)), AssetsDir: filepath.Join(unpackedPath, filepath.FromSlash(manifest.AssetsDir)), DBPath: filepath.Join(s.dataDir, "sites", site.ID, "db", "app.sqlite"), Dev: true, HealthPath: manifest.SmokePath, Capabilities: hostruntime.DefaultCapabilities()}
	if quota != nil {
		spec.DBSoftMaxBytes = quota.DBSoftMaxBytes
		spec.DBHardMaxBytes = quota.DBHardMaxBytes
		spec.RequestTimeoutMS = quota.RequestTimeoutMS
	}
	if err := s.supervisor.Activate(ctx, spec); err != nil {
		return nil, err
	}
	if err := s.store.MarkDeploymentActive(ctx, site.ID, dep.ID); err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: actorType, ActorID: actorID, Action: "deployment.activate", ResourceType: "deployment", ResourceID: dep.ID})
	return s.store.GetDeployment(ctx, dep.ID)
}

func (s *DeploymentService) Rollback(ctx context.Context, actorUserID, siteID string) (*store.Deployment, error) {
	site, err := s.store.GetSite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if err := ensureDeployRole(ctx, s.store, actorUserID, site.OrgID); err != nil {
		return nil, err
	}
	previous, err := s.store.PreviousValidatedDeployment(ctx, siteID, site.ActiveDeploymentID)
	if err != nil {
		return nil, err
	}
	dep, err := s.Activate(ctx, actorUserID, previous.ID)
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: site.OrgID, ActorType: "user", ActorID: actorUserID, Action: "deployment.rollback", ResourceType: "deployment", ResourceID: dep.ID})
	return dep, nil
}

func (s *DeploymentService) siteCapabilityPolicy(ctx context.Context, siteID string) (map[string]bool, error) {
	caps, err := s.store.ListSiteCapabilities(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if len(caps) == 0 {
		return deploy.SafeCapabilities, nil
	}
	policy := map[string]bool{}
	for _, cap := range caps {
		policy[cap.Capability] = cap.Enabled
	}
	return policy, nil
}

func manifestFromDeployment(dep *store.Deployment) (deploy.Manifest, error) {
	var manifest deploy.Manifest
	if err := json.Unmarshal(dep.ManifestJSON, &manifest); err != nil {
		return manifest, err
	}
	return manifest, nil
}

func ensureOwnerRole(ctx context.Context, st *store.Store, userID, orgID string) error {
	role, err := st.MembershipRole(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if role == store.RoleOrgOwner {
		return nil
	}
	return ErrPermissionDenied
}

func ensureDeployRole(ctx context.Context, st *store.Store, userID, orgID string) error {
	role, err := st.MembershipRole(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if role == store.RoleOrgOwner || role == store.RoleOrgDeveloper {
		return nil
	}
	return ErrPermissionDenied
}

func ensureViewRole(ctx context.Context, st *store.Store, userID, orgID string) error {
	role, err := st.MembershipRole(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if role != "" {
		return nil
	}
	return ErrPermissionDenied
}
