// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// FilterRelationships returns a slice of Element IDs returned by the given filter closure. The closure is passed
// one relationship at a time, and it can return an ElementID or nil.
func FilterRelationships(doc *spdx.Document, filter func(*spdx.Relationship) *common.ElementID) ([]common.ElementID, error) {
	elementIDs := []common.ElementID{}

	for _, relationship := range doc.Relationships {
		if id := filter(relationship); id != nil {
			elementIDs = append(elementIDs, *id)
		}
	}

	return elementIDs, nil
}
