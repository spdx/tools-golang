package spdx

import "testing"

func TestCanExtractExternalDocumentReference(t *testing.T) {
	refstring := "DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 SHA1:d6a770ba38583ed4bb4525bd96e50461655d2759"
	wantDocumentRefID := "DocumentRef-spdx-tool-1.2"
	wantURI := "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"
	wantAlg := ChecksumAlgorithm("SHA1")
	wantChecksum := "d6a770ba38583ed4bb4525bd96e50461655d2759"

	var edr ExternalDocumentRef2_2
	err := edr.FromString(refstring)
	if err != nil {
		t.Errorf("got non-nil error: %v", err)
	}
	if wantDocumentRefID != edr.DocumentRefID.String() {
		t.Errorf("wanted document ref ID %s, got %s", wantDocumentRefID, edr.DocumentRefID)
	}
	if wantURI != edr.URI {
		t.Errorf("wanted URI %s, got %s", wantURI, edr.URI)
	}
	if wantAlg != edr.Checksum.Algorithm {
		t.Errorf("wanted alg %s, got %s", wantAlg, edr.Checksum.Algorithm)
	}
	if wantChecksum != edr.Checksum.Value {
		t.Errorf("wanted checksum %s, got %s", wantChecksum, edr.Checksum.Value)
	}
}

func TestCanExtractExternalDocumentReferenceWithExtraWhitespace(t *testing.T) {
	refstring := "   DocumentRef-spdx-tool-1.2    \t http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 \t SHA1:  \t   d6a770ba38583ed4bb4525bd96e50461655d2759"
	wantDocumentRefID := "DocumentRef-spdx-tool-1.2"
	wantURI := "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"
	wantAlg := ChecksumAlgorithm("SHA1")
	wantChecksum := "d6a770ba38583ed4bb4525bd96e50461655d2759"

	var edr ExternalDocumentRef2_2
	err := edr.FromString(refstring)
	if err != nil {
		t.Errorf("got non-nil error: %v", err)
	}
	if wantDocumentRefID != edr.DocumentRefID.String() {
		t.Errorf("wanted document ref ID %s, got %s", wantDocumentRefID, edr.DocumentRefID)
	}
	if wantURI != edr.URI {
		t.Errorf("wanted URI %s, got %s", wantURI, edr.URI)
	}
	if wantAlg != edr.Checksum.Algorithm {
		t.Errorf("wanted alg %s, got %s", wantAlg, edr.Checksum.Algorithm)
	}
	if wantChecksum != edr.Checksum.Value {
		t.Errorf("wanted checksum %s, got %s", wantChecksum, edr.Checksum.Value)
	}
}

func TestFailsExternalDocumentReferenceWithInvalidFormats(t *testing.T) {
	invalidRefs := []string{
		"whoops",
		"DocumentRef-",
		"DocumentRef-   ",
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 d6a770ba38583ed4bb4525bd96e50461655d2759",
		"DocumentRef-spdx-tool-1.2",
	}
	var edr ExternalDocumentRef2_2
	for _, refstring := range invalidRefs {
		err := edr.FromString(refstring)
		if err == nil {
			t.Errorf("expected non-nil error for %s, got nil", refstring)
		}
	}
}
