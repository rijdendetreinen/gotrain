package parsers

import (
	"testing"

	"github.com/rijdendetreinen/gotrain/models"
)

func TestDvs2Recognition(t *testing.T) {
	departure, err := ParseDvsMessage(testFileReader(t, "dvs2/departure.xml"))

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	// Verify DVS version
	if departure.DvsVersion != models.DvsVersion2 {
		t.Errorf("Wrong DVS version: expected %d, but got %d", models.DvsVersion2, departure.DvsVersion)
	}
}

func TestDvs3Recognition(t *testing.T) {
	departure, err := ParseDvsMessage(testFileReader(t, "dvs3/example.xml"))

	if err != nil && err.Error() != "not implemented" {
		t.Fatalf("Parser error: %v", err)
	}

	// Verify DVS version
	if departure.DvsVersion != models.DvsVersion3 {
		t.Errorf("Wrong DVS version: expected %d, but got %d", models.DvsVersion3, departure.DvsVersion)
	}
}

func TestHandleInvalidXml(t *testing.T) {
	_, err := ParseDvsMessage(testFileReader(t, "invalid.xml"))

	if err == nil {
		t.Error("Should return an error for invalid XML")
	}
}

func TestHandleInvalidDvs(t *testing.T) {
	departure, err := ParseDvsMessage(testFileReader(t, "arrival.xml"))

	if departure.DvsVersion != models.DvsVersionUnknown {
		t.Errorf("Wrong DVS version: expected %d, but got %d", models.DvsVersionUnknown, departure.DvsVersion)
	}

	if err == nil {
		t.Error("Should return an error for an invalid DVS message")
	}
}
