package store

import (
	"context"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

const (
	DeploymentStatusUploaded   = "uploaded"
	DeploymentStatusValidated  = "validated"
	DeploymentStatusRejected   = "rejected"
	DeploymentStatusActive     = "active"
	DeploymentStatusSuperseded = "superseded"
)

type CreateDeploymentInput struct {
	SiteID         string
	Status         string
	BundleRef      string
	UnpackedPath   string
	ManifestJSON   []byte
	ValidationJSON []byte
	CreatedByType  string
	CreatedByID    string
}

type Deployment struct {
	ID             string
	SiteID         string
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

func (s *Store) CreateDeployment(ctx context.Context, input CreateDeploymentInput) (*Deployment, error) {
	version, err := s.q.NextDeploymentVersion(ctx, input.SiteID)
	if err != nil {
		return nil, err
	}
	row, err := s.q.CreateDeployment(ctx, storedb.CreateDeploymentParams{
		ID:             newID("dep"),
		SiteID:         input.SiteID,
		Version:        version,
		Status:         input.Status,
		BundleRef:      input.BundleRef,
		UnpackedPath:   input.UnpackedPath,
		ManifestJson:   input.ManifestJSON,
		ValidationJson: input.ValidationJSON,
		CreatedByType:  input.CreatedByType,
		CreatedByID:    input.CreatedByID,
		CreatedAt:      pgTime(now()),
	})
	if err != nil {
		return nil, err
	}
	return deploymentFromDB(row), nil
}

func (s *Store) GetDeployment(ctx context.Context, id string) (*Deployment, error) {
	row, err := s.q.GetDeployment(ctx, id)
	if err != nil {
		return nil, err
	}
	return deploymentFromDB(row), nil
}

func (s *Store) ListDeploymentsBySite(ctx context.Context, siteID string) ([]Deployment, error) {
	rows, err := s.q.ListDeploymentsBySite(ctx, siteID)
	if err != nil {
		return nil, err
	}
	out := make([]Deployment, 0, len(rows))
	for _, row := range rows {
		out = append(out, *deploymentFromDB(row))
	}
	return out, nil
}

func (s *Store) UpdateDeploymentArtifacts(ctx context.Context, id, status, bundleRef, unpackedPath string, manifestJSON, validationJSON []byte) error {
	return s.q.UpdateDeploymentArtifacts(ctx, storedb.UpdateDeploymentArtifactsParams{ID: id, Status: status, BundleRef: bundleRef, UnpackedPath: unpackedPath, ManifestJson: manifestJSON, ValidationJson: validationJSON})
}

func (s *Store) UpdateDeploymentStatus(ctx context.Context, id, status string, validationJSON []byte) error {
	return s.q.UpdateDeploymentStatus(ctx, storedb.UpdateDeploymentStatusParams{ID: id, Status: status, ValidationJson: validationJSON})
}

func (s *Store) MarkDeploymentActive(ctx context.Context, siteID, deploymentID string) error {
	if err := s.q.SupersedeActiveDeployments(ctx, storedb.SupersedeActiveDeploymentsParams{SiteID: siteID, ID: deploymentID}); err != nil {
		return err
	}
	if err := s.q.ActivateDeployment(ctx, storedb.ActivateDeploymentParams{ID: deploymentID, ActivatedAt: pgTime(now())}); err != nil {
		return err
	}
	return s.UpdateSiteActiveDeployment(ctx, siteID, deploymentID)
}

func (s *Store) PreviousValidatedDeployment(ctx context.Context, siteID, currentDeploymentID string) (*Deployment, error) {
	row, err := s.q.PreviousValidatedDeployment(ctx, storedb.PreviousValidatedDeploymentParams{SiteID: siteID, ID: currentDeploymentID})
	if err != nil {
		return nil, err
	}
	return deploymentFromDB(row), nil
}

func deploymentFromDB(row storedb.Deployment) *Deployment {
	return &Deployment{
		ID:             row.ID,
		SiteID:         row.SiteID,
		Version:        int(row.Version),
		Status:         row.Status,
		BundleRef:      row.BundleRef,
		UnpackedPath:   row.UnpackedPath,
		ManifestJSON:   row.ManifestJson,
		ValidationJSON: row.ValidationJson,
		CreatedByType:  row.CreatedByType,
		CreatedByID:    row.CreatedByID,
		CreatedAt:      fromPgTime(row.CreatedAt).Format(timeFormat),
		ActivatedAt:    fromPgTime(row.ActivatedAt).Format(timeFormat),
	}
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
