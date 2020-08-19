// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"errors"
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	gordfWriter "github.com/RishabhBhatnagar/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx"
)

// returns a new instance of rdfParser2_2 given the gordf object and nodeToTriples mapping
func NewParser2_2(gordfParserObj *gordfParser.Parser, nodeToTriples map[string][]*gordfParser.Triple) *rdfParser2_2 {
	parser := rdfParser2_2{
		gordfParserObj: gordfParserObj,
		nodeToTriples:  nodeToTriples,
		doc: &spdx.Document2_2{
			CreationInfo:    &spdx.CreationInfo2_2{},
			Packages:        map[spdx.ElementID]*spdx.Package2_2{},
			UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{},
			OtherLicenses:   []*spdx.OtherLicense2_2{},
			Relationships:   []*spdx.Relationship2_2{},
			Annotations:     []*spdx.Annotation2_2{},
			Reviews:         []*spdx.Review2_2{},
		},
		files:            map[spdx.ElementID]*spdx.File2_2{},
		packages:         map[spdx.ElementID]*spdx.Package2_2{},
		assocWithPackage: map[spdx.ElementID]bool{},
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
			snippet, err := parser.getSnippetInformationFromTriple2_2(typeTriples[0])
			if err != nil {
				return nil, err
			}
			err = parser.setSnippetToFileWithID(snippet, snippet.SnippetFromFileSPDXIdentifier.ElementRefID)
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
	// UnpackagedFiles attribute of the document
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
			parser.nodeToTriples[rootNode.String()], // triples
			&rootNode.ID,                            // Subject
			&RDF_TYPE,                               // Predicate
			nil,                                     // Object
		)
		if len(typeTriples) != 1 {
			return nil, fmt.Errorf("rootNode (%v) must be associated with exactly one"+
				" triple of predicate rdf:type, found %d triples", rootNode, len(typeTriples))
		}
		if typeTriples[0].Object.ID == SPDX_SPDX_DOCUMENT_CAPITALIZED {
			// we found a SpdxDocument Node
			// checking if we've already found a node and it is not same as the current one.
			if spdxDocNode != nil && spdxDocNode.ID != typeTriples[0].Subject.ID {
				return nil, fmt.Errorf("found more than one SpdxDocument Node (%v and %v)", spdxDocNode, typeTriples[0].Subject)
			}
			spdxDocNode = typeTriples[0].Subject
		}
	}
	if spdxDocNode == nil {
		return nil, fmt.Errorf("RDF files must be associated with a SpdxDocument tag. No tag found")
	}
	return spdxDocNode, nil
}

// unused
func (parser *rdfParser2_2) setFiles() error {
	allFilesTriples, err := parser.filterTriplesByRegex(parser.gordfParserObj.Triples, ".*", RDF_TYPE+"$", SPDX_FILE+"$")
	if err != nil {
		return err
	}
	for _, fileTriple := range allFilesTriples {
		file, err := parser.getFileFromNode(fileTriple.Subject)
		if err != nil {
			return fmt.Errorf("error setting a file: %v", err)
		}
		parser.files[file.FileSPDXIdentifier] = file
	}
	return nil
}

// unused
func (parser *rdfParser2_2) setPackages() error {
	allPackagesTriples, err := parser.filterTriplesByRegex(parser.gordfParserObj.Triples, ".*", RDF_TYPE+"$", SPDX_PACKAGE+"$")
	if err != nil {
		return err
	}
	for _, pkgTriple := range allPackagesTriples {
		pkg, err := parser.getPackageFromNode(pkgTriple.Subject)
		if err != nil {
			return fmt.Errorf("error setting a package: %v", err)
		}
		parser.packages[pkg.PackageSPDXIdentifier] = pkg
	}
	return nil
}

// unused
// assumes that the document's namespace is already set.
func (parser *rdfParser2_2) setSnippetToDoc(si *spdx.Snippet2_2, snippetNode *gordfParser.Node) (err error) {
	if parser.doc == nil || parser.doc.CreationInfo == nil {
		return errors.New("document namespace not set yet")
	}
	docNS := parser.doc.CreationInfo.DocumentNamespace
	snippetNS := stripLastPartOfUri(snippetNode.ID)
	if !isUriSame(docNS, snippetNS) {
		// found a snippet which doesn't belong to current document being set
		return fmt.Errorf("document namespace(%s) and snippet namespace(%s) doesn't match", docNS, snippetNS)
	}

	return nil
}

// unused
func (parser *rdfParser2_2) setAnnotations(spdxDocNode *gordfParser.Node) error {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_ANNOTATION, nil)
	for _, triple := range triples {
		err := parser.parseAnnotationFromNode(triple.Object)
		if err != nil {
			return err
		}
	}
	return nil
}

// unused
func (parser *rdfParser2_2) getSpecVersion(spdxDocNode *gordfParser.Node) (string, error) {
	specVersionTriples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_SPEC_VERSION, nil)
	n := len(specVersionTriples)
	if n == 0 {
		return "", fmt.Errorf("no specVersion found for the given spdxNode")
	}
	if n > 1 {
		return "", fmt.Errorf("there must be exactly one specVersion. found %d specVersion", n)
	}
	return specVersionTriples[0].Object.ID, nil
}

// unused
func (parser *rdfParser2_2) getDataLicense(spdxDocNode *gordfParser.Node) (string, error) {
	dataLicenseTriples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_DATA_LICENSE, nil)
	n := len(dataLicenseTriples)
	if n == 0 {
		return "", fmt.Errorf("no dataLicense found for the given spdxNode")
	}
	if n > 1 {
		return "", fmt.Errorf("there must be exactly one dataLicense. found %d dataLicense", n)
	}
	return parser.getLicenseFromTriple(dataLicenseTriples[0])
}

// unused
func (parser *rdfParser2_2) getDocumentName(spdxDocNode *gordfParser.Node) (string, error) {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_NAME, nil)
	n := len(triples)
	if n == 0 {
		return "", fmt.Errorf("no documentName found for the given spdxNode")
	}
	if n > 1 {
		return "", fmt.Errorf("there must be exactly one documentName. found %d documentName", n)
	}
	return triples[0].Object.ID, nil
}

// unused
func (parser *rdfParser2_2) setCreationInfo(spdxDocNode *gordfParser.Node, ci *spdx.CreationInfo2_2) error {
	docNS := stripJoiningChars(stripLastPartOfUri(spdxDocNode.ID))
	allCreationInfoTriples, err := parser.filterTriplesByRegex(parser.gordfParserObj.Triples, docNS+".*", SPDX_CREATION_INFO, ".*")
	if err != nil {
		return err
	}
	n := len(allCreationInfoTriples)
	if n > 1 {
		return fmt.Errorf("document(%s) must have exactly one creation info. found %d", docNS, n)
	}
	if n == 0 {
		return fmt.Errorf("no creation info found for the document identified by %s", docNS)
	}
	err = parser.parseCreationInfoFromNode(ci, allCreationInfoTriples[0].Object)
	if err != nil {
		return err
	}
	parser.doc.CreationInfo = ci
	return nil
}

// unused
func (parser *rdfParser2_2) getDocumentComment(spdxDocNode *gordfParser.Node) (string, error) {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_NAME, nil)
	n := len(triples)
	if n > 1 {
		return "", fmt.Errorf("there must be atmost one documentComment. found %d documentComment", n)
	}
	if n == 0 {
		return triples[0].Object.ID, nil
	}
	return "", nil
}

// unused
func (parser *rdfParser2_2) setReviewed(spdxDocNode *gordfParser.Node) error {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_REVIEWED, nil)
	for _, triple := range triples {
		err := parser.setReviewFromNode(triple.Object)
		if err != nil {
			return err
		}
	}
	return nil
}

// unused
func (parser *rdfParser2_2) setDescribesPackage(spdxDocNode *gordfParser.Node) error {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_DESCRIBES_PACKAGE, nil)
	for _, triple := range triples {
		pkg, err := parser.getPackageFromNode(triple.Object)
		if err != nil {
			return err
		}
		parser.doc.Packages[pkg.PackageSPDXIdentifier] = pkg
	}
	return nil
}

// unused
func (parser *rdfParser2_2) setExtractedLicensingInfo(spdxDocNode *gordfParser.Node) error {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_HAS_EXTRACTED_LICENSING_INFO, nil)
	for _, triple := range triples {
		err := parser.parseOtherLicenseFromNode(triple.Object)
		if err != nil {
			return err
		}
	}
	return nil
}

// unused
func (parser *rdfParser2_2) setRelationships(spdxDocNode *gordfParser.Node) error {
	triples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, &spdxDocNode.ID, &SPDX_RELATIONSHIP, nil)
	for _, triple := range triples {
		err := parser.parseRelationship(triple)
		if err != nil {
			return err
		}
	}
	return nil
}
