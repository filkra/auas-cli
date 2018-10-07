package group

import (
	"fmt"
	"github.com/filkra/auas-cli/api"
	"github.com/filkra/auas-cli/yml"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var groupUpdateCommand = &cobra.Command{
	Use:   "update [file]",
	Short: "Updates all groups using the specified file",
	SilenceErrors: true,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		// Open the file containing all groups to import
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}

		// Parse the file
		groups, err := yml.ReadGroups(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		// Delete all groups using the client
		err = client.UpdateGroups(groups)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Groups updated successfully")
	},
}
