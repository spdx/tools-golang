// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// ===== 2.1 Verification code functionality tests =====

func TestPackage2_1CanGetVerificationCode(t *testing.T) {
	files := []*v2_1.File{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums:          []common.Checksum{{Value: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File1",
			Checksums:          []common.Checksum{{Value: "3333333333bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums:          []common.Checksum{{Value: "8888888888bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file5.txt",
			FileSPDXIdentifier: "File3",
			Checksums:          []common.Checksum{{Value: "2222222222bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums:          []common.Checksum{{Value: "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa", Algorithm: common.SHA1}},
		},
	}

	wantCode := common.PackageVerificationCode{Value: "ac924b375119c81c1f08c3e2722044bfbbdcd3dc"}

	gotCode, err := GetVerificationCode2_1(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_1CanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := []*v2_1.File{
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File0",
			Checksums:          []common.Checksum{{Value: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File1",
			Checksums:          []common.Checksum{{Value: "3333333333bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "thisfile.spdx",
			FileSPDXIdentifier: "File2",
			Checksums:          []common.Checksum{{Value: "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File3",
			Checksums:          []common.Checksum{{Value: "8888888888bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums:          []common.Checksum{{Value: "2222222222bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
	}

	wantCode := common.PackageVerificationCode{Value: "17fab1bd18fe5c13b5d3983f1c17e5f88b8ff266"}

	gotCode, err := GetVerificationCode2_1(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}
}

func TestPackage2_1GetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := []*v2_1.File{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums:          []common.Checksum{{Value: "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
		nil,
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums:          []common.Checksum{{Value: "8888888888bbbbbbbbbbccccccccccdddddddddd", Algorithm: common.SHA1}},
		},
	}

	_, err := GetVerificationCode2_1(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.2 Verification code functionality tests =====

func TestPackage2_2CanGetVerificationCode(t *testing.T) {
	files := []*v2_2.File{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File1",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "3333333333bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file5.txt",
			FileSPDXIdentifier: "File3",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "2222222222bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
				},
			},
		},
	}

	wantCode := common.PackageVerificationCode{Value: "ac924b375119c81c1f08c3e2722044bfbbdcd3dc"}

	gotCode, err := GetVerificationCode2_2(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackage2_2CanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := []*v2_2.File{
		{
			FileName:           "file1.txt",
			FileSPDXIdentifier: "File0",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File1",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "3333333333bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "thisfile.spdx",
			FileSPDXIdentifier: "File2",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "bbbbbbbbbbccccccccccddddddddddaaaaaaaaaa",
				},
			},
		},
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File3",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "8888888888bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		{
			FileName:           "file4.txt",
			FileSPDXIdentifier: "File4",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "2222222222bbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
	}

	wantCode := common.PackageVerificationCode{Value: "17fab1bd18fe5c13b5d3983f1c17e5f88b8ff266"}

	gotCode, err := GetVerificationCode2_2(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}
}

func TestPackage2_2GetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := []*v2_2.File{
		{
			FileName:           "file2.txt",
			FileSPDXIdentifier: "File0",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
					Value:     "aaaaaaaaaabbbbbbbbbbccccccccccdddddddddd",
				},
			},
		},
		nil,
		{
			FileName:           "file3.txt",
			FileSPDXIdentifier: "File2",
			Checksums: []common.Checksum{
				{
					Algorithm: common.SHA1,
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
