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

type KeygenCommand struct{ *glazedcmds.CommandDescription }
type KeygenSettings struct {
	Config string `glazed:"config"`
	APIURL string `glazed:"api-url"`
}

var _ glazedcmds.GlazeCommand = (*KeygenCommand)(nil)

func NewKeygenCobraCommand() (*cobra.Command, error) {
	c, err := NewKeygenCommand()
	if err != nil {
		return nil, err
	}
	return BuildGlazedCobraCommand(c)
}

func NewKeygenCommand() (*KeygenCommand, error) {
	glazedSection, commandSettingsSection, err := standardSections()
	if err != nil {
		return nil, err
	}
	return &KeygenCommand{CommandDescription: glazedcmds.NewCommandDescription("keygen", glazedcmds.WithShort("Generate an Ed25519 agent key"), glazedcmds.WithLong(`Generate an Ed25519 key pair and store it in the local agent config.

Examples:
  go-go-host-agent keygen
  go-go-host-agent keygen --config ./agent.json --output json
`), glazedcmds.WithFlags(fields.New("config", fields.TypeString, fields.WithHelp("agent config path")), fields.New("api-url", fields.TypeString, fields.WithDefault(defaultAPIURL), fields.WithHelp("go-go-host API URL stored in config"))), glazedcmds.WithSections(glazedSection, commandSettingsSection))}, nil
}

func (c *KeygenCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &KeygenSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	pub, priv, err := generateKeyPair()
	if err != nil {
		return err
	}
	path := settings.Config
	if path == "" {
		path = defaultConfigPath()
	}
	cfg := AgentConfig{APIURL: settings.APIURL, PublicKey: pub, PrivateKey: priv}
	if err := saveConfig(path, cfg); err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(types.MRP("config", path), types.MRP("public_key", pub), types.MRP("status", "generated")))
}
