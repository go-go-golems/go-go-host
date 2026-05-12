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

type DeployCommand struct{ *glazedcmds.CommandDescription }
type DeploymentListCommand struct{ *glazedcmds.CommandDescription }
type DeploymentShowCommand struct{ *glazedcmds.CommandDescription }
type DeploymentActivateCommand struct{ *glazedcmds.CommandDescription }
type RollbackCommand struct{ *glazedcmds.CommandDescription }

type DeploymentSettings struct {
	APIURL       string `glazed:"api-url"`
	DevUser      string `glazed:"dev-user"`
	BearerToken  string `glazed:"bearer-token"`
	SiteID       string `glazed:"site-id"`
	DeploymentID string `glazed:"deployment-id"`
	Path         string `glazed:"path"`
	Message      string `glazed:"message"`
	Channel      string `glazed:"channel"`
}

type deploymentDTO struct {
	ID             string `json:"id"`
	SiteID         string `json:"siteId"`
	Version        int    `json:"version"`
	Status         string `json:"status"`
	BundleRef      string `json:"bundleRef"`
	UnpackedPath   string `json:"unpackedPath"`
	ManifestJSON   string `json:"manifestJson"`
	ValidationJSON string `json:"validationJson"`
	CreatedByType  string `json:"createdByType"`
	CreatedByID    string `json:"createdById"`
	CreatedAt      string `json:"createdAt"`
	ActivatedAt    string `json:"activatedAt"`
	BundleSHA256   string `json:"bundleSha256"`
}

type deployResponse struct {
	Deployment deploymentDTO  `json:"deployment"`
	Report     map[string]any `json:"report"`
	Manifest   map[string]any `json:"manifest"`
}

var _ glazedcmds.GlazeCommand = (*DeployCommand)(nil)
var _ glazedcmds.GlazeCommand = (*DeploymentListCommand)(nil)
var _ glazedcmds.GlazeCommand = (*DeploymentShowCommand)(nil)
var _ glazedcmds.GlazeCommand = (*DeploymentActivateCommand)(nil)
var _ glazedcmds.GlazeCommand = (*RollbackCommand)(nil)

func NewDeployCobraCommand() (*cobra.Command, error) {
	command, err := NewDeployCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func NewDeploymentsCobraCommand() (*cobra.Command, error) {
	root := &cobra.Command{Use: "deployments", Aliases: []string{"deployment"}, Short: "Manage deployments"}
	list, err := NewDeploymentListCommand()
	if err != nil {
		return nil, err
	}
	show, err := NewDeploymentShowCommand()
	if err != nil {
		return nil, err
	}
	activate, err := NewDeploymentActivateCommand()
	if err != nil {
		return nil, err
	}
	listCmd, err := BuildGlazedCobraCommand(list)
	if err != nil {
		return nil, err
	}
	showCmd, err := BuildGlazedCobraCommand(show)
	if err != nil {
		return nil, err
	}
	activateCmd, err := BuildGlazedCobraCommand(activate)
	if err != nil {
		return nil, err
	}
	root.AddCommand(listCmd, showCmd, activateCmd)
	return root, nil
}

func NewRollbackCobraCommand() (*cobra.Command, error) {
	command, err := NewRollbackCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(command)
}

func deploymentCommonFlags() []*fields.Definition {
	return []*fields.Definition{
		fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host daemon API base URL")),
		fields.New("dev-user", fields.TypeString, fields.WithHelp("dev auth user subject header")),
		fields.New("bearer-token", fields.TypeString, fields.WithHelp("OIDC bearer token for non-dev auth mode")),
	}
}

func NewDeployCommand() (*DeployCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := append(deploymentCommonFlags(),
		fields.New("site-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("site ID to deploy to")),
		fields.New("path", fields.TypeString, fields.WithRequired(true), fields.WithHelp("bundle archive path (.tar.gz or .zip)")),
		fields.New("message", fields.TypeString, fields.WithHelp("deployment message")),
		fields.New("channel", fields.TypeString, fields.WithHelp("deployment channel")),
	)
	return &DeployCommand{CommandDescription: glazedcmds.NewCommandDescription("deploy", glazedcmds.WithShort("Upload and validate a deployment bundle"), glazedcmds.WithLong(`Upload and validate a go-go-host deployment bundle.

Examples:
  go-go-host deploy --site-id site_123 --path ./hello.tar.gz --message "initial deploy" --dev-user alice
  go-go-host deploy --site-id site_123 --path ./hello.zip --output json
`), glazedcmds.WithFlags(flags...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func NewDeploymentListCommand() (*DeploymentListCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := append(deploymentCommonFlags(), fields.New("site-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("site ID")))
	return &DeploymentListCommand{CommandDescription: glazedcmds.NewCommandDescription("list", glazedcmds.WithShort("List deployments for a site"), glazedcmds.WithLong(`List deployments for a site.

Examples:
  go-go-host deployments list --site-id site_123 --dev-user alice
  go-go-host deployments list --site-id site_123 --output json
`), glazedcmds.WithFlags(flags...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func NewDeploymentShowCommand() (*DeploymentShowCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := append(deploymentCommonFlags(), fields.New("deployment-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("deployment ID")))
	return &DeploymentShowCommand{CommandDescription: glazedcmds.NewCommandDescription("show", glazedcmds.WithShort("Show deployment details"), glazedcmds.WithLong(`Show one deployment.

Examples:
  go-go-host deployments show --deployment-id dep_123 --dev-user alice --output yaml
`), glazedcmds.WithFlags(flags...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func NewDeploymentActivateCommand() (*DeploymentActivateCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := append(deploymentCommonFlags(), fields.New("deployment-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("deployment ID")))
	return &DeploymentActivateCommand{CommandDescription: glazedcmds.NewCommandDescription("activate", glazedcmds.WithShort("Activate a validated deployment"), glazedcmds.WithLong(`Activate a validated deployment.

Examples:
  go-go-host deployments activate --deployment-id dep_123 --dev-user alice
`), glazedcmds.WithFlags(flags...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func NewRollbackCommand() (*RollbackCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	flags := append(deploymentCommonFlags(), fields.New("site-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("site ID")))
	return &RollbackCommand{CommandDescription: glazedcmds.NewCommandDescription("rollback", glazedcmds.WithShort("Roll back a site to the previous deployment"), glazedcmds.WithLong(`Roll back a site by activating the previous validated/superseded deployment.

Examples:
  go-go-host rollback --site-id site_123 --dev-user alice
`), glazedcmds.WithFlags(flags...), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func (c *DeployCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeDeploymentSettings(vals)
	if err != nil {
		return err
	}
	var res deployResponse
	if err := postMultipartBundleWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/sites/%s/deployments", settings.SiteID), resolved.DevUser, resolved.BearerToken, settings.Path, map[string]string{"message": settings.Message, "channel": settings.Channel}, &res); err != nil {
		return err
	}
	return gp.AddRow(ctx, deploymentRow(res.Deployment, types.NewRow(types.MRP("validation_report", res.Report), types.MRP("manifest", res.Manifest))))
}

func (c *DeploymentListCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeDeploymentSettings(vals)
	if err != nil {
		return err
	}
	var deps []deploymentDTO
	if err := getJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/sites/%s/deployments", settings.SiteID), resolved.DevUser, resolved.BearerToken, &deps); err != nil {
		return err
	}
	for _, dep := range deps {
		if err := gp.AddRow(ctx, deploymentRow(dep)); err != nil {
			return err
		}
	}
	return nil
}

func (c *DeploymentShowCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeDeploymentSettings(vals)
	if err != nil {
		return err
	}
	var dep deploymentDTO
	if err := getJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/deployments/%s", settings.DeploymentID), resolved.DevUser, resolved.BearerToken, &dep); err != nil {
		return err
	}
	return gp.AddRow(ctx, deploymentRow(dep))
}

func (c *DeploymentActivateCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeDeploymentSettings(vals)
	if err != nil {
		return err
	}
	var dep deploymentDTO
	if err := postJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/deployments/%s/activate", settings.DeploymentID), resolved.DevUser, resolved.BearerToken, map[string]any{}, &dep); err != nil {
		return err
	}
	return gp.AddRow(ctx, deploymentRow(dep))
}

func (c *RollbackCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings, resolved, err := decodeDeploymentSettings(vals)
	if err != nil {
		return err
	}
	var dep deploymentDTO
	if err := postJSONWithAuth(resolved.APIURL, fmt.Sprintf("/api/v1/sites/%s/rollback", settings.SiteID), resolved.DevUser, resolved.BearerToken, map[string]any{}, &dep); err != nil {
		return err
	}
	return gp.AddRow(ctx, deploymentRow(dep))
}

func decodeDeploymentSettings(vals *values.Values) (*DeploymentSettings, CLIConfig, error) {
	settings := &DeploymentSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return nil, CLIConfig{}, err
	}
	resolved, err := resolveCLISettings(settings.APIURL, settings.DevUser, settings.BearerToken)
	return settings, resolved, err
}

func deploymentRow(dep deploymentDTO, extra ...types.Row) types.Row {
	row := types.NewRow(
		types.MRP("id", dep.ID), types.MRP("site_id", dep.SiteID), types.MRP("version", dep.Version), types.MRP("status", dep.Status),
		types.MRP("bundle_ref", dep.BundleRef), types.MRP("bundle_sha256", dep.BundleSHA256), types.MRP("unpacked_path", dep.UnpackedPath), types.MRP("created_by_type", dep.CreatedByType), types.MRP("created_by_id", dep.CreatedByID), types.MRP("created_at", dep.CreatedAt), types.MRP("activated_at", dep.ActivatedAt), types.MRP("manifest_json", dep.ManifestJSON), types.MRP("validation_json", dep.ValidationJSON),
	)
	for _, e := range extra {
		for pair := e.Oldest(); pair != nil; pair = pair.Next() {
			row.Set(pair.Key, pair.Value)
		}
	}
	return row
}
