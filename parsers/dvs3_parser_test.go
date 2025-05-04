package parsers

import (
	"testing"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestParseNormalDeparture_Dvs3(t *testing.T) {
	_, err := ParseDvs3Message(testFileReader(t, "dvs3/example.xml"))

	if err == nil {
		t.Error("Should return: Not implemented yet")
	}
}

func testParseDeparture_Dvs3(t *testing.T, name string) models.Departure {
	departure, err := ParseDvs3Message(testFileReader(t, name))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	return departure
}
