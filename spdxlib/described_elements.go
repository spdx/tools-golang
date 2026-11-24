// Package spdxlib contains convenience and utility functions for working
// with an SPDX document that has already been created in memory.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package spdxlib

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// GetDescribedPackageIDs returns a slice of ElementIDs for all Packages
// in this Document that it "describes," according to SPDX rules:
//   - If the document has only one Package, its ID is returned.
//   - If the document has 2+ Packages, it returns the IDs of those that have
//     a DESCRIBES (or DESCRIBED_BY) relationship to this DOCUMENT.
func GetDescribedPackageIDs(doc *spdx.Document) ([]common.ElementID, error) {
	// if nil Packages map or zero packages in it, return empty slice
	if doc.Packages == nil {
		return nil, fmt.Errorf("Packages map is nil")
	}

	// a nil entry is not a Package, so count only the real ones when deciding
	// which of the two SPDX rules below applies
	count := 0
	var sole *spdx.Package
	for _, pkg := range doc.Packages {
		if pkg == nil {
			continue
		}
		count++
		sole = pkg
	}

	if count == 0 {
		return nil, fmt.Errorf("no Packages in Document")
	}
	if count == 1 {
		// get first (only) one and return its ID
		return []common.ElementID{sole.PackageSPDXIdentifier}, nil
	}

	// two or more packages, so we need to go through the relationships,
	// find DESCRIBES or DESCRIBED_BY for this DOCUMENT, verify they are
	// valid IDs in this document's packages, and return them
	if doc.Relationships == nil {
		return nil, fmt.Errorf("multiple Packages in Document but Relationships slice is nil")
	}

	eIDs, err := FilterRelationships(doc, func(relationship *spdx.Relationship) *common.ElementID {
		refDocument := common.MakeDocElementID("", "DOCUMENT")

		if relationship.Relationship == "DESCRIBES" && relationship.RefA == refDocument {
			return &relationship.RefB.ElementRefID
		} else if relationship.Relationship == "DESCRIBED_BY" && relationship.RefB == refDocument {
			return &relationship.RefA.ElementRefID
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(eIDs) == 0 {
		return nil, fmt.Errorf("no DESCRIBES or DESCRIBED_BY relationships found for this Document")
	}

	eIDs = SortElementIDs(eIDs)

	return eIDs, nil
}
