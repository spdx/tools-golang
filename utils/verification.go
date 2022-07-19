// Package utils contains various utility functions to support the
// main tools-golang packages.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package utils

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// GetVerificationCode2_1 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.1, section 3.9.4.
func GetVerificationCode2_1(files []*v2_1.File, excludeFile string) (common.PackageVerificationCode, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return common.PackageVerificationCode{}, fmt.Errorf("got nil file for identifier %v", i)
		}
		if f.FileName != excludeFile {
			// find the SHA1 hash, if present
			for _, checksum := range f.Checksums {
				if checksum.Algorithm == common.SHA1 {
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

	code := common.PackageVerificationCode{
		Value:         fmt.Sprintf("%x", bs),
		ExcludedFiles: []string{excludeFile},
	}

	return code, nil
}

// GetVerificationCode2_2 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.2, section 3.9.4.
func GetVerificationCode2_2(files []*v2_2.File, excludeFile string) (common.PackageVerificationCode, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return common.PackageVerificationCode{}, fmt.Errorf("got nil file for identifier %v", i)
		}
		if f.FileName != excludeFile {
			// find the SHA1 hash, if present
			for _, checksum := range f.Checksums {
				if checksum.Algorithm == common.SHA1 {
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

	code := common.PackageVerificationCode{
		Value:         fmt.Sprintf("%x", bs),
		ExcludedFiles: []string{excludeFile},
	}

	return code, nil
}
