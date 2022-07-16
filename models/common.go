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

// Material is the physical train unit
type Material struct {
	NaterialType       string  `json:"type"`
	Number             string  `json:"number"`
	Position           int     `json:"position"`
	DestinationActual  Station `json:"destination_actual"`
	DestinationPlanned Station `json:"destination_planned"`
	Accessible         bool    `json:"accessible"`

	Closed        bool `json:"closed"`
	RemainsBehind bool `json:"remains_behind"`
	Added         bool `json:"added"`

	Modifications []Modification
}

// StoreItem is for shared fields like ID, timestamp etc.
type StoreItem struct {
	ID        string    `json:"-"`
	Timestamp time.Time `json:"-"`
	ProductID string    `json:"-"`
}

// NormalizedNumber returns a normal material number (i.e. it removes the 000000-0...-0 prefixes and suffixes).
// Example: 000000-09547-0 will be translated to 9547
func (material Material) NormalizedNumber() *string {
	if material.Number == "" {
		return nil
	}

	number := strings.TrimRight(strings.TrimRight(strings.TrimLeft(material.Number, "0-"), "0"), "0")
	number = strings.ReplaceAll(number, "-", "")

	// Translate TRAXX numbers (e.g., 186012 becomes 186.012)
	if len(number) == 6 {
		number = number[:3] + "." + number[3:]
	}

	return &number
}

func stationsLongString(stations []Station, separator string) string {
	stationsText := ""
	for index, station := range stations {
		if index > 0 {
			stationsText += separator
		}
		stationsText += station.NameLong
	}

	return stationsText
}

func stationsMediumString(stations []Station, separator string) string {
	stationsText := ""
	for index, station := range stations {
		if index > 0 {
			stationsText += separator
		}
		stationsText += station.NameMedium
	}

	return stationsText
}

func stationsShortString(stations []Station, separator string) string {
	stationsText := ""
	for index, station := range stations {
		if index > 0 {
			stationsText += separator
		}
		stationsText += station.NameShort
	}

	return stationsText
}

func stationCodes(stations []Station) []string {
	var stationCodes []string

	for _, station := range stations {
		stationCodes = append(stationCodes, station.Code)
	}

	return stationCodes
}
