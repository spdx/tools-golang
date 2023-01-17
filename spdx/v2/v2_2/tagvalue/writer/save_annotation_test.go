// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package writer

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// ===== Annotation section Saver tests =====
func TestSaverAnnotationSavesTextForPerson(t *testing.T) {
	ann := &spdx.Annotation{
		Annotator:                common.Annotator{AnnotatorType: "Person", Annotator: "John Doe"},
		AnnotationDate:           "2018-10-10T17:52:00Z",
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.MakeDocElementID("", "DOCUMENT"),
		AnnotationComment:        "This is an annotation about the SPDX document",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString(`Annotator: Person: John Doe
AnnotationDate: 2018-10-10T17:52:00Z
AnnotationType: REVIEW
SPDXREF: SPDXRef-DOCUMENT
AnnotationComment: This is an annotation about the SPDX document
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderAnnotation(ann, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaverAnnotationSavesTextForOrganization(t *testing.T) {
	ann := &spdx.Annotation{
		Annotator:                common.Annotator{AnnotatorType: "Organization", Annotator: "John Doe, Inc."},
		AnnotationDate:           "2018-10-10T17:52:00Z",
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.MakeDocElementID("", "DOCUMENT"),
		AnnotationComment:        "This is an annotation about the SPDX document",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString(`Annotator: Organization: John Doe, Inc.
AnnotationDate: 2018-10-10T17:52:00Z
AnnotationType: REVIEW
SPDXREF: SPDXRef-DOCUMENT
AnnotationComment: This is an annotation about the SPDX document
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderAnnotation(ann, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaverAnnotationSavesTextForTool(t *testing.T) {
	ann := &spdx.Annotation{
		Annotator:                common.Annotator{AnnotatorType: "Tool", Annotator: "magictool-1.1"},
		AnnotationDate:           "2018-10-10T17:52:00Z",
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.MakeDocElementID("", "DOCUMENT"),
		AnnotationComment:        "This is an annotation about the SPDX document",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString(`Annotator: Tool: magictool-1.1
AnnotationDate: 2018-10-10T17:52:00Z
AnnotationType: REVIEW
SPDXREF: SPDXRef-DOCUMENT
AnnotationComment: This is an annotation about the SPDX document
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderAnnotation(ann, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

// note that the annotation has no optional or multiple fields
