// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package core

import "github.com/spdx/tools-golang/spdx/common"

// Annotation is an Annotation section of an SPDX Document for version 2.2 of the spec.
type Annotation struct {
	// CoreMetadata includes additional information to map back and forth
	// to other SPDX types
	CoreMetadata Metadata

	// 8.1: Annotator
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	Annotator common.Annotator `json:"annotator"`

	// 8.2: Annotation Date: YYYY-MM-DDThh:mm:ssZ
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationDate string `json:"annotationDate"`

	// 8.3: Annotation Type: "REVIEW" or "OTHER"
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationType string `json:"annotationType"`

	// 8.4: SPDX Identifier Reference
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	// This field is not used in hierarchical data formats where the referenced element is clear, such as JSON or YAML.
	AnnotationSPDXIdentifier common.DocElementID `json:"-"`

	// 8.5: Annotation Comment
	// Cardinality: conditional (mandatory, one) if there is an Annotation
	AnnotationComment string `json:"comment"`
}
