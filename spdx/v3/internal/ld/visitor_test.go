package ld

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type value struct {
	Name string
}

func Test_VisitObjectGraph(t *testing.T) {
	t.Run("visits the root pointer", func(t *testing.T) {
		o := &value{}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("visits pointer fields", func(t *testing.T) {
		type outer struct {
			Inner *value
		}
		inner := &value{Name: "inner"}
		o := &outer{Inner: inner}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("visits pointers to pointers", func(t *testing.T) {
		type outer struct {
			PP **value
		}
		inner := &value{Name: "inner"}
		pp := &inner
		o := &outer{PP: pp}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("addrs structs in fields", func(t *testing.T) {
		type outer struct {
			Value value
		}
		o := &outer{Value: value{Name: "value"}}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("addrs structs in slices", func(t *testing.T) {
		s := []value{{Name: "a"}, {Name: "b"}}
		err := VisitObjectGraph(s, func(path []any, v any) error {
			if i, ok := v.(*value); ok {
				i.Name += "-visited"
			}
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, "a-visited", s[0].Name)
		require.Equal(t, "b-visited", s[1].Name)
	})

	t.Run("visits pointers in maps", func(t *testing.T) {
		inner := &value{Name: "inner"}
		m := map[string]*value{"k": inner}
		require.Len(t, collectVisits[*value](t, m), 1)
	})

	t.Run("visits pointers in interfaces", func(t *testing.T) {
		type outer struct {
			Iface any
		}
		inner := &value{Name: "inner"}
		o := &outer{Iface: inner}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("visits structs in interfaces", func(t *testing.T) {
		type outer struct {
			Iface any
		}
		o := &outer{Iface: value{Name: "inner"}}
		require.Len(t, collectVisits[value](t, o), 1)
	})

	t.Run("visits pointers within anonymous structs in interfaces", func(t *testing.T) {
		type outer struct {
			Iface any
		}
		type base struct {
			Inner *value
			Name  string
		}
		type derived struct {
			base
		}
		inner := &value{Name: "inner"}
		o := &outer{Iface: derived{base: base{Inner: inner}}}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("visits shared pointers once", func(t *testing.T) {
		enableDebug(t)
		type outer struct {
			A, B *value
		}
		shared := &value{Name: "shared"}
		o := &outer{A: shared, B: shared}
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("handles self-referential pointers", func(t *testing.T) {
		type node struct {
			Name string
			Next *node
		}
		n := &node{Name: "self"}
		n.Next = n
		require.Len(t, collectVisits[*node](t, n), 1)
	})

	t.Run("handles pointer cycles", func(t *testing.T) {
		type node struct {
			Name string
			Next *node
		}
		a := &node{Name: "a"}
		b := &node{Name: "b"}
		a.Next = b
		b.Next = a
		require.Len(t, collectVisits[*node](t, a), 2)
	})

	t.Run("handles self-referential slices", func(t *testing.T) {
		s := []any{nil}
		s[0] = s
		err := VisitObjectGraph(s, func(path []any, value any) error { return nil })
		require.NoError(t, err)
	})

	t.Run("handles self-referential maps", func(t *testing.T) {
		m := map[string]any{}
		m["self"] = m
		err := VisitObjectGraph(m, func(path []any, value any) error { return nil })
		require.NoError(t, err)
	})

	t.Run("handles cycles through unexported fields", func(t *testing.T) {
		type unexported struct {
			next *unexported
		}
		u := &unexported{}
		u.next = u
		err := VisitObjectGraph(u, func(path []any, value any) error { return nil })
		require.NoError(t, err)
	})

	t.Run("does not visit nil pointers", func(t *testing.T) {
		type outer struct {
			PtrSlice []*value
			PtrMap   map[string]*value
		}
		o := &outer{
			PtrSlice: []*value{nil},
			PtrMap:   map[string]*value{"k": nil},
		}
		visits := collectVisits[any](t, o)
		for _, v := range visits {
			rv := reflect.ValueOf(v.value)
			if rv.Kind() == reflect.Pointer {
				require.False(t, rv.IsNil(), "visited a nil pointer at: %s", v.path)
			}
		}
	})

	t.Run("does not visit nil pointers boxed in interfaces", func(t *testing.T) {
		type outer struct {
			Iface any
		}
		o := &outer{Iface: (*value)(nil)} // non-nil interface boxing a typed nil pointer
		err := VisitObjectGraph(o, func(path []any, v *value) error {
			require.NotNil(t, v, "visitor called with a nil pointer at: %s", StringifyPath(path))
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("does not visit pointers used as map keys", func(t *testing.T) {
		k := &value{Name: "key"}
		m := map[*value]string{k: "v"}
		require.Empty(t, collectVisits[*value](t, m))
	})

	t.Run("handles a typed nil pointer root", func(t *testing.T) {
		require.NotPanics(t, func() {
			require.Empty(t, collectVisits[*value](t, (*value)(nil)))
		})
	})

	t.Run("handles nil graph", func(t *testing.T) {
		require.NotPanics(t, func() {
			err := VisitObjectGraph(nil, func(path []any, value any) error { return nil })
			require.Error(t, err)
		})
	})

	t.Run("visits pointers reachable from a non-pointer root", func(t *testing.T) {
		type outer struct {
			Inner *value
		}
		inner := &value{Name: "inner"}
		o := outer{Inner: inner} // passed by value
		require.Len(t, collectVisits[*value](t, o), 1)
	})

	t.Run("stops traversing on StopTraversing", func(t *testing.T) {
		count := 0
		type outer struct {
			Inner    *value
			PtrSlice []*value
		}
		o := &outer{Inner: &value{}, PtrSlice: []*value{{}, {}}}
		err := VisitObjectGraph(o, func(path []any, value any) error {
			count++
			return StopTraversing
		})
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("returns visitor errors", func(t *testing.T) {
		type outer struct {
			Inner *value
		}
		boom := errors.New("boom")
		o := &outer{Inner: &value{}}
		err := VisitObjectGraph(o, func(path []any, value any) error {
			return boom
		})
		require.ErrorIs(t, err, boom)
	})

	t.Run("records the path to nested fields", func(t *testing.T) {
		type inner struct {
			Leaf *value
		}
		type outer struct {
			Mid *inner
		}
		o := &outer{Mid: &inner{Leaf: &value{Name: "leaf"}}}
		visits := collectVisits[*value](t, o)
		require.Len(t, visits, 1)
		require.Equal(t, "<outer>.Mid.Leaf", visits[0].path)
	})

	t.Run("records slice indices in the path", func(t *testing.T) {
		type outer struct {
			Items []*value
		}
		o := &outer{Items: []*value{{Name: "a"}, {Name: "b"}}}
		visits := collectVisits[*value](t, o)
		require.Len(t, visits, 2)
		require.Equal(t, "<outer>.Items[0]", visits[0].path)
		require.Equal(t, "<outer>.Items[1]", visits[1].path)
	})

	t.Run("records map keys in the path", func(t *testing.T) {
		type outer struct {
			M map[string]*value
		}
		o := &outer{M: map[string]*value{"k": {Name: "v"}}}
		visits := collectVisits[*value](t, o)
		require.Len(t, visits, 1)
		require.Equal(t, "<outer>.M/k", visits[0].path)
	})

	t.Run("omits anonymous fields from the path", func(t *testing.T) {
		type base struct {
			Leaf *value
		}
		type derived struct {
			base
		}
		o := &derived{base: base{Leaf: &value{Name: "leaf"}}}
		visits := collectVisits[*value](t, o)
		require.Len(t, visits, 1)
		require.Equal(t, "<derived>.Leaf", visits[0].path)
	})

	t.Run("visits a pointer but not its target when both match the visitor type", func(t *testing.T) {
		type leaf struct{} // no fields, so it produces no additional visits
		l := &leaf{}
		var pointerVisits, valueVisits int
		err := VisitObjectGraph(l, func(path []any, v any) error {
			switch v.(type) {
			case *leaf:
				pointerVisits++
			case leaf:
				valueVisits++
			}
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 1, pointerVisits)
		require.Equal(t, 0, valueVisits)
	})

	t.Run("supports into arrays", func(t *testing.T) {
		type outer struct {
			Arr [2]*value
		}
		o := &outer{Arr: [2]*value{{Name: "a"}, {Name: "b"}}}
		require.Len(t, collectVisits[*value](t, o), 2)
	})

	t.Run("accepts a reflect.Value as the graph", func(t *testing.T) {
		inner := &value{Name: "inner"}
		require.Len(t, collectVisits[*value](t, reflect.ValueOf(inner)), 1)
	})

	t.Run("handles nil interface fields", func(t *testing.T) {
		type outer struct {
			Iface any
		}
		o := &outer{Iface: nil}
		require.NotPanics(t, func() {
			require.Empty(t, collectVisits[*value](t, o))
		})
	})

	t.Run("does not pass values behind unexported fields to the visitor", func(t *testing.T) {
		type outer struct {
			hidden *value
		}
		o := &outer{hidden: &value{Name: "secret"}}
		require.Empty(t, collectVisits[*value](t, o))
	})

	t.Run("visits struct values in maps without addressing them", func(t *testing.T) {
		m := map[string]value{"k": {Name: "v"}}
		visits := collectVisits[value](t, m)
		require.Len(t, visits, 1)
		require.Equal(t, "v", visits[0].value.(value).Name)
	})

	t.Run("visits pointers in nested slices", func(t *testing.T) {
		type outer struct {
			Grid [][]*value
		}
		o := &outer{Grid: [][]*value{{{Name: "a"}}, {{Name: "b"}, {Name: "c"}}}}
		require.Len(t, collectVisits[*value](t, o), 3)
	})
}

func enableDebug(t *testing.T) {
	t.Helper()
	defaultDebug := Debug
	Debug = true
	t.Cleanup(func() {
		Debug = defaultDebug
	})
}

type visitRecord struct {
	path  string
	value any
}

// collectVisits traverses the graph, recording each visited value with its path rendered at visit time
func collectVisits[T any](t *testing.T, graph any) []visitRecord {
	t.Helper()
	var visits []visitRecord
	err := VisitObjectGraph(graph, func(path []any, value T) error {
		visits = append(visits, visitRecord{path: StringifyPath(path), value: value})
		return nil
	})
	require.NoError(t, err)
	return visits
}
