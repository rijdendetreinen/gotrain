package stores

import (
	"time"

	"github.com/rijdendetreinen/gotrain/models"
	log "github.com/sirupsen/logrus"
)

// The ServiceStore contains all services
type ServiceStore struct {
	Store
	services map[string]models.Service
}

// ProcessService adds or updates a service in a service store
func (store *ServiceStore) ProcessService(newService models.Service) {
	store.Counters.Received++

	// Check whether service already exists. If so, check whether this message is newer.
	if existingService, ok := store.services[newService.ID]; ok {
		// Check for duplicate:
		if existingService.ProductID == newService.ProductID {
			store.Counters.Duplicates++
			store.Counters.Processed++
			return
		}

		// Check whether newService is actually newer:
		if existingService.Timestamp.After(newService.Timestamp) {
			store.Counters.Outdated++
			store.Counters.Processed++
			return
		}
	}

	// Message is not duplicate or outdated, continue processing

	// Check message age (just for warning, always process):
	threshold := time.Now()
	threshold = threshold.Add(-10 * time.Second)

	if newService.Timestamp.Before(threshold) {
		store.Counters.TooLate++
	}

	store.services[newService.ID] = newService

	store.Counters.Processed++
}

// InitStore initializes the service store by creating the services map
func (store *ServiceStore) InitStore() {
	store.services = make(map[string]models.Service)
}

// GetNumberOfServices returns the number of services in the store (unfiltered)
func (store ServiceStore) GetNumberOfServices() int {
	return len(store.services)
}

// GetAllServices simply returns all services in the store
func (store ServiceStore) GetAllServices() map[string]models.Service {
	return store.services
}

// GetService retrieves a single service
func (store ServiceStore) GetService(serviceID, serviceDate string) *models.Service {
	id := serviceDate + "-" + serviceID

	if val, ok := store.services[id]; ok {
		return &val
	}

	return nil
}

// ReadStore reads the save store contents
func (store ServiceStore) ReadStore() error {
	return readGob("data/services.gob", &store.services)
}

// SaveStore saves the service store contents
func (store ServiceStore) SaveStore() error {
	return writeGob("data/services.gob", store.services)
}

// CleanUp removes outdated items
func (store *ServiceStore) CleanUp() {
	// Remove all services before date X:
	thresholdRemove := time.Now().AddDate(0, 0, -1)
	thresholdHide := time.Now()

	log.WithField("thresholdHide", thresholdHide).WithField("thresholdRemove", thresholdRemove).Debug("Cleaning up service store; hiding and removing all services before thresholds")

	for _, service := range store.services {
		if !service.Hidden && service.ValidUntil.Before(thresholdHide) {
			log.Debug("HIDE")
		} else if service.ValidUntil.Before(thresholdRemove) {
			log.Debug("REMOVE")
			log.Debug(service.ValidUntil)
		}
	}
}
