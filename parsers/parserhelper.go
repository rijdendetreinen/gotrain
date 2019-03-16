package parsers

import (
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

func ParseInfoPlusModifications(element *etree.Element) []models.Modification {
	var modifications []models.Modification
	// TODO:Implement

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

func ParseWhenAttribute(element *etree.Element, tag, attribute, value string) *etree.Element {
	for _, childElement := range element.SelectElements(tag) {
		if childElement.SelectAttrValue(attribute, "") == value {
			return childElement
		}
	}

	return nil
}

func ParseWhenAttributeMulti(element *etree.Element, tag, attribute, value string) []*etree.Element {
	var elements []*etree.Element

	for _, childElement := range element.SelectElements(tag) {
		if childElement.SelectAttrValue(attribute, "") == value {
			elements = append(elements, childElement)
		}
	}

	return elements
}

func ParseOptionalText(element *etree.Element) string {
	if element != nil {
		return element.Text()
	}

	return ""
}

func ParseInfoPlusDateTime(element *etree.Element) *time.Time {
	if element == nil {
		return nil
	}

	datetime, error := time.Parse(time.RFC3339, element.Text())

	if error != nil {
		return nil
	}
	return &datetime
}

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
