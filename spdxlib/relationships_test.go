// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

func TestFindsDependsOnRelationships(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
		Packages: []*spdx.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship{
			{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			{
				RefA:         common.MakeDocElementID("", "p4"),
				RefB:         common.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	eIDs, err := FilterRelationships(doc, func(relationship *spdx.Relationship) *common.ElementID {
		p1EID := common.MakeDocElementID("", "p1")
		if relationship.Relationship == "DEPENDS_ON" && relationship.RefA == p1EID {
			return &relationship.RefB.ElementRefID
		} else if relationship.Relationship == "DEPENDENCY_OF" && relationship.RefB == p1EID {
			return &relationship.RefA.ElementRefID
		}

		return nil
	})
	if err != nil {
		t.Fatalf("expected non-nil err, got: %s", err.Error())
	}

	if len(eIDs) != 1 {
		t.Fatalf("expected 1 ElementID, got: %v", eIDs)
	}

	if eIDs[0] != common.MakeDocElementID("", "p2").ElementRefID {
		t.Fatalf("received unexpected relationship: %v", eIDs[0])
	}
}
