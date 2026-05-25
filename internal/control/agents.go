package control

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-go-golems/go-go-host/internal/store"
)

type AgentService struct{ store *store.Store }

type AuditService struct{ store *store.Store }

func (s *AgentService) Create(ctx context.Context, actorUserID, orgID, name string) (*store.Agent, error) {
	if s.store == nil {
		return nil, errors.New("store is not configured")
	}
	if err := ensureDeployRole(ctx, s.store, actorUserID, orgID); err != nil {
		return nil, err
	}
	agent, err := s.store.CreateAgent(ctx, store.CreateAgentInput{OrgID: orgID, Name: name, CreatedByUserID: actorUserID})
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: actorUserID, Action: "agent.create", ResourceType: "agent", ResourceID: agent.ID})
	return agent, nil
}

func (s *AgentService) List(ctx context.Context, actorUserID, orgID string) ([]store.Agent, error) {
	if err := ensureViewRole(ctx, s.store, actorUserID, orgID); err != nil {
		return nil, err
	}
	return s.store.ListAgentsByOrg(ctx, orgID)
}

func (s *AgentService) CreateEnrollmentTokenForAgent(ctx context.Context, actorUserID, orgID, agentID string, expiresAt time.Time) (string, error) {
	if err := ensureDeployRole(ctx, s.store, actorUserID, orgID); err != nil {
		return "", err
	}
	agent, err := s.store.GetAgent(ctx, agentID)
	if err != nil {
		return "", err
	}
	if agent.OrgID != orgID || agent.Status != store.AgentStatusActive {
		return "", ErrPermissionDenied
	}
	if expiresAt.IsZero() {
		expiresAt = time.Now().UTC().Add(24 * time.Hour)
	}
	token, tokenHash, err := newSecret("enroll")
	if err != nil {
		return "", err
	}
	if err := s.store.CreateAgentEnrollmentToken(ctx, tokenHash, agent.ID, orgID, expiresAt); err != nil {
		return "", err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: actorUserID, Action: "agent.enrollment_token.create", ResourceType: "agent", ResourceID: agent.ID})
	return token, nil
}

func (s *AgentService) ListKeys(ctx context.Context, actorUserID, orgID, agentID string) ([]store.AgentKey, error) {
	if err := ensureViewRole(ctx, s.store, actorUserID, orgID); err != nil {
		return nil, err
	}
	agent, err := s.store.GetAgent(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if agent.OrgID != orgID {
		return nil, ErrPermissionDenied
	}
	return s.store.ListAgentKeys(ctx, agentID)
}

func (s *AgentService) RevokeKey(ctx context.Context, actorUserID, orgID, agentID, keyID, reason string) error {
	if err := ensureDeployRole(ctx, s.store, actorUserID, orgID); err != nil {
		return err
	}
	agent, err := s.store.GetAgent(ctx, agentID)
	if err != nil {
		return err
	}
	if agent.OrgID != orgID {
		return ErrPermissionDenied
	}
	key, err := s.store.GetAgentKey(ctx, keyID)
	if err != nil {
		return err
	}
	if key.AgentID != agentID {
		return ErrPermissionDenied
	}
	if err := s.store.RevokeAgentKey(ctx, keyID); err != nil {
		return err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: actorUserID, Action: "agent.key.revoke", ResourceType: "agent_key", ResourceID: keyID, MetadataJSON: fmt.Sprintf(`{"reason":%q}`, reason)})
	return nil
}

func (s *AgentService) Revoke(ctx context.Context, actorUserID, orgID, agentID string) error {
	if err := ensureDeployRole(ctx, s.store, actorUserID, orgID); err != nil {
		return err
	}
	agent, err := s.store.GetAgent(ctx, agentID)
	if err != nil {
		return err
	}
	if agent.OrgID != orgID {
		return ErrPermissionDenied
	}
	if err := s.store.UpdateAgentStatus(ctx, agentID, store.AgentStatusRevoked); err != nil {
		return err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: actorUserID, Action: "agent.revoke", ResourceType: "agent", ResourceID: agentID})
	return nil
}

func (s *AuditService) List(ctx context.Context, actorUserID string, filter store.AuditFilter) ([]store.AuditEvent, error) {
	if err := ensureViewRole(ctx, s.store, actorUserID, filter.OrgID); err != nil {
		return nil, err
	}
	return s.store.ListAuditEvents(ctx, filter)
}
