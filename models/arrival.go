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

	OriginActual  []Station
	OriginPlanned []Station
	ViaActual     []Station
	ViaPlanned    []Station

	PlatformActual  string
	PlatformPlanned string

	Modifications []Modification

	Hidden bool
}

// RealArrivalTime returns the actual arrival time, including delay
func (arrival Arrival) RealArrivalTime() time.Time {
	var delayDuration time.Duration
	delayDuration = time.Second * time.Duration(arrival.Delay)
	return arrival.ArrivalTime.Add(delayDuration)
}
