// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Annotation section tests =====
func TestParser2_1FailsIfAnnotationNotSet(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePairForAnnotation2_1("Annotator", "Person: John Doe (jdoe@example.com)")
	if err == nil {
		t.Errorf("expected error when calling parsePairFromAnnotation2_1 without setting ann pointer")
	}
}

func TestParser2_1FailsIfAnnotationTagUnknown(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	// start with valid annotator
	err := parser.parsePair2_1("Annotator", "Person: John Doe (jdoe@example.com)")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// parse invalid tag, using parsePairForAnnotation2_1(
	err = parser.parsePairForAnnotation2_1("blah", "oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfAnnotationFieldsWithoutAnnotation(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("AnnotationDate", "2018-09-15T17:25:00Z")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1 for AnnotationDate without Annotator first")
	}
	err = parser.parsePair2_1("AnnotationType", "REVIEW")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1 for AnnotationType without Annotator first")
	}
	err = parser.parsePair2_1("SPDXREF", "SPDXRef-45")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1 for SPDXREF without Annotator first")
	}
	err = parser.parsePair2_1("AnnotationComment", "comment whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1 for AnnotationComment without Annotator first")
	}
}

func TestParser2_1CanParseAnnotationTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// Annotator without email address
	err := parser.parsePair2_1("Annotator", "Person: John Doe")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.Annotator != "John Doe" {
		t.Errorf("got %v for Annotator, expected John Doe", parser.ann.Annotator)
	}
	if parser.ann.AnnotatorType != "Person" {
		t.Errorf("got %v for AnnotatorType, expected Person", parser.ann.AnnotatorType)
	}

	// Annotation Date
	dt := "2018-09-15T17:32:00Z"
	err = parser.parsePair2_1("AnnotationDate", dt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.AnnotationDate != dt {
		t.Errorf("got %v for AnnotationDate, expected %v", parser.ann.AnnotationDate, dt)
	}

	// Annotation type
	aType := "REVIEW"
	err = parser.parsePair2_1("AnnotationType", aType)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.AnnotationType != aType {
		t.Errorf("got %v for AnnotationType, expected %v", parser.ann.AnnotationType, aType)
	}

	// SPDX Identifier Reference
	ref := "SPDXRef-30"
	err = parser.parsePair2_1("SPDXREF", ref)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	deID := parser.ann.AnnotationSPDXIdentifier
	if deID.DocumentRefID != "" || deID.ElementRefID != "30" {
		t.Errorf("got %v for SPDXREF, expected %v", parser.ann.AnnotationSPDXIdentifier, "30")
	}

	// Annotation Comment
	cmt := "this is a comment"
	err = parser.parsePair2_1("AnnotationComment", cmt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.AnnotationComment != cmt {
		t.Errorf("got %v for AnnotationComment, expected %v", parser.ann.AnnotationComment, cmt)
	}
}

func TestParser2_1FailsIfAnnotatorInvalid(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("Annotator", "John Doe (jdoe@example.com)")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfAnnotatorTypeInvalid(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("Annotator", "Human: John Doe (jdoe@example.com)")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfAnnotationRefInvalid(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	// start with valid annotator
	err := parser.parsePair2_1("Annotator", "Person: John Doe (jdoe@example.com)")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePair2_1("SPDXREF", "blah:other")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

