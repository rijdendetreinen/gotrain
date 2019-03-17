package models

import "time"

// Arrival is an arriving train on a single station
type Arrival struct {
	StoreItem

	ServiceID   string
	ServiceDate string
	ServiceName string
	Station     Station

	Status          int
	ServiceNumber   string
	ServiceType     string
	ServiceTypeCode string
	Company         string

	ArrivalTime time.Time
	Delay       int

	ReservationRequired bool
	WithSupplement      bool
	SpecialTicket       bool
	RearPartRemains     bool
	DoNotBoard          bool
	Cancelled           bool
	NotRealTime         bool

	OriginActual    Station
	OriginPlanned   Station
	PlatformActual  string
	PlatformPlanned string

	// TODO: ViaActual, ViaPlanned, etc.

	Modifications []Modification

	Hidden bool
}
