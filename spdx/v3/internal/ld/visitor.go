package ld

import (
	"fmt"
	"reflect"
)

var StopTraversing = fmt.Errorf("stop-traversing-graph")

// VisitObjectGraph traverses the object graph, taking into account cycles, calling the visitor function for each
// step along the traversal, including field properties, pointer and subsequent struct values, elements in
// slices and both keys and values of maps, as well as some context such as the path within the graph and any
// containing struct field. The value is always able to have Interface() and Set() called.
func VisitObjectGraph(graph any, visitor func(path []any, value reflect.Value) error) error {
	t := reflect.TypeOf(graph)
	return visitObjectGraph(map[reflect.Value]struct{}{}, []any{baseType(t)}, reflect.ValueOf(graph), visitor)
}

func visitObjectGraph(visited map[reflect.Value]struct{}, path []any, v reflect.Value, visitor func([]any, reflect.Value) error) error {
	if !v.IsValid() {
		return nil
	}
	if _, ok := visited[v]; ok {
		return nil
	}
	visited[v] = struct{}{}

	var err error
	if v.CanInterface() {
		err = visitor(path, v)
		if err == StopTraversing {
			return nil
		} else if err != nil {
			return err
		}
	}

	t := v.Type()

	switch t.Kind() {
	case reflect.Interface:
		return visitObjectGraph(visited, path, v.Elem(), visitor)
	case reflect.Pointer:
		return visitObjectGraph(visited, path, v.Elem(), visitor)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			subPath := path[:]
			if !f.Anonymous {
				subPath = append(subPath, f)
			}
			fv := v.Field(i)
			err = visitObjectGraph(visited, subPath, fv, visitor)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := v.MapRange()
		if iter == nil {
			return nil
		}
		for iter.Next() {
			path := append(path[:], fmt.Sprintf("%v", iter.Key().Interface()))
			err = visitObjectGraph(visited, path, iter.Key(), visitor)
			if err != nil {
				return err
			}
			err = visitObjectGraph(visited, path, iter.Value(), visitor)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err = visitObjectGraph(visited, append(path[:], i), v.Index(i), visitor)
			if err != nil {
				return err
			}
		}
	default:
	}
	return nil
}
