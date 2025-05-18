package models

import (
	"reflect"
	"testing"
	"time"
)

// TestRealDepartureTime tests the real departure time based on the planned time plus delay
func TestRealDepartureTime(t *testing.T) {
	var departure Departure

	departure.DepartureTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:34:00+01:00")

	expectedTime, _ := time.Parse(time.RFC3339, "2019-01-27T12:34:00+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}

	departure.Delay = 30 // 30 seconds
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:34:30+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}

	departure.Delay = 3600 // 1 hour
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T13:34:00+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}

	departure.Delay = -120 // Negative: -2 minutes
	expectedTime, _ = time.Parse(time.RFC3339, "2019-01-27T12:32:00+01:00")

	if !departure.RealDepartureTime().Equal(expectedTime) {
		t.Errorf("Real departure time %v does not match expected value %v", departure.RealDepartureTime(), expectedTime)
	}
}

func TestDeparturePlatformChanged(t *testing.T) {
	tables := []struct {
		planned string
		actual  string
		changed bool
	}{
		{"", "", false},
		{"4", "4", false},
		{"4", "5", true},
		{"", "4", false},
		{"4", "", true},
	}

	for _, table := range tables {
		var departure Departure
		departure.PlatformPlanned = table.planned
		departure.PlatformActual = table.actual

		if table.changed {
			if !departure.PlatformChanged() {
				t.Errorf("Planned platform \"%s\" is different from actual platform \"%s\", but not changed", departure.PlatformPlanned, departure.PlatformActual)
			}
		} else {
			if departure.PlatformChanged() {
				t.Errorf("Planned platform \"%s\" is equal to actual platform \"%s\", but is marked as changed", departure.PlatformPlanned, departure.PlatformActual)
			}
		}
	}
}

func TestGetDepartureID(t *testing.T) {
	var departure Departure

	departure.ServiceDate = "2019-01-27"
	departure.ServiceNumber = "301234"
	departure.ServiceID = "1234"
	departure.Station.Code = "RTD"
	departure.GenerateID()

	expected := "2019-01-27-1234-RTD"
	if departure.ID != expected {
		t.Errorf("Wrong departure ID, expected %s, got %s", expected, departure.ID)
	}
}

func TestDepartureRemarks(t *testing.T) {
	tables := []struct {
		modifications     []Modification
		wingModifications []Modification
		remarks           []string
		tips              []string
	}{
		{
			[]Modification{{ModificationType: ModificationCancelledDeparture}},
			[]Modification{},
			[]string{"Trein rijdt niet"},
			[]string{},
		},
		{
			[]Modification{{ModificationType: ModificationDiverted}, {ModificationType: ModificationChangedDeparturePlatform}},
			[]Modification{},
			[]string{"Rijdt via een andere route", "Gewijzigd vertrekspoor"},
			[]string{},
		},
		{
			[]Modification{{ModificationType: ModificationDiverted}, {ModificationType: ModificationChangedDeparturePlatform}},
			[]Modification{{ModificationType: ModificationChangedDeparturePlatform}},
			[]string{"Rijdt via een andere route", "Gewijzigd vertrekspoor"},
			[]string{},
		},
	}

	for _, table := range tables {
		var departure Departure
		var wing TrainWing

		departure.Modifications = table.modifications
		wing.Modifications = table.wingModifications

		departure.TrainWings = append(departure.TrainWings, wing)
		remarks, tips := departure.GetRemarksTips("nl")

		if !reflect.DeepEqual(table.remarks, remarks) {
			t.Errorf("Remarks: expected %s, got %s", table.remarks, remarks)
		}

		if !reflect.DeepEqual(table.tips, tips) {
			t.Errorf("Remarks: expected %s, got %s", table.tips, tips)
		}
	}
}

func TestDepartureRemarksTips(t *testing.T) {
	tables := []struct {
		departure Departure
		remarks   []string
		tips      []string
	}{
		{
			Departure{
				DoNotBoard: true,
			},
			[]string{"Niet instappen"},
			[]string{},
		},
		{
			Departure{
				ReservationRequired: true,
			},
			[]string{},
			[]string{"Reservering verplicht"},
		},
		{
			Departure{
				ReservationRequired: false,
				WithSupplement:      true,
				SpecialTicket:       true,
			},
			[]string{},
			[]string{"Toeslag verplicht", "Bijzonder ticket"},
		},
	}

	for _, table := range tables {
		remarks, tips := table.departure.GetRemarksTips("nl")

		if !reflect.DeepEqual(table.remarks, remarks) {
			t.Errorf("Remarks: expected %s, got %s", table.remarks, remarks)
		}

		if !reflect.DeepEqual(table.tips, tips) {
			t.Errorf("Remarks: expected %s, got %s", table.tips, tips)
		}
	}
}

func TestDepartureMaterialRemarksRemainsBehind(t *testing.T) {
	var departure Departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType:  "ICM",
						Number:        "1234",
						RemainsBehind: true,
					},
				},
			},
		},
	}

	remarks, _ := departure.GetRemarksTips("nl")

	if remarks[0] != "Treinstel 1234 blijft achter in Rotterdam C." {
		t.Errorf("Remarks: expected %s, got %s", "Treinstel 1234 blijft achter in Rotterdam C.", remarks[0])
	}

	departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType:  "ICM",
						Number:        "1234",
						RemainsBehind: true,
					},
					{
						NaterialType:  "ICM",
						Number:        "4321",
						RemainsBehind: false,
					},
					{
						NaterialType:  "ICM",
						Number:        "2345",
						RemainsBehind: true,
					},
				},
			},
		},
	}

	remarks, _ = departure.GetRemarksTips("nl")

	if remarks[0] != "Treinstellen 1234, 2345 blijven in Rotterdam C." {
		t.Errorf("Remarks: expected %s, got %s", "Treinstellen 1234, 2345 blijven in Rotterdam C.", remarks[0])
	}
}

func TestDepartureMaterialTipsAdded(t *testing.T) {
	var departure Departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType: "ICM",
						Number:       "1234",
						Added:        true,
					},
				},
			},
		},
	}

	_, tips := departure.GetRemarksTips("nl")

	if tips[0] != "Treinstel 1234 wordt op dit station bijgeplaatst" {
		t.Errorf("Tips: expected %s, got %s", "Treinstel 1234 wordt op dit station bijgeplaatst", tips[0])
	}

	departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType: "ICM",
						Number:       "1234",
						Added:        true,
					},
					{
						NaterialType: "ICM",
						Number:       "4321",
						Added:        false,
					},
					{
						NaterialType: "ICM",
						Number:       "2345",
						Added:        true,
					},
				},
			},
		},
	}

	_, tips = departure.GetRemarksTips("nl")

	if tips[0] != "Treinstellen 1234, 2345 worden op dit station bijgeplaatst" {
		t.Errorf("Tips: expected %s, got %s", "Treinstellen 1234, 2345 worden op dit station bijgeplaatst", tips[0])
	}
}

func TestDepartureMaterialRemarksClosed(t *testing.T) {
	var departure Departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType: "ICM",
						Number:       "1234",
						Closed:       true,
					},
				},
			},
		},
	}

	remarks, _ := departure.GetRemarksTips("nl")

	if remarks[0] != "Treinstel 1234: niet instappen" {
		t.Errorf("Remarks: expected %s, got %s", "Treinstel 1234: niet instappen", remarks[0])
	}

	departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType: "ICM",
						Number:       "1234",
						Closed:       true,
					},
					{
						NaterialType: "ICM",
						Number:       "4321",
						Closed:       false,
					},
					{
						NaterialType: "ICM",
						Number:       "2345",
						Closed:       true,
					},
				},
			},
		},
	}

	remarks, _ = departure.GetRemarksTips("nl")

	if remarks[0] != "Treinstellen 1234, 2345: niet instappen" {
		t.Errorf("Remarks: expected %s, got %s", "Treinstellen 1234, 2345: niet instappen", remarks[0])
	}
}

func TestDepartureMaterialRemarksRemainsBehindAlreadyRemoved(t *testing.T) {
	var departure Departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType:   "ICM",
						Number:         "1234",
						RemainsBehind:  true,
						AlreadyRemoved: true,
					},
				},
			},
		},
	}

	remarks, _ := departure.GetRemarksTips("nl")

	if len(remarks) > 0 {
		t.Errorf("Remarks: expected no remarks, got %d", len(remarks))
	}

	departure = Departure{
		Station: Station{
			Code:       "RTD",
			NameMedium: "Rotterdam C.",
		},
		TrainWings: []TrainWing{
			{
				Material: []Material{
					{
						NaterialType:   "ICM",
						Number:         "1234",
						RemainsBehind:  true,
						AlreadyRemoved: true,
					},
					{
						NaterialType:  "ICM",
						Number:        "4321",
						RemainsBehind: false,
					},
					{
						NaterialType:  "ICM",
						Number:        "2345",
						RemainsBehind: true,
					},
				},
			},
		},
	}

	remarks, _ = departure.GetRemarksTips("nl")

	if remarks[0] != "Treinstel 2345 blijft achter in Rotterdam C." {
		t.Errorf("Remarks: expected %s, got %s", "Treinstel 2345 blijft achter in Rotterdam C.", remarks[0])
	}
}
