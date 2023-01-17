// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"
	"strings"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// parsing the relationship that exists in the rdf document.
// Relationship is of type RefA relationType RefB.
// parsing the relationship appends the relationship to the current document's
// Relationships Slice.
func (parser *rdfParser2_3) parseRelationship(triple *gordfParser.Triple) (err error) {
	reln := spdx.Relationship{}

	reln.RefA, err = getReferenceFromURI(triple.Subject.ID)
	if err != nil {
		return err
	}

	currState := parser.cache[triple.Object.ID]
	if currState == nil {
		// there is no entry about the state of current package node.
		// this is the first time we're seeing this node.
		parser.cache[triple.Object.ID] = &nodeState{
			object: reln,
			Color:  WHITE,
		}
	} else if currState.Color == GREY {
		// we have already started parsing this relationship node and we needn't parse it again.
		return nil
	}

	// setting color of the state to grey to indicate that we've started to
	// parse this node once.
	parser.cache[triple.Object.ID].Color = GREY

	// setting state color to black to indicate when we're done parsing this node.
	defer func() { parser.cache[triple.Object.ID].Color = BLACK }()

	for _, subTriple := range parser.nodeToTriples(triple.Object) {
		switch subTriple.Predicate.ID {
		case SPDX_RELATIONSHIP_TYPE:
			// cardinality: exactly 1
			reln.Relationship, err = getRelationshipTypeFromURI(subTriple.Object.ID)
		case RDF_TYPE:
			// cardinality: exactly 1
			continue
		case SPDX_RELATED_SPDX_ELEMENT:
			// cardinality: exactly 1
			// assumes: spdx-element is a uri
			reln.RefB, err = getReferenceFromURI(subTriple.Object.ID)
			if err != nil {
				return err
			}

			relatedSpdxElementTriples := parser.nodeToTriples(subTriple.Object)
			if len(relatedSpdxElementTriples) == 0 {
				continue
			}

			typeTriples := rdfwriter.FilterTriples(relatedSpdxElementTriples, &subTriple.Object.ID, &RDF_TYPE, nil)
			if len(typeTriples) != 1 {
				return fmt.Errorf("expected %s to have exactly one rdf:type triple. found %d triples", subTriple.Object, len(typeTriples))
			}
			err = parser.parseRelatedElementFromTriple(&reln, typeTriples[0])
			if err != nil {
				return err
			}
		case RDFS_COMMENT:
			// cardinality: max 1
			reln.RelationshipComment = subTriple.Object.ID
		default:
			return fmt.Errorf("unexpected predicate id: %s", subTriple.Predicate.ID)
		}
		if err != nil {
			return err
		}
	}
	parser.doc.Relationships = append(parser.doc.Relationships, &reln)
	return nil
}

func (parser *rdfParser2_3) parseRelatedElementFromTriple(reln *spdx.Relationship, triple *gordfParser.Triple) error {
	// iterate over relatedElement Type and check which SpdxElement it is.
	var err error
	switch triple.Object.ID {
	case SPDX_FILE:
		file, err := parser.getFileFromNode(triple.Subject)
		if err != nil {
			return fmt.Errorf("error setting a file: %v", err)
		}
		reln.RefB = common.DocElementID{
			DocumentRefID: "",
			ElementRefID:  file.FileSPDXIdentifier,
		}

	case SPDX_PACKAGE:
		pkg, err := parser.getPackageFromNode(triple.Subject)
		if err != nil {
			return fmt.Errorf("error setting a package inside a relationship: %v", err)
		}
		reln.RefB = common.DocElementID{
			DocumentRefID: "",
			ElementRefID:  pkg.PackageSPDXIdentifier,
		}

	case SPDX_SPDX_ELEMENT:
		// it shouldn't be associated with any other triple.
		// it must be a uri reference.
		reln.RefB, err = ExtractDocElementID(getLastPartOfURI(triple.Subject.ID))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("undefined relatedElement %s found while parsing relationship", triple.Object.ID)
	}
	return nil
}

// references like RefA and RefB of any relationship
func getReferenceFromURI(uri string) (common.DocElementID, error) {
	fragment := getLastPartOfURI(uri)
	switch strings.ToLower(strings.TrimSpace(fragment)) {
	case "noassertion", "none":
		return common.DocElementID{
			DocumentRefID: "",
			ElementRefID:  common.ElementID(strings.ToUpper(fragment)),
		}, nil
	}
	return ExtractDocElementID(fragment)
}

// note: relationshipType is case sensitive.
func getRelationshipTypeFromURI(relnTypeURI string) (string, error) {
	relnTypeURI = strings.TrimSpace(relnTypeURI)
	lastPart := getLastPartOfURI(relnTypeURI)
	if !strings.HasPrefix(lastPart, PREFIX_RELATIONSHIP_TYPE) {
		return "", fmt.Errorf("relationshipType must start with %s. found %s", PREFIX_RELATIONSHIP_TYPE, lastPart)
	}
	lastPart = strings.TrimPrefix(lastPart, PREFIX_RELATIONSHIP_TYPE)

	lastPart = strings.TrimSpace(lastPart)
	for _, validRelationshipType := range AllRelationshipTypes() {
		if lastPart == validRelationshipType {
			return lastPart, nil
		}
	}
	return "", fmt.Errorf("unknown relationshipType: '%s'", lastPart)
}
