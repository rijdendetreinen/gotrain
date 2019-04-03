package models

import "fmt"

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

	case ModificationCancelledArrival, ModificationCancelledDeparture, ModificationCancelledTrain:
		return modification.remarkWithCause("Trein rijdt niet", "Cancelled", language), true

	case ModificationChangedDeparturePlatform:
		// TODO: pass platform as argument
		return modification.remarkTranslation("Gewijzigd vertrekspoor", "Changed departure platform", language), true

	case ModificationChangedArrivalPlatform:
		// TODO: pass platform as argument
		return modification.remarkTranslation("Gewijzigd aankomstspoor", "Changed arrival platform", language), true

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

	case ModificationOriginRouteExtended, ModificationChangedOrigin, ModificationOriginRouteShortened:
		// TODO: pass origin station
		return modification.remarkWithCause("Afwijkende herkomst", "Different origin", language), true

	case ModificationChangedDestination:
		return modification.remarkWithStation("Let op, rijdt naar %s", "Attention, train goes to %s", language), true

	case ModificationNotActual:
		return modification.remarkTranslation("Geen actuele informatie", "Information is not real-time", language), true

	case ModificationBusReplacement:
		return modification.remarkTranslation("Bus in plaats van trein", "Bus replaces train", language), true

	}

	return "", false
}

func (modification Modification) remarkWithCause(remarkNL, remarkEN, language string) string {
	remark := modification.remarkTranslation(remarkNL, remarkEN, language)

	if modification.CauseLong != "" {
		cause := modification.CauseLong

		if language == "en" {
			// TODO: translate cause
		}

		remark = remark + " " + cause
	}

	return remark
}

func (modification Modification) remarkTranslation(remarkNL, remarkEN, language string) string {
	remark := remarkNL

	if language == "en" {
		remark = remarkEN
	}

	return remark
}

func (modification Modification) remarkWithStation(remarkNL, remarkEN, language string) string {
	remark := modification.remarkWithCause(remarkNL, remarkEN, language)

	return fmt.Sprintf(remark, modification.Station.NameLong)
}
