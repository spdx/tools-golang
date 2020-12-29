// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Snippet section Saver tests =====
func TestSaver2_2SnippetSavesText(t *testing.T) {
	sn := &spdx.Snippet2_2{
		SPDXIdentifier:                spdx.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File292"),
		ByteRangeStart:                17,
		ByteRangeEnd:                  209,
		LineRangeStart:                3,
		LineRangeEnd:                  8,
		LicenseConcluded:              "GPL-2.0-or-later",
		LicenseInfoInSnippet: []string{
			"GPL-2.0-or-later",
			"MIT",
		},
		LicenseComments:  "this is a comment(s) about the snippet license",
		CopyrightText:    "Copyright (c) John Doe 20x6",
		Comment:          "this is a snippet comment",
		Name:             "from John's program",
		AttributionTexts: []string{"some attributions"},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`SnippetSPDXIdentifier: SPDXRef-Snippet17
SnippetFromFileSPDXID: SPDXRef-File292
SnippetByteRange: 17:209
SnippetLineRange: 3:8
SnippetLicenseConcluded: GPL-2.0-or-later
LicenseInfoInSnippet: GPL-2.0-or-later
LicenseInfoInSnippet: MIT
SnippetLicenseComments: this is a comment(s) about the snippet license
SnippetCopyrightText: Copyright (c) John Doe 20x6
SnippetComment: this is a snippet comment
SnippetName: from John's program
SnippetAttributionText: some attributions

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderSnippet2_2(sn, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_2SnippetOmitsOptionalFieldsIfEmpty(t *testing.T) {
	sn := &spdx.Snippet2_2{
		SPDXIdentifier:                spdx.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File292"),
		ByteRangeStart:                17,
		ByteRangeEnd:                  209,
		LicenseConcluded:              "GPL-2.0-or-later",
		CopyrightText:                 "Copyright (c) John Doe 20x6",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`SnippetSPDXIdentifier: SPDXRef-Snippet17
SnippetFromFileSPDXID: SPDXRef-File292
SnippetByteRange: 17:209
SnippetLicenseConcluded: GPL-2.0-or-later
SnippetCopyrightText: Copyright (c) John Doe 20x6

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderSnippet2_2(sn, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
