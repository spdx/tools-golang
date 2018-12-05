// Package licensediff is used to generate a "diff" between the concluded
// licenses in two SPDX Packages, using the filename as the match point.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package licensediff

import (
	"github.com/swinslow/spdx-go/v0/spdx"
)

// LicensePair is a result set where we are talking about two license strings,
// potentially differing, for a single filename between two SPDX Packages.
type LicensePair struct {
	first  string
	second string
}

// MakePairs essentially just consolidates all files and LicenseConcluded
// strings into a single data structure.
func MakePairs(p1 *spdx.Package2_1, p2 *spdx.Package2_1) (map[string]LicensePair, error) {
	pairs := map[string]LicensePair{}

	// first, go through and add all files/licenses from p1
	for _, f := range p1.Files {
		pair := LicensePair{first: f.LicenseConcluded, second: ""}
		pairs[f.FileName] = pair
	}

	// now, go through all files/licenses from p2. If already
	// present, add as .second; if not, create new pair
	for _, f := range p2.Files {
		firstLic := ""
		existingPair, ok := pairs[f.FileName]
		if ok {
			// already present; update it
			firstLic = existingPair.first
		}
		// now, update what's there, either way
		pair := LicensePair{first: firstLic, second: f.LicenseConcluded}
		pairs[f.FileName] = pair
	}

	return pairs, nil
}
