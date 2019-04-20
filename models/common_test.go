package models

import "testing"

func TestNormalizedMaterialNumber(t *testing.T) {
	tables := []struct {
		materialNumber   string
		normalizedNumber string
	}{
		{"", ""},
		{"000000-09547-0", "9547"},
		{"000000-16475-0", "16475"},
		{"000001-86012-0", "186.012"},
		{"RdTrein", "RdTrein"},
	}

	for _, table := range tables {
		var material Material
		var normalizedNumber *string

		material.Number = table.materialNumber
		normalizedNumber = material.NormalizedNumber()

		if table.normalizedNumber == "" {
			if normalizedNumber != nil {
				t.Errorf("Invalid material number: got %s, expected nil", *normalizedNumber)
			}
		} else {
			if *normalizedNumber != table.normalizedNumber {
				t.Errorf("Invalid material number: got %s, expected %s", *normalizedNumber, table.normalizedNumber)
			}
		}
	}
}

func TestStationsString(t *testing.T) {
	stationRtd := Station{Code: "RTD", NameShort: "R'dam C.", NameMedium: "Rotterdam C.", NameLong: "Rotterdam Centraal"}
	stationAsd := Station{Code: "ASD", NameShort: "A'dam C.", NameMedium: "Amsterdam C.", NameLong: "Amsterdam Centraal"}
	stationUt := Station{Code: "UT", NameShort: "Utrecht C.", NameMedium: "Utrecht C.", NameLong: "Utrecht Centraal"}

	tables := []struct {
		stations     []Station
		separator    string
		outputShort  string
		outputMedium string
		outputLong   string
	}{
		{nil, "/", "", "", ""},
		{[]Station{stationRtd}, "/", "R'dam C.", "Rotterdam C.", "Rotterdam Centraal"},
		{[]Station{stationRtd, stationUt}, "/", "R'dam C./Utrecht C.", "Rotterdam C./Utrecht C.", "Rotterdam Centraal/Utrecht Centraal"},
		{[]Station{stationRtd, stationAsd, stationUt}, ", ", "R'dam C., A'dam C., Utrecht C.", "Rotterdam C., Amsterdam C., Utrecht C.", "Rotterdam Centraal, Amsterdam Centraal, Utrecht Centraal"},
	}

	for _, table := range tables {
		stationsStringShort := stationsShortString(table.stations, table.separator)
		stationsStringMedium := stationsMediumString(table.stations, table.separator)
		stationsStringLong := stationsLongString(table.stations, table.separator)

		if stationsStringShort != table.outputShort {
			t.Errorf("Invalid stationsShortString: got %s, expected %s", stationsStringShort, table.outputShort)
		}
		if stationsStringMedium != table.outputMedium {
			t.Errorf("Invalid stationsMediumString: got %s, expected %s", stationsStringMedium, table.outputMedium)
		}
		if stationsStringLong != table.outputLong {
			t.Errorf("Invalid stationsLongString: got %s, expected %s", stationsStringLong, table.outputLong)
		}

		// Test station codes:
		stationCodes := stationCodes(table.stations)

		if len(stationCodes) != len(table.stations) {
			t.Error("Number of station codes does not match")
		}

		for index, station := range table.stations {
			if stationCodes[index] != station.Code {
				t.Errorf("Station code %s does not match given code %s", stationCodes[index], station.Code)
			}
		}
	}
}
