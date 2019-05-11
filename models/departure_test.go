package models

import (
	"testing"
	"time"
)

// TestRealDepartureTime tests the real departure time based on the planned time plus delay
func TestRealDepartureTime(t *testing.T) {
	var departure Departure

	departure.DepartureTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:34:00+01:00")

	expectedTime, _ := time.Parse(time.RFC3339, "2019-01-27T12:34:00+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}

	departure.Delay = 30 // 30 seconds
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:34:30+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}

	departure.Delay = 3600 // 1 hour
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T13:34:00+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}

	departure.Delay = -120 // Negative: -2 minutes
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:32:00+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}
}

func TestDeparturePlatformChanged(t *testing.T) {
	tables := []struct {
		planned string
		actual  string
		changed bool
	}{
		{"", "", false},
		{"4", "4", false},
		{"4", "5", true},
		{"", "4", true},
		{"4", "", true},
	}

	for _, table := range tables {
		var departure Departure
		departure.PlatformPlanned = table.planned
		departure.PlatformActual = table.actual

		if table.changed {
			if !departure.PlatformChanged() {
				t.Errorf("Planned platform %s is different from actual platform %s, but not changed", departure.PlatformPlanned, departure.PlatformActual)
			}
		} else {
			if departure.PlatformChanged() {
				t.Errorf("Planned platform %s is equal to actual platform %s, but is marked as changed", departure.PlatformPlanned, departure.PlatformActual)
			}
		}
	}
}

func TestGetDepartureID(t *testing.T) {
	var departure Departure

	departure.ServiceDate = "2019-01-27"
	departure.ServiceNumber = "301234"
	departure.ServiceID = "1234"
	departure.Station.Code = "RTD"
	departure.GenerateID()

	expected := "2019-01-27-1234-RTD"
	if departure.ID != expected {
		t.Errorf("Wrong departure ID, expected %s, got %s", expected, departure.ID)
	}
}
