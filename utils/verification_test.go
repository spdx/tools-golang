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
			Name:           "file2.txt",
			SPDXIdentifier: "File0",
			ChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		"File1": &spdx.File2_1{
			Name:           "file1.txt",
			SPDXIdentifier: "File1",
			ChecksumSHA1:   "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		"File2": &spdx.File2_1{
			Name:           "file3.txt",
			SPDXIdentifier: "File2",
			ChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		"File3": &spdx.File2_1{
			Name:           "file5.txt",
			SPDXIdentifier: "File3",
			ChecksumSHA1:   "2222222222bbbbbbbbbbccccccccccdddddddddd",
		},
		"File4": &spdx.File2_1{
			Name:           "file4.txt",
			SPDXIdentifier: "File4",
			ChecksumSHA1:   "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
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
			Name:           "file1.txt",
			SPDXIdentifier: "File0",
			ChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File1"): &spdx.File2_1{
			Name:           "file2.txt",
			SPDXIdentifier: "File1",
			ChecksumSHA1:   "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File2"): &spdx.File2_1{
			Name:           "thisfile.spdx",
			SPDXIdentifier: "File2",
			ChecksumSHA1:   "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
		},
		spdx.ElementID("File3"): &spdx.File2_1{
			Name:           "file3.txt",
			SPDXIdentifier: "File3",
			ChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File4"): &spdx.File2_1{
			Name:           "file4.txt",
			SPDXIdentifier: "File4",
			ChecksumSHA1:   "2222222222bbbbbbbbbbccccccccccdddddddddd",
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
			Name:           "file2.txt",
			SPDXIdentifier: "File0",
			ChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File1"): nil,
		spdx.ElementID("File2"): &spdx.File2_1{
			Name:           "file3.txt",
			SPDXIdentifier: "File2",
			ChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
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
			Name:           "file2.txt",
			SPDXIdentifier: "File0",
			ChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		"File1": &spdx.File2_2{
			Name:           "file1.txt",
			SPDXIdentifier: "File1",
			ChecksumSHA1:   "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		"File2": &spdx.File2_2{
			Name:           "file3.txt",
			SPDXIdentifier: "File2",
			ChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		"File3": &spdx.File2_2{
			Name:           "file5.txt",
			SPDXIdentifier: "File3",
			ChecksumSHA1:   "2222222222bbbbbbbbbbccccccccccdddddddddd",
		},
		"File4": &spdx.File2_2{
			Name:           "file4.txt",
			SPDXIdentifier: "File4",
			ChecksumSHA1:   "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
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
			Name:           "file1.txt",
			SPDXIdentifier: "File0",
			ChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File1"): &spdx.File2_2{
			Name:           "file2.txt",
			SPDXIdentifier: "File1",
			ChecksumSHA1:   "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File2"): &spdx.File2_2{
			Name:           "thisfile.spdx",
			SPDXIdentifier: "File2",
			ChecksumSHA1:   "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
		},
		spdx.ElementID("File3"): &spdx.File2_2{
			Name:           "file3.txt",
			SPDXIdentifier: "File3",
			ChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File4"): &spdx.File2_2{
			Name:           "file4.txt",
			SPDXIdentifier: "File4",
			ChecksumSHA1:   "2222222222bbbbbbbbbbccccccccccdddddddddd",
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
			Name:           "file2.txt",
			SPDXIdentifier: "File0",
			ChecksumSHA1:   "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		spdx.ElementID("File1"): nil,
		spdx.ElementID("File2"): &spdx.File2_2{
			Name:           "file3.txt",
			SPDXIdentifier: "File2",
			ChecksumSHA1:   "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
	}

	_, err := GetVerificationCode2_2(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
