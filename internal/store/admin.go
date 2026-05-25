package store

import (
	"context"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type AdminOrg struct {
	ID              string
	Slug            string
	Name            string
	CreatedAt       string
	MemberCount     int64
	SiteCount       int64
	DeploymentCount int64
}

type AdminUser struct {
	ID            string
	Email         string
	DisplayName   string
	CreatedAt     string
	LastLoginAt   string
	PlatformAdmin bool
	OrgCount      int64
}

type AdminSite struct {
	ID                 string
	OrgID              string
	OrgSlug            string
	OrgName            string
	Slug               string
	Name               string
	PrimaryHost        string
	Status             string
	ActiveDeploymentID string
	CreatedAt          string
	RuntimeStatus      string
	RequestsTotal      int64
	ErrorsTotal        int64
	LastError          string
}

type AdminAgent struct {
	ID              string
	OrgID           string
	OrgSlug         string
	OrgName         string
	Name            string
	Status          string
	CreatedByUserID string
	CreatedAt       string
	LastSeenAt      string
	GrantCount      int64
}

type AdminDeployment struct {
	ID             string
	SiteID         string
	SiteSlug       string
	PrimaryHost    string
	OrgID          string
	OrgSlug        string
	OrgName        string
	Version        int
	Status         string
	BundleRef      string
	UnpackedPath   string
	ManifestJSON   []byte
	ValidationJSON []byte
	CreatedByType  string
	CreatedByID    string
	CreatedAt      string
	ActivatedAt    string
	BundleSHA256   string
}

type AdminDeploymentFilter struct {
	OrgID  string
	SiteID string
	Status string
	Limit  int
}

type AdminAgentFilter struct {
	OrgID  string
	Status string
}

type AdminQuota struct {
	SiteID           string
	SiteSlug         string
	PrimaryHost      string
	OrgID            string
	OrgSlug          string
	OrgName          string
	BundleMaxBytes   int64
	DBSoftMaxBytes   int64
	DBHardMaxBytes   int64
	RequestTimeoutMS int
	UpdatedAt        string
	RequestsTotal    int64
	ErrorsTotal      int64
}

type AdminCapability struct {
	SiteID     string
	SiteSlug   string
	OrgID      string
	OrgSlug    string
	OrgName    string
	Capability string
	Enabled    bool
	ConfigJSON []byte
	UpdatedAt  string
}

type AdminDomain struct {
	ID                string
	SiteID            string
	SiteSlug          string
	OrgID             string
	OrgSlug           string
	OrgName           string
	Hostname          string
	Status            string
	VerificationToken string
	VerifiedAt        string
	CreatedAt         string
}

func (s *Store) ListAdminOrgs(ctx context.Context) ([]AdminOrg, error) {
	rows, err := s.q.ListAdminOrgs(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminOrg, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminOrg{ID: row.ID, Slug: row.Slug, Name: row.Name, CreatedAt: formatDBTime(row.CreatedAt), MemberCount: row.MemberCount, SiteCount: row.SiteCount, DeploymentCount: row.DeploymentCount})
	}
	return out, nil
}

func (s *Store) ListAdminUsers(ctx context.Context) ([]AdminUser, error) {
	rows, err := s.q.ListAdminUsers(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminUser, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminUser{ID: row.ID, Email: row.Email, DisplayName: row.DisplayName, CreatedAt: formatDBTime(row.CreatedAt), LastLoginAt: formatDBTime(row.LastLoginAt), PlatformAdmin: row.PlatformAdmin, OrgCount: row.OrgCount})
	}
	return out, nil
}

func (s *Store) ListAdminSites(ctx context.Context) ([]AdminSite, error) {
	rows, err := s.q.ListAdminSites(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminSite, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminSite{ID: row.ID, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Slug: row.Slug, Name: row.Name, PrimaryHost: row.PrimaryHost, Status: row.Status, ActiveDeploymentID: row.ActiveDeploymentID, CreatedAt: formatDBTime(row.CreatedAt), RuntimeStatus: row.RuntimeStatus, RequestsTotal: row.RequestsTotal, ErrorsTotal: row.ErrorsTotal, LastError: row.LastError})
	}
	return out, nil
}

func (s *Store) ListAdminQuotas(ctx context.Context) ([]AdminQuota, error) {
	rows, err := s.q.ListAdminQuotas(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminQuota, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminQuota{SiteID: row.SiteID, SiteSlug: row.SiteSlug, PrimaryHost: row.PrimaryHost, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, BundleMaxBytes: row.BundleMaxBytes, DBSoftMaxBytes: row.DbSoftMaxBytes, DBHardMaxBytes: row.DbHardMaxBytes, RequestTimeoutMS: int(row.RequestTimeoutMs), UpdatedAt: formatDBTime(row.UpdatedAt), RequestsTotal: row.RequestsTotal, ErrorsTotal: row.ErrorsTotal})
	}
	return out, nil
}

func (s *Store) ListAdminCapabilities(ctx context.Context) ([]AdminCapability, error) {
	rows, err := s.q.ListAdminCapabilities(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminCapability, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminCapability{SiteID: row.SiteID, SiteSlug: row.SiteSlug, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Capability: row.Capability, Enabled: row.Enabled, ConfigJSON: row.ConfigJson, UpdatedAt: formatDBTime(row.UpdatedAt)})
	}
	return out, nil
}

func (s *Store) ListAdminDomains(ctx context.Context) ([]AdminDomain, error) {
	rows, err := s.q.ListAdminDomains(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AdminDomain, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminDomain{ID: row.ID, SiteID: row.SiteID, SiteSlug: row.SiteSlug, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Hostname: row.Hostname, Status: row.Status, VerificationToken: row.VerificationToken, VerifiedAt: formatDBTime(row.VerifiedAt), CreatedAt: formatDBTime(row.CreatedAt)})
	}
	return out, nil
}

func (s *Store) GetAdminDeployment(ctx context.Context, id string) (*AdminDeployment, error) {
	row, err := s.q.GetAdminDeployment(ctx, id)
	if err != nil {
		return nil, err
	}
	return adminDeploymentFromRow(adminDeploymentRow{
		ID: row.ID, SiteID: row.SiteID, SiteSlug: row.SiteSlug, PrimaryHost: row.PrimaryHost, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName,
		Version: row.Version, Status: row.Status, BundleRef: row.BundleRef, UnpackedPath: row.UnpackedPath, ManifestJson: row.ManifestJson, ValidationJson: row.ValidationJson,
		CreatedByType: row.CreatedByType, CreatedByID: row.CreatedByID, CreatedAt: row.CreatedAt, ActivatedAt: row.ActivatedAt, BundleSha256: row.BundleSha256,
	})
}

func (s *Store) ListAdminAgents(ctx context.Context, filter AdminAgentFilter) ([]AdminAgent, error) {
	rows, err := s.q.ListAdminAgents(ctx, storedb.ListAdminAgentsParams{OrgID: optionalText(filter.OrgID), Status: optionalText(filter.Status)})
	if err != nil {
		return nil, err
	}
	out := make([]AdminAgent, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminAgent{ID: row.ID, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Name: row.Name, Status: row.Status, CreatedByUserID: row.CreatedByUserID, CreatedAt: formatDBTime(row.CreatedAt), LastSeenAt: formatDBTime(row.LastSeenAt), GrantCount: row.GrantCount})
	}
	return out, nil
}

func (s *Store) ListAdminAuditEvents(ctx context.Context, filter AuditFilter) ([]AuditEvent, error) {
	limit := boundedListLimit(filter.Limit)
	rows, err := s.q.ListAdminAuditEvents(ctx, storedb.ListAdminAuditEventsParams{OrgID: optionalText(filter.OrgID), ResourceID: optionalText(filter.ResourceID), ActorType: optionalText(filter.ActorType), ActorID: optionalText(filter.ActorID), Action: optionalText(filter.Action), Limit: limit})
	if err != nil {
		return nil, err
	}
	events := make([]AuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, AuditEvent{ID: row.ID, OrgID: row.OrgID, ActorType: row.ActorType, ActorID: row.ActorID, Action: row.Action, ResourceType: row.ResourceType, ResourceID: row.ResourceID, IPAddress: row.IpAddress, UserAgent: row.UserAgent, MetadataJSON: string(row.MetadataJson), CreatedAt: fromPgTime(row.CreatedAt)})
	}
	return events, nil
}

func (s *Store) ListAdminDeployments(ctx context.Context, filter AdminDeploymentFilter) ([]AdminDeployment, error) {
	limit := boundedListLimit(filter.Limit)
	rows, err := s.q.ListAdminDeployments(ctx, storedb.ListAdminDeploymentsParams{OrgID: optionalText(filter.OrgID), SiteID: optionalText(filter.SiteID), Status: optionalText(filter.Status), Limit: limit})
	if err != nil {
		return nil, err
	}
	out := make([]AdminDeployment, 0, len(rows))
	for _, row := range rows {
		dep, err := adminDeploymentFromRow(adminDeploymentRow{ID: row.ID, SiteID: row.SiteID, SiteSlug: row.SiteSlug, PrimaryHost: row.PrimaryHost, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Version: row.Version, Status: row.Status, BundleRef: row.BundleRef, UnpackedPath: row.UnpackedPath, ManifestJson: row.ManifestJson, ValidationJson: row.ValidationJson, CreatedByType: row.CreatedByType, CreatedByID: row.CreatedByID, CreatedAt: row.CreatedAt, ActivatedAt: row.ActivatedAt, BundleSha256: row.BundleSha256})
		if err != nil {
			return nil, err
		}
		out = append(out, *dep)
	}
	return out, nil
}

type adminDeploymentRow struct {
	ID             string
	SiteID         string
	SiteSlug       string
	PrimaryHost    string
	OrgID          string
	OrgSlug        string
	OrgName        string
	Version        int32
	Status         string
	BundleRef      string
	UnpackedPath   string
	ManifestJson   []byte
	ValidationJson []byte
	CreatedByType  string
	CreatedByID    string
	CreatedAt      pgtype.Timestamptz
	ActivatedAt    pgtype.Timestamptz
	BundleSha256   string
}

func adminDeploymentFromRow(row adminDeploymentRow) (*AdminDeployment, error) {
	return &AdminDeployment{ID: row.ID, SiteID: row.SiteID, SiteSlug: row.SiteSlug, PrimaryHost: row.PrimaryHost, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Version: int(row.Version), Status: row.Status, BundleRef: row.BundleRef, UnpackedPath: row.UnpackedPath, ManifestJSON: row.ManifestJson, ValidationJSON: row.ValidationJson, CreatedByType: row.CreatedByType, CreatedByID: row.CreatedByID, CreatedAt: formatDBTime(row.CreatedAt), ActivatedAt: formatDBTime(row.ActivatedAt), BundleSHA256: row.BundleSha256}, nil
}

func optionalText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

func formatDBTime(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.UTC().Format(timeFormat)
}
