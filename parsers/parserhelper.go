package parsers

import (
	"strconv"
	"time"

	"github.com/beevik/etree"
	"github.com/rickb777/date/period"
	"github.com/rijdendetreinen/gotrain/models"
)

// ParseInfoPlusBoolean returns true when this is a InfoPlus boolean type which is true
func ParseInfoPlusBoolean(element *etree.Element) bool {
	if element == nil {
		return false
	}

	return element.Text() == "J"
}

// ParseInfoPlusModifications parses a list of modifications
func ParseInfoPlusModifications(element *etree.Element) []models.Modification {
	var modifications []models.Modification

	for _, modificationElement := range element.SelectElements("Wijziging") {
		var modification models.Modification

		modification.ModificationType, _ = strconv.Atoi(modificationElement.SelectElement("WijzigingType").Text())

		causeShort := modificationElement.SelectElement("WijzigingOorzaakKort")
		causeLong := modificationElement.SelectElement("WijzigingOorzaakLang")
		station := modificationElement.SelectElement("WijzigingStation")

		if causeShort != nil {
			modification.CauseShort = causeShort.Text()
			modification.CauseLong = causeLong.Text()
		}

		if station != nil {
			modification.Station = ParseInfoPlusStation(station)
		}

		modifications = append(modifications, modification)
	}

	return modifications
}

// ParseInfoPlusStation translates an XML InfoPlus station to a Station object
func ParseInfoPlusStation(element *etree.Element) models.Station {
	var station models.Station

	station.Code = element.SelectElement("StationCode").Text()
	station.NameShort = element.SelectElement("KorteNaam").Text()
	station.NameMedium = element.SelectElement("MiddelNaam").Text()
	station.NameLong = element.SelectElement("LangeNaam").Text()

	return station
}

// ParseInfoPlusStations process multiple station elements and returns them as a slice
func ParseInfoPlusStations(elements []*etree.Element) (stations []models.Station) {
	for _, element := range elements {
		stations = append(stations, ParseInfoPlusStation(element))
	}

	return stations
}

// ParseWhenAttribute filters a list of elements on an attribute with a given value. Returns a single element, or nil.
func ParseWhenAttribute(element *etree.Element, tag, attribute, value string) *etree.Element {
	for _, childElement := range element.SelectElements(tag) {
		if childElement.SelectAttrValue(attribute, "") == value {
			return childElement
		}
	}

	return nil
}

// ParseWhenAttributeMulti filters a list of elements on an attribute with a given value. Returns a slice with elements
func ParseWhenAttributeMulti(element *etree.Element, tag, attribute, value string) []*etree.Element {
	var elements []*etree.Element

	for _, childElement := range element.SelectElements(tag) {
		if childElement.SelectAttrValue(attribute, "") == value {
			elements = append(elements, childElement)
		}
	}

	return elements
}

// ParseOptionalText gets the text from an element, or returns an empty string when the element is nil.
func ParseOptionalText(element *etree.Element) string {
	if element != nil {
		return element.Text()
	}

	return ""
}

// ParseInfoPlusDateTime translates an element with a date/time to a time.Time struct
func ParseInfoPlusDateTime(element *etree.Element) time.Time {

	if element == nil {
		return time.Time{}
	}

	datetime, error := time.Parse(time.RFC3339, element.Text())

	if error != nil {
		return time.Time{}
	}
	return datetime
}

// ParseInfoPlusPlatform translates a platform element to a string
func ParseInfoPlusPlatform(elements []*etree.Element) string {
	if len(elements) == 0 {
		return ""
	}

	platform := ""

	for index, element := range elements {
		if index > 0 {
			platform = platform + "/"
		}
		platform = platform + element.SelectElement("SpoorNummer").Text()

		phaseLetter := element.SelectElement("SpoorFase")
		if phaseLetter != nil {
			platform = platform + phaseLetter.Text()
		}
	}

	return platform
}

// ParseInfoPlusDuration translates an element with a duration (i.e., delays) to seconds
func ParseInfoPlusDuration(element *etree.Element) int {
	if element == nil {
		return 0
	}

	delay, error := period.Parse(element.Text())

	if error != nil {
		return 0
	}

	return delay.Seconds() + delay.Minutes()*60 + delay.Hours()*3600
}
