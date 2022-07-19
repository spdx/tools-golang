// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
)

// ===== Creation Info section Saver tests =====
func TestSaver2_1CISavesText(t *testing.T) {
	ci := &v2_1.CreationInfo{
		LicenseListVersion: "2.0",
		Creators: []common.Creator{
			{Creator: "John Doe", CreatorType: "Person"},
			{Creator: "Jane Doe (janedoe@example.com)", CreatorType: "Person"},
			{Creator: "John Doe, Inc.", CreatorType: "Organization"},
			{Creator: "Jane Doe LLC", CreatorType: "Organization"},
			{Creator: "magictool1-1.0", CreatorType: "Tool"},
			{Creator: "magictool2-1.0", CreatorType: "Tool"},
			{Creator: "magictool3-1.0", CreatorType: "Tool"},
		},
		Created:        "2018-10-10T06:20:00Z",
		CreatorComment: "this is a creator comment",
	}

	// what we want to get, as a buffer of bytes
	want := bytes.NewBufferString(`LicenseListVersion: 2.0
Creator: Person: John Doe
Creator: Person: Jane Doe (janedoe@example.com)
Creator: Organization: John Doe, Inc.
Creator: Organization: Jane Doe LLC
Creator: Tool: magictool1-1.0
Creator: Tool: magictool2-1.0
Creator: Tool: magictool3-1.0
Created: 2018-10-10T06:20:00Z
CreatorComment: this is a creator comment

`)

	// render as buffer of bytes
	var got bytes.Buffer
	err := renderCreationInfo2_1(ci, &got)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c := bytes.Compare(want.Bytes(), got.Bytes())
	if c != 0 {
		t.Errorf("Expected %v, got %v", want.String(), got.String())
	}
}

func TestSaver2_1CIOmitsOptionalFieldsIfEmpty(t *testing.T) {
	// --- need at least one creator; do first for Persons ---
	ci1 := &v2_1.CreationInfo{
		Creators: []common.Creator{
			{Creator: "John Doe", CreatorType: "Person"},
		},
		Created: "2018-10-10T06:20:00Z",
	}

	// what we want to get, as a buffer of bytes
	want1 := bytes.NewBufferString(`Creator: Person: John Doe
Created: 2018-10-10T06:20:00Z

`)

	// render as buffer of bytes
	var got1 bytes.Buffer
	err := renderCreationInfo2_1(ci1, &got1)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c1 := bytes.Compare(want1.Bytes(), got1.Bytes())
	if c1 != 0 {
		t.Errorf("Expected %v, got %v", want1.String(), got1.String())
	}

	// --- need at least one creator; now switch to organization ---
	ci2 := &v2_1.CreationInfo{
		Creators: []common.Creator{
			{Creator: "John Doe, Inc.", CreatorType: "Organization"},
		},
		Created: "2018-10-10T06:20:00Z",
	}

	// what we want to get, as a buffer of bytes
	want2 := bytes.NewBufferString(`Creator: Organization: John Doe, Inc.
Created: 2018-10-10T06:20:00Z

`)

	// render as buffer of bytes
	var got2 bytes.Buffer
	err = renderCreationInfo2_1(ci2, &got2)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// check that they match
	c2 := bytes.Compare(want2.Bytes(), got2.Bytes())
	if c2 != 0 {
		t.Errorf("Expected %v, got %v", want2.String(), got2.String())
	}
}
