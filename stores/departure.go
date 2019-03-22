package stores

import (
	"time"

	"github.com/rijdendetreinen/gotrain/models"
	log "github.com/sirupsen/logrus"
)

// The DepartureStore contains all departures
type DepartureStore struct {
	Store
	departures map[string]models.Departure
}

// ProcessDeparture adds or updates a departure in a departure store
func (store *DepartureStore) ProcessDeparture(newDeparture models.Departure) {
	store.Counters.Received++

	// Check whether departure already exists. If so, check whether this message is newer.
	if existingDeparture, ok := store.departures[newDeparture.ID]; ok {
		// Check for duplicate:
		if existingDeparture.ProductID == newDeparture.ProductID {
			store.Counters.Duplicates++
			store.Counters.Processed++
			return
		}

		// Check whether newDeparture is actually newer:
		if existingDeparture.Timestamp.After(newDeparture.Timestamp) {
			store.Counters.Outdated++
			store.Counters.Processed++
			return
		}
	}

	// Message is not duplicate or outdated, continue processing

	// Check message age (just for warning, always process):
	threshold := time.Now()
	threshold = threshold.Add(-10 * time.Second)

	if newDeparture.Timestamp.Before(threshold) {
		store.Counters.TooLate++
	}

	store.departures[newDeparture.ID] = newDeparture

	store.Counters.Processed++
}

// InitStore initializes the departure store by creating the departures map
func (store *DepartureStore) InitStore() {
	store.departures = make(map[string]models.Departure)
}

// GetNumberOfDepartures returns the number of departures in the store (unfiltered)
func (store *DepartureStore) GetNumberOfDepartures() int {
	return len(store.departures)
}

// GetAllDepartures simply returns all departures in the store
func (store DepartureStore) GetAllDepartures() map[string]models.Departure {
	return store.departures
}

// ReadStore reads the save store contents
func (store DepartureStore) ReadStore() error {
	return readGob("data/departures.gob", &store.departures)
}

// SaveStore saves the departures store contents
func (store DepartureStore) SaveStore() error {
	return writeGob("data/departures.gob", store.departures)
}

// CleanUp removes outdated items
func (store *DepartureStore) CleanUp() {
	// Remove departures which should have departured 4 hours ago:
	thresholdRemove := time.Now().Add(-4 * time.Hour)

	// Hide departures which should have departed 10 minutes ago:
	thresholdHide := time.Now().Add(-10 * time.Minute)

	log.Debug("Cleaning up departure store")

	for departureID, departure := range store.departures {
		if !departure.Hidden && departure.RealDepartureTime().Before(thresholdHide) {
			log.WithField("DepartureID", departureID).Debug("Hiding departure")

			departure.Hidden = true
			store.departures[departureID] = departure
		} else if departure.Hidden && departure.RealDepartureTime().Before(thresholdRemove) {
			log.WithField("DepartureID", departureID).Debug("Removing departure")

			delete(store.departures, departureID)
		}
	}
}
