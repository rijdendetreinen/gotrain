package models

import "fmt"

// ModificationDelayedDeparture is the status for delays on departing trains
const ModificationDelayedDeparture = 10

// ModificationDelayedArrival is the status for delays on arriving trains
const ModificationDelayedArrival = 11

// ModificationChangedDepartureTime when the departure time has been changed
const ModificationChangedDepartureTime = 12

// ModificationChangedArrivalTime when the arrival time has been changed
const ModificationChangedArrivalTime = 13

// ModificationChangedDeparturePlatform is for changed departure platforms (compared to the schedule)
const ModificationChangedDeparturePlatform = 20

// ModificationChangedArrivalPlatform is for changed arrival platforms (compared to the schedule)
const ModificationChangedArrivalPlatform = 21

// ModificationDeparturePlatformAllocated when a departure platform is allocated (only DVM stations, i.e. Schiphol Airport)
const ModificationDeparturePlatformAllocated = 22

// ModificationArrivalPlatformAllocated when a arrival platform is allocated (only DVM stations, i.e. Schiphol Airport)
const ModificationArrivalPlatformAllocated = 23

// ModificationExtraTrain for extra trains (not in the normal schedule)
const ModificationExtraTrain = 24

// ModificationCancelledTrain Train is cancelled
const ModificationCancelledTrain = 25

// ModificationChangedStopPattern when an extra stop is added or a stop is changed
const ModificationChangedStopPattern = 30

// ModificationExtraDeparture extra departure for this station (train usually doesn't depart here)
const ModificationExtraDeparture = 31

// ModificationCancelledDeparture cancelled departure (train was scheduled to stop here)
const ModificationCancelledDeparture = 32

// ModificationDiverted train is diverted, takes another route
const ModificationDiverted = 33

// ModificationRouteShortened train terminates early
const ModificationRouteShortened = 34

// ModificationRouteExtended train continues beyond normal destination
const ModificationRouteExtended = 35

// ModificationOriginRouteShortened departs from a later stop than usual
const ModificationOriginRouteShortened = 36

// ModificationOriginRouteExtended train departs from an earlier stop than the normal origin station
const ModificationOriginRouteExtended = 37

// ModificationExtraArrival extra arrival, i.e. train doesn't normally arrive at this station
const ModificationExtraArrival = 38

// ModificationCancelledArrival train was scheduled to arrive, but didn't
const ModificationCancelledArrival = 39

// ModificationStatusChange train status changed (i.e. from 'at platform' to 'departed')
const ModificationStatusChange = 40

// ModificationChangedDestination destination has been changed (not extended nor shortened, it's a new route)
const ModificationChangedDestination = 41

// ModificationChangedOrigin origin has been changed (not an extension or shortening of the normal route)
const ModificationChangedOrigin = 42

// ModificationExtraThroughTrain extra train which just passes by (doesn't stop at this station)
const ModificationExtraThroughTrain = 43

// ModificationCancelledThroughTrain passing-through train has been cancelled (train would not stop at this station)
const ModificationCancelledThroughTrain = 44

// ModificationNotActual indicator for non-realtime information, e.g. replacement buses or some foreign trains
const ModificationNotActual = 50

// ModificationBusReplacement service is a bus replacement service
const ModificationBusReplacement = 51

// ModificationSprinterRunsAsIntercity when an intercity train is running with Sprinter stock
const ModificationSprinterRunsAsIntercity = 80

// ModificationSprinterRunsAsIntercity when an Sprinter train is running with intercity stock
const ModificationIntercityRunsAsSprinter = 81

// ModificationMaterialClosed when material is closed for passengers
const ModificationMaterialClosed = 82

// ModificationMaterialAdded when material is added on a station
const ModificationMaterialAdded = 83

// ModificationMaterialLeftBehind when material is left behind on a station
const ModificationMaterialLeftBehind = 84

// ModificationMaterialAlreadyRemoved when the material that was left behind has already been removed (driven away)
const ModificationMaterialAlreadyRemoved = 85

// Modification is a change (to the schedule) which is communicated to travellers
type Modification struct {
	ModificationType int     `json:"type"`
	CauseShort       string  `json:"cause_short"`
	CauseLong        string  `json:"cause_long"`
	Station          Station `json:"station"`
}

// Remark translates a Modification object to a translated text message
func (modification Modification) Remark(language string) (string, bool) {
	switch modification.ModificationType {
	case ModificationDelayedDeparture:
		if modification.CauseLong != "" {
			// Only translate when there is a cause for the delay:
			return modification.remarkWithCause("Later vertrek", "Delayed", language), true
		}

	case ModificationDelayedArrival:
		if modification.CauseLong != "" {
			// Only translate when there is a cause for the delay:
			return modification.remarkWithCause("Latere aankomst", "Delayed", language), true
		}

	case ModificationChangedDepartureTime:
		return modification.remarkWithCause("Vertrektijd gewijzigd", "Departure time changed", language), true

	case ModificationChangedArrivalTime:
		return modification.remarkWithCause("Aankomsttijd gewijzigd", "Arrival time changed", language), true

	case ModificationCancelledArrival, ModificationCancelledDeparture, ModificationCancelledTrain:
		return modification.remarkWithCause("Trein rijdt niet", "Cancelled", language), true

	case ModificationChangedDeparturePlatform:
		// TODO: pass platform as argument
		return Translate("Gewijzigd vertrekspoor", "Changed departure platform", language), true

	case ModificationChangedArrivalPlatform:
		// TODO: pass platform as argument
		return Translate("Gewijzigd aankomstspoor", "Changed arrival platform", language), true

	case ModificationChangedStopPattern:
		return modification.remarkWithCause("Gewijzigde dienstregeling", "Schedule changed", language), true

	case ModificationExtraArrival, ModificationExtraDeparture, ModificationExtraTrain:
		return modification.remarkWithCause("Extra trein", "Extra train", language), true

	case ModificationDiverted:
		return modification.remarkWithCause("Rijdt via een andere route", "Train is diverted", language), true

	case ModificationRouteShortened:
		return modification.remarkWithStation("Rijdt niet verder dan %s", "Terminates at %s", language), true

	case ModificationRouteExtended:
		return modification.remarkWithStation("Rijdt verder naar %s", "Continues to %s", language), true

	case ModificationOriginRouteShortened:
		return modification.remarkWithStation("Trein begint op %s", "Train starts at %s", language), true

	case ModificationOriginRouteExtended, ModificationChangedOrigin:
		// TODO: pass origin station
		return modification.remarkWithCause("Afwijkende herkomst", "Different origin", language), true

	case ModificationChangedDestination:
		return modification.remarkWithStation("Let op, rijdt naar %s", "Attention, train goes to %s", language), true

	case ModificationNotActual:
		return Translate("Geen actuele informatie", "Information is not real-time", language), true

	case ModificationBusReplacement:
		return Translate("Bus in plaats van trein", "Bus replaces train", language), true

	case ModificationSprinterRunsAsIntercity:
		return Translate("Trein rijdt met sprintermaterieel", "This train runs with Sprinter stock", language), true

	case ModificationIntercityRunsAsSprinter:
		return Translate("Trein rijdt met intercity-materieel", "This train runs with intercity stock", language), true

	}

	return "", false
}

func (modification Modification) remarkWithCause(remarkNL, remarkEN, language string) string {
	remark := Translate(remarkNL, remarkEN, language)

	if modification.CauseLong != "" {
		cause := modification.CauseLong

		if language == "en" {
			cause = TranslateCause(cause)
		}

		remark = remark + " " + cause
	}

	return remark
}

func (modification Modification) remarkWithStation(remarkNL, remarkEN, language string) string {
	remark := modification.remarkWithCause(remarkNL, remarkEN, language)

	return fmt.Sprintf(remark, modification.Station.NameLong)
}

// GetRemarks translates a slice of Modification structs to remarks in the requested language
func GetRemarks(modifications []Modification, language string) []string {
	remarks := make([]string, 0)

	for _, modification := range modifications {
		remark, hasRemark := modification.Remark(language)

		if hasRemark {
			remarks = append(remarks, remark)
		}
	}

	return remarks
}

// GetFilteredRemarks filters out modification types for arrival/departure platforms (useful for wings)
func GetFilteredRemarks(modifications []Modification, language string) []string {
	remarks := make([]string, 0)

	for _, modification := range modifications {
		remark, hasRemark := modification.Remark(language)

		if modification.ModificationType == ModificationChangedArrivalPlatform || modification.ModificationType == ModificationChangedDeparturePlatform {
			continue
		}

		if hasRemark {
			remarks = append(remarks, remark)
		}
	}

	return remarks
}
