package cmds

import (
	"context"

	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/spf13/cobra"
)

type LogoutCommand struct{ *glazedcmds.CommandDescription }

var _ glazedcmds.GlazeCommand = (*LogoutCommand)(nil)

func NewLogoutCobraCommand() (*cobra.Command, error) {
	command, err := NewLogoutCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewLogoutCommand() (*LogoutCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &LogoutCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"logout",
		glazedcmds.WithShort("Clear local go-go-host CLI authentication tokens"),
		glazedcmds.WithLong(`Clear local go-go-host CLI authentication settings.

When the current login was created with OAuth Device Authorization Grant, logout
best-effort revokes the refresh token with Keycloak before deleting local token
state. Local token state is cleared even when revocation fails.
`),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func (c *LogoutCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	_ = vals
	cfg, err := loadCLIConfig()
	if err != nil {
		return err
	}
	path, err := defaultCLIConfigPath()
	if err != nil {
		return err
	}
	revoked := false
	revokeError := ""
	if cfg.OIDC != nil && cfg.OIDC.RefreshToken != "" {
		if ok, err := revokeOIDCToken(ctx, cfg.OIDC); err != nil {
			revokeError = err.Error()
		} else {
			revoked = ok
		}
	}
	cfg.DevUser = ""
	cfg.BearerToken = ""
	cfg.OIDC = nil
	if _, err := saveCLIConfig(cfg); err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(
		types.MRP("config_path", path),
		types.MRP("cleared", true),
		types.MRP("revoked", revoked),
		types.MRP("revoke_error", revokeError),
	))
}
