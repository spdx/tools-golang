// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"reflect"
	"testing"
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
