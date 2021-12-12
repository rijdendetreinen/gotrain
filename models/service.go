package models

import (
	"time"
)

// Service is the train service containing all parts, stops etc.
type Service struct {
	StoreItem

	ValidUntil      time.Time
	ServiceNumber   string
	ServiceDate     string
	ServiceType     string
	ServiceTypeCode string
	LineNumber      string
	Company         string

	ServiceParts []ServicePart

	ReservationRequired bool
	WithSupplement      bool
	SpecialTicket       bool
	JourneyPlanner      bool

	Modifications []Modification

	Hidden bool
}

// ServicePart is a single part of a train service (a service usually contains just one part, but may contain more)
type ServicePart struct {
	ServiceNumber string
	Stops         []ServiceStop
	Modifications []Modification
}

// ServiceStop is a stops which is called by a train service.
type ServiceStop struct {
	Station Station

	StationAccessible   bool
	AssistanceAvailable bool

	DestinationActual        string
	DestinationPlanned       string
	ArrivalPlatformActual    string
	ArrivalPlatformPlanned   string
	DeparturePlatformActual  string
	DeparturePlatformPlanned string

	StoppingActual  bool
	StoppingPlanned bool
	StopType        string
	DoNotBoard      bool

	ArrivalTime    time.Time
	ArrivalDelay   int
	DepartureTime  time.Time
	DepartureDelay int

	ArrivalCancelled   bool
	DepartureCancelled bool

	Modifications []Modification
	Material      []Material
}

// GetStoppingStations filters out ServiceStops which are not called by the service.
func (servicePart ServicePart) GetStoppingStations() (stops []ServiceStop) {
	for _, stop := range servicePart.Stops {
		if stop.IsStopping() {
			stops = append(stops, stop)
		}
	}

	return
}

// GenerateID generates an ID for this service
func (service *Service) GenerateID() {
	service.ID = service.ServiceDate + "-" + service.ServiceNumber
}

// GetStops returns all stops (from all service parts) as a map, with the station code as key
func (service *Service) GetStops() map[string]ServiceStop {
	stops := make(map[string]ServiceStop)

	for _, part := range service.ServiceParts {
		for _, stop := range part.Stops {
			if stop.IsStopping() {
				stops[stop.Station.Code] = stop
			}
		}
	}

	return stops
}

// IsStopping checks whether the service is stopping at this stop or whether is was planned to do so
func (stop *ServiceStop) IsStopping() bool {
	return stop.StoppingActual || stop.StoppingPlanned
}

// ArrivalPlatformChanged returns true when the planned platform is not the actual one
func (stop *ServiceStop) ArrivalPlatformChanged() bool {
	return stop.ArrivalPlatformPlanned != stop.ArrivalPlatformActual
}

// DeparturePlatformChanged returns true when the planned platform is not the actual one
func (stop *ServiceStop) DeparturePlatformChanged() bool {
	return stop.DeparturePlatformPlanned != stop.DeparturePlatformActual
}
