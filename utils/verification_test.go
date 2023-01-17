// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

func TestPackageCanGetVerificationCode(t *testing.T) {
	files := []*spdx.File{
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

	gotCode, err := GetVerificationCode(files, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}

}

func TestPackageCanGetVerificationCodeIgnoringExcludesFile(t *testing.T) {
	files := []*spdx.File{
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

	gotCode, err := GetVerificationCode(files, "thisfile.spdx")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if wantCode.Value != gotCode.Value {
		t.Errorf("expected %v, got %v", wantCode, gotCode)
	}
}

func TestPackageGetVerificationCodeFailsIfNilFileInSlice(t *testing.T) {
	files := []*spdx.File{
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

	_, err := GetVerificationCode(files, "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
