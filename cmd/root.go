package cmd

import (
	"github.com/filkra/auas-cli/cmd/group"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use:   "auas-cli",
	Short: "A simple command line interface for managing AUAS",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCommand.AddCommand(group.RootCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
