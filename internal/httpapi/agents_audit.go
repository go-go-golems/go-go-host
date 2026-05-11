package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/go-go-golems/go-go-host/internal/store"
)

type createAgentRequest struct {
	Name string `json:"name"`
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
		agent, err := core.Agents.Create(r.Context(), p.User.ID, r.PathValue("org_id"), req.Name)
		if err != nil {
			writeDeploymentError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, agentToDTO(*agent))
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

func agentToDTO(agent store.Agent) agentDTO {
	return agentDTO{ID: agent.ID, OrgID: agent.OrgID, Name: agent.Name, Status: agent.Status, CreatedByUserID: agent.CreatedByUserID, CreatedAt: agent.CreatedAt.Format(time.RFC3339), LastSeenAt: agent.LastSeenAt.Format(time.RFC3339)}
}

func auditToDTO(event store.AuditEvent) auditDTO {
	return auditDTO{ID: event.ID, OrgID: event.OrgID, ActorType: event.ActorType, ActorID: event.ActorID, Action: event.Action, ResourceType: event.ResourceType, ResourceID: event.ResourceID, IPAddress: event.IPAddress, UserAgent: event.UserAgent, MetadataJSON: event.MetadataJSON, CreatedAt: event.CreatedAt.Format(time.RFC3339)}
}
