package main

import (
	"log"

	"github.com/spf13/cobra"
)

var datastore string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var rootCmd = &cobra.Command{Use: "nfl"}
	rootCmd.AddCommand(ScheduleCmd)
	rootCmd.AddCommand(TestDBCmd)
	rootCmd.Execute()
}
