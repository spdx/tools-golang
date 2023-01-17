// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"sort"

	"github.com/spdx/tools-golang/spdx/v2/common"
)

// SortElementIDs sorts and returns the given slice of ElementIDs
func SortElementIDs(eIDs []common.ElementID) []common.ElementID {
	sort.Slice(eIDs, func(i, j int) bool {
		return eIDs[i] < eIDs[j]
	})

	return eIDs
}
