// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== File section Saver tests =====
func TestSaver2_1FileSavesText(t *testing.T) {
	f := &spdx.File2_1{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: spdx.ElementID("File123"),
		FileTypes: []string{
			"TEXT",
			"DOCUMENTATION",
		},
		Checksums: []spdx.Checksum{
			{Algorithm: spdx.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
			{Algorithm: spdx.SHA256, Value: "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"},
			{Algorithm: spdx.MD5, Value: "624c1abb3664f4b35547e7c73864ad24"},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
			"Apache-1.1",
		},
		LicenseComments:   "this is a license comment(s)",
		FileCopyrightText: "Copyright (c) Jane Doe",
		ArtifactOfProjects: []*spdx.ArtifactOfProject2_1{
			&spdx.ArtifactOfProject2_1{
				Name:     "project1",
				HomePage: "http://example.com/1/",
				URI:      "http://example.com/1/uri.whatever",
			},
			&spdx.ArtifactOfProject2_1{
				Name: "project2",
			},
			&spdx.ArtifactOfProject2_1{
				Name:     "project3",
				HomePage: "http://example.com/3/",
			},
			&spdx.ArtifactOfProject2_1{
				Name: "project4",
				URI:  "http://example.com/4/uri.whatever",
			},
		},
		FileComment: "this is a file comment",
		FileNotice:  "This file may be used under either Apache-2.0 or Apache-1.1.",
		FileContributors: []string{
			"John Doe jdoe@example.com",
			"EvilCorp",
		},
		FileDependencies: []string{
			"f-1.txt",
			"g.txt",
		},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`FileName: /tmp/whatever.txt
SPDXID: SPDXRef-File123
FileTypes: TEXT
FileTypes: DOCUMENTATION
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
FileChecksum: SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd
FileChecksum: MD5: 624c1abb3664f4b35547e7c73864ad24
LicenseConcluded: Apache-2.0
LicenseInfoInFiles: Apache-2.0
LicenseInfoInFiles: Apache-1.1
LicenseComments: this is a license comment(s)
FileCopyrightText: Copyright (c) Jane Doe
ArtifactOfProjectName: project1
ArtifactOfProjectHomePage: http://example.com/1/
ArtifactOfProjectURI: http://example.com/1/uri.whatever
ArtifactOfProjectName: project2
ArtifactOfProjectName: project3
ArtifactOfProjectHomePage: http://example.com/3/
ArtifactOfProjectName: project4
ArtifactOfProjectURI: http://example.com/4/uri.whatever
FileComment: this is a file comment
FileNotice: This file may be used under either Apache-2.0 or Apache-1.1.
FileContributors: John Doe jdoe@example.com
FileContributors: EvilCorp
FileDependency: f-1.txt
FileDependency: g.txt

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_1(f, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1FileSavesSnippetsAlso(t *testing.T) {
	sn1 := &spdx.Snippet2_1{
		SnippetSPDXIdentifier:         spdx.ElementID("Snippet19"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File123").ElementRefID,
		Ranges:                        []spdx.SnippetRange{{StartPointer: spdx.SnippetRangePointer{Offset: 17}, EndPointer: spdx.SnippetRangePointer{Offset: 209}}},
		SnippetLicenseConcluded:       "GPL-2.0-or-later",
		SnippetCopyrightText:          "Copyright (c) John Doe 20x6",
	}

	sn2 := &spdx.Snippet2_1{
		SnippetSPDXIdentifier:         spdx.ElementID("Snippet20"),
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "File123").ElementRefID,
		Ranges:                        []spdx.SnippetRange{{StartPointer: spdx.SnippetRangePointer{Offset: 268}, EndPointer: spdx.SnippetRangePointer{Offset: 309}}},
		SnippetLicenseConcluded:       "WTFPL",
		SnippetCopyrightText:          "NOASSERTION",
	}

	sns := map[spdx.ElementID]*spdx.Snippet2_1{
		spdx.ElementID("Snippet19"): sn1,
		spdx.ElementID("Snippet20"): sn2,
	}

	f := &spdx.File2_1{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: spdx.ElementID("File123"),
		Checksums: []spdx.Checksum{
			{Algorithm: spdx.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
		},
		FileCopyrightText: "Copyright (c) Jane Doe",
		Snippets:          sns,
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`FileName: /tmp/whatever.txt
SPDXID: SPDXRef-File123
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
LicenseConcluded: Apache-2.0
LicenseInfoInFiles: Apache-2.0
FileCopyrightText: Copyright (c) Jane Doe

SnippetSPDXID: SPDXRef-Snippet19
SnippetFromFileSPDXID: SPDXRef-File123
SnippetByteRange: 17:209
SnippetLicenseConcluded: GPL-2.0-or-later
SnippetCopyrightText: Copyright (c) John Doe 20x6

SnippetSPDXID: SPDXRef-Snippet20
SnippetFromFileSPDXID: SPDXRef-File123
SnippetByteRange: 268:309
SnippetLicenseConcluded: WTFPL
SnippetCopyrightText: NOASSERTION

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_1(f, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1FileOmitsOptionalFieldsIfEmpty(t *testing.T) {
	f := &spdx.File2_1{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: spdx.ElementID("File123"),
		Checksums: []spdx.Checksum{
			{Algorithm: spdx.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
		},
		FileCopyrightText: "Copyright (c) Jane Doe",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`FileName: /tmp/whatever.txt
SPDXID: SPDXRef-File123
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
LicenseConcluded: Apache-2.0
LicenseInfoInFiles: Apache-2.0
FileCopyrightText: Copyright (c) Jane Doe

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_1(f, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1FileWrapsCopyrightMultiLine(t *testing.T) {
	f := &spdx.File2_1{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: spdx.ElementID("File123"),
		Checksums: []spdx.Checksum{
			{Algorithm: spdx.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
		},
		FileCopyrightText: `Copyright (c) Jane Doe
Copyright (c) John Doe`,
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`FileName: /tmp/whatever.txt
SPDXID: SPDXRef-File123
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
LicenseConcluded: Apache-2.0
LicenseInfoInFiles: Apache-2.0
FileCopyrightText: <text>Copyright (c) Jane Doe
Copyright (c) John Doe</text>

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_1(f, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1FileWrapsCommentsAndNoticesMultiLine(t *testing.T) {
	f := &spdx.File2_1{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: spdx.ElementID("File123"),
		Checksums: []spdx.Checksum{
			{Algorithm: spdx.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
		},
		LicenseComments: `this is a
multi-line license comment`,
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
		},
		FileCopyrightText: "Copyright (c) Jane Doe",
		FileComment: `this is a
multi-line file comment`,
		FileNotice: `This file may be used
under either Apache-2.0 or Apache-1.1.`,
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`FileName: /tmp/whatever.txt
SPDXID: SPDXRef-File123
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
LicenseConcluded: Apache-2.0
LicenseInfoInFiles: Apache-2.0
LicenseComments: <text>this is a
multi-line license comment</text>
FileCopyrightText: Copyright (c) Jane Doe
FileComment: <text>this is a
multi-line file comment</text>
FileNotice: <text>This file may be used
under either Apache-2.0 or Apache-1.1.</text>

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_1(f, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
