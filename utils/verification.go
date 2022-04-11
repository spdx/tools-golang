// Package utils contains various utility functions to support the
// main tools-golang packages.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package utils

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	"github.com/spdx/tools-golang/spdx"
)

// GetVerificationCode2_1 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.1, section 3.9.4.
func GetVerificationCode2_1(files []*spdx.File2_1, excludeFile string) (spdx.PackageVerificationCode, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return spdx.PackageVerificationCode{}, fmt.Errorf("got nil file for identifier %v", i)
		}
		if f.FileName != excludeFile {
			// find the SHA1 hash, if present
			for _, checksum := range f.Checksums {
				if checksum.Algorithm == spdx.SHA1 {
					shas = append(shas, checksum.Value)
				}
			}
		}
	}

	// sort the strings
	sort.Strings(shas)

	// concatenate them into one string, with no trailing separators
	shasConcat := strings.Join(shas, "")

	// and get its SHA1 value
	hsha1 := sha1.New()
	hsha1.Write([]byte(shasConcat))
	bs := hsha1.Sum(nil)

	code := spdx.PackageVerificationCode{
		Value:         fmt.Sprintf("%x", bs),
		ExcludedFiles: []string{excludeFile},
	}

	return code, nil
}

// GetVerificationCode2_2 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.2, section 3.9.4.
func GetVerificationCode2_2(files []*spdx.File2_2, excludeFile string) (spdx.PackageVerificationCode, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return spdx.PackageVerificationCode{}, fmt.Errorf("got nil file for identifier %v", i)
		}
		if f.FileName != excludeFile {
			// find the SHA1 hash, if present
			for _, checksum := range f.Checksums {
				if checksum.Algorithm == spdx.SHA1 {
					shas = append(shas, checksum.Value)
				}
			}
		}
	}

	// sort the strings
	sort.Strings(shas)

	// concatenate them into one string, with no trailing separators
	shasConcat := strings.Join(shas, "")

	// and get its SHA1 value
	hsha1 := sha1.New()
	hsha1.Write([]byte(shasConcat))
	bs := hsha1.Sum(nil)

	code := spdx.PackageVerificationCode{
		Value:         fmt.Sprintf("%x", bs),
		ExcludedFiles: []string{excludeFile},
	}

	return code, nil
}
