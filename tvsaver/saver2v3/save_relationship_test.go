// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v3

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

// ===== Relationship section Saver tests =====
func TestSaver2_3RelationshipSavesText(t *testing.T) {
	rln := &v2_3.Relationship{
		RefA:                common.MakeDocElementID("", "DOCUMENT"),
		RefB:                common.MakeDocElementID("", "2"),
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
	err := renderRelationship2_3(rln, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_3RelationshipOmitsOptionalFieldsIfEmpty(t *testing.T) {
	rln := &v2_3.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", "2"),
		Relationship: "DESCRIBES",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString("Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-2\n")

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship2_3(rln, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_3RelationshipCanHaveNONEOnRight(t *testing.T) {
	rln := &v2_3.Relationship{
		RefA:         common.MakeDocElementID("", "PackageA"),
		RefB:         common.MakeDocElementSpecial("NONE"),
		Relationship: "DEPENDS_ON",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString("Relationship: SPDXRef-PackageA DEPENDS_ON NONE\n")

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship2_3(rln, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_3RelationshipCanHaveNOASSERTIONOnRight(t *testing.T) {
	rln := &v2_3.Relationship{
		RefA:         common.MakeDocElementID("", "PackageA"),
		RefB:         common.MakeDocElementSpecial("NOASSERTION"),
		Relationship: "DEPENDS_ON",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString("Relationship: SPDXRef-PackageA DEPENDS_ON NOASSERTION\n")

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship2_3(rln, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_3RelationshipWrapsCommentMultiLine(t *testing.T) {
	rln := &v2_3.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", "2"),
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
	err := renderRelationship2_3(rln, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("expected %v, got %v", want.String(), got.String())
	}
}
