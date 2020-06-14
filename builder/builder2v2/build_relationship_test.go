// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Relationship section builder tests =====
func TestBuilder2_2CanBuildRelationshipSection(t *testing.T) {
	packageName := "project17"

	rln, err := BuildRelationshipSection2_2(packageName)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rln == nil {
		t.Fatalf("expected non-nil relationship, got nil")
	}
	if rln.RefA != spdx.MakeDocElementID("", "DOCUMENT") {
		t.Errorf("expected %v, got %v", "DOCUMENT", rln.RefA)
	}
	if rln.RefB != spdx.MakeDocElementID("", "Package-project17") {
		t.Errorf("expected %v, got %v", "Package-project17", rln.RefB)
	}
	if rln.Relationship != "DESCRIBES" {
		t.Errorf("expected %v, got %v", "DESCRIBES", rln.Relationship)
	}

}
