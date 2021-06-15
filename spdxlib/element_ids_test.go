package spdxlib

import (
    "github.com/spdx/tools-golang/spdx"
    "reflect"
    "testing"
)

func TestSortElementIDs(t *testing.T) {
    eIDs := []spdx.ElementID{"def", "abc", "123"}
    eIDs = SortElementIDs(eIDs)

    if !reflect.DeepEqual(eIDs, []spdx.ElementID{"123", "abc", "def"}) {
        t.Fatalf("expected sorted ElementIDs, got: %v", eIDs)
    }
}