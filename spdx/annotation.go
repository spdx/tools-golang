// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"encoding/json"
	"fmt"
	"github.com/spdx/tools-golang/utils"
	"strings"
)

type Annotator struct {
	Annotator string
	// including AnnotatorType: one of "Person", "Organization" or "Tool"
	AnnotatorType string
}

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (a Annotator) Validate() error {
	if a.Annotator == "" || a.AnnotatorType == "" {
		return fmt.Errorf("invalid Annotator, missing fields. %+v", a)
	}

	return nil
}

// FromString parses an Annotator string into an Annotator struct.
func (a *Annotator) FromString(value string) error {
	annotatorType, annotator, err := utils.ExtractSubs(value)
	if err != nil {
		return err
	}

	a.AnnotatorType = annotatorType
	a.Annotator = annotator

	return nil
}

// String converts the receiver into a string.
func (a Annotator) String() string {
	return fmt.Sprintf("%s: %s", a.AnnotatorType, a.Annotator)
}

// UnmarshalJSON takes an annotator in the typical one-line format and parses it into an Annotator struct.
// This function is also used when unmarshalling YAML
func (a *Annotator) UnmarshalJSON(data []byte) error {
	// annotator will simply be a string
	annotatorStr := string(data)
	annotatorStr = strings.Trim(annotatorStr, "\"")

	return a.FromString(annotatorStr)
}

// MarshalJSON converts the receiver into a slice of bytes representing an Annotator in string form.
// This function is also used when marshalling to YAML
func (a Annotator) MarshalJSON() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(a.String())
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
	AnnotationSPDXIdentifier DocElementID `json:"-"`

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
	AnnotationSPDXIdentifier DocElementID `json:"-"`

	// 8.5: Annotation Comment
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationComment string `json:"comment"`
}
