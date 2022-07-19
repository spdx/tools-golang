// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// ===== 2.1 tests =====

func Test2_1FilterForDependencies(t *testing.T) {
	// set up document and some packages and relationships
	doc := &v2_1.Document{
		SPDXVersion:    "SPDX-2.1",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_1.CreationInfo{},
		Packages: []*v2_1.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*v2_1.Relationship{
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
			{
				RefA:         common.MakeDocElementID("", "p4"),
				RefB:         common.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
			{
				RefA:         common.MakeDocElementID("", "p3"),
				RefB:         common.MakeDocElementID("", "p4"),
				Relationship: "DEPENDENCY_OF",
			},
		},
	}

	eIDs, err := FilterRelationships2_1(doc, func(relationship *v2_1.Relationship) *common.ElementID {
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

// ===== 2.2 tests =====

func Test2_2FindsDependsOnRelationships(t *testing.T) {
	// set up document and some packages and relationships
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.2",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_2.CreationInfo{},
		Packages: []*v2_2.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*v2_2.Relationship{
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

	eIDs, err := FilterRelationships2_2(doc, func(relationship *v2_2.Relationship) *common.ElementID {
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
