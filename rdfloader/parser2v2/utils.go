// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"errors"
	"fmt"
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
	urilib "github.com/spdx/gordf/uri"
	"github.com/spdx/tools-golang/spdx"
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
	return parts[len(parts)-1]
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

// returns which boolean was given as an input
// string(bool) is the only possible input for which it will not raise any error.
func boolFromString(boolString string) (bool, error) {
	switch strings.ToLower(boolString) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("boolean string can be either true/false")
	}
}

/* Function Below this line is taken from the tvloader/parser2v2/utils.go */

// used to extract DocumentRef and SPDXRef values from an SPDX Identifier
// which can point either to this document or to a different one
func ExtractDocElementID(value string) (spdx.DocElementID, error) {
	docRefID := ""
	idStr := value

	// check prefix to see if it's a DocumentRef ID
	if strings.HasPrefix(idStr, "DocumentRef-") {
		// extract the part that comes between "DocumentRef-" and ":"
		strs := strings.Split(idStr, ":")
		// should be exactly two, part before and part after
		if len(strs) < 2 {
			return spdx.DocElementID{}, fmt.Errorf("no colon found although DocumentRef- prefix present")
		}
		if len(strs) > 2 {
			return spdx.DocElementID{}, fmt.Errorf("more than one colon found")
		}

		// trim the prefix and confirm non-empty
		docRefID = strings.TrimPrefix(strs[0], "DocumentRef-")
		if docRefID == "" {
			return spdx.DocElementID{}, fmt.Errorf("document identifier has nothing after prefix")
		}
		// and use remainder for element ID parsing
		idStr = strs[1]
	}

	// check prefix to confirm it's got the right prefix for element IDs
	if !strings.HasPrefix(idStr, "SPDXRef-") {
		return spdx.DocElementID{}, fmt.Errorf("missing SPDXRef- prefix for element identifier")
	}

	// make sure no colons are present
	if strings.Contains(idStr, ":") {
		// we know this means there was no DocumentRef- prefix, because
		// we would have handled multiple colons above if it was
		return spdx.DocElementID{}, fmt.Errorf("invalid colon in element identifier")
	}

	// trim the prefix and confirm non-empty
	eltRefID := strings.TrimPrefix(idStr, "SPDXRef-")
	if eltRefID == "" {
		return spdx.DocElementID{}, fmt.Errorf("element identifier has nothing after prefix")
	}

	// we're good
	return spdx.DocElementID{DocumentRefID: docRefID, ElementRefID: spdx.ElementID(eltRefID)}, nil
}

// used to extract SPDXRef values only from an SPDX Identifier which can point
// to this document only. Use extractDocElementID for parsing IDs that can
// refer either to this document or a different one.
func ExtractElementID(value string) (spdx.ElementID, error) {
	// check prefix to confirm it's got the right prefix for element IDs
	if !strings.HasPrefix(value, "SPDXRef-") {
		return spdx.ElementID(""), fmt.Errorf("missing SPDXRef- prefix for element identifier")
	}

	// make sure no colons are present
	if strings.Contains(value, ":") {
		return spdx.ElementID(""), fmt.Errorf("invalid colon in element identifier")
	}

	// trim the prefix and confirm non-empty
	eltRefID := strings.TrimPrefix(value, "SPDXRef-")
	if eltRefID == "" {
		return spdx.ElementID(""), fmt.Errorf("element identifier has nothing after prefix")
	}

	// we're good
	return spdx.ElementID(eltRefID), nil
}

// used to extract key / value from embedded substrings
// returns subkey, subvalue, nil if no error, or "", "", error otherwise
func ExtractSubs(value string, sep string) (string, string, error) {
	// parse the value to see if it's a valid subvalue format
	sp := strings.SplitN(value, sep, 2)
	if len(sp) == 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (no %s found)", value, sep)
	}

	subkey := strings.TrimSpace(sp[0])
	subvalue := strings.TrimSpace(sp[1])

	return subkey, subvalue, nil
}
