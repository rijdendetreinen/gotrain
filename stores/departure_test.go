package stores

import (
	"strconv"
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestDeparturesCount(t *testing.T) {
	var store DepartureStore
	store.InitStore()

	if store.GetNumberOfDepartures() != 0 {
		t.Error("Wrong number of departures")
	}

	store.ProcessDeparture(generateDeparture())

	if store.GetNumberOfDepartures() != 1 {
		t.Error("Wrong number of departures")
	}

	if store.GetNumberOfDepartures() != len(store.GetAllDepartures()) {
		t.Error("Reported number of departures does not match with actual inventory count")
	}
}

func TestRetrieveDeparture(t *testing.T) {
	var store DepartureStore

	departure := generateDeparture()

	store.InitStore()
	store.ProcessDeparture(departure)

	departureInStore := store.GetDeparture("1234", "2019-01-27", "UT")

	if departureInStore == nil {
		t.Error("Could not retrieve departure from store")
	}
}

func TestDuplicateDeparture(t *testing.T) {
	var store DepartureStore

	departure := generateDeparture()

	store.InitStore()
	store.ProcessDeparture(departure)

	if store.Counters.Duplicates != 0 {
		t.Fatal("Wrong number of duplicates for counters")
	}

	// Process again (forcing duplicate)
	store.ProcessDeparture(departure)

	if store.Counters.Duplicates != 1 {
		t.Error("Should increment counter for duplicates")
	}
}

func TestDeparturesProcessing(t *testing.T) {
	var store DepartureStore

	departure := generateDeparture()

	store.InitStore()
	store.ProcessDeparture(departure)

	// Older:
	departure2 := departure

	// Earlier than previous message, so should be ignored:
	departure2.ProductID = "12344"
	departure2.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 68, time.UTC)

	store.ProcessDeparture(departure2)
	departureInStore := store.GetDeparture("1234", "2019-01-27", "UT")

	if departureInStore.ProductID != "12345" {
		t.Error("Should not update departure with earlier departure")
	}
	if store.Counters.Outdated != 1 {
		t.Error("Should increase counter for outdated messages")
	}

	departure3 := departure
	departure3.ProductID = "12343"
	departure3.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 98, time.UTC)

	store.ProcessDeparture(departure3)
	departureInStore = store.GetDeparture("1234", "2019-01-27", "UT")

	if departureInStore.ProductID != "12343" {
		t.Error("Should update departure with later message")
	}
}

func TestRetrieveDeparturesFromStation(t *testing.T) {
	var store DepartureStore

	store.InitStore()

	// Fake some departures
	departure1 := generateDeparture()

	departure2 := generateDeparture()
	departure2.ServiceID = "54321"
	departure2.GenerateID()

	departure3 := generateDeparture()
	departure3.ServiceID = "23456"
	departure3.Status = 5 // Departed
	departure3.GenerateID()

	departure4 := generateDeparture()
	departure4.ServiceID = "23456"
	departure4.Station.Code = "UTO"
	departure4.Station.NameLong = "Utrecht Overvecht"
	departure4.GenerateID()

	store.ProcessDeparture(departure1)
	store.ProcessDeparture(departure2)
	store.ProcessDeparture(departure3)
	store.ProcessDeparture(departure4)

	departuresUT := store.GetStationDepartures("UT", false)

	if len(departuresUT) != 2 {
		t.Error("Wrong number of departures for UT")
	}

	departuresUT = store.GetStationDepartures("UT", true)

	if len(departuresUT) != 3 {
		t.Error("Wrong number of (hidden+non-hidden) departures for UT")
	}

	departuresUTO := store.GetStationDepartures("UTO", false)

	if len(departuresUTO) != 1 {
		t.Error("Wrong number of departures for UTO")
	}

	// Station without departures:
	departuresHRY := store.GetStationDepartures("HRY", false)

	if len(departuresHRY) != 0 {
		t.Error("Wrong number of departures for HRY")
	}
}

func generateDeparture() models.Departure {
	var departure models.Departure

	departure.ProductID = "12345"
	departure.ServiceID = "1234"
	departure.ServiceNumber = "1234"
	departure.Station.Code = "UT"
	departure.Station.NameLong = "Utrecht Centraal"
	departure.ServiceDate = "2019-01-27"
	departure.GenerateID()
	departure.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)
	departure.DepartureTime = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)

	return departure
}

func TestCleanupDepartures(t *testing.T) {
	var store DepartureStore

	store.InitStore()

	// Fake some departures
	departure1 := generateDeparture()
	departure1.Hidden = true

	departure2 := generateDeparture()
	departure2.ServiceID = "54321"
	departure2.GenerateID()

	departure3 := generateDeparture()
	departure3.ServiceID = "99999"
	departure3.DepartureTime = time.Date(2099, time.January, 27, 12, 34, 56, 78, time.UTC)
	departure3.GenerateID()

	departure4 := generateDeparture()
	departure4.ServiceID = "66666"
	departure4.GenerateID()
	departure4.NotRealTime = true
	departure4.DepartureTime = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)

	store.ProcessDeparture(departure1)
	store.ProcessDeparture(departure2)
	store.ProcessDeparture(departure3)
	store.ProcessDeparture(departure4)

	// Verify that we have 4 departures in store:
	if store.GetNumberOfDepartures() != 4 {
		t.Error("Wrong number of departures")
	}

	// Cleanup, first pass:
	store.CleanUp(time.Date(2019, time.January, 27, 12, 30, 56, 78, time.UTC))

	// Nothing should be removed just yet (including the hidden departure, whose departure time
	// is in the future now, and should stay for at least 4 hours)

	// Verify that we have 4 departures in store:
	if store.GetNumberOfDepartures() != 4 {
		t.Error("Wrong number of departures")
	}

	// Cleanup at 12:35. Non-realtime departure should be hidden. The realtime departure should still be visible.
	store.CleanUp(time.Date(2019, time.January, 27, 12, 36, 56, 78, time.UTC))

	if store.GetDeparture(departure4.ServiceID, departure4.ServiceDate, departure4.Station.Code).Hidden == false {
		t.Error("Non-realtime train should be hidden after CleanUp #1")
	}
	if store.GetDeparture(departure2.ServiceID, departure2.ServiceDate, departure2.Station.Code).Hidden == true {
		t.Error("Realtime train should not be hidden after CleanUp #1")
	}
	if store.GetNumberOfDepartures() != 4 {
		t.Error("Wrong number of departures")
	}

	// Now cleanup 10 mins after planned departure, departure2 should be hidden too now
	store.CleanUp(time.Date(2019, time.January, 27, 12, 45, 56, 78, time.UTC))

	if store.GetDeparture(departure4.ServiceID, departure4.ServiceDate, departure4.Station.Code).Hidden == false {
		t.Error("Non-realtime train should be hidden after CleanUp #2")
	}
	if store.GetDeparture(departure2.ServiceID, departure2.ServiceDate, departure2.Station.Code).Hidden == false {
		t.Error("Realtime train should be hidden after CleanUp #2")
	}
	if store.GetNumberOfDepartures() != 4 {
		t.Error("Wrong number of departures")
	}

	// Cleanup #3, everything must be removed
	store.CleanUp(time.Date(2019, time.January, 27, 16, 45, 56, 78, time.UTC))

	// The hidden departure should be gone by now. The second departure should be hidden by now.
	// The third departure should still be visible.
	if store.GetNumberOfDepartures() != 1 {
		t.Fatal("Hidden departures not removed")
	}

	// Verify departure3 is still visible.
	if store.GetDeparture(departure3.ServiceID, departure3.ServiceDate, departure3.Station.Code).Hidden == true {
		t.Error("Train which departs in 2099 should not be hidden already")
	}
}

func TestSaveDepartureStore(t *testing.T) {
	var store, store2 DepartureStore

	store.InitStore()

	for i := 0; i < 40000; i++ {
		departure := generateDeparture()
		departure.ServiceID = strconv.Itoa(i)
		departure.GenerateID()

		store.ProcessDeparture(departure)
	}

	if store.GetNumberOfDepartures() != 40000 {
		t.Error("Wrong number of departures")
	}

	// Save
	error := store.SaveStore()

	if error != nil {
		t.Fatalf("%s", error)
	}

	// Load in empty store:
	store2.InitStore()
	store2.ReadStore()

	if store2.GetNumberOfDepartures() != 40000 {
		t.Error("Wrong number of departures")
	}
}
