package models

import "fmt"

// Translate returns the appropriate translation based on language
func Translate(remarkNL, remarkEN, language string) string {
	remark := remarkNL

	if language == "en" {
		remark = remarkEN
	}

	return remark
}

// TranslateStations returns the appropriate translation based on language
func TranslateStations(remarkNL, remarkEN string, stations []Station, language string) string {
	translation := Translate(remarkNL, remarkEN, language)

	return fmt.Sprintf(translation, stationsStringTranslated(stations, language))
}

func stationsStringTranslated(stations []Station, language string) string {
	stationsText := ""

	for index, station := range stations {
		if index > 0 {
			if index < len(stations)-1 {
				stationsText += ", "
			} else {
				if language == "en" {
					stationsText += " and "
				} else {
					stationsText += " en "
				}
			}
		}
		stationsText += station.NameLong
	}

	return stationsText
}
