package store

import (
	"context"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

type RuntimeEvent struct {
	ID           string
	SiteID       string
	OrgID        string
	DeploymentID string
	EventType    string
	Status       string
	Message      string
	MetadataJSON string
	CreatedAt    string
}

func (s *Store) InsertRuntimeEvent(ctx context.Context, event RuntimeEvent) (*RuntimeEvent, error) {
	if event.MetadataJSON == "" {
		event.MetadataJSON = "{}"
	}
	row, err := s.q.InsertRuntimeEvent(ctx, storedb.InsertRuntimeEventParams{ID: newID("rte"), SiteID: event.SiteID, OrgID: event.OrgID, DeploymentID: event.DeploymentID, EventType: event.EventType, Status: event.Status, Message: event.Message, MetadataJson: []byte(event.MetadataJSON), CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return runtimeEventFromRow(row), nil
}

func (s *Store) ListRuntimeEventsBySite(ctx context.Context, siteID string, limit int32) ([]RuntimeEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.q.ListRuntimeEventsBySite(ctx, storedb.ListRuntimeEventsBySiteParams{SiteID: siteID, Limit: limit})
	if err != nil {
		return nil, err
	}
	out := make([]RuntimeEvent, 0, len(rows))
	for _, row := range rows {
		out = append(out, *runtimeEventFromListRow(row))
	}
	return out, nil
}

func runtimeEventFromRow(row storedb.RuntimeEvent) *RuntimeEvent {
	return &RuntimeEvent{ID: row.ID, SiteID: row.SiteID, OrgID: row.OrgID, DeploymentID: row.DeploymentID, EventType: row.EventType, Status: row.Status, Message: row.Message, MetadataJSON: string(row.MetadataJson), CreatedAt: fromPgTime(row.CreatedAt).Format(timeFormat)}
}

func runtimeEventFromListRow(row storedb.RuntimeEvent) *RuntimeEvent {
	return &RuntimeEvent{ID: row.ID, SiteID: row.SiteID, OrgID: row.OrgID, DeploymentID: row.DeploymentID, EventType: row.EventType, Status: row.Status, Message: row.Message, MetadataJSON: string(row.MetadataJson), CreatedAt: fromPgTime(row.CreatedAt).Format(timeFormat)}
}
