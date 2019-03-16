package stores

import (
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

// The ServiceStore contains all services
type ServiceStore struct {
	Store
	services map[string]models.Service
}

// ProcessService adds or updates a service in a service store
func ProcessService(store *ServiceStore, newService models.Service) {
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

// InitServiceStore initializes the service store by creating the services map
func InitServiceStore(store *ServiceStore) {
	store.services = make(map[string]models.Service)
}

func GetAllServices(store *ServiceStore) map[string]models.Service {
	return store.services
}
