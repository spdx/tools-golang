// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxlib

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== 2.1 tests =====

func Test2_1ValidDocumentPassesValidation(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): {PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): {PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): {PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			spdx.ElementID("p4"): {PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			spdx.ElementID("p5"): {PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship2_1{
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			{
				RefA:         spdx.MakeDocElementID("", "p4"),
				RefB:         spdx.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	err := ValidateDocument2_1(doc)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err.Error())
	}
}

func Test2_1InvalidDocumentFailsValidation(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:    "SPDX-2.1",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_1{
			spdx.ElementID("p1"): {PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): {PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): {PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
		},
		Relationships: []*spdx.Relationship2_1{
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p2"),
				Relationship: "DESCRIBES",
			},
			// invalid ID p99
			{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p99"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	err := ValidateDocument2_1(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.2 tests =====

func Test2_2ValidDocumentPassesValidation(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			spdx.ElementID("p1"): {PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): {PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): {PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
			spdx.ElementID("p4"): {PackageName: "pkg4", PackageSPDXIdentifier: "p4"},
			spdx.ElementID("p5"): {PackageName: "pkg5", PackageSPDXIdentifier: "p5"},
		},
		Relationships: []*spdx.Relationship2_2{
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// inverse relationship -- should also get detected
			{
				RefA:         spdx.MakeDocElementID("", "p4"),
				RefB:         spdx.MakeDocElementID("", "DOCUMENT"),
				Relationship: "DESCRIBED_BY",
			},
			// different relationship
			{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p2"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	err := ValidateDocument2_2(doc)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err.Error())
	}
}

func Test2_2InvalidDocumentFailsValidation(t *testing.T) {
	// set up document and some packages and relationships
	doc := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:    "SPDX-2.2",
			DataLicense:    "CC0-1.0",
			SPDXIdentifier: spdx.ElementID("DOCUMENT"),
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			spdx.ElementID("p1"): {PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			spdx.ElementID("p2"): {PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			spdx.ElementID("p3"): {PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
		},
		Relationships: []*spdx.Relationship2_2{
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			{
				RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
				RefB:         spdx.MakeDocElementID("", "p5"),
				Relationship: "DESCRIBES",
			},
			// invalid ID p99
			{
				RefA:         spdx.MakeDocElementID("", "p1"),
				RefB:         spdx.MakeDocElementID("", "p99"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	err := ValidateDocument2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
