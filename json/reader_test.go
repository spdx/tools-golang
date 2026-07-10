package json

import (
	"os"
	"testing"
)

// TestRead tests that the SPDX Reader can still parse json documents correctly
// this protects against any of the custom unmarshalling code breaking given a new change set
func TestRead(t *testing.T) {
	tt := []struct {
		filename string
	}{
		{"test_fixtures/spdx2_3.json"},
	}

	for _, tc := range tt {
		t.Run(tc.filename, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			if err != nil {
				t.Errorf("error opening %s: %v", tc.filename, err)
			}
			defer file.Close()
			_, err = Read(file)
			if err != nil {
				t.Errorf("error reading %s: %v", tc.filename, err)
			}
		})
	}
}

// TestReadNullPackage makes sure a document with a null entry in the packages
// array doesn't panic. The custom unmarshaller walks d.Packages to build
// hasFiles relationships, and a null element decodes to a nil *Package that
// used to be dereferenced. This mirrors the null-relationships handling added
// for #238.
func TestReadNullPackage(t *testing.T) {
	fixtures := []string{
		"test_fixtures/spdx2_2_null_package.json",
		"test_fixtures/spdx2_3_null_package.json",
	}

	for _, filename := range fixtures {
		t.Run(filename, func(t *testing.T) {
			file, err := os.Open(filename)
			if err != nil {
				t.Fatalf("error opening %s: %v", filename, err)
			}
			defer file.Close()

			// Read must return cleanly here; before the nil guard it panicked.
			if _, err = Read(file); err != nil {
				t.Errorf("error reading %s: %v", filename, err)
			}
		})
	}
}
