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
