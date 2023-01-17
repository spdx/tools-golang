// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"
	"strconv"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// Snippet Information
// Cardinality: Optional, Many
func (parser *rdfParser2_3) getSnippetInformationFromNode2_3(node *gordfParser.Node) (si *spdx.Snippet, err error) {
	si = &spdx.Snippet{}

	err = setSnippetID(node.ID, si)
	if err != nil {
		return nil, err
	}

	for _, siTriple := range parser.nodeToTriples(node) {
		switch siTriple.Predicate.ID {
		case RDF_TYPE:
			// cardinality: exactly 1
		case SPDX_SNIPPET_FROM_FILE:
			// cardinality: exactly 1
			// file which is associated with the snippet
			_, err := parser.getFileFromNode(siTriple.Object)
			if err != nil {
				return nil, err
			}
			docElemID, err := ExtractDocElementID(getLastPartOfURI(siTriple.Object.ID))
			si.SnippetFromFileSPDXIdentifier = docElemID.ElementRefID
		case SPDX_RANGE:
			// cardinality: min 1
			err = parser.setSnippetRangeFromNode(siTriple.Object, si)
			if err != nil {
				return nil, err
			}
		case SPDX_LICENSE_INFO_IN_SNIPPET:
			// license info in snippet can be NONE, NOASSERTION or SimpleLicensingInfo
			// using AnyLicenseInfo because it can redirect the request and
			// can handle NONE & NOASSERTION
			var anyLicense AnyLicenseInfo
			anyLicense, err = parser.getAnyLicenseFromNode(siTriple.Object)
			if err != nil {
				return nil, fmt.Errorf("error parsing license info in snippet: %v", err)
			}
			si.LicenseInfoInSnippet = append(si.LicenseInfoInSnippet, anyLicense.ToLicenseString())
		case SPDX_NAME:
			si.SnippetName = siTriple.Object.ID
		case SPDX_COPYRIGHT_TEXT:
			si.SnippetCopyrightText = siTriple.Object.ID
		case SPDX_LICENSE_COMMENTS:
			si.SnippetLicenseComments = siTriple.Object.ID
		case RDFS_COMMENT:
			si.SnippetComment = siTriple.Object.ID
		case SPDX_LICENSE_CONCLUDED:
			var anyLicense AnyLicenseInfo
			anyLicense, err = parser.getAnyLicenseFromNode(siTriple.Object)
			if err != nil {
				return nil, fmt.Errorf("error parsing license info in snippet: %v", err)
			}
			si.SnippetLicenseConcluded = anyLicense.ToLicenseString()
		default:
			return nil, fmt.Errorf("unknown predicate %v", siTriple.Predicate.ID)
		}
	}
	return si, nil
}

// given is the id of the file, sets the snippet to the file in parser.
func (parser *rdfParser2_3) setSnippetToFileWithID(snippet *spdx.Snippet, fileID common.ElementID) error {
	if parser.files[fileID] == nil {
		return fmt.Errorf("snippet refers to an undefined file with ID: %s", fileID)
	}

	// initializing snippet of the files if it is not defined already
	if parser.files[fileID].Snippets == nil {
		parser.files[fileID].Snippets = map[common.ElementID]*spdx.Snippet{}
	}

	// setting the snippet to the file.
	parser.files[fileID].Snippets[snippet.SnippetSPDXIdentifier] = snippet

	return nil
}

func (parser *rdfParser2_3) setSnippetRangeFromNode(node *gordfParser.Node, si *spdx.Snippet) error {
	// for a range object, we can have only 3 associated triples:
	//		node -> RDF_TYPE     -> Object
	//      node -> startPointer -> Object
	//      node -> endPointer   -> Object
	associatedTriples := parser.nodeToTriples(node)
	if len(associatedTriples) != 3 {
		return fmt.Errorf("range should be associated with exactly 3 triples, got %d", len(associatedTriples))
	}

	// Triple 1: Predicate=RDF_TYPE
	typeTriple := rdfwriter.FilterTriples(associatedTriples, &node.ID, &RDF_TYPE, nil)
	if len(typeTriple) != 1 {
		// we had 3 associated triples. out of which 2 is start and end pointer,
		// if we do not have the rdf:type triple as the third one,
		// we have either extra or undefined predicate.
		return fmt.Errorf("every object node must be associated with exactly one rdf:type triple, found: %d", len(typeTriple))
	}

	// getting start pointer
	startPointerTriples := rdfwriter.FilterTriples(associatedTriples, &node.ID, &PTR_START_POINTER, nil)
	if len(startPointerTriples) != 1 {
		return fmt.Errorf("range object must be associated with exactly 1 startPointer, got %d", len(startPointerTriples))
	}
	startRangeType, start, err := parser.getPointerFromNode(startPointerTriples[0].Object, si)
	if err != nil {
		return fmt.Errorf("error parsing startPointer: %v", err)
	}

	// getting end pointer
	endPointerTriples := rdfwriter.FilterTriples(associatedTriples, &node.ID, &PTR_END_POINTER, nil)
	if len(startPointerTriples) != 1 {
		return fmt.Errorf("range object must be associated with exactly 1 endPointer, got %d", len(endPointerTriples))
	}
	endRangeType, end, err := parser.getPointerFromNode(endPointerTriples[0].Object, si)
	if err != nil {
		return fmt.Errorf("error parsing endPointer: %v", err)
	}

	// return error when start and end pointer type is not same.
	if startRangeType != endRangeType {
		return fmt.Errorf("start and end range type doesn't match")
	}

	si.Ranges = []common.SnippetRange{{
		StartPointer: common.SnippetRangePointer{FileSPDXIdentifier: si.SnippetFromFileSPDXIdentifier},
		EndPointer:   common.SnippetRangePointer{FileSPDXIdentifier: si.SnippetFromFileSPDXIdentifier},
	}}

	if startRangeType == LINE_RANGE {
		si.Ranges[0].StartPointer.LineNumber = start
		si.Ranges[0].EndPointer.LineNumber = end
	} else {
		si.Ranges[0].StartPointer.Offset = start
		si.Ranges[0].EndPointer.Offset = end
	}
	return nil
}

func (parser *rdfParser2_3) getPointerFromNode(node *gordfParser.Node, si *spdx.Snippet) (rt RangeType, number int, err error) {
	for _, triple := range parser.nodeToTriples(node) {
		switch triple.Predicate.ID {
		case RDF_TYPE:
		case PTR_REFERENCE:
			err = parser.parseRangeReference(triple.Object, si)
		case PTR_OFFSET:
			number, err = strconv.Atoi(triple.Object.ID)
			rt = BYTE_RANGE
		case PTR_LINE_NUMBER:
			number, err = strconv.Atoi(triple.Object.ID)
			rt = LINE_RANGE
		default:
			err = fmt.Errorf("undefined predicate (%s) for a pointer", triple.Predicate)
		}
		if err != nil {
			return
		}
	}
	if rt == "" {
		err = fmt.Errorf("range type not defined for a pointer")
	}
	return
}

func (parser *rdfParser2_3) parseRangeReference(node *gordfParser.Node, snippet *spdx.Snippet) error {
	// reference is supposed to be either a resource reference to an already
	// defined or a new file. Unfortunately, I didn't find field where this can be set in the tools-golang data model.
	// todo: set this reference to the snippet
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	if len(associatedTriples) == 0 {
		return nil
	}
	_, err := parser.getFileFromNode(node)
	if err != nil {
		return fmt.Errorf("error parsing a new file in a reference: %v", err)
	}
	return nil
}

func setSnippetID(uri string, si *spdx.Snippet) (err error) {
	fragment := getLastPartOfURI(uri)
	si.SnippetSPDXIdentifier, err = ExtractElementID(fragment)
	if err != nil {
		return fmt.Errorf("error setting snippet identifier: %v", uri)
	}
	return nil
}
