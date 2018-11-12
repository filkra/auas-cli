package group

import (
	"fmt"
	"github.com/filkra/auas-cli/api"
	"github.com/spf13/cobra"
	"log"
)

var groupCloseCommand = &cobra.Command{
	Use:   "close [courseId] [groupName...]",
	Short: "Closes specific groups within the specified course",
	SilenceErrors: true,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch all groups using the client
		groups, err := client.GetGroups(args[0])
		if err != nil {
			log.Fatal(err)
		}

		// Filter out groups and set the maximum participant count to 0
		groupNames := args[1:]
		filteredGroups := make([]api.GroupInformation, 0)
		for _, group := range groups {
			if contains(groupNames, group.Name)  {
				group.Space = "0"
				filteredGroups = append(filteredGroups, group)
			}
		}

		// Close the filtered groups
		err = client.WriteGroups(args[0], filteredGroups)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Groups closed successfully")
	},
}

func contains(elements []string, value string) bool {
	for _, element := range elements {
		if element == value {
			return true
		}
	}
	return false
}
