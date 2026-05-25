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

type DeployCommand struct{ *glazedcmds.CommandDescription }
type DeploySettings struct {
	Config     string `glazed:"config"`
	SiteID     string `glazed:"site-id"`
	Channel    string `glazed:"channel"`
	BundlePath string `glazed:"bundle-path"`
	Path       string `glazed:"path"`
	Bundle     string `glazed:"bundle"`
	Activate   bool   `glazed:"activate"`
}

type deployRunResponse struct {
	ID          string `json:"id"`
	SiteID      string `json:"siteId"`
	Status      string `json:"status"`
	UploadToken string `json:"uploadToken"`
	ExpiresAt   string `json:"expiresAt"`
}
type uploadResponse struct {
	DeployRunID string `json:"deployRunId"`
	Activated   bool   `json:"activated"`
	Deployment  struct {
		ID      string `json:"id"`
		SiteID  string `json:"siteId"`
		Status  string `json:"status"`
		Version int    `json:"version"`
	} `json:"deployment"`
	Report struct {
		Valid  bool     `json:"valid"`
		Errors []string `json:"errors"`
	} `json:"report"`
}

var _ glazedcmds.GlazeCommand = (*DeployCommand)(nil)

func NewDeployCobraCommand() (*cobra.Command, error) {
	c, err := NewDeployCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(c)
}
func NewDeployCommand() (*DeployCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &DeployCommand{CommandDescription: glazedcmds.NewCommandDescription("deploy", glazedcmds.WithShort("Deploy a bundle with signed agent credentials"), glazedcmds.WithLong(`Create a signed deploy run and upload a bundle archive.

Examples:
  go-go-host-agent deploy --bundle ./site.tar.gz --bundle-path bundles/site.tar.gz --site-id site_123
`), glazedcmds.WithFlags(fields.New("config", fields.TypeString, fields.WithHelp("agent config path")), fields.New("bundle", fields.TypeString, fields.WithRequired(true), fields.WithHelp("local bundle tar.gz or zip path")), fields.New("site-id", fields.TypeString, fields.WithHelp("site ID; defaults to enrolled grant site")), fields.New("channel", fields.TypeString, fields.WithDefault("default"), fields.WithHelp("deployment channel")), fields.New("bundle-path", fields.TypeString, fields.WithHelp("logical bundle path checked against grant")), fields.New("path", fields.TypeString, fields.WithHelp("deprecated alias for --bundle-path")), fields.New("activate", fields.TypeBool, fields.WithHelp("request auto-activation; requires agent grant can_activate"))), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func (c *DeployCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &DeploySettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	cfg, err := loadConfig(settings.Config)
	if err != nil {
		return err
	}
	siteID := settings.SiteID
	if siteID == "" {
		siteID = cfg.SiteID
	}
	logicalPath := settings.BundlePath
	if logicalPath == "" {
		logicalPath = settings.Path
	}
	if logicalPath == "" {
		logicalPath = settings.Bundle
	}
	var run deployRunResponse
	if err := signedPostJSON(cfg.APIURL, "/api/v1/agent/deploy-runs", map[string]any{"siteId": siteID, "channel": settings.Channel, "bundlePath": logicalPath, "action": "deploy", "activate": settings.Activate}, cfg, &run); err != nil {
		return err
	}
	var upload uploadResponse
	if err := uploadBundle(cfg.APIURL, run.ID, run.UploadToken, settings.Bundle, &upload); err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(types.MRP("deploy_run_id", run.ID), types.MRP("deployment_id", upload.Deployment.ID), types.MRP("site_id", upload.Deployment.SiteID), types.MRP("version", upload.Deployment.Version), types.MRP("status", upload.Deployment.Status), types.MRP("valid", upload.Report.Valid), types.MRP("activated", upload.Activated), types.MRP("errors", upload.Report.Errors)))
}
