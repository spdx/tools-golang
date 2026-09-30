package ld_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
)

// convenience for writing json as code:

type o = map[string]any // object
type l = []any          // list
type a = any            // any

var testJsonLdContext = o{"@context": o{
	"type":             "@type",
	"spdxId":           "@id",
	"spdx":             "https://spdx.org/rdf/3.0.1/terms/",
	"software_Package": "https://spdx.org/rdf/3.0.1/terms/Software/Package",
	"software_File":    "https://spdx.org/rdf/3.0.1/terms/Software/File",
	"Relationship":     "https://spdx.org/rdf/3.0.1/terms/Core/relationship",
	"software_primaryPurpose": o{
		"@context": o{
			"@vocab": "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/",
		},
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Software/primaryPurpose",
		"@type": "@vocab",
	},
	"software_additionalPurpose": o{
		"@context": o{
			"@vocab": "https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose/",
		},
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Software/additionalPurpose",
		"@type": "@vocab",
	},
	"from": o{
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Core/from",
		"@type": "@vocab",
	},
	"to": o{
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Core/to",
		"@type": "@vocab",
	},
	"specVersion": o{
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Core/specVersion",
		"@type": "http://www.w3.org/2001/XMLSchema#string",
	},
	"Element": "https://spdx.org/rdf/3.0.1/terms/Core/Element",
	"element": o{
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Core/element",
		"@type": "@vocab",
	},
	"contents": o{
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Software/File/contents",
		"@type": "http://www.w3.org/2001/XMLSchema#string",
	},
	"name": o{
		"@id":   "https://spdx.org/rdf/3.0.1/terms/Core/name",
		"@type": "http://www.w3.org/2001/XMLSchema#string",
	},
}}

type Document struct {
	_        ld.Type     `iri:"https://spdx.org/rdf/3.0.1/terms/Core/Document"`
	ID       string      `iri:"@id"`
	Elements ElementList `iri:"https://spdx.org/rdf/3.0.1/terms/Core/element"`
}

type AnyElement interface {
	asElement() *Element
}

type ElementList []AnyElement

type Element struct {
	_    ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/Element"`
	ID   string  `iri:"@id"`
	Name string  `iri:"https://spdx.org/rdf/3.0.1/terms/Core/name"`
}

func (e *Element) asElement() *Element {
	return e
}

type SoftwarePurpose struct {
	_  ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Software/SoftwarePurpose"`
	ID string  `iri:"@id"`
}

type Package struct {
	_ ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Software/Package"`
	Element
	SoftwarePurpose    SoftwarePurpose   `iri:"https://spdx.org/rdf/3.0.1/terms/Software/primaryPurpose"`
	AdditionalPurposes []SoftwarePurpose `iri:"https://spdx.org/rdf/3.0.1/terms/Software/primaryPurpose"`
}

type File struct {
	_ ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Software/File"`
	Element
	Contents string `iri:"https://spdx.org/rdf/3.0.1/terms/Software/File/contents"`
}

type Relationship struct {
	_    ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/relationship"`
	ID   string  `iri:"@id"`
	From any     `iri:"https://spdx.org/rdf/3.0.1/terms/Core/from"`
	To   []any   `iri:"https://spdx.org/rdf/3.0.1/terms/Core/to"`
}

type AnyRelationship interface {
	asRelationship() *Relationship
}

func (r *Relationship) asRelationship() *Relationship {
	return r
}

var File_DevNull = &File{Element: Element{ID: "https://example.org/iri/file/dev/null"}}

// SubRelationship implements inheritance by embedding
type SubRelationship struct {
	_ ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/sub-relationship"`
	Relationship
}

type ExternalIRI struct {
	ExternalID string `iri:"@id"`
}

func (r *ExternalIRI) asRelationship() *Relationship {
	return nil
}

func testGraph(t *testing.T, graph l) []any {
	ctx, contextURL := testContext()

	in := o{
		"@context": contextURL,
		"@graph":   graph,
	}

	buf := bytes.Buffer{}

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(in)
	require.NoError(t, err)

	graph, err = ctx.FromJSON(&buf)
	if err != nil {
		t.Fatal(err)
	}
	return graph
}

func testContext() (ld.Context, string) {
	contextURL := "https://example.org/test-context"
	return ld.NewContext().Register(contextURL, testJsonLdContext,
		Package{},
		File{},
		File_DevNull,
		Relationship{},
	), contextURL
}

func toJSON(t *testing.T, o any) string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(o)
	require.NoError(t, err)
	return buf.String()
}
