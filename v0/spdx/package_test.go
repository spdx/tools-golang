// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"testing"
)

// ===== Utility functions =====
func TestPackage2_1CanGetVerificationCode(t *testing.T) {
	files := []*File2_1{
		&File2_1{
			FileName:         "file2.txt",
			FileChecksumSHA1: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "file1.txt",
			FileChecksumSHA1: "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "file3.txt",
			FileChecksumSHA1: "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "file5.txt",
			FileChecksumSHA1: "2222222222bbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "file4.txt",
			FileChecksumSHA1: "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
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
	files := []*File2_1{
		&File2_1{
			FileName:         "file1.txt",
			FileChecksumSHA1: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "file2.txt",
			FileChecksumSHA1: "3333333333bbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "thisfile.spdx",
			FileChecksumSHA1: "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
		},
		&File2_1{
			FileName:         "file3.txt",
			FileChecksumSHA1: "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
		&File2_1{
			FileName:         "file4.txt",
			FileChecksumSHA1: "2222222222bbbbbbbbbbccccccccccdddddddddd",
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
	files := []*File2_1{
		&File2_1{
			FileName:         "file2.txt",
			FileChecksumSHA1: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
		},
		nil,
		&File2_1{
			FileName:         "file3.txt",
			FileChecksumSHA1: "8888888888bbbbbbbbbbccccccccccdddddddddd",
		},
	}

	_, err := GetVerificationCode2_1(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
