package store

import (
	"context"
	"time"

	storedb "github.com/go-go-golems/go-go-host/internal/store/db"
)

type CreateAgentInput struct {
	OrgID           string
	Name            string
	CreatedByUserID string
}

type UpsertAgentGrantInput struct {
	AgentID         string
	SiteID          string
	CanDeploy       bool
	CanRollback     bool
	AllowedChannels []string
	AllowedPaths    []string
	ExpiresAt       time.Time
}

func (s *Store) CreateAgent(ctx context.Context, input CreateAgentInput) (*Agent, error) {
	row, err := s.q.CreateAgent(ctx, storedb.CreateAgentParams{ID: newID("agt"), OrgID: input.OrgID, Name: input.Name, Status: AgentStatusActive, CreatedByUserID: input.CreatedByUserID, CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return agentFromDB(row), nil
}

func (s *Store) GetAgent(ctx context.Context, id string) (*Agent, error) {
	row, err := s.q.GetAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	return agentFromDB(row), nil
}

func (s *Store) ListAgentsByOrg(ctx context.Context, orgID string) ([]Agent, error) {
	rows, err := s.q.ListAgentsByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	out := make([]Agent, 0, len(rows))
	for _, row := range rows {
		out = append(out, *agentFromDB(row))
	}
	return out, nil
}

func (s *Store) UpdateAgentStatus(ctx context.Context, id, status string) error {
	return s.q.UpdateAgentStatus(ctx, storedb.UpdateAgentStatusParams{ID: id, Status: status})
}

func (s *Store) UpsertAgentSiteGrant(ctx context.Context, input UpsertAgentGrantInput) (*AgentSiteGrant, error) {
	row, err := s.q.UpsertAgentSiteGrant(ctx, storedb.UpsertAgentSiteGrantParams{AgentID: input.AgentID, SiteID: input.SiteID, CanDeploy: input.CanDeploy, CanRollback: input.CanRollback, AllowedChannels: input.AllowedChannels, AllowedPaths: input.AllowedPaths, ExpiresAt: pgTime(input.ExpiresAt), CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return agentGrantFromDB(row), nil
}

func (s *Store) ListAgentSiteGrants(ctx context.Context, agentID string) ([]AgentSiteGrant, error) {
	rows, err := s.q.ListAgentSiteGrants(ctx, agentID)
	if err != nil {
		return nil, err
	}
	out := make([]AgentSiteGrant, 0, len(rows))
	for _, row := range rows {
		out = append(out, *agentGrantFromDB(row))
	}
	return out, nil
}

func agentFromDB(row storedb.Agent) *Agent {
	return &Agent{ID: row.ID, OrgID: row.OrgID, Name: row.Name, Status: row.Status, CreatedByUserID: row.CreatedByUserID, CreatedAt: fromPgTime(row.CreatedAt), LastSeenAt: fromPgTime(row.LastSeenAt)}
}

func agentGrantFromDB(row storedb.AgentSiteGrant) *AgentSiteGrant {
	return &AgentSiteGrant{AgentID: row.AgentID, SiteID: row.SiteID, CanDeploy: row.CanDeploy, CanRollback: row.CanRollback, AllowedChannels: row.AllowedChannels, AllowedPaths: row.AllowedPaths, ExpiresAt: fromPgTime(row.ExpiresAt), CreatedAt: fromPgTime(row.CreatedAt), UpdatedAt: fromPgTime(row.UpdatedAt)}
}
