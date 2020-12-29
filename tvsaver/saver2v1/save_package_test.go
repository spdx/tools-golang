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
		// no ExternalRefComment for this one
	}

	per3 := &spdx.PackageExternalReference2_1{
		Category: "OTHER",
		RefType:  "anything",
		Locator:  "anything-without-spaces-can-go-here",
		// no ExternalRefComment for this one
	}

	pkg := &spdx.Package2_1{
		Name:                         "p1",
		SPDXIdentifier:               spdx.ElementID("p1"),
		Version:                      "0.1.0",
		FileName:                     "p1-0.1.0-master.tar.gz",
		SupplierOrganization:         "John Doe, Inc.",
		OriginatorPerson:             "John Doe",
		DownloadLocation:             "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:                true,
		IsFilesAnalyzedTagPresent:    true,
		VerificationCode:             "0123456789abcdef0123456789abcdef01234567",
		VerificationCodeExcludedFile: "p1-0.1.0.spdx",
		ChecksumSHA1:                 "85ed0817af83a24ad8da68c2b5094de69833983c",
		ChecksumSHA256:               "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
		ChecksumMD5:                  "624c1abb3664f4b35547e7c73864ad24",
		HomePage:                     "http://example.com/p1",
		SourceInfo:                   "this is a source comment",
		LicenseConcluded:             "GPL-2.0-or-later",
		LicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		LicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		LicenseComments: "this is a license comment(s)",
		CopyrightText:   "Copyright (c) John Doe, Inc.",
		Summary:         "this is a summary comment",
		Description:     "this is a description comment",
		Comment:         "this is a comment comment",
		ExternalReferences: []*spdx.PackageExternalReference2_1{
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
PackageVerificationCode: 0123456789abcdef0123456789abcdef01234567 (excludes p1-0.1.0.spdx)
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
		Name:                      "p1",
		SPDXIdentifier:            spdx.ElementID("p1"),
		Version:                   "0.1.0",
		FileName:                  "p1-0.1.0-master.tar.gz",
		SupplierNOASSERTION:       true,
		OriginatorOrganization:    "John Doe, Inc.",
		DownloadLocation:          "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:             true,
		IsFilesAnalyzedTagPresent: false,
		VerificationCode:          "0123456789abcdef0123456789abcdef01234567",
		ChecksumSHA1:              "85ed0817af83a24ad8da68c2b5094de69833983c",
		ChecksumSHA256:            "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
		ChecksumMD5:               "624c1abb3664f4b35547e7c73864ad24",
		HomePage:                  "http://example.com/p1",
		SourceInfo:                "this is a source comment",
		LicenseConcluded:          "GPL-2.0-or-later",
		LicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		LicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		LicenseComments: "this is a license comment(s)",
		CopyrightText:   "Copyright (c) John Doe, Inc.",
		Summary:         "this is a summary comment",
		Description:     "this is a description comment",
		Comment:         "this is a comment comment",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageVersion: 0.1.0
PackageFileName: p1-0.1.0-master.tar.gz
PackageSupplier: NOASSERTION
PackageOriginator: Organization: John Doe, Inc.
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
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

	pkg := &spdx.Package2_1{
		Name:                      "p1",
		SPDXIdentifier:            spdx.ElementID("p1"),
		Version:                   "0.1.0",
		FileName:                  "p1-0.1.0-master.tar.gz",
		SupplierPerson:            "John Doe",
		OriginatorNOASSERTION:     true,
		DownloadLocation:          "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:             false,
		IsFilesAnalyzedTagPresent: true,
		// NOTE that verification code MUST be omitted from output
		// since FilesAnalyzed is false
		VerificationCode: "0123456789abcdef0123456789abcdef01234567",
		ChecksumSHA1:     "85ed0817af83a24ad8da68c2b5094de69833983c",
		ChecksumSHA256:   "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
		ChecksumMD5:      "624c1abb3664f4b35547e7c73864ad24",
		HomePage:         "http://example.com/p1",
		SourceInfo:       "this is a source comment",
		LicenseConcluded: "GPL-2.0-or-later",
		// NOTE that license info from files MUST be omitted from output
		// since FilesAnalyzed is false
		LicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		LicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		LicenseComments: "this is a license comment(s)",
		CopyrightText:   "Copyright (c) John Doe, Inc.",
		Summary:         "this is a summary comment",
		Description:     "this is a description comment",
		Comment:         "this is a comment comment",
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
PackageChecksum: SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c
PackageChecksum: SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd
PackageChecksum: MD5: 624c1abb3664f4b35547e7c73864ad24
PackageHomePage: http://example.com/p1
PackageSourceInfo: this is a source comment
PackageLicenseConcluded: GPL-2.0-or-later
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
	pkg := &spdx.Package2_1{
		Name:                      "p1",
		SPDXIdentifier:            spdx.ElementID("p1"),
		DownloadLocation:          "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:             false,
		IsFilesAnalyzedTagPresent: true,
		// NOTE that verification code MUST be omitted from output,
		// even if present in model, since FilesAnalyzed is false
		LicenseConcluded: "GPL-2.0-or-later",
		// NOTE that license info from files MUST be omitted from output
		// even if present in model, since FilesAnalyzed is false
		LicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		LicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		CopyrightText:   "Copyright (c) John Doe, Inc.",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: false
PackageLicenseConcluded: GPL-2.0-or-later
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

	pkg := &spdx.Package2_1{
		Name:                      "p1",
		SPDXIdentifier:            spdx.ElementID("p1"),
		DownloadLocation:          "http://example.com/p1/p1-0.1.0-master.tar.gz",
		FilesAnalyzed:             false,
		IsFilesAnalyzedTagPresent: true,
		// NOTE that verification code MUST be omitted from output,
		// even if present in model, since FilesAnalyzed is false
		LicenseConcluded: "GPL-2.0-or-later",
		// NOTE that license info from files MUST be omitted from output
		// even if present in model, since FilesAnalyzed is false
		LicenseInfoFromFiles: []string{
			"Apache-1.1",
			"Apache-2.0",
			"GPL-2.0-or-later",
		},
		LicenseDeclared: "Apache-2.0 OR GPL-2.0-or-later",
		CopyrightText:   "Copyright (c) John Doe, Inc.",
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID("File1231"): f1,
			spdx.ElementID("File1232"): f2,
		},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`PackageName: p1
SPDXID: SPDXRef-p1
PackageDownloadLocation: http://example.com/p1/p1-0.1.0-master.tar.gz
FilesAnalyzed: false
PackageLicenseConcluded: GPL-2.0-or-later
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
