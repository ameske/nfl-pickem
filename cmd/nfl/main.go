package main

import (
	"log"

	"github.com/spf13/cobra"
)

var datastore string
var verbose bool

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var rootCmd = &cobra.Command{Use: "nfl"}
	rootCmd.AddCommand(ScheduleCmd)
	rootCmd.AddCommand(TestCmd)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&datastore, "db", "d", "", "path to datastore")
	rootCmd.Execute()
}
