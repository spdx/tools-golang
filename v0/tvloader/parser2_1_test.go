// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"testing"

	"github.com/swinslow/spdx-go/v0/spdx"
)

func TestParser2_1InitCreatesResetStatus(t *testing.T) {
	parser := tvParser2_1{}
	if parser.st != psStart2_1 {
		t.Errorf("parser did not begin in start state")
	}
	if parser.doc != nil {
		t.Errorf("parser did not begin with nil document")
	}
}

func TestParser2_1HasDocumentAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_1{}
	err := parser.parsePair2_1("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.doc == nil {
		t.Errorf("doc is still nil after parsing first pair")
	}
}

func TestParser2_1MovesToCreationInfoStateAfterParsingFirstTag(t *testing.T) {
	parser := tvParser2_1{}
	err := parser.parsePairFromStart2_1("a", "b")
	if err != nil {
		t.Errorf("got error when calling parsePairFromStart2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v after parsing first pair", parser.st)
	}
}

func TestParser2_1HasCreationInfoAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePairFromCreationInfo2_1("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("got error when calling parsePairFromCreationInfo2_1: %v", err)
	}
	if parser.doc.CreationInfo == nil {
		t.Errorf("doc.CreationInfo is still nil after parsing first pair")
	}
}

func TestCanExtractSubvalues(t *testing.T) {
	subkey, subvalue, err := extractSubs("SHA1: abc123")
	if err != nil {
		t.Errorf("got error when calling extractSubs: %v", err)
	}
	if subkey != "SHA1" {
		t.Errorf("got %v for subkey", subkey)
	}
	if subvalue != "abc123" {
		t.Errorf("got %v for subvalue", subvalue)
	}
}

func TestReturnsErrorForInvalidSubvalueFormat(t *testing.T) {
	_, _, err := extractSubs("blah")
	if err == nil {
		t.Errorf("expected error when calling extractSubs for invalid format (0 colons), got nil")
	}

	_, _, err = extractSubs("MD5: 12390834: other")
	if err == nil {
		t.Errorf("expected error when calling extractSubs for invalid format (>1 colon), got nil")
	}
}

func TestParser2_1CanParseCreationInfoTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// SPDX Version
	err := parser.parsePairFromCreationInfo2_1("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.SPDXVersion != "SPDX-2.1" {
		t.Errorf("got %v for SPDXVersion", parser.doc.CreationInfo.SPDXVersion)
	}

	// Data License
	err = parser.parsePairFromCreationInfo2_1("DataLicense", "CC0-1.0")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.DataLicense != "CC0-1.0" {
		t.Errorf("got %v for DataLicense", parser.doc.CreationInfo.DataLicense)
	}

	// SPDX Identifier
	err = parser.parsePairFromCreationInfo2_1("SPDXID", "SPDXRef-DOCUMENT")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.SPDXIdentifier != "SPDXRef-DOCUMENT" {
		t.Errorf("got %v for SPDXIdentifier", parser.doc.CreationInfo.SPDXIdentifier)
	}

	// Document Name
	err = parser.parsePairFromCreationInfo2_1("DocumentName", "xyz-2.1.5")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.DocumentName != "xyz-2.1.5" {
		t.Errorf("got %v for DocumentName", parser.doc.CreationInfo.DocumentName)
	}

	// Document Namespace
	err = parser.parsePairFromCreationInfo2_1("DocumentNamespace", "http://example.com/xyz-2.1.5.spdx")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.DocumentNamespace != "http://example.com/xyz-2.1.5.spdx" {
		t.Errorf("got %v for DocumentNamespace", parser.doc.CreationInfo.DocumentNamespace)
	}

	// External Document Reference
	refs := []string{
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 SHA1: d6a770ba38583ed4bb4525bd96e50461655d2759",
		"DocumentRef-xyz-2.1.2 http://example.com/xyz-2.1.2 SHA1: d6a770ba38583ed4bb4525bd96e50461655d2760",
	}
	err = parser.parsePairFromCreationInfo2_1("ExternalDocumentRef", refs[0])
	if err != nil {
		t.Errorf("nil returned")
	}
	err = parser.parsePairFromCreationInfo2_1("ExternalDocumentRef", refs[1])
	if err != nil {
		t.Errorf("nil returned")
	}
	if len(parser.doc.CreationInfo.ExternalDocumentReferences) != 2 ||
		parser.doc.CreationInfo.ExternalDocumentReferences[0] != refs[0] ||
		parser.doc.CreationInfo.ExternalDocumentReferences[1] != refs[1] {
		t.Errorf("got %v for ExternalDocumentReferences", parser.doc.CreationInfo.ExternalDocumentReferences)
	}

	// License List Version
	err = parser.parsePairFromCreationInfo2_1("LicenseListVersion", "2.2")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.LicenseListVersion != "2.2" {
		t.Errorf("got %v for LicenseListVersion", parser.doc.CreationInfo.LicenseListVersion)
	}

	// Creators: Persons
	refPersons := []string{
		"Person: Person A",
		"Person: Person B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refPersons[0])
	if err != nil {
		t.Errorf("nil returned")
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refPersons[1])
	if err != nil {
		t.Errorf("nil returned")
	}
	if len(parser.doc.CreationInfo.CreatorPersons) != 2 ||
		parser.doc.CreationInfo.CreatorPersons[0] != "Person A" ||
		parser.doc.CreationInfo.CreatorPersons[1] != "Person B" {
		t.Errorf("got %v for CreatorPersons", parser.doc.CreationInfo.CreatorPersons)
	}

	// Creators: Organizations
	refOrgs := []string{
		"Organization: Organization A",
		"Organization: Organization B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refOrgs[0])
	if err != nil {
		t.Errorf("nil returned")
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refOrgs[1])
	if err != nil {
		t.Errorf("nil returned")
	}
	if len(parser.doc.CreationInfo.CreatorOrganizations) != 2 ||
		parser.doc.CreationInfo.CreatorOrganizations[0] != "Organization A" ||
		parser.doc.CreationInfo.CreatorOrganizations[1] != "Organization B" {
		t.Errorf("got %v for CreatorOrganizations", parser.doc.CreationInfo.CreatorOrganizations)
	}

	// Creators: Tools
	refTools := []string{
		"Tool: Tool A",
		"Tool: Tool B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refTools[0])
	if err != nil {
		t.Errorf("nil returned")
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refTools[1])
	if err != nil {
		t.Errorf("nil returned")
	}
	if len(parser.doc.CreationInfo.CreatorTools) != 2 ||
		parser.doc.CreationInfo.CreatorTools[0] != "Tool A" ||
		parser.doc.CreationInfo.CreatorTools[1] != "Tool B" {
		t.Errorf("got %v for CreatorTools", parser.doc.CreationInfo.CreatorTools)
	}

	// Created date
	err = parser.parsePairFromCreationInfo2_1("Created", "2018-09-10T11:46:00Z")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.Created != "2018-09-10T11:46:00Z" {
		t.Errorf("got %v for Created", parser.doc.CreationInfo.Created)
	}

	// Creator Comment
	err = parser.parsePairFromCreationInfo2_1("CreatorComment", "Blah whatever")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.CreatorComment != "Blah whatever" {
		t.Errorf("got %v for CreatorComment", parser.doc.CreationInfo.CreatorComment)
	}

	// Document COmment
	err = parser.parsePairFromCreationInfo2_1("DocumentComment", "Blah whatever")
	if err != nil {
		t.Errorf("nil returned")
	}
	if parser.doc.CreationInfo.DocumentComment != "Blah whatever" {
		t.Errorf("got %v for DocumentComment", parser.doc.CreationInfo.DocumentComment)
	}

}
