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
	store.RLock()
	existingDeparture, departureExists := store.departures[newDeparture.ID]
	store.RUnlock()

	if departureExists {
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

	store.Lock()
	store.departures[newDeparture.ID] = newDeparture
	store.Unlock()

	store.Counters.Processed++
}

// InitStore initializes the departure store by creating the departures map
func (store *DepartureStore) InitStore() {
	store.departures = make(map[string]models.Departure)
}

// GetNumberOfDepartures returns the number of departures in the store (unfiltered)
func (store *DepartureStore) GetNumberOfDepartures() int {
	store.RLock()
	count := len(store.departures)
	store.RUnlock()

	return count
}

// GetAllDepartures simply returns all departures in the store
func (store *DepartureStore) GetAllDepartures() map[string]models.Departure {
	store.RLock()
	departures := store.departures
	store.RUnlock()

	return departures
}

// GetDeparture retrieves a single departure
func (store *DepartureStore) GetDeparture(serviceID, serviceDate string, station string) *models.Departure {
	id := serviceDate + "-" + serviceID + "-" + station

	store.RLock()
	if departure, ok := store.departures[id]; ok {
		return &departure
	}
	store.RUnlock()

	return nil
}

// ReadStore reads the save store contents
func (store *DepartureStore) ReadStore() error {
	return readGob("data/departures.gob", &store.departures)
}

// SaveStore saves the departures store contents
func (store *DepartureStore) SaveStore() error {
	return writeGob("data/departures.gob", store.departures)
}

// hideDeparture hides a departure
func (store *DepartureStore) hideDeparture(ID string) {
	store.Lock()
	departure := store.departures[ID]
	departure.Hidden = true
	store.departures[ID] = departure
	store.Unlock()
}

// deleteDeparture deletes a departure
func (store *DepartureStore) deleteDeparture(ID string) {
	store.Lock()
	delete(store.departures, ID)
	store.Unlock()
}

// CleanUp removes outdated items
func (store *DepartureStore) CleanUp() {
	// Remove departures which should have departured 4 hours ago:
	thresholdRemove := time.Now().Add(-4 * time.Hour)

	// Hide departures which should have departed 10 minutes ago:
	thresholdHide := time.Now().Add(-10 * time.Minute)

	log.Debug("Cleaning up departure store")

	store.RLock()
	defer store.RUnlock()

	for departureID, departure := range store.departures {
		store.RUnlock()

		if !departure.Hidden && departure.RealDepartureTime().Before(thresholdHide) {
			log.WithField("DepartureID", departureID).Debug("Hiding departure")

			store.hideDeparture(departureID)
		} else if departure.Hidden && departure.RealDepartureTime().Before(thresholdRemove) {
			log.WithField("DepartureID", departureID).Debug("Removing departure")

			store.deleteDeparture(departureID)
		}

		store.RLock()
	}
}
