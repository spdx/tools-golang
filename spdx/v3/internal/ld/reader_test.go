package ld

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func Test_nestedDeserialization(t *testing.T) {
	type t1 struct {
		_     Type   `iri:"https://example.org/test-1"`
		Id    string `iri:"@id"`
		Child any    `iri:"https://example.org/test-1/child"`
	}

	contextURI := "http://example.org/uri"

	ctx := NewContext().Register(contextURI, o{"@context": o{
		"t": "https://example.org/test-1",
		"child": o{
			"@id":   "https://example.org/test-1/child",
			"@type": "@vocab",
		},
	}}, t1{})

	graph, err := ctx.(*context).fromMaps(o{
		"@context": contextURI,
		"@graph": l{
			o{
				"@id":   "_:t1",
				"@type": "t",
				"child": "_:t2",
			},
			o{
				"@id":   "_:t2",
				"@type": "t",
				"child": "_:t3",
			},
			o{
				"@id":   "_:t3",
				"@type": "t",
				"child": "_:t1",
			},
		},
	})
	got := graph[0].(*t1)
	require.NoError(t, err)
	require.True(t, got.Child.(*t1).Child != nil)
}

func Test_readerAliasFields(t *testing.T) {
	type typ struct {
		_     Type      `iri:"https://example.org/test-iri"`
		Id    string    `iri:"@id"`
		Str   string    `iri:"https://example.org/test-iri/str-iri"`
		Bool  bool      `iri:"https://example.org/test-iri/bool-iri"`
		Int   int       `iri:"https://example.org/test-iri/int-iri"`
		Float float64   `iri:"https://example.org/test-iri/float-iri"`
		Time  time.Time `iri:"https://example.org/test-iri/time-iri"`
	}

	typContext := o{
		"t": "https://example.org/test-iri",
		"s": o{
			"@id":   "https://example.org/test-iri/str-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#string",
		},
		"b": o{
			"@id":   "https://example.org/test-iri/bool-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#boolean",
		},
		"i": o{
			"@id":   "https://example.org/test-iri/int-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#integer",
		},
		"f": o{
			"@id":   "https://example.org/test-iri/float-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#decimal",
		},
		"tm": o{
			"@id":   "https://example.org/test-iri/time-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#dateTimeStamp",
		},
	}

	idTypeAliasedContext := merge(typContext, o{
		"myId":   "@id",
		"myType": "@type",
	})

	tests := []struct {
		name     string
		context  o
		graph    any
		expected func() any
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name:    "all aliases, @context prop",
			context: typContext,
			graph: o{
				"@type": "t",
				"s":     "joe",
				"b":     true,
				"i":     12,
				"f":     39.11,
				"tm":    mar25noon.Format(time.RFC3339),
			},
			expected: func() any {
				return &typ{
					Str:   "joe",
					Bool:  true,
					Int:   12,
					Float: 39.11,
					Time:  mar25noon,
				}
			},
		},
		{
			name:    "all aliases, no @context prop",
			context: typContext,
			graph: o{
				"@type": "t",
				"s":     "joe",
				"b":     true,
				"i":     12,
				"f":     39.11,
				"tm":    mar25noon.Format(time.RFC3339),
			},
			expected: func() any {
				return &typ{
					Str:   "joe",
					Bool:  true,
					Int:   12,
					Float: 39.11,
					Time:  mar25noon,
				}
			},
		},
		{
			name:    "full IRI, no aliases",
			context: o{},
			graph: o{
				"@type": "https://example.org/test-iri",
				"https://example.org/test-iri/str-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#string",
					"@value": "joe",
				},
				"https://example.org/test-iri/bool-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#boolean",
					"@value": true,
				},
				"https://example.org/test-iri/int-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#integer",
					"@value": 12,
				},
				"https://example.org/test-iri/float-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#decimal",
					"@value": 39.11,
				},
				"https://example.org/test-iri/time-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#dateTimeStamp",
					"@value": mar25noon.Format(time.RFC3339),
				},
			},
			expected: func() any {
				return &typ{
					Str:   "joe",
					Bool:  true,
					Int:   12,
					Float: 39.11,
					Time:  mar25noon,
				}
			},
		},
		{
			name:    "full IRI, all aliases",
			context: typContext,
			graph: o{
				"@type": "https://example.org/test-iri",
				"https://example.org/test-iri/str-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#string",
					"@value": "joe",
				},
				"https://example.org/test-iri/bool-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#boolean",
					"@value": true,
				},
				"https://example.org/test-iri/int-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#integer",
					"@value": 12,
				},
				"https://example.org/test-iri/float-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#decimal",
					"@value": 39.11,
				},
				"https://example.org/test-iri/time-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#dateTimeStamp",
					"@value": mar25noon.Format(time.RFC3339),
				},
			},
			expected: func() any {
				return &typ{
					Str:   "joe",
					Bool:  true,
					Int:   12,
					Float: 39.11,
					Time:  mar25noon,
				}
			},
		},
		{
			name:    "mixed IRI and aliases",
			context: typContext,
			graph: o{
				"@type": "t",
				"https://example.org/test-iri/str-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#string",
					"@value": "joe",
				},
				"b": true,
				"https://example.org/test-iri/int-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#integer",
					"@value": 12,
				},
				"f": 39.11,
				"https://example.org/test-iri/time-iri": o{
					"@type":  "http://www.w3.org/2001/XMLSchema#dateTimeStamp",
					"@value": mar25noon.Format(time.RFC3339),
				},
			},
			expected: func() any {
				return &typ{
					Str:   "joe",
					Bool:  true,
					Int:   12,
					Float: 39.11,
					Time:  mar25noon,
				}
			},
		},
		{
			name:    "id and type aliases",
			context: idTypeAliasedContext,
			graph: o{
				"myType": "t",
				"s":      "joe",
				"b":      true,
				"i":      12,
				"f":      39.11,
				"tm":     mar25noon.Format(time.RFC3339),
			},
			expected: func() any {
				return &typ{
					Str:   "joe",
					Bool:  true,
					Int:   12,
					Float: 39.11,
					Time:  mar25noon,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.expected()
			contextURI := "http://example.org/uri"

			ctx := NewContext().Register(contextURI, o{"@context": tt.context},
				// register an empty instance of the returned type:
				reflect.New(reflect.TypeOf(expected).Elem()).Interface())
			graph := tt.graph
			if _, ok := tt.graph.(l); !ok {
				graph = l{graph}
			}
			var got any
			gotList, err := ctx.(*context).fromMaps(o{
				"@context": contextURI,
				"@graph":   graph,
			})

			wantErr := require.NoError
			if tt.wantErr != nil {
				wantErr = tt.wantErr
			}
			wantErr(t, err)

			got = gotList
			if _, ok := expected.(l); !ok {
				if len(gotList) > 0 {
					got = gotList[0]
				}
			}

			d := cmp.Diff(expected, got)
			if d != "" {
				t.Fatal(d)
			}
		})
	}
}

var mar25noon = get(time.Parse(time.RFC3339, "2025-03-25T12:00:00Z"))

func get[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

type o = map[string]any
type l = []any
