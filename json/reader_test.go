package json

import (
	"os"
	"testing"

	"github.com/spdx/tools-golang/spdx/v3/v3_0"
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

// TestReadV3 tests that the SPDX Reader can read json documents into the SPDX 3 model,
// including older 3.0.x documents, without altering their content
func TestReadV3(t *testing.T) {
	tt := []struct {
		filename        string
		wantID          string
		wantElements    int
		wantSpecVersion string
	}{
		{
			filename: "test_fixtures/spdx2_3.json",
		},
		{
			filename:        "../spdx/v3/v3_0/testdata/test.json",
			wantID:          "http://spdx.example.com/Document1",
			wantElements:    5,
			wantSpecVersion: "3.0.1",
		},
		{
			filename:        "test_fixtures/spdx3_0_0.json",
			wantID:          "http://spdx.example.com/Document1",
			wantElements:    5,
			wantSpecVersion: "3.0.0",
		},
		{
			filename:        "test_fixtures/spdx3_0_0_graph_first.json",
			wantID:          "http://spdx.example.com/Document1",
			wantElements:    5,
			wantSpecVersion: "3.0.0",
		},
	}

	for _, tc := range tt {
		t.Run(tc.filename, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			if err != nil {
				t.Fatalf("error opening %s: %v", tc.filename, err)
			}
			defer file.Close()
			doc := v3_0.Document{}
			err = ReadInto(file, &doc)
			if err != nil {
				t.Fatalf("error reading %s: %v", tc.filename, err)
			}
			if tc.wantID != "" && doc.ID != tc.wantID {
				t.Errorf("document ID: got %q, want %q", doc.ID, tc.wantID)
			}
			if tc.wantElements != 0 && len(doc.Elements) != tc.wantElements {
				t.Errorf("elements: got %d, want %d", len(doc.Elements), tc.wantElements)
			}
			if tc.wantSpecVersion != "" && doc.CreationInfo.GetSpecVersion() != tc.wantSpecVersion {
				t.Errorf("specVersion: got %q, want %q", doc.CreationInfo.GetSpecVersion(), tc.wantSpecVersion)
			}
		})
	}
}
