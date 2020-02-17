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
			stop.StoppingActual = true
		} else {
			stop.StoppingActual = false
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
		if stopStation.StoppingActual == false && stopStation.StoppingPlanned == false {
			t.Error("Through station in GetStoppingStations")
		}
	}
}

func TestGetStops(t *testing.T) {
	var service Service
	var servicePart ServicePart

	for i := 0; i <= 10; i++ {
		var stop ServiceStop

		stop.Station.Code = "S" + strconv.Itoa(i)
		stop.Station.NameLong = "station " + strconv.Itoa(i)
		stop.StoppingPlanned = true
		servicePart.Stops = append(servicePart.Stops, stop)
	}

	service.ServiceParts = append(service.ServiceParts, servicePart)

	stops := service.GetStops()

	for i := 0; i <= 10; i++ {
		_, exists := stops["S"+strconv.Itoa(i)]
		if !exists {
			t.Errorf("Station S %d not returned", i)
		}
	}
}

func TestGetServiceID(t *testing.T) {
	var service Service

	service.ServiceDate = "2019-01-27"
	service.ServiceNumber = "12345"
	service.GenerateID()

	expected := "2019-01-27-12345"
	if service.ID != expected {
		t.Errorf("Wrong service ID, expected %s, got %s", expected, service.ID)
	}
}
