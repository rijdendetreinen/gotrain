package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rijdendetreinen/gotrain/stores"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Server status",
	Long:  `Get the current server status`,
	Run: func(cmd *cobra.Command, args []string) {
		baseURL, _ := cmd.Flags().GetString("url")

		if baseURL == "" {
			baseURL = viper.GetString("api.address")
		}

		timeout := time.Duration(5 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		url := baseURL + "/v2/status"

		response, err := client.Get(url)

		if err != nil {
			fmt.Printf("UNKNOWN - Error while getting status: %s\n", err)

			os.Exit(3)
		}

		if response.StatusCode != 200 {
			fmt.Printf("CRITICAL - Wrong status code: %d\n", response.StatusCode)

			os.Exit(2)
		}

		var status struct {
			Arrivals   string `json:"arrivals"`
			Departures string `json:"departures"`
			Services   string `json:"services"`
		}

		json.NewDecoder(response.Body).Decode(&status)

		if status.Arrivals == stores.StatusUp && status.Departures == stores.StatusUp && status.Services == stores.StatusUp {
			fmt.Printf("OK - Status: arrivals=%s, departures=%s, services=%s\n", status.Arrivals, status.Departures, status.Services)
			os.Exit(0)
		}

		if status.Arrivals == stores.StatusDown || status.Departures == stores.StatusDown || status.Services == stores.StatusDown {
			fmt.Printf("CRITICAL - Status: arrivals=%s, departures=%s, services=%s\n", status.Arrivals, status.Departures, status.Services)
			os.Exit(2)
		}

		if status.Arrivals == stores.StatusUnknown || status.Departures == stores.StatusUnknown || status.Services == stores.StatusUnknown {
			fmt.Printf("CRITICAL - Status: arrivals=%s, departures=%s, services=%s\n", status.Arrivals, status.Departures, status.Services)
			os.Exit(2)
		}

		if status.Arrivals == stores.StatusRecovering || status.Departures == stores.StatusRecovering || status.Services == stores.StatusRecovering {
			fmt.Printf("WARNING - Status: arrivals=%s, departures=%s, services=%s\n", status.Arrivals, status.Departures, status.Services)
			os.Exit(1)
		}

		os.Exit(2)
	},
}

func init() {
	RootCmd.AddCommand(statusCommand)
	statusCommand.Flags().StringP("url", "u", "", "Server URL")
}
