// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v2

import (
	"testing"

	"github.com/spdx/tools-golang/tvloader/reader"
)

// ===== Parser exported entry point tests =====
func TestParser2_2CanParseTagValues(t *testing.T) {
	var tvPairs []reader.TagValuePair

	// create some pairs
	tvPair1 := reader.TagValuePair{Tag: "SPDXVersion", Value: "SPDX-2.2"}
	tvPairs = append(tvPairs, tvPair1)
	tvPair2 := reader.TagValuePair{Tag: "DataLicense", Value: "CC0-1.0"}
	tvPairs = append(tvPairs, tvPair2)
	tvPair3 := reader.TagValuePair{Tag: "SPDXID", Value: "SPDXRef-DOCUMENT"}
	tvPairs = append(tvPairs, tvPair3)

	// now parse them
	doc, err := ParseTagValues(tvPairs)
	if err != nil {
		t.Errorf("got error when calling ParseTagValues: %v", err)
	}
	if doc.CreationInfo.SPDXVersion != "SPDX-2.2" {
		t.Errorf("expected SPDXVersion to be SPDX-2.2, got %v", doc.CreationInfo.SPDXVersion)
	}
	if doc.CreationInfo.DataLicense != "CC0-1.0" {
		t.Errorf("expected DataLicense to be CC0-1.0, got %v", doc.CreationInfo.DataLicense)
	}
	if doc.CreationInfo.SPDXIdentifier != "DOCUMENT" {
		t.Errorf("expected SPDXIdentifier to be DOCUMENT, got %v", doc.CreationInfo.SPDXIdentifier)
	}

}

// ===== Parser initialization tests =====
func TestParser2_2InitCreatesResetStatus(t *testing.T) {
	parser := tvParser2_2{}
	if parser.st != psStart2_2 {
		t.Errorf("parser did not begin in start state")
	}
	if parser.doc != nil {
		t.Errorf("parser did not begin with nil document")
	}
}

func TestParser2_2HasDocumentAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_2{}
	err := parser.parsePair2_2("SPDXVersion", "SPDX-2.2")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.doc == nil {
		t.Errorf("doc is still nil after parsing first pair")
	}
}

// ===== Parser start state change tests =====
func TestParser2_2StartMovesToCreationInfoStateAfterParsingFirstTag(t *testing.T) {
	parser := tvParser2_2{}
	err := parser.parsePair2_2("SPDXVersion", "b")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}
}

func TestParser2_2StartFailsToParseIfInInvalidState(t *testing.T) {
	parser := tvParser2_2{st: psReview2_2}
	err := parser.parsePairFromStart2_2("SPDXVersion", "SPDX-2.2")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_2FilesWithoutSpdxIdThrowErrorAtCompleteParse(t *testing.T) {
	// case: Checks the last file
	// Last unpackaged file with no packages in doc
	// Last file of last package in the doc
	tvPairs := []reader.TagValuePair{
		{Tag: "SPDXVersion", Value: "SPDX-2.2"},
		{Tag: "DataLicense", Value: "CC0-1.0"},
		{Tag: "SPDXID", Value: "SPDXRef-DOCUMENT"},
		{Tag: "FileName", Value: "f1"},
	}
	_, err := ParseTagValues(tvPairs)
	if err == nil {
		t.Errorf("file without SPDX Identifier getting accepted")
	}
}

func TestParser2_2PackageWithoutSpdxIdThrowErrorAtCompleteParse(t *testing.T) {
	// case: Checks the last package
	tvPairs := []reader.TagValuePair{
		{Tag: "SPDXVersion", Value: "SPDX-2.2"},
		{Tag: "DataLicense", Value: "CC0-1.0"},
		{Tag: "SPDXID", Value: "SPDXRef-DOCUMENT"},
		{Tag: "PackageName", Value: "p1"},
	}
	_, err := ParseTagValues(tvPairs)
	if err == nil {
		t.Errorf("package without SPDX Identifier getting accepted")
	}
}
