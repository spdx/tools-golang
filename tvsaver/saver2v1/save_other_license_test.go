// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Other License section Saver tests =====
func TestSaver2_1OtherLicenseSavesText(t *testing.T) {
	ol := &spdx.OtherLicense2_1{
		Identifier: "LicenseRef-1",
		ExtractedText: `License 1 text
blah blah blah
blah blah blah blah`,
		Name: "License 1",
		CrossReferences: []string{
			"http://example.com/License1/",
			"http://example.com/License1AnotherURL/",
		},
		Comment: "this is a license comment",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`LicenseID: LicenseRef-1
ExtractedText: <text>License 1 text
blah blah blah
blah blah blah blah</text>
LicenseName: License 1
LicenseCrossReference: http://example.com/License1/
LicenseCrossReference: http://example.com/License1AnotherURL/
LicenseComment: this is a license comment

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderOtherLicense2_1(ol, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1OtherLicenseOmitsOptionalFieldsIfEmpty(t *testing.T) {
	ol := &spdx.OtherLicense2_1{
		Identifier: "LicenseRef-1",
		ExtractedText: `License 1 text
blah blah blah
blah blah blah blah`,
		Name: "License 1",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`LicenseID: LicenseRef-1
ExtractedText: <text>License 1 text
blah blah blah
blah blah blah blah</text>
LicenseName: License 1

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderOtherLicense2_1(ol, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
