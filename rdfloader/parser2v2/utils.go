// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"errors"
	"fmt"
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
	urilib "github.com/spdx/gordf/uri"
	"strings"
)

// a uri is of type baseURI#fragment or baseFragment/subFragment
// returns fragment or subFragment when given as an input.
func getLastPartOfURI(uri string) string {
	if strings.Contains(uri, "#") {
		parts := strings.Split(uri, "#")
		return parts[len(parts)-1]
	}
	parts := strings.Split(uri, "/")
	return strings.TrimSpace(parts[len(parts)-1])
}

func isUriValid(uri string) bool {
	_, err := urilib.NewURIRef(uri)
	return err == nil
}

func getNodeTypeFromTriples(triples []*gordfParser.Triple, node *gordfParser.Node) (string, error) {
	if node == nil {
		return "", errors.New("empty node passed to find node type")
	}
	typeTriples := rdfwriter.FilterTriples(triples, &node.ID, &RDF_TYPE, nil)
	switch len(typeTriples) {
	case 0:
		return "", fmt.Errorf("node{%v} not associated with any type triple", node)
	case 1:
		return typeTriples[0].Object.ID, nil
	default:
		return "", fmt.Errorf("node{%v} is associated with more than one type triples", node)
	}
}

func (parser *rdfParser2_2) nodeToTriples(node *gordfParser.Node) []*gordfParser.Triple {
	if node == nil {
		return []*gordfParser.Triple{}
	}
	return parser.nodeStringToTriples[node.String()]
}
