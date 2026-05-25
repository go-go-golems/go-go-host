package cmds

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/spf13/cobra"
)

type AuditListCommand struct{ *glazedcmds.CommandDescription }

type AuditSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
	OrgID       string `glazed:"org-id"`
	ResourceID  string `glazed:"resource-id"`
	ActorType   string `glazed:"actor-type"`
	ActorID     string `glazed:"actor-id"`
	Action      string `glazed:"action"`
	Limit       int    `glazed:"limit"`
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

var _ glazedcmds.GlazeCommand = (*AuditListCommand)(nil)

func NewAuditCobraCommand() (*cobra.Command, error) {
	root := &cobra.Command{Use: "audit", Short: "Inspect audit events"}
	list, err := NewAuditListCommand()
	if err != nil {
		return nil, err
	}
	listCmd, err := BuildGlazedCobraCommand(list)
	if err != nil {
		return nil, err
	}
	root.AddCommand(listCmd)
	return root, nil
}

func NewAuditListCommand() (*AuditListCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &AuditListCommand{CommandDescription: glazedcmds.NewCommandDescription("list", glazedcmds.WithShort("List audit events for an org"), glazedcmds.WithLong(`List audit events for an organization with optional filters.

Examples:
  go-go-host audit list --org-id org_123 --dev-user alice
  go-go-host audit list --org-id org_123 --action deployment.activate --limit 20 --output json
`), glazedcmds.WithFlags(
		fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
		fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
		fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
		fields.New("org-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("organization ID")),
		fields.New("resource-id", fields.TypeString, fields.WithHelp("resource ID filter")),
		fields.New("actor-type", fields.TypeString, fields.WithHelp("actor type filter")),
		fields.New("actor-id", fields.TypeString, fields.WithHelp("actor ID filter")),
		fields.New("action", fields.TypeString, fields.WithHelp("action filter")),
		fields.New("limit", fields.TypeInteger, fields.WithDefault(100), fields.WithHelp("maximum events to return")),
	), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func (c *AuditListCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &AuditSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	q := url.Values{}
	if settings.ResourceID != "" {
		q.Set("resource_id", settings.ResourceID)
	}
	if settings.ActorType != "" {
		q.Set("actor_type", settings.ActorType)
	}
	if settings.ActorID != "" {
		q.Set("actor_id", settings.ActorID)
	}
	if settings.Action != "" {
		q.Set("action", settings.Action)
	}
	if settings.Limit > 0 {
		q.Set("limit", strconv.Itoa(settings.Limit))
	}
	path := fmt.Sprintf("/api/v1/orgs/%s/audit", settings.OrgID)
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var events []auditDTO
	if err := getJSONWithAuth(resolved.APIURL, path, resolved.DevUser, resolved.BearerToken, &events); err != nil {
		return err
	}
	for _, event := range events {
		if err := gp.AddRow(ctx, auditRow(event)); err != nil {
			return err
		}
	}
	return nil
}

func auditRow(event auditDTO) types.Row {
	return types.NewRow(types.MRP("id", event.ID), types.MRP("org_id", event.OrgID), types.MRP("actor_type", event.ActorType), types.MRP("actor_id", event.ActorID), types.MRP("action", event.Action), types.MRP("resource_type", event.ResourceType), types.MRP("resource_id", event.ResourceID), types.MRP("ip_address", event.IPAddress), types.MRP("user_agent", event.UserAgent), types.MRP("metadata_json", event.MetadataJSON), types.MRP("created_at", event.CreatedAt))
}
