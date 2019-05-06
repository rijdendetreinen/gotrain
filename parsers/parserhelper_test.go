package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func testFileReader(t *testing.T, name string) *os.File {
	path := filepath.Join("testdata", name)

	file, err := os.Open(path)

	if err != nil {
		t.Fatal(err)
	}

	return file
}
