// Package spdxlib contains convenience and utility functions for working
// with an SPDX document that has already been created in memory.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package spdxlib

import (
	"fmt"
	"sort"

	"github.com/spdx/tools-golang/spdx"
)

// GetDescribedPackageIDs2_1 returns a slice of ElementIDs for all Packages
// in this Document that it "describes," according to SPDX rules:
// - If the document has only one Package, its ID is returned.
// - If the document has 2+ Packages, it returns the IDs of those that have
//   a DESCRIBES (or DESCRIBED_BY) relationship to this DOCUMENT. If no
// -
func GetDescribedPackageIDs2_1(doc *spdx.Document2_1) ([]spdx.ElementID, error) {
	// if nil Packages map or zero packages in it, return empty slice
	if doc.Packages == nil {
		return nil, fmt.Errorf("Packages map is nil")
	}
	if len(doc.Packages) == 0 {
		return nil, fmt.Errorf("no Packages in Document")
	}
	if len(doc.Packages) == 1 {
		// get first (only) one and return its ID
		for i := range doc.Packages {
			return []spdx.ElementID{i}, nil
		}
	}

	// two or more packages, so we need to go through the relationships,
	// find DESCRIBES or DESCRIBED_BY for this DOCUMENT, verify they are
	// valid IDs in this document's packages, and return them
	if doc.Relationships == nil {
		return nil, fmt.Errorf("multiple Packages in Document but Relationships slice is nil")
	}
	// collect IDs as strings so we can sort them easily
	eIDStrs := []string{}
	for _, rln := range doc.Relationships {
		if rln.Relationship == "DESCRIBES" && rln.RefA == spdx.MakeDocElementID("", "DOCUMENT") {
			// confirm RefB is actually a package in this document
			if _, ok := doc.Packages[rln.RefB.ElementRefID]; !ok {
				// if it's an unpackaged file, that's valid (no error) but don't return it
				if _, ok2 := doc.UnpackagedFiles[rln.RefB.ElementRefID]; !ok2 {
					return nil, fmt.Errorf("Document DESCRIBES %s but no such Package or unpackaged File", string(rln.RefB.ElementRefID))
				}
			}
			eIDStrs = append(eIDStrs, string(rln.RefB.ElementRefID))
		}
		if rln.Relationship == "DESCRIBED_BY" && rln.RefB == spdx.MakeDocElementID("", "DOCUMENT") {
			// confirm RefA is actually a package in this document
			// if it's an unpackaged file, that's valid (no error) but don't return it
			if _, ok := doc.Packages[rln.RefA.ElementRefID]; !ok {
				// if it's an unpackaged file, that's valid (no error) but don't return it
				if _, ok2 := doc.UnpackagedFiles[rln.RefA.ElementRefID]; !ok2 {
					return nil, fmt.Errorf("%s DESCRIBED_BY Document but no such Package or unpackaged File", string(rln.RefA.ElementRefID))
				}
			}
			eIDStrs = append(eIDStrs, string(rln.RefA.ElementRefID))
		}
	}
	if len(eIDStrs) == 0 {
		return nil, fmt.Errorf("no DESCRIBES or DESCRIBED_BY relationships found for this Document")
	}
	// sort them, convert back to ElementIDs and return
	sort.Strings(eIDStrs)
	eIDs := []spdx.ElementID{}
	for _, eIDStr := range eIDStrs {
		eIDs = append(eIDs, spdx.ElementID(eIDStr))
	}
	return eIDs, nil
}

// GetDescribedPackageIDs2_2 returns a slice of ElementIDs for all Packages
// in this Document that it "describes," according to SPDX rules:
// - If the document has only one Package, its ID is returned.
// - If the document has 2+ Packages, it returns the IDs of those that have
//   a DESCRIBES (or DESCRIBED_BY) relationship to this DOCUMENT. If no
// -
func GetDescribedPackageIDs2_2(doc *spdx.Document2_2) ([]spdx.ElementID, error) {
	// if nil Packages map or zero packages in it, return empty slice
	if doc.Packages == nil {
		return nil, fmt.Errorf("Packages map is nil")
	}
	if len(doc.Packages) == 0 {
		return nil, fmt.Errorf("no Packages in Document")
	}
	if len(doc.Packages) == 1 {
		// get first (only) one and return its ID
		for i := range doc.Packages {
			return []spdx.ElementID{i}, nil
		}
	}

	// two or more packages, so we need to go through the relationships,
	// find DESCRIBES or DESCRIBED_BY for this DOCUMENT, verify they are
	// valid IDs in this document's packages, and return them
	if doc.Relationships == nil {
		return nil, fmt.Errorf("multiple Packages in Document but Relationships slice is nil")
	}
	// collect IDs as strings so we can sort them easily
	eIDStrs := []string{}
	for _, rln := range doc.Relationships {
		if rln.Relationship == "DESCRIBES" && rln.RefA == spdx.MakeDocElementID("", "DOCUMENT") {
			// confirm RefB is actually a package in this document
			if _, ok := doc.Packages[rln.RefB.ElementRefID]; !ok {
				// if it's an unpackaged file, that's valid (no error) but don't return it
				if _, ok2 := doc.UnpackagedFiles[rln.RefB.ElementRefID]; !ok2 {
					return nil, fmt.Errorf("Document DESCRIBES %s but no such Package or unpackaged File", string(rln.RefB.ElementRefID))
				}
			}
			eIDStrs = append(eIDStrs, string(rln.RefB.ElementRefID))
		}
		if rln.Relationship == "DESCRIBED_BY" && rln.RefB == spdx.MakeDocElementID("", "DOCUMENT") {
			// confirm RefA is actually a package in this document
			// if it's an unpackaged file, that's valid (no error) but don't return it
			if _, ok := doc.Packages[rln.RefA.ElementRefID]; !ok {
				// if it's an unpackaged file, that's valid (no error) but don't return it
				if _, ok2 := doc.UnpackagedFiles[rln.RefA.ElementRefID]; !ok2 {
					return nil, fmt.Errorf("%s DESCRIBED_BY Document but no such Package or unpackaged File", string(rln.RefA.ElementRefID))
				}
			}
			eIDStrs = append(eIDStrs, string(rln.RefA.ElementRefID))
		}
	}
	if len(eIDStrs) == 0 {
		return nil, fmt.Errorf("no DESCRIBES or DESCRIBED_BY relationships found for this Document")
	}
	// sort them, convert back to ElementIDs and return
	sort.Strings(eIDStrs)
	eIDs := []spdx.ElementID{}
	for _, eIDStr := range eIDStrs {
		eIDs = append(eIDs, spdx.ElementID(eIDStr))
	}
	return eIDs, nil
}
