// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Package section Saver tests =====
func TestSaver2_1PackageSavesTextCombo1(t *testing.T) {
	// include package external refs
	// test Supplier:Organization, Originator:Person
	// FilesAnalyzed true, IsFilesAnalyzedTagPresent true
	// PackageVerificationCodeExcludedFile has string

	// NOTE, this is an entirely made up CPE and the format is likely invalid
	per1 := &spdx.PackageExternalReference2_1{
		Category:           "SECURITY",
		RefType:            "cpe22Type",
		Locator:            "cpe:/a:john_doe_inc:p1:0.1.0",
		ExternalRefComment: "this is an external ref comment #1",
	}

	// NOTE, this is an entirely made up NPM
	per2 := &spdx.PackageExternalReference2_1{
		Category: "PACKAGE-MANAGER",
		RefType:  "npm",
		Locator:  "p1@0.1.0",
		ExternalRefComment: `this is a
multi-line external ref comment`,
	}

	per3 := &spdx.PackageExternalReference2_1{
		Category: "OTHER",
		RefType:  "anything",
		Locator:  "anything-without-spaces-can-go-here",
		// no ExternalRefComment for this one
	}

	truthy := true
	pkg := &spdx.Package2_1{
		PackageName:             "p1",
		PackageSPDXIdentifier:   spdx.ElementID("p1"),
		PackageVersion:          "0.1.0",
		PackageFileName:         "p1-0.1.0-master.tar.gz",
		PackageSupplier:         &spdx.Supplier{SupplierType: "Organization", Supplier: "John Doe, Inc."},
		PackageOriginator:       &spdx.Originator{Originator: "John Doe", OriginatorType: "Person"},
		PackageDownloadLocation: "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:           &truthy,
		PackageVerificationCode: spdx.PackageVerificationCode{
			Value:         "0123456789abcdef0123456789abcdef01234567",
			ExcludedFiles: []string{"p1-0.1.0.spdx"},
		},
		PackageChecksums: []spdx.Checksum{
			{
				Algorithm: spdx.SHA1,
				Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
			},
			{
				Algorithm: spdx.SHA256,
				Value:     "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
			},
			{
				Algorithm: spdx.MD5,
				Value:     "624c1abb3664f4b35547e7c73864ad24",
			},
		},
		PackageHomePage:         "http://example.com/p1",
		PackageSourceInfo:       "this is a source comment",
		PackageLicenseConcluded: "GPL-2.0-or-later",
		PackageLicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageLicenseComments: "this is a license comment(s)",
		PackageCopyrightText:   "Copyright (c) John Doe, Inc.",
		PackageSummary:         "this is a summary comment",
		PackageDescription:     "this is a description comment",
		PackageComment:         "this is a comment comment",
		PackageExternalReferences: []*spdx.PackageExternalReference2_1{
			per1,
			per2,
			per3,
		},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageVersion: 0.1.0
PackageFileName: p1-0.1.0-master.tar.gz
PackageSupplier: Organization: John Doe, Inc.
PackageOriginator: Person: John Doe
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: true
PackageVerificationCode: 0123456789abcdef0123456789abcdef01234567 (excludes: p1-0.1.0.spdx)
PackageChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
PackageChecksum: SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd
PackageChecksum: MD5: 624c1abb3664f4b35547e7c73864ad24
PackageHomePage: http://example.com/p1
PackageSourceInfo: this is a source comment
PackageLicenseConcluded: GPL-2.0-or-later
PackageLicenseInfoFromFiles: Apache-1.1
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageLicenseComments: this is a license comment(s)
PackageCopyrightText: Copyright (c) John Doe, Inc.
PackageSummary: this is a summary comment
PackageDescription: this is a description comment
PackageComment: this is a comment comment
ExternalRef: SECURITY cpe22Type cpe:/a:john_doe_inc:p1:0.1.0
ExternalRefComment: this is an external ref comment #1
ExternalRef: PACKAGE-MANAGER npm p1@0.1.0
ExternalRefComment: <text>this is a
multi-line external ref comment</text>
ExternalRef: OTHER anything anything-without-spaces-can-go-here

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderPackage2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1PackageSavesTextCombo2(t *testing.T) {
	// no package external refs
	// test Supplier:NOASSERTION, Originator:Organization
	// FilesAnalyzed true, IsFilesAnalyzedTagPresent false
	// PackageVerificationCodeExcludedFile is empty

	pkg := &spdx.Package2_1{
		PackageName:             "p1",
		PackageSPDXIdentifier:   spdx.ElementID("p1"),
		PackageVersion:          "0.1.0",
		PackageFileName:         "p1-0.1.0-master.tar.gz",
		PackageSupplier:         &spdx.Supplier{Supplier: "NOASSERTION"},
		PackageOriginator:       &spdx.Originator{OriginatorType: "Organization", Originator: "John Doe, Inc."},
		PackageDownloadLocation: "http://example.com/p1/p1-0.1.0-master.tar.gz",
		PackageVerificationCode: spdx.PackageVerificationCode{Value: "0123456789abcdef0123456789abcdef01234567"},
		PackageChecksums: []spdx.Checksum{
			{
				Algorithm: spdx.SHA1,
				Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
			},
			{
				Algorithm: spdx.SHA256,
				Value:     "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
			},
			{
				Algorithm: spdx.MD5,
				Value:     "624c1abb3664f4b35547e7c73864ad24",
			},
		},
		PackageHomePage:         "http://example.com/p1",
		PackageSourceInfo:       "this is a source comment",
		PackageLicenseConcluded: "GPL-2.0-or-later",
		PackageLicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageLicenseComments: "this is a license comment(s)",
		PackageCopyrightText:   "Copyright (c) John Doe, Inc.",
		PackageSummary:         "this is a summary comment",
		PackageDescription:     "this is a description comment",
		PackageComment:         "this is a comment comment",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageVersion: 0.1.0
PackageFileName: p1-0.1.0-master.tar.gz
PackageSupplier: NOASSERTION
PackageOriginator: Organization: John Doe, Inc.
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: true
PackageVerificationCode: 0123456789abcdef0123456789abcdef01234567
PackageChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
PackageChecksum: SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd
PackageChecksum: MD5: 624c1abb3664f4b35547e7c73864ad24
PackageHomePage: http://example.com/p1
PackageSourceInfo: this is a source comment
PackageLicenseConcluded: GPL-2.0-or-later
PackageLicenseInfoFromFiles: Apache-1.1
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageLicenseComments: this is a license comment(s)
PackageCopyrightText: Copyright (c) John Doe, Inc.
PackageSummary: this is a summary comment
PackageDescription: this is a description comment
PackageComment: this is a comment comment

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderPackage2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1PackageSavesTextCombo3(t *testing.T) {
	// no package external refs
	// test Supplier:Person, Originator:NOASSERTION
	// FilesAnalyzed false, IsFilesAnalyzedTagPresent true
	// PackageVerificationCodeExcludedFile is empty

	falsy := false
	pkg := &spdx.Package2_1{
		PackageName:             "p1",
		PackageSPDXIdentifier:   spdx.ElementID("p1"),
		PackageVersion:          "0.1.0",
		PackageFileName:         "p1-0.1.0-master.tar.gz",
		PackageSupplier:         &spdx.Supplier{Supplier: "John Doe", SupplierType: "Person"},
		PackageOriginator:       &spdx.Originator{Originator: "NOASSERTION"},
		PackageDownloadLocation: "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:           &falsy,
		// NOTE that verification code MUST be omitted from output
		// since FilesAnalyzed is false
		PackageVerificationCode: spdx.PackageVerificationCode{Value: "0123456789abcdef0123456789abcdef01234567"},
		PackageChecksums: []spdx.Checksum{
			{
				Algorithm: spdx.SHA1,
				Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
			},
			{
				Algorithm: spdx.SHA256,
				Value:     "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
			},
			{
				Algorithm: spdx.MD5,
				Value:     "624c1abb3664f4b35547e7c73864ad24",
			},
		},
		PackageHomePage:         "http://example.com/p1",
		PackageSourceInfo:       "this is a source comment",
		PackageLicenseConcluded: "GPL-2.0-or-later",
		// NOTE that license info from files MUST be omitted from output
		// since FilesAnalyzed is false
		PackageLicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageLicenseComments: "this is a license comment(s)",
		PackageCopyrightText:   "Copyright (c) John Doe, Inc.",
		PackageSummary:         "this is a summary comment",
		PackageDescription:     "this is a description comment",
		PackageComment:         "this is a comment comment",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageVersion: 0.1.0
PackageFileName: p1-0.1.0-master.tar.gz
PackageSupplier: Person: John Doe
PackageOriginator: NOASSERTION
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: false
PackageVerificationCode: 0123456789abcdef0123456789abcdef01234567
PackageChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
PackageChecksum: SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd
PackageChecksum: MD5: 624c1abb3664f4b35547e7c73864ad24
PackageHomePage: http://example.com/p1
PackageSourceInfo: this is a source comment
PackageLicenseConcluded: GPL-2.0-or-later
PackageLicenseInfoFromFiles: Apache-1.1
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageLicenseComments: this is a license comment(s)
PackageCopyrightText: Copyright (c) John Doe, Inc.
PackageSummary: this is a summary comment
PackageDescription: this is a description comment
PackageComment: this is a comment comment

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderPackage2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1PackageSaveOmitsOptionalFieldsIfEmpty(t *testing.T) {
	falsy := false
	pkg := &spdx.Package2_1{
		PackageName:             "p1",
		PackageSPDXIdentifier:   spdx.ElementID("p1"),
		PackageDownloadLocation: "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:           &falsy,
		PackageLicenseConcluded: "GPL-2.0-or-later",
		PackageLicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageCopyrightText:   "Copyright (c) John Doe, Inc.",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: false
PackageLicenseConcluded: GPL-2.0-or-later
PackageLicenseInfoFromFiles: Apache-1.1
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageCopyrightText: Copyright (c) John Doe, Inc.

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderPackage2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1PackageSavesFilesIfPresent(t *testing.T) {
	f1 := &spdx.File2_1{
		FileName:           "/tmp/whatever1.txt",
		FileSPDXIdentifier: spdx.ElementID("File1231"),
		Checksums: []spdx.Checksum{
			{
				Algorithm: spdx.SHA1,
				Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
			},
		},
		LicenseConcluded:   "Apache-2.0",
		LicenseInfoInFiles: []string{"Apache-2.0"},
		FileCopyrightText:  "Copyright (c) Jane Doe",
	}

	f2 := &spdx.File2_1{
		FileName:           "/tmp/whatever2.txt",
		FileSPDXIdentifier: spdx.ElementID("File1232"),
		Checksums: []spdx.Checksum{
			{
				Algorithm: spdx.SHA1,
				Value:     "85ed0817af83a24ad8da68c2b5094de69833983d",
			},
		},
		LicenseConcluded:   "MIT",
		LicenseInfoInFiles: []string{"MIT"},
		FileCopyrightText:  "Copyright (c) John Doe",
	}

	falsy := false
	pkg := &spdx.Package2_1{
		PackageName:             "p1",
		PackageSPDXIdentifier:   spdx.ElementID("p1"),
		PackageDownloadLocation: "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:           &falsy,
		PackageLicenseConcluded: "GPL-2.0-or-later",
		PackageLicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageCopyrightText:   "Copyright (c) John Doe, Inc.",
		Files: []*spdx.File2_1{
			f1,
			f2,
		},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: false
PackageLicenseConcluded: GPL-2.0-or-later
PackageLicenseInfoFromFiles: Apache-1.1
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageCopyrightText: Copyright (c) John Doe, Inc.

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

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderPackage2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1PackageWrapsCopyrightMultiLine(t *testing.T) {
	falsy := false
	pkg := &spdx.Package2_1{
		PackageName:             "p1",
		PackageSPDXIdentifier:   spdx.ElementID("p1"),
		PackageDownloadLocation: "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:           &falsy,
		PackageLicenseConcluded: "GPL-2.0-or-later",
		PackageLicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		PackageLicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		PackageCopyrightText: `Copyright (c) John Doe, Inc.
Copyright Jane Doe`,
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: false
PackageLicenseConcluded: GPL-2.0-or-later
PackageLicenseInfoFromFiles: Apache-1.1
PackageLicenseInfoFromFiles: Apache-2.0
PackageLicenseInfoFromFiles: GPL-2.0-or-later
PackageLicenseDeclared: Apache-2.0 OR GPL-2.0-or-later
PackageCopyrightText: <text>Copyright (c) John Doe, Inc.
Copyright Jane Doe</text>

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderPackage2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
