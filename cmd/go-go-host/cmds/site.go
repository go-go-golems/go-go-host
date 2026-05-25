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

type SiteListCommand struct{ *glazedcmds.CommandDescription }
type SiteCreateCommand struct{ *glazedcmds.CommandDescription }
type SiteRuntimeCommand struct{ *glazedcmds.CommandDescription }

type SiteListSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
	OrgID       string `glazed:"org-id"`
}

type SiteCreateSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
	OrgID       string `glazed:"org-id"`
	Slug        string `glazed:"slug"`
	Name        string `glazed:"name"`
}

type SiteRuntimeSettings struct {
	APIURL      string `glazed:"api-url"`
	DevUser     string `glazed:"dev-user"`
	BearerToken string `glazed:"bearer-token"`
	SiteID      string `glazed:"site-id"`
}

type runtimeStatusDTO struct {
	SiteID        string   `json:"siteId"`
	OrgID         string   `json:"orgId"`
	DeploymentID  string   `json:"deploymentId"`
	Hosts         []string `json:"hosts"`
	Status        string   `json:"status"`
	StartedAt     string   `json:"startedAt"`
	LastError     string   `json:"lastError"`
	RequestsTotal uint64   `json:"requestsTotal"`
	ErrorsTotal   uint64   `json:"errorsTotal"`
}

type siteDTO struct {
	ID                 string `json:"id"`
	OrgID              string `json:"orgId"`
	Slug               string `json:"slug"`
	Name               string `json:"name"`
	PrimaryHost        string `json:"primaryHost"`
	Status             string `json:"status"`
	ActiveDeploymentID string `json:"activeDeploymentId"`
}

var _ glazedcmds.GlazeCommand = (*SiteListCommand)(nil)
var _ glazedcmds.GlazeCommand = (*SiteCreateCommand)(nil)
var _ glazedcmds.GlazeCommand = (*SiteRuntimeCommand)(nil)

func NewSiteCobraCommand() (*cobra.Command, error) {
	siteCmd := &cobra.Command{Use: "site", Aliases: []string{"sites"}, Short: "Manage sites"}
	listCmd, err := NewSiteListCobraCommand()
	if err != nil {
		return nil, err
	}
	createCmd, err := NewSiteCreateCobraCommand()
	if err != nil {
		return nil, err
	}
	runtimeCmd, err := NewSiteRuntimeCobraCommand()
	if err != nil {
		return nil, err
	}
	siteCmd.AddCommand(listCmd, createCmd, runtimeCmd)
	return siteCmd, nil
}

func NewSiteListCobraCommand() (*cobra.Command, error) {
	command, err := NewSiteListCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewSiteCreateCobraCommand() (*cobra.Command, error) {
	command, err := NewSiteCreateCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewSiteRuntimeCobraCommand() (*cobra.Command, error) {
	command, err := NewSiteRuntimeCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewSiteListCommand() (*SiteListCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &SiteListCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"list",
		glazedcmds.WithShort("List sites in an organization"),
		glazedcmds.WithFlags(commonSiteFlags(true)...),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func NewSiteRuntimeCommand() (*SiteRuntimeCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := append(commonSiteFlags(false), fields.New("site-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("site ID")))
	return &SiteRuntimeCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"runtime",
		glazedcmds.WithShort("Show runtime status for a site"),
		glazedcmds.WithLong(`Show runtime status for a site.

Examples:
  go-go-host site runtime --site-id site_123 --dev-user alice
  go-go-host site runtime --site-id site_123 --output json
`),
		glazedcmds.WithFlags(flags...),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func NewSiteCreateCommand() (*SiteCreateCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := commonSiteFlags(true)
	flags = append(flags,
		fields.New("slug", fields.TypeString, fields.WithRequired(true), fields.WithHelp("site slug")),
		fields.New("name", fields.TypeString, fields.WithRequired(true), fields.WithHelp("site display name")),
	)
	return &SiteCreateCommand{CommandDescription: glazedcmds.NewCommandDescription(
		"create",
		glazedcmds.WithShort("Create a site in an organization"),
		glazedcmds.WithFlags(flags...),
		glazedcmds.WithSections(glazedSection, commandSettingsSection),
	)}, nil
}

func commonSiteFlags(requireOrg bool) []*fields.Definition {
	orgField := fields.New("org-id", fields.TypeString, fields.WithHelp("organization ID"))
	if requireOrg {
		orgField = fields.New("org-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("organization ID"))
	}
	return []*fields.Definition{
		fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
		fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
		fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
		orgField,
	}
}

func (c *SiteListCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &SiteListSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	var sites []siteDTO
	if err := getJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/orgs/%s/sites", settings.OrgID), resolved.DevUser, resolved.BearerToken, &sites); err != nil {
		return err
	}
	for _, site := range sites {
		if err := gp.AddRow(ctx, siteRow(site)); err != nil {
			return err
		}
	}
	return nil
}

func (c *SiteRuntimeCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &SiteRuntimeSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	var status runtimeStatusDTO
	if err := getJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/sites/%s/runtime", settings.SiteID), resolved.DevUser, resolved.BearerToken, &status); err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(
		types.MRP("site_id", status.SiteID),
		types.MRP("org_id", status.OrgID),
		types.MRP("deployment_id", status.DeploymentID),
		types.MRP("hosts", status.Hosts),
		types.MRP("status", status.Status),
		types.MRP("started_at", status.StartedAt),
		types.MRP("last_error", status.LastError),
		types.MRP("requests_total", status.RequestsTotal),
		types.MRP("errors_total", status.ErrorsTotal),
	))
}

func (c *SiteCreateCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &SiteCreateSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	if err != nil {
		return err
	}
	var site siteDTO
	if err := postJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/orgs/%s/sites", settings.OrgID), resolved.DevUser, resolved.BearerToken, map[string]string{"slug": settings.Slug, "name": settings.Name}, &site); err != nil {
		return err
	}
	return gp.AddRow(ctx, siteRow(site))
}

func siteRow(site siteDTO) types.Row {
	return types.NewRow(
		types.MRP("id", site.ID),
		types.MRP("org_id", site.OrgID),
		types.MRP("slug", site.Slug),
		types.MRP("name", site.Name),
		types.MRP("primary_host", site.PrimaryHost),
		types.MRP("status", site.Status),
		types.MRP("active_deployment_id", site.ActiveDeploymentID),
	)
}
