// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Package section builder tests =====
func TestBuilder2_2CanBuildPackageSection(t *testing.T) {
	packageName := "project1"
	dirRoot := "../../testdata/project1/"

	wantVerificationCode := "fc9ac4a370af0a471c2e52af66d6b4cf4e2ba12b"

	pkg, err := BuildPackageSection2_2(packageName, dirRoot, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if pkg == nil {
		t.Fatalf("expected non-nil Package, got nil")
	}
	if pkg.Name != "project1" {
		t.Errorf("expected %v, got %v", "project1", pkg.Name)
	}
	if pkg.SPDXIdentifier != spdx.ElementID("Package-project1") {
		t.Errorf("expected %v, got %v", "Package-project1", pkg.SPDXIdentifier)
	}
	if pkg.DownloadLocation != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.DownloadLocation)
	}
	if pkg.FilesAnalyzed != true {
		t.Errorf("expected %v, got %v", true, pkg.FilesAnalyzed)
	}
	if pkg.IsFilesAnalyzedTagPresent != true {
		t.Errorf("expected %v, got %v", true, pkg.IsFilesAnalyzedTagPresent)
	}
	if pkg.VerificationCode != wantVerificationCode {
		t.Errorf("expected %v, got %v", wantVerificationCode, pkg.VerificationCode)
	}
	if pkg.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.LicenseConcluded)
	}
	if len(pkg.LicenseInfoFromFiles) != 0 {
		t.Errorf("expected %v, got %v", 0, len(pkg.LicenseInfoFromFiles))
	}
	if pkg.LicenseDeclared != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.LicenseDeclared)
	}
	if pkg.CopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.CopyrightText)
	}

	// and make sure we got the right number of files, and check the first one
	if pkg.Files == nil {
		t.Fatalf("expected non-nil pkg.Files, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected %d, got %d", 5, len(pkg.Files))
	}
	fileEmpty := pkg.Files[spdx.ElementID("File0")]
	if fileEmpty == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if fileEmpty.Name != "/emptyfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "/emptyfile.testdata.txt", fileEmpty.Name)
	}
	if fileEmpty.SPDXIdentifier != spdx.ElementID("File0") {
		t.Errorf("expected %v, got %v", "File0", fileEmpty.SPDXIdentifier)
	}
	if fileEmpty.ChecksumSHA1 != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Errorf("expected %v, got %v", "da39a3ee5e6b4b0d3255bfef95601890afd80709", fileEmpty.ChecksumSHA1)
	}
	if fileEmpty.ChecksumSHA256 != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Errorf("expected %v, got %v", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", fileEmpty.ChecksumSHA256)
	}
	if fileEmpty.ChecksumMD5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("expected %v, got %v", "d41d8cd98f00b204e9800998ecf8427e", fileEmpty.ChecksumMD5)
	}
	if fileEmpty.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.LicenseConcluded)
	}
	if len(fileEmpty.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(fileEmpty.LicenseInfoInFile))
	}
	if fileEmpty.CopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.CopyrightText)
	}
}

func TestBuilder2_2CanIgnoreFiles(t *testing.T) {
	packageName := "project3"
	dirRoot := "../../testdata/project3/"
	pathsIgnored := []string{
		"**/ignoredir/",
		"/excludedir/",
		"**/ignorefile.txt",
		"/alsoEXCLUDEthis.txt",
	}
	pkg, err := BuildPackageSection2_2(packageName, dirRoot, pathsIgnored)
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
	got := pkg.Files[spdx.ElementID("File0")].Name
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/keep/keep.txt"
	got = pkg.Files[spdx.ElementID("File1")].Name
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/keep.txt"
	got = pkg.Files[spdx.ElementID("File2")].Name
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/subdir/keep/dontscan.txt"
	got = pkg.Files[spdx.ElementID("File3")].Name
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "/subdir/keep/keep.txt"
	got = pkg.Files[spdx.ElementID("File4")].Name
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}
