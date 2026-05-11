package main

import (
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	commands "github.com/go-go-golems/go-go-host/cmd/go-go-host-agent/cmds"
	"github.com/go-go-golems/go-go-host/cmd/go-go-host-agent/doc"
	"github.com/spf13/cobra"
)

func newRootCommand() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "go-go-host-agent",
		Short: "Headless deploy-agent CLI for go-go-host",
		Long:  "go-go-host-agent enrolls machine identities and deploys Goja sites without human credentials.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return logging.InitLoggerFromCobra(cmd)
		},
	}
	if err := logging.AddLoggingSectionToRootCommand(rootCmd, "go-go-host-agent"); err != nil {
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
	rootCmd.AddCommand(statusCmd)
	return rootCmd, nil
}
