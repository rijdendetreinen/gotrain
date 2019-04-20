package models

import (
	"strconv"
	"testing"
)

func TestGetStoppingStations(t *testing.T) {
	var servicePart ServicePart

	for i := 0; i <= 10; i++ {
		var stop ServiceStop

		stop.Station.Code = "S" + strconv.Itoa(i)
		stop.Station.NameLong = "station " + strconv.Itoa(i)

		if i%2 == 0 {
			stop.StopType = "X"
		} else {
			stop.StopType = "D"
		}

		servicePart.Stops = append(servicePart.Stops, stop)
	}

	stopStations := servicePart.GetStoppingStations()

	// We expect 6 stations:
	if len(stopStations) != 6 {
		t.Fatal("Wrong number of stopping stations")
	}

	// Ensure that only StopType X is present:
	for _, stopStation := range stopStations {
		if stopStation.StopType == "D" {
			t.Error("Through station in GetStoppingStations")
		}
	}
}
