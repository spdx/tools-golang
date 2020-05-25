// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Relationship section tests =====
func TestParser2_1FailsIfRelationshipNotSet(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePairForRelationship2_1("Relationship", "SPDXRef-A CONTAINS SPDXRef-B")
	if err == nil {
		t.Errorf("expected error when calling parsePairFromRelationship2_1 without setting rln pointer")
	}
}

func TestParser2_1FailsIfRelationshipCommentWithoutRelationship(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("RelationshipComment", "comment whatever")
	if err == nil {
		t.Errorf("expected error when calling parsePair2_1 for RelationshipComment without Relationship first")
	}
}

func TestParser2_1CanParseRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// Relationship
	err := parser.parsePair2_1("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
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
	err = parser.parsePair2_1("RelationshipComment", cmt)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.rln.RelationshipComment != cmt {
		t.Errorf("got %v for RelationshipComment, expected %v", parser.rln.RelationshipComment, cmt)
	}
}

func TestParser2_1InvalidRelationshipTagsNoValueFail(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// no items
	parser.rln = nil
	err := parser.parsePair2_1("Relationship", "")
	if err == nil {
		t.Errorf("expected error for empty items in relationship, got nil")
	}
}

func TestParser2_1InvalidRelationshipTagsOneValueFail(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// one item
	parser.rln = nil
	err := parser.parsePair2_1("Relationship", "DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only one item in relationship, got nil")
	}
}

func TestParser2_1InvalidRelationshipTagsTwoValuesFail(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// two items
	parser.rln = nil
	err := parser.parsePair2_1("Relationship", "SPDXRef-DOCUMENT DESCRIBES")
	if err == nil {
		t.Errorf("expected error for only two items in relationship, got nil")
	}
}

func TestParser2_1InvalidRelationshipTagsThreeValuesSucceed(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// three items but with interspersed additional whitespace
	parser.rln = nil
	err := parser.parsePair2_1("Relationship", "  SPDXRef-DOCUMENT \t   DESCRIBES  SPDXRef-something-else    ")
	if err != nil {
		t.Errorf("expected pass for three items in relationship w/ extra whitespace, got: %v", err)
	}
}

func TestParser2_1InvalidRelationshipTagsFourValuesFail(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair2_1("Relationship", "SPDXRef-a DESCRIBES SPDXRef-b SPDXRef-c")
	if err == nil {
		t.Errorf("expected error for more than three items in relationship, got nil")
	}
}

func TestParser2_1InvalidRelationshipTagsInvalidRefIDs(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// four items
	parser.rln = nil
	err := parser.parsePair2_1("Relationship", "SPDXRef-a DESCRIBES b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}

	parser.rln = nil
	err = parser.parsePair2_1("Relationship", "a DESCRIBES SPDXRef-b")
	if err == nil {
		t.Errorf("expected error for missing SPDXRef- prefix, got nil")
	}
}

func TestParser2_1FailsToParseUnknownTagInRelationshipSection(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// Relationship
	err := parser.parsePair2_1("Relationship", "SPDXRef-something CONTAINS DocumentRef-otherdoc:SPDXRef-something-else")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid tag
	err = parser.parsePairForRelationship2_1("blah", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
