// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// FilterRelationships returns a slice of Element IDs returned by the given filter closure. The closure is passed
// one relationship at a time, and it can return an ElementID or nil. Nil entries in the Relationships slice are
// skipped, so the closure is only ever handed a relationship it can read.
func FilterRelationships(doc *spdx.Document, filter func(*spdx.Relationship) *common.ElementID) ([]common.ElementID, error) {
	elementIDs := []common.ElementID{}

	for _, relationship := range doc.Relationships {
		if relationship == nil {
			continue
		}

		if id := filter(relationship); id != nil {
			elementIDs = append(elementIDs, *id)
		}
	}

	return elementIDs, nil
}
