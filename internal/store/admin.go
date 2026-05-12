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
}

type AdminDeploymentFilter struct {
	OrgID  string
	SiteID string
	Status string
	Limit  int
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

func (s *Store) ListAdminDeployments(ctx context.Context, filter AdminDeploymentFilter) ([]AdminDeployment, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.q.ListAdminDeployments(ctx, storedb.ListAdminDeploymentsParams{OrgID: optionalText(filter.OrgID), SiteID: optionalText(filter.SiteID), Status: optionalText(filter.Status), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]AdminDeployment, 0, len(rows))
	for _, row := range rows {
		out = append(out, AdminDeployment{ID: row.ID, SiteID: row.SiteID, SiteSlug: row.SiteSlug, PrimaryHost: row.PrimaryHost, OrgID: row.OrgID, OrgSlug: row.OrgSlug, OrgName: row.OrgName, Version: int(row.Version), Status: row.Status, BundleRef: row.BundleRef, UnpackedPath: row.UnpackedPath, ManifestJSON: row.ManifestJson, ValidationJSON: row.ValidationJson, CreatedByType: row.CreatedByType, CreatedByID: row.CreatedByID, CreatedAt: formatDBTime(row.CreatedAt), ActivatedAt: formatDBTime(row.ActivatedAt)})
	}
	return out, nil
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
