package stores

import (
	"time"

	"github.com/rijdendetreinen/gotrain/models"
	log "github.com/sirupsen/logrus"
)

// The ArrivalStore contains all arrivals
type ArrivalStore struct {
	Store
	arrivals map[string]models.Arrival
	stations map[string]map[string]struct{}
}

// ProcessArrival adds or updates an arrival in an arrival store
func (store *ArrivalStore) ProcessArrival(newArrival models.Arrival) {
	store.Counters.Received++

	// Check whether an arrival already exists. If so, check whether this message is newer.
	store.RLock()
	existingArrival, arrivalExists := store.arrivals[newArrival.ID]
	store.RUnlock()

	if arrivalExists {
		// Check for duplicate:
		if existingArrival.ProductID == newArrival.ProductID {
			store.Counters.Duplicates++
			store.Counters.Processed++
			return
		}

		// Check whether newArrival is actually newer:
		if existingArrival.Timestamp.After(newArrival.Timestamp) {
			store.Counters.Outdated++
			store.Counters.Processed++
			return
		}
	}

	// Message is not duplicate or outdated, continue processing

	// Check message age (just for warning, always process):
	threshold := time.Now()
	threshold = threshold.Add(-10 * time.Second)

	if newArrival.Timestamp.Before(threshold) {
		store.Counters.TooLate++
	}

	store.Lock()
	store.arrivals[newArrival.ID] = newArrival
	store.updateStationReference(newArrival.Station.Code, newArrival.ID)
	store.Unlock()

	store.Counters.Processed++
}

func (store *ArrivalStore) updateStationReference(station, ID string) {
	_, stationExists := store.stations[station]
	if !stationExists {
		store.stations[station] = make(map[string]struct{})
	}

	_, exists := store.stations[station][ID]
	if !exists {
		store.stations[station][ID] = struct{}{}
	}

}

// InitStore initializes the arrival store by creating the arrivals map
func (store *ArrivalStore) InitStore() {
	store.arrivals = make(map[string]models.Arrival)
	store.stations = make(map[string]map[string]struct{})
}

// GetNumberOfArrivals returns the number of arrivals in the store (unfiltered)
func (store *ArrivalStore) GetNumberOfArrivals() int {
	store.RLock()
	count := len(store.arrivals)
	store.RUnlock()

	return count
}

// GetAllArrivals simply returns all arrivals in the store
func (store *ArrivalStore) GetAllArrivals() map[string]models.Arrival {
	store.RLock()
	arrivals := store.arrivals
	store.RUnlock()

	return arrivals
}

// GetStationArrivals returns all arrivals for a given station
func (store *ArrivalStore) GetStationArrivals(station string, includeHidden bool) []models.Arrival {
	var arrivals []models.Arrival

	store.RLock()
	stationServices := store.stations[station]

	for ID := range stationServices {
		arrival, found := store.arrivals[ID]

		if found {
			if includeHidden || !arrival.Hidden {
				arrivals = append(arrivals, arrival)
			}
		}
	}
	store.RUnlock()

	return arrivals
}

// GetArrival retrieves a single arrival
func (store *ArrivalStore) GetArrival(serviceID, serviceDate string, station string) *models.Arrival {
	id := serviceDate + "-" + serviceID + "-" + station

	store.RLock()
	arrival, found := store.arrivals[id]
	store.RUnlock()

	if found {
		return &arrival
	}

	return nil
}

// ReadStore reads the save store contents
func (store *ArrivalStore) ReadStore() error {
	err := readGob("data/arrivals.gob", &store.arrivals)

	if err != nil {
		return err
	}

	for _, arrival := range store.arrivals {
		store.updateStationReference(arrival.Station.Code, arrival.ID)
	}

	return nil
}

// SaveStore saves the arrivals store contents
func (store *ArrivalStore) SaveStore() error {
	return writeGob("data/arrivals.gob", store.arrivals)
}

// hideArrival hides an arrival
func (store *ArrivalStore) hideArrival(ID string) {
	store.Lock()
	arrival := store.arrivals[ID]
	arrival.Hidden = true
	store.arrivals[ID] = arrival
	store.Unlock()
}

// deleteArrival deletes an arrival
func (store *ArrivalStore) deleteArrival(arrival models.Arrival) {
	store.Lock()
	delete(store.arrivals, arrival.ID)

	_, stationExists := store.stations[arrival.Station.Code]

	if stationExists {
		_, exists := store.stations[arrival.Station.Code][arrival.ID]
		if exists {
			delete(store.stations[arrival.Station.Code], arrival.ID)
		}
	}

	store.Unlock()
}

// CleanUp removes outdated items
func (store *ArrivalStore) CleanUp() {
	// Remove arrivals which have arrived 4 hours ago:
	thresholdRemove := time.Now().Add(-4 * time.Hour)

	// Hide arrivals which should have arrived 30 minutes ago:
	thresholdHide := time.Now().Add(-30 * time.Minute)

	log.Debug("Cleaning up arrival store")

	store.RLock()
	defer store.RUnlock()

	for arrivalID, arrival := range store.arrivals {
		store.RUnlock()

		if !arrival.Hidden && arrival.RealArrivalTime().Before(thresholdHide) {
			log.WithField("ArrivalID", arrivalID).Debug("Hiding arrival")

			store.hideArrival(arrivalID)
		} else if arrival.Hidden && arrival.RealArrivalTime().Before(thresholdRemove) {
			log.WithField("ArrivalID", arrivalID).Debug("Removing arrival")

			store.deleteArrival(arrival)
		}

		store.RLock()
	}
}
