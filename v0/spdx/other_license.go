// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

// OtherLicense2_1 is an Other License Information section of an
// SPDX Document for version 2.1 of the spec.
type OtherLicense2_1 struct {

	// 6.1: License Identifier: "LicenseRef-[idstring]"
	// Cardinality: conditional (mandatory, one) if license is not
	//              on SPDX License List
	LicenseIdentifier string

	// 6.2: Extracted Text
	// Cardinality: conditional (mandatory, one) if there is a
	//              License Identifier assigned
	ExtractedText string

	// 6.3: License Name: single line of text or "NOASSERTION"
	// Cardinality: conditional (mandatory, one) if license is not
	//              on SPDX License List
	LicenseName string

	// 6.4: License Cross Reference
	// Cardinality: conditional (optiona, one or many) if license
	//              is not on SPDX License List
	LicenseCrossReferences []string

	// 6.5: License Comment
	// Cardinality: optional, one
	LicenseComment string
}
