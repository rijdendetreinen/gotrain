package stores

import (
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

func TestCleanup(t *testing.T) {
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

	store.ProcessDeparture(departure1)
	store.ProcessDeparture(departure2)
	store.ProcessDeparture(departure3)

	// Verify that we have 3 departures in store:
	if store.GetNumberOfDepartures() != 3 {
		t.Error("Wrong number of departures")
	}

	// Cleanup, first pass:
	// (We expect that the testing system is already beyond January 27th 2019...)
	store.CleanUp()

	// Teh hidden departure should be gone by now. The second departure should be hidden by now.
	// The third departure should still be visible.
	if store.GetNumberOfDepartures() > 2 {
		t.Fatal("Hidden departure not removed")
	} else if store.GetNumberOfDepartures() < 2 {
		t.Fatal("Non-hidden departure already removed")
	}

	// Verify departure2 is hidden by now:
	if store.GetDeparture(departure2.ServiceID, departure2.ServiceDate, departure2.Station.Code).Hidden == false {
		t.Error("Departed train should be hidden after CleanUp")
	}

	// Verify departure3 is still visible.
	// That is, if you're not testing this code in year 2099 or later (hello from the past!)
	if store.GetDeparture(departure3.ServiceID, departure3.ServiceDate, departure3.Station.Code).Hidden == true {
		t.Error("Train which departs in 2099 should not be hidden already")
	}

	// Second pass for cleaning up.
	// After that, departure2 should be gone, departure3 still be visible.
	store.CleanUp()

	if store.GetDeparture(departure2.ServiceID, departure2.ServiceDate, departure2.Station.Code) != nil {
		t.Error("Departure2 should have been deleted by now")
	}

	if store.GetDeparture(departure3.ServiceID, departure3.ServiceDate, departure3.Station.Code).Hidden == true {
		t.Error("Train which departs in 2099 should not be hidden already")
	}

}
