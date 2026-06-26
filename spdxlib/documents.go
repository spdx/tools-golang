// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package spdxlib

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// ValidateDocument returns an error if the Document is found to be invalid, or nil if the Document is valid.
// Currently, this only verifies that all Element IDs mentioned in Relationships exist in the Document as either a
// Package or an UnpackagedFile.
func ValidateDocument(doc *spdx.Document) error {
	// cache a map of package IDs for quick lookups
	validElementIDs := make(map[common.ElementID]bool)
	for _, docPackage := range doc.Packages {
		validElementIDs[docPackage.PackageSPDXIdentifier] = true
	}

	for _, unpackagedFile := range doc.Files {
		validElementIDs[unpackagedFile.FileSPDXIdentifier] = true
	}

	// add the Document element ID
	validElementIDs[common.MakeDocElementID("", "DOCUMENT").ElementRefID] = true

	for _, relationship := range doc.Relationships {
		// a special value such as NONE or NOASSERTION (SPDX 2.3 section 11.1)
		// is not an element reference, so skip the existence check for it
		if relationship.RefA.SpecialID == "" && !validElementIDs[relationship.RefA.ElementRefID] {
			return fmt.Errorf("%s used in relationship but no such package exists", string(relationship.RefA.ElementRefID))
		}

		if relationship.RefB.SpecialID == "" && !validElementIDs[relationship.RefB.ElementRefID] {
			return fmt.Errorf("%s used in relationship but no such package exists", string(relationship.RefB.ElementRefID))
		}
	}

	return nil
}
