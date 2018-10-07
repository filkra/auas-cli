package group

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "groups",
	Short: "Provides functions for managing groups within courses",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(groupAddCommand)
	RootCommand.AddCommand(groupDeleteCommand)
	RootCommand.AddCommand(groupListCommand)
	RootCommand.AddCommand(groupUpdateCommand)
}
