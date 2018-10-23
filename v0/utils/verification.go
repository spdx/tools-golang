// Package utils contains various utility functions to support the
// main spdx-go packages.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package utils

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// GetVerificationCode2_1 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.1, section 3.9.4.
func GetVerificationCode2_1(files []*spdx.File2_1, excludeFile string) (string, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return "", fmt.Errorf("got nil file in index %d", i)
		}
		if f.FileName != excludeFile {
			shas = append(shas, f.FileChecksumSHA1)
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
	code := fmt.Sprintf("%x", bs)

	return code, nil
}
