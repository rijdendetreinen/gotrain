package stores

import (
	"testing"
	"time"
)

func TestMeasurements(t *testing.T) {
	var store Store

	store.ResetCounters()
	store.ResetStatus()

	// Verify that we have no measurements yet:
	if len(store.measurements) != 0 {
		t.Fatal("Measurements in store, cannot test measurements")
	}

	if store.Status != StatusUnknown {
		t.Fatalf("Store status should be %s", StatusUnknown)
	}

	t.Log("Store our first measurement on 2019-01-01 12:00:00, and continue from there on:")

	// Store our second and further measurements on 2019-01-01 12:xx:xx (every 30s)
	// Increment Processed counter every time by 100 messages, until
	for i := 0; i < 20; i++ {
		time := time.Date(2019, time.January, 1, 12, 0, i*30, 0, time.UTC)
		t.Logf("Store measurement on %v, 150 msg processed", time)
		store.Counters.Processed = 1000 + i*150
		store.newMeasurement(time)

		if len(store.measurements) != i+1 {
			t.Error("Measurement not stored")
		}
		if store.messagesAverage != 0 {
			t.Errorf("Messages/minute is %f, should be 0", store.messagesAverage)
		}
		if store.Status != StatusUnknown {
			t.Fatalf("Store status should be %s", StatusUnknown)
		}
	}

	for i := 0; i < 20; i++ {
		time := time.Date(2019, time.January, 1, 12, 10, i*30, 0, time.UTC)

		store.Counters.Processed = 4000
		if i == 0 {
			t.Logf("Store measurement on %v, 150 msg processed", time)
		} else {
			t.Logf("Store measurement on %v, 0 msg processed", time)
		}

		store.newMeasurement(time)

		// We expect messagesAverage (avg/minute) to be (4000 - 1000 = 3000) / 600 = 5 for the first round
		expected := float64(4000-1000-(i*150)) / 600

		if store.messagesAverage != expected {
			t.Errorf("Messages/minute is %f, should be %f", store.messagesAverage, expected)
		}
	}
}
