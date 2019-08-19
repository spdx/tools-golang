package rdf2v1

import (
	"fmt"
	"strings"

	"github.com/deltamobile/goraptor"
)

const (
	BaseUri    = "http://spdx.org/rdf/terms#"
	LicenseUri = "http://spdx.org/licenses/"
)

var rdfPrefixes = map[string]string{
	"ns:":   "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"rdf:":  "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"doap:": "http://usefulinc.com/ns/doap#",
	"rdfs:": "http://www.w3.org/2000/01/rdf-schema#",
	"j.0:":  "http://www.w3.org/2009/pointers#",
	"":      BaseUri,
}

// simple key value pair struct
type Pair struct {
	key, val string
}

// Converts typeX to its full URI accorinding to rdfPrefixes,
// if no : is found in the string it'll assume it as "spdx:" and expand to baseUri
func Prefix(k string) *goraptor.Uri {
	var pref string = BaseUri
	rest := k
	if i := strings.Index(k, ":"); i >= 0 {
		pref = k[:i+1]
		rest = k[i+1:]
	}
	if long, ok := rdfPrefixes[pref]; ok {
		pref = long
	}
	uri := goraptor.Uri(pref + rest)
	return &uri
}

// Change the RDF prefixes to their short forms.
func ShortPrefix(t goraptor.Term) string {
	str := termStr(t)
	for short, long := range rdfPrefixes {
		if strings.HasPrefix(str, long) {
			str = strings.Replace(str, long, short, 1)
			return strings.Replace(str, long, short, 1)
		}
	}
	return str
}

// Converts goraptor.Term (Subject, Predicate and Object) to string.
func termStr(term goraptor.Term) string {
	switch t := term.(type) {
	case *goraptor.Uri:
		return string(*t)
	case *goraptor.Blank:
		return string(*t)
	case *goraptor.Literal:
		return t.Value
	default:
		return ""
	}
}

// Uri, Literal and Blank are goraptors named types
// Return *goraptor.Uri
func Uri(uri string) *goraptor.Uri {
	return (*goraptor.Uri)(&uri)
}

// Return *goraptor.Literal
func Literal(lit string) *goraptor.Literal {
	return &goraptor.Literal{Value: lit}
}

// Return *goraptor.Blank from string
func Blank(b string) *goraptor.Blank {
	return (*goraptor.Blank)(&b)
}

func extractSubs(value string) (string, string, error) {
	// parse the value to see if it's a valid subvalue format
	sp := strings.SplitN(value, ":", 2)
	if len(sp) == 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (no colon found)", value)
	}

	subkey := strings.TrimSpace(sp[0])
	subvalue := strings.TrimSpace(sp[1])

	return subkey, subvalue, nil
}

func ExtractValueType(value string, t string) string {
	subkey, subvalue, _ := extractSubs(value)
	if subkey == t {
		return subvalue
	}
	return ""
}
func ExtractKey(value string) string {
	subkey, _, _ := extractSubs(value)
	return subkey
}
func ExtractKeyValue(value string, sub string) string {
	// parse the value to see if it's a valid subvalue format
	subkey, subvalue, _ := extractSubs(value)
	if sub == "subkey" {
		return subkey
	} else if sub == "subvalue" {
		return subvalue
	} else {
		return ""
	}
}

func ExtractCodeAndExcludes(value string) (string, string) {

	sp := strings.SplitN(value, "(excludes:", 2)
	if len(sp) < 2 {
		return value, ""
	}
	code := strings.TrimSpace(sp[0])
	parsedSp := strings.SplitN(sp[1], ")", 2)
	fileName := strings.TrimSpace(parsedSp[0])
	return code, fileName
}

func ExtractNs(value string) (string, string, error) {
	s := strings.SplitN(value, "#", 2)
	if len(s) == 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (no # found)", value)
	}
	return s[0], s[1], nil
}

func ExtractId(value string) string {
	s := strings.SplitN(value, "#", 2)
	if len(s) == 1 {
		return ""
	}
	return s[1]
}

func ExtractRelType(value string) string {
	s := strings.SplitN(value, "_", 2)
	if len(s) == 1 {
		return ""
	}
	return s[1]
}
func InsertId(value string) ValueStr {
	s := value + "#SPDXRef-DOCUMENT"
	vs := Str(s)
	return vs
}
