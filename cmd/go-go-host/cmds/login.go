package cmds

import (
	"context"
	"fmt"
	"os"

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
	ClientID    string `glazed:"client-id"`
	Scopes      string `glazed:"scopes"`
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

For local development, use --dev-user. For non-dev one-off smoke tests, paste
an OIDC token with --bearer-token. For production, omit both and the CLI starts
OAuth 2.0 Device Authorization Grant: it prints a Keycloak URL and user code,
then stores the returned access and refresh tokens after browser approval.

Examples:
  go-go-host login --api-url http://127.0.0.1:8080 --dev-user manuel
  go-go-host login --api-url https://hosting.yolo.scapegoat.dev
  go-go-host login --api-url https://host.example --bearer-token "$TOKEN"
`),
		glazedcmds.WithFlags(
			fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
			fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header to store")),
			fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token to store")),
			fields.New("client-id", fields.TypeString, fields.WithHelp("OIDC client ID override for device flow")),
			fields.New("scopes", fields.TypeString, fields.WithHelp("OIDC scopes for device flow, separated by spaces or commas")),
		),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func (c *LoginCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &LoginSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	if settings.DevUser != "" || settings.BearerToken != "" {
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
			types.MRP("auth_mode", bridgeAuthMode(cfg)),
		))
	}
	return c.runDeviceLogin(ctx, settings, gp)
}

func (c *LoginCommand) runDeviceLogin(ctx context.Context, settings *LoginSettings, gp middlewares.Processor) error {
	apiCfg, err := fetchPublicConfig(ctx, settings.APIURL)
	if err != nil {
		return err
	}
	clientID := settings.ClientID
	if clientID == "" && apiCfg.OIDC != nil {
		clientID = apiCfg.OIDC.DeviceClientID
	}
	if clientID == "" && apiCfg.OIDC != nil {
		clientID = apiCfg.OIDC.ClientID
	}
	if clientID == "" {
		return fmt.Errorf("OIDC client id is required for device flow")
	}
	scopes := []string{"openid", "profile", "email"}
	if apiCfg.OIDC != nil {
		scopes = defaultScopes(apiCfg.OIDC.Scopes)
	}
	scopes = scopesFromString(settings.Scopes, scopes)
	discovery, err := discoverOIDC(ctx, apiCfg.OIDC.Issuer)
	if err != nil {
		return err
	}
	device, err := startDeviceAuthorization(ctx, discovery.DeviceAuthorizationEndpoint, clientID, scopes)
	if err != nil {
		return err
	}
	printDeviceInstructions(device)
	tok, err := pollDeviceToken(ctx, discovery.TokenEndpoint, clientID, device)
	if err != nil {
		return err
	}
	cfg := CLIConfig{APIURL: settings.APIURL, OIDC: sessionWithToken(apiCfg.OIDC.Issuer, clientID, scopes, tok, "")}
	path, err := saveCLIConfig(cfg)
	if err != nil {
		return err
	}
	var me meResponse
	if err := getJSONWithAuth(cfg.APIURL, "/api/v1/me", "", cfg.OIDC.AccessToken, &me); err != nil {
		return fmt.Errorf("login succeeded but API token validation failed: %w", err)
	}
	return gp.AddRow(ctx, types.NewRow(
		types.MRP("config_path", path),
		types.MRP("api_url", cfg.APIURL),
		types.MRP("auth_mode", "device"),
		types.MRP("issuer", cfg.OIDC.Issuer),
		types.MRP("client_id", cfg.OIDC.ClientID),
		types.MRP("has_refresh_token", cfg.OIDC.RefreshToken != ""),
		types.MRP("expires_at", cfg.OIDC.ExpiresAt),
		types.MRP("user_id", me.User.ID),
		types.MRP("email", me.User.Email),
		types.MRP("display_name", me.User.DisplayName),
	))
}

func printDeviceInstructions(device deviceAuthorizationResponse) {
	verification := device.VerificationURIComplete
	if verification == "" {
		verification = device.VerificationURI
	}
	fmt.Fprintf(os.Stderr, "\nOpen this URL in your browser:\n  %s\n\nEnter this code if prompted:\n  %s\n\nWaiting for browser authorization...\n", verification, device.UserCode)
}

func bridgeAuthMode(cfg CLIConfig) string {
	if cfg.DevUser != "" {
		return "dev-user"
	}
	if cfg.BearerToken != "" {
		return "bearer-token"
	}
	return "unknown"
}
