package spdx

import (
	"encoding/json"
	"testing"
)

func TestMakeDocElementID(t *testing.T) {
	// without DocRef
	docElementID := MakeDocElementID("", "Package")
	if docElementID.String() != "SPDXRef-Package" {
		t.Errorf("expected 'SPDXRef-Package', got %s", docElementID)
		return
	}

	// with DocRef
	docElementID = MakeDocElementID("OtherDoc", "Package")
	if docElementID.String() != "DocumentRef-OtherDoc:SPDXRef-Package" {
		t.Errorf("expected 'DocumentRef-OtherDoc:SPDXRef-Package', got %s", docElementID)
		return
	}
}

func TestDocElementID_UnmarshalJSON(t *testing.T) {
	rawJSON := json.RawMessage("\"DocumentRef-some-doc\"")
	docElementID := DocElementID{}

	err := json.Unmarshal(rawJSON, &docElementID)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if docElementID.DocumentRefID != "some-doc" {
		t.Errorf("Bad!")
		return
	}
}