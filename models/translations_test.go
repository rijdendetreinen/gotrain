package models

import (
	"testing"
)

func TestTranslate(t *testing.T) {
	tables := []struct {
		inputNL        string
		inputEN        string
		language       string
		outputLanguage string
	}{
		{"Hallo wereld", "Hello world", "nl", "nl"},
		{"Hallo wereld", "Hello world", "en", "en"},
		{"Hallo wereld", "Hello world", "", "nl"},
		{"Hallo wereld", "Hello world", "lala", "nl"},
	}

	for _, table := range tables {
		translation := Translate(table.inputNL, table.inputEN, table.language)

		if table.outputLanguage == "en" {
			if translation != table.inputEN {
				t.Errorf("Expected English translation, received: %s", translation)
			}
		} else {
			if translation != table.inputNL {
				t.Errorf("Expected Dutch translation, received: %s", translation)
			}
		}
	}
}

func TestTranslateStations(t *testing.T) {
	stationRtd := Station{Code: "RTD", NameShort: "R'dam C.", NameMedium: "Rotterdam C.", NameLong: "Rotterdam Centraal"}
	stationAsd := Station{Code: "ASD", NameShort: "A'dam C.", NameMedium: "Amsterdam C.", NameLong: "Amsterdam Centraal"}
	stationUt := Station{Code: "UT", NameShort: "Utrecht C.", NameMedium: "Utrecht C.", NameLong: "Utrecht Centraal"}

	tables := []struct {
		inputNL  string
		inputEN  string
		language string
		output   string
		stations []Station
	}{
		{"Welkom in %s!", "Welcome to %s!", "nl", "Welkom in Rotterdam Centraal!", []Station{stationRtd}},
		{"Welkom in %s!", "Welcome to %s!", "en", "Welcome to Rotterdam Centraal!", []Station{stationRtd}},
		{"Trein naar %s", "Train to %s", "nl", "Trein naar Rotterdam Centraal en Amsterdam Centraal", []Station{stationRtd, stationAsd}},
		{"Trein naar %s", "Train to %s", "en", "Train to Rotterdam Centraal and Amsterdam Centraal", []Station{stationRtd, stationAsd}},
		{"Trein naar %s", "Train to %s", "nl", "Trein naar Rotterdam Centraal, Utrecht Centraal en Amsterdam Centraal", []Station{stationRtd, stationUt, stationAsd}},
		{"Trein naar %s", "Train to %s", "en", "Train to Rotterdam Centraal, Utrecht Centraal and Amsterdam Centraal", []Station{stationRtd, stationUt, stationAsd}},
	}

	for _, table := range tables {
		translation := TranslateStations(table.inputNL, table.inputEN, table.stations, table.language)

		if table.output != translation {
			t.Errorf("Wrong translation: expected %s, received: %s", table.output, translation)
		}
	}
}

func TestTranslateCause(t *testing.T) {
	tables := []struct {
		cause       string
		translation string
	}{
		{"door werkzaamheden", "due to engineering work"},
		{"door het onschadelijk maken van een bom uit de Tweede Wereldoorlog", "due to defusing a bomb from World War II"},
		{"door een onbekende vertaling", "door een onbekende vertaling"},
		{"", ""},
	}

	for _, table := range tables {
		translation := TranslateCause(table.cause)

		if translation != table.translation {
			t.Errorf("Wrong translation for cause: expected %s, received: %s", table.translation, translation)
		}
	}
}
