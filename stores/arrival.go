package stores

import (
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

// The ArrivalStore contains all arrivals
type ArrivalStore struct {
	Store
	arrivals map[string]models.Arrival
}

// ProcessArrival adds or updates an arrival in an arrival store
func (store *ArrivalStore) ProcessArrival(newArrival models.Arrival) {
	store.Counters.Received++

	// Check whether an arrival already exists. If so, check whether this message is newer.
	if existingArrival, ok := store.arrivals[newArrival.ID]; ok {
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

	store.arrivals[newArrival.ID] = newArrival

	store.Counters.Processed++
}

// InitStore initializes the arrival store by creating the arrivals map
func (store *ArrivalStore) InitStore() {
	store.arrivals = make(map[string]models.Arrival)
}

// GetNumberOfArrivals returns the number of arrivals in the store (unfiltered)
func (store *ArrivalStore) GetNumberOfArrivals() int {
	return len(store.arrivals)
}

// GetAllArrivals simply returns all arrivals in the store
func (store ArrivalStore) GetAllArrivals() map[string]models.Arrival {
	return store.arrivals
}

// SaveStore saves the departures store contents
func (store ArrivalStore) SaveStore() error {
	return writeGob("data/arrivals.gob", store.arrivals)
}
