// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"testing"
)

// ===== Filesystem and hash functionality tests =====
func TestBuilder2_1CanGetSliceOfFolderContents(t *testing.T) {
	dirRoot := "../../testdata/project1/"

	filePaths, err := getAllFilePaths(dirRoot)
	if err != nil {
		t.Fatalf("expected filePaths, got error: %v", err)
	}

	if filePaths == nil {
		t.Fatalf("expected non-nil filePaths, got nil")
	}
	if len(filePaths) != 5 {
		t.Fatalf("expected %v, got %v", 5, len(filePaths))
	}

	// should be in alphabetical order, with files prefixed with '/'
	if filePaths[0] != "/emptyfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "/emptyfile.testdata.txt", filePaths[0])
	}
	if filePaths[1] != "/file1.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file1.testdata.txt", filePaths[1])
	}
	if filePaths[2] != "/file2.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file2.testdata.txt", filePaths[2])
	}
	if filePaths[3] != "/file3.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file3.testdata.txt", filePaths[3])
	}
	if filePaths[4] != "/folder1/file4.testdata.txt" {
		t.Errorf("expected %v, got %v", "/folder1/file4.testdata.txt", filePaths[4])
	}
}

func TestBuilder2_1GetAllFilePathsFailsForNonExistentDirectory(t *testing.T) {
	dirRoot := "./does/not/exist/"

	_, err := getAllFilePaths(dirRoot)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// FIXME add test to make sure we get an error for a directory without
// FIXME appropriate permissions to read its (sub)contents
