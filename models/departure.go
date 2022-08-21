package models

import (
	"fmt"
	"strings"
	"time"
)

// Departure is a train service which departs from a single station
type Departure struct {
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
	StationsPlanned    []Station
	Material           []Material
	Modifications      []Modification
}

// GenerateID generates an ID for this departure
func (departure *Departure) GenerateID() {
	departure.ID = departure.ServiceDate + "-" + departure.ServiceID + "-" + departure.Station.Code
}

// RealDepartureTime returns the actual departure time, including delay
func (departure Departure) RealDepartureTime() time.Time {
	delayDuration := time.Second * time.Duration(departure.Delay)
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

// ActualViaStationsString returns a string of all actual via stations (medium name)
func (departure Departure) ActualViaStationsString() string {
	return stationsMediumString(departure.ViaActual, ", ")
}

// PlannedViaStationsString returns a string of all planned via stations (medium name)
func (departure Departure) PlannedViaStationsString() string {
	return stationsMediumString(departure.ViaPlanned, ", ")
}

// DestinationActualString returns a string of all actual destination stations (medium name)
func (wing TrainWing) DestinationActualString() string {
	return stationsMediumString(wing.DestinationActual, ", ")
}

// DestinationPlannedString returns a string of all planned destination stations (medium name)
func (wing TrainWing) DestinationPlannedString() string {
	return stationsMediumString(wing.DestinationPlanned, ", ")
}

// Translation provides a translation for this tip
func (tip BoardingTip) Translation(language string) string {
	translation := Translate("%s %s naar %s (spoor %s) is eerder in %s", "%s %s to %s (platform %s) reaches %s sooner", language)

	return fmt.Sprintf(translation, tip.TrainTypeCode, tip.DepartureTime.Local().Format("15:04"), tip.Destination.NameLong, tip.DeparturePlatform, tip.ExitStation.NameLong)
}

// Translation provides a translation for this tip
func (tip ChangeTip) Translation(language string) string {
	translation := Translate("Voor %s overstappen in %s", "For %s, change at %s", language)

	return fmt.Sprintf(translation, tip.Destination.NameLong, tip.ChangeStation.NameLong)
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

// GetRemarksTips returns all (translated) remarks and tips, both travel tips and tips based on departure flags
func (departure Departure) GetRemarksTips(language string) (remarks, tips []string) {
	remarks = GetRemarks(departure.Modifications, language)
	tips = make([]string, 0)

	if !departure.Cancelled {
		if departure.DoNotBoard {
			remarks = append(remarks, Translate("Niet instappen", "Do not board", language))
		}
		if departure.ReservationRequired {
			tips = append(tips, Translate("Reservering verplicht", "Reservation required", language))
		}
		if departure.WithSupplement {
			tips = append(tips, Translate("Toeslag verplicht", "Supplement required", language))
		}
		if departure.SpecialTicket {
			tips = append(tips, Translate("Bijzonder ticket", "Special ticket", language))
		}

		// Wing remarks:
		for _, wing := range departure.TrainWings {
			wingRemarks := GetFilteredRemarks(wing.Modifications, language)

			for _, wingRemark := range wingRemarks {
				if len(departure.TrainWings) > 1 {
					wingRemark = wing.DestinationPlanned[0].NameMedium + ": " + wingRemark
				}

				remarks = append(remarks, wingRemark)
			}
		}

		// Translate all tips:
		for _, tip := range departure.TravelTips {
			tips = append(tips, tip.Translation(language))
		}
		for _, tip := range departure.BoardingTips {
			tips = append(tips, tip.Translation(language))
		}
		for _, tip := range departure.ChangeTips {
			tips = append(tips, tip.Translation(language))
		}

		// Check for closed material:
		var closedMaterialUnits = make([]string, 0)
		var leftBehindMaterialUnits = make([]string, 0)
		var addedMaterialUnits = make([]string, 0)

		for _, wing := range departure.TrainWings {
			for _, material := range wing.Material {
				if material.Number == "" {
					continue
				}

				if material.Closed {
					closedMaterialUnits = append(closedMaterialUnits, *material.NormalizedNumber())
				}
				if material.RemainsBehind && !material.AlreadyRemoved {
					leftBehindMaterialUnits = append(leftBehindMaterialUnits, *material.NormalizedNumber())
				}
				if material.Added {
					addedMaterialUnits = append(addedMaterialUnits, *material.NormalizedNumber())
				}
			}
		}

		if len(closedMaterialUnits) > 0 {
			if len(closedMaterialUnits) == 1 {
				remarks = append(remarks, fmt.Sprintf(Translate("Treinstel %s: niet instappen", "Coach %s: do not board", language), strings.Join(closedMaterialUnits, ", ")))
			} else {
				remarks = append(remarks, fmt.Sprintf(Translate("Treinstellen %s: niet instappen", "Coaches %s: do not board", language), strings.Join(closedMaterialUnits, ", ")))
			}
		}

		if len(leftBehindMaterialUnits) > 0 {
			if len(leftBehindMaterialUnits) == 1 {
				remarks = append(remarks, fmt.Sprintf(Translate("Treinstel %s blijft achter in %s", "Coach %s remains at %s", language), strings.Join(leftBehindMaterialUnits, ", "), departure.Station.NameMedium))
			} else {
				remarks = append(remarks, fmt.Sprintf(Translate("Treinstellen %s blijven in %s", "Coaches %s remain at %s", language), strings.Join(leftBehindMaterialUnits, ", "), departure.Station.NameMedium))
			}
		}

		if len(addedMaterialUnits) > 0 {
			if len(addedMaterialUnits) == 1 {
				tips = append(tips, fmt.Sprintf(Translate("Trein wordt op dit station verlengd. Treinstel %s wordt op dit station bijgeplaatst", "Coach %s is added to the train at this station", language), strings.Join(addedMaterialUnits, ", ")))
			} else {
				tips = append(tips, fmt.Sprintf(Translate("Trein wordt op dit station verlengd. Treinstellen %s worden op dit station bijgeplaatst", "Coaches %s are added to the train at this station", language), strings.Join(addedMaterialUnits, ", ")))
			}
		}

		// Check for material destinations:
		for _, wing := range departure.TrainWings {
			differentTerminus := make(map[string][]Material)

			for _, material := range wing.Material {
				if len(wing.DestinationActual) > 0 && material.DestinationActual.Code != wing.DestinationActual[0].Code {
					// Different terminus, and not explained by modifications 82-84:
					if material.Number != "" && !material.Closed && !material.RemainsBehind {
						terminus := material.DestinationActual.NameLong

						differentTerminus[terminus] = append(differentTerminus[terminus], material)
					}
				}
			}

			for terminus, units := range differentTerminus {
				if len(units) == 1 {
					remarks = append(remarks, fmt.Sprintf(Translate("Treinstel %s rijdt niet verder dan %s", "Coach %s terminates at %s", language), *units[0].NormalizedNumber(), terminus))
				} else {
					strs := make([]string, len(units))
					for i, v := range units {
						strs[i] = *v.NormalizedNumber()
					}
					unitsString := strings.Join(strs, ", ")

					remarks = append(remarks, fmt.Sprintf(Translate("Treinstellen %s rijden niet verder dan %s", "Coaches %s terminate at %s", language), unitsString, terminus))
				}
			}
		}
	}

	if departure.ServiceName != "" {
		tips = append(tips, departure.ServiceName)
	}

	return
}
