package ld_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_decode(t *testing.T) {
	tests := []struct {
		name  string
		graph l
		want  func() a
	}{
		{
			name: "simple graph",
			graph: l{
				o{"@type": "software_Package", "@id": "pkg-1", "name": "pkg 1"},
				o{"@type": "software_File", "@id": "file-1", "contents": "file 1"},
				o{"@type": "Relationship", "from": "file-1", "to": l{"pkg-1"}},
			},
			want: func() a {
				p := &Package{Element: Element{ID: "pkg-1", Name: "pkg 1"}}
				f := &File{Element: Element{ID: "file-1"}, Contents: "file 1"}
				r := &Relationship{From: f, To: l{p}}
				return l{p, f, r}
			},
		},
		{
			name: "top level named individual",
			graph: l{
				o{
					"@id": "https://example.org/iri/file/dev/null",
				},
			},
			want: func() a {
				return l{File_DevNull}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := testGraph(t, tt.graph)

			want := tt.want()

			if diff := cmp.Diff(want, graph); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}
