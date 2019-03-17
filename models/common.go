package models

import "time"

// Station is a station in the railway network. It has a code and 3 names (short, medium, long)
type Station struct {
	Code       string
	NameShort  string
	NameMedium string
	NameLong   string
}

// Modification is a change (to the schedule) which is communicated to travellers
type Modification struct {
	ModificationType int
	CauseShort       string
	CauseLong        string
	Station          Station
}

// Material is the physical train unit
type Material struct {
	NaterialType       string
	Number             string
	Position           int
	DestinationActual  Station
	DestinationPlanned Station
	Accessible         bool
	RemainsBehind      bool
}

// StoreItem is for shared fields like ID, timestamp etc.
type StoreItem struct {
	ID        string
	Timestamp time.Time
	ProductID string
}
