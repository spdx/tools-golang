package ld

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_visitor(t *testing.T) {
	type typ struct {
		Name  string
		Slice []int
		Map   map[string]int
	}

	v := &typ{
		Name:  "a name",
		Slice: []int{2, 4, 6},
		Map: map[string]int{
			"key1": 10,
			"key2": 20,
		},
	}

	err := VisitObjectGraph(v, func(path []any, value reflect.Value) error {
		f := path[len(path)-1]
		if f, ok := f.(reflect.StructField); ok {
			if f.Name == "Name" {
				require.Equal(t, value.String(), "a name")
				value.SetString("a new name")
			}
		}
		//fmt.Println(path)
		//if intSlice, ok := value.Interface().([]int); ok {
		//	fmt.Println(intSlice)
		//	//intSlice = append(intSlice, 12)
		//	newValue := reflect.Append(value, reflect.ValueOf(12))
		//	value.Set(newValue)
		//}
		if mapVal, ok := value.Interface().(map[string]int); ok {
			fmt.Println(mapVal)
			mapVal["new"] = 100
		}
		return nil
	})
	fmt.Println(v)
	require.NoError(t, err)
}
