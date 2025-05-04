package parsers

import (
	"testing"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestNormalService(t *testing.T) {
	service := testParseService(t, "service.xml")

	if len(service.ServiceParts) != 1 {
		t.Error("Wrong number of service parts")
	}

	if len(service.ServiceParts[0].Stops) != 48 {
		t.Error("Wrong number of stops")
	}
}

func TestCancelledService(t *testing.T) {
	service := testParseService(t, "service_cancelled.xml")

	if len(service.ServiceParts) != 1 {
		t.Error("Wrong number of service parts")
	}

	for index, stop := range service.ServiceParts[0].Stops {
		if index < 13 && stop.ArrivalCancelled {
			t.Errorf("Arrival should not be cancelled for stop %s", stop.Station.Code)
		} else if index >= 13 && !stop.ArrivalCancelled {
			t.Errorf("Arrival should be cancelled for stop %s", stop.Station.Code)
		}

		if index < 12 && stop.DepartureCancelled {
			t.Errorf("Departure should not be cancelled for stop %s", stop.Station.Code)
		} else if index == 12 && !stop.DepartureCancelled {
			t.Errorf("Departure should be cancelled for stop %s", stop.Station.Code)
		}
	}
}

func TestDelayedService(t *testing.T) {
	service := testParseService(t, "service_delay.xml")

	// PT3M14S
	if service.ServiceParts[0].Stops[0].DepartureDelay != 3*60+14 {
		t.Errorf("Wrong departure delay for stop %d, expected %d but got %d", 0, 3*60+14, service.ServiceParts[0].Stops[0].DepartureDelay)
	}

	hasDelayModification := false

	for _, modification := range service.ServiceParts[0].Stops[0].Modifications {
		if modification.ModificationType == models.ModificationDelayedDeparture {
			hasDelayModification = true
		}
	}

	if !hasDelayModification {
		t.Error("Departure delay modification missing")
	}

	// PT3M
	if service.ServiceParts[0].Stops[1].ArrivalDelay != 3*60 {
		t.Errorf("Wrong arrival delay for stop %d, expected %d but got %d", 0, 3*60, service.ServiceParts[0].Stops[1].ArrivalDelay)
	}

	hasDelayModification = false

	for _, modification := range service.ServiceParts[0].Stops[1].Modifications {
		if modification.ModificationType == models.ModificationDelayedArrival {
			hasDelayModification = true
		}
	}

	if !hasDelayModification {
		t.Error("Arrival delay modification missing")
	}
}

func TestInvalidService(t *testing.T) {
	_, err := ParseRitMessage(testFileReader(t, "invalid.xml"))

	if err == nil {
		t.Error("Should return an error for invalid XML")
	}

	_, err = ParseRitMessage(testFileReader(t, "dvs2/departure.xml"))

	if err == nil {
		t.Error("Should return an error for a Departure message")
	}
}

func testParseService(t *testing.T, name string) models.Service {
	service, err := ParseRitMessage(testFileReader(t, name))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	return service
}
