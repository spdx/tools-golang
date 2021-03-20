// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== 2.1 Verification code functionality tests =====

func TestPackage2_1CanGetVerificationCode(t *testing.T) {
	files := map[spdx.ElementID]*spdx.File2_1{
		"File0": &spdx.File2_1{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			FileChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		"File1": &spdx.File2_1{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File1",
			FileChecksumSHA1:   "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		"File2": &spdx.File2_1{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			FileChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		"File3": &spdx.File2_1{
			FileName:           "file5.txt",
			FileSPDXIdentifier: "File3",
			FileChecksumSHA1:   "2222222222bbbbbbbbbbccccccccccdddddddddd",
		},
		"File4": &spdx.File2_1{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			FileChecksumSHA1:   "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
		},
	}

	wantCode := "ac924b375119c81c1f08c3e2722044bfbbdcd3dc"

	gotCode, err := GetVerificationCode2_1(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode != gotCode {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_1CanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := map[spdx.ElementID]*spdx.File2_1{
		spdx.ElementID("File0"): &spdx.File2_1{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File0",
			FileChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File1"): &spdx.File2_1{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File1",
			FileChecksumSHA1:   "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File2"): &spdx.File2_1{
			FileName:           "thisfile.spdx",
			FileSPDXIdentifier: "File2",
			FileChecksumSHA1:   "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
		},
		spdx.ElementID("File3"): &spdx.File2_1{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File3",
			FileChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File4"): &spdx.File2_1{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			FileChecksumSHA1:   "2222222222bbbbbbbbbbccccccccccdddddddddd",
		},
	}

	wantCode := "17fab1bd18fe5c13b5d3983f1c17e5f88b8ff266"

	gotCode, err := GetVerificationCode2_1(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode != gotCode {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_1GetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := map[spdx.ElementID]*spdx.File2_1{
		spdx.ElementID("File0"): &spdx.File2_1{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			FileChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File1"): nil,
		spdx.ElementID("File2"): &spdx.File2_1{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			FileChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
	}

	_, err := GetVerificationCode2_1(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.2 Verification code functionality tests =====

func TestPackage2_2CanGetVerificationCode(t *testing.T) {
	files := map[spdx.ElementID]*spdx.File2_2{
		"File0": &spdx.File2_2{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		"File1": &spdx.File2_2{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File1",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "3333333333bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		"File2": &spdx.File2_2{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		"File3": &spdx.File2_2{
			FileName:           "file5.txt",
			FileSPDXIdentifier: "File3",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "2222222222bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		"File4": &spdx.File2_2{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
				},
			},
		},
	}

	wantCode := "ac924b375119c81c1f08c3e2722044bfbbdcd3dc"

	gotCode, err := GetVerificationCode2_2(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode != gotCode {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_2CanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := map[spdx.ElementID]*spdx.File2_2{
		spdx.ElementID("File0"): &spdx.File2_2{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File0",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		spdx.ElementID("File1"): &spdx.File2_2{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File1",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "3333333333bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		spdx.ElementID("File2"): &spdx.File2_2{
			FileName:           "thisfile.spdx",
			FileSPDXIdentifier: "File2",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
				},
			},
		},
		spdx.ElementID("File3"): &spdx.File2_2{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File3",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		spdx.ElementID("File4"): &spdx.File2_2{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "2222222222bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
	}

	wantCode := "17fab1bd18fe5c13b5d3983f1c17e5f88b8ff266"

	gotCode, err := GetVerificationCode2_2(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode != gotCode {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_2GetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := map[spdx.ElementID]*spdx.File2_2{
		spdx.ElementID("File0"): &spdx.File2_2{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		spdx.ElementID("File1"): nil,
		spdx.ElementID("File2"): &spdx.File2_2{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				spdx.SHA1: spdx.Checksum{
					Algorithm: spdx.SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
	}

	_, err := GetVerificationCode2_2(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
