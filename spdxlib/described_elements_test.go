// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

// ===== 2.1 tests =====

func Test2_1CanGetIDsOfDescribedPackages(t *testing.T) {
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
			&v2_1.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			&v2_1.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			&v2_1.Relationship{
				RefA:         common.MakeDocElementID("", "p4"),
				RefB:         common.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			&v2_1.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs2_1(doc)
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

func Test2_1GetDescribedPackagesReturnsSinglePackageIfOnlyOne(t *testing.T) {
	// set up document and one package, but no relationships
	// b/c only one package
	doc := &v2_1.Document{
		SPDXVersion:    "SPDX-2.1",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_1.CreationInfo{},
		Packages: []*v2_1.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs2_1(doc)
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

func Test2_1FailsToGetDescribedPackagesIfMoreThanOneWithoutDescribesRelationship(t *testing.T) {
	// set up document and multiple packages, but no DESCRIBES relationships
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
			// different relationship
			&v2_1.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_1FailsToGetDescribedPackagesIfMoreThanOneWithNilRelationships(t *testing.T) {
	// set up document and multiple packages, but no relationships slice
	doc := &v2_1.Document{
		SPDXVersion:    "SPDX-2.1",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_1.CreationInfo{},
		Packages: []*v2_1.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
		},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_1FailsToGetDescribedPackagesIfZeroPackagesInMap(t *testing.T) {
	// set up document but no packages
	doc := &v2_1.Document{
		SPDXVersion:    "SPDX-2.1",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_1.CreationInfo{},
		Packages:       []*v2_1.Package{},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_1FailsToGetDescribedPackagesIfNilMap(t *testing.T) {
	// set up document but no packages
	doc := &v2_1.Document{
		SPDXVersion:    "SPDX-2.1",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_1.CreationInfo{},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.2 tests =====

func Test2_2CanGetIDsOfDescribedPackages(t *testing.T) {
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
			&v2_2.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			&v2_2.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			&v2_2.Relationship{
				RefA:         common.MakeDocElementID("", "p4"),
				RefB:         common.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			&v2_2.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs2_2(doc)
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

func Test2_2GetDescribedPackagesReturnsSinglePackageIfOnlyOne(t *testing.T) {
	// set up document and one package, but no relationships
	// b/c only one package
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.2",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_2.CreationInfo{},
		Packages: []*v2_2.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs2_2(doc)
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

func Test2_2FailsToGetDescribedPackagesIfMoreThanOneWithoutDescribesRelationship(t *testing.T) {
	// set up document and multiple packages, but no DESCRIBES relationships
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
			// different relationship
			&v2_2.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_2FailsToGetDescribedPackagesIfMoreThanOneWithNilRelationships(t *testing.T) {
	// set up document and multiple packages, but no relationships slice
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.2",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_2.CreationInfo{},
		Packages: []*v2_2.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
		},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_2FailsToGetDescribedPackagesIfZeroPackagesInMap(t *testing.T) {
	// set up document but no packages
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.2",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_2.CreationInfo{},
		Packages:       []*v2_2.Package{},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_2FailsToGetDescribedPackagesIfNilMap(t *testing.T) {
	// set up document but no packages
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.2",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_2.CreationInfo{},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.3 tests =====

func Test2_3CanGetIDsOfDescribedPackages(t *testing.T) {
	// set up document and some packages and relationships
	doc := &v2_3.Document{
		SPDXVersion:    "SPDX-2.3",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_3.CreationInfo{},
		Packages: []*v2_3.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*v2_3.Relationship{
			&v2_3.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			&v2_3.Relationship{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			&v2_3.Relationship{
				RefA:         common.MakeDocElementID("", "p4"),
				RefB:         common.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			&v2_3.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs2_3(doc)
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

func Test2_3GetDescribedPackagesReturnsSinglePackageIfOnlyOne(t *testing.T) {
	// set up document and one package, but no relationships
	// b/c only one package
	doc := &v2_3.Document{
		SPDXVersion:    "SPDX-2.3",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_3.CreationInfo{},
		Packages: []*v2_3.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
		},
	}

	// request IDs for DESCRIBES / DESCRIBED_BY relationships
	describedPkgIDs, err := GetDescribedPackageIDs2_3(doc)
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

func Test2_3FailsToGetDescribedPackagesIfMoreThanOneWithoutDescribesRelationship(t *testing.T) {
	// set up document and multiple packages, but no DESCRIBES relationships
	doc := &v2_3.Document{
		SPDXVersion:    "SPDX-2.3",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_3.CreationInfo{},
		Packages: []*v2_3.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*v2_3.Relationship{
			// different relationship
			&v2_3.Relationship{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	_, err := GetDescribedPackageIDs2_3(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_3FailsToGetDescribedPackagesIfMoreThanOneWithNilRelationships(t *testing.T) {
	// set up document and multiple packages, but no relationships slice
	doc := &v2_3.Document{
		SPDXVersion:    "SPDX-2.3",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_3.CreationInfo{},
		Packages: []*v2_3.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
		},
	}

	_, err := GetDescribedPackageIDs2_3(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_3FailsToGetDescribedPackagesIfZeroPackagesInMap(t *testing.T) {
	// set up document but no packages
	doc := &v2_3.Document{
		SPDXVersion:    "SPDX-2.3",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_3.CreationInfo{},
		Packages:       []*v2_3.Package{},
	}

	_, err := GetDescribedPackageIDs2_3(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_3FailsToGetDescribedPackagesIfNilMap(t *testing.T) {
	// set up document but no packages
	doc := &v2_3.Document{
		SPDXVersion:    "SPDX-2.3",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_3.CreationInfo{},
	}

	_, err := GetDescribedPackageIDs2_3(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
