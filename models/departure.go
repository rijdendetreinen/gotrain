package models

import "time"

// Departure is a train service which departs from a single station
type Departure struct {
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

	DepartureTime time.Time
	Delay         int

	ReservationRequired bool
	WithSupplement      bool
	SpecialTicket       bool
	RearPartRemains     bool
	DoNotBoard          bool
	Cancelled           bool
	NotRealTime         bool

	DestinationActual  []Station
	DestinationPlanned []Station
	ViaActual          []Station
	ViaPlanned         []Station

	PlatformActual  string
	PlatformPlanned string

	// TODO: ViaActual, ViaPlanned, Wings, BoardingTips, TravelTips, ChangeTips

	Modifications []Modification

	Hidden bool
}

// RealDepartureTime returns the actual departure time, including delay
func (departure Departure) RealDepartureTime() time.Time {
	var delayDuration time.Duration
	delayDuration = time.Second * time.Duration(departure.Delay)
	return departure.DepartureTime.Add(delayDuration)
}
