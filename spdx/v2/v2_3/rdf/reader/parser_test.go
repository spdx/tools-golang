// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"testing"
)

func TestNewParser2_3(t *testing.T) {
	// testing if the attributes are initialised well and no top-level is left uninitialized.
	// primarily, checking if all the maps are initialized because
	// uninitialized slices are by default slices of length 0
	p, _ := parserFromBodyContent(``)
	parser := NewParser2_3(p.gordfParserObj, p.nodeStringToTriples)
	if parser.files == nil {
		t.Errorf("files should've been initialised, got %v", parser.files)
	}
	if parser.assocWithPackage == nil {
		t.Errorf("assocWithPackage should've been initialised, got %v", parser.assocWithPackage)
	}
	if parser.doc.CreationInfo == nil {
		t.Errorf("doc.CreationInfo should've been initialised, got %v", parser.doc.CreationInfo)
	}
	if parser.doc.Packages == nil {
		t.Errorf("doc.Packages should've been initialised, got %v", parser.doc.Packages)
	}
	if parser.doc.Files == nil {
		t.Errorf("doc.Files should've been initialised, got %v", parser.doc.Files)
	}
}

func TestLoadFromGoRDFParser(t *testing.T) {
	var parser *rdfParser2_3
	var err error

	// TestCase 1: gordfparser without a SpdxDocument node triple:
	parser, _ = parserFromBodyContent("")
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err == nil {
		t.Errorf("expected an error because of absence of SpdxDocument node, got %v", err)
	}

	// TestCase 2: invalid SpdxDocumentNode
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301/Document">
			<spdx:invalidTag />
		</spdx:SpdxDocument>
	`)
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err == nil {
		t.Errorf("expected an error because of absence of SpdxDocument node, got %v", err)
	}

	// TestCase 3: >1 type triples for subnode of a SpdxDocument:
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:Snippet rdf:about="#Snippet"/>
		<spdx:CreationInfo rdf:about="#Snippet"/>
	`)
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err == nil {
		t.Errorf("expected an error due to more than one type triples, got %v", err)
	}

	// TestCase 4: invalid snippet must raise an error.
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:Snippet rdf:about="#Snippet"/>
	`)
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err == nil {
		t.Errorf("expected an error due to invalid Snippet, got %v", err)
	}

	// TestCase 5: invalid snippet not associated with any File must raise an error.
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:Snippet rdf:about="#SPDXRef-Snippet"/>
	`)
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err == nil {
		t.Errorf("expected an error due to invalid Snippet File, got %v", err)
	}

	// TestCase 6: other Tag alongwith the SpdxDocument node mustn't raise any error.
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:review/>
	`)
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// TestCase 5: everything valid:
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:name>from linux kernel</spdx:name>
			<spdx:copyrightText>Copyright 2008-2010 John Smith</spdx:copyrightText>
			<spdx:licenseComments>The concluded license was taken from package xyz, from which the snippet was copied into the current file. The concluded license information was found in the COPYING.txt file in package xyz.</spdx:licenseComments>
			<spdx:snippetFromFile>
				<spdx:File rdf:about="#SPDXRef-DoapSource">
					<spdx:copyrightText>Copyright 2010, 2011 Source Auditor Inc.</spdx:copyrightText>
					<spdx:fileContributor>Open Logic Inc.</spdx:fileContributor>
					<spdx:fileName>./src/org/spdx/parser/DOAPProject.java</spdx:fileName>
					<spdx:fileContributor>Black Duck Software In.c</spdx:fileContributor>
					<spdx:fileType rdf:resource="http://spdx.org/rdf/terms#fileType_source"/>
					<spdx:licenseInfoInFile rdf:resource="http://spdx.org/licenses/Apache-2.0"/>
				</spdx:File>
			</spdx:snippetFromFile>
		</spdx:Snippet>
	`)
	_, err = LoadFromGoRDFParser(parser.gordfParserObj)
	if err != nil {
		t.Errorf("error parsing a valid example: %v", err)
	}
}

func Test_rdfParser2_3_getSpdxDocNode(t *testing.T) {
	var parser *rdfParser2_3
	var err error

	// TestCase 1: more than one association type for a single node.
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:Snippet rdf:about="#SPDXRef-Document"/>
	`)
	_, err = parser.getSpdxDocNode()
	t.Log(err)
	if err == nil {
		t.Errorf("expected and error due to more than one type triples for the SpdxDocument Node, got %v", err)
	}

	// TestCase 2: must be associated with exactly one rdf:type.
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document"/>
		<spdx:Snippet rdf:about="#SPDXRef-Document"/>
		<spdx:File rdf:about="#SPDXRef-DoapSource"/>
	`)
	_, err = parser.getSpdxDocNode()
	t.Log(err)
	if err == nil {
		t.Errorf("rootNode  must be associated with exactly one triple of predicate rdf:type, got %v", err)
	}

	// TestCase 3: two different spdx nodes found in a single document.
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document-1"/>
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document-2"/>
	`)
	_, err = parser.getSpdxDocNode()
	if err == nil {
		t.Errorf("expected and error due to more than one type SpdxDocument Node, got %v", err)
	}

	// TestCase 4: no spdx document
	parser, _ = parserFromBodyContent(``)
	_, err = parser.getSpdxDocNode()
	if err == nil {
		t.Errorf("expected and error due to no SpdxDocument Node, got %v", err)
	}

	// TestCase 5: valid spdxDocument node
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document-1"/>
	`)
	_, err = parser.getSpdxDocNode()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
