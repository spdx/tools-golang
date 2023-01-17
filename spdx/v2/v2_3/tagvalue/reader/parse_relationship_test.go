// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import (
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// ===== Relationship section tests =====
func TestParserFailsIfRelationshipNotSet(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}
	err := parser.parsePairForRelationship("Relationship", "SPDXRef-A CONTAINS SPDXRef-B")
	if err == nil {
		t.Errorf("expected error when calling parsePairFromRelationship without setting rln pointer")
	}
}

func TestParserFailsIfRelationshipCommentWithoutRelationship(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}
	err := parser.parsePair("RelationshipComment", "comment whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair for RelationshipComment without Relationship first")
	}
}

func TestParserCanParseRelationshipTags(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// Relationship
	err := parser.parsePair("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
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
	err = parser.parsePair("RelationshipComment", cmt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rln.RelationshipComment != cmt {
		t.Errorf("got %v for RelationshipComment, expected %v", parser.rln.RelationshipComment, cmt)
	}
}

func TestParserInvalidRelationshipTagsNoValueFail(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// no items
	parser.rln = nil
	err := parser.parsePair("Relationship", "")
	if err == nil {
		t.Errorf("expected error for empty items in relationship, got nil")
	}
}

func TestParserInvalidRelationshipTagsOneValueFail(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// one item
	parser.rln = nil
	err := parser.parsePair("Relationship", "DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only one item in relationship, got nil")
	}
}

func TestParserInvalidRelationshipTagsTwoValuesFail(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// two items
	parser.rln = nil
	err := parser.parsePair("Relationship", "SPDXRef-DOCUMENT DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only two items in relationship, got nil")
	}
}

func TestParserInvalidRelationshipTagsThreeValuesSucceed(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// three items but with interspersed additional whitespace
	parser.rln = nil
	err := parser.parsePair("Relationship", "  SPDXRef-DOCUMENT \t   DESCRIBES  SPDXRef-something-else    ")
	if err != nil {
		t.Errorf("expected pass for three items in relationship w/ extra whitespace, got: %v", err)
	}
}

func TestParserInvalidRelationshipTagsFourValuesFail(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair("Relationship", "SPDXRef-a DESCRIBES SPDXRef-b SPDXRef-c")
	if err == nil {
		t.Errorf("expected error for more than three items in relationship, got nil")
	}
}

func TestParserInvalidRelationshipTagsInvalidRefIDs(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair("Relationship", "SPDXRef-a DESCRIBES b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}

	parser.rln = nil
	err = parser.parsePair("Relationship", "a DESCRIBES SPDXRef-b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}
}

func TestParserSpecialValuesValidForRightSideOfRelationship(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// NONE in right side of relationship should pass
	err := parser.parsePair("Relationship", "SPDXRef-a CONTAINS NONE")
	if err != nil {
		t.Errorf("expected nil error for CONTAINS NONE, got %v", err)
	}

	// NOASSERTION in right side of relationship should pass
	err = parser.parsePair("Relationship", "SPDXRef-a CONTAINS NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error for CONTAINS NOASSERTION, got %v", err)
	}

	// NONE in left side of relationship should fail
	err = parser.parsePair("Relationship", "NONE CONTAINS SPDXRef-a")
	if err == nil {
		t.Errorf("expected non-nil error for NONE CONTAINS, got nil")
	}

	// NOASSERTION in left side of relationship should fail
	err = parser.parsePair("Relationship", "NOASSERTION CONTAINS SPDXRef-a")
	if err == nil {
		t.Errorf("expected non-nil error for NOASSERTION CONTAINS, got nil")
	}
}

func TestParserFailsToParseUnknownTagInRelationshipSection(t *testing.T) {
	parser := tvParser{
		doc: &spdx.Document{},
		st:  psCreationInfo,
	}

	// Relationship
	err := parser.parsePair("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid tag
	err = parser.parsePairForRelationship("blah", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
