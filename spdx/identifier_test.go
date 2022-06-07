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

func TestCanExtractElementRefsOnlyFromID(t *testing.T) {
	// test with valid ID in this document
	helperForExtractElementID(t, "SPDXRef-file1", false, "file1")
	// test with valid ID in another document
	helperForExtractElementID(t, "DocumentRef-doc2:SPDXRef-file2", true, "")
	// test with invalid ID in this document
	helperForExtractElementID(t, "a:SPDXRef-file1", true, "")
	helperForExtractElementID(t, "file1", true, "")
	helperForExtractElementID(t, "SPDXRef-", true, "")
	helperForExtractElementID(t, "SPDXRef-file1:", true, "")
	// test with invalid ID in another document
	helperForExtractElementID(t, "DocumentRef-doc2", true, "")
	helperForExtractElementID(t, "DocumentRef-doc2:", true, "")
	helperForExtractElementID(t, "DocumentRef-doc2:SPDXRef-", true, "")
	helperForExtractElementID(t, "DocumentRef-doc2:a", true, "")
	helperForExtractElementID(t, "DocumentRef-:", true, "")
	helperForExtractElementID(t, "DocumentRef-:SPDXRef-file1", true, "")
}

func helperForExtractElementID(t *testing.T, tst string, wantErr bool, wantElt string) {
	var eID ElementID
	err := eID.FromString(tst)
	if err != nil && wantErr == false {
		t.Errorf("testing %v: expected nil error, got %v", tst, err)
	}
	if err == nil && wantErr == true {
		t.Errorf("testing %v: expected non-nil error, got nil", tst)
	}
	if eID != ElementID(wantElt) {
		if wantElt == "" {
			t.Errorf("testing %v: want emptyString for ElementRefID, got %v", tst, eID)
		} else {
			t.Errorf("testing %v: want %v for ElementRefID, got %v", tst, wantElt, eID)
		}
	}
}

func TestCanExtractDocumentAndElementRefsFromID(t *testing.T) {
	// test with valid ID in this document
	helperForExtractDocElementID(t, "SPDXRef-file1", false, "", "file1")
	// test with valid ID in another document
	helperForExtractDocElementID(t, "DocumentRef-doc2", false, "doc2", "")
	helperForExtractDocElementID(t, "DocumentRef-doc2:SPDXRef-file2", false, "doc2", "file2")
	// test with invalid ID in this document
	helperForExtractDocElementID(t, "a:SPDXRef-file1", true, "", "")
	helperForExtractDocElementID(t, "file1", true, "", "")
	helperForExtractDocElementID(t, "SPDXRef-", true, "", "")
	helperForExtractDocElementID(t, "SPDXRef-file1:", true, "", "")
	// test with invalid ID in another document
	helperForExtractDocElementID(t, "DocumentRef-doc2:", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-doc2:SPDXRef-", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-doc2:a", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-:", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-:SPDXRef-file1", true, "", "")
	// test with invalid formats
	helperForExtractDocElementID(t, "DocumentRef-doc2:SPDXRef-file1:file2", true, "", "")
}

func helperForExtractDocElementID(t *testing.T, tst string, wantErr bool, wantDoc string, wantElt string) {
	var deID DocElementID
	err := deID.FromString(tst)
	if err != nil && !wantErr {
		t.Errorf("testing %v: expected nil error, got %v", tst, err)
	} else if err == nil && wantErr {
		t.Errorf("testing %v: expected non-nil error, got nil: %+v", tst, deID)
		return
	} else if wantErr {
		// bail early if an error was expected; the behavior when an error occurs is undefined
		return
	}

	if deID.DocumentRefID != wantDoc {
		if wantDoc == "" {
			t.Errorf("testing %v: want empty string for DocumentRefID, got %v", tst, deID.DocumentRefID)
		} else {
			t.Errorf("testing %v: want %v for DocumentRefID, got %v", tst, wantDoc, deID.DocumentRefID)
		}
	}
	if deID.ElementRefID != ElementID(wantElt) {
		if wantElt == "" {
			t.Errorf("testing %v: want empty string for ElementRefID, got %v", tst, deID.ElementRefID)
		} else {
			t.Errorf("testing %v: want %v for ElementRefID, got %v", tst, wantElt, deID.ElementRefID)
		}
	}
}

func TestCanExtractSpecialDocumentIDs(t *testing.T) {
	helperForExtractDocElementSpecial(t, "NONE", false, "", "", "NONE")
	helperForExtractDocElementSpecial(t, "NOASSERTION", false, "", "", "NOASSERTION")
	// test with invalid other words not on permitted list
	helperForExtractDocElementSpecial(t, "FOO", true, "", "", "")
}

func helperForExtractDocElementSpecial(t *testing.T, tst string, wantErr bool, wantDoc string, wantElt string, wantSpecial string) {
	var deID DocElementID
	err := deID.FromString(tst)
	if err != nil && wantErr == false {
		t.Errorf("testing %v: expected nil error, got %v", tst, err)
	}
	if err == nil && wantErr == true {
		t.Errorf("testing %v: expected non-nil error, got nil", tst)
	}
	if deID.DocumentRefID != wantDoc {
		if wantDoc == "" {
			t.Errorf("testing %v: want empty string for DocumentRefID, got %v", tst, deID.DocumentRefID)
		} else {
			t.Errorf("testing %v: want %v for DocumentRefID, got %v", tst, wantDoc, deID.DocumentRefID)
		}
	}
	if deID.ElementRefID != ElementID(wantElt) {
		if wantElt == "" {
			t.Errorf("testing %v: want empty string for ElementRefID, got %v", tst, deID.ElementRefID)
		} else {
			t.Errorf("testing %v: want %v for ElementRefID, got %v", tst, wantElt, deID.ElementRefID)
		}
	}
	if deID.SpecialID != wantSpecial {
		if wantSpecial == "" {
			t.Errorf("testing %v: want empty string for SpecialID, got %v", tst, deID.SpecialID)
		} else {
			t.Errorf("testing %v: want %v for SpecialID, got %v", tst, wantSpecial, deID.SpecialID)
		}
	}
}
