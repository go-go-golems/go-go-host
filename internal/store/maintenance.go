package store

import (
	"context"
	"time"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

func (s *Store) ListPrunableDeployments(ctx context.Context, siteID, activeDeploymentID string, statuses []string, olderThan time.Time) ([]Deployment, error) {
	rows, err := s.q.ListPrunableDeployments(ctx, storedb.ListPrunableDeploymentsParams{SiteID: siteID, ID: activeDeploymentID, Column3: statuses, CreatedAt: pgTime(olderThan)})
	if err != nil {
		return nil, err
	}
	out := make([]Deployment, 0, len(rows))
	for _, row := range rows {
		out = append(out, *deploymentFromDB(row))
	}
	return out, nil
}

func (s *Store) DeleteDeployment(ctx context.Context, id string) error {
	return s.q.DeleteDeployment(ctx, id)
}

func (s *Store) DeleteAuditEventsBefore(ctx context.Context, olderThan time.Time) (int64, error) {
	return s.q.DeleteAuditEventsBefore(ctx, pgTime(olderThan))
}
