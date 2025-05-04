package parsers

import (
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestNormalArrival(t *testing.T) {
	arrival := testParseArrival(t, "arrival.xml")

	if arrival.Cancelled {
		t.Error("Arrival should not be cancelled")
	}

	expectedArrivalTime, _ := time.Parse(time.RFC3339, "2018-09-04T09:30:00+02:00")

	if !arrival.ArrivalTime.Equal(expectedArrivalTime) {
		t.Errorf("Expected arrival time %v does not match %v", expectedArrivalTime, arrival.ArrivalTime)
	}

	if arrival.OriginActual[0].Code != "GVC" {
		t.Errorf("Expected destination %s does not match %s", "GVC", arrival.OriginActual[0].Code)
	}

	if arrival.PlatformActual != "12" {
		t.Errorf("Expected platform %s does not match %s", "2", arrival.PlatformActual)
	}

	if len(arrival.ViaActual) != 1 {
		t.Error("Wrong number of via stations")
	}
}

func TestCancelledArrival(t *testing.T) {
	arrival := testParseArrival(t, "arrival_cancelled.xml")

	if !arrival.Cancelled {
		t.Error("Arrival should be cancelled")
	}
}

func TestDelayedArrival(t *testing.T) {
	arrival := testParseArrival(t, "arrival_delay.xml")

	if arrival.Cancelled {
		t.Error("Arrival should not be cancelled")
	}

	// PT6M39S
	if arrival.Delay != 6*60+39 {
		t.Error("Wrong amount of delay")
	}

	hasDelayModification := false

	for _, modification := range arrival.Modifications {
		if modification.ModificationType == models.ModificationDelayedArrival {
			hasDelayModification = true
		}
	}

	if !hasDelayModification {
		t.Error("Delay modification missing")
	}
}

func TestArrivalTrainName(t *testing.T) {
	arrival := testParseArrival(t, "arrival_train-name.xml")

	if arrival.ServiceName != "Spoorwegmuseum" {
		t.Errorf("Train name should be '%s', but is '%s'", "Spoorwegmuseum", arrival.ServiceName)
	}
}

func TestArrivalModification(t *testing.T) {
	arrival := testParseArrival(t, "arrival_modification-cause.xml")

	if len(arrival.Modifications) != 2 {
		t.Error("Wrong number of modifications")
	}

	foundModification := false

	for _, modification := range arrival.Modifications {
		if modification.ModificationType == models.ModificationDiverted {
			foundModification = true
			expectedShort := "wisselstoring"
			expectedLong := "door een wisselstoring"

			if modification.CauseShort != expectedShort {
				t.Errorf("Wrong CauseShort for modification, expected '%s', but got '%s'", expectedShort, modification.CauseShort)
			}

			if modification.CauseLong != expectedLong {
				t.Errorf("Wrong CauseLong for modification, expected '%s', but got '%s'", expectedLong, modification.CauseLong)
			}

			if modification.Station.Code != "" {
				t.Error("Should not have a station for this modification")
			}
		}
	}

	if !foundModification {
		t.Error("Did not find modification")
	}
}

func TestInvalidArrival(t *testing.T) {
	_, err := ParseDasMessage(testFileReader(t, "invalid.xml"))

	if err == nil {
		t.Error("Should return an error for invalid XML")
	}

	_, err = ParseDasMessage(testFileReader(t, "dvs2/departure.xml"))

	if err == nil {
		t.Error("Should return an error for a Departure message")
	}
}

func TestArrivalMultiplePlatforms(t *testing.T) {
	arrival := testParseArrival(t, "arrival_multiple-platforms.xml")

	if arrival.PlatformActual != "5/6" {
		t.Errorf("Wrong platform: expected '%s', but got '%s'", "5/6", arrival.PlatformActual)
	}
}

func testParseArrival(t *testing.T, name string) models.Arrival {
	arrival, err := ParseDasMessage(testFileReader(t, name))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	return arrival
}
