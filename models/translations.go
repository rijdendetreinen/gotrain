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

	return fmt.Sprintf(translation, stationsMediumString(stations, ", "))
}
