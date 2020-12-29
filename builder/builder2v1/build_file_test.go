// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== File section builder tests =====
func TestBuilder2_1CanBuildFileSection(t *testing.T) {
	filePath := "/file1.testdata.txt"
	prefix := "../../testdata/project1/"
	fileNumber := 17

	file1, err := BuildFileSection2_1(filePath, prefix, fileNumber)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if file1 == nil {
		t.Fatalf("expected non-nil file, got nil")
	}
	if file1.Name != "/file1.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file1.testdata.txt", file1.Name)
	}
	if file1.SPDXIdentifier != spdx.ElementID("File17") {
		t.Errorf("expected %v, got %v", "File17", file1.SPDXIdentifier)
	}
	if file1.ChecksumSHA1 != "024f870eb6323f532515f7a09d5646a97083b819" {
		t.Errorf("expected %v, got %v", "024f870eb6323f532515f7a09d5646a97083b819", file1.ChecksumSHA1)
	}
	if file1.ChecksumSHA256 != "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf" {
		t.Errorf("expected %v, got %v", "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf", file1.ChecksumSHA256)
	}
	if file1.ChecksumMD5 != "37c8208479dfe42d2bb29debd6e32d4a" {
		t.Errorf("expected %v, got %v", "37c8208479dfe42d2bb29debd6e32d4a", file1.ChecksumMD5)
	}
	if file1.LicenseConcluded != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.LicenseConcluded)
	}
	if len(file1.LicenseInfoInFile) != 0 {
		t.Errorf("expected %v, got %v", 0, len(file1.LicenseInfoInFile))
	}
	if file1.CopyrightText != "NOASSERTION" {
		t.Errorf("expected %v, got %v", "NOASSERTION", file1.CopyrightText)
	}

}

func TestBuilder2_1BuildFileSectionFailsForInvalidFilePath(t *testing.T) {
	filePath := "/file1.testdata.txt"
	prefix := "oops/wrong/path"
	fileNumber := 11

	_, err := BuildFileSection2_1(filePath, prefix, fileNumber)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
