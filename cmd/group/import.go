package group

import (
	"fmt"
	"github.com/filkra/auas-cli/api"
	"github.com/filkra/auas-cli/yml"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var groupAddCommand = &cobra.Command{
	Use:   "import [file]",
	Short: "Imports groups from a YAML file into AUAS",
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

		// Import all groups using the client
		err = client.ImportGroups(groups)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Groups imported successfully")
	},
}
