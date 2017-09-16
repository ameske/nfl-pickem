package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ameske/nfl-pickem/parser/results"
	"github.com/ameske/nfl-pickem/parser/schedule"
	"github.com/ameske/nfl-pickem/sqlite3"
	"github.com/spf13/cobra"
)

var scheduleYear, scheduleWeek uint
var scheduleFile string

func init() {
	ScheduleCmd.AddCommand(scheduleDownloadCmd)
	ScheduleCmd.AddCommand(scheduleResultsCmd)
	ScheduleCmd.AddCommand(scheduleImportCmd)

	scheduleDownloadCmd.Flags().UintVarP(&scheduleYear, "year", "y", 0, "NFL season year")
	scheduleDownloadCmd.Flags().UintVarP(&scheduleWeek, "week", "w", 0, "NFL season week")

	scheduleResultsCmd.Flags().UintVarP(&scheduleYear, "year", "y", 0, "NFL season year")
	scheduleResultsCmd.Flags().UintVarP(&scheduleWeek, "week", "w", 0, "NFL season week")

	scheduleResultsCmd.AddCommand(scheduleResultsImportCmd)
	scheduleResultsImportCmd.Flags().UintVarP(&scheduleYear, "year", "y", 0, "NFL season year")
	scheduleResultsImportCmd.Flags().UintVarP(&scheduleWeek, "week", "w", 0, "NFL season week")

	scheduleImportCmd.Flags().UintVarP(&scheduleYear, "year", "y", 0, "NFL season year")
	scheduleImportCmd.Flags().UintVarP(&scheduleWeek, "week", "w", 0, "NFL season week")
	scheduleImportCmd.Flags().StringVarP(&scheduleFile, "file", "f", "", "use file for schedule JSON")
	scheduleImportCmd.Flags().StringVarP(&datastore, "db", "d", "", "path to datastore")
}

var ScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "query or modify the schedule information",
	Long:  "query or modify the schedule information",
}

var scheduleResultsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "import results into the provided database",
	Long:  "import results into the provided database",
	Run: func(cmd *cobra.Command, args []string) {
		if scheduleYear == 0 || scheduleWeek == 0 {
			log.Fatal("year and week must be set via command line")
		}

		// Get a handle to the datastore
		if datastore == "" {
			log.Fatal("db flag is required")
		}

		db, err := sqlite3.NewDatastore(datastore)
		if err != nil {
			log.Fatal(err)
		}

		results, err := getResultsFromNFL(scheduleYear, scheduleWeek)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			for _, r := range results {
				fmt.Fprintln(os.Stderr, r)
			}
		}

		for _, r := range results {
			err := db.UpdateGame(int(scheduleWeek), int(scheduleYear), r.Home, r.HomeScore, r.AwayScore)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var scheduleResultsCmd = &cobra.Command{
	Use:   "results",
	Short: "downlaod schedule results from the NFL's website",
	Long:  "download schedule results from the NFL's website",
	Run: func(cmd *cobra.Command, args []string) {
		if scheduleYear == 0 || scheduleWeek == 0 {
			log.Fatal("year and week must be set via command line")
		}

		results, err := getResultsFromNFL(scheduleYear, scheduleWeek)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			for _, r := range results {
				fmt.Fprintln(os.Stderr, r)
			}
		}

		err = json.NewEncoder(os.Stdout).Encode(&results)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var scheduleDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download schedules from the NFL's website",
	Long:  "download schedules from the NFL's website",
	Run: func(cmd *cobra.Command, args []string) {
		if scheduleYear == 0 || scheduleWeek == 0 {
			log.Fatal("year and week must be set via command line")
		}

		games, err := getScheduleFromNFL(scheduleYear, scheduleWeek)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			for _, g := range games {
				fmt.Fprintln(os.Stderr, g)
			}
		}

		err = json.NewEncoder(os.Stdout).Encode(&games)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var scheduleImportCmd = &cobra.Command{
	Use:   "import",
	Short: "import schedule into a datastore",
	Long:  "import schedule into a datastore",
	Run: func(cmd *cobra.Command, args []string) {
		// Get a handle to the datastore
		if datastore == "" {
			log.Fatal("db flag is required")
		}

		db, err := sqlite3.NewDatastore(datastore)
		if err != nil {
			log.Fatal(err)
		}

		// Load a []schedule.Matchup from the NFL or a file
		var games []schedule.Matchup
		if scheduleFile != "" {
			fd, err := os.Open(scheduleFile)
			if err != nil {
				log.Fatal(err)
			}

			err = json.NewDecoder(fd).Decode(&games)
			if err != nil {
				log.Fatal(err)
			}
		} else if scheduleYear != 0 && scheduleWeek != 0 {
			games, err = getScheduleFromNFL(scheduleYear, scheduleWeek)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("either a file containing matchups or a year/week flag is required")
		}

		// Load the games into the datastore
		for _, g := range games {
			err := db.AddGame(g.Date, g.Home, g.Away)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

// getScheduleFromNFL creates a []schedule.Matchup from the NFL's website
// for the given week of the season.
func getScheduleFromNFL(year, week uint) ([]schedule.Matchup, error) {
	r, err := getScheduleHTML(year, week)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	p := schedule.NewParser(int(year), r)

	return p.Parse()
}

func getResultsFromNFL(year, week uint) ([]results.Result, error) {
	r, err := getScheduleHTML(year, week)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	p := results.NewParser(r)

	return p.Parse()
}

// getScheduleHTML returns an io.ReadCloser whose contents are the NFL website's
// schedule for the given year and week in HTML format. The ReadCloser MUST be
// closed by the caller.
func getScheduleHTML(year, week uint) (io.ReadCloser, error) {
	url := fmt.Sprintf("http://www.nfl.com/schedules/%d/REG%d", year, week)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
