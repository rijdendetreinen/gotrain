package models

import (
	"testing"
	"time"
)

// TestRealArrivalTime tests the real arrival time based on the planned time plus delay
func TestRealArrivalTime(t *testing.T) {
	var arrival Arrival

	arrival.ArrivalTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:34:00+01:00")

	expectedTime, _ := time.Parse(time.RFC3339, "2019-01-27T12:34:00+01:00")

	if !arrival.RealArrivalTime().Equal(expectedTime) {
		t.Errorf("Real arrival time %v does not match expected value %v", arrival.RealArrivalTime(), expectedTime)
	}

	arrival.Delay = 30 // 30 seconds
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:34:30+01:00")

	if !arrival.RealArrivalTime().Equal(expectedTime) {
		t.Errorf("Real arrival time %v does not match expected value %v", arrival.RealArrivalTime(), expectedTime)
	}

	arrival.Delay = 3600 // 1 hour
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T13:34:00+01:00")

	if !arrival.RealArrivalTime().Equal(expectedTime) {
		t.Errorf("Real arrival time %v does not match expected value %v", arrival.RealArrivalTime(), expectedTime)
	}

	arrival.Delay = -120 // Negative: -2 minutes
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:32:00+01:00")

	if !arrival.RealArrivalTime().Equal(expectedTime) {
		t.Errorf("Real arrival time %v does not match expected value %v", arrival.RealArrivalTime(), expectedTime)
	}
}

func TestArrivalPlatformChanged(t *testing.T) {
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
		var arrival Arrival
		arrival.PlatformPlanned = table.planned
		arrival.PlatformActual = table.actual

		if table.changed {
			if !arrival.PlatformChanged() {
				t.Errorf("Planned platform %s is different from actual platform %s, but not changed", arrival.PlatformPlanned, arrival.PlatformActual)
			}
		} else {
			if arrival.PlatformChanged() {
				t.Errorf("Planned platform %s is equal to actual platform %s, but is marked as changed", arrival.PlatformPlanned, arrival.PlatformActual)
			}
		}
	}
}
