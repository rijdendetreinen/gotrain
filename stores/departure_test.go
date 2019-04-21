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

func generateDeparture() models.Departure {
	var departure models.Departure

	departure.ProductID = "12345"
	departure.ServiceID = "1234"
	departure.ServiceNumber = "1234"
	departure.Station.Code = "UT"
	departure.Station.NameLong = "Utrecht Centraal"
	departure.ServiceDate = "2019-01-27"
	departure.ID = "2019-01-27-1234-UT"
	departure.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)

	return departure
}
