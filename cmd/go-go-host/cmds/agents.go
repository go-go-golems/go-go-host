package cmds

import (
	"context"
	"fmt"

	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/spf13/cobra"
)

type AgentListCommand struct{ *glazedcmds.CommandDescription }
type AgentCreateCommand struct{ *glazedcmds.CommandDescription }

type AgentSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
	OrgID       string `glazed:"org-id"`
	Name        string `glazed:"name"`
}

type agentDTO struct {
	ID              string `json:"id"`
	OrgID           string `json:"orgId"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	CreatedByUserID string `json:"createdByUserId"`
	CreatedAt       string `json:"createdAt"`
	LastSeenAt      string `json:"lastSeenAt"`
}

var _ glazedcmds.GlazeCommand = (*AgentListCommand)(nil)
var _ glazedcmds.GlazeCommand = (*AgentCreateCommand)(nil)

func NewAgentsCobraCommand() (*cobra.Command, error) {
	root := &cobra.Command{Use: "agents", Aliases: []string{"agent"}, Short: "Manage deployment agents"}
	list, err := NewAgentListCommand()
	if err != nil {
		return nil, err
	}
	create, err := NewAgentCreateCommand()
	if err != nil {
		return nil, err
	}
	listCmd, err := BuildGlazedCobraCommand(list)
	if err != nil {
		return nil, err
	}
	createCmd, err := BuildGlazedCobraCommand(create)
	if err != nil {
		return nil, err
	}
	root.AddCommand(listCmd, createCmd)
	return root, nil
}

func agentFlags(requireName bool) []*fields.Definition {
	flags := []*fields.Definition{
		fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
		fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
		fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
		fields.New("org-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("organization ID")),
	}
	if requireName {
		flags = append(flags, fields.New("name", fields.TypeString, fields.WithRequired(true), fields.WithHelp("agent display name")))
	}
	return flags
}

func NewAgentListCommand() (*AgentListCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &AgentListCommand{CommandDescription: glazedcmds.NewCommandDescription("list", glazedcmds.WithShort("List deployment agents for an org"), glazedcmds.WithLong(`List deployment agents for an organization.

Examples:
  go-go-host agents list --org-id org_123 --dev-user alice
  go-go-host agents list --org-id org_123 --output json
`), glazedcmds.WithFlags(agentFlags(false)...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func NewAgentCreateCommand() (*AgentCreateCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &AgentCreateCommand{CommandDescription: glazedcmds.NewCommandDescription("create", glazedcmds.WithShort("Create a deployment agent record"), glazedcmds.WithLong(`Create a deployment agent record.

Examples:
  go-go-host agents create --org-id org_123 --name ci-bot --dev-user alice
`), glazedcmds.WithFlags(agentFlags(true)...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func (c *AgentListCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeAgentSettings(vals)
	if err != nil {
		return err
	}
	var agents []agentDTO
	if err := getJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/orgs/%s/agents", settings.OrgID), resolved.DevUser, resolved.BearerToken, &agents); err != nil {
		return err
	}
	for _, agent := range agents {
		if err := gp.AddRow(ctx, agentRow(agent)); err != nil {
			return err
		}
	}
	return nil
}

func (c *AgentCreateCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeAgentSettings(vals)
	if err != nil {
		return err
	}
	var agent agentDTO
	if err := postJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/orgs/%s/agents", settings.OrgID), resolved.DevUser, resolved.BearerToken, map[string]string{"name": settings.Name}, &agent); err != nil {
		return err
	}
	return gp.AddRow(ctx, agentRow(agent))
}

func decodeAgentSettings(vals *values.Values) (*AgentSettings, CLIConfig, error) {
	settings := &AgentSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return nil, CLIConfig{}, err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	return settings, resolved, err
}

func agentRow(agent agentDTO) types.Row {
	return types.NewRow(types.MRP("id", agent.ID), types.MRP("org_id", agent.OrgID), types.MRP("name", agent.Name), types.MRP("status", agent.Status), types.MRP("created_by_user_id", agent.CreatedByUserID), types.MRP("created_at", agent.CreatedAt), types.MRP("last_seen_at", agent.LastSeenAt))
}
