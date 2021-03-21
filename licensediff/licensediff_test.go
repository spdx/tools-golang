// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package licensediff

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== 2.1 License diff top-level function tests =====
func Test2_1DifferCanCreateDiffPairs(t *testing.T) {
	// create files to be used in diff
	// f1 will be identical in both
	f1 := &spdx.File2_1{
		FileName:           "/project/file1.txt",
		FileSPDXIdentifier: spdx.ElementID("File561"),
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
		FileSPDXIdentifier: spdx.ElementID("File562"),
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
		FileSPDXIdentifier: spdx.ElementID("File563"),
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
		FileSPDXIdentifier: spdx.ElementID("File564"),
		FileChecksumSHA1:   "bc417a575ceae93435bcb7bfd382ac28cbdaa8b5",
		LicenseConcluded:   "MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f4_2 := &spdx.File2_1{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: spdx.ElementID("File564"),
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
		FileSPDXIdentifier: spdx.ElementID("File565"),
		FileChecksumSHA1:   "ba226db943bbbf455da77afab6f16dbab156d000",
		LicenseConcluded:   "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f5_2 := &spdx.File2_1{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: spdx.ElementID("File565"),
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
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksumSHA1:   "ba226db943bbbf455da77afab6f16dbab156d000",
		LicenseConcluded:   "CC0-1.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f6_2 := &spdx.File2_1{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksumSHA1:   "b6e0ec7d085c5699b46f6f8d425413702652874d",
		LicenseConcluded:   "Unlicense",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// create Packages
	p1 := &spdx.Package2_1{
		PackageName:               "p1",
		PackageSPDXIdentifier:     spdx.ElementID("p1"),
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
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f2.FileSPDXIdentifier):   f2,
			spdx.ElementID(f4_1.FileSPDXIdentifier): f4_1,
			spdx.ElementID(f5_1.FileSPDXIdentifier): f5_1,
			spdx.ElementID(f6_1.FileSPDXIdentifier): f6_1,
		},
	}
	p2 := &spdx.Package2_1{
		PackageName:               "p2",
		PackageSPDXIdentifier:     spdx.ElementID("p2"),
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
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f3.FileSPDXIdentifier):   f3,
			spdx.ElementID(f4_2.FileSPDXIdentifier): f4_2,
			spdx.ElementID(f5_2.FileSPDXIdentifier): f5_2,
			spdx.ElementID(f6_2.FileSPDXIdentifier): f6_2,
		},
	}

	// run the diff between the two packages
	diffMap, err := MakePairs2_1(p1, p2)
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
	if pair1.First != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f1.LicenseConcluded, pair1.First)
	}
	if pair1.Second != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, pair1.Second)
	}

	// pair 2 -- only in first
	pair2, ok := diffMap["/project/file2.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair2")
	}
	if pair2.First != f2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, pair2.First)
	}
	if pair2.Second != "" {
		t.Errorf("Expected %s, got %s", "", pair2.Second)
	}

	// pair 3 -- only in second
	pair3, ok := diffMap["/project/file3.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair3")
	}
	if pair3.First != "" {
		t.Errorf("Expected %s, got %s", "", pair3.First)
	}
	if pair3.Second != f3.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f3.LicenseConcluded, pair3.Second)
	}

	// pair 4 -- in both but different license
	pair4, ok := diffMap["/project/file4.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair4")
	}
	if pair4.First != f4_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_1.LicenseConcluded, pair4.First)
	}
	if pair4.Second != f4_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_2.LicenseConcluded, pair4.Second)
	}

	// pair 5 -- in both but different hash, same license
	pair5, ok := diffMap["/project/file5.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair5")
	}
	if pair5.First != f5_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_1.LicenseConcluded, pair5.First)
	}
	if pair5.Second != f5_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_2.LicenseConcluded, pair5.Second)
	}

	// pair 6 -- in both but different hash, different license
	pair6, ok := diffMap["/project/file6.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair6")
	}
	if pair6.First != f6_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_1.LicenseConcluded, pair6.First)
	}
	if pair6.Second != f6_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_2.LicenseConcluded, pair6.Second)
	}
}

func Test2_1DifferCanCreateDiffStructuredResults(t *testing.T) {
	// create files to be used in diff
	// f1 will be identical in both
	f1 := &spdx.File2_1{
		FileName:           "/project/file1.txt",
		FileSPDXIdentifier: spdx.ElementID("File561"),
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
		FileSPDXIdentifier: spdx.ElementID("File562"),
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
		FileSPDXIdentifier: spdx.ElementID("File563"),
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
		FileSPDXIdentifier: spdx.ElementID("File564"),
		FileChecksumSHA1:   "bc417a575ceae93435bcb7bfd382ac28cbdaa8b5",
		LicenseConcluded:   "MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f4_2 := &spdx.File2_1{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: spdx.ElementID("File564"),
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
		FileSPDXIdentifier: spdx.ElementID("File565"),
		FileChecksumSHA1:   "ba226db943bbbf455da77afab6f16dbab156d000",
		LicenseConcluded:   "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f5_2 := &spdx.File2_1{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: spdx.ElementID("File565"),
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
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksumSHA1:   "ba226db943bbbf455da77afab6f16dbab156d000",
		LicenseConcluded:   "CC0-1.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f6_2 := &spdx.File2_1{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksumSHA1:   "b6e0ec7d085c5699b46f6f8d425413702652874d",
		LicenseConcluded:   "Unlicense",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// create Packages
	p1 := &spdx.Package2_1{
		PackageName:               "p1",
		PackageSPDXIdentifier:     spdx.ElementID("p1"),
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
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f2.FileSPDXIdentifier):   f2,
			spdx.ElementID(f4_1.FileSPDXIdentifier): f4_1,
			spdx.ElementID(f5_1.FileSPDXIdentifier): f5_1,
			spdx.ElementID(f6_1.FileSPDXIdentifier): f6_1,
		},
	}
	p2 := &spdx.Package2_1{
		PackageName:               "p2",
		PackageSPDXIdentifier:     spdx.ElementID("p2"),
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
		Files: map[spdx.ElementID]*spdx.File2_1{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f3.FileSPDXIdentifier):   f3,
			spdx.ElementID(f4_2.FileSPDXIdentifier): f4_2,
			spdx.ElementID(f5_2.FileSPDXIdentifier): f5_2,
			spdx.ElementID(f6_2.FileSPDXIdentifier): f6_2,
		},
	}

	// run the diff between the two packages
	diffMap, err := MakePairs2_1(p1, p2)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	// now, create the LicenseDiff structured results from the pairs
	diffResults, err := MakeResults(diffMap)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	// check that the diff results are the expected lengths
	if len(diffResults.InBothChanged) != 2 {
		t.Fatalf("Expected %d, got %d", 2, len(diffResults.InBothChanged))
	}
	if len(diffResults.InBothSame) != 2 {
		t.Fatalf("Expected %d, got %d", 2, len(diffResults.InBothSame))
	}
	if len(diffResults.InFirstOnly) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(diffResults.InFirstOnly))
	}
	if len(diffResults.InSecondOnly) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(diffResults.InSecondOnly))
	}

	// check each filename is present where it belongs, and check license(s)

	// in both and different license: f4 and f6
	// filename will map to a LicensePair
	check4, ok := diffResults.InBothChanged["/project/file4.txt"]
	if !ok {
		t.Fatalf("Couldn't get check4")
	}
	if check4.First != f4_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_1.LicenseConcluded, check4.First)
	}
	if check4.Second != f4_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_2.LicenseConcluded, check4.Second)
	}
	check6, ok := diffResults.InBothChanged["/project/file6.txt"]
	if !ok {
		t.Fatalf("Couldn't get check6")
	}
	if check6.First != f6_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_1.LicenseConcluded, check6.First)
	}
	if check6.Second != f6_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_2.LicenseConcluded, check6.Second)
	}

	// in both and same license: f1 and f5
	// filename will map to a string
	check1, ok := diffResults.InBothSame["/project/file1.txt"]
	if !ok {
		t.Fatalf("Couldn't get check1")
	}
	if check1 != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f1.LicenseConcluded, check1)
	}
	check5, ok := diffResults.InBothSame["/project/file5.txt"]
	if !ok {
		t.Fatalf("Couldn't get check5")
	}
	if check5 != f5_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_1.LicenseConcluded, check5)
	}
	if check5 != f5_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_2.LicenseConcluded, check5)
	}

	// in first only: f2
	// filename will map to a string
	check2, ok := diffResults.InFirstOnly["/project/file2.txt"]
	if !ok {
		t.Fatalf("Couldn't get check2")
	}
	if check2 != f2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, check2)
	}

	// in second only: f3
	// filename will map to a string
	check3, ok := diffResults.InSecondOnly["/project/file3.txt"]
	if !ok {
		t.Fatalf("Couldn't get check3")
	}
	if check3 != f3.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f3.LicenseConcluded, check2)
	}

}

// ===== 2.2 License diff top-level function tests =====
func Test2_2DifferCanCreateDiffPairs(t *testing.T) {
	// create files to be used in diff
	// f1 will be identical in both
	f1 := &spdx.File2_2{
		FileName:           "/project/file1.txt",
		FileSPDXIdentifier: spdx.ElementID("File561"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFile: []string{
			"LicenseRef-We-will-ignore-LicenseInfoInFile",
		},
		FileCopyrightText: "We'll ignore copyright values",
	}

	// f2 will only appear in the first Package
	f2 := &spdx.File2_2{
		FileName:           "/project/file2.txt",
		FileSPDXIdentifier: spdx.ElementID("File562"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "GPL-2.0-or-later",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f3 will only appear in the second Package
	f3 := &spdx.File2_2{
		FileName:           "/project/file3.txt",
		FileSPDXIdentifier: spdx.ElementID("File563"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "MPL-2.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f4_1 and f4_2 will appear in first and second,
	// with same name, same hash and different license
	f4_1 := &spdx.File2_2{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: spdx.ElementID("File564"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f4_2 := &spdx.File2_2{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: spdx.ElementID("File564"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "Apache-2.0 AND MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f5_1 and f5_2 will appear in first and second,
	// with same name, different hash and same license
	f5_1 := &spdx.File2_2{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: spdx.ElementID("File565"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f5_2 := &spdx.File2_2{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: spdx.ElementID("File565"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f6_1 and f6_2 will appear in first and second,
	// with same name, different hash and different license
	f6_1 := &spdx.File2_2{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "CC0-1.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f6_2 := &spdx.File2_2{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "Unlicense",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// create Packages
	p1 := &spdx.Package2_2{
		PackageName:               "p1",
		PackageSPDXIdentifier:     spdx.ElementID("p1"),
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
		Files: map[spdx.ElementID]*spdx.File2_2{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f2.FileSPDXIdentifier):   f2,
			spdx.ElementID(f4_1.FileSPDXIdentifier): f4_1,
			spdx.ElementID(f5_1.FileSPDXIdentifier): f5_1,
			spdx.ElementID(f6_1.FileSPDXIdentifier): f6_1,
		},
	}
	p2 := &spdx.Package2_2{
		PackageName:               "p2",
		PackageSPDXIdentifier:     spdx.ElementID("p2"),
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
		Files: map[spdx.ElementID]*spdx.File2_2{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f3.FileSPDXIdentifier):   f3,
			spdx.ElementID(f4_2.FileSPDXIdentifier): f4_2,
			spdx.ElementID(f5_2.FileSPDXIdentifier): f5_2,
			spdx.ElementID(f6_2.FileSPDXIdentifier): f6_2,
		},
	}

	// run the diff between the two packages
	diffMap, err := MakePairs2_2(p1, p2)
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
	if pair1.First != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f1.LicenseConcluded, pair1.First)
	}
	if pair1.Second != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, pair1.Second)
	}

	// pair 2 -- only in first
	pair2, ok := diffMap["/project/file2.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair2")
	}
	if pair2.First != f2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, pair2.First)
	}
	if pair2.Second != "" {
		t.Errorf("Expected %s, got %s", "", pair2.Second)
	}

	// pair 3 -- only in second
	pair3, ok := diffMap["/project/file3.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair3")
	}
	if pair3.First != "" {
		t.Errorf("Expected %s, got %s", "", pair3.First)
	}
	if pair3.Second != f3.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f3.LicenseConcluded, pair3.Second)
	}

	// pair 4 -- in both but different license
	pair4, ok := diffMap["/project/file4.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair4")
	}
	if pair4.First != f4_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_1.LicenseConcluded, pair4.First)
	}
	if pair4.Second != f4_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_2.LicenseConcluded, pair4.Second)
	}

	// pair 5 -- in both but different hash, same license
	pair5, ok := diffMap["/project/file5.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair5")
	}
	if pair5.First != f5_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_1.LicenseConcluded, pair5.First)
	}
	if pair5.Second != f5_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_2.LicenseConcluded, pair5.Second)
	}

	// pair 6 -- in both but different hash, different license
	pair6, ok := diffMap["/project/file6.txt"]
	if !ok {
		t.Fatalf("Couldn't get pair6")
	}
	if pair6.First != f6_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_1.LicenseConcluded, pair6.First)
	}
	if pair6.Second != f6_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_2.LicenseConcluded, pair6.Second)
	}
}

func Test2_2DifferCanCreateDiffStructuredResults(t *testing.T) {
	// create files to be used in diff
	// f1 will be identical in both
	f1 := &spdx.File2_2{
		FileName:           "/project/file1.txt",
		FileSPDXIdentifier: spdx.ElementID("File561"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "Apache-2.0",
		LicenseInfoInFile: []string{
			"LicenseRef-We-will-ignore-LicenseInfoInFile",
		},
		FileCopyrightText: "We'll ignore copyright values",
	}

	// f2 will only appear in the first Package
	f2 := &spdx.File2_2{
		FileName:           "/project/file2.txt",
		FileSPDXIdentifier: spdx.ElementID("File562"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "GPL-2.0-or-later",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f3 will only appear in the second Package
	f3 := &spdx.File2_2{
		FileName:           "/project/file3.txt",
		FileSPDXIdentifier: spdx.ElementID("File563"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "MPL-2.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f4_1 and f4_2 will appear in first and second,
	// with same name, same hash and different license
	f4_1 := &spdx.File2_2{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: spdx.ElementID("File564"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f4_2 := &spdx.File2_2{
		FileName:           "/project/file4.txt",
		FileSPDXIdentifier: spdx.ElementID("File564"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "Apache-2.0 AND MIT",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f5_1 and f5_2 will appear in first and second,
	// with same name, different hash and same license
	f5_1 := &spdx.File2_2{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: spdx.ElementID("File565"),
		LicenseConcluded:   "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f5_2 := &spdx.File2_2{
		FileName:           "/project/file5.txt",
		FileSPDXIdentifier: spdx.ElementID("File565"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},

		LicenseConcluded: "BSD-3-Clause",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// f6_1 and f6_2 will appear in first and second,
	// with same name, different hash and different license
	f6_1 := &spdx.File2_2{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "CC0-1.0",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}
	f6_2 := &spdx.File2_2{
		FileName:           "/project/file6.txt",
		FileSPDXIdentifier: spdx.ElementID("File566"),
		FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
			spdx.SHA1: spdx.Checksum{
				Algorithm: spdx.SHA1,
				Value:     "6c92dc8bc462b6889d9b1c0bc16c54d19a2cbdd3",
			},
		},
		LicenseConcluded: "Unlicense",
		LicenseInfoInFile: []string{
			"NOASSERTION",
		},
		FileCopyrightText: "NOASSERTION",
	}

	// create Packages
	p1 := &spdx.Package2_2{
		PackageName:               "p1",
		PackageSPDXIdentifier:     spdx.ElementID("p1"),
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
		Files: map[spdx.ElementID]*spdx.File2_2{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f2.FileSPDXIdentifier):   f2,
			spdx.ElementID(f4_1.FileSPDXIdentifier): f4_1,
			spdx.ElementID(f5_1.FileSPDXIdentifier): f5_1,
			spdx.ElementID(f6_1.FileSPDXIdentifier): f6_1,
		},
	}
	p2 := &spdx.Package2_2{
		PackageName:               "p2",
		PackageSPDXIdentifier:     spdx.ElementID("p2"),
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
		Files: map[spdx.ElementID]*spdx.File2_2{
			spdx.ElementID(f1.FileSPDXIdentifier):   f1,
			spdx.ElementID(f3.FileSPDXIdentifier):   f3,
			spdx.ElementID(f4_2.FileSPDXIdentifier): f4_2,
			spdx.ElementID(f5_2.FileSPDXIdentifier): f5_2,
			spdx.ElementID(f6_2.FileSPDXIdentifier): f6_2,
		},
	}

	// run the diff between the two packages
	diffMap, err := MakePairs2_2(p1, p2)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	// now, create the LicenseDiff structured results from the pairs
	diffResults, err := MakeResults(diffMap)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	// check that the diff results are the expected lengths
	if len(diffResults.InBothChanged) != 2 {
		t.Fatalf("Expected %d, got %d", 2, len(diffResults.InBothChanged))
	}
	if len(diffResults.InBothSame) != 2 {
		t.Fatalf("Expected %d, got %d", 2, len(diffResults.InBothSame))
	}
	if len(diffResults.InFirstOnly) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(diffResults.InFirstOnly))
	}
	if len(diffResults.InSecondOnly) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(diffResults.InSecondOnly))
	}

	// check each filename is present where it belongs, and check license(s)

	// in both and different license: f4 and f6
	// filename will map to a LicensePair
	check4, ok := diffResults.InBothChanged["/project/file4.txt"]
	if !ok {
		t.Fatalf("Couldn't get check4")
	}
	if check4.First != f4_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_1.LicenseConcluded, check4.First)
	}
	if check4.Second != f4_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f4_2.LicenseConcluded, check4.Second)
	}
	check6, ok := diffResults.InBothChanged["/project/file6.txt"]
	if !ok {
		t.Fatalf("Couldn't get check6")
	}
	if check6.First != f6_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_1.LicenseConcluded, check6.First)
	}
	if check6.Second != f6_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f6_2.LicenseConcluded, check6.Second)
	}

	// in both and same license: f1 and f5
	// filename will map to a string
	check1, ok := diffResults.InBothSame["/project/file1.txt"]
	if !ok {
		t.Fatalf("Couldn't get check1")
	}
	if check1 != f1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f1.LicenseConcluded, check1)
	}
	check5, ok := diffResults.InBothSame["/project/file5.txt"]
	if !ok {
		t.Fatalf("Couldn't get check5")
	}
	if check5 != f5_1.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_1.LicenseConcluded, check5)
	}
	if check5 != f5_2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f5_2.LicenseConcluded, check5)
	}

	// in first only: f2
	// filename will map to a string
	check2, ok := diffResults.InFirstOnly["/project/file2.txt"]
	if !ok {
		t.Fatalf("Couldn't get check2")
	}
	if check2 != f2.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f2.LicenseConcluded, check2)
	}

	// in second only: f3
	// filename will map to a string
	check3, ok := diffResults.InSecondOnly["/project/file3.txt"]
	if !ok {
		t.Fatalf("Couldn't get check3")
	}
	if check3 != f3.LicenseConcluded {
		t.Errorf("Expected %s, got %s", f3.LicenseConcluded, check2)
	}

}
