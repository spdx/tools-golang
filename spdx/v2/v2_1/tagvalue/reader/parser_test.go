// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import (
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_1"
	"github.com/spdx/tools-golang/tagvalue/reader"
)

// ===== Parser exported entry point tests =====
func TestParserCanParseTagValues(t *testing.T) {
	var tvPairs []reader.TagValuePair

	// create some pairs
	tvPair1 := reader.TagValuePair{Tag: "SPDXVersion", Value: "SPDX-2.1"}
	tvPairs = append(tvPairs, tvPair1)
	tvPair2 := reader.TagValuePair{Tag: "DataLicense", Value: spdx.DataLicense}
	tvPairs = append(tvPairs, tvPair2)
	tvPair3 := reader.TagValuePair{Tag: "SPDXID", Value: "SPDXRef-DOCUMENT"}
	tvPairs = append(tvPairs, tvPair3)

	// now parse them
	doc, err := ParseTagValues(tvPairs)
	if err != nil {
		t.Errorf("got error when calling ParseTagValues: %v", err)
	}
	if doc.SPDXVersion != "SPDX-2.1" {
		t.Errorf("expected SPDXVersion to be SPDX-2.1, got %v", doc.SPDXVersion)
	}
	if doc.DataLicense != spdx.DataLicense {
		t.Errorf("expected DataLicense to be CC0-1.0, got %v", doc.DataLicense)
	}
	if doc.SPDXIdentifier != "DOCUMENT" {
		t.Errorf("expected SPDXIdentifier to be DOCUMENT, got %v", doc.SPDXIdentifier)
	}

}

// ===== Parser initialization tests =====
func TestParserInitCreatesResetStatus(t *testing.T) {
	parser := tvParser{}
	if parser.st != psStart {
		t.Errorf("parser did not begin in start state")
	}
	if parser.doc != nil {
		t.Errorf("parser did not begin with nil document")
	}
}

func TestParserHasDocumentAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser{}
	err := parser.parsePair("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.doc == nil {
		t.Errorf("doc is still nil after parsing first pair")
	}
}

func TestParserStartFailsToParseIfInInvalidState(t *testing.T) {
	parser := tvParser{st: psReview}
	err := parser.parsePairFromStart("SPDXVersion", "SPDX-2.1")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFilesWithoutSpdxIdThrowErrorAtCompleteParse(t *testing.T) {
	// case: checks the last file
	// Last unpackaged file no packages in doc
	// Last file of last package in the doc
	tvPairs := []reader.TagValuePair{
		{Tag: "SPDXVersion", Value: "SPDX-2.1"},
		{Tag: "DataLicense", Value: spdx.DataLicense},
		{Tag: "SPDXID", Value: "SPDXRef-DOCUMENT"},
		{Tag: "FileName", Value: "f1"},
	}
	_, err := ParseTagValues(tvPairs)
	if err == nil {
		t.Errorf("file without SPDX Identifier getting accepted")
	}
}

func TestParserPackageWithoutSpdxIdThrowErrorAtCompleteParse(t *testing.T) {
	// case: Checks the last package
	tvPairs := []reader.TagValuePair{
		{Tag: "SPDXVersion", Value: "SPDX-2.1"},
		{Tag: "DataLicense", Value: spdx.DataLicense},
		{Tag: "SPDXID", Value: "SPDXRef-DOCUMENT"},
		{Tag: "PackageName", Value: "p1"},
	}
	_, err := ParseTagValues(tvPairs)
	if err == nil {
		t.Errorf("package without SPDX Identifier getting accepted")
	}
}
