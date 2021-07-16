// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Relationship section Saver tests =====
func TestSaver2_1RelationshipSavesText(t *testing.T) {
	rln := &spdx.Relationship2_1{
		RefA:                spdx.MakeDocElementID("", "DOCUMENT"),
		RefB:                spdx.MakeDocElementID("", "2"),
		Relationship:        "DESCRIBES",
		RelationshipComment: "this is a comment",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString(`Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-2
RelationshipComment: this is a comment
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship2_1(rln, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1RelationshipOmitsOptionalFieldsIfEmpty(t *testing.T) {
	rln := &spdx.Relationship2_1{
		RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
		RefB:         spdx.MakeDocElementID("", "2"),
		Relationship: "DESCRIBES",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString("Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-2\n")

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship2_1(rln, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1RelationshipWrapsCommentMultiLine(t *testing.T) {
	rln := &spdx.Relationship2_1{
		RefA:         spdx.MakeDocElementID("", "DOCUMENT"),
		RefB:         spdx.MakeDocElementID("", "2"),
		Relationship: "DESCRIBES",
		RelationshipComment: `this is a
multi-line comment`,
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString(`Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-2
RelationshipComment: <text>this is a
multi-line comment</text>
`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship2_1(rln, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
