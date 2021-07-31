// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Snippet section Saver tests =====
func TestSaver2_1SnippetSavesText(t *testing.T) {
	sn := &spdx.Snippet2_1{
		SnippetSPDXIdentifier:         spdx.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File292"),
		SnippetByteRangeStart:         17,
		SnippetByteRangeEnd:           209,
		SnippetLineRangeStart:         3,
		SnippetLineRangeEnd:           8,
		SnippetLicenseConcluded:       "GPL-2.0-or-later",
		LicenseInfoInSnippet: []string{
			"GPL-2.0-or-later",
			"MIT",
		},
		SnippetLicenseComments: "this is a comment(s) about the snippet license",
		SnippetCopyrightText:   "Copyright (c) John Doe 20x6",
		SnippetComment:         "this is a snippet comment",
		SnippetName:            "from John's program",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`SnippetSPDXID: SPDXRef-Snippet17
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

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderSnippet2_1(sn, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1SnippetOmitsOptionalFieldsIfEmpty(t *testing.T) {
	sn := &spdx.Snippet2_1{
		SnippetSPDXIdentifier:         spdx.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File292"),
		SnippetByteRangeStart:         17,
		SnippetByteRangeEnd:           209,
		SnippetLicenseConcluded:       "GPL-2.0-or-later",
		SnippetCopyrightText:          "Copyright (c) John Doe 20x6",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`SnippetSPDXID: SPDXRef-Snippet17
SnippetFromFileSPDXID: SPDXRef-File292
SnippetByteRange: 17:209
SnippetLicenseConcluded: GPL-2.0-or-later
SnippetCopyrightText: Copyright (c) John Doe 20x6

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderSnippet2_1(sn, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1SnippetWrapsCopyrightMultiline(t *testing.T) {
	sn := &spdx.Snippet2_1{
		SnippetSPDXIdentifier:         spdx.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File292"),
		SnippetByteRangeStart:         17,
		SnippetByteRangeEnd:           209,
		SnippetLicenseConcluded:       "GPL-2.0-or-later",
		SnippetCopyrightText: `Copyright (c) John Doe 20x6
Copyright (c) John Doe 20x6`,
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`SnippetSPDXID: SPDXRef-Snippet17
SnippetFromFileSPDXID: SPDXRef-File292
SnippetByteRange: 17:209
SnippetLicenseConcluded: GPL-2.0-or-later
SnippetCopyrightText: <text>Copyright (c) John Doe 20x6
Copyright (c) John Doe 20x6</text>

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderSnippet2_1(sn, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
