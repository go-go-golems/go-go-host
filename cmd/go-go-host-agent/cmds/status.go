package cmds

import (
	"context"

	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/spf13/cobra"
)

type StatusCommand struct {
	*glazedcmds.CommandDescription
}

type StatusSettings struct {
	APIURL string `glazed:"api-url"`
	Config string `glazed:"config"`
}

var _ glazedcmds.GlazeCommand = (*StatusCommand)(nil)

func NewStatusCobraCommand() (*cobra.Command, error) {
	command, err := NewStatusCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewStatusCommand() (*StatusCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &StatusCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"status",
		glazedcmds.WithShort("Check go-go-host daemon reachability for the agent CLI"),
		glazedcmds.WithLong(`Check the go-go-host daemon before agent enrollment or deployment.

Examples:
  go-go-host-agent status --api-url http://127.0.0.1:8080
  go-go-host-agent status --output json
`),
		glazedcmds.WithFlags(
			fields.New("api-url", fields.TypeString,
				fields.WithDefault(defaultAPIURL),
				fields.WithHelp("go-go-host daemon API base URL")),
			fields.New("config", fields.TypeString,
				fields.WithHelp("agent config path")),
		),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func (c *StatusCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &StatusSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	var health struct {
		Status string `json:"status"`
	}
	if err := getJSON(settings.APIURL, "/healthz", &health); err != nil {
		return err
	}
	var version struct {
		Version string `json:"version"`
	}
	if err := getJSON(settings.APIURL, "/api/v1/version", &version); err != nil {
		return err
	}
	cfg, _ := loadConfig(settings.Config)
	return gp.AddRow(ctx, types.NewRow(
		types.MRP("api_url", settings.APIURL),
		types.MRP("status", health.Status),
		types.MRP("version", version.Version),
		types.MRP("agent_id", cfg.AgentID),
		types.MRP("key_id", cfg.KeyID),
		types.MRP("org_id", cfg.OrgID),
		types.MRP("site_id", cfg.SiteID),
		types.MRP("enrolled", cfg.AgentID != "" && cfg.KeyID != ""),
	))
}
