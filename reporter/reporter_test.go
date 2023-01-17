// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reporter

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestReporterCanMakeReportFromPackage(t *testing.T) {
	pkg := &spdx.Package{
		FilesAnalyzed: true,
		Files: []*spdx.File{
			{FileSPDXIdentifier: "File0", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File1", LicenseConcluded: "NOASSERTION"},
			{FileSPDXIdentifier: "File2", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File3", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File4", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File5", LicenseConcluded: "NOASSERTION"},
			{FileSPDXIdentifier: "File6", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File7", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File8", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File9", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File10", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File11", LicenseConcluded: "NOASSERTION"},
		},
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`   9  License found
   3  License not found
  12  TOTAL

  5  GPL-2.0-only
  4  MIT
  9  TOTAL FOUND
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := Generate(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestReporterReturnsErrorIfPackageFilesNotAnalyzed(t *testing.T) {
	pkg := &spdx.Package{
		FilesAnalyzed: false,
	}

	// render as buffer of bytes
	var got bytes.Buffer
	err := Generate(pkg, &got)
	if err == nil {
		t.Errorf("Expected non-nil error, got nil")
	}
}

// ===== 2.3 Utility functions =====

func TestCanGetCountsOfLicenses(t *testing.T) {
	pkg := &spdx.Package{
		FilesAnalyzed: true,
		Files: []*spdx.File{
			{FileSPDXIdentifier: "File0", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File1", LicenseConcluded: "NOASSERTION"},
			{FileSPDXIdentifier: "File2", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File3", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File4", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File5", LicenseConcluded: "NOASSERTION"},
			{FileSPDXIdentifier: "File6", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File7", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File8", LicenseConcluded: "MIT"},
			{FileSPDXIdentifier: "File9", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File10", LicenseConcluded: "GPL-2.0-only"},
			{FileSPDXIdentifier: "File11", LicenseConcluded: "NOASSERTION"},
		},
	}

	totalFound, totalNotFound, foundCounts := countLicenses(pkg)
	if totalFound != 9 {
		t.Errorf("expected %v, got %v", 9, totalFound)
	}
	if totalNotFound != 3 {
		t.Errorf("expected %v, got %v", 3, totalNotFound)
	}
	if len(foundCounts) != 2 {
		t.Fatalf("expected %v, got %v", 2, len(foundCounts))
	}

	// foundCounts is a map of license ID to count of licenses
	// confirm that the results are as expected
	if foundCounts["GPL-2.0-only"] != 5 {
		t.Errorf("expected %v, got %v", 5, foundCounts["GPL-2.0-only"])
	}
	if foundCounts["MIT"] != 4 {
		t.Errorf("expected %v, got %v", 4, foundCounts["MIT"])
	}
}

func TestNilPackageReturnsZeroCountsOfLicenses(t *testing.T) {
	totalFound, totalNotFound, foundCounts := countLicenses(nil)
	if totalFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalFound)
	}
	if totalNotFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalNotFound)
	}
	if len(foundCounts) != 0 {
		t.Fatalf("expected %v, got %v", 0, len(foundCounts))
	}

	pkg := &spdx.Package{}
	totalFound, totalNotFound, foundCounts = countLicenses(pkg)
	if totalFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalFound)
	}
	if totalNotFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalNotFound)
	}
	if len(foundCounts) != 0 {
		t.Fatalf("expected %v, got %v", 0, len(foundCounts))
	}
}
