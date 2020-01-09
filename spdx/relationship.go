// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

// Relationship2_1 is a Relationship section of an SPDX Document for
// version 2.1 of the spec.
type Relationship2_1 struct {

	// 7.1: Relationship
	// Cardinality: optional, one or more; one per Relationship2_1
	//              one mandatory for SPDX Document with multiple packages
	// RefA and RefB are first and second item
	// Relationship is type from 7.1.1
	RefA         string
	RefB         string
	Relationship string

	// 7.2: Relationship Comment
	// Cardinality: optional, one
	RelationshipComment string
}
