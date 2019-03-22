package stores

import (
	"encoding/gob"
	"os"
	"time"
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

// StoreCollection is the collection of all stores
type StoreCollection struct {
	ArrivalStore   ArrivalStore
	DepartureStore DepartureStore
	ServiceStore   ServiceStore
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

// Store is the generic store struct
type Store struct {
	Counters         Counters
	status           string
	messagesAverage  float32
	lastStatusChange time.Time
}

// ResetCounters resets all store counters
func (store *Store) ResetCounters() {
	store.Counters.Received = 0
	store.Counters.Processed = 0
	store.Counters.Error = 0
	store.Counters.Duplicates = 0
	store.Counters.Outdated = 0
	store.Counters.TooLate = 0
}

// ResetStatus resets the status and counters of a store
func (store *Store) ResetStatus() {
	store.ResetCounters()

	store.status = StatusUnknown
	store.messagesAverage = 0
	store.lastStatusChange = time.Now()
}

// InitializeStores initializes all stores and resets their counters/status
func InitializeStores() StoreCollection {
	Stores.ArrivalStore.ResetStatus()
	Stores.ArrivalStore.InitStore()

	Stores.DepartureStore.ResetStatus()
	Stores.DepartureStore.InitStore()

	Stores.ServiceStore.ResetStatus()
	Stores.ServiceStore.InitStore()

	return Stores
}

// CleanUp cleans all stores and removes outdated items
func CleanUp() {
	Stores.ServiceStore.CleanUp()
	Stores.DepartureStore.CleanUp()
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
func SaveStores() error {
	servicesError := Stores.ServiceStore.SaveStore()
	departuresError := Stores.DepartureStore.SaveStore()
	arrivalsError := Stores.ArrivalStore.SaveStore()

	if servicesError != nil {
		return servicesError
	} else if departuresError != nil {
		return departuresError
	} else if arrivalsError != nil {
		return arrivalsError
	}

	return nil
}

// Encode a GOB file
func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Read a GOB file
func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
