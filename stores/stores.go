package stores

import (
	"bufio"
	"encoding/gob"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// StatusUnknown for when no status has been determined (not enough information)
const StatusUnknown = "UNKNOWN"

// StatusDown for when system is down
const StatusDown = "DOWN"

// StatusRecovering for when system is recovering
const StatusRecovering = "RECOVERING"

// StatusUp for when system is up and has all data
const StatusUp = "UP"

// Stores is the stores collection. Initialize with InitializeStores()
var Stores StoreCollection

// StoresDataDirectory is the location where stores are saved
var StoresDataDirectory = "data/"

// StoreCollection is the collection of all stores
type StoreCollection struct {
	ArrivalStore   ArrivalStore
	DepartureStore DepartureStore
	ServiceStore   ServiceStore
}

// Store is the generic store struct
type Store struct {
	sync.RWMutex
	Counters          Counters
	Status            string
	measurements      []Measurement
	MessagesAverage   float64
	LastStatusChange  time.Time
	DowntimeDetection DowntimeDetectionConfig
}

// Counters stores some interesting counters for a store
type Counters struct {
	Received   int `json:"received"`
	Processed  int `json:"processed"`
	Error      int `json:"error"`
	Duplicates int `json:"duplicate"`
	Outdated   int `json:"outdated"`
	TooLate    int `json:"too_late"`
}

// DowntimeDetectionConfig contains the configuration for this store's downtime detection
type DowntimeDetectionConfig struct {
	MinAverage      float64 // Minimum average messages per second
	MinAverageNight float64 // Minimum average messages per second during night time
	NightStartHour  int     // Start hour of night
	NightEndHour    int     // End hour of night
	RecoveryTime    int     // Recovery time in minutes
}

// Measurement is a struct to store the number of received and processed messages
type Measurement struct {
	Time      time.Time
	Processed int
}

// CurrentMinimumAverage returns the minimum average, based on the current time
func (downtimeDetection DowntimeDetectionConfig) CurrentMinimumAverage(time time.Time) float64 {
	if time.Hour() >= downtimeDetection.NightStartHour && time.Hour() < downtimeDetection.NightEndHour {
		return downtimeDetection.MinAverageNight
	}

	return downtimeDetection.MinAverage
}

// ResetCounters resets all store counters
func (store *Store) ResetCounters() {
	store.Counters.Received = 0
	store.Counters.Processed = 0
	store.Counters.Error = 0
	store.Counters.Duplicates = 0
	store.Counters.Outdated = 0
	store.Counters.TooLate = 0

	store.measurements = make([]Measurement, 0)
}

// TakeMeasurement takes a new measurement. This method is expected to be called approximately every 20s.
// This function re-calculates the average messages per minute if enough data is available and updates the
// store status accordingly.
func (store *Store) TakeMeasurement() {
	store.newMeasurement(time.Now())
	store.updateStatus(time.Now())
}

// newMeasurement stores a new measurement, and re-calculates the average messages per minute if enough data is available.
// The store status is updated based on the average messages that are processed
func (store *Store) newMeasurement(time time.Time) {
	var measurement Measurement

	measurement.Time = time
	measurement.Processed = store.Counters.Processed

	store.measurements = append(store.measurements, measurement)

	store.MessagesAverage = -1

	if len(store.measurements) > 1 {
		foundMeasurement := false
		var firstMeasurement Measurement
		popMeasurements := 0

		// Loop over measurements, until they do not meet our condition anymore:
		for index, earlierMeasurement := range store.measurements {
			duration := measurement.Time.Sub(earlierMeasurement.Time)

			if duration.Seconds() >= 600 {
				foundMeasurement = true
				firstMeasurement = earlierMeasurement
				popMeasurements = index
			} else {
				break
			}
		}

		if foundMeasurement {
			duration := measurement.Time.Sub(firstMeasurement.Time)
			store.MessagesAverage = float64(measurement.Processed-firstMeasurement.Processed) / duration.Seconds()

			if popMeasurements > 0 {
				store.measurements = store.measurements[popMeasurements:]
			}
		}
	}
}

// Update the store status based on the current messagesAverage
func (store *Store) updateStatus(currentTime time.Time) {
	// Determine whether we are currently receiving messages:
	isReceiving := store.MessagesAverage >= store.DowntimeDetection.CurrentMinimumAverage(currentTime)

	// Determine possible status changes:
	if isReceiving && (store.Status == StatusUnknown || store.Status == StatusDown) {
		// Status was DOWN or UNKNOWN, but we are currently receiving.
		// Change to RECOVERING
		store.Status = StatusRecovering
		store.LastStatusChange = currentTime
	} else if isReceiving && store.Status == StatusRecovering {
		// We are currently receiving and our status is RECOVERING
		// Check last update time to see if we can change to UP:
		if currentTime.Sub(store.LastStatusChange).Seconds() >= float64(store.DowntimeDetection.RecoveryTime*60) {
			store.Status = StatusUp
			store.LastStatusChange = currentTime
		}
	} else if isReceiving && store.Status == StatusUp {
		// We are currently receiving and our status was already UP.
		// Keep up the good job!
	} else if !isReceiving {
		// We are not receiving.
		if store.MessagesAverage == -1 {
			// Average of -1 implies not enough data to determine avg. number of messages
			// Change status to UNKNOWN (if it wasn't already UNKNOWN):
			if store.Status != StatusUnknown {
				store.Status = StatusUnknown
				store.LastStatusChange = currentTime
			}
		} else if store.Status != StatusDown {
			// Average is not -1, so we have valid data about our received messages.
			// Change status to DOWN (if it wasn't already DOWN):
			store.Status = StatusDown
			store.LastStatusChange = currentTime
		}
	}
}

// ResetStatus resets the status and counters of a store
func (store *Store) ResetStatus() {
	store.ResetCounters()

	store.Status = StatusUnknown
	store.MessagesAverage = 0
	store.LastStatusChange = time.Now()
}

// InitializeStores initializes all stores and resets their counters/status
func InitializeStores() *StoreCollection {
	Stores.ArrivalStore.ResetStatus()
	Stores.ArrivalStore.InitStore()

	Stores.DepartureStore.ResetStatus()
	Stores.DepartureStore.InitStore()

	Stores.ServiceStore.ResetStatus()
	Stores.ServiceStore.InitStore()

	return &Stores
}

// CleanUp cleans all stores and removes outdated items
func CleanUp() {
	currentTime := time.Now()

	Stores.ArrivalStore.CleanUp(currentTime)
	Stores.DepartureStore.CleanUp(currentTime)
	Stores.ServiceStore.CleanUp(currentTime)
}

// TakeMeasurements takes measurements for all stores and updates downtime status
func TakeMeasurements() {
	Stores.ArrivalStore.TakeMeasurement()
	Stores.DepartureStore.TakeMeasurement()
	Stores.ServiceStore.TakeMeasurement()
}

// LoadStores reads all store content files
func LoadStores() error {
	servicesError := Stores.ServiceStore.ReadStore()
	departuresError := Stores.DepartureStore.ReadStore()
	arrivalsError := Stores.ArrivalStore.ReadStore()

	if servicesError != nil {
		return servicesError
	} else if departuresError != nil {
		return departuresError
	} else if arrivalsError != nil {
		return arrivalsError
	}

	return nil
}

// SaveStores saves all stores
func SaveStores() {
	servicesError := Stores.ServiceStore.SaveStore()
	departuresError := Stores.DepartureStore.SaveStore()
	arrivalsError := Stores.ArrivalStore.SaveStore()

	if servicesError != nil {
		log.WithError(servicesError).Error("Error while saving services store")
	}
	if departuresError != nil {
		log.WithError(departuresError).Error("Error while saving departures store")
	}
	if arrivalsError != nil {
		log.WithError(arrivalsError).Error("Error while saving arrivals store")
	}
}

// Encode a GOB file
func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(getDataDirectory() + filePath)

	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	enc := gob.NewEncoder(w)
	err = enc.Encode(object)

	w.Flush()
	file.Close()

	return err
}

// Read a GOB file
func readGob(filePath string, object interface{}) error {
	file, err := os.Open(getDataDirectory() + filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func getDataDirectory() string {
	return StoresDataDirectory
}
