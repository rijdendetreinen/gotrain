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

	TrainWings []TrainWing

	BoardingTips []BoardingTip
	TravelTips   []TravelTip
	ChangeTips   []ChangeTip

	Modifications []Modification

	Hidden bool
}

// BoardingTip is a tip for passengers to board another train for certain destinations
type BoardingTip struct {
	ExitStation       Station
	Destination       Station
	TrainType         string
	TrainTypeCode     string
	DeparturePlatform string
	DepartureTime     time.Time
}

// TravelTip is a tip that a service calls (or doesn't call) at a specific station
type TravelTip struct {
	TipCode  string
	Stations []Station
}

// ChangeTip is a tip to change trains at ChangeStation for the given destination
type ChangeTip struct {
	Destination   Station
	ChangeStation Station
}

// TrainWing is a part of a train departure with a single destination
type TrainWing struct {
	DestinationActual  []Station
	DestinationPlanned []Station
	Stations           []Station
	Material           []Material
	Modifications      []Modification
}

// RealDepartureTime returns the actual departure time, including delay
func (departure Departure) RealDepartureTime() time.Time {
	var delayDuration time.Duration
	delayDuration = time.Second * time.Duration(departure.Delay)
	return departure.DepartureTime.Add(delayDuration)
}
