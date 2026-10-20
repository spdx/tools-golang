// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"strings"
	"testing"

	spdxjson "github.com/spdx/tools-golang/json"
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

func TestGetDescribedPackageIDsSkipsNilPackages(t *testing.T) {
	// a nil entry is not a Package, so a document carrying one real Package
	// next to it is still a single-Package document
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
		Packages: []*spdx.Package{
			nil,
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
		},
	}

	describedPkgIDs, err := GetDescribedPackageIDs(doc)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(describedPkgIDs) != 1 {
		t.Fatalf("expected %d packages, got %d", 1, len(describedPkgIDs))
	}
	if describedPkgIDs[0] != common.ElementID("p1") {
		t.Errorf("expected %v, got %v", common.ElementID("p1"), describedPkgIDs[0])
	}
}

func TestFailsToGetDescribedPackagesIfOnlyPackageIsNil(t *testing.T) {
	// set up document whose only Packages entry is nil
	doc := &spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &spdx.CreationInfo{},
		Packages:       []*spdx.Package{nil},
	}

	_, err := GetDescribedPackageIDs(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// Reading JSON and then asking for the described packages is the sequence in
// examples/1-load, and json.Read keeps a null packages entry as a nil pointer.
func TestGetDescribedPackageIDsAfterReadingNullPackage(t *testing.T) {
	content := `{
		"spdxVersion": "SPDX-2.3",
		"dataLicense": "CC0-1.0",
		"SPDXID": "SPDXRef-DOCUMENT",
		"name": "null-package",
		"documentNamespace": "https://example.com/null-package",
		"creationInfo": {"created": "2024-01-01T00:00:00Z", "creators": ["Tool: test"]},
		"packages": [
			null,
			{"SPDXID": "SPDXRef-p1", "name": "pkg1", "downloadLocation": "NOASSERTION"}
		]
	}`

	doc, err := spdxjson.Read(strings.NewReader(content))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(doc.Packages) != 2 || doc.Packages[0] != nil {
		t.Fatalf("expected the reader to keep the null entry, got %v", doc.Packages)
	}

	describedPkgIDs, err := GetDescribedPackageIDs(doc)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(describedPkgIDs) != 1 || describedPkgIDs[0] != common.ElementID("p1") {
		t.Errorf("expected [p1], got %v", describedPkgIDs)
	}
}

// v2_1.Document has no UnmarshalJSON of its own, so unlike 2.2 and 2.3 a null
// relationships entry is not dropped at parse time and reaches the filter
// closure that GetDescribedPackageIDs passes to FilterRelationships.
func TestGetDescribedPackageIDsAfterReadingNullRelationship(t *testing.T) {
	content := `{
		"spdxVersion": "SPDX-2.1",
		"dataLicense": "CC0-1.0",
		"SPDXID": "SPDXRef-DOCUMENT",
		"name": "null-relationship",
		"documentNamespace": "https://example.com/null-relationship",
		"creationInfo": {"created": "2024-01-01T00:00:00Z", "creators": ["Tool: test"]},
		"packages": [
			{"SPDXID": "SPDXRef-p1", "name": "pkg1", "downloadLocation": "NOASSERTION"},
			{"SPDXID": "SPDXRef-p2", "name": "pkg2", "downloadLocation": "NOASSERTION"}
		],
		"relationships": [
			null,
			{"spdxElementId": "SPDXRef-DOCUMENT", "relatedSpdxElement": "SPDXRef-p2", "relationshipType": "DESCRIBES"}
		]
	}`

	doc, err := spdxjson.Read(strings.NewReader(content))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(doc.Relationships) != 2 || doc.Relationships[0] != nil {
		t.Fatalf("expected the reader to keep the null entry, got %v", doc.Relationships)
	}

	describedPkgIDs, err := GetDescribedPackageIDs(doc)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(describedPkgIDs) != 1 || describedPkgIDs[0] != common.ElementID("p2") {
		t.Errorf("expected [p2], got %v", describedPkgIDs)
	}
}
