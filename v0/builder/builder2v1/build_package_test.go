// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"testing"
)

// ===== Package section builder tests =====
func TestBuilder2_1CanBuildPackageSection(t *testing.T) {
	packageName := "project1"
	dirRoot := "../../../testdata/project1/"

	wantVerificationCode := "fc9ac4a370af0a471c2e52af66d6b4cf4e2ba12b"

	pkg, err := BuildPackageSection2_1(packageName, dirRoot, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if pkg == nil {
		t.Fatalf("expected non-nil Package, got nil")
	}
	if pkg.IsUnpackaged {
		t.Errorf("expected %v, got %v", false, pkg.IsUnpackaged)
	}
	if pkg.PackageName != "project1" {
		t.Errorf("expected %v, got %v", "project1", pkg.PackageName)
	}
	if pkg.PackageSPDXIdentifier != "SPDXRef-Package-project1" {
		t.Errorf("expected %v, got %v", "SPDXRef-Package-project1", pkg.PackageSPDXIdentifier)
	}
	if pkg.PackageDownloadLocation != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.PackageDownloadLocation)
	}
	if pkg.FilesAnalyzed != true {
		t.Errorf("expected %v, got %v", true, pkg.FilesAnalyzed)
	}
	if pkg.IsFilesAnalyzedTagPresent != true {
		t.Errorf("expected %v, got %v", true, pkg.IsFilesAnalyzedTagPresent)
	}
	if pkg.PackageVerificationCode != wantVerificationCode {
		t.Errorf("expected %v, got %v", wantVerificationCode, pkg.PackageVerificationCode)
	}
	if pkg.PackageLicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.PackageLicenseConcluded)
	}
	if len(pkg.PackageLicenseInfoFromFiles) != 0 {
		t.Errorf("expected %v, got %v", 0, len(pkg.PackageLicenseInfoFromFiles))
	}
	if pkg.PackageLicenseDeclared != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.PackageLicenseDeclared)
	}
	if pkg.PackageCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.PackageCopyrightText)
	}

	// and make sure we got the right number of files, and check the first one
	if pkg.Files == nil {
		t.Fatalf("expected non-nil pkg.Files, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected %d, got %d", 5, len(pkg.Files))
	}
	fileEmpty := pkg.Files[0]
	if fileEmpty == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if fileEmpty.FileName != "/emptyfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "/emptyfile.testdata.txt", fileEmpty.FileName)
	}
	if fileEmpty.FileSPDXIdentifier != "SPDXRef-File0" {
		t.Errorf("expected %v, got %v", "SPDXRef-File0", fileEmpty.FileSPDXIdentifier)
	}
	if fileEmpty.FileChecksumSHA1 != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Errorf("expected %v, got %v", "da39a3ee5e6b4b0d3255bfef95601890afd80709", fileEmpty.FileChecksumSHA1)
	}
	if fileEmpty.FileChecksumSHA256 != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Errorf("expected %v, got %v", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", fileEmpty.FileChecksumSHA256)
	}
	if fileEmpty.FileChecksumMD5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("expected %v, got %v", "d41d8cd98f00b204e9800998ecf8427e", fileEmpty.FileChecksumMD5)
	}
	if fileEmpty.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.LicenseConcluded)
	}
	if len(fileEmpty.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(fileEmpty.LicenseInfoInFile))
	}
	if fileEmpty.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.FileCopyrightText)
	}
}

func TestBuilder2_1CanIgnoreFiles(t *testing.T) {
	packageName := "project3"
	dirRoot := "../../../testdata/project3/"
	pathsIgnored := []string{
		"**/ignoredir/",
		"/excludedir/",
		"**/ignorefile.txt",
		"/alsoEXCLUDEthis.txt",
	}
	pkg, err := BuildPackageSection2_1(packageName, dirRoot, pathsIgnored)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// make sure we got the right files
	if pkg.Files == nil {
		t.Fatalf("expected non-nil pkg.Files, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected %d, got %d", 5, len(pkg.Files))
	}

	want := "/dontscan.txt"
	got := pkg.Files[0].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/keep/keep.txt"
	got = pkg.Files[1].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/keep.txt"
	got = pkg.Files[2].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/subdir/keep/dontscan.txt"
	got = pkg.Files[3].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/subdir/keep/keep.txt"
	got = pkg.Files[4].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}
