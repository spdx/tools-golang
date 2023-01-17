// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/common"
)

// ===== Relationship section builder tests =====
func TestBuilderCanBuildRelationshipSection(t *testing.T) {
	packageName := "project17"

	rln, err := BuildRelationshipSection(packageName)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rln == nil {
		t.Fatalf("expected non-nil relationship, got nil")
	}
	if rln.RefA != common.MakeDocElementID("", "DOCUMENT") {
		t.Errorf("expected %v, got %v", "DOCUMENT", rln.RefA)
	}
	if rln.RefB != common.MakeDocElementID("", "Package-project17") {
		t.Errorf("expected %v, got %v", "Package-project17", rln.RefB)
	}
	if rln.Relationship != "DESCRIBES" {
		t.Errorf("expected %v, got %v", "DESCRIBES", rln.Relationship)
	}

}
