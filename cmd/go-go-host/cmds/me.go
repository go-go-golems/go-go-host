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

type MeCommand struct{ *glazedcmds.CommandDescription }

type MeSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
}

type meResponse struct {
	User struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	Memberships []struct {
		OrgID   string `json:"orgId"`
		OrgSlug string `json:"orgSlug"`
		OrgName string `json:"orgName"`
		Role    string `json:"role"`
	} `json:"memberships"`
	PlatformAdmin bool `json:"platformAdmin"`
}

var _ glazedcmds.GlazeCommand = (*MeCommand)(nil)

func NewMeCobraCommand() (*cobra.Command, error) {
	command, err := NewMeCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewMeCommand() (*MeCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &MeCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"me",
		glazedcmds.WithShort("Show current go-go-host user and memberships"),
		glazedcmds.WithFlags(
			fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
			fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
			fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
		),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func (c *MeCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &MeSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	var me meResponse
	if err := getJSONWithAuth(resolved.APIURL, "/api/v1/me", resolved.DevUser, resolved.BearerToken, &me); err != nil {
		return err
	}
	if len(me.Memberships) == 0 {
		return gp.AddRow(ctx, types.NewRow(
			types.MRP("user_id", me.User.ID),
			types.MRP("email", me.User.Email),
			types.MRP("display_name", me.User.DisplayName),
			types.MRP("platform_admin", me.PlatformAdmin),
		))
	}
	for _, m := range me.Memberships {
		if err := gp.AddRow(ctx, types.NewRow(
			types.MRP("user_id", me.User.ID),
			types.MRP("email", me.User.Email),
			types.MRP("display_name", me.User.DisplayName),
			types.MRP("platform_admin", me.PlatformAdmin),
			types.MRP("org_id", m.OrgID),
			types.MRP("org_slug", m.OrgSlug),
			types.MRP("org_name", m.OrgName),
			types.MRP("role", m.Role),
		)); err != nil {
			return err
		}
	}
	return nil
}
