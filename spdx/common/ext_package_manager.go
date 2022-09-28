// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package common

// ReferenceType is an [idstring] as defined in Appendix VI;
type ReferenceType string

// *Type is the enumerable of ReferenceType
const (
	Cpe23Type ReferenceType = "cpe23Type"
	PurlType  ReferenceType = "purl"
	Ref2Type  ReferenceType = "LocationRef-acmeforge"
)
