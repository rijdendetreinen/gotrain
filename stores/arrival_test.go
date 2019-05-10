package stores

import (
	"strconv"
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestArrivalsCount(t *testing.T) {
	var store ArrivalStore
	store.InitStore()

	if store.GetNumberOfArrivals() != 0 {
		t.Error("Wrong number of arrivals")
	}

	store.ProcessArrival(generateArrival())

	if store.GetNumberOfArrivals() != 1 {
		t.Error("Wrong number of arrivals")
	}

	if store.GetNumberOfArrivals() != len(store.GetAllArrivals()) {
		t.Error("Reported number of arrivals does not match with actual inventory count")
	}
}

func TestRetrieveArrival(t *testing.T) {
	var store ArrivalStore

	arrival := generateArrival()

	store.InitStore()
	store.ProcessArrival(arrival)

	arrivalInStore := store.GetArrival("1234", "2019-01-27", "UT")

	if arrivalInStore == nil {
		t.Error("Could not retrieve arrival from store")
	}
}

func TestDuplicateArrival(t *testing.T) {
	var store ArrivalStore

	arrival := generateArrival()

	store.InitStore()
	store.ProcessArrival(arrival)

	if store.Counters.Duplicates != 0 {
		t.Fatal("Wrong number of duplicates for counters")
	}

	// Process again (forcing duplicate)
	store.ProcessArrival(arrival)

	if store.Counters.Duplicates != 1 {
		t.Error("Should increment counter for duplicates")
	}
}

func TestArrivalsProcessing(t *testing.T) {
	var store ArrivalStore

	arrival := generateArrival()

	store.InitStore()
	store.ProcessArrival(arrival)

	// Older:
	arrival2 := arrival

	// Earlier than previous message, so should be ignored:
	arrival2.ProductID = "12344"
	arrival2.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 68, time.UTC)

	store.ProcessArrival(arrival2)
	arrivalInStore := store.GetArrival("1234", "2019-01-27", "UT")

	if arrivalInStore.ProductID != "12345" {
		t.Error("Should not update arrival with earlier arrival")
	}
	if store.Counters.Outdated != 1 {
		t.Error("Should increase counter for outdated messages")
	}

	arrival3 := arrival
	arrival3.ProductID = "12343"
	arrival3.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 98, time.UTC)

	store.ProcessArrival(arrival3)
	arrivalInStore = store.GetArrival("1234", "2019-01-27", "UT")

	if arrivalInStore.ProductID != "12343" {
		t.Error("Should update arrival with later message")
	}
}

func TestRetrieveArrivalsFromStation(t *testing.T) {
	var store ArrivalStore

	store.InitStore()

	// Fake some arrivals
	arrival1 := generateArrival()

	arrival2 := generateArrival()
	arrival2.ServiceID = "54321"
	arrival2.GenerateID()

	arrival3 := generateArrival()
	arrival3.ServiceID = "23456"
	arrival3.Hidden = true
	arrival3.GenerateID()

	arrival4 := generateArrival()
	arrival4.ServiceID = "23456"
	arrival4.Station.Code = "UTO"
	arrival4.Station.NameLong = "Utrecht Overvecht"
	arrival4.GenerateID()

	store.ProcessArrival(arrival1)
	store.ProcessArrival(arrival2)
	store.ProcessArrival(arrival3)
	store.ProcessArrival(arrival4)

	arrivalsUT := store.GetStationArrivals("UT", false)

	if len(arrivalsUT) != 2 {
		t.Error("Wrong number of arrivals for UT")
	}

	arrivalsUT = store.GetStationArrivals("UT", true)

	if len(arrivalsUT) != 3 {
		t.Error("Wrong number of (hidden+non-hidden) arrivals for UT")
	}

	arrivalsUTO := store.GetStationArrivals("UTO", false)

	if len(arrivalsUTO) != 1 {
		t.Error("Wrong number of arrivals for UTO")
	}

	// Station without arrivals:
	arrivalsHRY := store.GetStationArrivals("HRY", false)

	if len(arrivalsHRY) != 0 {
		t.Error("Wrong number of arrivals for HRY")
	}
}

func generateArrival() models.Arrival {
	var arrival models.Arrival

	arrival.ProductID = "12345"
	arrival.ServiceID = "1234"
	arrival.ServiceNumber = "1234"
	arrival.Station.Code = "UT"
	arrival.Station.NameLong = "Utrecht Centraal"
	arrival.ServiceDate = "2019-01-27"
	arrival.GenerateID()
	arrival.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)
	arrival.ArrivalTime = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)

	return arrival
}

func TestCleanupArrivals(t *testing.T) {
	var store ArrivalStore

	store.InitStore()

	// Fake some arrivals
	arrival1 := generateArrival()
	arrival1.Hidden = true

	arrival2 := generateArrival()
	arrival2.ServiceID = "54321"
	arrival2.GenerateID()

	arrival3 := generateArrival()
	arrival3.ServiceID = "99999"
	arrival3.ArrivalTime = time.Date(2099, time.January, 27, 12, 34, 56, 78, time.UTC)
	arrival3.GenerateID()

	arrival4 := generateArrival()
	arrival4.ServiceID = "66666"
	arrival4.GenerateID()
	arrival4.ArrivalTime = time.Date(2019, time.January, 27, 11, 34, 56, 78, time.UTC)

	store.ProcessArrival(arrival1)
	store.ProcessArrival(arrival2)
	store.ProcessArrival(arrival3)
	store.ProcessArrival(arrival4)

	// Verify that we have 4 arrivals in store:
	if store.GetNumberOfArrivals() != 4 {
		t.Error("Wrong number of arrivals")
	}

	// Cleanup, first pass:
	store.CleanUp(time.Date(2019, time.January, 27, 12, 30, 56, 78, time.UTC))

	// Nothing should be removed just yet (including the hidden arrival, whose arrival time
	// is in the future now, and should stay for at least 4 hours)

	// Verify that we have 4 arrivals in store:
	if store.GetNumberOfArrivals() != 4 {
		t.Error("Wrong number of arrivals")
	}

	// Cleanup at 12:35
	store.CleanUp(time.Date(2019, time.January, 27, 12, 36, 56, 78, time.UTC))

	if store.GetArrival(arrival4.ServiceID, arrival4.ServiceDate, arrival4.Station.Code).Hidden == false {
		t.Error("Old arrival should be hidden after CleanUp #1")
	}
	if store.GetNumberOfArrivals() != 4 {
		t.Error("Wrong number of arrivals")
	}

	// Now cleanup 31 mins after planned arrival, arrival2 should be hidden too now
	store.CleanUp(time.Date(2019, time.January, 27, 13, 05, 56, 78, time.UTC))

	if store.GetArrival(arrival4.ServiceID, arrival4.ServiceDate, arrival4.Station.Code).Hidden == false {
		t.Error("Old arrival should be hidden after CleanUp #2")
	}
	if store.GetArrival(arrival2.ServiceID, arrival2.ServiceDate, arrival2.Station.Code).Hidden == false {
		t.Error("Arrival should be hidden after 30 mins")
	}
	if store.GetNumberOfArrivals() != 4 {
		t.Error("Wrong number of arrivals")
	}

	// Cleanup #3, everything must be removed
	store.CleanUp(time.Date(2019, time.January, 27, 16, 45, 56, 78, time.UTC))

	// The hidden arrival should be gone by now. The second arrival should be hidden by now.
	// The third arrival should still be visible.
	if store.GetNumberOfArrivals() != 1 {
		t.Fatal("Hidden arrivals not removed")
	}

	// Verify arrival3 is still visible.
	if store.GetArrival(arrival3.ServiceID, arrival3.ServiceDate, arrival3.Station.Code).Hidden == true {
		t.Error("Train which arrives in 2099 should not be hidden already")
	}
}

func TestSaveArrivalStore(t *testing.T) {
	var store, store2 ArrivalStore

	store.InitStore()

	for i := 0; i < 40000; i++ {
		arrival := generateArrival()
		arrival.ServiceID = strconv.Itoa(i)
		arrival.GenerateID()

		store.ProcessArrival(arrival)
	}

	if store.GetNumberOfArrivals() != 40000 {
		t.Error("Wrong number of arrivals")
	}

	// Save
	error := store.SaveStore()

	if error != nil {
		t.Fatalf("%s", error)
	}

	// Load in empty store:
	store2.InitStore()
	store2.ReadStore()

	if store2.GetNumberOfArrivals() != 40000 {
		t.Error("Wrong number of arrivals")
	}
}
