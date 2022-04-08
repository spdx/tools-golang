// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Annotator struct {
	Annotator string
	// including AnnotatorType: one of "Person", "Organization" or "Tool"
	AnnotatorType string
}

// UnmarshalJSON takes an annotator in the typical one-line format and parses it into an Annotator struct.
// This function is also used when unmarshalling YAML
func (a *Annotator) UnmarshalJSON(data []byte) error {
	// annotator will simply be a string
	annotatorStr := string(data)
	annotatorStr = strings.Trim(annotatorStr, "\"")

	annotatorFields := strings.SplitN(annotatorStr, ": ", 2)

	if len(annotatorFields) != 2 {
		return fmt.Errorf("failed to parse Annotator '%s'", annotatorStr)
	}

	a.AnnotatorType = annotatorFields[0]
	a.Annotator = annotatorFields[1]

	return nil
}

// MarshalJSON converts the receiver into a slice of bytes representing an Annotator in string form.
// This function is also used when marshalling to YAML
func (a Annotator) MarshalJSON() ([]byte, error) {
	if a.Annotator != "" {
		return json.Marshal(fmt.Sprintf("%s: %s", a.AnnotatorType, a.Annotator))
	}

	return []byte{}, nil
}

// Annotation2_1 is an Annotation section of an SPDX Document for version 2.1 of the spec.
type Annotation2_1 struct {
	// 8.1: Annotator
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	Annotator Annotator `json:"annotator"`

	// 8.2: Annotation Date: YYYY-MM-DDThh:mm:ssZ
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationDate string `json:"annotationDate"`

	// 8.3: Annotation Type: "REVIEW" or "OTHER"
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationType string `json:"annotationType"`

	// 8.4: SPDX Identifier Reference
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	// This field is not used in hierarchical data formats where the referenced element is clear, such as JSON or YAML.
	AnnotationSPDXIdentifier DocElementID

	// 8.5: Annotation Comment
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationComment string `json:"comment"`
}

// Annotation2_2 is an Annotation section of an SPDX Document for version 2.2 of the spec.
type Annotation2_2 struct {
	// 8.1: Annotator
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	Annotator Annotator `json:"annotator"`

	// 8.2: Annotation Date: YYYY-MM-DDThh:mm:ssZ
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationDate string `json:"annotationDate"`

	// 8.3: Annotation Type: "REVIEW" or "OTHER"
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationType string `json:"annotationType"`

	// 8.4: SPDX Identifier Reference
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	// This field is not used in hierarchical data formats where the referenced element is clear, such as JSON or YAML.
	AnnotationSPDXIdentifier DocElementID

	// 8.5: Annotation Comment
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationComment string `json:"comment"`
}
