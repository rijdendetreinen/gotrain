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
		if store.MessagesAverage != -1 {
			t.Errorf("Messages/minute is %f, should be -1", store.MessagesAverage)
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

		// We expect MessagesAverage (avg/minute) to be (4000 - 1000 = 3000) / 600 = 5 for the first round
		expected := float64(4000-1000-(i*150)) / 600

		if store.MessagesAverage != expected {
			t.Errorf("Messages/minute is %f, should be %f", store.MessagesAverage, expected)
		}
	}
}

func TestExtremeMeasurements(t *testing.T) {
	var store Store

	store.ResetCounters()
	store.ResetStatus()

	time1 := time.Date(2019, time.January, 1, 12, 0, 0, 0, time.UTC)
	store.Counters.Processed = 0
	store.newMeasurement(time1)

	time2 := time.Date(2019, time.January, 1, 13, 0, 0, 0, time.UTC)
	store.Counters.Processed = 10000
	store.newMeasurement(time2)

	if store.MessagesAverage != float64(10000)/3600 {
		t.Errorf("MessagesAverage = %f, expected %f", store.MessagesAverage, float64(10000)/3600)
	}
}

func TestStatusRecovery(t *testing.T) {
	var store Store

	store.ResetCounters()
	store.ResetStatus()

	// Typical downtime config for departures:
	store.DowntimeDetection.MinAverage = float64(1) / 60 // One message per minute
	store.DowntimeDetection.RecoveryTime = 70            // 70 mins recovery time

	// Overrule status (to counter ResetStatus(), which already sets it to UNKNOWN):
	store.Status = StatusUp

	// Status should be UNKNOWN
	time1 := time.Date(2019, time.January, 1, 12, 0, 0, 0, time.UTC)
	store.MessagesAverage = -1
	store.updateStatus(time1)

	if store.Status != StatusUnknown {
		t.Errorf("Status should be %s, but is %s", StatusUnknown, store.Status)
	}
	if !time1.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time1, store.LastStatusChange)
	}

	time1b := time.Date(2019, time.January, 1, 12, 5, 0, 0, time.UTC)
	store.MessagesAverage = -1
	store.updateStatus(time1b)

	if store.Status != StatusUnknown {
		t.Errorf("Status should be %s, but is %s", StatusUnknown, store.Status)
	}
	if !time1.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time1, store.LastStatusChange)
	}

	// Recovery starts now, status should change to RECOVERING
	time2 := time.Date(2019, time.January, 1, 12, 30, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time2)

	if store.Status != StatusRecovering {
		t.Errorf("Status should be %s, but is %s", StatusRecovering, store.Status)
	}
	if !time2.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time2, store.LastStatusChange)
	}

	// Currently 20 mins in recovery, status should still be recovering, LastStatusChange should not have changed
	time3 := time.Date(2019, time.January, 1, 12, 50, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time3)

	if store.Status != StatusRecovering {
		t.Errorf("Status should be %s, but is %s", StatusRecovering, store.Status)
	}
	if !time2.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time2, store.LastStatusChange)
	}

	// Currently 70 mins in recovery, status should change to UP
	time4 := time.Date(2019, time.January, 1, 13, 40, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time4)

	if store.Status != StatusUp {
		t.Errorf("Status should be %s, but is %s", StatusUp, store.Status)
	}
	if !time4.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time4, store.LastStatusChange)
	}
}

func TestStatusDown(t *testing.T) {
	var store Store

	store.ResetCounters()
	store.ResetStatus()

	store.DowntimeDetection.MinAverage = float64(1) / 60 // One message per minute
	store.DowntimeDetection.RecoveryTime = 40            // 40 mins recovery time

	// Store status is UP:
	store.Status = StatusUp

	// Status should be DOWN
	time1 := time.Date(2019, time.January, 1, 12, 0, 0, 0, time.UTC)
	store.MessagesAverage = 0
	store.updateStatus(time1)

	if store.Status != StatusDown {
		t.Errorf("Status should be %f, but is %f", store.DowntimeDetection.MinAverage, store.MessagesAverage)
		t.Errorf("Status should be %s, but is %s", StatusDown, store.Status)
	}
	if !time1.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time1, store.LastStatusChange)
	}

	// Recovery starts now, status should change to RECOVERING
	time2 := time.Date(2019, time.January, 1, 12, 30, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time2)

	if store.Status != StatusRecovering {
		t.Errorf("Status should be %s, but is %s", StatusRecovering, store.Status)
	}
	if !time2.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time2, store.LastStatusChange)
	}

	// Currently 20 mins in recovery, status should still be recovering, LastStatusChange should not have changed
	time3 := time.Date(2019, time.January, 1, 12, 50, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time3)

	if store.Status != StatusRecovering {
		t.Errorf("Status should be %s, but is %s", StatusRecovering, store.Status)
	}
	if !time2.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time2, store.LastStatusChange)
	}

	// Currently 40 mins in recovery, status should change to UP
	time4 := time.Date(2019, time.January, 1, 13, 10, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time4)

	if store.Status != StatusUp {
		t.Errorf("Status should be %s, but is %s", StatusUp, store.Status)
	}
	if !time4.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time4, store.LastStatusChange)
	}
}

func TestStatusRecoveryFailure(t *testing.T) {
	var store Store

	store.ResetCounters()
	store.ResetStatus()

	store.DowntimeDetection.MinAverage = float64(1) / 60 // One message per minute
	store.DowntimeDetection.RecoveryTime = 70            // 70 mins recovery time

	// Status should be DOWN
	time1 := time.Date(2019, time.January, 1, 12, 0, 0, 0, time.UTC)
	store.MessagesAverage = 0
	store.updateStatus(time1)

	if store.Status != StatusDown {
		t.Errorf("Status should be %s, but is %s", StatusDown, store.Status)
	}
	if !time1.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time1, store.LastStatusChange)
	}

	// Recovery starts now, status should change to RECOVERING
	time2 := time.Date(2019, time.January, 1, 12, 30, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time2)

	if store.Status != StatusRecovering {
		t.Errorf("Status should be %s, but is %s", StatusRecovering, store.Status)
	}
	if !time2.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time2, store.LastStatusChange)
	}

	// Currently 20 mins in recovery, status should still be recovering, LastStatusChange should not have changed
	time3 := time.Date(2019, time.January, 1, 12, 50, 0, 0, time.UTC)
	store.MessagesAverage = 10
	store.updateStatus(time3)

	if store.Status != StatusRecovering {
		t.Errorf("Status should be %s, but is %s", StatusRecovering, store.Status)
	}
	if !time2.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time2, store.LastStatusChange)
	}

	// Currently 60 mins in recovery, average messages dropped to 0.005.
	// Status should change to DOWN
	time4 := time.Date(2019, time.January, 1, 13, 30, 0, 0, time.UTC)
	store.MessagesAverage = 0.005
	store.updateStatus(time4)

	if store.Status != StatusDown {
		t.Errorf("Status should be %s, but is %s", StatusDown, store.Status)
	}
	if !time4.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time4, store.LastStatusChange)
	}

	// Currently 70 mins after recovery started, still not receiving:
	// Status should still be DOWN:
	time5 := time.Date(2019, time.January, 1, 13, 40, 0, 0, time.UTC)
	store.MessagesAverage = 0.0
	store.updateStatus(time5)

	if store.Status != StatusDown {
		t.Errorf("Status should be %s, but is %s", StatusDown, store.Status)
	}
	if !time4.Equal(store.LastStatusChange) {
		t.Errorf("Last status change should be %s, but is %s", time4, store.LastStatusChange)
	}
}
