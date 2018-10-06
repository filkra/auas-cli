package group

import (
	"fmt"
	"github.com/filkra/auas-cli/api"
	"github.com/spf13/cobra"
	"log"
)

var groupDeleteCommand = &cobra.Command{
	Use:   "delete [courseId]",
	Short: "Deletes all groups from the specified course",
	SilenceErrors: true,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		// Import all groups using the client
		groups, err := client.GetGroups(args[0])
		if err != nil {
			log.Fatal(err)
		}

		// Delete all groups using the client
		err = client.DeleteGroups(groups)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Groups deleted successfully")
	},
}
