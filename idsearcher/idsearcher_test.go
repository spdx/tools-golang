// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package idsearcher

import (
	"testing"
)

// ===== 2.1 Searcher top-level function tests =====
func Test2_1SearcherCanFillInIDs(t *testing.T) {
	packageName := "project2"
	dirRoot := "../testdata/project2/"
	config := &Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
	}

	doc, err := BuildIDsDocument2_1(packageName, dirRoot, config)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil Document, got nil")
	}

	// not checking all contents of doc, see builder tests for those

	// get the package and its files, checking size of each
	if doc.Packages == nil {
		t.Fatalf("expected non-nil Packages, got nil")
	}
	if len(doc.Packages) != 1 {
		t.Fatalf("expected Packages len to be 1, got %d", len(doc.Packages))
	}
	pkg := doc.Packages[0]
	if pkg == nil {
		t.Fatalf("expected non-nil pkg, got nil")
	}

	if pkg.Files == nil {
		t.Fatalf("expected non-nil Files, got nil")
	}
	if len(pkg.Files) != 6 {
		t.Fatalf("expected Files len to be 6, got %d", len(pkg.Files))
	}

	fileInFolder := pkg.Files[0]
	if fileInFolder.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileInFolder.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileInFolder.LicenseInfoInFiles))
	}
	if fileInFolder.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileInFolder.LicenseInfoInFiles[0])
	}
	if fileInFolder.LicenseConcluded != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileInFolder.LicenseConcluded)
	}

	fileTrailingComment := pkg.Files[1]
	if fileTrailingComment.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileTrailingComment.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileTrailingComment.LicenseInfoInFiles))
	}
	if fileTrailingComment.LicenseInfoInFiles[0] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", fileTrailingComment.LicenseInfoInFiles[0])
	}
	if fileTrailingComment.LicenseConcluded != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", fileTrailingComment.LicenseConcluded)
	}

	fileHasDuplicateID := pkg.Files[2]
	if fileHasDuplicateID.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileHasDuplicateID.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileHasDuplicateID.LicenseInfoInFiles))
	}
	if fileHasDuplicateID.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileHasDuplicateID.LicenseInfoInFiles[0])
	}
	if fileHasDuplicateID.LicenseConcluded != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileHasDuplicateID.LicenseConcluded)
	}

	fileHasID := pkg.Files[3]
	if fileHasID.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileHasID.LicenseInfoInFiles) != 2 {
		t.Fatalf("expected LicenseInfoInFiles len to be 2, got %d", len(fileHasID.LicenseInfoInFiles))
	}
	if fileHasID.LicenseInfoInFiles[0] != "Apache-2.0" {
		t.Errorf("expected %v, got %v", "Apache-2.0", fileHasID.LicenseInfoInFiles[0])
	}
	if fileHasID.LicenseInfoInFiles[1] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", fileHasID.LicenseInfoInFiles[1])
	}
	if fileHasID.LicenseConcluded != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "Apache-2.0 OR GPL-2.0-or-later", fileHasID.LicenseConcluded)
	}

	fileMultipleIDs := pkg.Files[4]
	if fileMultipleIDs.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileMultipleIDs.LicenseInfoInFiles) != 5 {
		t.Fatalf("expected LicenseInfoInFiles len to be 5, got %d", len(fileMultipleIDs.LicenseInfoInFiles))
	}
	if fileMultipleIDs.LicenseInfoInFiles[0] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", fileMultipleIDs.LicenseInfoInFiles[0])
	}
	if fileMultipleIDs.LicenseInfoInFiles[1] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", fileMultipleIDs.LicenseInfoInFiles[1])
	}
	// here, DO NOT keep the +
	if fileMultipleIDs.LicenseInfoInFiles[2] != "EPL-1.0" {
		t.Errorf("expected %v, got %v", "EPL-1.0", fileMultipleIDs.LicenseInfoInFiles[2])
	}
	if fileMultipleIDs.LicenseInfoInFiles[3] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", fileMultipleIDs.LicenseInfoInFiles[3])
	}
	if fileMultipleIDs.LicenseInfoInFiles[4] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileMultipleIDs.LicenseInfoInFiles[4])
	}
	if fileMultipleIDs.LicenseConcluded != "((MIT AND BSD-3-Clause) OR ISC) AND BSD-2-Clause AND EPL-1.0+" {
		t.Errorf("expected %v, got %v", "((MIT AND BSD-3-Clause) OR ISC) AND BSD-2-Clause AND EPL-1.0+", fileMultipleIDs.LicenseConcluded)
	}

	fileNoID := pkg.Files[5]
	if fileNoID.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileNoID.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileNoID.LicenseInfoInFiles))
	}
	if fileNoID.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileNoID.LicenseInfoInFiles[0])
	}
	if fileNoID.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileNoID.LicenseConcluded)
	}

	// and finally, the package should have all of these licenses
	if pkg.PackageLicenseInfoFromFiles == nil {
		t.Fatalf("expected non-nil PackageLicenseInfoFromFiles, got nil")
	}
	if len(pkg.PackageLicenseInfoFromFiles) != 7 {
		t.Fatalf("expected PackageLicenseInfoFromFiles len to be 7, got %d", len(pkg.PackageLicenseInfoFromFiles))
	}
	if pkg.PackageLicenseInfoFromFiles[0] != "Apache-2.0" {
		t.Errorf("expected %v, got %v", "Apache-2.0", pkg.PackageLicenseInfoFromFiles[0])
	}
	if pkg.PackageLicenseInfoFromFiles[1] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", pkg.PackageLicenseInfoFromFiles[1])
	}
	if pkg.PackageLicenseInfoFromFiles[2] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", pkg.PackageLicenseInfoFromFiles[2])
	}
	// here, DO NOT keep the +
	if pkg.PackageLicenseInfoFromFiles[3] != "EPL-1.0" {
		t.Errorf("expected %v, got %v", "EPL-1.0", pkg.PackageLicenseInfoFromFiles[3])
	}
	if pkg.PackageLicenseInfoFromFiles[4] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", pkg.PackageLicenseInfoFromFiles[4])
	}
	if pkg.PackageLicenseInfoFromFiles[5] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", pkg.PackageLicenseInfoFromFiles[5])
	}
	if pkg.PackageLicenseInfoFromFiles[6] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", pkg.PackageLicenseInfoFromFiles[6])
	}

}

func Test2_1SearcherCanFillInIDsAndIgnorePaths(t *testing.T) {
	packageName := "project3"
	dirRoot := "../testdata/project3/"
	config := &Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
		BuilderPathsIgnored: []string{
			"**/ignoredir/",
			"/excludedir/",
			"**/ignorefile.txt",
			"/alsoEXCLUDEthis.txt",
		},
		SearcherPathsIgnored: []string{
			"**/dontscan.txt",
		},
	}

	doc, err := BuildIDsDocument2_1(packageName, dirRoot, config)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil Document, got nil")
	}

	// not checking all contents of doc, see builder tests for those

	// get the package and its files, checking licenses for each, and
	// confirming NOASSERTION for those that are skipped
	pkg := doc.Packages[0]
	if pkg == nil {
		t.Fatalf("expected non-nil pkg, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected len %d, got %d", 5, len(pkg.Files))
	}

	f := pkg.Files[0]
	if f.FileName != "./dontscan.txt" {
		t.Errorf("expected %v, got %v", "./dontscan.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseConcluded)
	}

	f = pkg.Files[1]
	if f.FileName != "./keep/keep.txt" {
		t.Errorf("expected %v, got %v", "./keep/keep.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseConcluded)
	}

	f = pkg.Files[2]
	if f.FileName != "./keep.txt" {
		t.Errorf("expected %v, got %v", "./keep.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseConcluded)
	}

	f = pkg.Files[3]
	if f.FileName != "./subdir/keep/dontscan.txt" {
		t.Errorf("expected %v, got %v", "./subdir/keep/dontscan.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseConcluded)
	}

	f = pkg.Files[4]
	if f.FileName != "./subdir/keep/keep.txt" {
		t.Errorf("expected %v, got %v", "./subdir/keep/keep.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseConcluded)
	}
}

func Test2_1SearcherFailsWithInvalidPath(t *testing.T) {
	packageName := "project2"
	dirRoot := "./oops/invalid"
	config := &Config2_1{
		NamespacePrefix: "whatever",
	}

	_, err := BuildIDsDocument2_1(packageName, dirRoot, config)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.2 Searcher top-level function tests =====
func Test2_2SearcherCanFillInIDs(t *testing.T) {
	packageName := "project2"
	dirRoot := "../testdata/project2/"
	config := &Config2_2{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
	}

	doc, err := BuildIDsDocument2_2(packageName, dirRoot, config)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil Document, got nil")
	}

	// not checking all contents of doc, see builder tests for those

	// get the package and its files, checking size of each
	if doc.Packages == nil {
		t.Fatalf("expected non-nil Packages, got nil")
	}
	if len(doc.Packages) != 1 {
		t.Fatalf("expected Packages len to be 1, got %d", len(doc.Packages))
	}
	pkg := doc.Packages[0]
	if pkg == nil {
		t.Fatalf("expected non-nil pkg, got nil")
	}

	if pkg.Files == nil {
		t.Fatalf("expected non-nil Files, got nil")
	}
	if len(pkg.Files) != 6 {
		t.Fatalf("expected Files len to be 6, got %d", len(pkg.Files))
	}

	fileInFolder := pkg.Files[0]
	if fileInFolder.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileInFolder.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileInFolder.LicenseInfoInFiles))
	}
	if fileInFolder.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileInFolder.LicenseInfoInFiles[0])
	}
	if fileInFolder.LicenseConcluded != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileInFolder.LicenseConcluded)
	}

	fileTrailingComment := pkg.Files[1]
	if fileTrailingComment.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileTrailingComment.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileTrailingComment.LicenseInfoInFiles))
	}
	if fileTrailingComment.LicenseInfoInFiles[0] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", fileTrailingComment.LicenseInfoInFiles[0])
	}
	if fileTrailingComment.LicenseConcluded != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", fileTrailingComment.LicenseConcluded)
	}

	fileHasDuplicateID := pkg.Files[2]
	if fileHasDuplicateID.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileHasDuplicateID.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileHasDuplicateID.LicenseInfoInFiles))
	}
	if fileHasDuplicateID.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileHasDuplicateID.LicenseInfoInFiles[0])
	}
	if fileHasDuplicateID.LicenseConcluded != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileHasDuplicateID.LicenseConcluded)
	}

	fileHasID := pkg.Files[3]
	if fileHasID.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileHasID.LicenseInfoInFiles) != 2 {
		t.Fatalf("expected LicenseInfoInFiles len to be 2, got %d", len(fileHasID.LicenseInfoInFiles))
	}
	if fileHasID.LicenseInfoInFiles[0] != "Apache-2.0" {
		t.Errorf("expected %v, got %v", "Apache-2.0", fileHasID.LicenseInfoInFiles[0])
	}
	if fileHasID.LicenseInfoInFiles[1] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", fileHasID.LicenseInfoInFiles[1])
	}
	if fileHasID.LicenseConcluded != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "Apache-2.0 OR GPL-2.0-or-later", fileHasID.LicenseConcluded)
	}

	fileMultipleIDs := pkg.Files[4]
	if fileMultipleIDs.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileMultipleIDs.LicenseInfoInFiles) != 5 {
		t.Fatalf("expected LicenseInfoInFiles len to be 5, got %d", len(fileMultipleIDs.LicenseInfoInFiles))
	}
	if fileMultipleIDs.LicenseInfoInFiles[0] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", fileMultipleIDs.LicenseInfoInFiles[0])
	}
	if fileMultipleIDs.LicenseInfoInFiles[1] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", fileMultipleIDs.LicenseInfoInFiles[1])
	}
	// here, DO NOT keep the +
	if fileMultipleIDs.LicenseInfoInFiles[2] != "EPL-1.0" {
		t.Errorf("expected %v, got %v", "EPL-1.0", fileMultipleIDs.LicenseInfoInFiles[2])
	}
	if fileMultipleIDs.LicenseInfoInFiles[3] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", fileMultipleIDs.LicenseInfoInFiles[3])
	}
	if fileMultipleIDs.LicenseInfoInFiles[4] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", fileMultipleIDs.LicenseInfoInFiles[4])
	}
	if fileMultipleIDs.LicenseConcluded != "((MIT AND BSD-3-Clause) OR ISC) AND BSD-2-Clause AND EPL-1.0+" {
		t.Errorf("expected %v, got %v", "((MIT AND BSD-3-Clause) OR ISC) AND BSD-2-Clause AND EPL-1.0+", fileMultipleIDs.LicenseConcluded)
	}

	fileNoID := pkg.Files[5]
	if fileNoID.LicenseInfoInFiles == nil {
		t.Fatalf("expected non-nil LicenseInfoInFiles, got nil")
	}
	if len(fileNoID.LicenseInfoInFiles) != 1 {
		t.Fatalf("expected LicenseInfoInFiles len to be 1, got %d", len(fileNoID.LicenseInfoInFiles))
	}
	if fileNoID.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileNoID.LicenseInfoInFiles[0])
	}
	if fileNoID.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileNoID.LicenseConcluded)
	}

	// and finally, the package should have all of these licenses
	if pkg.PackageLicenseInfoFromFiles == nil {
		t.Fatalf("expected non-nil PackageLicenseInfoFromFiles, got nil")
	}
	if len(pkg.PackageLicenseInfoFromFiles) != 7 {
		t.Fatalf("expected PackageLicenseInfoFromFiles len to be 7, got %d", len(pkg.PackageLicenseInfoFromFiles))
	}
	if pkg.PackageLicenseInfoFromFiles[0] != "Apache-2.0" {
		t.Errorf("expected %v, got %v", "Apache-2.0", pkg.PackageLicenseInfoFromFiles[0])
	}
	if pkg.PackageLicenseInfoFromFiles[1] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", pkg.PackageLicenseInfoFromFiles[1])
	}
	if pkg.PackageLicenseInfoFromFiles[2] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", pkg.PackageLicenseInfoFromFiles[2])
	}
	// here, DO NOT keep the +
	if pkg.PackageLicenseInfoFromFiles[3] != "EPL-1.0" {
		t.Errorf("expected %v, got %v", "EPL-1.0", pkg.PackageLicenseInfoFromFiles[3])
	}
	if pkg.PackageLicenseInfoFromFiles[4] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", pkg.PackageLicenseInfoFromFiles[4])
	}
	if pkg.PackageLicenseInfoFromFiles[5] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", pkg.PackageLicenseInfoFromFiles[5])
	}
	if pkg.PackageLicenseInfoFromFiles[6] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", pkg.PackageLicenseInfoFromFiles[6])
	}

}

func Test2_2SearcherCanFillInIDsAndIgnorePaths(t *testing.T) {
	packageName := "project3"
	dirRoot := "../testdata/project3/"
	config := &Config2_2{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
		BuilderPathsIgnored: []string{
			"**/ignoredir/",
			"/excludedir/",
			"**/ignorefile.txt",
			"/alsoEXCLUDEthis.txt",
		},
		SearcherPathsIgnored: []string{
			"**/dontscan.txt",
		},
	}

	doc, err := BuildIDsDocument2_2(packageName, dirRoot, config)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil Document, got nil")
	}

	// not checking all contents of doc, see builder tests for those

	// get the package and its files, checking licenses for each, and
	// confirming NOASSERTION for those that are skipped
	pkg := doc.Packages[0]
	if pkg == nil {
		t.Fatalf("expected non-nil pkg, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected len %d, got %d", 5, len(pkg.Files))
	}

	f := pkg.Files[0]
	if f.FileName != "./dontscan.txt" {
		t.Errorf("expected %v, got %v", "./dontscan.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseConcluded)
	}

	f = pkg.Files[1]
	if f.FileName != "./keep/keep.txt" {
		t.Errorf("expected %v, got %v", "./keep/keep.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseConcluded)
	}

	f = pkg.Files[2]
	if f.FileName != "./keep.txt" {
		t.Errorf("expected %v, got %v", "./keep.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseConcluded)
	}

	f = pkg.Files[3]
	if f.FileName != "./subdir/keep/dontscan.txt" {
		t.Errorf("expected %v, got %v", "./subdir/keep/dontscan.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %s, got %s", "NOASSERTION", f.LicenseConcluded)
	}

	f = pkg.Files[4]
	if f.FileName != "./subdir/keep/keep.txt" {
		t.Errorf("expected %v, got %v", "./subdir/keep/keep.txt", f.FileName)
	}
	if len(f.LicenseInfoInFiles) != 1 {
		t.Errorf("expected len to be %d, got %d", 1, len(f.LicenseInfoInFiles))
	}
	if f.LicenseInfoInFiles[0] != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseInfoInFiles[0])
	}
	if f.LicenseConcluded != "MIT" {
		t.Errorf("expected %s, got %s", "MIT", f.LicenseConcluded)
	}
}

func Test2_2SearcherFailsWithInvalidPath(t *testing.T) {
	packageName := "project2"
	dirRoot := "./oops/invalid"
	config := &Config2_2{
		NamespacePrefix: "whatever",
	}

	_, err := BuildIDsDocument2_2(packageName, dirRoot, config)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== Searcher utility tests =====
func TestCanFindShortFormIDWhenPresent(t *testing.T) {
	filePath := "../testdata/project2/has-id.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("expected len 1, got %d", len(ids))
	}

	if ids[0] != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "Apache-2.0 OR GPL-2.0-or-later", ids[0])
	}
}

func TestCanFindMultipleShortFormIDsWhenPresent(t *testing.T) {
	filePath := "../testdata/project2/has-multiple-ids.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("expected len 3, got %d", len(ids))
	}

	if ids[0] != "(MIT AND BSD-3-Clause) OR ISC" {
		t.Errorf("expected %v, got %v", "(MIT AND BSD-3-Clause) OR ISC", ids[0])
	}
	if ids[1] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", ids[1])
	}
	if ids[2] != "EPL-1.0+" {
		t.Errorf("expected %v, got %v", "EPL-1.0+", ids[2])
	}
}

func TestCanCollapseDuplicateShortFormIDsWhenPresent(t *testing.T) {
	filePath := "../testdata/project2/has-duplicate-ids.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("expected len 1, got %d", len(ids))
	}

	if ids[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", ids[0])
	}
}

func TestCanStripTrailingStarSlash(t *testing.T) {
	filePath := "../testdata/project2/folder/has-trailing-comment-marker.c"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("expected len 1, got %d", len(ids))
	}

	if ids[0] != "GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "GPL-2.0-or-later", ids[0])
	}
}

func TestCanIgnoreShortFormIDWhenTooManyPrefixChars(t *testing.T) {
	filePath := "../testdata/project4/has-id-to-ignore.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 0 {
		t.Fatalf("expected len 0, got %d", len(ids))
	}
}

func TestCanPickJustTheRightID(t *testing.T) {
	filePath := "../testdata/project4/has-mix-of-ids.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("expected len 1, got %d", len(ids))
	}

	if ids[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", ids[0])
	}
}

func TestCannotFindShortFormIDWhenAbsent(t *testing.T) {
	filePath := "../testdata/project2/no-id.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 0 {
		t.Fatalf("expected len 0, got %d", len(ids))
	}
}

func TestCanExcludeTrashCharactersFromID(t *testing.T) {
	lid := "Apac\",he-2.0"
	want := "Apache-2.0"
	got := stripTrash(lid)
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	lid = "Apache-2.0"
	want = "Apache-2.0"
	got = stripTrash(lid)
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestSearchFileIDsFailsWithInvalidFilePath(t *testing.T) {
	filePath := "./oops/nm/invalid"

	_, err := searchFileIDs(filePath)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestWillParenthesizeIfNeeded(t *testing.T) {
	licID := "MIT OR BSD-3-Clause"
	retval := makeElement(licID)
	if retval != "(MIT OR BSD-3-Clause)" {
		t.Errorf("expected %v, got %v", "(MIT OR BSD-3-Clause)", retval)
	}

	licID = "ISC AND HPND"
	retval = makeElement(licID)
	if retval != "(ISC AND HPND)" {
		t.Errorf("expected %v, got %v", "(ISC AND HPND)", retval)
	}
}

func TestWillNotParenthesizeIfNotNeeded(t *testing.T) {
	lic := "MIT"
	retval := makeElement(lic)
	if retval != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", retval)
	}

	lic = "GPL-2.0-only WITH Classpath-exception-2.0"
	retval = makeElement(lic)
	if retval != "GPL-2.0-only WITH Classpath-exception-2.0" {
		t.Errorf("expected %v, got %v", "GPL-2.0-only WITH Classpath-exception-2.0", retval)
	}
}

func TestCanGetIndividualLicenses(t *testing.T) {
	// single license
	lic := "MIT"
	lics := getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	if len(lics) != 1 {
		t.Fatalf("expected lics to have len 1, got %d", len(lics))
	}
	if lics[0] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", lics[0])
	}

	// two-license combo
	lic = "ISC AND BSD-3-Clause"
	lics = getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	if len(lics) != 2 {
		t.Fatalf("expected lics to have len 2, got %d", len(lics))
	}
	// should be sorted alphabetically
	if lics[0] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", lics[0])
	}
	if lics[1] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", lics[1])
	}

	// license WITH exception
	lic = "GPL-2.0-only WITH Classpath-exception-2.0"
	lics = getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	if len(lics) != 2 {
		t.Fatalf("expected lics to have len 2, got %d", len(lics))
	}
	// exception should be listed separately
	if lics[0] != "Classpath-exception-2.0" {
		t.Errorf("expected %v, got %v", "Classpath-exception-2.0", lics[0])
	}
	if lics[1] != "GPL-2.0-only" {
		t.Errorf("expected %v, got %v", "GPL-2.0-only", lics[1])
	}

	// two-license combo with parens
	lic = "(JSON OR BSD-2-Clause)"
	lics = getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	if len(lics) != 2 {
		t.Fatalf("expected lics to have len 2, got %d", len(lics))
	}
	// parens should get dropped
	if lics[0] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", lics[0])
	}
	if lics[1] != "JSON" {
		t.Errorf("expected %v, got %v", "JSON", lics[1])
	}

	// multi-license combo with nested parens
	lic = "GPL-2.0-only AND ((EPL-1.0 AND BSD-4-Clause) OR MIT)"
	lics = getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	if len(lics) != 4 {
		t.Fatalf("expected lics to have len 4, got %d", len(lics))
	}
	if lics[0] != "BSD-4-Clause" {
		t.Errorf("expected %v, got %v", "BSD-4-Clause", lics[0])
	}
	if lics[1] != "EPL-1.0" {
		t.Errorf("expected %v, got %v", "EPL-1.0", lics[1])
	}
	if lics[2] != "GPL-2.0-only" {
		t.Errorf("expected %v, got %v", "GPL-2.0-only", lics[2])
	}
	if lics[3] != "MIT" {
		t.Errorf("expected %v, got %v", "MIT", lics[3])
	}
}

func TestCanGetIndividualLicensesIgnoringOperatorCase(t *testing.T) {
	// two-license combo with lowercase 'and'
	lic := "ISC and BSD-3-Clause"
	lics := getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	// should be sorted alphabetically; 'and' should not appear
	if len(lics) != 2 {
		t.Fatalf("expected lics to have len 2, got %d", len(lics))
	}
	if lics[0] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", lics[0])
	}
	if lics[1] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", lics[1])
	}

	// two-license combo with lowercase 'or'
	lic = "ISC or BSD-3-Clause"
	lics = getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	// should be sorted alphabetically; 'or' should not appear
	if len(lics) != 2 {
		t.Fatalf("expected lics to have len 2, got %d", len(lics))
	}
	if lics[0] != "BSD-3-Clause" {
		t.Errorf("expected %v, got %v", "BSD-3-Clause", lics[0])
	}
	if lics[1] != "ISC" {
		t.Errorf("expected %v, got %v", "ISC", lics[1])
	}

	// two-license combo with lowercase 'with'
	lic = "GPL-2.0-only with Classpath-exception-2.0"
	lics = getIndividualLicenses(lic)
	if lics == nil {
		t.Fatalf("expected non-nil lics, got nil")
	}
	// should be sorted alphabetically; 'with' should not appear
	if len(lics) != 2 {
		t.Fatalf("expected lics to have len 2, got %d", len(lics))
	}
	if lics[0] != "Classpath-exception-2.0" {
		t.Errorf("expected %v, got %v", "Classpath-exception-2.0", lics[0])
	}
	if lics[1] != "GPL-2.0-only" {
		t.Errorf("expected %v, got %v", "GPL-2.0-only", lics[1])
	}

}
