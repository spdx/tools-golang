// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx"
	"strings"
)

// parsing the relationship that exists in the rdf document.
// Relationship is of type RefA relationType RefB.
// parsing the relationship appends the relationship to the current document's
// Relationships Slice.
func (parser *rdfParser2_2) parseRelationship(triple *gordfParser.Triple) (err error) {
	reln := spdx.Relationship2_2{}

	reln.RefA, err = getReferenceFromURI(triple.Subject.ID)
	if err != nil {
		return err
	}

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

func (parser *rdfParser2_2) parseRelatedElementFromTriple(reln *spdx.Relationship2_2, triple *gordfParser.Triple) error {
	// iterate over relatedElement Type and check which SpdxElement it is.
	var err error
	switch triple.Object.ID {
	case SPDX_FILE:
		file, err := parser.getFileFromNode(triple.Subject)
		if err != nil {
			return fmt.Errorf("error setting a file: %v", err)
		}
		reln.RefB = spdx.DocElementID{
			DocumentRefID: "",
			ElementRefID:  file.FileSPDXIdentifier,
		}

	case SPDX_PACKAGE:
		pkg, err := parser.getPackageFromNode(triple.Subject)
		if err != nil {
			return fmt.Errorf("error setting a package inside a relationship: %v", err)
		}
		reln.RefB = spdx.DocElementID{
			DocumentRefID: "",
			ElementRefID:  pkg.PackageSPDXIdentifier,
		}
		parser.doc.Packages[pkg.PackageSPDXIdentifier] = pkg

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
func getReferenceFromURI(uri string) (spdx.DocElementID, error) {
	fragment := getLastPartOfURI(uri)
	switch strings.ToLower(strings.TrimSpace(fragment)) {
	case "noassertion", "none":
		return spdx.DocElementID{
			DocumentRefID: "",
			ElementRefID:  spdx.ElementID(strings.ToUpper(fragment)),
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
