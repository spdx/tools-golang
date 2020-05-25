// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

// ElementID represents the identifier string portion of an SPDX element
// identifier. DocElementID should be used for any attributes which can
// contain identifiers defined in a different SPDX document.
// ElementIDs should NOT contain the mandatory 'SPDXRef-' portion.
type ElementID string

// DocElementID represents an SPDX element identifier that could be defined
// in a different SPDX document, and therefore could have a "DocumentRef-"
// portion, such as Relationships and Annotations.
// ElementID is used for attributes in which a "DocumentRef-" portion cannot
// appear, such as a Package or File definition (since it is necessarily
// being defined in the present document).
// DocumentRefID will be the empty string for elements defined in the
// present document.
// DocElementIDs should NOT contain the mandatory 'DocumentRef-' or
// 'SPDXRef-' portions.
type DocElementID struct {
	DocumentRefID string
	ElementRefID  ElementID
}

// TODO: add equivalents for LicenseRef- identifiers

// MakeDocElementID takes strings (without prefixes) for the DocumentRef-
// and SPDXRef- identifiers, and returns a DocElementID. An empty string
// should be used for the DocumentRef- portion if it is referring to the
// present document.
func MakeDocElementID(docRef string, eltRef string) DocElementID {
	return DocElementID{
		DocumentRefID: docRef,
		ElementRefID:  ElementID(eltRef),
	}
}

// RenderElementID takes an ElementID and returns the string equivalent,
// with the SPDXRef- prefix reinserted.
func RenderElementID(eID ElementID) string {
	return "SPDXRef-" + string(eID)
}

// RenderDocElementID takes a DocElementID and returns the string equivalent,
// with the SPDXRef- prefix (and, if applicable, the DocumentRef- prefix)
// reinserted.
func RenderDocElementID(deID DocElementID) string {
	prefix := ""
	if deID.DocumentRefID != "" {
		prefix = "DocumentRef-" + deID.DocumentRefID + ":"
	}
	return prefix + "SPDXRef-" + string(deID.ElementRefID)
}
