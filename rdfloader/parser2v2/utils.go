package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	urilib "github.com/RishabhBhatnagar/gordf/uri"
	"github.com/spdx/tools-golang/spdx"
	"regexp"
	"strings"
)

func getLastPartOfURI(uri string) string {
	if strings.Contains(uri, "#") {
		parts := strings.Split(uri, "#")
		return parts[len(parts)-1]
	}
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1]
}

func stripLastPartOfUri(uri string) string {
	lastPart := getLastPartOfURI(uri)
	uri = strings.TrimSuffix(uri, lastPart)
	return uri
}

func stripJoiningChars(uri string) string {
	return strings.TrimSuffix(strings.TrimSuffix(uri, "/"), "#")
}

func isUriSame(uri1, uri2 string) bool {
	return stripJoiningChars(uri1) == stripJoiningChars(uri2)
}

func (parser *rdfParser2_2) filterAllTriplesByString(subject, predicate, object string) (retTriples []*gordfParser.Triple) {
	for _, triple := range parser.gordfParserObj.Triples {
		if triple.Subject.ID == subject && triple.Predicate.ID == predicate && triple.Object.ID == object {
			retTriples = append(retTriples, triple)
		}
	}
	return retTriples
}

func (parser *rdfParser2_2) filterTriplesByRegex(triples []*gordfParser.Triple, subject, predicate, object string) (retTriples []*gordfParser.Triple, err error) {
	var subjectCompiled, objectCompiled, predicateCompiled *regexp.Regexp
	subjectCompiled, err = regexp.Compile(subject)
	if err != nil {
		return
	}
	predicateCompiled, err = regexp.Compile(predicate)
	if err != nil {
		return
	}
	objectCompiled, err = regexp.Compile(object)
	if err != nil {
		return
	}
	for _, triple := range triples {
		if subjectCompiled.MatchString(triple.Subject.ID) && predicateCompiled.MatchString(triple.Predicate.ID) && objectCompiled.MatchString(triple.Object.ID) {
			retTriples = append(retTriples, triple)
		}
	}
	return
}

func isUriValid(uri string) bool {
	_, err := urilib.NewURIRef(uri)
	return err == nil
}


// Function Below this line is taken from the tvloader/parser2v2/utils.go

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