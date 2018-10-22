package sql

import (
	"bufio"
	"fmt"
	"github.com/filkra/auas-cli/sql"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const (
	Prompt = "SQL> "
	TerminalCharacter = ';'
)

var sqlQueryCommand = &cobra.Command{
	Use:   "query",
	Short: "Performs a raw SQL query using stdin",
	SilenceErrors: true,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new SQL client
		client, err := sql.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		// Print prompt
		fmt.Print(Prompt)

		// Read all lines until user types in ';'
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString(TerminalCharacter)
		if err != nil {
			log.Fatal(err)
		}

		// Perform the SQL query
		response, err := client.Query(input)
		if err != nil {
			log.Fatal(err)
		}

		// Print the result
		fmt.Println(response)
	},
}
