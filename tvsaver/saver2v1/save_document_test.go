// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== entire Document Saver tests =====
func TestSaver2_1DocumentSavesText(t *testing.T) {

	// Creation Info section
	ci := &spdx.CreationInfo2_1{
		SPDXVersion:       "SPDX-2.1",
		DataLicense:       "CC0-1.0",
		SPDXIdentifier:    spdx.ElementID("DOCUMENT"),
		DocumentName:      "spdx-go-0.0.1.abcdef",
		DocumentNamespace: "https://github.com/swinslow/spdx-docs/spdx-go/spdx-go-0.0.1.abcdef.whatever",
		CreatorPersons: []string{
			"John Doe",
		},
		Created: "2018-10-10T06:20:00Z",
	}

	// unpackaged files
	f1 := &spdx.File2_1{
		Name:              "/tmp/whatever1.txt",
		SPDXIdentifier:    spdx.ElementID("File1231"),
		ChecksumSHA1:      "85ed0817af83a24ad8da68c2b5094de69833983c",
		LicenseConcluded:  "Apache-2.0",
		LicenseInfoInFile: []string{"Apache-2.0"},
		CopyrightText:     "Copyright (c) Jane Doe",
	}

	f2 := &spdx.File2_1{
		Name:              "/tmp/whatever2.txt",
		SPDXIdentifier:    spdx.ElementID("File1232"),
		ChecksumSHA1:      "85ed0817af83a24ad8da68c2b5094de69833983d",
		LicenseConcluded:  "MIT",
		LicenseInfoInFile: []string{"MIT"},
		CopyrightText:     "Copyright (c) John Doe",
	}

	unFiles := map[spdx.ElementID]*spdx.File2_1{
		spdx.ElementID("File1231"): f1,
		spdx.ElementID("File1232"): f2,
	}

	// Package 1: packaged files with snippets
	sn1 := &spdx.Snippet2_1{
		SPDXIdentifier:                "Snippet19",
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "FileHasSnippets"),
		ByteRangeStart:                17,
		ByteRangeEnd:                  209,
		LicenseConcluded:              "GPL-2.0-or-later",
		CopyrightText:                 "Copyright (c) John Doe 20x6",
	}

	sn2 := &spdx.Snippet2_1{
		SPDXIdentifier:                "Snippet20",
		SnippetFromFileSPDXIdentifier: spdx.MakeDocElementID("", "FileHasSnippets"),
		ByteRangeStart:                268,
		ByteRangeEnd:                  309,
		LicenseConcluded:              "WTFPL",
		CopyrightText:                 "NOASSERTION",
	}

	f3 := &spdx.File2_1{
		Name:             "/tmp/file-with-snippets.txt",
		SPDXIdentifier:   spdx.ElementID("FileHasSnippets"),
		ChecksumSHA1:     "85ed0817af83a24ad8da68c2b5094de69833983e",
		LicenseConcluded: "GPL-2.0-or-later AND WTFPL",
		LicenseInfoInFile: []string{
			"Apache-2.0",
			"GPL-2.0-or-later",
			"WTFPL",
		},
		CopyrightText: "Copyright (c) Jane Doe",
		Snippets: map[spdx.ElementID]*spdx.Snippet2_1{
			spdx.ElementID("Snippet19"): sn1,
			spdx.ElementID("Snippet20"): sn2,
		},
	}

	f4 := &spdx.File2_1{
		Name:              "/tmp/another-file.txt",
		SPDXIdentifier:    spdx.ElementID("FileAnother"),
		ChecksumSHA1:      "85ed0817af83a24ad8da68c2b5094de69833983f",
		LicenseConcluded:  "BSD-3-Clause",
		LicenseInfoInFile: []string{"BSD-3-Clause"},
		CopyrightText:     "Copyright (c) Jane Doe LLC",
	}

	pkgWith := &spdx.Package2_1{
		Name:                      "p1",
		SPDXIdentifier:            spdx.ElementID("p1"),
		DownloadLocation:          "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:             true,
		IsFilesAnalyzedTagPresent: true,
		VerificationCode:          "0123456789abcdef0123456789abcdef01234567",
		LicenseConcluded:          "GPL-2.0-or-later AND BSD-3-Clause AND WTFPL",
		LicenseInfoFromFiles: []string{
			"Apache-2.0",
			"GPL-2.0-or-later",
			"WTFPL",
			"BSD-3-Clause",
		},
		LicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		CopyrightText:   "Copyright (c) John Doe, Inc.",
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID("FileHasSnippets"): f3,
			spdx.ElementID("FileAnother"):     f4,
		},
	}

	// Other Licenses 1 and 2
	ol1 := &spdx.OtherLicense2_1{
		Identifier: "LicenseRef-1",
		ExtractedText: `License 1 text
blah blah blah
blah blah blah blah`,
		Name: "License 1",
	}

	ol2 := &spdx.OtherLicense2_1{
		Identifier:    "LicenseRef-2",
		ExtractedText: `License 2 text - this is a license that does some stuff`,
		Name:          "License 2",
	}

	// Relationships
	rln1 := &spdx.Relationship2_1{
		RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
		RefB:         spdx.MakeDocElementID("", "p1"),
		Relationship: "DESCRIBES",
	}

	rln2 := &spdx.Relationship2_1{
		RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
		RefB:         spdx.MakeDocElementID("", "File1231"),
		Relationship: "DESCRIBES",
	}

	rln3 := &spdx.Relationship2_1{
		RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
		RefB:         spdx.MakeDocElementID("", "File1232"),
		Relationship: "DESCRIBES",
	}

	// Annotations
	ann1 := &spdx.Annotation2_1{
		Annotator:      "John Doe",
		AnnotatorType:  "Person",
		Date:           "2018-10-10T17:52:00Z",
		Type:           "REVIEW",
		SPDXIdentifier: spdx.MakeDocElementID("", "DOCUMENT"),
		Comment:        "This is an annotation about the SPDX document",
	}

	ann2 := &spdx.Annotation2_1{
		Annotator:      "John Doe, Inc.",
		AnnotatorType:  "Organization",
		Date:           "2018-10-10T17:52:00Z",
		Type:           "REVIEW",
		SPDXIdentifier: spdx.MakeDocElementID("", "p1"),
		Comment:        "This is an annotation about Package p1",
	}

	// Reviews
	rev1 := &spdx.Review2_1{
		Reviewer:     "John Doe",
		ReviewerType: "Person",
		Date:         "2018-10-14T10:28:00Z",
	}
	rev2 := &spdx.Review2_1{
		Reviewer:     "Jane Doe LLC",
		ReviewerType: "Organization",
		Date:         "2018-10-14T10:28:00Z",
		Comment:      "I have reviewed this SPDX document and it is awesome",
	}

	// now, build the document
	doc := &spdx.Document2_1{
		CreationInfo: ci,
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): pkgWith,
		},
		UnpackagedFiles: unFiles,
		OtherLicenses: []*spdx.OtherLicense2_1{
			ol1,
			ol2,
		},
		Relationships: []*spdx.Relationship2_1{
			rln1,
			rln2,
			rln3,
		},
		Annotations: []*spdx.Annotation2_1{
			ann1,
			ann2,
		},
		Reviews: []*spdx.Review2_1{
			rev1,
			rev2,
		},
	}

	want := bytes.NewBufferString(`SPDXVersion: SPDX-2.1
DataLicense: CC0-1.0
SPDXID: SPDXRef-DOCUMENT
DocumentName: spdx-go-0.0.1.abcdef
DocumentNamespace: https://github.com/swinslow/spdx-docs/spdx-go/spdx-go-0.0.1.abcdef.whatever
Creator: Person: John Doe
Created: 2018-10-10T06:20:00Z

##### Unpackaged files

FileName: /tmp/whatever1.txt
SPDXID: SPDXRef-File1231
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
LicenseConcluded: Apache-2.0
LicenseInfoInFile: Apache-2.0
FileCopyrightText: Copyright (c) Jane Doe

FileName: /tmp/whatever2.txt
SPDXID: SPDXRef-File1232
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983d
LicenseConcluded: MIT
LicenseInfoInFile: MIT
FileCopyrightText: Copyright (c) John Doe

##### Package: p1

PackageName: p1
SPDXID: SPDXRef-p1
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: true
PackageVerificationCode: 0123456789abcdef0123456789abcdef01234567
PackageLicenseConcluded: GPL-2.0-or-later AND BSD-3-Clause AND WTFPL
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseInfoFromFiles: WTFPL
PackageLicenseInfoFromFiles: BSD-3-Clause
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageCopyrightText: Copyright (c) John Doe, Inc.

FileName: /tmp/another-file.txt
SPDXID: SPDXRef-FileAnother
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983f
LicenseConcluded: BSD-3-Clause
LicenseInfoInFile: BSD-3-Clause
FileCopyrightText: Copyright (c) Jane Doe LLC

FileName: /tmp/file-with-snippets.txt
SPDXID: SPDXRef-FileHasSnippets
FileChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983e
LicenseConcluded: GPL-2.0-or-later AND WTFPL
LicenseInfoInFile: Apache-2.0
LicenseInfoInFile: GPL-2.0-or-later
LicenseInfoInFile: WTFPL
FileCopyrightText: Copyright (c) Jane Doe

SnippetSPDXIdentifier: SPDXRef-Snippet19
SnippetFromFileSPDXID: SPDXRef-FileHasSnippets
SnippetByteRange: 17:209
SnippetLicenseConcluded: GPL-2.0-or-later
SnippetCopyrightText: Copyright (c) John Doe 20x6

SnippetSPDXIdentifier: SPDXRef-Snippet20
SnippetFromFileSPDXID: SPDXRef-FileHasSnippets
SnippetByteRange: 268:309
SnippetLicenseConcluded: WTFPL
SnippetCopyrightText: NOASSERTION

##### Other Licenses

LicenseID: LicenseRef-1
ExtractedText: <text>License 1 text
blah blah blah
blah blah blah blah</text>
LicenseName: License 1

LicenseID: LicenseRef-2
ExtractedText: License 2 text - this is a license that does some stuff
LicenseName: License 2

##### Relationships

Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-p1
Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-File1231
Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-File1232

##### Annotations

Annotator: Person: John Doe
AnnotationDate: 2018-10-10T17:52:00Z
AnnotationType: REVIEW
SPDXREF: SPDXRef-DOCUMENT
AnnotationComment: This is an annotation about the SPDX document

Annotator: Organization: John Doe, Inc.
AnnotationDate: 2018-10-10T17:52:00Z
AnnotationType: REVIEW
SPDXREF: SPDXRef-p1
AnnotationComment: This is an annotation about Package p1

##### Reviews

Reviewer: Person: John Doe
ReviewDate: 2018-10-14T10:28:00Z

Reviewer: Organization: Jane Doe LLC
ReviewDate: 2018-10-14T10:28:00Z
ReviewComment: I have reviewed this SPDX document and it is awesome

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := RenderDocument2_1(doc, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected {{{%v}}}, got {{{%v}}}", want.String(), got.String())
	}

}

func TestSaver2_1DocumentReturnsErrorIfNilCreationInfo(t *testing.T) {
	doc := &spdx.Document2_1{}

	var got bytes.Buffer
	err := RenderDocument2_1(doc, &got)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
