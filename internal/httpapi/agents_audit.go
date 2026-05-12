package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type createAgentRequest struct {
	Name            string   `json:"name"`
	SiteID          string   `json:"siteId"`
	AllowedChannels []string `json:"allowedChannels"`
	AllowedPaths    []string `json:"allowedPaths"`
	CanActivate     bool     `json:"canActivate"`
}

type upsertAgentGrantRequest struct {
	SiteID          string   `json:"siteId"`
	CanDeploy       bool     `json:"canDeploy"`
	CanRollback     bool     `json:"canRollback"`
	CanActivate     bool     `json:"canActivate"`
	AllowedChannels []string `json:"allowedChannels"`
	AllowedPaths    []string `json:"allowedPaths"`
	ExpiresAt       string   `json:"expiresAt"`
}

type createAgentResponse struct {
	Agent           agentDTO       `json:"agent"`
	EnrollmentToken string         `json:"enrollmentToken,omitempty"`
	Grant           *agentGrantDTO `json:"grant,omitempty"`
}

type agentDTO struct {
	ID              string `json:"id"`
	OrgID           string `json:"orgId"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	CreatedByUserID string `json:"createdByUserId"`
	CreatedAt       string `json:"createdAt"`
	LastSeenAt      string `json:"lastSeenAt,omitempty"`
}

type agentGrantDTO struct {
	AgentID         string   `json:"agentId"`
	SiteID          string   `json:"siteId"`
	CanDeploy       bool     `json:"canDeploy"`
	CanRollback     bool     `json:"canRollback"`
	CanActivate     bool     `json:"canActivate"`
	AllowedChannels []string `json:"allowedChannels"`
	AllowedPaths    []string `json:"allowedPaths"`
	ExpiresAt       string   `json:"expiresAt,omitempty"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

type enrollAgentRequest struct {
	Token     string `json:"token"`
	PublicKey string `json:"publicKey"`
}

type enrollAgentResponse struct {
	Agent agentDTO       `json:"agent"`
	KeyID string         `json:"keyId"`
	Grant *agentGrantDTO `json:"grant,omitempty"`
}

type createDeployRunRequest struct {
	SiteID   string `json:"siteId"`
	Channel  string `json:"channel"`
	Path     string `json:"path"`
	Action   string `json:"action"`
	Activate bool   `json:"activate"`
}

type createDeployRunResponse struct {
	ID           string   `json:"id"`
	SiteID       string   `json:"siteId"`
	Status       string   `json:"status"`
	UploadToken  string   `json:"uploadToken"`
	ExpiresAt    string   `json:"expiresAt"`
	AllowedPaths []string `json:"allowedPaths"`
}

type auditDTO struct {
	ID           string `json:"id"`
	OrgID        string `json:"orgId"`
	ActorType    string `json:"actorType"`
	ActorID      string `json:"actorId"`
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceId"`
	IPAddress    string `json:"ipAddress"`
	UserAgent    string `json:"userAgent"`
	MetadataJSON string `json:"metadataJson"`
	CreatedAt    string `json:"createdAt"`
}

func handleListAgents(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		agents, err := core.Agents.List(r.Context(), p.User.ID, r.PathValue("org_id"))
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		out := make([]agentDTO, 0, len(agents))
		for _, agent := range agents {
			out = append(out, agentToDTO(agent))
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func handleCreateAgent(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req createAgentRequest
		if err := decodeJSONBody(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := core.Agents.CreateWithEnrollmentToken(r.Context(), control.CreateAgentWithTokenInput{ActorUserID: p.User.ID, OrgID: r.PathValue("org_id"), Name: req.Name, SiteID: req.SiteID, AllowedChannels: req.AllowedChannels, AllowedPaths: req.AllowedPaths, CanActivate: req.CanActivate})
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, createAgentResponse{Agent: agentToDTO(*result.Agent), EnrollmentToken: result.EnrollmentToken, Grant: grantToDTO(result.Grant)})
	}
}

func handleRevokeAgent(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if err := core.Agents.Revoke(r.Context(), p.User.ID, r.PathValue("org_id"), r.PathValue("agent_id")); err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "revoked", "agentId": r.PathValue("agent_id")})
	}
}

func handleUpsertAgentGrant(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		var req upsertAgentGrantRequest
		if err := decodeJSONBody(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		var expires time.Time
		if strings.TrimSpace(req.ExpiresAt) != "" {
			expires, err = time.Parse(time.RFC3339, req.ExpiresAt)
			if err != nil {
				writeError(w, http.StatusBadRequest, "expiresAt must be RFC3339")
				return
			}
		}
		grant, err := core.Agents.UpsertGrant(r.Context(), p.User.ID, r.PathValue("org_id"), store.UpsertAgentGrantInput{AgentID: r.PathValue("agent_id"), SiteID: req.SiteID, CanDeploy: req.CanDeploy, CanRollback: req.CanRollback, CanActivate: req.CanActivate, AllowedChannels: req.AllowedChannels, AllowedPaths: req.AllowedPaths, ExpiresAt: expires})
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, grantToDTO(grant))
	}
}

func handleEnrollAgent(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req enrollAgentRequest
		if err := decodeJSONBody(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := core.Agents.Enroll(r.Context(), control.EnrollAgentInput{Token: req.Token, PublicKey: req.PublicKey})
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, enrollAgentResponse{Agent: agentToDTO(*result.Agent), KeyID: result.Key.ID, Grant: grantToDTO(result.Grant)})
	}
}

func handleCreateAgentDeployRun(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		var req createDeployRunRequest
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		agent, err := verifyAgentRequest(core, r, body)
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		result, err := core.Agents.CreateDeployRun(r.Context(), control.CreateDeployRunInput{AgentID: agent.ID, SiteID: req.SiteID, Channel: req.Channel, Path: req.Path, Action: req.Action, Activate: req.Activate})
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, createDeployRunResponse{ID: result.Run.ID, SiteID: result.Run.SiteID, Status: result.Run.Status, UploadToken: result.UploadToken, ExpiresAt: result.Run.ExpiresAt.Format(time.RFC3339), AllowedPaths: result.Run.AllowedPaths})
	}
}

func handleListAudit(core *control.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := requirePrincipal(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		events, err := core.Audit.List(r.Context(), p.User.ID, store.AuditFilter{OrgID: r.PathValue("org_id"), ResourceID: r.URL.Query().Get("resource_id"), ActorType: r.URL.Query().Get("actor_type"), ActorID: r.URL.Query().Get("actor_id"), Action: r.URL.Query().Get("action"), Limit: limit})
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		out := make([]auditDTO, 0, len(events))
		for _, event := range events {
			out = append(out, auditToDTO(event))
		}
		writeJSON(w, http.StatusOK, out)
	}
}

func decodeJSONBody(r *http.Request, out any) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return fmt.Errorf("invalid JSON body")
	}
	return nil
}

func verifyAgentRequest(core *control.Core, r *http.Request, body []byte) (*store.Agent, error) {
	timestamp, err := time.Parse(time.RFC3339, r.Header.Get("X-Go-Go-Agent-Timestamp"))
	if err != nil {
		return nil, fmt.Errorf("invalid X-Go-Go-Agent-Timestamp")
	}
	pathQuery := r.URL.RequestURI()
	return core.Agents.VerifySignedRequest(r.Context(), control.SignedAgentRequest{AgentID: r.Header.Get("X-Go-Go-Agent-ID"), KeyID: r.Header.Get("X-Go-Go-Agent-Key-ID"), Method: r.Method, PathQuery: pathQuery, BodySHA256: control.HashBody(body), Timestamp: timestamp, Nonce: r.Header.Get("X-Go-Go-Agent-Nonce"), Signature: r.Header.Get("X-Go-Go-Agent-Signature")})
}

func agentToDTO(agent store.Agent) agentDTO {
	return agentDTO{ID: agent.ID, OrgID: agent.OrgID, Name: agent.Name, Status: agent.Status, CreatedByUserID: agent.CreatedByUserID, CreatedAt: agent.CreatedAt.Format(time.RFC3339), LastSeenAt: agent.LastSeenAt.Format(time.RFC3339)}
}

func grantToDTO(grant *store.AgentSiteGrant) *agentGrantDTO {
	if grant == nil {
		return nil
	}
	return &agentGrantDTO{AgentID: grant.AgentID, SiteID: grant.SiteID, CanDeploy: grant.CanDeploy, CanRollback: grant.CanRollback, CanActivate: grant.CanActivate, AllowedChannels: grant.AllowedChannels, AllowedPaths: grant.AllowedPaths, ExpiresAt: grant.ExpiresAt.Format(time.RFC3339), CreatedAt: grant.CreatedAt.Format(time.RFC3339), UpdatedAt: grant.UpdatedAt.Format(time.RFC3339)}
}

func auditToDTO(event store.AuditEvent) auditDTO {
	return auditDTO{ID: event.ID, OrgID: event.OrgID, ActorType: event.ActorType, ActorID: event.ActorID, Action: event.Action, ResourceType: event.ResourceType, ResourceID: event.ResourceID, IPAddress: event.IPAddress, UserAgent: event.UserAgent, MetadataJSON: event.MetadataJSON, CreatedAt: event.CreatedAt.Format(time.RFC3339)}
}
