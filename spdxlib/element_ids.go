package spdxlib

import (
    "github.com/spdx/tools-golang/spdx"
    "sort"
)

// SortElementIDs sorts and returns the given slice of ElementIDs
func SortElementIDs(eIDs []spdx.ElementID) []spdx.ElementID {
    sort.Slice(eIDs, func(i, j int) bool {
        return eIDs[i] < eIDs[j]
    })

    return eIDs
}