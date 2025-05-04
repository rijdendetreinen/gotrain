package parsers

import (
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
)

// testParseDeparture_Dvs3 is a helper function to parse a DVS3 message and check for errors
// It returns the parsed Departure object or fails the test if an error occurs
// It also checks if the DVS version is correct
func testParseDeparture_Dvs3(t *testing.T, name string) models.Departure {
	departure, err := ParseDvs3Message(testFileReader(t, name))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if departure.DvsVersion != models.DvsVersion2 {
		t.Errorf("Wrong DVS version: expected %d, but got %d", 2, departure.DvsVersion)
	}

	return departure
}

func TestParseNormalDeparture_Dvs3(t *testing.T) {
	departure, err := ParseDvs3Message(testFileReader(t, "dvs3/example.xml"))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	// Verify that the departure is parsed correctly
	if departure.ServiceID != "32437" {
		t.Errorf("Wrong ServiceID: expected %s, but got %s", "32437", departure.ServiceID)
	}

	if departure.ServiceDate != "2024-12-19" {
		t.Errorf("Wrong ServiceDate: expected %s, but got %s", "2024-12-19", departure.ServiceDate)
	}

	if departure.Station.Code != "MT" {
		t.Errorf("Wrong Station Code: expected %s, but got %s", "MT", departure.Station.Code)
	}
	if departure.Station.NameShort != "Maastricht" {
		t.Errorf("Wrong Station NameShort: expected %s, but got %s", "Maastricht", departure.Station.NameShort)
	}
	if departure.Station.NameMedium != "Maastricht" {
		t.Errorf("Wrong Station NameMedium: expected %s, but got %s", "Maastricht", departure.Station.NameMedium)
	}
	if departure.Station.NameLong != "Maastricht" {
		t.Errorf("Wrong Station NameLong: expected %s, but got %s", "Maastricht", departure.Station.NameLong)
	}

	// Verify that the service ID is generated correctly
	if departure.ID != "2024-12-19-32437-MT" {
		t.Errorf("Wrong ID: expected %s, but got %s", "2024-12-19-32437-MT", departure.ID)
	}

	// Service type and line number
	if departure.ServiceType != "Stoptrein" {
		t.Errorf("Wrong ServiceType: expected %s, but got %s", "Stoptrein", departure.ServiceType)
	}
	if departure.ServiceTypeCode != "ST" {
		t.Errorf("Wrong ServiceTypeCode: expected %s, but got %s", "ST", departure.ServiceTypeCode)
	}
	if departure.LineNumber != "RS12" {
		t.Errorf("Wrong LineNumber: expected %s, but got %s", "RS12", departure.LineNumber)
	}
	if departure.Company != "Arriva" {
		t.Errorf("Wrong Company: expected %s, but got %s", "Arriva", departure.Company)
	}

	// Departure status
	if departure.Status != models.DepartureStatusUnknown {
		t.Errorf("Wrong Status: expected %d, but got %d", models.DepartureStatusUnknown, departure.Status)
	}

	// Departure time, delay, destination(s), platform
	expectedDepartureTime, _ := time.Parse(time.RFC3339, "2024-12-19T11:52:00+01:00")

	if !departure.DepartureTime.Equal(expectedDepartureTime) {
		t.Errorf("Expected departure time %v does not match %v", expectedDepartureTime, departure.DepartureTime)
	}

	if departure.Delay != 189 {
		t.Errorf("Wrong Delay: expected %d, but got %d", 189, departure.Delay)
	}
}

func TestHandleWrongNamespace_Dvs3(t *testing.T) {
	departure, err := ParseDvs3Message(testFileReader(t, "dvs3/wrong_namespace.xml"))

	if departure.DvsVersion != models.DvsVersionUnknown {
		t.Errorf("Wrong DVS version: expected %d, but got %d", models.DvsVersionUnknown, departure.DvsVersion)
	}

	if err == nil {
		t.Error("Should return an error for an invalid DVS message")
	}
}
func TestInvalidDeparture_Dvs3(t *testing.T) {
	_, err := ParseDvs3Message(testFileReader(t, "invalid.xml"))

	if err == nil {
		t.Error("Should return an error for invalid XML")
	}

	_, err = ParseDvs3Message(testFileReader(t, "arrival.xml"))

	if err == nil {
		t.Error("Should return an error for an Arrival message")
	}
}
