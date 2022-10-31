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

func Test2_1ValidDocumentPassesValidation(t *testing.T) {
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

	err := ValidateDocument2_1(doc)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err.Error())
	}
}

func Test2_1InvalidDocumentFailsValidation(t *testing.T) {
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
		},
		Relationships: []*v2_1.Relationship{
			{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p1"),
				Relationship: "DESCRIBES",
			},
			{
				RefA:         common.MakeDocElementID("", "DOCUMENT"),
				RefB:         common.MakeDocElementID("", "p2"),
				Relationship: "DESCRIBES",
			},
			// invalid ID p99
			{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p99"),
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
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.1",
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

	err := ValidateDocument2_2(doc)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err.Error())
	}
}

func Test2_2InvalidDocumentFailsValidation(t *testing.T) {
	// set up document and some packages and relationships
	doc := &v2_2.Document{
		SPDXVersion:    "SPDX-2.1",
		DataLicense:    "CC0-1.0",
		SPDXIdentifier: common.ElementID("DOCUMENT"),
		CreationInfo:   &v2_2.CreationInfo{},
		Packages: []*v2_2.Package{
			{PackageName: "pkg1", PackageSPDXIdentifier: "p1"},
			{PackageName: "pkg2", PackageSPDXIdentifier: "p2"},
			{PackageName: "pkg3", PackageSPDXIdentifier: "p3"},
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
			// invalid ID p99
			{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p99"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	err := ValidateDocument2_2(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

// ===== 2.3 tests =====

func Test2_3ValidDocumentPassesValidation(t *testing.T) {
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

	err := ValidateDocument2_3(doc)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err.Error())
	}
}

func Test2_3InvalidDocumentFailsValidation(t *testing.T) {
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
		},
		Relationships: []*v2_3.Relationship{
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
			// invalid ID p99
			{
				RefA:         common.MakeDocElementID("", "p1"),
				RefB:         common.MakeDocElementID("", "p99"),
				Relationship: "DEPENDS_ON",
			},
		},
	}

	err := ValidateDocument2_3(doc)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
