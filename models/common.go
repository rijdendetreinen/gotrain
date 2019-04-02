package models

import (
	"strings"
	"time"
)

// Station is a station in the railway network. It has a code and 3 names (short, medium, long)
type Station struct {
	Code       string `json:"code"`
	NameShort  string `json:"short"`
	NameMedium string `json:"medium"`
	NameLong   string `json:"long"`
}

// Modification is a change (to the schedule) which is communicated to travellers
type Modification struct {
	ModificationType int     `json:"type"`
	CauseShort       string  `json:"cause_short"`
	CauseLong        string  `json:"cause_long"`
	Station          Station `json:"station"`
}

// Material is the physical train unit
type Material struct {
	NaterialType       string  `json:"type"`
	Number             string  `json:"number"`
	Position           int     `json:"position"`
	DestinationActual  Station `json:"destination_actual"`
	DestinationPlanned Station `json:"destination_planned"`
	Accessible         bool    `json:"accesible"`
	RemainsBehind      bool    `json:"remains_behind"`
}

// StoreItem is for shared fields like ID, timestamp etc.
type StoreItem struct {
	ID        string    `json:"-"`
	Timestamp time.Time `json:"-"`
	ProductID string    `json:"-"`
}

func (material Material) NormalizedNumber() *string {
	if material.Number == "" {
		return nil
	}

	number := strings.TrimRight(strings.TrimLeft(material.Number, "0-"), "0-")
	return &number
}

const ModificationDelayedDeparture = 10
const ModificationDelayedArrival = 11
const ModificationChangedDeparturePlatform = 20
const ModificationChangedArrivalPlatform = 21
const ModificationDeparturePlatformAllocated = 22
const ModificationArrivalPlatformAllocated = 23
const ModificationExtraTrain = 24
const ModificationCancelledTrain = 25
const ModificationChangedStopPattern = 30
const ModificationExtraDeparture = 31
const ModificationCancelledDeparture = 32
const ModificationDiverted = 33
const ModificationRouteShortened = 34
const ModificationRouteExtended = 35
const ModificationOriginRouteShortened = 36
const ModificationOriginRouteExtended = 37
const ModificationExtraArrival = 38
const ModificationCancelledArrival = 39
const ModificationStatusChange = 40
const ModificationChangedDestination = 41
const ModificationChangedOrigin = 42
const ModificationExtraThroughTrain = 43
const ModificationCancelledThroughTrain = 44
const ModificationNotActual = 50
const ModificationBusReplacement = 51
