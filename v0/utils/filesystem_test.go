// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package utils

import (
	"testing"
)

// ===== Filesystem and hash functionality tests =====
func TestFilesystemCanGetSliceOfFolderContents(t *testing.T) {
	dirRoot := "../../testdata/project1/"

	filePaths, err := GetAllFilePaths(dirRoot, nil)
	if err != nil {
		t.Fatalf("expected filePaths, got error: %v", err)
	}
	if filePaths == nil {
		t.Fatalf("expected non-nil filePaths, got nil")
	}
	// should only be 5 files
	// symbolic link in project1/symbolic-link should be ignored
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
	if filePaths[2] != "/file3.testdata.txt" {
		t.Errorf("expected %v, got %v", "/file3.testdata.txt", filePaths[2])
	}
	if filePaths[3] != "/folder1/file4.testdata.txt" {
		t.Errorf("expected %v, got %v", "/folder1/file4.testdata.txt", filePaths[3])
	}
	if filePaths[4] != "/lastfile.testdata.txt" {
		t.Errorf("expected %v, got %v", "/lastfile.testdata.txt", filePaths[4])
	}
}

func TestFilesystemGetAllFilePathsFailsForNonExistentDirectory(t *testing.T) {
	dirRoot := "./does/not/exist/"

	_, err := GetAllFilePaths(dirRoot, nil)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestFilesystemCanIgnoreFilesWhenGettingFilePaths(t *testing.T) {
	dirRoot := "../../testdata/project3/"
	pathsIgnored := []string{
		"**/ignoredir/",
		"/excludedir/",
		"**/ignorefile.txt",
		"/alsoEXCLUDEthis.txt",
	}

	filePaths, err := GetAllFilePaths(dirRoot, pathsIgnored)
	if err != nil {
		t.Fatalf("expected filePaths, got error: %v", err)
	}
	if filePaths == nil {
		t.Fatalf("expected non-nil filePaths, got nil")
	}

	// should only be 5 files
	if len(filePaths) != 5 {
		t.Fatalf("expected %v, got %v", 5, len(filePaths))
	}

	// should be in alphabetical order, with files prefixed with '/'
	if filePaths[0] != "/dontscan.txt" {
		t.Errorf("expected %v, got %v", "/dontscan.txt", filePaths[0])
	}
	if filePaths[1] != "/keep/keep.txt" {
		t.Errorf("expected %v, got %v", "/keep/keep.txt", filePaths[1])
	}
	if filePaths[2] != "/keep.txt" {
		t.Errorf("expected %v, got %v", "/keep.txt", filePaths[2])
	}
	if filePaths[3] != "/subdir/keep/dontscan.txt" {
		t.Errorf("expected %v, got %v", "/subdir/keep/dontscan.txt", filePaths[3])
	}
	if filePaths[4] != "/subdir/keep/keep.txt" {
		t.Errorf("expected %v, got %v", "/subdir/keep/keep.txt", filePaths[4])
	}

}

// FIXME add test to make sure we get an error for a directory without
// FIXME appropriate permissions to read its (sub)contents

func TestFilesystemGetsHashesForFilePath(t *testing.T) {
	f := "../../testdata/project1/file1.testdata.txt"

	ssha1, ssha256, smd5, err := GetHashesForFilePath(f)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if ssha1 != "024f870eb6323f532515f7a09d5646a97083b819" {
		t.Errorf("expected %v, got %v", "024f870eb6323f532515f7a09d5646a97083b819", ssha1)
	}
	if ssha256 != "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf" {
		t.Errorf("expected %v, got %v", "b14e44284ca477b4c0db34b15ca4c454b2947cce7883e22321cf2984050e15bf", ssha256)
	}
	if smd5 != "37c8208479dfe42d2bb29debd6e32d4a" {
		t.Errorf("expected %v, got %v", "37c8208479dfe42d2bb29debd6e32d4a", smd5)
	}
}

func TestFilesystemGetsErrorWhenRequestingHashesForInvalidFilePath(t *testing.T) {
	f := "./does/not/exist"

	_, _, _, err := GetHashesForFilePath(f)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// FIXME add test to make sure we get an error for hashes for a file without
// FIXME appropriate permissions to read its contents

func TestFilesystemExcludesForIgnoredPaths(t *testing.T) {
	// one specific file
	pathsIgnored := []string{"/file.txt"}
	fileName := "/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/fileNope.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}

	// two specific files
	pathsIgnored = []string{"/file.txt", "/file2.txt"}
	fileName = "/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/fileNope.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/file2.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}

	// one specific file in specific subdirectory
	pathsIgnored = []string{"/subdir/file.txt"}
	fileName = "/subdir/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/file.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/something/subdir/file.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}

	// one specific file in any directory
	pathsIgnored = []string{"**/file.txt"}
	fileName = "/subdir/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/something/subdir/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/something/fileNope.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}

	// all files in one specific subdirectory (and its subdirectories)
	pathsIgnored = []string{"/subdir/"}
	fileName = "/subdir/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/file.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/subdir/sub2/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/nope/subdir/file.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}

	// all files in subdirectory with this name, wherever it appears
	pathsIgnored = []string{"**/subdir/"}
	fileName = "/subdir/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/file.txt"
	if ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/subdir/sub2/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}
	fileName = "/nope/subdir/file.txt"
	if !ShouldIgnore(fileName, pathsIgnored) {
		t.Errorf("incorrect for %v, ignoring %v", fileName, pathsIgnored)
	}

}
