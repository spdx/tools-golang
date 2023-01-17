// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package writer

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_1"
)

// ===== Relationship section Saver tests =====
func TestSaverRelationshipSavesText(t *testing.T) {
	rln := &spdx.Relationship{
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
	err := renderRelationship(rln, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaverRelationshipOmitsOptionalFieldsIfEmpty(t *testing.T) {
	rln := &spdx.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", "2"),
		Relationship: "DESCRIBES",
	}

	// what we want to get, as a buffer of bytes
	// no trailing blank newline
	want := bytes.NewBufferString("Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-2\n")

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderRelationship(rln, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaverRelationshipWrapsCommentMultiLine(t *testing.T) {
	rln := &spdx.Relationship{
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
	err := renderRelationship(rln, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}
