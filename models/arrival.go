package models

import "time"

// Arrival is an arriving train on a single station
type Arrival struct {
	StoreItem

	ServiceID   string
	ServiceDate string
	ServiceName string
	Station     Station
	LineNumber  string

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

// GenerateID generates an ID for this arrival
func (arrival *Arrival) GenerateID() {
	arrival.ID = arrival.ServiceDate + "-" + arrival.ServiceID + "-" + arrival.Station.Code
}

// RealArrivalTime returns the actual arrival time, including delay
func (arrival Arrival) RealArrivalTime() time.Time {
	var delayDuration time.Duration
	delayDuration = time.Second * time.Duration(arrival.Delay)
	return arrival.ArrivalTime.Add(delayDuration)
}

// PlatformChanged returns true when the platform has been changed
func (arrival Arrival) PlatformChanged() bool {
	return arrival.PlatformActual != arrival.PlatformPlanned
}

// ActualOriginString returns a string of all actual origins (long name)
func (arrival Arrival) ActualOriginString() string {
	return stationsLongString(arrival.OriginActual, "/")
}

// PlannedOriginString returns a string of all planned origins (long name)
func (arrival Arrival) PlannedOriginString() string {
	return stationsLongString(arrival.OriginPlanned, "/")
}

// ActualOriginCodes returns a slice of all actual origin station codes
func (arrival Arrival) ActualOriginCodes() []string {
	return stationCodes(arrival.OriginActual)
}

// ViaStationsString returns a string of all actual via stations (medium name)
func (arrival Arrival) ViaStationsString() string {
	return stationsMediumString(arrival.ViaActual, ", ")
}
