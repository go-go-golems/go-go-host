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
	CanActivate     bool
	AllowedChannels []string
	AllowedPaths    []string
	ExpiresAt       time.Time
}

type CreateDeployRunInput struct {
	AgentID         string
	SiteID          string
	AllowedActions  []string
	AllowedChannels []string
	AllowedPaths    []string
	UploadTokenHash string
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
	row, err := s.q.UpsertAgentSiteGrant(ctx, storedb.UpsertAgentSiteGrantParams{AgentID: input.AgentID, SiteID: input.SiteID, CanDeploy: input.CanDeploy, CanRollback: input.CanRollback, CanActivate: input.CanActivate, AllowedChannels: input.AllowedChannels, AllowedPaths: input.AllowedPaths, ExpiresAt: pgTime(input.ExpiresAt), CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return agentGrantFromUpsertRow(row), nil
}

func (s *Store) ListAgentSiteGrants(ctx context.Context, agentID string) ([]AgentSiteGrant, error) {
	rows, err := s.q.ListAgentSiteGrants(ctx, agentID)
	if err != nil {
		return nil, err
	}
	out := make([]AgentSiteGrant, 0, len(rows))
	for _, row := range rows {
		out = append(out, *agentGrantFromListRow(row))
	}
	return out, nil
}

func (s *Store) CreateAgentEnrollmentToken(ctx context.Context, tokenHash, agentID, orgID string, expiresAt time.Time) error {
	return s.q.CreateAgentEnrollmentToken(ctx, storedb.CreateAgentEnrollmentTokenParams{TokenHash: tokenHash, AgentID: agentID, OrgID: orgID, Status: AgentEnrollmentTokenStatusActive, ExpiresAt: pgTime(expiresAt), CreatedAt: pgTime(now())})
}

func (s *Store) GetAgentEnrollmentToken(ctx context.Context, tokenHash string) (*AgentEnrollmentToken, error) {
	row, err := s.q.GetAgentEnrollmentToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return enrollmentTokenFromDB(row), nil
}

func (s *Store) MarkAgentEnrollmentTokenUsed(ctx context.Context, tokenHash string) error {
	return s.q.MarkAgentEnrollmentTokenUsed(ctx, storedb.MarkAgentEnrollmentTokenUsedParams{TokenHash: tokenHash, UsedAt: pgTime(now())})
}

func (s *Store) CreateAgentKey(ctx context.Context, agentID, publicKey string) (*AgentKey, error) {
	row, err := s.q.CreateAgentKey(ctx, storedb.CreateAgentKeyParams{ID: newID("ak"), AgentID: agentID, PublicKey: publicKey, Status: AgentKeyStatusActive, CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return agentKeyFromDB(row), nil
}

func (s *Store) GetAgentKey(ctx context.Context, id string) (*AgentKey, error) {
	row, err := s.q.GetAgentKey(ctx, id)
	if err != nil {
		return nil, err
	}
	return agentKeyFromDB(row), nil
}

func (s *Store) ListAgentKeys(ctx context.Context, agentID string) ([]AgentKey, error) {
	rows, err := s.q.ListAgentKeys(ctx, agentID)
	if err != nil {
		return nil, err
	}
	out := make([]AgentKey, 0, len(rows))
	for _, row := range rows {
		out = append(out, *agentKeyFromDB(row))
	}
	return out, nil
}

func (s *Store) RevokeAgentKey(ctx context.Context, id string) error {
	return s.q.RevokeAgentKey(ctx, storedb.RevokeAgentKeyParams{ID: id, RevokedAt: pgTime(now())})
}

func (s *Store) TouchAgentKeyLastUsed(ctx context.Context, keyID string) error {
	return s.q.TouchAgentKeyLastUsed(ctx, storedb.TouchAgentKeyLastUsedParams{ID: keyID, LastUsedAt: pgTime(now())})
}

func (s *Store) TouchAgentLastSeen(ctx context.Context, agentID string) error {
	return s.q.TouchAgentLastSeen(ctx, storedb.TouchAgentLastSeenParams{ID: agentID, LastSeenAt: pgTime(now())})
}

func (s *Store) InsertAgentNonce(ctx context.Context, agentID, nonce string) error {
	return s.q.InsertAgentNonce(ctx, storedb.InsertAgentNonceParams{AgentID: agentID, Nonce: nonce, SeenAt: pgTime(now())})
}

func (s *Store) CreateDeployRun(ctx context.Context, input CreateDeployRunInput) (*DeployRun, error) {
	row, err := s.q.CreateDeployRun(ctx, storedb.CreateDeployRunParams{ID: newID("dr"), SiteID: input.SiteID, ActorID: input.AgentID, Status: DeployRunStatusPending, AllowedActions: input.AllowedActions, AllowedChannels: input.AllowedChannels, AllowedPaths: input.AllowedPaths, UploadTokenHash: input.UploadTokenHash, ExpiresAt: pgTime(input.ExpiresAt), CreatedAt: pgTime(now())})
	if err != nil {
		return nil, err
	}
	return deployRunFromDB(row), nil
}

func (s *Store) GetDeployRun(ctx context.Context, id string) (*DeployRun, error) {
	row, err := s.q.GetDeployRun(ctx, id)
	if err != nil {
		return nil, err
	}
	return deployRunFromDB(row), nil
}

func (s *Store) FinishDeployRun(ctx context.Context, id, status string) error {
	return s.q.FinishDeployRun(ctx, storedb.FinishDeployRunParams{ID: id, Status: status, FinishedAt: pgTime(now())})
}

func agentFromDB(row storedb.Agent) *Agent {
	return &Agent{ID: row.ID, OrgID: row.OrgID, Name: row.Name, Status: row.Status, CreatedByUserID: row.CreatedByUserID, CreatedAt: fromPgTime(row.CreatedAt), LastSeenAt: fromPgTime(row.LastSeenAt)}
}

func agentGrantFromListRow(row storedb.ListAgentSiteGrantsRow) *AgentSiteGrant {
	return &AgentSiteGrant{AgentID: row.AgentID, SiteID: row.SiteID, CanDeploy: row.CanDeploy, CanRollback: row.CanRollback, CanActivate: row.CanActivate, AllowedChannels: row.AllowedChannels, AllowedPaths: row.AllowedPaths, ExpiresAt: fromPgTime(row.ExpiresAt), CreatedAt: fromPgTime(row.CreatedAt), UpdatedAt: fromPgTime(row.UpdatedAt)}
}

func agentGrantFromUpsertRow(row storedb.UpsertAgentSiteGrantRow) *AgentSiteGrant {
	return &AgentSiteGrant{AgentID: row.AgentID, SiteID: row.SiteID, CanDeploy: row.CanDeploy, CanRollback: row.CanRollback, CanActivate: row.CanActivate, AllowedChannels: row.AllowedChannels, AllowedPaths: row.AllowedPaths, ExpiresAt: fromPgTime(row.ExpiresAt), CreatedAt: fromPgTime(row.CreatedAt), UpdatedAt: fromPgTime(row.UpdatedAt)}
}

func enrollmentTokenFromDB(row storedb.AgentEnrollmentToken) *AgentEnrollmentToken {
	return &AgentEnrollmentToken{TokenHash: row.TokenHash, AgentID: row.AgentID, OrgID: row.OrgID, Status: row.Status, ExpiresAt: fromPgTime(row.ExpiresAt), CreatedAt: fromPgTime(row.CreatedAt), UsedAt: fromPgTime(row.UsedAt)}
}

func agentKeyFromDB(row storedb.AgentKey) *AgentKey {
	return &AgentKey{ID: row.ID, AgentID: row.AgentID, PublicKey: row.PublicKey, Status: row.Status, CreatedAt: fromPgTime(row.CreatedAt), RevokedAt: fromPgTime(row.RevokedAt), LastUsedAt: fromPgTime(row.LastUsedAt)}
}

func deployRunFromDB(row storedb.DeployRun) *DeployRun {
	return &DeployRun{ID: row.ID, SiteID: row.SiteID, ActorType: row.ActorType, ActorID: row.ActorID, AgentID: row.AgentID, RequestedByUserID: row.RequestedByUserID, Status: row.Status, AllowedActions: row.AllowedActions, AllowedChannels: row.AllowedChannels, AllowedPaths: row.AllowedPaths, UploadTokenHash: row.UploadTokenHash, ExpiresAt: fromPgTime(row.ExpiresAt), CreatedAt: fromPgTime(row.CreatedAt), FinishedAt: fromPgTime(row.FinishedAt)}
}
