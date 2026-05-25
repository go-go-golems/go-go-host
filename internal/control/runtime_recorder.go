package control

import (
	"context"
	"math"

	hostruntime "github.com/go-go-golems/go-go-host/internal/runtime"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type runtimeStatusRecorder struct{ store *store.Store }

func (r runtimeStatusRecorder) RecordRuntimeStatus(ctx context.Context, status hostruntime.RuntimeStatus) error {
	if r.store == nil {
		return nil
	}
	err := r.store.UpsertRuntimeStatus(ctx, store.RuntimeStatus{
		SiteID:        status.SiteID,
		OrgID:         status.OrgID,
		DeploymentID:  status.DeploymentID,
		Hosts:         status.Hosts,
		Status:        string(status.Status),
		StartedAt:     status.StartedAt,
		LastError:     status.LastError,
		RequestsTotal: boundedUint64ToInt64(status.RequestsTotal),
		ErrorsTotal:   boundedUint64ToInt64(status.ErrorsTotal),
	})
	_, _ = r.store.InsertRuntimeEvent(ctx, store.RuntimeEvent{SiteID: status.SiteID, OrgID: status.OrgID, DeploymentID: status.DeploymentID, EventType: "runtime.status", Status: string(status.Status), Message: status.LastError})
	return err
}

func boundedUint64ToInt64(v uint64) int64 {
	if v > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(v)
}
