// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// ===== File section Saver tests =====
func TestSaver2_2FileSavesText(t *testing.T) {
	f := &v2_2.File{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: common.ElementID("File123"),
		FileTypes: []string{
			"TEXT",
			"DOCUMENTATION",
		},
		Checksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
			{Algorithm: common.SHA256, Value: "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"},
			{Algorithm: common.MD5, Value: "624c1abb3664f4b35547e7c73864ad24"},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
			"Apache-1.1",
		},
		LicenseComments:   "this is a license comment(s)",
		FileCopyrightText: "Copyright (c) Jane Doe",
		ArtifactOfProjects: []*v2_2.ArtifactOfProject{
			&v2_2.ArtifactOfProject{
				Name:     "project1",
				HomePage: "http://example.com/1/",
				URI:      "http://example.com/1/uri.whatever",
			},
			&v2_2.ArtifactOfProject{
				Name: "project2",
			},
			&v2_2.ArtifactOfProject{
				Name:     "project3",
				HomePage: "http://example.com/3/",
			},
			&v2_2.ArtifactOfProject{
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
		FileAttributionTexts: []string{
			"attributions",
			`multi-line
attribution`,
		},
		FileDependencies: []string{
			"f-1.txt",
			"g.txt",
		},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`FileName: /tmp/whatever.txt
				SPDXID: SPDXRef-File123
				FileType: TEXT
				FileType: DOCUMENTATION
				FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
				FileChecksum: SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd
				FileChecksum: MD5: 624c1abb3664f4b35547e7c73864ad24
				LicenseConcluded: Apache-2.0
				LicenseInfoInFile: Apache-2.0
				LicenseInfoInFile: Apache-1.1
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
				FileContributor: John Doe jdoe@example.com
				FileContributor: EvilCorp
				FileAttributionText: attributions
				FileAttributionText: <text>multi-line
				attribution</text>
				FileDependency: f-1.txt
FileDependency: g.txt
				
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_2(f, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_2FileSavesSnippetsAlso(t *testing.T) {
	sn1 := &v2_2.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("Snippet19"),
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "File123").ElementRefID,
		Ranges:                        []common.SnippetRange{{StartPointer: common.SnippetRangePointer{Offset: 17}, EndPointer: common.SnippetRangePointer{Offset: 209}}},
		SnippetLicenseConcluded:       "GPL-2.0-or-later",
		SnippetCopyrightText:          "Copyright (c) John Doe 20x6",
	}

	sn2 := &v2_2.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("Snippet20"),
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "File123").ElementRefID,
		Ranges:                        []common.SnippetRange{{StartPointer: common.SnippetRangePointer{Offset: 268}, EndPointer: common.SnippetRangePointer{Offset: 309}}},
		SnippetLicenseConcluded:       "WTFPL",
		SnippetCopyrightText:          "NOASSERTION",
	}

	sns := map[common.ElementID]*v2_2.Snippet{
		common.ElementID("Snippet19"): sn1,
		common.ElementID("Snippet20"): sn2,
	}

	f := &v2_2.File{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: common.ElementID("File123"),
		Checksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
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
LicenseInfoInFile: Apache-2.0
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
	err := renderFile2_2(f, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_2FileOmitsOptionalFieldsIfEmpty(t *testing.T) {
	f := &v2_2.File{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: common.ElementID("File123"),
		Checksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
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
LicenseInfoInFile: Apache-2.0
FileCopyrightText: Copyright (c) Jane Doe

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_2(f, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_2FileWrapsCopyrightMultiLine(t *testing.T) {
	f := &v2_2.File{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: common.ElementID("File123"),
		Checksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
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
LicenseInfoInFile: Apache-2.0
FileCopyrightText: <text>Copyright (c) Jane Doe
Copyright (c) John Doe</text>

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderFile2_2(f, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_2FileWrapsCommentsAndNoticesMultiLine(t *testing.T) {
	f := &v2_2.File{
		FileName:           "/tmp/whatever.txt",
		FileSPDXIdentifier: common.ElementID("File123"),
		Checksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "85ed0817af83a24ad8da68c2b5094de69833983c"},
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
LicenseInfoInFile: Apache-2.0
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
	err := renderFile2_2(f, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}
