// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reporter

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== 2.1 Reporter top-level function tests =====
func Test2_1ReporterCanMakeReportFromPackage(t *testing.T) {
	pkg := &spdx.Package2_1{
		FilesAnalyzed: true,
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID("File0"):  &spdx.File2_1{FileSPDXIdentifier: "File0", LicenseConcluded: "MIT"},
			spdx.ElementID("File1"):  &spdx.File2_1{FileSPDXIdentifier: "File1", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File2"):  &spdx.File2_1{FileSPDXIdentifier: "File2", LicenseConcluded: "MIT"},
			spdx.ElementID("File3"):  &spdx.File2_1{FileSPDXIdentifier: "File3", LicenseConcluded: "MIT"},
			spdx.ElementID("File4"):  &spdx.File2_1{FileSPDXIdentifier: "File4", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File5"):  &spdx.File2_1{FileSPDXIdentifier: "File5", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File6"):  &spdx.File2_1{FileSPDXIdentifier: "File6", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File7"):  &spdx.File2_1{FileSPDXIdentifier: "File7", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File8"):  &spdx.File2_1{FileSPDXIdentifier: "File8", LicenseConcluded: "MIT"},
			spdx.ElementID("File9"):  &spdx.File2_1{FileSPDXIdentifier: "File9", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File10"): &spdx.File2_1{FileSPDXIdentifier: "File10", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File11"): &spdx.File2_1{FileSPDXIdentifier: "File11", LicenseConcluded: "NOASSERTION"},
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
	err := Generate2_1(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func Test2_1ReporterReturnsErrorIfPackageFilesNotAnalyzed(t *testing.T) {
	pkg := &spdx.Package2_1{
		FilesAnalyzed: false,
	}

	// render as buffer of bytes
	var got bytes.Buffer
	err := Generate2_1(pkg, &got)
	if err == nil {
		t.Errorf("Expected non-nil error, got nil")
	}
}

// ===== 2.1 Utility functions =====

func Test2_1CanGetCountsOfLicenses(t *testing.T) {
	pkg := &spdx.Package2_1{
		FilesAnalyzed: true,
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID("File0"):  &spdx.File2_1{FileSPDXIdentifier: "File0", LicenseConcluded: "MIT"},
			spdx.ElementID("File1"):  &spdx.File2_1{FileSPDXIdentifier: "File1", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File2"):  &spdx.File2_1{FileSPDXIdentifier: "File2", LicenseConcluded: "MIT"},
			spdx.ElementID("File3"):  &spdx.File2_1{FileSPDXIdentifier: "File3", LicenseConcluded: "MIT"},
			spdx.ElementID("File4"):  &spdx.File2_1{FileSPDXIdentifier: "File4", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File5"):  &spdx.File2_1{FileSPDXIdentifier: "File5", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File6"):  &spdx.File2_1{FileSPDXIdentifier: "File6", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File7"):  &spdx.File2_1{FileSPDXIdentifier: "File7", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File8"):  &spdx.File2_1{FileSPDXIdentifier: "File8", LicenseConcluded: "MIT"},
			spdx.ElementID("File9"):  &spdx.File2_1{FileSPDXIdentifier: "File9", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File10"): &spdx.File2_1{FileSPDXIdentifier: "File10", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File11"): &spdx.File2_1{FileSPDXIdentifier: "File11", LicenseConcluded: "NOASSERTION"},
		},
	}

	totalFound, totalNotFound, foundCounts := countLicenses2_1(pkg)
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

func Test2_1NilPackageReturnsZeroCountsOfLicenses(t *testing.T) {
	totalFound, totalNotFound, foundCounts := countLicenses2_1(nil)
	if totalFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalFound)
	}
	if totalNotFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalNotFound)
	}
	if len(foundCounts) != 0 {
		t.Fatalf("expected %v, got %v", 0, len(foundCounts))
	}

	pkg := &spdx.Package2_1{}
	totalFound, totalNotFound, foundCounts = countLicenses2_1(pkg)
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

// ===== 2.2 Reporter top-level function tests =====
func Test2_2ReporterCanMakeReportFromPackage(t *testing.T) {
	pkg := &spdx.Package2_2{
		FilesAnalyzed: true,
		Files: map[spdx.ElementID]*spdx.File2_2{
			spdx.ElementID("File0"):  &spdx.File2_2{FileSPDXIdentifier: "File0", LicenseConcluded: "MIT"},
			spdx.ElementID("File1"):  &spdx.File2_2{FileSPDXIdentifier: "File1", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File2"):  &spdx.File2_2{FileSPDXIdentifier: "File2", LicenseConcluded: "MIT"},
			spdx.ElementID("File3"):  &spdx.File2_2{FileSPDXIdentifier: "File3", LicenseConcluded: "MIT"},
			spdx.ElementID("File4"):  &spdx.File2_2{FileSPDXIdentifier: "File4", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File5"):  &spdx.File2_2{FileSPDXIdentifier: "File5", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File6"):  &spdx.File2_2{FileSPDXIdentifier: "File6", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File7"):  &spdx.File2_2{FileSPDXIdentifier: "File7", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File8"):  &spdx.File2_2{FileSPDXIdentifier: "File8", LicenseConcluded: "MIT"},
			spdx.ElementID("File9"):  &spdx.File2_2{FileSPDXIdentifier: "File9", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File10"): &spdx.File2_2{FileSPDXIdentifier: "File10", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File11"): &spdx.File2_2{FileSPDXIdentifier: "File11", LicenseConcluded: "NOASSERTION"},
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
	err := Generate2_2(pkg, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func Test2_2ReporterReturnsErrorIfPackageFilesNotAnalyzed(t *testing.T) {
	pkg := &spdx.Package2_2{
		FilesAnalyzed: false,
	}

	// render as buffer of bytes
	var got bytes.Buffer
	err := Generate2_2(pkg, &got)
	if err == nil {
		t.Errorf("Expected non-nil error, got nil")
	}
}

// ===== 2.2 Utility functions =====

func Test2_2CanGetCountsOfLicenses(t *testing.T) {
	pkg := &spdx.Package2_2{
		FilesAnalyzed: true,
		Files: map[spdx.ElementID]*spdx.File2_2{
			spdx.ElementID("File0"):  &spdx.File2_2{FileSPDXIdentifier: "File0", LicenseConcluded: "MIT"},
			spdx.ElementID("File1"):  &spdx.File2_2{FileSPDXIdentifier: "File1", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File2"):  &spdx.File2_2{FileSPDXIdentifier: "File2", LicenseConcluded: "MIT"},
			spdx.ElementID("File3"):  &spdx.File2_2{FileSPDXIdentifier: "File3", LicenseConcluded: "MIT"},
			spdx.ElementID("File4"):  &spdx.File2_2{FileSPDXIdentifier: "File4", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File5"):  &spdx.File2_2{FileSPDXIdentifier: "File5", LicenseConcluded: "NOASSERTION"},
			spdx.ElementID("File6"):  &spdx.File2_2{FileSPDXIdentifier: "File6", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File7"):  &spdx.File2_2{FileSPDXIdentifier: "File7", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File8"):  &spdx.File2_2{FileSPDXIdentifier: "File8", LicenseConcluded: "MIT"},
			spdx.ElementID("File9"):  &spdx.File2_2{FileSPDXIdentifier: "File9", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File10"): &spdx.File2_2{FileSPDXIdentifier: "File10", LicenseConcluded: "GPL-2.0-only"},
			spdx.ElementID("File11"): &spdx.File2_2{FileSPDXIdentifier: "File11", LicenseConcluded: "NOASSERTION"},
		},
	}

	totalFound, totalNotFound, foundCounts := countLicenses2_2(pkg)
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

func Test2_2NilPackageReturnsZeroCountsOfLicenses(t *testing.T) {
	totalFound, totalNotFound, foundCounts := countLicenses2_2(nil)
	if totalFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalFound)
	}
	if totalNotFound != 0 {
		t.Errorf("expected %v, got %v", 0, totalNotFound)
	}
	if len(foundCounts) != 0 {
		t.Fatalf("expected %v, got %v", 0, len(foundCounts))
	}

	pkg := &spdx.Package2_2{}
	totalFound, totalNotFound, foundCounts = countLicenses2_2(pkg)
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
