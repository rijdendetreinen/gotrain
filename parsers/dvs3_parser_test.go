package parsers

import (
	"testing"
	"time"

	"github.com/rijdendetreinen/gotrain/models"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, 189, departure.Delay, "Wrong Delay: expected %d, but got %d", 189, departure.Delay)

	if departure.Delay != 189 {
		t.Errorf("Wrong Delay: expected %d, but got %d", 189, departure.Delay)
	}

	// Verify that the destination is parsed correctly
	assert.Len(t, departure.DestinationActual, 1, "Wrong number of actual destinations")
	// assert.Len(t, departure.DestinationPlanned, 1, "Wrong number of planned destinations")

	if len(departure.DestinationActual) > 0 {
		// Verify that the destination station is parsed correctly
		assert.Equal(t, "MTR", departure.DestinationActual[0].Code, "Wrong DestinationActual Code")
		assert.Equal(t, "Randwyck", departure.DestinationActual[0].NameShort, "Wrong DestinationActual NameShort")
		assert.Equal(t, "Randwyck", departure.DestinationActual[0].NameMedium, "Wrong DestinationActual NameMedium")
		assert.Equal(t, "Maastricht Randwyck", departure.DestinationActual[0].NameLong, "Wrong DestinationActual NameLong")
	}

	// if len(departure.DestinationPlanned) > 0 {
	// 	// Verify that the destination station is parsed correctly
	// 	assert.Equal(t, "MTR", departure.DestinationPlanned[0].Code, "Wrong DestinationPlanned Code")
	// 	assert.Equal(t, "Randwyck", departure.DestinationPlanned[0].NameShort, "Wrong DestinationPlanned NameShort")
	// 	assert.Equal(t, "Randwyck", departure.DestinationPlanned[0].NameMedium, "Wrong DestinationPlanned NameMedium")
	// 	assert.Equal(t, "Maastricht Randwyck", departure.DestinationPlanned[0].NameLong, "Wrong DestinationPlanned NameLong")
	// }

	// Verify departure platform
	assert.Equal(t, "4b", departure.PlatformActual, "Wrong actual departure platform")
	assert.False(t, departure.PlatformChanged(), "Platform should not be changed")

	// Verify number of wings (1)
	assert.Len(t, departure.TrainWings, 1, "Wrong number of wings")
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
