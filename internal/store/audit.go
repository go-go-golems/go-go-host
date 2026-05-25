package store

import (
	"context"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

func (s *Store) InsertAuditEvent(ctx context.Context, event AuditEvent) (*AuditEvent, error) {
	if event.ID == "" {
		event.ID = newID("aud")
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now()
	}
	if event.MetadataJSON == "" {
		event.MetadataJSON = "{}"
	}
	row, err := s.q.InsertAuditEvent(ctx, storedb.InsertAuditEventParams{ID: event.ID, OrgID: event.OrgID, ActorType: event.ActorType, ActorID: event.ActorID, Action: event.Action, ResourceType: event.ResourceType, ResourceID: event.ResourceID, IpAddress: event.IPAddress, UserAgent: event.UserAgent, MetadataJson: []byte(event.MetadataJSON), CreatedAt: pgTime(event.CreatedAt)})
	if err != nil {
		return nil, err
	}
	return auditFromDB(row), nil
}

type AuditFilter struct {
	OrgID      string
	ResourceID string
	ActorType  string
	ActorID    string
	Action     string
	Limit      int
}

func (s *Store) ListAuditEventsForOrg(ctx context.Context, orgID string, limit int) ([]AuditEvent, error) {
	boundedLimit := boundedListLimit(limit)
	rows, err := s.q.ListAuditEventsForOrg(ctx, storedb.ListAuditEventsForOrgParams{OrgID: orgID, Limit: boundedLimit})
	if err != nil {
		return nil, err
	}
	events := make([]AuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, *auditFromDB(row))
	}
	return events, nil
}

func (s *Store) ListAuditEvents(ctx context.Context, filter AuditFilter) ([]AuditEvent, error) {
	limit := boundedListLimit(filter.Limit)
	rows, err := s.q.ListAuditEventsFiltered(ctx, storedb.ListAuditEventsFilteredParams{OrgID: filter.OrgID, Column2: filter.ResourceID, Column3: filter.ActorType, Column4: filter.ActorID, Column5: filter.Action, Limit: limit})
	if err != nil {
		return nil, err
	}
	events := make([]AuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, *auditFromDB(row))
	}
	return events, nil
}

func auditFromDB(row storedb.AuditLog) *AuditEvent {
	return &AuditEvent{ID: row.ID, OrgID: row.OrgID, ActorType: row.ActorType, ActorID: row.ActorID, Action: row.Action, ResourceType: row.ResourceType, ResourceID: row.ResourceID, IPAddress: row.IpAddress, UserAgent: row.UserAgent, MetadataJSON: string(row.MetadataJson), CreatedAt: fromPgTime(row.CreatedAt)}
}
