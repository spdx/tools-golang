// Package licensediff is used to generate a "diff" between the concluded
// licenses in two SPDX Packages, using the filename as the match point.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package licensediff

import (
	"github.com/spdx/tools-golang/spdx"
)

// LicensePair is a result set where we are talking about two license strings,
// potentially differing, for a single filename between two SPDX Packages.
type LicensePair struct {
	First  string
	Second string
}

// MakePairs essentially just consolidates all files and LicenseConcluded
// strings into a single data structure.
func MakePairs(p1 *spdx.Package, p2 *spdx.Package) (map[string]LicensePair, error) {
	pairs := map[string]LicensePair{}

	// first, go through and add all files/licenses from p1
	for _, f := range p1.Files {
		pair := LicensePair{First: f.LicenseConcluded, Second: ""}
		pairs[f.FileName] = pair
	}

	// now, go through all files/licenses from p2. If already
	// present, add as .second; if not, create new pair
	for _, f := range p2.Files {
		firstLic := ""
		existingPair, ok := pairs[f.FileName]
		if ok {
			// already present; update it
			firstLic = existingPair.First
		}
		// now, update what's there, either way
		pair := LicensePair{First: firstLic, Second: f.LicenseConcluded}
		pairs[f.FileName] = pair
	}

	return pairs, nil
}

// LicenseDiff is a structured version of the output of MakePairs. It is
// meant to make it easier to find and report on, e.g., just the files that
// have different licenses, or those that are in just one scan.
type LicenseDiff struct {
	InBothChanged map[string]LicensePair
	InBothSame    map[string]string
	InFirstOnly   map[string]string
	InSecondOnly  map[string]string
}

// MakeResults creates a more structured set of results from the output
// of MakePairs.
func MakeResults(pairs map[string]LicensePair) (*LicenseDiff, error) {
	diff := &LicenseDiff{
		InBothChanged: map[string]LicensePair{},
		InBothSame:    map[string]string{},
		InFirstOnly:   map[string]string{},
		InSecondOnly:  map[string]string{},
	}

	// walk through pairs and allocate them where they belong
	for filename, pair := range pairs {
		if pair.First == pair.Second {
			diff.InBothSame[filename] = pair.First
		} else {
			if pair.First == "" {
				diff.InSecondOnly[filename] = pair.Second
			} else if pair.Second == "" {
				diff.InFirstOnly[filename] = pair.First
			} else {
				diff.InBothChanged[filename] = pair
			}
		}
	}

	return diff, nil
}
