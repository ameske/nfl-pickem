package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ameske/nfl-pickem/parser/schedule"
	"github.com/spf13/cobra"
)

var year, week uint
var verbose bool

func init() {
	ScheduleCmd.AddCommand(scheduleDownloadCmd)
	scheduleDownloadCmd.Flags().UintVarP(&year, "year", "y", 0, "NFL season year")
	scheduleDownloadCmd.Flags().UintVarP(&week, "week", "w", 0, "NFL season week")
	scheduleDownloadCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output to stderr")
}

var ScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "query or modify the schedule information",
	Long:  "query or modify the schedule information",
}

var scheduleDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download schedules from the NFL's website",
	Long:  "download schedules from the NFL's website",
	Run: func(cmd *cobra.Command, args []string) {
		if year == 0 || week == 0 {
			log.Fatal("year and week must be set via command line")
		}

		r, err := getScheduleHTML(year, week)
		if err != nil {
			log.Fatal(err)
		}

		defer r.Close()

		p := schedule.NewParser(int(year), r)

		games, err := p.Parse()
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
