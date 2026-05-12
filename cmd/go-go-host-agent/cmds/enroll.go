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

type EnrollCommand struct{ *glazedcmds.CommandDescription }
type EnrollSettings struct {
	Config string `glazed:"config"`
	APIURL string `glazed:"api-url"`
	Token  string `glazed:"token"`
}

type enrollResponse struct {
	Agent struct {
		ID     string `json:"id"`
		OrgID  string `json:"orgId"`
		Status string `json:"status"`
		Name   string `json:"name"`
	} `json:"agent"`
	KeyID string `json:"keyId"`
	Grant *struct {
		SiteID          string   `json:"siteId"`
		AllowedChannels []string `json:"allowedChannels"`
		AllowedPaths    []string `json:"allowedPaths"`
	} `json:"grant"`
}

var _ glazedcmds.GlazeCommand = (*EnrollCommand)(nil)

func NewEnrollCobraCommand() (*cobra.Command, error) {
	c, err := NewEnrollCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(c)
}
func NewEnrollCommand() (*EnrollCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &EnrollCommand{CommandDescription: glazedcmds.NewCommandDescription("enroll", glazedcmds.WithShort("Enroll this machine as a go-go-host deployment agent"), glazedcmds.WithLong(`Exchange a one-time enrollment token for an agent key registration.

Examples:
  go-go-host-agent enroll --token enroll_... --config ./agent.json
`), glazedcmds.WithFlags(fields.New("config", fields.TypeString, fields.WithHelp("agent config path")), fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host API URL")), fields.New("token", fields.TypeString, fields.WithRequired(true), fields.WithHelp("one-time enrollment token"))), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func (c *EnrollCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &EnrollSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	path := settings.Config
	if path == "" {
		path = defaultConfigPath()
	}
	cfg, err := loadConfig(path)
	if err != nil {
		return err
	}
	apiURL := settings.APIURL
	if cfg.APIURL != "" && settings.APIURL == defaultAPIURL {
		apiURL = cfg.APIURL
	}
	var resp enrollResponse
	if err := postJSON(apiURL, "/api/v1/agent/enroll", map[string]string{"token": settings.Token, "publicKey": cfg.PublicKey}, &resp); err != nil {
		return err
	}
	cfg.APIURL = apiURL
	cfg.AgentID = resp.Agent.ID
	cfg.KeyID = resp.KeyID
	cfg.OrgID = resp.Agent.OrgID
	if resp.Grant != nil {
		cfg.SiteID = resp.Grant.SiteID
	}
	if err := saveConfig(path, cfg); err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(types.MRP("config", path), types.MRP("agent_id", cfg.AgentID), types.MRP("key_id", cfg.KeyID), types.MRP("org_id", cfg.OrgID), types.MRP("site_id", cfg.SiteID), types.MRP("status", resp.Agent.Status)))
}
