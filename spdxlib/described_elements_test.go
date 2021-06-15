// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== 2.1 tests =====

func Test2_1CanGetIDsOfDescribedPackages(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): &spdx.Package2_1{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): &spdx.Package2_1{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): &spdx.Package2_1{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			spdx.ElementID("p4"): &spdx.Package2_1{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			spdx.ElementID("p5"): &spdx.Package2_1{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship2_1{
			&spdx.Relationship2_1{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			&spdx.Relationship2_1{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			&spdx.Relationship2_1{
				RefA:         spdx.MakeDocElementID("", "p4"),
				RefB:         spdx.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			&spdx.Relationship2_1{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p2"),
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
	if describedPkgIDs[0] != spdx.ElementID("p1") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p1"), describedPkgIDs[0])
	}
	if describedPkgIDs[1] != spdx.ElementID("p4") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p4"), describedPkgIDs[1])
	}
	if describedPkgIDs[2] != spdx.ElementID("p5") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p5"), describedPkgIDs[2])
	}
}

func Test2_1GetDescribedPackagesReturnsSinglePackageIfOnlyOne(t *testing.T) {
	// set up document and one package, but no relationships
	// b/c only one package
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): &spdx.Package2_1{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
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
	if describedPkgIDs[0] != spdx.ElementID("p1") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p1"), describedPkgIDs[0])
	}
}

func Test2_1FailsToGetDescribedPackagesIfMoreThanOneWithoutDescribesRelationship(t *testing.T) {
	// set up document and multiple packages, but no DESCRIBES relationships
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): &spdx.Package2_1{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): &spdx.Package2_1{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): &spdx.Package2_1{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			spdx.ElementID("p4"): &spdx.Package2_1{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			spdx.ElementID("p5"): &spdx.Package2_1{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship2_1{
			// different relationship
			&spdx.Relationship2_1{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p2"),
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
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): &spdx.Package2_1{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): &spdx.Package2_1{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
		},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_1FailsToGetDescribedPackagesIfZeroPackagesInMap(t *testing.T) {
	// set up document but no packages
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_1FailsToGetDescribedPackagesIfNilMap(t *testing.T) {
	// set up document but no packages
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
	}

	_, err := GetDescribedPackageIDs2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.2 tests =====

func Test2_2CanGetIDsOfDescribedPackages(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			spdx.ElementID("p1"): &spdx.Package2_2{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): &spdx.Package2_2{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): &spdx.Package2_2{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			spdx.ElementID("p4"): &spdx.Package2_2{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			spdx.ElementID("p5"): &spdx.Package2_2{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship2_2{
			&spdx.Relationship2_2{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			&spdx.Relationship2_2{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			&spdx.Relationship2_2{
				RefA:         spdx.MakeDocElementID("", "p4"),
				RefB:         spdx.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			&spdx.Relationship2_2{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p2"),
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
	if describedPkgIDs[0] != spdx.ElementID("p1") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p1"), describedPkgIDs[0])
	}
	if describedPkgIDs[1] != spdx.ElementID("p4") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p4"), describedPkgIDs[1])
	}
	if describedPkgIDs[2] != spdx.ElementID("p5") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p5"), describedPkgIDs[2])
	}
}

func Test2_2GetDescribedPackagesReturnsSinglePackageIfOnlyOne(t *testing.T) {
	// set up document and one package, but no relationships
	// b/c only one package
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			spdx.ElementID("p1"): &spdx.Package2_2{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
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
	if describedPkgIDs[0] != spdx.ElementID("p1") {
		t.Errorf("expected %v, got %v", spdx.ElementID("p1"), describedPkgIDs[0])
	}
}

func Test2_2FailsToGetDescribedPackagesIfMoreThanOneWithoutDescribesRelationship(t *testing.T) {
	// set up document and multiple packages, but no DESCRIBES relationships
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			spdx.ElementID("p1"): &spdx.Package2_2{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): &spdx.Package2_2{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): &spdx.Package2_2{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			spdx.ElementID("p4"): &spdx.Package2_2{PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			spdx.ElementID("p5"): &spdx.Package2_2{PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship2_2{
			// different relationship
			&spdx.Relationship2_2{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p2"),
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
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			spdx.ElementID("p1"): &spdx.Package2_2{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): &spdx.Package2_2{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
		},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_2FailsToGetDescribedPackagesIfZeroPackagesInMap(t *testing.T) {
	// set up document but no packages
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func Test2_2FailsToGetDescribedPackagesIfNilMap(t *testing.T) {
	// set up document but no packages
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
	}

	_, err := GetDescribedPackageIDs2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
