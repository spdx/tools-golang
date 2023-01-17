package rdf

import (
	"io"
	"strings"
	"testing"
)

func Test_Read(t *testing.T) {
	var reader io.Reader
	var err error

	// TestCase 1: invalid rdf/xml must raise an error
	reader = strings.NewReader("")
	_, err = Read(reader)
	if err == nil {
		t.Errorf("expected an EOF error reading an empty file, got %v", err)
	}

	// TestCase 2: Valid rdf/xml but invalid spdx document must raise an error
	reader = strings.NewReader(`
		<rdf:RDF
			xmlns:spdx="http://spdx.org/rdf/terms#"
			xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
			xmlns="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#">
		</rdf:RDF>
	`)
	_, err = Read(reader)
	if err == nil {
		t.Errorf("expected an error due to no SpdxDocument Node in the document")
	}
}
