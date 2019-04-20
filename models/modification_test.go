package models

import (
	"testing"
)

func TestRemarks(t *testing.T) {
	tables := []struct {
		modificationType int
		cause            string
		hasRemark        bool
		language         string
		remark           string
	}{
		{123, "", false, "", ""},                          // Does not exist
		{ModificationDelayedArrival, "", false, "", ""},   // No cause
		{ModificationDelayedDeparture, "", false, "", ""}, // No cause
		{ModificationCancelledDeparture, "", true, "nl", "Trein rijdt niet"},
		{ModificationCancelledDeparture, "", true, "en", "Cancelled"},
		{ModificationCancelledDeparture, "door een seinstoring", true, "nl", "Trein rijdt niet door een seinstoring"},
		{ModificationCancelledDeparture, "due to a signal failure", true, "en", "Cancelled due to a signal failure"},
		{ModificationBusReplacement, "door een seinstoring", true, "nl", "Bus in plaats van trein"},
		{ModificationStatusChange, "", false, "", ""}, // No remark
	}

	for _, table := range tables {
		var modification Modification
		modification.ModificationType = table.modificationType
		modification.CauseLong = table.cause

		remark, hasRemark := modification.Remark(table.language)

		if hasRemark && remark == "" {
			t.Error("Should not return an empty string when hasRemark is set")
		} else if !hasRemark && remark != "" {
			t.Error("Should return an empty string when hasRemark is false")
		}

		if hasRemark && !table.hasRemark {
			t.Error("Should not have a remark")
		} else if !hasRemark && table.hasRemark {
			t.Error("Should have a remark, but no remark returned")
		} else if hasRemark && table.hasRemark {
			if remark != table.remark {
				t.Errorf("Expected remark %s does not match with actual value %s", table.remark, remark)
			}
		}
	}
}

func TestDelayRemark(t *testing.T) {
	var modification Modification
	modification.ModificationType = ModificationDelayedDeparture

	_, hasRemark := modification.Remark("nl")

	if hasRemark {
		t.Error("Delayed modification should not have a remark unless a cause is given")
	}

	modification.CauseLong = "unit testing"
	modification.CauseShort = "testing"
	_, hasRemark = modification.Remark("nl")

	if !hasRemark {
		t.Error("Delayed modification should have a remark when a cause is given")
	}
}

func TestRemarksWithStation(t *testing.T) {
	tables := []struct {
		modificationType int
		stationName      string
		cause            string
		language         string
		remark           string
	}{
		{ModificationRouteShortened, "Rotterdam Centraal", "", "nl", "Rijdt niet verder dan Rotterdam Centraal"},
		{ModificationChangedDestination, "Rotterdam Centraal", "", "en", "Attention, train goes to Rotterdam Centraal"},
		{ModificationRouteExtended, "Rotterdam Centraal", "door werkzaamheden", "nl", "Rijdt verder naar Rotterdam Centraal door werkzaamheden"},
	}

	for _, table := range tables {
		var modification Modification
		modification.ModificationType = table.modificationType
		modification.CauseLong = table.cause
		modification.Station.NameLong = table.stationName

		remark, _ := modification.Remark(table.language)

		if remark != table.remark {
			t.Errorf("Expected remark %s does not match with actual value %s", table.remark, remark)
		}
	}
}

func TestGetRemarks(t *testing.T) {
	statusChange := Modification{ModificationType: ModificationStatusChange}
	delayWithCause := Modification{ModificationType: ModificationDelayedDeparture, CauseLong: "door testen"}
	delayNoCause := Modification{ModificationType: ModificationDelayedDeparture}
	cancelled := Modification{ModificationType: ModificationCancelledDeparture}

	tables := []struct {
		modifications []Modification
		remarks       []string
	}{
		{[]Modification{}, []string{}},
		{[]Modification{statusChange}, []string{}},
		{[]Modification{statusChange, delayNoCause}, []string{}},
		{[]Modification{statusChange, delayWithCause}, []string{"Later vertrek door testen"}},
		{[]Modification{cancelled, delayWithCause}, []string{"Trein rijdt niet", "Later vertrek door testen"}},
	}

	for _, table := range tables {
		remarks := GetRemarks(table.modifications, "nl")

		if len(remarks) != len(table.remarks) {
			t.Error("Number of expected remarks different from actual remarks returned")
		}

		for index, remark := range remarks {
			if remark != table.remarks[index] {
				t.Errorf("Expected remark %s does not match with actual value %s", table.remarks[index], remark)
			}
		}
	}
}
