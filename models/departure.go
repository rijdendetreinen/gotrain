package models

import (
	"fmt"
	"time"
)

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

// PlatformChanged returns true when the platform has been changed
func (departure Departure) PlatformChanged() bool {
	return departure.PlatformActual != departure.PlatformPlanned
}

// ActualDestinationString returns a string of all actual destinations (long name)
func (departure Departure) ActualDestinationString() string {
	return stationsLongString(departure.DestinationActual, "/")
}

// PlannedDestinationString returns a string of all planned destinations (long name)
func (departure Departure) PlannedDestinationString() string {
	return stationsLongString(departure.DestinationPlanned, "/")
}

// ActualDestinationCodes returns a slice of all actual destination station codes
func (departure Departure) ActualDestinationCodes() []string {
	return stationCodes(departure.DestinationActual)
}

// ViaStationsString returns a string of all actual via stations (medium name)
func (departure Departure) ViaStationsString() string {
	return stationsMediumString(departure.ViaActual, ", ")
}

// Translation provides a translation for this tip
func (tip BoardingTip) Translation(language string) string {
	translation := Translate("%s %s naar %s (spoor %s) is eerder in %s", "%s %s to %s (platform %s) reaches %s sooner", language)

	return fmt.Sprintf(translation, tip.TrainTypeCode, tip.DepartureTime.Local().Format("15:04"), tip.Destination.NameLong, tip.DeparturePlatform, tip.ExitStation.NameLong)
}

// Translation provides a translation for this tip
func (tip ChangeTip) Translation(language string) string {
	translation := Translate("Voor %s overstappen in %s", "For %s, change at %s", language)

	return fmt.Sprintf(translation, tip.Destination, tip.ChangeStation)
}

// Translation provides a translation for this tip
func (tip TravelTip) Translation(language string) string {
	switch tip.TipCode {
	case "STNS":
		return TranslateStations("Stopt niet in %s", "Does not call at %s", tip.Stations, language)
	case "STO":
		return TranslateStations("Stopt ook in %s", "Also calls at %s", tip.Stations, language)
	case "STVA":
		return TranslateStations("Stopt vanaf %s op alle tussengelegen stations", "Calls at all stations after %s", tip.Stations, language)
	case "STNVA":
		return TranslateStations("Stopt vanaf %s niet op tussengelegen stations", "Does not call at intermediate stations after %s", tip.Stations, language)
	case "STT":
		return TranslateStations("Stopt tot %s op alle tussengelegen stations", "Calls at all stations until %s", tip.Stations, language)
	case "STNT":
		return TranslateStations("Stopt tot %s niet op tussengelegen stations", "First stop at %s", tip.Stations, language)
	case "STAL":
		return Translate("Stopt op alle tussengelegen stations", "Calls at all stations", language)
	case "STN":
		return Translate("Stopt niet op tussengelegen stations", "Does not call at intermediate stations", language)
	}

	// Fallback:
	return tip.TipCode
}
