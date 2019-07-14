package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "GoTrain version",
	Long:  `Version information about GoTrain`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoTrain %v (%v), built: %v\n", Version.Version, Version.Commit, Version.Date)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
