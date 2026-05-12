package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

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

func (s *Store) UpsertSiteCapability(ctx context.Context, siteID, capability string, enabled bool, configJSON []byte) error {
	return s.q.UpsertSiteCapability(ctx, storedb.UpsertSiteCapabilityParams{SiteID: siteID, Capability: capability, Enabled: enabled, ConfigJson: configJSON, UpdatedAt: pgTime(now())})
}

func (s *Store) GetSiteQuota(ctx context.Context, siteID string) (*SiteQuota, error) {
	row, err := s.q.GetSiteQuota(ctx, siteID)
	if err != nil {
		return nil, err
	}
	return &SiteQuota{SiteID: row.SiteID, BundleMaxBytes: row.BundleMaxBytes, DBSoftMaxBytes: row.DbSoftMaxBytes, DBHardMaxBytes: row.DbHardMaxBytes, RequestTimeoutMS: int(row.RequestTimeoutMs), UpdatedAt: fromPgTime(row.UpdatedAt)}, nil
}

func (s *Store) ListSiteConfig(ctx context.Context, siteID string) ([]SiteConfigItem, error) {
	rows, err := s.q.ListSiteConfig(ctx, siteID)
	if err != nil {
		return nil, err
	}
	out := make([]SiteConfigItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, SiteConfigItem{SiteID: row.SiteID, Key: row.Key, ValueJSON: row.ValueJson, UpdatedAt: fromPgTime(row.UpdatedAt)})
	}
	return out, nil
}

func (s *Store) UpsertSiteConfig(ctx context.Context, siteID, key string, valueJSON []byte) error {
	return s.q.UpsertSiteConfig(ctx, storedb.UpsertSiteConfigParams{SiteID: siteID, Key: key, ValueJson: valueJSON, UpdatedAt: pgTime(now())})
}

func (s *Store) DeleteSiteConfig(ctx context.Context, siteID, key string) error {
	return s.q.DeleteSiteConfig(ctx, storedb.DeleteSiteConfigParams{SiteID: siteID, Key: key})
}

func (s *Store) CreateSiteDomain(ctx context.Context, siteID, hostname string) (*SiteDomain, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}
	row, err := s.q.CreateSiteDomain(ctx, storedb.CreateSiteDomainParams{ID: newID("dom"), SiteID: siteID, Hostname: strings.ToLower(strings.TrimSuffix(hostname, ".")), Status: "pending", VerificationToken: "ggh-verify-" + hex.EncodeToString(token), CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return siteDomainFromDB(row), nil
}

func (s *Store) ListSiteDomains(ctx context.Context, siteID string) ([]SiteDomain, error) {
	rows, err := s.q.ListSiteDomains(ctx, siteID)
	if err != nil {
		return nil, err
	}
	out := make([]SiteDomain, 0, len(rows))
	for _, row := range rows {
		out = append(out, *siteDomainFromDB(row))
	}
	return out, nil
}

func (s *Store) ListVerifiedSiteDomains(ctx context.Context, siteID string) ([]SiteDomain, error) {
	rows, err := s.q.ListVerifiedSiteDomains(ctx, siteID)
	if err != nil {
		return nil, err
	}
	out := make([]SiteDomain, 0, len(rows))
	for _, row := range rows {
		out = append(out, *siteDomainFromDB(row))
	}
	return out, nil
}

func (s *Store) GetSiteDomain(ctx context.Context, id string) (*SiteDomain, error) {
	row, err := s.q.GetSiteDomain(ctx, id)
	if err != nil {
		return nil, err
	}
	return siteDomainFromDB(row), nil
}

func (s *Store) VerifySiteDomain(ctx context.Context, id string) (*SiteDomain, error) {
	row, err := s.q.VerifySiteDomain(ctx, storedb.VerifySiteDomainParams{ID: id, VerifiedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return siteDomainFromDB(row), nil
}

func (s *Store) DeleteSiteDomain(ctx context.Context, id string) error {
	return s.q.DeleteSiteDomain(ctx, id)
}

func siteFromDB(row storedb.Site) *Site {
	return &Site{ID: row.ID, OrgID: row.OrgID, Slug: row.Slug, Name: row.Name, PrimaryHost: row.PrimaryHost, Status: row.Status, ActiveDeploymentID: row.ActiveDeploymentID, CreatedAt: fromPgTime(row.CreatedAt)}
}

func siteDomainFromDB(row storedb.SiteDomain) *SiteDomain {
	return &SiteDomain{ID: row.ID, SiteID: row.SiteID, Hostname: row.Hostname, Status: row.Status, VerificationToken: row.VerificationToken, VerifiedAt: fromPgTime(row.VerifiedAt), CreatedAt: fromPgTime(row.CreatedAt)}
}
