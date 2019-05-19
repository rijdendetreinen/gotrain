package stores

import (
	"strconv"
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestServicesCount(t *testing.T) {
	var store ServiceStore
	store.InitStore()

	if store.GetNumberOfServices() != 0 {
		t.Error("Wrong number of services")
	}

	store.ProcessService(generateService())

	if store.GetNumberOfServices() != 1 {
		t.Error("Wrong number of services")
	}

	if store.GetNumberOfServices() != len(store.GetAllServices()) {
		t.Error("Reported number of services does not match with actual inventory count")
	}
}

func TestRetrieveService(t *testing.T) {
	var store ServiceStore

	service := generateService()

	store.InitStore()
	store.ProcessService(service)

	serviceInStore := store.GetService("1234", "2019-01-27")

	if serviceInStore == nil {
		t.Error("Could not retrieve service from store")
	}
}

func TestDuplicateService(t *testing.T) {
	var store ServiceStore

	service := generateService()

	store.InitStore()
	store.ProcessService(service)

	if store.Counters.Duplicates != 0 {
		t.Fatal("Wrong number of services for counters")
	}

	if store.GetNumberOfServices() != 1 {
		t.Error("Wrong number of services")
	}

	// Process again (forcing duplicate)
	store.ProcessService(service)

	if store.GetNumberOfServices() != 1 {
		t.Error("Wrong number of services")
	}

	if store.Counters.Duplicates != 1 {
		t.Error("Should increment counter for duplicates")
	}
}

func TestServiceProcessing(t *testing.T) {
	var store ServiceStore

	service := generateService()

	store.InitStore()
	store.ProcessService(service)

	// Older:
	service2 := service

	// Earlier than previous message, so should be ignored:
	service2.ProductID = "12344"
	service2.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 68, time.UTC)

	store.ProcessService(service2)
	serviceInStore := store.GetService("1234", "2019-01-27")

	if serviceInStore.ProductID != "12345" {
		t.Error("Should not update service with earlier service")
	}
	if store.Counters.Outdated != 1 {
		t.Error("Should increase counter for outdated messages")
	}

	service3 := service
	service3.ProductID = "12343"
	service3.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 98, time.UTC)

	store.ProcessService(service3)
	serviceInStore = store.GetService("1234", "2019-01-27")

	if serviceInStore.ProductID != "12343" {
		t.Error("Should update service with later message")
	}
}

func generateService() models.Service {
	var service models.Service

	service.ProductID = "12345"
	service.ServiceNumber = "1234"
	service.ServiceDate = "2019-01-27"
	service.GenerateID()
	service.Timestamp = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)
	service.ValidUntil = time.Date(2019, time.January, 27, 13, 34, 56, 78, time.UTC)

	var servicePart models.ServicePart
	var stop1, stop2 models.ServiceStop

	stop1.Station.Code = "UT"
	stop1.Station.NameLong = "Utrecht Centraal"
	stop1.DepartureTime = time.Date(2019, time.January, 27, 12, 34, 56, 78, time.UTC)

	stop2.Station.Code = "GVC"
	stop2.Station.NameLong = "Den Haag Centraal"
	stop2.ArrivalTime = time.Date(2019, time.January, 27, 13, 34, 56, 78, time.UTC)

	servicePart.Stops = append(servicePart.Stops, stop1)
	servicePart.Stops = append(servicePart.Stops, stop2)
	service.ServiceParts = append(service.ServiceParts, servicePart)

	return service
}

func TestCleanupServices(t *testing.T) {
	var store ServiceStore

	store.InitStore()

	// Fake some services
	service1 := generateService()
	service1.Hidden = true

	service2 := generateService()
	service2.ServiceNumber = "54321"
	service2.GenerateID()

	service3 := generateService()
	service3.ServiceNumber = "99999"
	service3.ValidUntil = time.Date(2099, time.January, 27, 12, 34, 56, 78, time.UTC)
	service3.GenerateID()

	store.ProcessService(service1)
	store.ProcessService(service2)
	store.ProcessService(service3)

	// Verify that we have 3 services in store:
	if store.GetNumberOfServices() != 3 {
		t.Error("Wrong number of services")
	}

	// Cleanup, first pass:
	// (We expect that the testing system is already beyond January 27th 2019...)
	store.CleanUp(time.Date(2019, time.February, 27, 12, 44, 56, 78, time.UTC))

	// The hidden service should be gone by now. The second service should be hidden by now.
	// The third service should still be visible.
	if store.GetNumberOfServices() > 2 {
		t.Fatal("Hidden service not removed")
	} else if store.GetNumberOfServices() < 2 {
		t.Fatal("Non-hidden service already removed")
	}

	// Verify service2 is hidden by now:
	if store.GetService(service2.ServiceNumber, service2.ServiceDate).Hidden == false {
		t.Error("Departed train should be hidden after CleanUp")
	}

	// Verify service3 is still visible.
	// That is, if you're not testing this code in year 2099 or later (hello from the past!)
	if store.GetService(service3.ServiceNumber, service3.ServiceDate).Hidden == true {
		t.Error("Train which is valid until 2099 should not be hidden already")
	}

	// Second pass for cleaning up.
	// After that, service2 should be gone, service3 still be visible.
	store.CleanUp(time.Date(2019, time.March, 27, 12, 44, 56, 78, time.UTC))

	if store.GetService(service2.ServiceNumber, service2.ServiceDate) != nil {
		t.Error("Service2 should have been deleted by now")
	}

	if store.GetService(service3.ServiceNumber, service3.ServiceDate).Hidden == true {
		t.Error("Train which departs in 2099 should not be hidden already")
	}
}

func TestSaveServiceStore(t *testing.T) {
	var store, store2 ServiceStore

	store.InitStore()

	for i := 0; i < 40000; i++ {
		service := generateService()
		service.ServiceNumber = strconv.Itoa(i)
		service.GenerateID()

		store.ProcessService(service)
	}

	if store.GetNumberOfServices() != 40000 {
		t.Errorf("Wrong number of services: expected %d, got %d", 40000, store.GetNumberOfServices())
	}

	// Save
	error := store.SaveStore()

	if error != nil {
		t.Fatalf("%s", error)
	}

	// Load in empty store:
	store2.InitStore()
	store2.ReadStore()

	if store2.GetNumberOfServices() != 40000 {
		t.Errorf("Wrong number of services: expected %d, got %d", 40000, store2.GetNumberOfServices())
	}
}

func TestSaveServiceStore2(t *testing.T) {
	var store, store2 ServiceStore

	store.InitStore()

	for i := 0; i < 40000; i++ {
		service := generateService()
		service.ServiceNumber = strconv.Itoa(i)
		service.GenerateID()

		store.ProcessService(service)
	}

	for i := 0; i < 40000; i = i + 2 {
		serviceID := "2019-01-27-" + strconv.Itoa(i)

		store.deleteService(serviceID)
	}

	if store.GetNumberOfServices() != 20000 {
		t.Errorf("Wrong number of services: expected %d, got %d", 20000, store.GetNumberOfServices())
	}

	// Save
	error := store.SaveStore()

	if error != nil {
		t.Fatalf("%s", error)
	}

	// Load in empty store:
	store2.InitStore()
	store2.ReadStore()

	if store2.GetNumberOfServices() != 20000 {
		t.Errorf("Wrong number of services: expected %d, got %d", 20000, store2.GetNumberOfServices())
	}
}
