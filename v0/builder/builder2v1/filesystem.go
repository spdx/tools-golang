// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getAllFilePaths(dirRoot string, pathsIgnored []string) ([]string, error) {
	// paths is a _pointer_ to a slice -- not just a slice.
	// this is so that it can be appropriately modified by append
	// in the sub-function.
	paths := &[]string{}
	prefix := strings.TrimSuffix(dirRoot, "/")

	err := filepath.Walk(dirRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// don't include path if it's a directory
		if fi.IsDir() {
			return nil
		}
		// don't include path if it's a symbolic link
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}

		shortPath := strings.TrimPrefix(path, prefix)

		// don't include path if it should be ignored
		if pathsIgnored != nil && shouldIgnore(shortPath, pathsIgnored) {
			return nil
		}

		// if we got here, record the path
		*paths = append(*paths, shortPath)
		return nil
	})

	return *paths, err
}

func getHashesForFilePath(p string) (string, string, string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()

	var ssha1, ssha256, smd5 string
	hSHA1 := sha1.New()
	hSHA256 := sha256.New()
	hMD5 := md5.New()
	hMulti := io.MultiWriter(hSHA1, hSHA256, hMD5)

	if _, err := io.Copy(hMulti, f); err != nil {
		f.Close()
		return "", "", "", err
	}
	ssha1 = fmt.Sprintf("%x", hSHA1.Sum(nil))
	ssha256 = fmt.Sprintf("%x", hSHA256.Sum(nil))
	smd5 = fmt.Sprintf("%x", hMD5.Sum(nil))

	return ssha1, ssha256, smd5, nil
}

func shouldIgnore(fileName string, pathsIgnored []string) bool {
	fDirs, fFile := filepath.Split(fileName)

	for _, pattern := range pathsIgnored {
		// split into dir(s) and filename
		patternDirs, patternFile := filepath.Split(pattern)
		patternDirStars := strings.HasPrefix(patternDirs, "**")
		if patternDirStars {
			patternDirs = patternDirs[2:]
		}

		// case 1: specific file
		if !patternDirStars && patternDirs == fDirs && patternFile != "" && patternFile == fFile {
			return true
		}

		// case 2: all files in specific directory
		if !patternDirStars && strings.HasPrefix(fDirs, patternDirs) && patternFile == "" {
			return true
		}

		// case 3: specific file in any dir
		if patternDirStars && patternDirs == "/" && patternFile != "" && patternFile == fFile {
			return true
		}

		// case 4: specific file in any matching subdir
		if patternDirStars && strings.Contains(fDirs, patternDirs) && patternFile != "" && patternFile == fFile {
			return true
		}

		// case 5: any file in any matching subdir
		if patternDirStars && strings.Contains(fDirs, patternDirs) && patternFile == "" {
			return true
		}

	}

	// if no match, don't ignore
	return false
}
