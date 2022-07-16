package parsers

import (
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestParseNormalDeparture(t *testing.T) {
	departure := testParseDeparture(t, "departure.xml")

	if departure.Cancelled {
		t.Error("Departure should not be cancelled")
	}

	expectedDepartureTime, _ := time.Parse(time.RFC3339, "2019-04-06T22:44:00+01:00")

	if !departure.DepartureTime.Equal(expectedDepartureTime) {
		t.Errorf("Expected departure time %v does not match %v", expectedDepartureTime, departure.DepartureTime)
	}

	if departure.DestinationActual[0].Code != "RHN" {
		t.Errorf("Expected destination %s does not match %s", "RHN", departure.DestinationActual[0].Code)
	}

	if departure.PlatformActual != "2" {
		t.Errorf("Expected platform %s does not match %s", "2", departure.PlatformActual)
	}

	if len(departure.ViaActual) != 3 {
		t.Error("Wrong number of via stations")
	}

	if len(departure.TrainWings) != 1 {
		t.Error("Wrong number of train wings")
	}

	if len(departure.TrainWings[0].Stations) != 6 {
		t.Error("Wrong number of wing stations")
	}
}

func TestCancelledDeparture(t *testing.T) {
	departure := testParseDeparture(t, "departure_cancelled.xml")

	if !departure.Cancelled {
		t.Error("Departure should be cancelled")
	}

	if len(departure.TrainWings[0].Stations) != 1 {
		t.Error("Wrong number of actual stations")
	}

	if len(departure.TrainWings[0].StationsPlanned) != 5 {
		t.Error("Wrong number of planned stations")
	}
}

func TestDelayedDeparture(t *testing.T) {
	departure := testParseDeparture(t, "departure_delay.xml")

	if departure.Cancelled {
		t.Error("Departure should not be cancelled")
	}

	if departure.Delay != 63 {
		t.Error("Wrong amount of delay")
	}

	hasDelayModification := false

	for _, modification := range departure.Modifications {
		if modification.ModificationType == models.ModificationDelayedDeparture {
			hasDelayModification = true
		}
	}

	if !hasDelayModification {
		t.Error("Delay modification missing")
	}
}

func TestDepartureTravelTips(t *testing.T) {
	departure := testParseDeparture(t, "departure_travel-tips.xml")

	if len(departure.TravelTips) != 2 {
		t.Error("Wrong number of travel tips")
	}
}

func TestDepartureBoardingTips(t *testing.T) {
	departure := testParseDeparture(t, "departure_boarding-tips.xml")

	if len(departure.BoardingTips) != 1 {
		t.Error("Wrong number of boarding tips")
	}
}

func TestDepartureNotRealtime(t *testing.T) {
	departure := testParseDeparture(t, "departure_not-realtime.xml")

	if !departure.NotRealTime {
		t.Error("Departure should be flagged as NotRealTime")
	}
}

func TestDepartureTrainName(t *testing.T) {
	departure := testParseDeparture(t, "departure_train-name.xml")

	if departure.ServiceName != "Spoorwegmuseum" {
		t.Errorf("Train name should be '%s', but is '%s'", "Spoorwegmuseum", departure.ServiceName)
	}
}

func TestDepartureModification(t *testing.T) {
	departure := testParseDeparture(t, "departure_modification-cause.xml")

	if len(departure.Modifications) != 1 {
		t.Error("Wrong number of modifications")
	}

	foundModification := false

	for _, modification := range departure.TrainWings[0].Modifications {
		if modification.ModificationType == models.ModificationRouteShortened {
			foundModification = true
			expectedShort := "herstelwerkzaamheden"
			expectedLong := "door herstelwerkzaamheden"
			expectedStation := "VNDC"

			if modification.CauseShort != expectedShort {
				t.Errorf("Wrong CauseShort for modification, expected '%s', but got '%s'", expectedShort, modification.CauseShort)
			}

			if modification.CauseLong != expectedLong {
				t.Errorf("Wrong CauseLong for modification, expected '%s', but got '%s'", expectedLong, modification.CauseLong)
			}

			if modification.Station.Code == "" {
				t.Error("Should have a station for this modification")
			} else if modification.Station.Code != expectedStation {
				t.Errorf("Wrong Station.Code for modification, expected '%s', but got '%s'", expectedStation, modification.Station.Code)
			}
		}
	}

	if !foundModification {
		t.Error("Did not find modification")
	}
}

func TestInvalidDeparture(t *testing.T) {
	_, err := ParseDvsMessage(testFileReader(t, "invalid.xml"))

	if err == nil {
		t.Error("Should return an error for invalid XML")
	}

	_, err = ParseDvsMessage(testFileReader(t, "arrival.xml"))

	if err == nil {
		t.Error("Should return an error for an Arrival message")
	}
}

func TestDepartureMultiplePlatforms(t *testing.T) {
	departure := testParseDeparture(t, "departure_multiple-platforms.xml")

	if departure.PlatformActual != "5/6" {
		t.Errorf("Wrong platform: expected '%s', but got '%s'", "5/6", departure.PlatformActual)
	}
}

func testParseDeparture(t *testing.T, name string) models.Departure {
	departure, err := ParseDvsMessage(testFileReader(t, name))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	return departure
}

func TestParseMaterialLeftBehind(t *testing.T) {
	departure := testParseDeparture(t, "departure_material-left-behind.xml")

	if departure.TrainWings[0].Material[0].RemainsBehind != true {
		t.Errorf("Wrong Material.RemainsBehind: expected '%v', but got '%v'", true, departure.TrainWings[0].Material[0].RemainsBehind)
	}
}

func TestParseMaterialModifications(t *testing.T) {
	departure := testParseDeparture(t, "departure_material-added.xml")

	if departure.TrainWings[0].Material[0].Added != true {
		t.Errorf("Wrong Material.RemainsBehind: expected '%v', but got '%v'", true, departure.TrainWings[0].Material[0].Added)
	}
}
