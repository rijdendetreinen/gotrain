package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rijdendetreinen/gotrain/parsers"
)

func main() {
	filename := os.Args[1]

	fmt.Printf("Opening %s\n", filename)

	f, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	service := parsers.ParseRitMessage(f)

	fmt.Println("")
	fmt.Printf("%+v\n", service)
	fmt.Println("")

	loc, _ := time.LoadLocation("Europe/Amsterdam")

	for index, part := range service.ServiceParts {
		fmt.Printf("  ** Service part %d  service=%s\n", index+1, part.ServiceNumber)

		for stopIndex, stop := range part.Stops {
			fmt.Printf("    ** Stop %02d %7s = %s\n", stopIndex+1, stop.Station.Code, stop.Station.NameLong)
			if stop.ArrivalTime != nil {
				fmt.Printf("       A: %s +%d\n", stop.ArrivalTime.In(loc).Format("15:04"), stop.ArrivalDelay)
			}
			if stop.DepartureTime != nil {
				fmt.Printf("       V: %s +%d\n", stop.DepartureTime.In(loc).Format("15:04"), stop.DepartureDelay)
			}
			if len(stop.Material) > 0 {
				fmt.Print("       Material: ")

				for _, material := range stop.Material {
					fmt.Printf("%s[%s]>%s ", material.NaterialType, material.Number, material.DestinationActual.Code)
				}

				fmt.Print("\n")
			}
			// fmt.Printf("    service number: %s\n", stop.Station.NameLong)
		}
	}

	f.Close()
}
