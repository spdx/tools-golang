package ld

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func Test_registeredInstances(t *testing.T) {
	type AType struct {
		_      Type   `iri:"https://example.org/context/thing-a"`
		ID     string `iri:"@id"`
		Name   string `iri:"https://example.org/context/name"`
		Values []any  `iri:"https://example.org/context/values"`
	}

	inst := &AType{
		ID:   "https://example.org/context/an-instance",
		Name: "an-instance",
	}

	contextName := "https://example.org/context"
	ctx := NewContext().Register(contextName, o{"@context": o{
		"thing-a": "https://example.org/context/thing-a",
		"name": o{
			"@type": "http://www.w3.org/2001/XMLSchema#string",
			"@id":   "https://example.org/context/name",
		},
		"values": o{
			"@type": "@vocab",
			"@id":   "https://example.org/context/values",
		},
	}}, inst)

	g, err := ctx.(*context).fromMaps(o{
		"@context": contextName,
		//"@graph": l{
		//	o{
		"@type": "thing-a",
		"name":  "some-name",
		"values": l{
			"https://example.org/context/an-instance",
		},
		//	},
		//},
	})
	require.NoError(t, err)

	got := g[0].(*AType)
	require.Equal(t, "some-name", got.Name)
	require.Len(t, got.Values, 1)
	require.Equal(t, inst, got.Values[0])
}

func Test_MultiRegistration(t *testing.T) {
	type AType struct {
		_    Type   `iri:"https://example.org/context/thing-a"`
		AnID string `iri:"@id"`
		Name string `iri:"https://example.org/context/name"`
	}

	contextName := "https://example.org/context"
	ctx := NewContext().Register(contextName, o{"@context": o{
		"name": o{
			"@type": "http://www.w3.org/2001/XMLSchema#string",
			"@id":   "https://example.org/context/name",
		},
	}}, AType{})
	require.Len(t, ctx.(*context).contextMap, 1)

	type BType struct {
		_         Type   `iri:"https://example.org/context/thing-b"`
		AnotherID string `iri:"@id"`
		Name      string `iri:"https://example.org/context/name"`
	}

	ctx = ctx.Register(contextName, o{"@context": o{}}, BType{})
	require.Len(t, ctx.(*context).contextMap, 1)

	maps, err := ctx.(*context).toMaps(&AType{
		AnID: "id1",
		Name: "A",
	}, &BType{
		AnotherID: "id2",
		Name:      "B",
	})
	require.NoError(t, err)

	diff := cmp.Diff(o{
		"@context": "https://example.org/context",
		"@graph": l{
			o{"@type": "https://example.org/context/thing-a", "name": "A", "@id": "id1"},
			o{"@type": "https://example.org/context/thing-b", "name": "B", "@id": "id2"},
		},
	}, maps)
	if diff != "" {
		t.Fatal(diff)
	}
}
