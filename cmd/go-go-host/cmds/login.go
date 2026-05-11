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

type LoginCommand struct{ *glazedcmds.CommandDescription }

type LoginSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
}

var _ glazedcmds.GlazeCommand = (*LoginCommand)(nil)

func NewLoginCobraCommand() (*cobra.Command, error) {
	command, err := NewLoginCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewLoginCommand() (*LoginCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &LoginCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"login",
		glazedcmds.WithShort("Store local CLI connection and auth settings"),
		glazedcmds.WithLong(`Store local go-go-host CLI settings.

This is a Phase 2 bridge command. For local development, use --dev-user. For
non-dev auth smoke tests, paste an OIDC ID token with --bearer-token. A browser
OAuth flow can later reuse the same config file.

Examples:
  go-go-host login --api-url http://127.0.0.1:8080 --dev-user manuel
  go-go-host login --api-url https://host.example --bearer-token "$TOKEN"
`),
		glazedcmds.WithFlags(
			fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
			fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header to store")),
			fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token to store")),
		),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func (c *LoginCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &LoginSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	if settings.DevUser == "" && settings.BearerToken == "" {
		return fmt.Errorf("one of --dev-user or --bearer-token is required")
	}
	cfg := CLIConfig{APIURL: settings.APIURL, DevUser: settings.DevUser, BearerToken: settings.BearerToken}
	path, err := saveCLIConfig(cfg)
	if err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(
		types.MRP("config_path", path),
		types.MRP("api_url", cfg.APIURL),
		types.MRP("dev_user", cfg.DevUser),
		types.MRP("has_bearer_token", cfg.BearerToken != ""),
	))
}
