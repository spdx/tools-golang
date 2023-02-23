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
		if !validElementIDs[relationship.RefA.ElementRefID] {
			return fmt.Errorf("%s used in relationship but no such package exists", string(relationship.RefA.ElementRefID))
		}

		if !validElementIDs[relationship.RefB.ElementRefID] {
			return fmt.Errorf("%s used in relationship but no such package exists", string(relationship.RefB.ElementRefID))
		}
	}

	return nil
}

// PopulateJsonSchemaFields populates the JSON schema specific fields (such as those providing a convenience fields aliasing
// relationships (like documentDescribes that represents all identifiers identified by the document).
func PopulateJsonSchemaFields(doc *spdx.Document) error {
	// TODO: do the same for other similar JSON schema fields like "hasFiles"
	documentDescribes := []string{}
	for _, r := range doc.Relationships {
		switch r.Relationship {
		case common.TypeRelationshipDescribe:
			if r.RefA.ElementRefID == doc.SPDXIdentifier {
				documentDescribes = append(documentDescribes, common.RenderDocElementID(r.RefB))
			}
		case common.TypeRelationshipDescribeBy:
			if r.RefB.ElementRefID == doc.SPDXIdentifier {
				documentDescribes = append(documentDescribes, common.RenderDocElementID(r.RefA))
			}
		}
	}

	doc.DocumentDescribesJSON = documentDescribes

	return nil
}

// StripJsonSchemaFields strips the JSON schema specific fields
func StripJsonSchemaFields(doc *spdx.Document) error {
	doc.DocumentDescribesJSON = nil
	return nil
}
