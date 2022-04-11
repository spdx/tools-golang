// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"errors"
	"fmt"
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	gordfWriter "github.com/spdx/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx"
)

// returns a new instance of rdfParser2_2 given the gordf object and nodeToTriples mapping
func NewParser2_2(gordfParserObj *gordfParser.Parser, nodeToTriples map[string][]*gordfParser.Triple) *rdfParser2_2 {
	parser := rdfParser2_2{
		gordfParserObj:      gordfParserObj,
		nodeStringToTriples: nodeToTriples,
		doc: &spdx.Document2_2{
			ExternalDocumentReferences: []spdx.ExternalDocumentRef2_2{},
			CreationInfo:               &spdx.CreationInfo2_2{},
			Packages:                   []*spdx.Package2_2{},
			Files:                      []*spdx.File2_2{},
			OtherLicenses:              []*spdx.OtherLicense2_2{},
			Relationships:              []*spdx.Relationship2_2{},
			Annotations:                []*spdx.Annotation2_2{},
			Reviews:                    []*spdx.Review2_2{},
		},
		files:            map[spdx.ElementID]*spdx.File2_2{},
		assocWithPackage: map[spdx.ElementID]bool{},
		cache:            map[string]*nodeState{},
	}
	return &parser
}

// main function which takes in a gordfParser and returns
// a spdxDocument model or the error encountered while parsing it
func LoadFromGoRDFParser(gordfParserObj *gordfParser.Parser) (*spdx.Document2_2, error) {
	// nodeToTriples is a mapping from a node to list of triples.
	// for every node in the set of subjects of all the triples,
	// it provides a list of triples that are associated with that subject node.
	nodeToTriples := gordfWriter.GetNodeToTriples(gordfParserObj.Triples)
	parser := NewParser2_2(gordfParserObj, nodeToTriples)

	spdxDocumentNode, err := parser.getSpdxDocNode()
	if err != nil {
		return nil, err
	}

	err = parser.parseSpdxDocumentNode(spdxDocumentNode)
	if err != nil {
		return nil, err
	}

	// parsing other root elements
	for _, rootNode := range gordfWriter.GetRootNodes(parser.gordfParserObj.Triples) {
		typeTriples := gordfWriter.FilterTriples(gordfParserObj.Triples, &rootNode.ID, &RDF_TYPE, nil)
		if len(typeTriples) != 1 {
			return nil, fmt.Errorf("every node must be associated with exactly 1 type Triple. found %d type triples", len(typeTriples))
		}
		switch typeTriples[0].Object.ID {
		case SPDX_SPDX_DOCUMENT_CAPITALIZED:
			continue // it is already parsed.
		case SPDX_SNIPPET:
			snippet, err := parser.getSnippetInformationFromNode2_2(typeTriples[0].Subject)
			if err != nil {
				return nil, fmt.Errorf("error parsing a snippet: %v", err)
			}
			err = parser.setSnippetToFileWithID(snippet, snippet.SnippetFromFileSPDXIdentifier)
			if err != nil {
				return nil, err
			}
		// todo: check other root node attributes.
		default:
			continue
			// because in rdf it is quite possible that the root node is an
			// element that has been used in the some other element as a child
		}
	}

	// parsing packages and files sets the files to a files variable which is
	// associated with the parser and not the document. following method is
	// necessary to transfer the files which are not set in the packages to the
	// Files attribute of the document
	// WARNING: do not relocate following function call. It must be at the end of the function
	parser.setUnpackagedFiles()
	return parser.doc, nil
}

// from the given parser object, returns the SpdxDocument Node defined in the root elements.
// returns error if the document is associated with no SpdxDocument or
// associated with more than one SpdxDocument node.
func (parser *rdfParser2_2) getSpdxDocNode() (node *gordfParser.Node, err error) {
	/* Possible Questions:
	1. why are you traversing the root nodes only? why not directly filter out
	   all the triples with rdf:type=spdx:SpdxDocument?
	Ans: It is quite possible that the relatedElement or any other attribute
		 to have dependency of another SpdxDocument. In that case, that
		 element will reference the dependency using SpdxDocument tag which will
		 cause false positives when direct filtering is done.
	*/
	// iterate over root nodes and find the node which has a property of rdf:type=spdx:SpdxDocument
	var spdxDocNode *gordfParser.Node
	for _, rootNode := range gordfWriter.GetRootNodes(parser.gordfParserObj.Triples) {
		typeTriples := gordfWriter.FilterTriples(
			parser.nodeToTriples(rootNode), // triples
			&rootNode.ID,                   // Subject
			&RDF_TYPE,                      // Predicate
			nil,                            // Object
		)

		if typeTriples[0].Object.ID == SPDX_SPDX_DOCUMENT_CAPITALIZED {
			// we found a SpdxDocument Node

			// must be associated with exactly one rdf:type.
			if len(typeTriples) != 1 {
				return nil, fmt.Errorf("rootNode (%v) must be associated with exactly one"+
					" triple of predicate rdf:type, found %d triples", rootNode, len(typeTriples))
			}

			// checking if we've already found a node and it is not same as the current one.
			if spdxDocNode != nil && spdxDocNode.ID != typeTriples[0].Subject.ID {
				return nil, fmt.Errorf("found more than one SpdxDocument Node (%v and %v)", spdxDocNode, typeTriples[0].Subject)
			}
			spdxDocNode = typeTriples[0].Subject
		}
	}
	if spdxDocNode == nil {
		return nil, errors.New("RDF files must be associated with a SpdxDocument tag. No tag found")
	}
	return spdxDocNode, nil
}
