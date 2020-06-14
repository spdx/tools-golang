// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v2

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Annotation section tests =====
func TestParser2_2FailsIfAnnotationNotSet(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePairForAnnotation2_2("Annotator", "Person: John Doe (jdoe@example.com)")
	if err == nil {
		t.Errorf("expected error when calling parsePairFromAnnotation2_2 without setting ann pointer")
	}
}

func TestParser2_2FailsIfAnnotationTagUnknown(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	// start with valid annotator
	err := parser.parsePair2_2("Annotator", "Person: John Doe (jdoe@example.com)")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// parse invalid tag, using parsePairForAnnotation2_2(
	err = parser.parsePairForAnnotation2_2("blah", "oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_2FailsIfAnnotationFieldsWithoutAnnotation(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("AnnotationDate", "2018-09-15T17:25:00Z")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_2 for AnnotationDate without Annotator first")
	}
	err = parser.parsePair2_2("AnnotationType", "REVIEW")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_2 for AnnotationType without Annotator first")
	}
	err = parser.parsePair2_2("SPDXREF", "SPDXRef-45")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_2 for SPDXREF without Annotator first")
	}
	err = parser.parsePair2_2("AnnotationComment", "comment whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_2 for AnnotationComment without Annotator first")
	}
}

func TestParser2_2CanParseAnnotationTags(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// Annotator without email address
	err := parser.parsePair2_2("Annotator", "Person: John Doe")
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
	err = parser.parsePair2_2("AnnotationDate", dt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.AnnotationDate != dt {
		t.Errorf("got %v for AnnotationDate, expected %v", parser.ann.AnnotationDate, dt)
	}

	// Annotation type
	aType := "REVIEW"
	err = parser.parsePair2_2("AnnotationType", aType)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.AnnotationType != aType {
		t.Errorf("got %v for AnnotationType, expected %v", parser.ann.AnnotationType, aType)
	}

	// SPDX Identifier Reference
	ref := "SPDXRef-30"
	err = parser.parsePair2_2("SPDXREF", ref)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	deID := parser.ann.AnnotationSPDXIdentifier
	if deID.DocumentRefID != "" || deID.ElementRefID != "30" {
		t.Errorf("got %v for SPDXREF, expected %v", parser.ann.AnnotationSPDXIdentifier, "30")
	}

	// Annotation Comment
	cmt := "this is a comment"
	err = parser.parsePair2_2("AnnotationComment", cmt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.ann.AnnotationComment != cmt {
		t.Errorf("got %v for AnnotationComment, expected %v", parser.ann.AnnotationComment, cmt)
	}
}

func TestParser2_2FailsIfAnnotatorInvalid(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("Annotator", "John Doe (jdoe@example.com)")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_2FailsIfAnnotatorTypeInvalid(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("Annotator", "Human: John Doe (jdoe@example.com)")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_2FailsIfAnnotationRefInvalid(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	// start with valid annotator
	err := parser.parsePair2_2("Annotator", "Person: John Doe (jdoe@example.com)")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePair2_2("SPDXREF", "blah:other")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
