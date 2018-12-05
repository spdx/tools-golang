// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package licensediff

import (
	"testing"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// ===== License diff top-level function tests =====
func TestDifferCanCreateDiffPairs(t *testing.T) {
	// create files to be used in diff
	// f1 will be identical in both
	f1 := &spdx.File2_1{
		FileName:           "/project/file1.txt",
		FileSPDXIdentifier: "SPDXRef-File561",
		FileChecksumSHA1:   "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
		LicenseConcluded:   "Apache-2.0",
		LicenseInfoInFile: []string{
			"LicenseRef-We-will-ignore-LicenseInfoInFile",
		},
		FileCopyrightText: "We'll ignore copyright values",
	}

	// f2 will only appear in the first Package
	f2 := &spdx.File2_1{
		FileName:           "/project/file2.txt",
		FileSPDXIdentifier: "SPDXRef-File562",
		FileChecksumSHA1:   "066c5139bd9a43d15812ec1a1755b08ccf199824",
		LicenseConcluded:   "GPL-2.0-or-later",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f3 will only appear in the second Package
	f3 := &spdx.File2_1{
		FileName:           "/project/file3.txt",
		FileSPDXIdentifier: "SPDXRef-File563",
		FileChecksumSHA1:   "bd0f4863b15fad2b79b35303af54fcb5baaf7c68",
		LicenseConcluded:   "MPL-2.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f4_1 and f4_2 will appear in first and second,
	// with same name, same hash and different license
	f4_1 := &spdx.File2_1{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: "SPDXRef-File564",
		FileChecksumSHA1:   "bc417a575ceae93435bcb7bfd382ac28cbdaa8b5",
		LicenseConcluded:   "MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f4_2 := &spdx.File2_1{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: "SPDXRef-File564",
		FileChecksumSHA1:   "bc417a575ceae93435bcb7bfd382ac28cbdaa8b5",
		LicenseConcluded:   "Apache-2.0 AND MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f5_1 and f5_2 will appear in first and second,
	// with same name, different hash and same license
	f5_1 := &spdx.File2_1{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: "SPDXRef-File565",
		FileChecksumSHA1:   "ba226db943bbbf455da77afab6f16dbab156d000",
		LicenseConcluded:   "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f5_2 := &spdx.File2_1{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: "SPDXRef-File565",
		FileChecksumSHA1:   "b6e0ec7d085c5699b46f6f8d425413702652874d",
		LicenseConcluded:   "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f6_1 and f6_2 will appear in first and second,
	// with same name, different hash and different license
	f6_1 := &spdx.File2_1{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: "SPDXRef-File566",
		FileChecksumSHA1:   "ba226db943bbbf455da77afab6f16dbab156d000",
		LicenseConcluded:   "CC0-1.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f6_2 := &spdx.File2_1{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: "SPDXRef-File566",
		FileChecksumSHA1:   "b6e0ec7d085c5699b46f6f8d425413702652874d",
		LicenseConcluded:   "Unlicense",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// create Packages
	p1 := &spdx.Package2_1{
		IsUnpackaged:              false,
		PackageName:               "p1",
		PackageSPDXIdentifier:     "SPDXRef-p1",
		PackageDownloadLocation:   "NOASSERTION",
		FilesAnalyzed:             true,
		IsFilesAnalyzedTagPresent: true,
		// fake the verification code for present purposes
		PackageVerificationCode: "abc123abc123",
		PackageLicenseConcluded: "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{
			"NOASSERTION",
		},
		PackageLicenseDeclared: "NOASSERTION",
		PackageCopyrightText:   "NOASSERTION",
		Files: []*spdx.File2_1{
			f1,
			f2,
			f4_1,
			f5_1,
			f6_1,
		},
	}
	p2 := &spdx.Package2_1{
		IsUnpackaged:              false,
		PackageName:               "p2",
		PackageSPDXIdentifier:     "SPDXRef-p2",
		PackageDownloadLocation:   "NOASSERTION",
		FilesAnalyzed:             true,
		IsFilesAnalyzedTagPresent: true,
		// fake the verification code for present purposes
		PackageVerificationCode: "def456def456",
		PackageLicenseConcluded: "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{
			"NOASSERTION",
		},
		PackageLicenseDeclared: "NOASSERTION",
		PackageCopyrightText:   "NOASSERTION",
		Files: []*spdx.File2_1{
			f1,
			f3,
			f4_2,
			f5_2,
			f6_2,
		},
	}

	// run the diff between the two packages
	diffMap, err := MakePairs(p1, p2)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	// check that the diff results are what we expect
	// there should be 6 entries, one for each unique filename
	if len(diffMap) != 6 {
		t.Fatalf("Expected %d, got %d", 6, len(diffMap))
	}

	// check each filename is present, and check its pair
	// pair 1 -- same in both
	pair1, ok := diffMap["/project/file1.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair1")
	}
	if pair1.first != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f1.LicenseConcluded, pair1.first)
	}
	if pair1.second != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, pair1.second)
	}

	// pair 2 -- only in first
	pair2, ok := diffMap["/project/file2.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair2")
	}
	if pair2.first != f2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, pair2.first)
	}
	if pair2.second != "" {
		t.Errorf("Expected %s, got %s", "", pair2.second)
	}

	// pair 3 -- only in second
	pair3, ok := diffMap["/project/file3.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair3")
	}
	if pair3.first != "" {
		t.Errorf("Expected %s, got %s", "", pair3.first)
	}
	if pair3.second != f3.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f3.LicenseConcluded, pair3.second)
	}

	// pair 4 -- in both but different license
	pair4, ok := diffMap["/project/file4.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair4")
	}
	if pair4.first != f4_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_1.LicenseConcluded, pair4.first)
	}
	if pair4.second != f4_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_2.LicenseConcluded, pair4.second)
	}

	// pair 5 -- in both but different hash, same license
	pair5, ok := diffMap["/project/file5.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair5")
	}
	if pair5.first != f5_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_1.LicenseConcluded, pair5.first)
	}
	if pair5.second != f5_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_2.LicenseConcluded, pair5.second)
	}

	// pair 6 -- in both but different hash, different license
	pair6, ok := diffMap["/project/file6.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair6")
	}
	if pair6.first != f6_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_1.LicenseConcluded, pair6.first)
	}
	if pair6.second != f6_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_2.LicenseConcluded, pair6.second)
	}
}
