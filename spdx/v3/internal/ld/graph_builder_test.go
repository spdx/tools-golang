package ld

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func Test_graphBuilder(t *testing.T) {
	type typ struct {
		_     Type      `iri:"http://t-iri"`
		Id    string    `iri:"@id"`
		Str   string    `iri:"http://str-iri"`
		Bool  bool      `iri:"http://bool-iri"`
		Int   int       `iri:"http://int-iri"`
		Float float64   `iri:"http://float-iri"`
		Time  time.Time `iri:"http://time-iri"`
	}

	type typ2 struct {
		_          Type   `iri:"http://t2-iri"`
		Identifier string `iri:"@id"`
		T1         *typ   `iri:"http://t2-to-t1-iri"`
	}

	typContext := o{
		"t":  "http://t-iri",
		"t2": "http://t2-iri",
		"s": o{
			"@id":   "http://str-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#string",
		},
		"b": o{
			"@id":   "http://bool-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#boolean",
		},
		"i": o{
			"@id":   "http://int-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#integer",
		},
		"f": o{
			"@id":   "http://float-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#decimal",
		},
		"tm": o{
			"@id":   "http://time-iri",
			"@type": "http://www.w3.org/2001/XMLSchema#dateTimeStamp",
		},
		"t2t1": "http://t2-to-t1-iri",
	}

	contextWithIdTypeOverride := merge(typContext, o{
		"aliased-type": "@type",
		"aliased-id":   "@id",
	})

	contextURI := "http://example.org/uri"

	tests := []struct {
		name     string
		context  o
		graph    func() any
		expected any
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name:    "basic no context full iri",
			context: o{},
			graph: func() any {
				return &typ{
					Str:   "str-val",
					Bool:  true,
					Int:   101,
					Float: 940.33,
					Time:  mar25noon,
				}
			},
			expected: o{
				"@context":         contextURI,
				"@type":            "http://t-iri",
				"@id":              "_:typ-1",
				"http://str-iri":   expanded("str-val"),
				"http://bool-iri":  expanded(true),
				"http://int-iri":   expanded(101),
				"http://float-iri": expanded(940.33),
				"http://time-iri":  expanded(mar25noon),
			},
		},
		{
			name:    "basic all aliases context",
			context: typContext,
			graph: func() any {
				return &typ{
					Str:   "str-val",
					Bool:  true,
					Int:   101,
					Float: 940.33,
					Time:  mar25noon,
				}
			},
			expected: o{
				"@context": contextURI,
				"@type":    "t",
				"@id":      "_:typ-1",
				"s":        "str-val",
				"b":        true,
				"i":        101,
				"f":        940.33,
				"tm":       mar25noon.Format(time.RFC3339),
			},
		},
		{
			name:    "all aliases overridden id type",
			context: contextWithIdTypeOverride,
			graph: func() any {
				return &typ{}
			},
			expected: o{
				"@context":     contextURI,
				"aliased-id":   "_:typ-1",
				"aliased-type": "t",
			},
		},
		{
			name:    "multiple refs gets id",
			context: contextWithIdTypeOverride,
			graph: func() any {
				t1 := &typ{
					Str: "a-val",
				}
				return l{
					&typ2{
						T1: t1,
					},
					&typ2{
						T1: t1,
					},
				}
			},
			expected: o{
				"@context": contextURI,
				"@graph": l{
					o{
						"aliased-type": "t",
						"aliased-id":   "_:typ-1",
						"s":            "a-val",
					},
					o{
						"aliased-type": "t2",
						"aliased-id":   "_:typ2-1",
						"t2t1":         "_:typ-1",
					},
					o{
						"aliased-type": "t2",
						"aliased-id":   "_:typ2-2",
						"t2t1":         "_:typ-1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			ctx.Register(contextURI, o{"@context": tt.context}, typ{}, typ2{})
			expected := tt.expected

			var graphs l
			graph := tt.graph()
			if g, ok := graph.(l); ok {
				graphs = g
			} else {
				graphs = append(graphs, graph)
			}
			var got any
			got, err := ctx.(*context).toMaps(graphs...)

			wantErr := require.NoError
			if tt.wantErr != nil {
				wantErr = tt.wantErr
			}
			wantErr(t, err)

			d := cmp.Diff(expected, got)
			if d != "" {
				t.Fatal(d)
			}
		})
	}
}

func expanded[T any](value T) any {
	var v any = value
	t := ""
	var out any
	switch v := v.(type) {
	case string:
		t = "http://www.w3.org/2001/XMLSchema#string"
		out = v
	case bool:
		t = "http://www.w3.org/2001/XMLSchema#boolean"
		out = v
	case float32, float64:
		t = "http://www.w3.org/2001/XMLSchema#decimal"
		out = v
	case byte, int, int8, int16, int32, int64:
		t = "http://www.w3.org/2001/XMLSchema#integer"
		out = v
	case time.Time:
		t = "http://www.w3.org/2001/XMLSchema#dateTimeStamp"
		out = v.Format(time.RFC3339)
	default:
		panic("unsupported type")
	}

	return o{
		JsonTypeProp:  t,
		JsonValueProp: out,
	}
}
