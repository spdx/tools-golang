// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v2

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Relationship section tests =====
func TestParser2_2FailsIfRelationshipNotSet(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePairForRelationship2_2("Relationship", "SPDXRef-A CONTAINS SPDXRef-B")
	if err == nil {
		t.Errorf("expected error when calling parsePairFromRelationship2_2 without setting rln pointer")
	}
}

func TestParser2_2FailsIfRelationshipCommentWithoutRelationship(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("RelationshipComment", "comment whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_2 for RelationshipComment without Relationship first")
	}
}

func TestParser2_2CanParseRelationshipTags(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// Relationship
	err := parser.parsePair2_2("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
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
	err = parser.parsePair2_2("RelationshipComment", cmt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rln.RelationshipComment != cmt {
		t.Errorf("got %v for RelationshipComment, expected %v", parser.rln.RelationshipComment, cmt)
	}
}

func TestParser2_2InvalidRelationshipTagsNoValueFail(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// no items
	parser.rln = nil
	err := parser.parsePair2_2("Relationship", "")
	if err == nil {
		t.Errorf("expected error for empty items in relationship, got nil")
	}
}

func TestParser2_2InvalidRelationshipTagsOneValueFail(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// one item
	parser.rln = nil
	err := parser.parsePair2_2("Relationship", "DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only one item in relationship, got nil")
	}
}

func TestParser2_2InvalidRelationshipTagsTwoValuesFail(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// two items
	parser.rln = nil
	err := parser.parsePair2_2("Relationship", "SPDXRef-DOCUMENT DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only two items in relationship, got nil")
	}
}

func TestParser2_2InvalidRelationshipTagsThreeValuesSucceed(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// three items but with interspersed additional whitespace
	parser.rln = nil
	err := parser.parsePair2_2("Relationship", "  SPDXRef-DOCUMENT \t   DESCRIBES  SPDXRef-something-else    ")
	if err != nil {
		t.Errorf("expected pass for three items in relationship w/ extra whitespace, got: %v", err)
	}
}

func TestParser2_2InvalidRelationshipTagsFourValuesFail(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair2_2("Relationship", "SPDXRef-a DESCRIBES SPDXRef-b SPDXRef-c")
	if err == nil {
		t.Errorf("expected error for more than three items in relationship, got nil")
	}
}

func TestParser2_2InvalidRelationshipTagsInvalidRefIDs(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair2_2("Relationship", "SPDXRef-a DESCRIBES b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}

	parser.rln = nil
	err = parser.parsePair2_2("Relationship", "a DESCRIBES SPDXRef-b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}
}

func TestParser2_2SpecialValuesValidForRightSideOfRelationship(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// NONE in right side of relationship should pass
	err := parser.parsePair2_2("Relationship", "SPDXRef-a CONTAINS NONE")
	if err != nil {
		t.Errorf("expected nil error for CONTAINS NONE, got %v", err)
	}

	// NOASSERTION in right side of relationship should pass
	err = parser.parsePair2_2("Relationship", "SPDXRef-a CONTAINS NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error for CONTAINS NOASSERTION, got %v", err)
	}

	// NONE in left side of relationship should fail
	err = parser.parsePair2_2("Relationship", "NONE CONTAINS SPDXRef-a")
	if err == nil {
		t.Errorf("expected non-nil error for NONE CONTAINS, got nil")
	}

	// NOASSERTION in left side of relationship should fail
	err = parser.parsePair2_2("Relationship", "NOASSERTION CONTAINS SPDXRef-a")
	if err == nil {
		t.Errorf("expected non-nil error for NOASSERTION CONTAINS, got nil")
	}
}

func TestParser2_2FailsToParseUnknownTagInRelationshipSection(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// Relationship
	err := parser.parsePair2_2("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid tag
	err = parser.parsePairForRelationship2_2("blah", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
