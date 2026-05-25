package control

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-host/internal/store"
)

const agentSignatureSkew = 5 * time.Minute

type AgentCreateResult struct {
	Agent           *store.Agent
	EnrollmentToken string
	Grant           *store.AgentSiteGrant
}

type CreateAgentWithTokenInput struct {
	ActorUserID     string
	OrgID           string
	Name            string
	SiteID          string
	AllowedChannels []string
	AllowedPaths    []string
	CanActivate     bool
	ExpiresAt       time.Time
}

type EnrollAgentInput struct {
	Token     string
	PublicKey string
}

type EnrollAgentResult struct {
	Agent *store.Agent
	Key   *store.AgentKey
	Grant *store.AgentSiteGrant
}

type SignedAgentRequest struct {
	AgentID    string
	KeyID      string
	Method     string
	PathQuery  string
	BodySHA256 string
	Timestamp  time.Time
	Nonce      string
	Signature  string
}

type CreateDeployRunInput struct {
	AgentID  string
	SiteID   string
	Channel  string
	Path     string
	Action   string
	Activate bool
}

type CreateDeployRunResult struct {
	Run         *store.DeployRun
	UploadToken string
}

func (s *AgentService) CreateWithEnrollmentToken(ctx context.Context, input CreateAgentWithTokenInput) (*AgentCreateResult, error) {
	if input.ExpiresAt.IsZero() {
		input.ExpiresAt = time.Now().UTC().Add(24 * time.Hour)
	}
	agent, err := s.Create(ctx, input.ActorUserID, input.OrgID, input.Name)
	if err != nil {
		return nil, err
	}
	token, tokenHash, err := newSecret("enroll")
	if err != nil {
		return nil, err
	}
	if err := s.store.CreateAgentEnrollmentToken(ctx, tokenHash, agent.ID, input.OrgID, input.ExpiresAt); err != nil {
		return nil, err
	}
	var grant *store.AgentSiteGrant
	if input.SiteID != "" {
		grant, err = s.UpsertGrant(ctx, input.ActorUserID, input.OrgID, store.UpsertAgentGrantInput{AgentID: agent.ID, SiteID: input.SiteID, CanDeploy: true, CanRollback: false, CanActivate: input.CanActivate, AllowedChannels: defaultStrings(input.AllowedChannels, []string{"default"}), AllowedPaths: defaultStrings(input.AllowedPaths, []string{"**"}), ExpiresAt: input.ExpiresAt})
		if err != nil {
			return nil, err
		}
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: input.OrgID, ActorType: "user", ActorID: input.ActorUserID, Action: "agent.enrollment_token.create", ResourceType: "agent", ResourceID: agent.ID})
	return &AgentCreateResult{Agent: agent, EnrollmentToken: token, Grant: grant}, nil
}

func (s *AgentService) UpsertGrant(ctx context.Context, actorUserID, orgID string, input store.UpsertAgentGrantInput) (*store.AgentSiteGrant, error) {
	if err := ensureDeployRole(ctx, s.store, actorUserID, orgID); err != nil {
		return nil, err
	}
	if input.CanActivate {
		if err := ensureOwnerRole(ctx, s.store, actorUserID, orgID); err != nil {
			return nil, err
		}
	}
	agent, err := s.store.GetAgent(ctx, input.AgentID)
	if err != nil {
		return nil, err
	}
	if agent.OrgID != orgID {
		return nil, ErrPermissionDenied
	}
	site, err := s.store.GetSite(ctx, input.SiteID)
	if err != nil {
		return nil, err
	}
	if site.OrgID != orgID {
		return nil, ErrPermissionDenied
	}
	grant, err := s.store.UpsertAgentSiteGrant(ctx, input)
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "user", ActorID: actorUserID, Action: "agent.grant.upsert", ResourceType: "agent", ResourceID: input.AgentID})
	return grant, nil
}

func (s *AgentService) Enroll(ctx context.Context, input EnrollAgentInput) (*EnrollAgentResult, error) {
	if strings.TrimSpace(input.Token) == "" || strings.TrimSpace(input.PublicKey) == "" {
		return nil, fmt.Errorf("token and publicKey are required")
	}
	tokenHash := hashSecret(input.Token)
	tok, err := s.store.GetAgentEnrollmentToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if tok.Status != store.AgentEnrollmentTokenStatusActive || time.Now().UTC().After(tok.ExpiresAt) {
		return nil, fmt.Errorf("enrollment token is expired or already used")
	}
	agent, err := s.store.GetAgent(ctx, tok.AgentID)
	if err != nil {
		return nil, err
	}
	if agent.Status != store.AgentStatusActive {
		return nil, ErrPermissionDenied
	}
	if _, err := decodePublicKey(input.PublicKey); err != nil {
		return nil, err
	}
	key, err := s.store.CreateAgentKey(ctx, agent.ID, input.PublicKey)
	if err != nil {
		return nil, err
	}
	if err := s.store.MarkAgentEnrollmentTokenUsed(ctx, tokenHash); err != nil {
		return nil, err
	}
	grants, _ := s.store.ListAgentSiteGrants(ctx, agent.ID)
	var grant *store.AgentSiteGrant
	if len(grants) > 0 {
		grant = &grants[0]
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: agent.OrgID, ActorType: "agent", ActorID: agent.ID, Action: "agent.enroll", ResourceType: "agent_key", ResourceID: key.ID})
	return &EnrollAgentResult{Agent: agent, Key: key, Grant: grant}, nil
}

func (s *AgentService) VerifySignedRequest(ctx context.Context, req SignedAgentRequest) (*store.Agent, error) {
	if req.AgentID == "" || req.KeyID == "" || req.Nonce == "" || req.Signature == "" {
		s.auditAgentSecurity(ctx, "agent.signature.missing_headers", req.AgentID, "", "missing required signature headers")
		return nil, fmt.Errorf("missing agent signature headers")
	}
	if d := time.Since(req.Timestamp); d > agentSignatureSkew || d < -agentSignatureSkew {
		s.auditAgentSecurity(ctx, "agent.signature.timestamp_skew", req.AgentID, "", "agent signature timestamp outside allowed skew")
		return nil, fmt.Errorf("agent signature timestamp outside allowed skew")
	}
	agent, err := s.store.GetAgent(ctx, req.AgentID)
	if err != nil {
		return nil, err
	}
	if agent.Status != store.AgentStatusActive {
		s.auditAgentSecurity(ctx, "agent.signature.revoked_agent", req.AgentID, agent.OrgID, "agent is not active")
		return nil, ErrPermissionDenied
	}
	key, err := s.store.GetAgentKey(ctx, req.KeyID)
	if err != nil {
		return nil, err
	}
	if key.AgentID != agent.ID || key.Status != store.AgentKeyStatusActive {
		s.auditAgentSecurity(ctx, "agent.signature.revoked_key", req.AgentID, agent.OrgID, "agent key is not active")
		return nil, ErrPermissionDenied
	}
	pub, err := decodePublicKey(key.PublicKey)
	if err != nil {
		return nil, err
	}
	sig, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		s.auditAgentSecurity(ctx, "agent.signature.invalid", req.AgentID, agent.OrgID, "invalid signature encoding")
		return nil, fmt.Errorf("invalid signature encoding")
	}
	canonical := AgentCanonicalString(req.Method, req.PathQuery, req.BodySHA256, req.Timestamp.Format(time.RFC3339), req.Nonce)
	if !ed25519.Verify(pub, []byte(canonical), sig) {
		s.auditAgentSecurity(ctx, "agent.signature.invalid", req.AgentID, agent.OrgID, "signature verification failed")
		return nil, ErrPermissionDenied
	}
	if err := s.store.InsertAgentNonce(ctx, agent.ID, req.Nonce); err != nil {
		s.auditAgentSecurity(ctx, "agent.signature.nonce_replay", req.AgentID, agent.OrgID, "agent nonce replay detected")
		return nil, fmt.Errorf("agent nonce replay detected")
	}
	_ = s.store.TouchAgentLastSeen(ctx, agent.ID)
	_ = s.store.TouchAgentKeyLastUsed(ctx, key.ID)
	return agent, nil
}

func (s *AgentService) CreateDeployRun(ctx context.Context, input CreateDeployRunInput) (*CreateDeployRunResult, error) {
	agent, err := s.store.GetAgent(ctx, input.AgentID)
	if err != nil {
		return nil, err
	}
	if agent.Status != store.AgentStatusActive {
		return nil, ErrPermissionDenied
	}
	grant, err := s.findGrant(ctx, input.AgentID, input.SiteID)
	if err != nil {
		return nil, err
	}
	if !grant.CanDeploy || grantExpired(grant) || !allowedString(input.Channel, grant.AllowedChannels) || !allowedPath(input.Path, grant.AllowedPaths) {
		s.auditAgentSecurity(ctx, "agent.grant.denied", input.AgentID, agent.OrgID, "deploy grant denied")
		return nil, ErrPermissionDenied
	}
	requestedAction := defaultString(input.Action, "deploy")
	wantsActivation := input.Activate || requestedAction == "activate"
	if wantsActivation && !grant.CanActivate {
		s.auditAgentSecurity(ctx, "agent.grant.denied", input.AgentID, agent.OrgID, "activation grant denied")
		return nil, ErrPermissionDenied
	}
	token, tokenHash, err := newSecret("upload")
	if err != nil {
		return nil, err
	}
	expires := time.Now().UTC().Add(30 * time.Minute)
	if !grant.ExpiresAt.IsZero() && grant.ExpiresAt.Before(expires) {
		expires = grant.ExpiresAt
	}
	actions := []string{requestedAction}
	if input.Activate && requestedAction != "activate" {
		actions = append(actions, "activate")
	}
	run, err := s.store.CreateDeployRun(ctx, store.CreateDeployRunInput{AgentID: input.AgentID, SiteID: input.SiteID, AllowedActions: actions, AllowedChannels: []string{defaultString(input.Channel, "default")}, AllowedPaths: grant.AllowedPaths, UploadTokenHash: tokenHash, ExpiresAt: expires})
	if err != nil {
		return nil, err
	}
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: agent.OrgID, ActorType: "agent", ActorID: agent.ID, Action: "deploy_run.create", ResourceType: "deploy_run", ResourceID: run.ID})
	return &CreateDeployRunResult{Run: run, UploadToken: token}, nil
}

func (s *AgentService) ValidateUploadToken(ctx context.Context, runID, token string) (*store.DeployRun, *store.Agent, error) {
	run, err := s.store.GetDeployRun(ctx, runID)
	if err != nil {
		return nil, nil, err
	}
	if run.Status != store.DeployRunStatusPending || time.Now().UTC().After(run.ExpiresAt) || hashSecret(token) != run.UploadTokenHash {
		s.auditAgentSecurity(ctx, "agent.upload_token.invalid", run.AgentID, "", "upload token invalid, expired, or already used")
		return nil, nil, ErrPermissionDenied
	}
	run, err = s.store.BeginDeployRunUpload(ctx, runID)
	if err != nil {
		return nil, nil, ErrPermissionDenied
	}
	agent, err := s.store.GetAgent(ctx, run.AgentID)
	if err != nil {
		return nil, nil, err
	}
	if agent.Status != store.AgentStatusActive {
		return nil, nil, ErrPermissionDenied
	}
	return run, agent, nil
}

func (s *AgentService) findGrant(ctx context.Context, agentID, siteID string) (*store.AgentSiteGrant, error) {
	grants, err := s.store.ListAgentSiteGrants(ctx, agentID)
	if err != nil {
		return nil, err
	}
	for _, grant := range grants {
		if grant.SiteID == siteID {
			return &grant, nil
		}
	}
	return nil, ErrPermissionDenied
}

func AgentCanonicalString(method, pathQuery, bodySHA256, timestamp, nonce string) string {
	return strings.ToUpper(method) + "\n" + pathQuery + "\n" + bodySHA256 + "\n" + timestamp + "\n" + nonce
}

func HashBody(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func hashSecret(secret string) string { return HashBody([]byte(secret)) }

func newSecret(prefix string) (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	secret := prefix + "_" + base64.RawURLEncoding.EncodeToString(buf)
	return secret, hashSecret(secret), nil
}

func decodePublicKey(s string) (ed25519.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("publicKey must be base64 encoded")
	}
	if len(b) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("publicKey must decode to %d bytes", ed25519.PublicKeySize)
	}
	return ed25519.PublicKey(b), nil
}

func defaultStrings(v, fallback []string) []string {
	if len(v) == 0 {
		return fallback
	}
	return v
}

func defaultString(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func grantExpired(grant *store.AgentSiteGrant) bool {
	return !grant.ExpiresAt.IsZero() && time.Now().UTC().After(grant.ExpiresAt)
}

func allowedString(v string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	v = defaultString(v, "default")
	for _, a := range allowed {
		if a == "*" || a == v {
			return true
		}
	}
	return false
}

func allowedPath(v string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	v = filepath.ToSlash(v)
	for _, a := range allowed {
		a = filepath.ToSlash(a)
		if a == "" || a == "*" || a == "**" || a == v {
			return true
		}
		if strings.HasSuffix(a, "/**") && strings.HasPrefix(v, strings.TrimSuffix(a, "**")) {
			return true
		}
		if ok, err := filepath.Match(a, v); err == nil && ok {
			return true
		}
	}
	return false
}

func (s *AgentService) auditAgentSecurity(ctx context.Context, action, agentID, orgID, message string) {
	if s == nil || s.store == nil {
		return
	}
	metadata, _ := json.Marshal(map[string]string{"message": message})
	_, _ = s.store.InsertAuditEvent(ctx, store.AuditEvent{OrgID: orgID, ActorType: "agent", ActorID: agentID, Action: action, ResourceType: "agent", ResourceID: agentID, MetadataJSON: string(metadata)})
}

var ErrAgentSignature = errors.New("invalid agent signature")
