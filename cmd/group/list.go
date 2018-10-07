package group

import (
	"fmt"
	"github.com/filkra/auas-cli/api"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/tabwriter"
)

var groupListCommand = &cobra.Command{
	Use:   "list [courseId]",
	Short: "Lists all groups within the specified course",
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, getTableHeader())
		for _, row := range groups {
			fmt.Fprintln(w, rowToString(&row))
		}

		w.Flush()
	},
}

func rowToString(information *api.GroupInformation) string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t",
		information.Id,
		information.Name,
		information.Room,
		information.Day,
		information.Time,
		information.Tutor,
		information.Participants,
		information.Space,
	)
}

func getTableHeader() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t",
		"Id",
		"Name",
		"Raum",
		"Wochentag",
		"Uhrzeit",
		"Tutor",
		"Teilnehmer",
		"Plätze",
	)
}
