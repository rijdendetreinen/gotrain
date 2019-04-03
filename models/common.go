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
