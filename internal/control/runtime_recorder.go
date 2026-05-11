package control

import (
	"context"

	hostruntime "github.com/go-go-golems/go-go-host/internal/runtime"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type runtimeStatusRecorder struct{ store *store.Store }

func (r runtimeStatusRecorder) RecordRuntimeStatus(ctx context.Context, status hostruntime.RuntimeStatus) error {
	if r.store == nil {
		return nil
	}
	return r.store.UpsertRuntimeStatus(ctx, store.RuntimeStatus{
		SiteID:        status.SiteID,
		OrgID:         status.OrgID,
		DeploymentID:  status.DeploymentID,
		Hosts:         status.Hosts,
		Status:        string(status.Status),
		StartedAt:     status.StartedAt,
		LastError:     status.LastError,
		RequestsTotal: int64(status.RequestsTotal),
		ErrorsTotal:   int64(status.ErrorsTotal),
	})
}
