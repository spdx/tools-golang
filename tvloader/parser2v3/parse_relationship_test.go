// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v3

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/v2_3"
)

// ===== Relationship section tests =====
func TestParser2_3FailsIfRelationshipNotSet(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}
	err := parser.parsePairForRelationship2_3("Relationship", "SPDXRef-A CONTAINS SPDXRef-B")
	if err == nil {
		t.Errorf("expected error when calling parsePairFromRelationship2_3 without setting rln pointer")
	}
}

func TestParser2_3FailsIfRelationshipCommentWithoutRelationship(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}
	err := parser.parsePair2_3("RelationshipComment", "comment whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_3 for RelationshipComment without Relationship first")
	}
}

func TestParser2_3CanParseRelationshipTags(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// Relationship
	err := parser.parsePair2_3("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rln.RefA.DocumentRefID != "" || parser.rln.RefA.ElementRefID != "something" {
		t.Errorf("got %v for first part of Relationship, expected something", parser.rln.RefA)
	}
	if parser.rln.RefB.DocumentRefID != "otherdoc" || parser.rln.RefB.ElementRefID != "something-else" {
		t.Errorf("got %v for second part of Relationship, expected otherdoc:something-else", parser.rln.RefB)
	}
	if parser.rln.Relationship != "CONTAINS" {
		t.Errorf("got %v for Relationship type, expected CONTAINS", parser.rln.Relationship)
	}

	// Relationship Comment
	cmt := "this is a comment"
	err = parser.parsePair2_3("RelationshipComment", cmt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rln.RelationshipComment != cmt {
		t.Errorf("got %v for RelationshipComment, expected %v", parser.rln.RelationshipComment, cmt)
	}
}

func TestParser2_3InvalidRelationshipTagsNoValueFail(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// no items
	parser.rln = nil
	err := parser.parsePair2_3("Relationship", "")
	if err == nil {
		t.Errorf("expected error for empty items in relationship, got nil")
	}
}

func TestParser2_3InvalidRelationshipTagsOneValueFail(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// one item
	parser.rln = nil
	err := parser.parsePair2_3("Relationship", "DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only one item in relationship, got nil")
	}
}

func TestParser2_3InvalidRelationshipTagsTwoValuesFail(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// two items
	parser.rln = nil
	err := parser.parsePair2_3("Relationship", "SPDXRef-DOCUMENT DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only two items in relationship, got nil")
	}
}

func TestParser2_3InvalidRelationshipTagsThreeValuesSucceed(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// three items but with interspersed additional whitespace
	parser.rln = nil
	err := parser.parsePair2_3("Relationship", "  SPDXRef-DOCUMENT \t   DESCRIBES  SPDXRef-something-else    ")
	if err != nil {
		t.Errorf("expected pass for three items in relationship w/ extra whitespace, got: %v", err)
	}
}

func TestParser2_3InvalidRelationshipTagsFourValuesFail(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair2_3("Relationship", "SPDXRef-a DESCRIBES SPDXRef-b SPDXRef-c")
	if err == nil {
		t.Errorf("expected error for more than three items in relationship, got nil")
	}
}

func TestParser2_3InvalidRelationshipTagsInvalidRefIDs(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair2_3("Relationship", "SPDXRef-a DESCRIBES b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}

	parser.rln = nil
	err = parser.parsePair2_3("Relationship", "a DESCRIBES SPDXRef-b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}
}

func TestParser2_3SpecialValuesValidForRightSideOfRelationship(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// NONE in right side of relationship should pass
	err := parser.parsePair2_3("Relationship", "SPDXRef-a CONTAINS NONE")
	if err != nil {
		t.Errorf("expected nil error for CONTAINS NONE, got %v", err)
	}

	// NOASSERTION in right side of relationship should pass
	err = parser.parsePair2_3("Relationship", "SPDXRef-a CONTAINS NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error for CONTAINS NOASSERTION, got %v", err)
	}

	// NONE in left side of relationship should fail
	err = parser.parsePair2_3("Relationship", "NONE CONTAINS SPDXRef-a")
	if err == nil {
		t.Errorf("expected non-nil error for NONE CONTAINS, got nil")
	}

	// NOASSERTION in left side of relationship should fail
	err = parser.parsePair2_3("Relationship", "NOASSERTION CONTAINS SPDXRef-a")
	if err == nil {
		t.Errorf("expected non-nil error for NOASSERTION CONTAINS, got nil")
	}
}

func TestParser2_3FailsToParseUnknownTagInRelationshipSection(t *testing.T) {
	parser := tvParser2_3{
		doc: &v2_3.Document{},
		st:  psCreationInfo2_3,
	}

	// Relationship
	err := parser.parsePair2_3("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid tag
	err = parser.parsePairForRelationship2_3("blah", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
