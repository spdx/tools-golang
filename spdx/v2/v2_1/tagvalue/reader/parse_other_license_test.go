// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/v2_1"
)

// ===== Parser other license section state change tests =====
func TestParserOLStartsNewOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	// create the first other license
	olid1 := "LicenseRef-Lic11"
	olname1 := "License 11"

	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_1.OtherLicense{
			LicenseIdentifier: olid1,
			LicenseName:       olname1,
		},
	}
	olic1 := parser.otherLic
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	// the Document's OtherLicenses should have this one only
	if parser.doc.OtherLicenses[0] != olic1 {
		t.Errorf("Expected other license %v in OtherLicenses[0], got %v", olic1, parser.doc.OtherLicenses[0])
	}
	if parser.doc.OtherLicenses[0].LicenseName != olname1 {
		t.Errorf("expected other license name %s in OtherLicenses[0], got %s", olname1, parser.doc.OtherLicenses[0].LicenseName)
	}

	// now add a new other license
	olid2 := "LicenseRef-22"
	olname2 := "License 22"
	err := parser.parsePair("LicenseID", olid2)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psOtherLicense {
		t.Errorf("expected state to be %v, got %v", psOtherLicense, parser.st)
	}
	// and an other license should be created
	if parser.otherLic == nil {
		t.Fatalf("parser didn't create new other license")
	}
	// also parse the new license's name
	err = parser.parsePair("LicenseName", olname2)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should still be correct
	if parser.st != psOtherLicense {
		t.Errorf("expected state to be %v, got %v", psOtherLicense, parser.st)
	}
	// and the other license name should be as expected
	if parser.otherLic.LicenseName != olname2 {
		t.Errorf("expected other license name %s, got %s", olname2, parser.otherLic.LicenseName)
	}
	// and the Document's Other Licenses should be of size 2 and have these two
	if len(parser.doc.OtherLicenses) != 2 {
		t.Fatalf("Expected OtherLicenses to have len 2, got %d", len(parser.doc.OtherLicenses))
	}
	if parser.doc.OtherLicenses[0] != olic1 {
		t.Errorf("Expected other license %v in OtherLicenses[0], got %v", olic1, parser.doc.OtherLicenses[0])
	}
	if parser.doc.OtherLicenses[0].LicenseIdentifier != olid1 {
		t.Errorf("expected other license ID %s in OtherLicenses[0], got %s", olid1, parser.doc.OtherLicenses[0].LicenseIdentifier)
	}
	if parser.doc.OtherLicenses[0].LicenseName != olname1 {
		t.Errorf("expected other license name %s in OtherLicenses[0], got %s", olname1, parser.doc.OtherLicenses[0].LicenseName)
	}
	if parser.doc.OtherLicenses[1] != parser.otherLic {
		t.Errorf("Expected other license %v in OtherLicenses[1], got %v", parser.otherLic, parser.doc.OtherLicenses[1])
	}
	if parser.doc.OtherLicenses[1].LicenseIdentifier != olid2 {
		t.Errorf("expected other license ID %s in OtherLicenses[1], got %s", olid2, parser.doc.OtherLicenses[1].LicenseIdentifier)
	}
	if parser.doc.OtherLicenses[1].LicenseName != olname2 {
		t.Errorf("expected other license name %s in OtherLicenses[1], got %s", olname2, parser.doc.OtherLicenses[1].LicenseName)
	}
}

func TestParserOLMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	err := parser.parsePair("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("expected state to be %v, got %v", psReview, parser.st)
	}
}

func TestParserOtherLicenseStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_1.OtherLicense{
			LicenseIdentifier: "LicenseRef-whatever",
			LicenseName:       "the whatever license",
		},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	err := parser.parsePair("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should remain unchanged
	if parser.st != psOtherLicense {
		t.Errorf("expected state to be %v, got %v", psOtherLicense, parser.st)
	}
	// and the relationship should be in the Document's Relationships
	if len(parser.doc.Relationships) != 1 {
		t.Fatalf("expected doc.Relationships to have len 1, got %d", len(parser.doc.Relationships))
	}
	deID := parser.doc.Relationships[0].RefA
	if deID.DocumentRefID != "" || deID.ElementRefID != "blah" {
		t.Errorf("expected RefA to be %s, got %s", "blah", parser.doc.Relationships[0].RefA)
	}

	err = parser.parsePair("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psOtherLicense {
		t.Errorf("expected state to be %v, got %v", psOtherLicense, parser.st)
	}
}

func TestParserOtherLicenseStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_1.OtherLicense{
			LicenseIdentifier: "LicenseRef-whatever",
			LicenseName:       "the whatever license",
		},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	err := parser.parsePair("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense)
	}

	err = parser.parsePair("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense)
	}

	err = parser.parsePair("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense)
	}

	err = parser.parsePair("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense)
	}

	err = parser.parsePair("AnnotationComment", "i guess i had something to say about this particular file")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense)
	}

	// and the annotation should be in the Document's Annotations
	if len(parser.doc.Annotations) != 1 {
		t.Fatalf("expected doc.Annotations to have len 1, got %d", len(parser.doc.Annotations))
	}
	if parser.doc.Annotations[0].Annotator.Annotator != "John Doe ()" {
		t.Errorf("expected Annotator to be %s, got %s", "John Doe ()", parser.doc.Annotations[0].Annotator)
	}
}

func TestParserOLFailsAfterParsingOtherSectionTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_1.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	// can't go back to old sections
	err := parser.parsePair("SPDXVersion", "SPDX-2.1")
	if err == nil {
		t.Errorf("expected error when calling parsePair, got nil")
	}
	err = parser.parsePair("PackageName", "whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair, got nil")
	}
	err = parser.parsePair("FileName", "whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair, got nil")
	}
}

// ===== Other License data section tests =====
func TestParserCanParseOtherLicenseTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	// License Identifier
	err := parser.parsePairFromOtherLicense("LicenseID", "LicenseRef-Lic11")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.otherLic.LicenseIdentifier != "LicenseRef-Lic11" {
		t.Errorf("got %v for LicenseID", parser.otherLic.LicenseIdentifier)
	}

	// Extracted Text
	err = parser.parsePairFromOtherLicense("ExtractedText", "You are permitted to do anything with the software, hooray!")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.otherLic.ExtractedText != "You are permitted to do anything with the software, hooray!" {
		t.Errorf("got %v for ExtractedText", parser.otherLic.ExtractedText)
	}

	// License Name
	err = parser.parsePairFromOtherLicense("LicenseName", "License 11")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.otherLic.LicenseName != "License 11" {
		t.Errorf("got %v for LicenseName", parser.otherLic.LicenseName)
	}

	// License Cross Reference
	crossRefs := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	}
	for _, cr := range crossRefs {
		err = parser.parsePairFromOtherLicense("LicenseCrossReference", cr)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, refWant := range crossRefs {
		flagFound := false
		for _, refCheck := range parser.otherLic.LicenseCrossReferences {
			if refWant == refCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in LicenseCrossReferences", refWant)
		}
	}
	if len(crossRefs) != len(parser.otherLic.LicenseCrossReferences) {
		t.Errorf("expected %d types in LicenseCrossReferences, got %d", len(crossRefs),
			len(parser.otherLic.LicenseCrossReferences))
	}

	// License Comment
	err = parser.parsePairFromOtherLicense("LicenseComment", "this is a comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.otherLic.LicenseComment != "this is a comment" {
		t.Errorf("got %v for LicenseComment", parser.otherLic.LicenseComment)
	}
}

func TestParserOLUnknownTagFails(t *testing.T) {
	parser := tvParser{
		doc:  &v2_1.Document{Packages: []*v2_1.Package{}},
		st:   psOtherLicense,
		pkg:  &v2_1.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_1.File{}},
		file: &v2_1.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)

	err := parser.parsePairFromOtherLicense("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}
