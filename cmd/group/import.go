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
		// Get the username
		user, present := os.LookupEnv("AUAS_USER")
		if present == false {
			log.Fatal("Please specify a username within the environment variable AUAS_USER")
		}

		// Get the password
		password, present := os.LookupEnv("AUAS_PASS")
		if present == false {
			log.Fatal("Please specify a password within the environment variable AUAS_PASS")
		}

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
		groups, err := yml.ReadGroup(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		// Login using the client
		err = client.Login(user, password)
		if err != nil {
			log.Fatal(err)
		}

		// Import all groups using the client
		err = client.ImportGroups(groups)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Groups imported successfully")
	},
}
