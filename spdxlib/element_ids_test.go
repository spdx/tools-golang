// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
)

func TestSortElementIDs(t *testing.T) {
	eIDs := []common.ElementID{"def", "abc", "123"}
	eIDs = SortElementIDs(eIDs)

	if !reflect.DeepEqual(eIDs, []common.ElementID{"123", "abc", "def"}) {
		t.Fatalf("expected sorted ElementIDs, got: %v", eIDs)
	}
}
