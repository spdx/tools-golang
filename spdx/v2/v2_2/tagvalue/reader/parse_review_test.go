// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// ===== Parser review section state change tests =====
func TestParserReviewStartsNewReviewAfterParsingReviewerTag(t *testing.T) {
	// create the first review
	rev1 := "John Doe"
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{
			Reviewer:     rev1,
			ReviewerType: "Person",
		},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
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
	err := parser.parsePair("Reviewer", rp2)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psReview {
		t.Errorf("expected state to be %v, got %v", psReview, parser.st)
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

func TestParserReviewStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{
			Reviewer:     "Jane Doe",
			ReviewerType: "Person",
		},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePair("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should remain unchanged
	if parser.st != psReview {
		t.Errorf("expected state to be %v, got %v", psReview, parser.st)
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
	if parser.st != psReview {
		t.Errorf("expected state to be %v, got %v", psReview, parser.st)
	}
}

func TestParserReviewStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{
			Reviewer:     "Jane Doe",
			ReviewerType: "Person",
		},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePair("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview)
	}

	err = parser.parsePair("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview)
	}

	err = parser.parsePair("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview)
	}

	err = parser.parsePair("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview)
	}

	err = parser.parsePair("AnnotationComment", "i guess i had something to say about this particular file")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview)
	}

	// and the annotation should be in the Document's Annotations
	if len(parser.doc.Annotations) != 1 {
		t.Fatalf("expected doc.Annotations to have len 1, got %d", len(parser.doc.Annotations))
	}
	if parser.doc.Annotations[0].Annotator.Annotator != "John Doe ()" {
		t.Errorf("expected Annotator to be %s, got %s", "John Doe ()", parser.doc.Annotations[0].Annotator)
	}
}

func TestParserReviewFailsAfterParsingOtherSectionTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// can't go back to old sections
	err := parser.parsePair("SPDXVersion", "SPDX-2.2")
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
	err = parser.parsePair("LicenseID", "LicenseRef-Lic22")
	if err == nil {
		t.Errorf("expected error when calling parsePair, got nil")
	}
}

// ===== Review data section tests =====
func TestParserCanParseReviewTags(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer (DEPRECATED)
	// handled in subsequent subtests

	// Review Date (DEPRECATED)
	err := parser.parsePairFromReview("ReviewDate", "2018-09-23T08:30:00Z")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.ReviewDate != "2018-09-23T08:30:00Z" {
		t.Errorf("got %v for ReviewDate", parser.rev.ReviewDate)
	}

	// Review Comment (DEPRECATED)
	err = parser.parsePairFromReview("ReviewComment", "this is a comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rev.ReviewComment != "this is a comment" {
		t.Errorf("got %v for ReviewComment", parser.rev.ReviewComment)
	}
}

func TestParserCanParseReviewerPersonTag(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer: Person
	err := parser.parsePairFromReview("Reviewer", "Person: John Doe")
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

func TestParserCanParseReviewerOrganizationTag(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer: Organization
	err := parser.parsePairFromReview("Reviewer", "Organization: John Doe, Inc.")
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

func TestParserCanParseReviewerToolTag(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	// Reviewer: Tool
	err := parser.parsePairFromReview("Reviewer", "Tool: scannertool - 1.2.12")
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

func TestParserFailsIfReviewerInvalidFormat(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psReview,
		rev: &v2_2.Review{},
	}
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePairFromReview("Reviewer", "oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfReviewerUnknownType(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psReview,
		rev: &v2_2.Review{},
	}
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePairFromReview("Reviewer", "whoops: John Doe")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserReviewUnknownTagFails(t *testing.T) {
	parser := tvParser{
		doc:  &v2_2.Document{Packages: []*v2_2.Package{}},
		st:   psReview,
		pkg:  &v2_2.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*v2_2.File{}},
		file: &v2_2.File{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
		otherLic: &v2_2.OtherLicense{
			LicenseIdentifier: "LicenseRef-Lic11",
			LicenseName:       "License 11",
		},
		rev: &v2_2.Review{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
	parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)

	err := parser.parsePairFromReview("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}
