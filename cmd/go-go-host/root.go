package main

import (
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	commands "github.com/go-go-golems/go-go-host/cmd/go-go-host/cmds"
	"github.com/go-go-golems/go-go-host/cmd/go-go-host/doc"
	"github.com/spf13/cobra"
)

func newRootCommand() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "go-go-host",
		Short: "Manage go-go-host sites, deployments, agents, and runtime state",
		Long:  "go-go-host is the human CLI for the Goja hosting platform.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return logging.InitLoggerFromCobra(cmd)
		},
	}
	if err := logging.AddLoggingSectionToRootCommand(rootCmd, "go-go-host"); err != nil {
		return nil, err
	}
	helpSystem := help.NewHelpSystem()
	if err := doc.AddDocToHelpSystem(helpSystem); err != nil {
		return nil, err
	}
	help_cmd.SetupCobraRootCommand(helpSystem, rootCmd)

	statusCmd, err := commands.NewStatusCobraCommand()
	if err != nil {
		return nil, err
	}
	loginCmd, err := commands.NewLoginCobraCommand()
	if err != nil {
		return nil, err
	}
	meCmd, err := commands.NewMeCobraCommand()
	if err != nil {
		return nil, err
	}
	orgCmd, err := commands.NewOrgCobraCommand()
	if err != nil {
		return nil, err
	}
	siteCmd, err := commands.NewSiteCobraCommand()
	if err != nil {
		return nil, err
	}
	deployCmd, err := commands.NewDeployCobraCommand()
	if err != nil {
		return nil, err
	}
	deploymentsCmd, err := commands.NewDeploymentsCobraCommand()
	if err != nil {
		return nil, err
	}
	rollbackCmd, err := commands.NewRollbackCobraCommand()
	if err != nil {
		return nil, err
	}
	rootCmd.AddCommand(statusCmd, loginCmd, meCmd, orgCmd, siteCmd, deployCmd, deploymentsCmd, rollbackCmd)
	return rootCmd, nil
}
