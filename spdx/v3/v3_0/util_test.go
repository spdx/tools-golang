package v3_0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_collectAllElements(t *testing.T) {
	pkg := &Package{Name: "my-pkg"}
	file := &File{Name: "main.go"}
	rel := &Relationship{
		From: pkg,
		Type: RelationshipType_Contains,
		To:   ElementList{file},
	}

	// a sub-collection nested within the SBOM, with its own element
	bundlePkg := &Package{Name: "bundled-pkg"}
	bundle := &Bundle{
		RootElements: ElementList{bundlePkg},
	}

	sbom := &SBOM{
		RootElements: ElementList{pkg, file, rel, bundle},
	}

	doc := &Document{
		SpdxDocument: SpdxDocument{
			RootElements: ElementList{sbom},
		},
	}

	got := collectAllElements(&doc.SpdxDocument)

	// everything reachable from the document graph, except the document collection itself
	require.ElementsMatch(t, []AnyElement{sbom, pkg, file, rel, bundle, bundlePkg}, got)
}

func Test_mapKeys(t *testing.T) {
	p1 := &Package{Name: "p1"}
	p2 := &Package{Name: "p2"}

	tests := []struct {
		name string
		in   map[any]struct{}
		want []any
	}{
		{
			name: "empty map",
			in:   map[any]struct{}{},
			want: []any{},
		},
		{
			name: "string keys",
			in:   map[any]struct{}{"a": {}, "b": {}, "c": {}},
			want: []any{"a", "b", "c"},
		},
		{
			name: "zero-value key is dropped",
			in:   map[any]struct{}{"a": {}, "": {}, "b": {}},
			want: []any{"a", "b"},
		},
		{
			name: "nil pointer key is dropped",
			in:   map[any]struct{}{p1: {}, (*Package)(nil): {}, p2: {}},
			want: []any{p1, p2},
		},
		{
			name: "all keys dropped",
			in:   map[any]struct{}{"": {}, (*Package)(nil): {}},
			want: []any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapKeys(tt.in)
			// map iteration order is not deterministic
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func Test_notNil(t *testing.T) {
	p1 := &Package{Name: "p1"}
	p2 := &Package{Name: "p2"}

	t.Run("nil slice", func(t *testing.T) {
		require.Nil(t, notNil[*Package, []*Package](nil))
	})

	t.Run("no nils returns input unchanged", func(t *testing.T) {
		in := []*Package{p1, p2}
		require.Equal(t, in, notNil(in))
	})

	t.Run("drops nil pointers, preserves order", func(t *testing.T) {
		in := []*Package{p1, nil, p2}
		require.Equal(t, []*Package{p1, p2}, notNil(in))
	})

	t.Run("drops leading nil", func(t *testing.T) {
		in := []*Package{nil, p1, p2}
		require.Equal(t, []*Package{p1, p2}, notNil(in))
	})

	t.Run("zero values treated as nil", func(t *testing.T) {
		in := []string{"a", "", "b"}
		require.Equal(t, []string{"a", "b"}, notNil(in))
	})

	t.Run("all nil yields empty", func(t *testing.T) {
		in := []*Package{nil, nil}
		require.Empty(t, notNil(in))
	})
}
