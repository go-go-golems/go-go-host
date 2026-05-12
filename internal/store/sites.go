package store

import (
	"context"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

type CreateSiteInput struct {
	OrgID       string
	Slug        string
	Name        string
	PrimaryHost string
}

func (s *Store) CreateSite(ctx context.Context, input CreateSiteInput) (*Site, error) {
	name := input.Name
	if name == "" {
		name = input.Slug
	}
	row, err := s.q.CreateSite(ctx, storedb.CreateSiteParams{ID: newID("site"), OrgID: input.OrgID, Slug: input.Slug, Name: name, PrimaryHost: input.PrimaryHost, Status: SiteStatusProvisioning, CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return siteFromDB(row), nil
}

func (s *Store) GetSite(ctx context.Context, id string) (*Site, error) {
	row, err := s.q.GetSite(ctx, id)
	if err != nil {
		return nil, err
	}
	return siteFromDB(row), nil
}

func (s *Store) ListSitesByOrg(ctx context.Context, orgID string) ([]Site, error) {
	rows, err := s.q.ListSitesByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	sites := make([]Site, 0, len(rows))
	for _, row := range rows {
		sites = append(sites, *siteFromDB(row))
	}
	return sites, nil
}

func (s *Store) UpdateSiteStatus(ctx context.Context, siteID, status string) error {
	return s.q.UpdateSiteStatus(ctx, storedb.UpdateSiteStatusParams{ID: siteID, Status: status})
}

func (s *Store) UpdateSiteActiveDeployment(ctx context.Context, siteID, deploymentID string) error {
	return s.q.UpdateSiteActiveDeployment(ctx, storedb.UpdateSiteActiveDeploymentParams{ID: siteID, ActiveDeploymentID: deploymentID})
}

func (s *Store) CreateDefaultSiteQuota(ctx context.Context, siteID string) error {
	return s.q.CreateDefaultSiteQuota(ctx, storedb.CreateDefaultSiteQuotaParams{SiteID: siteID, BundleMaxBytes: int64(50 * 1024 * 1024), DbSoftMaxBytes: int64(50 * 1024 * 1024), DbHardMaxBytes: int64(100 * 1024 * 1024), RequestTimeoutMs: 2000, UpdatedAt: pgTime(now())})
}

func (s *Store) CreateDefaultSiteCapabilities(ctx context.Context, siteID string) error {
	for _, capability := range []string{"express", "ui.dsl", "database", "db", "time", "timer", "assets"} {
		if err := s.q.UpsertSiteCapability(ctx, storedb.UpsertSiteCapabilityParams{SiteID: siteID, Capability: capability, Enabled: true, ConfigJson: []byte("{}"), UpdatedAt: pgTime(now())}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListSiteCapabilities(ctx context.Context, siteID string) ([]SiteCapability, error) {
	rows, err := s.q.ListSiteCapabilities(ctx, siteID)
	if err != nil {
		return nil, err
	}
	out := make([]SiteCapability, 0, len(rows))
	for _, row := range rows {
		out = append(out, SiteCapability{SiteID: row.SiteID, Capability: row.Capability, Enabled: row.Enabled, ConfigJSON: row.ConfigJson, UpdatedAt: fromPgTime(row.UpdatedAt)})
	}
	return out, nil
}

func (s *Store) GetSiteQuota(ctx context.Context, siteID string) (*SiteQuota, error) {
	row, err := s.q.GetSiteQuota(ctx, siteID)
	if err != nil {
		return nil, err
	}
	return &SiteQuota{SiteID: row.SiteID, BundleMaxBytes: row.BundleMaxBytes, DBSoftMaxBytes: row.DbSoftMaxBytes, DBHardMaxBytes: row.DbHardMaxBytes, RequestTimeoutMS: int(row.RequestTimeoutMs), UpdatedAt: fromPgTime(row.UpdatedAt)}, nil
}

func siteFromDB(row storedb.Site) *Site {
	return &Site{ID: row.ID, OrgID: row.OrgID, Slug: row.Slug, Name: row.Name, PrimaryHost: row.PrimaryHost, Status: row.Status, ActiveDeploymentID: row.ActiveDeploymentID, CreatedAt: fromPgTime(row.CreatedAt)}
}
