// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"fmt"
	"testing"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

func TestBuildCreatesDocument(t *testing.T) {
	dirRoot := "../testdata/project1/"

	config := &Config{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-",
		CreatorType:     "Person",
		Creator:         "John Doe",
		TestValues:      make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-19T04:38:00Z"

	wantVerificationCode := common.PackageVerificationCode{Value: "fc9ac4a370af0a471c2e52af66d6b4cf4e2ba12b"}

	doc, err := Build("project1", dirRoot, config)
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
	if doc.SPDXVersion != spdx.Version {
		t.Errorf("expected %s, got %s", spdx.Version, doc.SPDXVersion)
	}
	if doc.DataLicense != spdx.DataLicense {
		t.Errorf("expected %s, got %s", spdx.DataLicense, doc.DataLicense)
	}
	if doc.SPDXIdentifier != common.ElementID("DOCUMENT") {
		t.Errorf("expected %s, got %v", "DOCUMENT", doc.SPDXIdentifier)
	}
	if doc.DocumentName != "project1" {
		t.Errorf("expected %s, got %s", "project1", doc.DocumentName)
	}
	wantNamespace := fmt.Sprintf("https://github.com/swinslow/spdx-docs/spdx-go/testdata-project1-%s", wantVerificationCode)
	if doc.DocumentNamespace != wantNamespace {
		t.Errorf("expected %s, got %s", wantNamespace, doc.DocumentNamespace)
	}
	if len(doc.CreationInfo.Creators) != 2 {
		t.Fatalf("expected %d, got %d", 2, len(doc.CreationInfo.Creators))
	}
	if doc.CreationInfo.Creators[1].Creator != "John Doe" {
		t.Errorf("expected %s, got %+v", "John Doe", doc.CreationInfo.Creators[1])
	}
	if doc.CreationInfo.Creators[0].Creator != "github.com/spdx/tools-golang/builder" {
		t.Errorf("expected %s, got %+v", "github.com/spdx/tools-golang/builder", doc.CreationInfo.Creators[0])
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
	if pkg == nil {
		t.Fatalf("expected non-nil pkg, got nil")
	}
	if pkg.PackageName != "project1" {
		t.Errorf("expected %v, got %v", "project1", pkg.PackageName)
	}
	if pkg.PackageSPDXIdentifier != common.ElementID("Package-project1") {
		t.Errorf("expected %v, got %v", "Package-project1", pkg.PackageSPDXIdentifier)
	}
	if pkg.PackageDownloadLocation != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", pkg.PackageDownloadLocation)
	}
	if pkg.FilesAnalyzed != true {
		t.Errorf("expected %v, got %v", true, pkg.FilesAnalyzed)
	}
	if pkg.PackageVerificationCode.Value != wantVerificationCode.Value {
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

	// files should be in order of identifier, which is numeric,
	// created based on alphabetical order of files:
	// emptyfile, file1, file3, folder/file4, lastfile

	// check emptyfile.testdata.txt
	fileEmpty := pkg.Files[0]
	if fileEmpty == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if fileEmpty.FileName != "./emptyfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "./emptyfile.testdata.txt", fileEmpty.FileName)
	}
	if fileEmpty.FileSPDXIdentifier != common.ElementID("File0") {
		t.Errorf("expected %v, got %v", "File0", fileEmpty.FileSPDXIdentifier)
	}
	for _, checksum := range fileEmpty.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
				t.Errorf("expected %v, got %v", "da39a3ee5e6b4b0d3255bfef95601890afd80709", checksum.Value)
			}
		case common.SHA256:
			if checksum.Value != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
				t.Errorf("expected %v, got %v", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", checksum.Value)
			}
		case common.MD5:
			if checksum.Value != "d41d8cd98f00b204e9800998ecf8427e" {
				t.Errorf("expected %v, got %v", "d41d8cd98f00b204e9800998ecf8427e", checksum.Value)
			}
		}
	}
	if fileEmpty.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.LicenseConcluded)
	}
	if len(fileEmpty.LicenseInfoInFiles) != 1 {
		t.Errorf("expected %v, got %v", 1, len(fileEmpty.LicenseInfoInFiles))
	} else {
		if fileEmpty.LicenseInfoInFiles[0] != "NOASSERTION" {
			t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.LicenseInfoInFiles[0])
		}
	}
	if fileEmpty.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", fileEmpty.FileCopyrightText)
	}

	// check file1.testdata.txt
	file1 := pkg.Files[1]
	if file1 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file1.FileName != "./file1.testdata.txt" {
		t.Errorf("expected %v, got %v", "./file1.testdata.txt", file1.FileName)
	}
	if file1.FileSPDXIdentifier != common.ElementID("File1") {
		t.Errorf("expected %v, got %v", "File1", file1.FileSPDXIdentifier)
	}
	for _, checksum := range file1.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != "024f870eb6323f532515f7a09d5646a97083b819" {
				t.Errorf("expected %v, got %v", "024f870eb6323f532515f7a09d5646a97083b819", checksum.Value)
			}
		case common.SHA256:
			if checksum.Value != "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf" {
				t.Errorf("expected %v, got %v", "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf", checksum.Value)
			}
		case common.MD5:
			if checksum.Value != "37c8208479dfe42d2bb29debd6e32d4a" {
				t.Errorf("expected %v, got %v", "37c8208479dfe42d2bb29debd6e32d4a", checksum.Value)
			}
		}
	}
	if file1.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.LicenseConcluded)
	}
	if len(file1.LicenseInfoInFiles) != 1 {
		t.Errorf("expected %v, got %v", 1, len(file1.LicenseInfoInFiles))
	} else {
		if file1.LicenseInfoInFiles[0] != "NOASSERTION" {
			t.Errorf("expected %v, got %v", "NOASSERTION", file1.LicenseInfoInFiles[0])
		}
	}
	if file1.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.FileCopyrightText)
	}

	// check file3.testdata.txt
	file3 := pkg.Files[2]
	if file3 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file3.FileName != "./file3.testdata.txt" {
		t.Errorf("expected %v, got %v", "./file3.testdata.txt", file3.FileName)
	}
	if file3.FileSPDXIdentifier != common.ElementID("File2") {
		t.Errorf("expected %v, got %v", "File2", file3.FileSPDXIdentifier)
	}
	for _, checksum := range file3.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != "a46114b70e163614f01c64adf44cdd438f158fce" {
				t.Errorf("expected %v, got %v", "a46114b70e163614f01c64adf44cdd438f158fce", checksum.Value)
			}
		case common.SHA256:
			if checksum.Value != "9fc181b9892720a15df1a1e561860318db40621bd4040ccdf18e110eb01d04b4" {
				t.Errorf("expected %v, got %v", "9fc181b9892720a15df1a1e561860318db40621bd4040ccdf18e110eb01d04b4", checksum.Value)
			}
		case common.MD5:
			if checksum.Value != "3e02d3ab9c58eec6911dbba37570934f" {
				t.Errorf("expected %v, got %v", "3e02d3ab9c58eec6911dbba37570934f", checksum.Value)
			}
		}
	}
	if file3.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file3.LicenseConcluded)
	}
	if len(file3.LicenseInfoInFiles) != 1 {
		t.Errorf("expected %v, got %v", 1, len(file3.LicenseInfoInFiles))
	} else {
		if file3.LicenseInfoInFiles[0] != "NOASSERTION" {
			t.Errorf("expected %v, got %v", "NOASSERTION", file3.LicenseInfoInFiles[0])
		}
	}
	if file3.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file3.FileCopyrightText)
	}

	// check folder1/file4.testdata.txt
	file4 := pkg.Files[3]
	if file4 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file4.FileName != "./folder1/file4.testdata.txt" {
		t.Errorf("expected %v, got %v", "./folder1/file4.testdata.txt", file4.FileName)
	}
	if file4.FileSPDXIdentifier != common.ElementID("File3") {
		t.Errorf("expected %v, got %v", "File3", file4.FileSPDXIdentifier)
	}
	for _, checksum := range file4.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != "e623d7d7d782a7c8323c4d436acee4afab34320f" {
				t.Errorf("expected %v, got %v", "e623d7d7d782a7c8323c4d436acee4afab34320f", checksum.Value)
			}
		case common.SHA256:
			if checksum.Value != "574fa42c5e0806c0f8906a44884166540206f021527729407cd5326838629c59" {
				t.Errorf("expected %v, got %v", "574fa42c5e0806c0f8906a44884166540206f021527729407cd5326838629c59", checksum.Value)
			}
		case common.MD5:
			if checksum.Value != "96e6a25d35df5b1c477710ef4d0c7210" {
				t.Errorf("expected %v, got %v", "96e6a25d35df5b1c477710ef4d0c7210", checksum.Value)
			}
		}
	}
	if file4.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file4.LicenseConcluded)
	}
	if len(file4.LicenseInfoInFiles) != 1 {
		t.Errorf("expected %v, got %v", 1, len(file4.LicenseInfoInFiles))
	} else {
		if file4.LicenseInfoInFiles[0] != "NOASSERTION" {
			t.Errorf("expected %v, got %v", "NOASSERTION", file4.LicenseInfoInFiles[0])
		}
	}
	if file4.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file4.FileCopyrightText)
	}

	// check lastfile.testdata.txt
	lastfile := pkg.Files[4]
	if lastfile == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if lastfile.FileName != "./lastfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "/lastfile.testdata.txt", lastfile.FileName)
	}
	if lastfile.FileSPDXIdentifier != common.ElementID("File4") {
		t.Errorf("expected %v, got %v", "File4", lastfile.FileSPDXIdentifier)
	}
	for _, checksum := range lastfile.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != "26d6221d682d9ba59116f9753a701f34271c8ce1" {
				t.Errorf("expected %v, got %v", "26d6221d682d9ba59116f9753a701f34271c8ce1", checksum.Value)
			}
		case common.SHA256:
			if checksum.Value != "0a4bdaf990e9b330ff72022dd78110ae98b60e08337cf2105b89856373416805" {
				t.Errorf("expected %v, got %v", "0a4bdaf990e9b330ff72022dd78110ae98b60e08337cf2105b89856373416805", checksum.Value)
			}
		case common.MD5:
			if checksum.Value != "f60baa793870d9085461ad6bbab50b7f" {
				t.Errorf("expected %v, got %v", "f60baa793870d9085461ad6bbab50b7f", checksum.Value)
			}
		}
	}
	if lastfile.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", lastfile.LicenseConcluded)
	}
	if len(lastfile.LicenseInfoInFiles) != 1 {
		t.Errorf("expected %v, got %v", 1, len(lastfile.LicenseInfoInFiles))
	} else {
		if lastfile.LicenseInfoInFiles[0] != "NOASSERTION" {
			t.Errorf("expected %v, got %v", "NOASSERTION", lastfile.LicenseInfoInFiles[0])
		}
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
	if rln.RefA != common.MakeDocElementID("", "DOCUMENT") {
		t.Errorf("expected %v, got %v", "DOCUMENT", rln.RefA)
	}
	if rln.RefB != common.MakeDocElementID("", "Package-project1") {
		t.Errorf("expected %v, got %v", "Package-project1", rln.RefB)
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

func TestBuildCanIgnoreFiles(t *testing.T) {
	dirRoot := "../testdata/project3/"

	config := &Config{
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

	doc, err := Build("project1", dirRoot, config)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	pkg := doc.Packages[0]
	if pkg == nil {
		t.Fatalf("expected non-nil pkg, got nil")
	}
	if len(pkg.Files) != 5 {
		t.Fatalf("expected len %d, got %d", 5, len(pkg.Files))
	}

	want := "./dontscan.txt"
	got := pkg.Files[0].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "./keep/keep.txt"
	got = pkg.Files[1].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "./keep.txt"
	got = pkg.Files[2].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "./subdir/keep/dontscan.txt"
	got = pkg.Files[3].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "./subdir/keep/keep.txt"
	got = pkg.Files[4].FileName
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}
