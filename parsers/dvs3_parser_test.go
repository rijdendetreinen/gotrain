package parsers

import (
	"testing"

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
	_, err := ParseDvs3Message(testFileReader(t, "dvs3/example.xml"))

	if err == nil {
		t.Error("Should return: Not implemented yet")
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
