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
	store.RLock()
	existingService, serviceExists := store.services[newService.ID]
	store.RUnlock()

	if serviceExists {
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

	store.Lock()
	store.services[newService.ID] = newService
	store.Unlock()

	store.Counters.Processed++
}

// InitStore initializes the service store by creating the services map
func (store *ServiceStore) InitStore() {
	store.services = make(map[string]models.Service)

	store.DowntimeDetection.MinAverage = float64(1) / 60       // One message per minute
	store.DowntimeDetection.MinAverageNight = float64(1) / 600 // One message per 10 minutes
	store.DowntimeDetection.NightStartHour = 2                 // Night starts at 02:00
	store.DowntimeDetection.NightEndHour = 5                   // Night ends at 05:00
	store.DowntimeDetection.RecoveryTime = 1                   // 1 minute recovery time
}

// GetNumberOfServices returns the number of services in the store (unfiltered)
func (store *ServiceStore) GetNumberOfServices() int {
	store.RLock()
	count := len(store.services)
	store.RUnlock()

	return count
}

// GetAllServices simply returns all services in the store
func (store *ServiceStore) GetAllServices() map[string]models.Service {
	store.RLock()
	services := store.services
	store.RUnlock()

	return services
}

// GetService retrieves a single service
func (store *ServiceStore) GetService(serviceID, serviceDate string) *models.Service {
	id := serviceDate + "-" + serviceID

	store.RLock()
	service, found := store.services[id]
	store.RUnlock()

	if found {
		return &service
	}

	return nil
}

// hideService hides a service
func (store *ServiceStore) hideService(serviceID string) {
	store.Lock()
	service := store.services[serviceID]
	service.Hidden = true
	store.services[serviceID] = service
	store.Unlock()
}

// deleteService deletes a service
func (store *ServiceStore) deleteService(serviceID string) {
	store.Lock()
	delete(store.services, serviceID)
	store.Unlock()
}

// ReadStore reads the save store contents
func (store *ServiceStore) ReadStore() error {
	return readGob("data/services.gob", &store.services)
}

// SaveStore saves the service store contents
func (store *ServiceStore) SaveStore() error {
	return writeGob("data/services.gob", store.services)
}

// CleanUp removes outdated items
func (store *ServiceStore) CleanUp() {
	// Remove all services before date X:
	thresholdRemove := time.Now().AddDate(0, 0, -2)
	thresholdHide := time.Now()

	log.Debug("Cleaning up service store")

	store.RLock()
	defer store.RUnlock()

	for serviceID, service := range store.services {
		store.RUnlock()

		if !service.Hidden && service.ValidUntil.Before(thresholdHide) {
			log.WithField("ServiceID", serviceID).Debug("Hiding service")

			store.hideService(serviceID)
		} else if service.Hidden && service.ValidUntil.Before(thresholdRemove) {
			log.WithField("ServiceID", serviceID).Debug("Removing service")

			store.deleteService(serviceID)
		}

		store.RLock()
	}
}
