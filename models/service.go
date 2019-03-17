package models

import "time"

// Service is the train service containing all parts, stops etc.
type Service struct {
	ID              string
	Timestamp       time.Time
	ProductID       string
	ValidUntil      time.Time
	ServiceID       string
	ServiceNumber   string
	ServiceDate     string
	ServiceType     string
	ServiceTypeCode string
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
	Station          Station
	StationAccesible bool

	AssistenceAvailable bool

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
