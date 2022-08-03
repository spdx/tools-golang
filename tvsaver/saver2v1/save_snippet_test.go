// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
)

// ===== Snippet section Saver tests =====
func TestSaver2_1SnippetSavesText(t *testing.T) {
	sn := &v2_1.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "File292").ElementRefID,
		Ranges: []common.SnippetRange{
			{
				StartPointer: common.SnippetRangePointer{LineNumber: 3},
				EndPointer:   common.SnippetRangePointer{LineNumber: 8},
			},
			{
				StartPointer: common.SnippetRangePointer{Offset: 17},
				EndPointer:   common.SnippetRangePointer{Offset: 209},
			},
		},
		SnippetLicenseConcluded: "GPL-2.0-or-later",
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
SnippetLineRange: 3:8
SnippetByteRange: 17:209
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
	sn := &v2_1.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "File292").ElementRefID,
		Ranges: []common.SnippetRange{
			{
				StartPointer: common.SnippetRangePointer{Offset: 17},
				EndPointer:   common.SnippetRangePointer{Offset: 209},
			},
		},
		SnippetLicenseConcluded: "GPL-2.0-or-later",
		SnippetCopyrightText:    "Copyright (c) John Doe 20x6",
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
	sn := &v2_1.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("Snippet17"),
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "File292").ElementRefID,
		Ranges: []common.SnippetRange{
			{
				StartPointer: common.SnippetRangePointer{Offset: 17},
				EndPointer:   common.SnippetRangePointer{Offset: 209},
			},
		},
		SnippetLicenseConcluded: "GPL-2.0-or-later",
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
