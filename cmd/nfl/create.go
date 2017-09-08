package main

import (
	"log"

	"github.com/ameske/nfl-pickem/sqlite3"
	"github.com/spf13/cobra"
)

var createYear uint

func init() {
	CreateCmd.AddCommand(createPicksCmd)

	createPicksCmd.Flags().UintVarP(&createYear, "year", "y", 0, "NFL season year")
}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create objects in the database",
	Long:  "create objects in the database",
}

var createPicksCmd = &cobra.Command{
	Use:   "picks",
	Short: "create picks for all users",
	Long:  "create picks for all users",
	Run: func(cmd *cobra.Command, args []string) {
		if createYear == 0 {
			log.Fatal("year must be set via command line")
		}

		if datastore == "" {
			log.Fatal("db flag is required")
		}

		db, err := sqlite3.NewDatastore(datastore)
		if err != nil {
			log.Fatal(err)
		}

		users, err := db.Users()
		if err != nil {
			log.Fatal(err)
		}

		for i := 1; i <= 17; i++ {
			for _, u := range users {
				err := db.CreatePicks(u, int(createYear), i)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}
