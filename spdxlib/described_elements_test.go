// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

func TestCanGetIDsOfDescribedPackages(t *testing.T) {
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
			&spdx.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			&spdx.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			&spdx.Relationship{
				RefA:         common.MakeDocElementID("", "p4"),
				RefB:         common.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			&spdx.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs(doc)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	// should be three of the five IDs, returned in alphabetical order
	if len(describedPkgIDs) != 3 {
		t.Fatalf("expected %d packages, got %d", 3, len(describedPkgIDs))
	}
	if describedPkgIDs[0] != common.ElementID("p1") {
		t.Errorf("expected %v, got %v", common.ElementID("p1"), describedPkgIDs[0])
	}
	if describedPkgIDs[1] != common.ElementID("p4") {
		t.Errorf("expected %v, got %v", common.ElementID("p4"), describedPkgIDs[1])
	}
	if describedPkgIDs[2] != common.ElementID("p5") {
		t.Errorf("expected %v, got %v", common.ElementID("p5"), describedPkgIDs[2])
	}
}

func TestGetDescribedPackagesReturnsSinglePackageIfOnlyOne(t *testing.T) {
	// set up document and one package, but no relationships
	// b/c only one package
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
		Packages: []*spdx.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs(doc)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	// should return the single package
	if len(describedPkgIDs) != 1 {
		t.Fatalf("expected %d package, got %d", 1, len(describedPkgIDs))
	}
	if describedPkgIDs[0] != common.ElementID("p1") {
		t.Errorf("expected %v, got %v", common.ElementID("p1"), describedPkgIDs[0])
	}
}

func TestFailsToGetDescribedPackagesIfMoreThanOneWithoutDescribesRelationship(t *testing.T) {
	// set up document and multiple packages, but no DESCRIBES relationships
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
			// different relationship
			&spdx.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	_, err := GetDescribedPackageIDs(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestFailsToGetDescribedPackagesIfMoreThanOneWithNilRelationships(t *testing.T) {
	// set up document and multiple packages, but no relationships slice
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
		Packages: []*spdx.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
		},
	}

	_, err := GetDescribedPackageIDs(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestFailsToGetDescribedPackagesIfZeroPackagesInMap(t *testing.T) {
	// set up document but no packages
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
		Packages:       []*spdx.Package{},
	}

	_, err := GetDescribedPackageIDs(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestFailsToGetDescribedPackagesIfNilMap(t *testing.T) {
	// set up document but no packages
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
	}

	_, err := GetDescribedPackageIDs(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
