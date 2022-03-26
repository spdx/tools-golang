// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== File section builder tests =====
func TestBuilder2_2CanBuildFileSection(t *testing.T) {
	filePath := "/file1.testdata.txt"
	prefix := "../../testdata/project1/"
	fileNumber := 17

	file1, err := BuildFileSection2_2(filePath, prefix, fileNumber)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if file1 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file1.FileName != "/file1.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file1.testdata.txt", file1.FileName)
	}
	if file1.FileSPDXIdentifier != spdx.ElementID("File17") {
		t.Errorf("expected %v, got %v", "File17", file1.FileSPDXIdentifier)
	}
	for _, checksum := range file1.FileChecksums {
		switch checksum.Algorithm {
		case spdx.SHA1:
			if checksum.Value != "024f870eb6323f532515f7a09d5646a97083b819" {
				t.Errorf("expected %v, got %v", "024f870eb6323f532515f7a09d5646a97083b819", checksum.Value)
			}
		case spdx.SHA256:
			if checksum.Value != "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf" {
				t.Errorf("expected %v, got %v", "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf", checksum.Value)
			}
		case spdx.MD5:
			if checksum.Value != "37c8208479dfe42d2bb29debd6e32d4a" {
				t.Errorf("expected %v, got %v", "37c8208479dfe42d2bb29debd6e32d4a", checksum.Value)
			}
		}
	}
	if file1.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.LicenseConcluded)
	}
	if len(file1.LicenseInfoInFile) != 1 {
		t.Errorf("expected %v, got %v", 1, len(file1.LicenseInfoInFile))
	} else {
		if file1.LicenseInfoInFile[0] != "NOASSERTION" {
			t.Errorf("expected %v, got %v", "NOASSERTION", file1.LicenseInfoInFile[0])
		}
	}
	if file1.FileCopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.FileCopyrightText)
	}

}

func TestBuilder2_2BuildFileSectionFailsForInvalidFilePath(t *testing.T) {
	filePath := "/file1.testdata.txt"
	prefix := "oops/wrong/path"
	fileNumber := 11

	_, err := BuildFileSection2_2(filePath, prefix, fileNumber)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
