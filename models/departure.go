package models

import "time"

// Departure is the train service containing all parts, stops etc.
type Departure struct {
	ID          string
	Timestamp   time.Time
	ProductID   string
	ServiceID   string
	ServiceDate string
	ServiceName string
	Station     Station

	Status          int
	ServiceNumber   string
	ServiceType     string
	ServiceTypeCode string
	Company         string

	DepartureTime time.Time
	Delay         int

	ReservationRequired bool
	WithSupplement      bool
	SpecialTicket       bool
	RearPartRemains     bool
	DoNotBoard          bool
	Cancelled           bool
	NotRealTime         bool

	DestinationActual  Station
	DestinationPlanned Station
	PlatformActual     string
	PlatformPlanned    string

	// TODO: ViaActual, ViaPlanned, Wings, BoardingTips, TravelTips, ChangeTips

	Modifications []Modification

	Hidden bool
}
