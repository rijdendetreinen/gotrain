package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rijdendetreinen/gotrain/models"
	"github.com/rijdendetreinen/gotrain/parsers"
	"github.com/spf13/cobra"
)

var inspectCommand = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect messages",
	Long:  `Inspect XML messages. Use a sub-command to specify the message type.`,
}

var inspectDepartureCommand = &cobra.Command{
	Use:   "departure [filename]",
	Short: "Inspect a departure message",
	Long:  `Inspect a departure XML message and print a summary of the content to the screen.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		showModifications, _ := cmd.Flags().GetBool("modifications")
		language, _ := cmd.Flags().GetString("language")

		f := openFile(args)

		departure, err := parsers.ParseDvsMessage(f)

		if err != nil {
			fmt.Println("Error while parsing departure")
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Printf("Product ID: %s\n", departure.ProductID)
		fmt.Printf("Timestamp: %s\n", departure.Timestamp.Local())
		fmt.Printf("Departure ID: %s\n", departure.ID)
		fmt.Printf("Service ID: %s\n", departure.ServiceID)
		fmt.Printf("Departure station: %s = %s\n", departure.Station.Code, departure.Station.NameLong)
		fmt.Printf("Service number: %s\n", departure.ServiceNumber)
		fmt.Printf("Service date: %s\n", departure.ServiceDate)
		fmt.Printf("Departure time: %s\n", departure.DepartureTime.Local())
		fmt.Printf("Delay: %ds\n", departure.Delay)
		fmt.Printf("Status: %d\n", departure.Status)
		fmt.Printf("Real departure time: %s\n", departure.RealDepartureTime().Local())
		fmt.Print("Actual destination(s): ")
		displayStations(departure.DestinationActual)

		fmt.Print("\nPlanned destination(s): ")
		displayStations(departure.DestinationPlanned)

		fmt.Printf("\nType: %s/%s\n", departure.ServiceTypeCode, departure.ServiceType)
		fmt.Printf("Company: %v\n", departure.Company)

		fmt.Printf("DoNotBoard: %v\n", departure.DoNotBoard)
		fmt.Printf("NotRealTime: %v\n", departure.NotRealTime)
		fmt.Printf("RearPartRemains: %v\n", departure.RearPartRemains)
		fmt.Printf("ReservationRequired: %v\n", departure.ReservationRequired)
		fmt.Printf("SpecialTicket: %v\n", departure.SpecialTicket)
		fmt.Printf("WithSupplement: %v\n", departure.WithSupplement)

		fmt.Print("Actual route station(s): ")
		displayStations(departure.ViaActual)
		fmt.Print("\nPlanned route station(s): ")
		displayStations(departure.ViaPlanned)
		fmt.Print("\n")

		fmt.Printf("Boarding tips: %v\n", departure.BoardingTips)
		fmt.Printf("Travel tips: %v\n", departure.TravelTips)
		fmt.Printf("Change tips: %v\n", departure.ChangeTips)

		fmt.Println("Service modifications:")
		displayModifications(departure.Modifications, 1, showModifications, language)
	},
}

var inspectServiceCommand = &cobra.Command{
	Use:   "service [filename]",
	Short: "Inspect a service message",
	Long:  `Inspect a service XML message and print a summary of the content to the screen.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		showModifications, _ := cmd.Flags().GetBool("modifications")
		language, _ := cmd.Flags().GetString("language")
		showStops, _ := cmd.Flags().GetBool("stops")

		f := openFile(args)

		service, err := parsers.ParseRitMessage(f)

		if err != nil {
			fmt.Println("Error while parsing service")
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Printf("Product ID: %s\n", service.ProductID)
		fmt.Printf("Timestamp: %s\n", service.Timestamp.Local())
		fmt.Printf("Validity: %s\n", service.ValidUntil.Local())
		fmt.Printf("Service ID: %s\n", service.ID)
		fmt.Printf("Service number: %s\n", service.ServiceNumber)
		fmt.Printf("Service date: %s\n", service.ServiceDate)
		fmt.Printf("Type: %s/%s\n", service.ServiceTypeCode, service.ServiceType)
		fmt.Printf("Company: %v\n", service.Company)
		fmt.Printf("JourneyPlanner: %v\n", service.JourneyPlanner)
		fmt.Printf("ReservationRequired: %v\n", service.ReservationRequired)
		fmt.Printf("SpecialTicket: %v\n", service.SpecialTicket)
		fmt.Printf("WithSupplement: %v\n", service.WithSupplement)

		fmt.Println("Service parts:")

		for index, part := range service.ServiceParts {
			fmt.Printf("  ** Service part %d  service=%s\n", index+1, part.ServiceNumber)

			if showStops {
				for stopIndex, stop := range part.Stops {
					fmt.Printf("    ** Stop %02d %7s = %s\n", stopIndex+1, stop.Station.Code, stop.Station.NameLong)
					if !stop.ArrivalTime.IsZero() {
						fmt.Printf("       A: %s +%d\n", stop.ArrivalTime.Local().Format("15:04"), stop.ArrivalDelay)
					}
					if !stop.DepartureTime.IsZero() {
						fmt.Printf("       V: %s +%d\n", stop.DepartureTime.Local().Format("15:04"), stop.DepartureDelay)
					}
					if len(stop.Material) > 0 {
						fmt.Print("       Material: ")

						for _, material := range stop.Material {
							fmt.Printf("%s[%s]>%s ", material.NaterialType, material.Number, material.DestinationActual.Code)
						}

						fmt.Print("\n")
					}
					fmt.Println("       Stop modifications:")
					displayModifications(stop.Modifications, 7, showModifications, language)
				}
			} else {
				fmt.Printf("     %d stop(s)\n", len(part.Stops))
			}

			fmt.Println("     Service part modifications:")
			displayModifications(part.Modifications, 5, showModifications, language)
		}

		fmt.Println("Service modifications:")
		displayModifications(service.Modifications, 1, showModifications, language)
	},
}

func displayModifications(modifications []models.Modification, level int, showModifications bool, language string) {
	if showModifications {
		if len(modifications) == 0 {
			fmt.Printf("%sno modifications\n", strings.Repeat(" ", level))
		}
		for index, modification := range modifications {
			remark, _ := modification.Remark(language)
			fmt.Printf("%s%d. [%d] %v\n", strings.Repeat(" ", level), index, modification.ModificationType, remark)
		}
	} else {
		fmt.Printf("%s%d modifications(s)\n", strings.Repeat(" ", level), len(modifications))
	}
}

func displayStations(stations []models.Station) {
	if len(stations) == 0 {
		fmt.Printf("none")
	} else {
		for index, station := range stations {
			if index > 0 {
				fmt.Print("; ")
			}
			fmt.Printf("%s = %s", station.Code, station.NameLong)
		}
	}
}

func openFile(args []string) *os.File {
	filename := args[0]

	f, err := os.Open(filename)

	if err != nil {
		fmt.Printf("Error opening %s", filename)
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s:\n", filename)

	return f
}

func init() {
	RootCmd.AddCommand(inspectCommand)
	inspectCommand.AddCommand(inspectDepartureCommand)
	inspectCommand.AddCommand(inspectServiceCommand)

	inspectDepartureCommand.Flags().BoolP("modifications", "m", false, "Show modifications")
	inspectDepartureCommand.Flags().BoolP("stops", "s", false, "Show stops")
	inspectDepartureCommand.Flags().StringP("language", "l", "nl", "Language")

	inspectServiceCommand.Flags().BoolP("modifications", "m", false, "Show modifications")
	inspectServiceCommand.Flags().BoolP("stops", "s", false, "Show stops")
	inspectServiceCommand.Flags().StringP("language", "l", "nl", "Language")
}
