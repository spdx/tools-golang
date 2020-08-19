// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
	"strconv"
	"strings"
)

// Snippet Information
// Cardinality: Optional, Many
func (parser *rdfParser2_2) getSnippetInformationFromTriple2_2(triple *gordfParser.Triple) (si *spdx.Snippet2_2, err error) {
	si = &spdx.Snippet2_2{}

	err = setSnippetID(triple.Subject.ID, si)
	if err != nil {
		return nil, err
	}

	for _, siTriple := range parser.nodeToTriples[triple.Subject.String()] {
		switch siTriple.Predicate.ID {
		case RDF_TYPE:
			// cardinality: exactly 1
		case SPDX_SNIPPET_FROM_FILE:
			// cardinality: exactly 1
			// file which is associated with the snippet
			file, err := parser.getFileFromNode(siTriple.Object)
			if err != nil {
				return nil, err
			}
			si.SnippetFromFileSPDXIdentifier, err = ExtractDocElementID(getLastPartOfURI(siTriple.Object.ID))
			parser.files[file.FileSPDXIdentifier] = file
		case SPDX_NAME:
			si.SnippetName = siTriple.Object.ID
		case SPDX_COPYRIGHT_TEXT:
			si.SnippetCopyrightText = siTriple.Object.ID
		case SPDX_LICENSE_COMMENTS:
			si.SnippetLicenseComments = siTriple.Object.ID
		case SPDX_LICENSE_INFO_IN_SNIPPET:
			si.LicenseInfoInSnippet = append(si.LicenseInfoInSnippet, siTriple.Object.ID)
		case RDFS_COMMENT:
			si.SnippetComment = siTriple.Object.ID
		case SPDX_LICENSE_CONCLUDED:
			si.SnippetLicenseConcluded = siTriple.Object.ID
		case SPDX_RANGE:
			// cardinality: min 1
			err = parser.setSnippetRangeFromNode(siTriple.Object, si)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown predicate %v", siTriple.Predicate.ID)
		}
	}
	return si, nil
}

// given is the id of the file, sets the snippet to the file in parser.
func (parser *rdfParser2_2) setSnippetToFileWithID(snippet *spdx.Snippet2_2, fileID spdx.ElementID) error {
	if parser.files[fileID] == nil {
		return fmt.Errorf("snippet refers to an undefined file with ID: %s", fileID)
	}

	// initializing snippet of the files if it is not defined already
	if parser.files[fileID].Snippets == nil {
		parser.files[fileID].Snippets = map[spdx.ElementID]*spdx.Snippet2_2{}
	}

	// setting the snippet to the file.
	parser.files[fileID].Snippets[snippet.SnippetSPDXIdentifier] = snippet

	return nil
}

func (parser *rdfParser2_2) setSnippetRangeFromNode(node *gordfParser.Node, si *spdx.Snippet2_2) error {
	// todo: apply DRY in this method.
	rangeType := 0 // 0: undefined range, 1: byte, 2: line
	var start, end string
	for _, t := range parser.nodeToTriples[node.String()] {
		switch t.Predicate.ID {
		case RDF_TYPE:
			if t.Object.ID != PTR_START_END_POINTER {
				return fmt.Errorf("expected range to have sub tag of type StartEndPointer, found %v", t.Object.ID)
			}
		case PTR_START_POINTER:
			for _, subTriple := range parser.nodeToTriples[t.Object.String()] {
				switch subTriple.Predicate.ID {
				case RDF_TYPE:
					switch subTriple.Object.ID {
					case PTR_BYTE_OFFSET_POINTER:
						if rangeType == 2 {
							return fmt.Errorf("byte offset pointer merged with line offset pointer")
						}
						rangeType = 1
					case PTR_LINE_CHAR_POINTER:
						if rangeType == 1 {
							return fmt.Errorf("byte offset pointer merged with line offset pointer")
						}
						rangeType = 2
					default:
						return fmt.Errorf("illegal pointer type %v", subTriple.Object.ID)
					}
				case PTR_REFERENCE:
					err := parser.parseRangeReference(subTriple.Object, si)
					if err != nil {
						return nil
					}
				case PTR_OFFSET, PTR_LINE_NUMBER:
					start = subTriple.Object.ID
				default:
					return fmt.Errorf("undefined predicate %v while parsing range", subTriple.Predicate.ID)
				}
			}
		case PTR_END_POINTER:
			subTriples := parser.nodeToTriples[t.Object.String()]
			for _, subTriple := range subTriples {
				switch subTriple.Predicate.ID {
				case RDF_TYPE:
					switch subTriple.Object.ID {
					case PTR_BYTE_OFFSET_POINTER:
						if rangeType == 2 {
							return fmt.Errorf("byte offset pointer merged with line offset pointer")
						}
						rangeType = 1
					case PTR_LINE_CHAR_POINTER:
						if rangeType == 1 {
							return fmt.Errorf("byte offset pointer merged with line offset pointer")
						}
						rangeType = 2
					default:
						return fmt.Errorf("illegal pointer type %v", subTriple.Object.ID)
					}
				case PTR_REFERENCE:
					err := parser.parseRangeReference(subTriple.Object, si)
					if err != nil {
						return nil
					}
				case PTR_OFFSET, PTR_LINE_NUMBER:
					end = subTriple.Object.ID
				}
			}
		default:
			return fmt.Errorf("unknown predicate %v", t.Predicate.ID)
		}
	}
	if rangeType != 1 && rangeType != 2 {
		return fmt.Errorf("undefined range type")
	}
	startNumber, err := strconv.Atoi(strings.TrimSpace(start))
	if err != nil {
		return fmt.Errorf("invalid number for range start: %v", start)
	}
	endNumber, err := strconv.Atoi(strings.TrimSpace(end))
	if err != nil {
		return fmt.Errorf("invalid number for range end: %v", end)
	}
	if rangeType == 1 {
		// byte range
		si.SnippetByteRangeStart = startNumber
		si.SnippetByteRangeEnd = endNumber
	} else {
		// line range
		si.SnippetLineRangeStart = startNumber
		si.SnippetLineRangeEnd = endNumber
	}
	return nil
}

func (parser *rdfParser2_2) parseRangeReference(node *gordfParser.Node, snippet *spdx.Snippet2_2) error {
	// reference is supposed to be either a resource reference to a file or a new file
	// Unfortunately, I didn't find field where this can be set in the tools-golang data model.
	// todo: set this reference to the snippet
	switch node.NodeType {
	case gordfParser.RESOURCELITERAL, gordfParser.LITERAL, gordfParser.BLANK:
		return nil
	}
	file, err := parser.getFileFromNode(node)
	if err != nil {
		return fmt.Errorf("error parsing a new file in a reference")
	}

	// a new file found within the pointer reference is an unpackaged file.
	if parser.doc.UnpackagedFiles == nil {
		parser.doc.UnpackagedFiles = map[spdx.ElementID]*spdx.File2_2{}
	}
	parser.doc.UnpackagedFiles[file.FileSPDXIdentifier] = file
	return nil
}

func setSnippetID(uri string, si *spdx.Snippet2_2) (err error) {
	fragment := getLastPartOfURI(uri)
	si.SnippetSPDXIdentifier, err = ExtractElementID(fragment)
	if err != nil {
		return fmt.Errorf("error setting snippet identifier: %v", uri)
	}
	return nil
}
