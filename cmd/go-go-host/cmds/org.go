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

type OrgListCommand struct{ *glazedcmds.CommandDescription }
type OrgCreateCommand struct{ *glazedcmds.CommandDescription }

type OrgListSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
}

type OrgCreateSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
	Slug        string `glazed:"slug"`
	Name        string `glazed:"name"`
}

type orgDTO struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type orgMembershipDTO struct {
	OrgID   string `json:"orgId"`
	OrgSlug string `json:"orgSlug"`
	OrgName string `json:"orgName"`
	Role    string `json:"role"`
}

var _ glazedcmds.GlazeCommand = (*OrgListCommand)(nil)
var _ glazedcmds.GlazeCommand = (*OrgCreateCommand)(nil)

func NewOrgCobraCommand() (*cobra.Command, error) {
	orgCmd := &cobra.Command{Use: "org", Short: "Manage organizations"}
	listCmd, err := NewOrgListCobraCommand()
	if err != nil {
		return nil, err
	}
	createCmd, err := NewOrgCreateCobraCommand()
	if err != nil {
		return nil, err
	}
	orgCmd.AddCommand(listCmd, createCmd)
	return orgCmd, nil
}

func NewOrgListCobraCommand() (*cobra.Command, error) {
	command, err := NewOrgListCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewOrgCreateCobraCommand() (*cobra.Command, error) {
	command, err := NewOrgCreateCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewOrgListCommand() (*OrgListCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &OrgListCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"list",
		glazedcmds.WithShort("List organizations for the current user"),
		glazedcmds.WithFlags(
			fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
			fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
			fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
		),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func NewOrgCreateCommand() (*OrgCreateCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &OrgCreateCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"create",
		glazedcmds.WithShort("Create an organization"),
		glazedcmds.WithFlags(
			fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
			fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
			fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
			fields.New("slug", fields.TypeString, fields.WithRequired(true), fields.WithHelp("organization slug")),
			fields.New("name", fields.TypeString, fields.WithRequired(true), fields.WithHelp("organization display name")),
		),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func (c *OrgListCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &OrgListSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	var orgs []orgMembershipDTO
	if err := getJSONWithAuth(resolved.APIURL, "/api/v1/orgs", resolved.DevUser, resolved.BearerToken, &orgs); err != nil {
		return err
	}
	for _, org := range orgs {
		if err := gp.AddRow(ctx, types.NewRow(types.MRP("org_id", org.OrgID), types.MRP("org_slug", org.OrgSlug), types.MRP("org_name", org.OrgName), types.MRP("role", org.Role))); err != nil {
			return err
		}
	}
	return nil
}

func (c *OrgCreateCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &OrgCreateSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	var org orgDTO
	if err := postJSONWithAuth(resolved.APIURL, "/api/v1/orgs", resolved.DevUser, resolved.BearerToken, map[string]string{"slug": settings.Slug, "name": settings.Name}, &org); err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(types.MRP("id", org.ID), types.MRP("slug", org.Slug), types.MRP("name", org.Name)))
}
