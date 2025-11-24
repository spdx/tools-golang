// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"reflect"
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

func Test_getLastPartOfURI(t *testing.T) {
	// uri of type baseFragment#fragment
	input := "baseFragment#fragment"
	expectedOutput := "fragment"
	output := getLastPartOfURI(input)
	if output != expectedOutput {
		t.Errorf("expected %s, found %s", expectedOutput, output)
	}

	// uri of type baseFragment/subFragment
	input = "baseFragment/subFragment"
	expectedOutput = "subFragment"
	output = getLastPartOfURI(input)
	if output != expectedOutput {
		t.Errorf("expected %s, found %s", expectedOutput, output)
	}

	// neither of the case mustn't raise any error.
	input = "www.github.com"
	expectedOutput = input
	output = getLastPartOfURI(input)
	if output != expectedOutput {
		t.Errorf("expected %s, found %s", expectedOutput, output)
	}
}

func Test_isUriValid(t *testing.T) {
	// TestCase 1: Valid Input URI
	input := "https://www.github.com"
	isValid := isUriValid(input)
	if !isValid {
		t.Errorf("valid input(%s) detected as invalid.", input)
	}

	// TestCase 2: Invalid Input URI
	input = `http\:www.github.com`
	isValid = isUriValid(input)
	if isValid {
		t.Errorf("invalid input(%s) detected as valid", input)
	}
}

func Test_rdfParser2_2_nodeToTriples(t *testing.T) {
	var parser *rdfParser2_2
	var output, expectedOutput []*gordfParser.Triple

	// TestCase 1: a nil node shouldn't raise any error or panic.
	parser, _ = parserFromBodyContent(``)
	output = parser.nodeToTriples(nil)
	if output == nil {
		t.Errorf("nil input should return an empty slice and not nil")
	}
	expectedOutput = []*gordfParser.Triple{}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %+v, got %+v", expectedOutput, output)
	}

	// TestCase 2: node should be addressable based on the node content and not the pointer.
	// It should allow new nodes same as the older ones to retrieve the associated triples.
	parser, _ = parserFromBodyContent(`
		  <spdx:Checksum rdf:about="#checksum">
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1" />
			<spdx:checksumValue>75068c26abbed3ad3980685bae21d7202d288317</spdx:checksumValue>
		  </spdx:Checksum>
	`)
	newNode := &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#checksum",
	}
	output = parser.nodeToTriples(newNode)

	// The output must have 3 triples:
	// 1. newNode rdf:type Checksum
	// 2. newNode spdx:algorithm http://spdx.org/rdf/terms#checksumAlgorithm_sha1
	// 3. newNode spdx:checksumValue 75068c26abbed3ad3980685bae21d7202d288317
	if len(output) != 3 {
		t.Errorf("expected output to have 3 triples, got %d", len(output))
	}
}

func Test_boolFromString(t *testing.T) {
	// TestCase 1: Valid Input: "true"
	// mustn't raise any error
	input := "true"
	val, err := boolFromString(input)
	if err != nil {
		t.Errorf("function raised an error for a valid input(%s): %s", input, err)
	}
	if val != true {
		t.Errorf("invalid output. Expected %v, found %v", true, val)
	}

	// TestCase 2: Valid Input: "true"
	// mustn't raise any error
	input = "false"
	val, err = boolFromString(input)
	if err != nil {
		t.Errorf("function raised an error for a valid input(%s): %s", input, err)
	}
	if val != false {
		t.Errorf("invalid output. Expected %v, found %v", false, val)
	}

	// TestCase 3: invalid input: ""
	// it must raise an error
	input = ""
	val, err = boolFromString(input)
	if err == nil {
		t.Errorf("invalid input should've raised an error")
	}
}

func Test_getNodeTypeFromTriples(t *testing.T) {
	var err error
	var node *gordfParser.Node
	var triples []*gordfParser.Triple
	var nodeType, expectedNodeType string

	// TestCase 1: nil node must raise an error because,
	// nil nodes cannot be associated with any rdf:type attribute.
	_, err = getNodeTypeFromTriples(triples, nil)
	if err == nil {
		t.Errorf("expected an error due to nil node, got %v", err)
	}

	// TestCase 2: none of the triples give information about the rdf:type of a node.
	node = &gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       "N0",
	}
	_, err = getNodeTypeFromTriples(triples, node)
	if err == nil {
		t.Errorf("expected an error saying no rdf:type found, got %v", err)
	}

	// TestCase 3: node is associated with exactly one rdf:type triples
	typeTriple := &gordfParser.Triple{
		Subject: node,
		Predicate: &gordfParser.Node{
			NodeType: gordfParser.IRI,
			ID:       RDF_TYPE,
		},
		Object: &gordfParser.Node{
			NodeType: gordfParser.IRI,
			ID:       "http://spdx.org/rdf/terms#Checksum",
		},
	}
	triples = append(triples, typeTriple)
	expectedNodeType = "http://spdx.org/rdf/terms#Checksum"
	nodeType, err = getNodeTypeFromTriples(triples, node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nodeType != expectedNodeType {
		t.Errorf("expected: %v, got: %v", nodeType, expectedNodeType)
	}

	// TestCase 4: node associated with more than one rdf:type triples must raise an error.
	typeTriple = &gordfParser.Triple{
		Subject: node,
		Predicate: &gordfParser.Node{
			NodeType: gordfParser.IRI,
			ID:       RDF_TYPE,
		},
		Object: &gordfParser.Node{
			NodeType: gordfParser.IRI,
			ID:       "http://spdx.org/rdf/terms#Snippet",
		},
	}
	triples = append(triples, typeTriple)
	_, err = getNodeTypeFromTriples(triples, node)
	if err == nil {
		t.Errorf("expected an error saying more than one rdf:type found, got %v", err)
	}
}

// following tests are copy pasted from tvloader/parser2v2/util_test.go

func TestCanExtractDocumentAndElementRefsFromID(t *testing.T) {
	// test with valid ID in this document
	helperForExtractDocElementID(t, "SPDXRef-file1", false, "", "file1")
	// test with valid ID in another document
	helperForExtractDocElementID(t, "DocumentRef-doc2:SPDXRef-file2", false, "doc2", "file2")
	// test with invalid ID in this document
	helperForExtractDocElementID(t, "a:SPDXRef-file1", true, "", "")
	helperForExtractDocElementID(t, "file1", true, "", "")
	helperForExtractDocElementID(t, "SPDXRef-", true, "", "")
	helperForExtractDocElementID(t, "SPDXRef-file1:", true, "", "")
	// test with invalid ID in another document
	helperForExtractDocElementID(t, "DocumentRef-doc2", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-doc2:", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-doc2:SPDXRef-", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-doc2:a", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-:", true, "", "")
	helperForExtractDocElementID(t, "DocumentRef-:SPDXRef-file1", true, "", "")
	// test with invalid formats
	helperForExtractDocElementID(t, "DocumentRef-doc2:SPDXRef-file1:file2", true, "", "")
}

func helperForExtractDocElementID(t *testing.T, tst string, wantErr bool, wantDoc string, wantElt string) {
	deID, err := ExtractDocElementID(tst)
	if err != nil && wantErr == false {
		t.Errorf("testing %v: expected nil error, got %v", tst, err)
	}
	if err == nil && wantErr == true {
		t.Errorf("testing %v: expected non-nil error, got nil", tst)
	}
	if deID.DocumentRefID != common.DocumentID(wantDoc) {
		if wantDoc == "" {
			t.Errorf("testing %v: want empty string for DocumentRefID, got %v", tst, deID.DocumentRefID)
		} else {
			t.Errorf("testing %v: want %v for DocumentRefID, got %v", tst, wantDoc, deID.DocumentRefID)
		}
	}
	if deID.ElementRefID != common.ElementID(wantElt) {
		if wantElt == "" {
			t.Errorf("testing %v: want emptyString for ElementRefID, got %v", tst, deID.ElementRefID)
		} else {
			t.Errorf("testing %v: want %v for ElementRefID, got %v", tst, wantElt, deID.ElementRefID)
		}
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
	eID, err := ExtractElementID(tst)
	if err != nil && wantErr == false {
		t.Errorf("testing %v: expected nil error, got %v", tst, err)
	}
	if err == nil && wantErr == true {
		t.Errorf("testing %v: expected non-nil error, got nil", tst)
	}
	if eID != common.ElementID(wantElt) {
		if wantElt == "" {
			t.Errorf("testing %v: want emptyString for ElementRefID, got %v", tst, eID)
		} else {
			t.Errorf("testing %v: want %v for ElementRefID, got %v", tst, wantElt, eID)
		}
	}
}

func TestCanExtractSubvalues(t *testing.T) {
	subkey, subvalue, err := ExtractSubs("SHA1: abc123", ":")
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
	_, _, err := ExtractSubs("blah", ":")
	if err == nil {
		t.Errorf("expected error when calling extractSubs for invalid format (0 colons), got nil")
	}
}
