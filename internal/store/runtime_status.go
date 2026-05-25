package store

import (
	"context"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

const staleRuntimeReconcileMessage = "reconciled on daemon startup: runtime state does not survive process restarts"

func (s *Store) UpsertRuntimeStatus(ctx context.Context, st RuntimeStatus) error {
	if st.UpdatedAt.IsZero() {
		st.UpdatedAt = now()
	}
	return s.q.UpsertRuntimeStatus(ctx, storedb.UpsertRuntimeStatusParams{
		SiteID:        st.SiteID,
		OrgID:         st.OrgID,
		DeploymentID:  st.DeploymentID,
		Hosts:         st.Hosts,
		Status:        st.Status,
		StartedAt:     pgTime(st.StartedAt),
		LastError:     st.LastError,
		RequestsTotal: st.RequestsTotal,
		ErrorsTotal:   st.ErrorsTotal,
		UpdatedAt:     pgTime(st.UpdatedAt),
	})
}

func (s *Store) GetRuntimeStatus(ctx context.Context, siteID string) (*RuntimeStatus, error) {
	row, err := s.q.GetRuntimeStatus(ctx, siteID)
	if err != nil {
		return nil, err
	}
	return runtimeStatusFromDB(row), nil
}

func (s *Store) ListRuntimeStatuses(ctx context.Context) ([]RuntimeStatus, error) {
	rows, err := s.q.ListRuntimeStatuses(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]RuntimeStatus, 0, len(rows))
	for _, row := range rows {
		out = append(out, *runtimeStatusFromDB(row))
	}
	return out, nil
}

func (s *Store) ReconcileStaleRuntimeStatuses(ctx context.Context) error {
	return s.q.ReconcileStaleRuntimeStatuses(ctx, storedb.ReconcileStaleRuntimeStatusesParams{LastError: staleRuntimeReconcileMessage, UpdatedAt: pgTime(now())})
}

func runtimeStatusFromDB(row storedb.RuntimeStatus) *RuntimeStatus {
	return &RuntimeStatus{
		SiteID:        row.SiteID,
		OrgID:         row.OrgID,
		DeploymentID:  row.DeploymentID,
		Hosts:         row.Hosts,
		Status:        row.Status,
		StartedAt:     fromPgTime(row.StartedAt),
		LastError:     row.LastError,
		RequestsTotal: row.RequestsTotal,
		ErrorsTotal:   row.ErrorsTotal,
		UpdatedAt:     fromPgTime(row.UpdatedAt),
	}
}
