// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package writer

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// ===== entire Document Saver tests =====
func TestSaverDocumentSavesText(t *testing.T) {

	// Creation Info section
	ci := &spdx.CreationInfo{
		Creators: []common.Creator{
			{Creator: "John Doe", CreatorType: "Person"},
		},
		Created: "2018-10-10T06:20:00Z",
	}

	// unpackaged files
	f1 := &spdx.File{
		FileName:           "/tmp/whatever1.txt",
		FileSPDXIdentifier: common.ElementID("File1231"),
		Checksums:          []common.Checksum{{Value: "85ed0817af83a24ad8da68c2b5094de69833983c", Algorithm: common.SHA1}},
		LicenseConcluded:   "Apache-2.0",
		LicenseInfoInFiles: []string{"Apache-2.0"},
		FileCopyrightText:  "Copyright (c) Jane Doe",
	}

	f2 := &spdx.File{
		FileName:           "/tmp/whatever2.txt",
		FileSPDXIdentifier: common.ElementID("File1232"),
		Checksums:          []common.Checksum{{Value: "85ed0817af83a24ad8da68c2b5094de69833983d", Algorithm: common.SHA1}},
		LicenseConcluded:   "MIT",
		LicenseInfoInFiles: []string{"MIT"},
		FileCopyrightText:  "Copyright (c) John Doe",
	}

	unFiles := []*spdx.File{
		f1,
		f2,
	}

	// Package 1: packaged files with snippets
	sn1 := &spdx.Snippet{
		SnippetSPDXIdentifier:         "Snippet19",
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "FileHasSnippets").ElementRefID,
		Ranges:                        []common.SnippetRange{{StartPointer: common.SnippetRangePointer{Offset: 17}, EndPointer: common.SnippetRangePointer{Offset: 209}}},
		SnippetLicenseConcluded:       "GPL-2.0-or-later",
		SnippetCopyrightText:          "Copyright (c) John Doe 20x6",
	}

	sn2 := &spdx.Snippet{
		SnippetSPDXIdentifier:         "Snippet20",
		SnippetFromFileSPDXIdentifier: common.MakeDocElementID("", "FileHasSnippets").ElementRefID,
		Ranges:                        []common.SnippetRange{{StartPointer: common.SnippetRangePointer{Offset: 268}, EndPointer: common.SnippetRangePointer{Offset: 309}}},
		SnippetLicenseConcluded:       "WTFPL",
		SnippetCopyrightText:          "NOASSERTION",
	}

	f3 := &spdx.File{
		FileName:           "/tmp/file-with-snippets.txt",
		FileSPDXIdentifier: common.ElementID("FileHasSnippets"),
		Checksums:          []common.Checksum{{Value: "85ed0817af83a24ad8da68c2b5094de69833983e", Algorithm: common.SHA1}},
		LicenseConcluded:   "GPL-2.0-or-later AND WTFPL",
		LicenseInfoInFiles: []string{
			"Apache-2.0",
			"GPL-2.0-or-later",
			"WTFPL",
		},
		FileCopyrightText: "Copyright (c) Jane Doe",
		Snippets: map[common.ElementID]*spdx.Snippet{
			common.ElementID("Snippet19"): sn1,
			common.ElementID("Snippet20"): sn2,
		},
	}

	f4 := &spdx.File{
		FileName:           "/tmp/another-file.txt",
		FileSPDXIdentifier: common.ElementID("FileAnother"),
		Checksums:          []common.Checksum{{Value: "85ed0817af83a24ad8da68c2b5094de69833983f", Algorithm: common.SHA1}},
		LicenseConcluded:   "BSD-3-Clause",
		LicenseInfoInFiles: []string{"BSD-3-Clause"},
		FileCopyrightText:  "Copyright (c) Jane Doe LLC",
	}

	pkgWith := &spdx.Package{
		PackageName:               "p1",
		PackageSPDXIdentifier:     common.ElementID("p1"),
		PackageDownloadLocation:   "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:             true,
		IsFilesAnalyzedTagPresent: true,
		PackageVerificationCode:   common.PackageVerificationCode{Value: "0123456789abcdef0123456789abcdef01234567"},
		PackageLicenseConcluded:   "GPL-2.0-or-later AND BSD-3-Clause AND WTFPL",
		PackageLicenseInfoFromFiles: []string{
			"Apache-2.0",
			"GPL-2.0-or-later",
			"WTFPL",
			"BSD-3-Clause",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageCopyrightText:   "Copyright (c) John Doe, Inc.",
		Files: []*spdx.File{
			f3,
			f4,
		},
	}

	// Other Licenses 1 and 2
	ol1 := &spdx.OtherLicense{
		LicenseIdentifier: "LicenseRef-1",
		ExtractedText: `License 1 text
blah blah blah
blah blah blah blah`,
		LicenseName: "License 1",
	}

	ol2 := &spdx.OtherLicense{
		LicenseIdentifier: "LicenseRef-2",
		ExtractedText:     `License 2 text - this is a license that does some stuff`,
		LicenseName:       "License 2",
	}

	// Relationships
	rln1 := &spdx.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", "p1"),
		Relationship: "DESCRIBES",
	}

	rln2 := &spdx.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", "File1231"),
		Relationship: "DESCRIBES",
	}

	rln3 := &spdx.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", "File1232"),
		Relationship: "DESCRIBES",
	}

	// Annotations
	ann1 := &spdx.Annotation{
		Annotator: common.Annotator{Annotator: "John Doe",
			AnnotatorType: "Person"},
		AnnotationDate:           "2018-10-10T17:52:00Z",
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.MakeDocElementID("", "DOCUMENT"),
		AnnotationComment:        "This is an annotation about the SPDX document",
	}

	ann2 := &spdx.Annotation{
		Annotator: common.Annotator{Annotator: "John Doe, Inc.",
			AnnotatorType: "Organization"},
		AnnotationDate:           "2018-10-10T17:52:00Z",
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.MakeDocElementID("", "p1"),
		AnnotationComment:        "This is an annotation about Package p1",
	}

	// Reviews
	rev1 := &spdx.Review{
		Reviewer:     "John Doe",
		ReviewerType: "Person",
		ReviewDate:   "2018-10-14T10:28:00Z",
	}
	rev2 := &spdx.Review{
		Reviewer:      "Jane Doe LLC",
		ReviewerType:  "Organization",
		ReviewDate:    "2018-10-14T10:28:00Z",
		ReviewComment: "I have reviewed this SPDX document and it is awesome",
	}

	// now, build the document
	doc := &spdx.Document{
		SPDXVersion:       "SPDX-2.2",
		DataLicense:       spdx.DataLicense,
		SPDXIdentifier:    common.ElementID("DOCUMENT"),
		DocumentName:      "tools-golang-0.0.1.abcdef",
		DocumentNamespace: "https://github.com/spdx/spdx-docs/tools-golang/tools-golang-0.0.1.abcdef.whatever",
		CreationInfo:      ci,
		Packages: []*spdx.Package{
			pkgWith,
		},
		Files: unFiles,
		OtherLicenses: []*spdx.OtherLicense{
			ol1,
			ol2,
		},
		Relationships: []*spdx.Relationship{
			rln1,
			rln2,
			rln3,
		},
		Annotations: []*spdx.Annotation{
			ann1,
			ann2,
		},
		Reviews: []*spdx.Review{
			rev1,
			rev2,
		},
	}

	want := bytes.NewBufferString(`SPDXVersion: SPDX-2.2
DataLicense: CC0-1.0
SPDXID: SPDXRef-DOCUMENT
DocumentName: tools-golang-0.0.1.abcdef
DocumentNamespace: https://github.com/spdx/spdx-docs/tools-golang/tools-golang-0.0.1.abcdef.whatever
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

SnippetSPDXID: SPDXRef-Snippet19
SnippetFromFileSPDXID: SPDXRef-FileHasSnippets
SnippetByteRange: 17:209
SnippetLicenseConcluded: GPL-2.0-or-later
SnippetCopyrightText: Copyright (c) John Doe 20x6

SnippetSPDXID: SPDXRef-Snippet20
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
	err := RenderDocument(doc, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected {{{%v}}}, got {{{%v}}}", want.String(), got.String())
	}

}

func TestSaverDocumentReturnsErrorIfNilCreationInfo(t *testing.T) {
	doc := &spdx.Document{}

	var got bytes.Buffer
	err := RenderDocument(doc, &got)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
