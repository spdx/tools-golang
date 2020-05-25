// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Parser review section state change tests =====
func TestParser2_1ReviewStartsNewReviewAfterParsingReviewerTag(t *testing.T) {
	// create the first review
	rev1 := "John Doe"
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{
			Reviewer:     rev1,
			ReviewerType: "Person",
		},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)
	r1 := parser.rev

	// the Document's Reviews should have this one only
	if len(parser.doc.Reviews) != 1 {
		t.Errorf("Expected only one review, got %d", len(parser.doc.Reviews))
	}
	if parser.doc.Reviews[0] != r1 {
		t.Errorf("Expected review %v in Reviews[0], got %v", r1, parser.doc.Reviews[0])
	}
	if parser.doc.Reviews[0].Reviewer != rev1 {
		t.Errorf("expected review name %s in Reviews[0], got %s", rev1, parser.doc.Reviews[0].Reviewer)
	}

	// now add a new review
	rev2 := "Steve"
	rp2 := "Person: Steve"
	err := parser.parsePair2_1("Reviewer", rp2)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psReview2_1 {
		t.Errorf("expected state to be %v, got %v", psReview2_1, parser.st)
	}
	// and a review should be created
	if parser.rev == nil {
		t.Fatalf("parser didn't create new review")
	}
	// and the reviewer's name should be as expected
	if parser.rev.Reviewer != rev2 {
		t.Errorf("expected reviewer name %s, got %s", rev2, parser.rev.Reviewer)
	}
	// and the Document's reviews should be of size 2 and have these two
	if len(parser.doc.Reviews) != 2 {
		t.Fatalf("Expected Reviews to have len 2, got %d", len(parser.doc.Reviews))
	}
	if parser.doc.Reviews[0] != r1 {
		t.Errorf("Expected review %v in Reviews[0], got %v", r1, parser.doc.Reviews[0])
	}
	if parser.doc.Reviews[0].Reviewer != rev1 {
		t.Errorf("expected reviewer name %s in Reviews[0], got %s", rev1, parser.doc.Reviews[0].Reviewer)
	}
	if parser.doc.Reviews[1] != parser.rev {
		t.Errorf("Expected review %v in Reviews[1], got %v", parser.rev, parser.doc.Reviews[1])
	}
	if parser.doc.Reviews[1].Reviewer != rev2 {
		t.Errorf("expected reviewer name %s in Reviews[1], got %s", rev2, parser.doc.Reviews[1].Reviewer)
	}

}

func TestParser2_1ReviewStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{
			Reviewer:     "Jane Doe",
			ReviewerType: "Person",
		},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePair2_1("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should remain unchanged
	if parser.st != psReview2_1 {
		t.Errorf("expected state to be %v, got %v", psReview2_1, parser.st)
	}
	// and the relationship should be in the Document's Relationships
	if len(parser.doc.Relationships) != 1 {
		t.Fatalf("expected doc.Relationships to have len 1, got %d", len(parser.doc.Relationships))
	}
	deID := parser.doc.Relationships[0].RefA
	if deID.DocumentRefID != "" || deID.ElementRefID != "blah" {
		t.Errorf("expected RefA to be %s, got %s", "blah", parser.doc.Relationships[0].RefA)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psReview2_1 {
		t.Errorf("expected state to be %v, got %v", psReview2_1, parser.st)
	}
}

func TestParser2_1ReviewStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{
			Reviewer:     "Jane Doe",
			ReviewerType: "Person",
		},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePair2_1("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_1)
	}

	err = parser.parsePair2_1("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_1)
	}

	err = parser.parsePair2_1("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_1)
	}

	err = parser.parsePair2_1("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_1)
	}

	err = parser.parsePair2_1("AnnotationComment", "i guess i had something to say about this particular file")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_1)
	}

	// and the annotation should be in the Document's Annotations
	if len(parser.doc.Annotations) != 1 {
		t.Fatalf("expected doc.Annotations to have len 1, got %d", len(parser.doc.Annotations))
	}
	if parser.doc.Annotations[0].Annotator != "John Doe ()" {
		t.Errorf("expected Annotator to be %s, got %s", "John Doe ()", parser.doc.Annotations[0].Annotator)
	}
}

func TestParser2_1ReviewFailsAfterParsingOtherSectionTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// can't go back to old sections
	err := parser.parsePair2_1("SPDXVersion", "SPDX-2.1")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1, got nil")
	}
	err = parser.parsePair2_1("PackageName", "whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1, got nil")
	}
	err = parser.parsePair2_1("FileName", "whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1, got nil")
	}
	err = parser.parsePair2_1("LicenseID", "LicenseRef-Lic22")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1, got nil")
	}
}

// ===== Review data section tests =====
func TestParser2_1CanParseReviewTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer (DEPRECATED)
	// handled in subsequent subtests

	// Review Date (DEPRECATED)
	err := parser.parsePairFromReview2_1("ReviewDate", "2018-09-23T08:30:00Z")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.ReviewDate != "2018-09-23T08:30:00Z" {
		t.Errorf("got %v for ReviewDate", parser.rev.ReviewDate)
	}

	// Review Comment (DEPRECATED)
	err = parser.parsePairFromReview2_1("ReviewComment", "this is a comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.ReviewComment != "this is a comment" {
		t.Errorf("got %v for ReviewComment", parser.rev.ReviewComment)
	}
}

func TestParser2_1CanParseReviewerPersonTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer: Person
	err := parser.parsePairFromReview2_1("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.Reviewer != "John Doe" {
		t.Errorf("got %v for Reviewer", parser.rev.Reviewer)
	}
	if parser.rev.ReviewerType != "Person" {
		t.Errorf("got %v for ReviewerType", parser.rev.ReviewerType)
	}
}

func TestParser2_1CanParseReviewerOrganizationTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer: Organization
	err := parser.parsePairFromReview2_1("Reviewer", "Organization: John Doe, Inc.")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.Reviewer != "John Doe, Inc." {
		t.Errorf("got %v for Reviewer", parser.rev.Reviewer)
	}
	if parser.rev.ReviewerType != "Organization" {
		t.Errorf("got %v for ReviewerType", parser.rev.ReviewerType)
	}
}

func TestParser2_1CanParseReviewerToolTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer: Tool
	err := parser.parsePairFromReview2_1("Reviewer", "Tool: scannertool - 1.2.12")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.Reviewer != "scannertool - 1.2.12" {
		t.Errorf("got %v for Reviewer", parser.rev.Reviewer)
	}
	if parser.rev.ReviewerType != "Tool" {
		t.Errorf("got %v for ReviewerType", parser.rev.ReviewerType)
	}
}

func TestParser2_1FailsIfReviewerInvalidFormat(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		rev: &spdx.Review2_1{},
	}
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePairFromReview2_1("Reviewer", "oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfReviewerUnknownType(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		rev: &spdx.Review2_1{},
	}
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePairFromReview2_1("Reviewer", "whoops: John Doe")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1ReviewUnknownTagFails(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:   psReview2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: map[spdx.ElementID]*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &spdx.OtherLicense2_1{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &spdx.Review2_1{},
	}
	parser.doc.Packages["test"] = parser.pkg
	parser.pkg.Files["f1"] = parser.file
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePairFromReview2_1("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}


