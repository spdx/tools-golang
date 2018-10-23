// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"fmt"
	"testing"
)

// ===== Builder top-level Document test =====
func TestBuild2_1CreatesDocument(t *testing.T) {
	dirRoot := "../../testdata/project1/"

	config := &Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
		CreatorType:     "Person",
		Creator:         "John Doe",
		TestValues:      make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-19T04:38:00Z"

	wantVerificationCode := "fc9ac4a370af0a471c2e52af66d6b4cf4e2ba12b"

	doc, err := Build2_1("project1", dirRoot, config)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil Document, got nil")
	}

	// check CI section
	if doc.CreationInfo == nil {
		t.Fatalf("expected non-nil CreationInfo section, got nil")
	}
	if doc.CreationInfo.SPDXVersion != "SPDX-2.1" {
		t.Errorf("expected %s, got %s", "SPDX-2.1", doc.CreationInfo.SPDXVersion)
	}
	if doc.CreationInfo.DataLicense != "CC0-1.0" {
		t.Errorf("expected %s, got %s", "CC0-1.0", doc.CreationInfo.DataLicense)
	}
	if doc.CreationInfo.SPDXIdentifier != "SPDXRef-DOCUMENT" {
		t.Errorf("expected %s, got %s", "SPDXRef-DOCUMENT", doc.CreationInfo.SPDXIdentifier)
	}
	if doc.CreationInfo.DocumentName != "project1" {
		t.Errorf("expected %s, got %s", "project1", doc.CreationInfo.DocumentName)
	}
	wantNamespace := fmt.Sprintf("https://github.com/swinslow/spdx-docs/spdx-go/testdata-project1-%s", wantVerificationCode)
	if doc.CreationInfo.DocumentNamespace != wantNamespace {
		t.Errorf("expected %s, got %s", wantNamespace, doc.CreationInfo.DocumentNamespace)
	}
	if len(doc.CreationInfo.CreatorPersons) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(doc.CreationInfo.CreatorPersons))
	}
	if doc.CreationInfo.CreatorPersons[0] != "John Doe" {
		t.Errorf("expected %s, got %s", "John Doe", doc.CreationInfo.CreatorPersons[0])
	}
	if len(doc.CreationInfo.CreatorTools) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(doc.CreationInfo.CreatorTools))
	}
	if doc.CreationInfo.CreatorTools[0] != "github.com/swinslow/spdx-go/v0/builder" {
		t.Errorf("expected %s, got %s", "github.com/swinslow/spdx-go/v0/builder", doc.CreationInfo.CreatorTools[0])
	}
	if doc.CreationInfo.Created != "2018-10-19T04:38:00Z" {
		t.Errorf("expected %s, got %s", "2018-10-19T04:38:00Z", doc.CreationInfo.Created)
	}

	// check Package section
	if doc.Packages == nil {
		t.Fatalf("expected non-nil doc.Packages, got nil")
	}
	if len(doc.Packages) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(doc.Packages))
	}
	pkg := doc.Packages[0]
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

	// check Files section
	if pkg.Files == nil {
		t.Fatalf("expected non-nil pkg.Files, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected %d, got %d", 5, len(pkg.Files))
	}

	// files should be in alphabetical order:
	// emptyfile, file1, file3, folder/file4, lastfile

	// check emptyfile.testdata.txt
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

	// check file1.testdata.txt
	file1 := pkg.Files[1]
	if file1 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file1.FileName != "/file1.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file1.testdata.txt", file1.FileName)
	}
	if file1.FileSPDXIdentifier != "SPDXRef-File1" {
		t.Errorf("expected %v, got %v", "SPDXRef-File1", file1.FileSPDXIdentifier)
	}
	if file1.FileChecksumSHA1 != "024f870eb6323f532515f7a09d5646a97083b819" {
		t.Errorf("expected %v, got %v", "024f870eb6323f532515f7a09d5646a97083b819", file1.FileChecksumSHA1)
	}
	if file1.FileChecksumSHA256 != "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf" {
		t.Errorf("expected %v, got %v", "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf", file1.FileChecksumSHA256)
	}
	if file1.FileChecksumMD5 != "37c8208479dfe42d2bb29debd6e32d4a" {
		t.Errorf("expected %v, got %v", "37c8208479dfe42d2bb29debd6e32d4a", file1.FileChecksumMD5)
	}
	if file1.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.LicenseConcluded)
	}
	if len(file1.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(file1.LicenseInfoInFile))
	}
	if file1.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.FileCopyrightText)
	}

	// check file3.testdata.txt
	file3 := pkg.Files[2]
	if file3 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file3.FileName != "/file3.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file3.testdata.txt", file3.FileName)
	}
	if file3.FileSPDXIdentifier != "SPDXRef-File2" {
		t.Errorf("expected %v, got %v", "SPDXRef-File2", file3.FileSPDXIdentifier)
	}
	if file3.FileChecksumSHA1 != "a46114b70e163614f01c64adf44cdd438f158fce" {
		t.Errorf("expected %v, got %v", "a46114b70e163614f01c64adf44cdd438f158fce", file3.FileChecksumSHA1)
	}
	if file3.FileChecksumSHA256 != "9fc181b9892720a15df1a1e561860318db40621bd4040ccdf18e110eb01d04b4" {
		t.Errorf("expected %v, got %v", "9fc181b9892720a15df1a1e561860318db40621bd4040ccdf18e110eb01d04b4", file3.FileChecksumSHA256)
	}
	if file3.FileChecksumMD5 != "3e02d3ab9c58eec6911dbba37570934f" {
		t.Errorf("expected %v, got %v", "3e02d3ab9c58eec6911dbba37570934f", file3.FileChecksumMD5)
	}
	if file3.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file3.LicenseConcluded)
	}
	if len(file3.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(file3.LicenseInfoInFile))
	}
	if file3.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file3.FileCopyrightText)
	}

	// check folder1/file4.testdata.txt
	file4 := pkg.Files[3]
	if file4 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file4.FileName != "/folder1/file4.testdata.txt" {
		t.Errorf("expected %v, got %v", "folder1/file4.testdata.txt", file4.FileName)
	}
	if file4.FileSPDXIdentifier != "SPDXRef-File3" {
		t.Errorf("expected %v, got %v", "SPDXRef-File3", file4.FileSPDXIdentifier)
	}
	if file4.FileChecksumSHA1 != "e623d7d7d782a7c8323c4d436acee4afab34320f" {
		t.Errorf("expected %v, got %v", "e623d7d7d782a7c8323c4d436acee4afab34320f", file4.FileChecksumSHA1)
	}
	if file4.FileChecksumSHA256 != "574fa42c5e0806c0f8906a44884166540206f021527729407cd5326838629c59" {
		t.Errorf("expected %v, got %v", "574fa42c5e0806c0f8906a44884166540206f021527729407cd5326838629c59", file4.FileChecksumSHA256)
	}
	if file4.FileChecksumMD5 != "96e6a25d35df5b1c477710ef4d0c7210" {
		t.Errorf("expected %v, got %v", "96e6a25d35df5b1c477710ef4d0c7210", file4.FileChecksumMD5)
	}
	if file4.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file4.LicenseConcluded)
	}
	if len(file4.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(file4.LicenseInfoInFile))
	}
	if file4.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file4.FileCopyrightText)
	}

	// check lastfile.testdata.txt
	lastfile := pkg.Files[4]
	if lastfile == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if lastfile.FileName != "/lastfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "/lastfile.testdata.txt", lastfile.FileName)
	}
	if lastfile.FileSPDXIdentifier != "SPDXRef-File4" {
		t.Errorf("expected %v, got %v", "SPDXRef-File4", lastfile.FileSPDXIdentifier)
	}
	if lastfile.FileChecksumSHA1 != "26d6221d682d9ba59116f9753a701f34271c8ce1" {
		t.Errorf("expected %v, got %v", "26d6221d682d9ba59116f9753a701f34271c8ce1", lastfile.FileChecksumSHA1)
	}
	if lastfile.FileChecksumSHA256 != "0a4bdaf990e9b330ff72022dd78110ae98b60e08337cf2105b89856373416805" {
		t.Errorf("expected %v, got %v", "0a4bdaf990e9b330ff72022dd78110ae98b60e08337cf2105b89856373416805", lastfile.FileChecksumSHA256)
	}
	if lastfile.FileChecksumMD5 != "f60baa793870d9085461ad6bbab50b7f" {
		t.Errorf("expected %v, got %v", "f60baa793870d9085461ad6bbab50b7f", lastfile.FileChecksumMD5)
	}
	if lastfile.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", lastfile.LicenseConcluded)
	}
	if len(lastfile.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(lastfile.LicenseInfoInFile))
	}
	if lastfile.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", lastfile.FileCopyrightText)
	}

	// check Relationship section -- should be a relationship for doc DESCRIBES pkg
	if doc.Relationships == nil {
		t.Fatalf("expected non-nil Relationships section, got nil")
	}
	if len(doc.Relationships) == 0 {
		t.Fatalf("expected %v, got %v", 0, len(doc.Relationships))
	}
	rln := doc.Relationships[0]
	if rln == nil {
		t.Fatalf("expected non-nil Relationship, got nil")
	}
	if rln.RefA != "SPDXRef-DOCUMENT" {
		t.Errorf("expected %v, got %v", "SPDXRef-DOCUMENT", rln.RefA)
	}
	if rln.RefB != "SPDXRef-Package-project1" {
		t.Errorf("expected %v, got %v", "SPDXRef-Package-project1", rln.RefB)
	}
	if rln.Relationship != "DESCRIBES" {
		t.Errorf("expected %v, got %v", "DESCRIBES", rln.Relationship)
	}

	// and check that other sections are present, but empty
	if doc.OtherLicenses != nil {
		t.Fatalf("expected nil OtherLicenses section, got non-nil")
	}
	if doc.Annotations != nil {
		t.Fatalf("expected nil Annotations section, got non-nil")
	}
	if doc.Reviews != nil {
		t.Fatalf("expected nil Reviews section, got non-nil")
	}

}

func TestBuild2_1CanIgnoreFiles(t *testing.T) {
	dirRoot := "../../testdata/project3/"

	config := &Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
		CreatorType:     "Person",
		Creator:         "John Doe",
		PathsIgnored: []string{
			"**/ignoredir/",
			"/excludedir/",
			"**/ignorefile.txt",
			"/alsoEXCLUDEthis.txt",
		},
		TestValues: make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-19T04:38:00Z"

	doc, err := Build2_1("project1", dirRoot, config)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	pkg := doc.Packages[0]
	if len(pkg.Files) != 5 {
		t.Fatalf("expected len %d, got %d", 5, len(pkg.Files))
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
