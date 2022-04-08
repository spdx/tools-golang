// Package spdx contains the struct definition for an SPDX Document
// and its constituent parts.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package spdx

// ExternalDocumentRef2_1 is a reference to an external SPDX document
// as defined in section 2.6 for version 2.1 of the spec.
type ExternalDocumentRef2_1 struct {
	// DocumentRefID is the ID string defined in the start of the
	// reference. It should _not_ contain the "DocumentRef-" part
	// of the mandatory ID string.
	DocumentRefID string `json:"externalDocumentId"`

	// URI is the URI defined for the external document
	URI string `json:"spdxDocument"`

	// Checksum is the actual hash data
	Checksum Checksum `json:"checksum"`
}

// ExternalDocumentRef2_2 is a reference to an external SPDX document
// as defined in section 2.6 for version 2.2 of the spec.
type ExternalDocumentRef2_2 struct {
	// DocumentRefID is the ID string defined in the start of the
	// reference. It should _not_ contain the "DocumentRef-" part
	// of the mandatory ID string.
	DocumentRefID string `json:"externalDocumentId"`

	// URI is the URI defined for the external document
	URI string `json:"spdxDocument"`

	// Checksum is the actual hash data
	Checksum Checksum `json:"checksum"`
}

// Document2_1 is an SPDX Document for version 2.1 of the spec.
// See https://spdx.org/sites/cpstandard/files/pages/files/spdxversion2.1.pdf
type Document2_1 struct {
	// 2.1: SPDX Version; should be in the format "SPDX-2.1"
	// Cardinality: mandatory, one
	SPDXVersion string `json:"spdxVersion"`

	// 2.2: Data License; should be "CC0-1.0"
	// Cardinality: mandatory, one
	DataLicense string `json:"dataLicense"`

	// 2.3: SPDX Identifier; should be "DOCUMENT" to represent
	//      mandatory identifier of SPDXRef-DOCUMENT
	// Cardinality: mandatory, one
	SPDXIdentifier ElementID `json:"SPDXID"`

	// 2.4: Document Name
	// Cardinality: mandatory, one
	DocumentName string `json:"name"`

	// 2.5: Document Namespace
	// Cardinality: mandatory, one
	DocumentNamespace string `json:"documentNamespace"`

	// 2.6: External Document References
	// Cardinality: optional, one or many
	ExternalDocumentReferences []ExternalDocumentRef2_1 `json:"externalDocumentRefs,omitempty"`

	// 2.11: Document Comment
	// Cardinality: optional, one
	DocumentComment string `json:"comment,omitempty"`

	CreationInfo  *CreationInfo2_1   `json:"creationInfo"`
	Packages      []*Package2_1      `json:"packages"`
	Files         []*File2_1         `json:"files"`
	OtherLicenses []*OtherLicense2_1 `json:"hasExtractedLicensingInfos"`
	Relationships []*Relationship2_1 `json:"relationships"`
	Annotations   []*Annotation2_1   `json:"annotations"`
	Snippets      []Snippet2_1       `json:"snippets"`

	// DEPRECATED in version 2.0 of spec
	Reviews []*Review2_1
}

// Document2_2 is an SPDX Document for version 2.2 of the spec.
// See https://spdx.github.io/spdx-spec/v2-draft/ (DRAFT)
type Document2_2 struct {
	// 2.1: SPDX Version; should be in the format "SPDX-2.2"
	// Cardinality: mandatory, one
	SPDXVersion string `json:"spdxVersion"`

	// 2.2: Data License; should be "CC0-1.0"
	// Cardinality: mandatory, one
	DataLicense string `json:"dataLicense"`

	// 2.3: SPDX Identifier; should be "DOCUMENT" to represent
	//      mandatory identifier of SPDXRef-DOCUMENT
	// Cardinality: mandatory, one
	SPDXIdentifier ElementID `json:"SPDXID"`

	// 2.4: Document Name
	// Cardinality: mandatory, one
	DocumentName string `json:"name"`

	// 2.5: Document Namespace
	// Cardinality: mandatory, one
	DocumentNamespace string `json:"documentNamespace"`

	// 2.6: External Document References
	// Cardinality: optional, one or many
	ExternalDocumentReferences []ExternalDocumentRef2_2 `json:"externalDocumentRefs,omitempty"`

	// 2.11: Document Comment
	// Cardinality: optional, one
	DocumentComment string `json:"comment,omitempty"`

	CreationInfo  *CreationInfo2_2   `json:"creationInfo"`
	Packages      []*Package2_2      `json:"packages"`
	Files         []*File2_2         `json:"files"`
	OtherLicenses []*OtherLicense2_2 `json:"hasExtractedLicensingInfos"`
	Relationships []*Relationship2_2 `json:"relationships"`
	Annotations   []*Annotation2_2   `json:"annotations"`
	Snippets      []Snippet2_2       `json:"snippets"`

	// DEPRECATED in version 2.0 of spec
	Reviews []*Review2_2
}
