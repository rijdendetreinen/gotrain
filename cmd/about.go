package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About GoTrain",
	Long:  `About GoTrain.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`GoTrain is a server application for receiving and processing real-time information about Dutch train services.
It is able to receive real-time departures, arrivals and service messages.

More information: https://github.com/rijdendetreinen/gotrain/`)
	},
}

func init() {
	RootCmd.AddCommand(aboutCmd)
}
