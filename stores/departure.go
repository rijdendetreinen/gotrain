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
	stations   map[string]map[string]struct{}
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

	// Hide departured trains:
	if newDeparture.Status == 5 {
		newDeparture.Hidden = true
	}

	store.Lock()
	store.departures[newDeparture.ID] = newDeparture
	store.updateStationReference(newDeparture.Station.Code, newDeparture.ID)
	store.Unlock()

	store.Counters.Processed++
}

func (store *DepartureStore) updateStationReference(station, ID string) {
	_, stationExists := store.stations[station]
	if !stationExists {
		store.stations[station] = make(map[string]struct{})
	}

	_, exists := store.stations[station][ID]
	if !exists {
		store.stations[station][ID] = struct{}{}
	}

}

// InitStore initializes the departure store by creating the departures map
// and sets the downtime detection config
func (store *DepartureStore) InitStore() {
	store.departures = make(map[string]models.Departure)
	store.stations = make(map[string]map[string]struct{})

	store.DowntimeDetection.MinAverage = float64(1) / 60       // One message per minute
	store.DowntimeDetection.MinAverageNight = float64(1) / 600 // One message per 10 minutes
	store.DowntimeDetection.NightStartHour = 2                 // Night starts at 02:00
	store.DowntimeDetection.NightEndHour = 5                   // Night ends at 05:00
	store.DowntimeDetection.RecoveryTime = 70                  // 70 mins recovery time
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

// GetStationDepartures returns all departures for a given station
func (store *DepartureStore) GetStationDepartures(station string, includeHidden bool) []models.Departure {
	var departures []models.Departure

	store.RLock()
	stationServices := store.stations[station]

	for ID := range stationServices {
		departure, found := store.departures[ID]

		if found {
			if includeHidden || !departure.Hidden {
				departures = append(departures, departure)
			}
		}
	}
	store.RUnlock()

	return departures
}

// GetDeparture retrieves a single departure
func (store *DepartureStore) GetDeparture(serviceID, serviceDate string, station string) *models.Departure {
	id := serviceDate + "-" + serviceID + "-" + station

	store.RLock()
	departure, found := store.departures[id]
	store.RUnlock()

	if found {
		return &departure
	}

	return nil
}

// ReadStore reads the save store contents
func (store *DepartureStore) ReadStore() error {
	err := readGob("departures.gob", &store.departures)

	if err != nil {
		return err
	}

	for _, departure := range store.departures {
		store.updateStationReference(departure.Station.Code, departure.ID)
	}

	return nil
}

// SaveStore saves the departures store contents
func (store *DepartureStore) SaveStore() error {
	store.RLock()

	err := writeGob("departures.gob", store.departures)

	store.RUnlock()
	return err
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
func (store *DepartureStore) deleteDeparture(departure models.Departure) {
	store.Lock()
	delete(store.departures, departure.ID)

	_, stationExists := store.stations[departure.Station.Code]

	if stationExists {
		_, exists := store.stations[departure.Station.Code][departure.ID]
		if exists {
			delete(store.stations[departure.Station.Code], departure.ID)
		}
	}

	store.Unlock()
}

// CleanUp removes outdated items
func (store *DepartureStore) CleanUp(currentTime time.Time) {
	// Remove departures which should have departured 4 hours ago:
	thresholdRemove := currentTime.Add(-4 * time.Hour)

	// Hide departures which should have departed 10 minutes ago:
	thresholdHide := currentTime.Add(-10 * time.Minute)

	// Hide departures which should have departed 1 minute ago if they are not realtime:
	thresholdHideNonRealtime := currentTime.Add(-1 * time.Minute)

	log.Debug("Cleaning up departure store")

	store.RLock()
	defer store.RUnlock()

	for departureID, departure := range store.departures {
		store.RUnlock()

		if !departure.Hidden && departure.RealDepartureTime().Before(thresholdHide) {
			log.WithField("DepartureID", departureID).Debug("Hiding departure")

			store.hideDeparture(departureID)
		} else if !departure.Hidden && (departure.NotRealTime || departure.Cancelled) && departure.RealDepartureTime().Before(thresholdHideNonRealtime) {
			log.WithField("DepartureID", departureID).Debug("Hiding non-realtime departure")

			store.hideDeparture(departureID)
		} else if departure.Hidden && departure.RealDepartureTime().Before(thresholdRemove) {
			log.WithField("DepartureID", departureID).Debug("Removing departure")

			store.deleteDeparture(departure)
		}

		store.RLock()
	}
}
