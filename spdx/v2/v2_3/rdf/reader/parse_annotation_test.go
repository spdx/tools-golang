// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func Test_setAnnotatorFromString(t *testing.T) {
	// TestCase 1: Empty String must raise an error
	ann := &spdx.Annotation{}
	input := ""
	err := setAnnotatorFromString(input, ann)
	if err == nil {
		t.Error("should've raised an error for an empty string")
	}

	// TestCase 2: Invalid annotator type
	ann = &spdx.Annotation{}
	input = "Company: some_company"
	err = setAnnotatorFromString(input, ann)
	if err == nil {
		t.Errorf("should've raised an error for an unknown annotator type")
	}

	// TestCase 3: Valid annotator
	ann = &spdx.Annotation{}
	input = "Person: Rishabh"
	err = setAnnotatorFromString(input, ann)
	if err != nil {
		t.Errorf("unexpected error for a valid annotator")
	}
	if ann.Annotator.AnnotatorType != "Person" {
		t.Errorf("wrnog annotator type: expected: %s, found: %s", "Person", ann.Annotator)
	}
	if ann.Annotator.Annotator != "Rishabh" {
		t.Errorf("wrong annotator: expected: %s, found: %s", "Rishabh", ann.Annotator)
	}
}

func Test_setAnnotationType(t *testing.T) {
	ann := &spdx.Annotation{}
	// TestCase 1: invalid input (empty annotationType)
	err := setAnnotationType("", ann)
	if err == nil {
		t.Errorf("expected an error for empty input")
	}

	// TestCase 2: invalid input (unknown annotation type)
	err = setAnnotationType(NS_SPDX+"annotationType_unknown", ann)
	if err == nil {
		t.Errorf("expected an error for invalid annotationType")
	}

	// TestCase 3: valid input (annotationType_other)
	err = setAnnotationType(SPDX_ANNOTATION_TYPE_OTHER, ann)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ann.AnnotationType != "OTHER" {
		t.Errorf("expected: OTHER, found: %s", ann.AnnotationType)
	}

	// TestCase 4: valid input (annotationType_review)
	err = setAnnotationType(SPDX_ANNOTATION_TYPE_REVIEW, ann)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ann.AnnotationType != "REVIEW" {
		t.Errorf("expected: REVIEW, found: %s", ann.AnnotationType)
	}
}

func Test_setAnnotationToParser(t *testing.T) {
	// TestCase 1: doc is nil (must raise an error)
	parser, _ := parserFromBodyContent(``)
	parser.doc = nil
	err := setAnnotationToParser(parser, &spdx.Annotation{})
	if err == nil {
		t.Errorf("empty doc should've raised an error")
	}

	// TestCase 2: empty annotations should create a new annotations
	//			   list and append the input to it.
	parser, _ = parserFromBodyContent(``)
	parser.doc.Annotations = nil
	err = setAnnotationToParser(parser, &spdx.Annotation{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(parser.doc.Annotations) != 1 {
		t.Errorf("expected doc to have 1 annotation, found %d", len(parser.doc.Annotations))
	}
}

func Test_rdfParser2_3_parseAnnotationFromNode(t *testing.T) {
	// TestCase 1: invalid annotator must raise an error
	parser, _ := parserFromBodyContent(`
		<spdx:Annotation>
			<spdx:annotationDate>2010-01-29T18:30:22Z</spdx:annotationDate>
			<rdfs:comment>Document level annotation</rdfs:comment>
			<spdx:annotator>Company: some company</spdx:annotator>
			<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_other"/>
		</spdx:Annotation>
	`)
	node := parser.gordfParserObj.Triples[0].Subject
	err := parser.parseAnnotationFromNode(node)
	if err == nil {
		t.Errorf("wrong annotator type should've raised an error")
	}

	// TestCase 2: wrong annotation type should raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Annotation>
			<spdx:annotationDate>2010-01-29T18:30:22Z</spdx:annotationDate>
			<rdfs:comment>Document level annotation</rdfs:comment>
			<spdx:annotator>Person: Jane Doe</spdx:annotator>
			<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_unknown"/>
		</spdx:Annotation>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseAnnotationFromNode(node)
	if err == nil {
		t.Errorf("wrong annotation type should've raised an error")
	}

	// TestCase 3: unknown predicate should also raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Annotation>
			<spdx:annotationDate>2010-01-29T18:30:22Z</spdx:annotationDate>
			<rdfs:comment>Document level annotation</rdfs:comment>
			<spdx:annotator>Person: Jane Doe</spdx:annotator>
			<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_unknown"/>
			<spdx:unknownPredicate />
		</spdx:Annotation>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseAnnotationFromNode(node)
	if err == nil {
		t.Errorf("unknown predicate must raise an error")
	}

	// TestCase 4: completely valid annotation
	parser, _ = parserFromBodyContent(`
		<spdx:Annotation>
			<spdx:annotationDate>2010-01-29T18:30:22Z</spdx:annotationDate>
			<rdfs:comment>Document level annotation</rdfs:comment>
			<spdx:annotator>Person: Jane Doe</spdx:annotator>
			<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_other"/>
		</spdx:Annotation>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseAnnotationFromNode(node)
	if err != nil {
		t.Errorf("error parsing valid a annotation")
	}
	if n := len(parser.doc.Annotations); n != 1 {
		t.Errorf("document should've had only one annotation, found %d", n)
	}
	ann := parser.doc.Annotations[0]
	// validating all the attributes of the annotations
	expectedComment := "Document level annotation"
	if ann.AnnotationComment != expectedComment {
		t.Errorf(`expected: "%s", found "%s"`, expectedComment, ann.AnnotationComment)
	}
	expectedDate := "2010-01-29T18:30:22Z"
	if expectedDate != ann.AnnotationDate {
		t.Errorf(`expected: "%s", found "%s"`, expectedDate, ann.AnnotationDate)
	}
	expectedAnnotator := "Jane Doe"
	if expectedAnnotator != ann.Annotator.Annotator {
		t.Errorf(`expected: "%s", found "%s"`, expectedAnnotator, ann.Annotator)
	}
	if ann.Annotator.AnnotatorType != "Person" {
		t.Errorf(`expected: "%s", found "%s"`, "Person", ann.Annotator.AnnotatorType)
	}
	expectedAnnotationType := "OTHER"
	if expectedAnnotationType != ann.AnnotationType {
		t.Errorf(`expected: "%s", found "%s"`, expectedAnnotationType, ann.AnnotationType)
	}
}
