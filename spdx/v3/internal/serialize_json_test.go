package internal

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
)

type SomeOtherStruct struct {
	_  ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/SomeOtherStruct" node-kind:"http://www.w3.org/ns/shacl#BlankNodeOrIRI"`
	ID string  `iri:"@id"`
}

type CreationInfo struct {
	_  ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/CreationInfo" node-kind:"http://www.w3.org/ns/shacl#BlankNodeOrIRI"`
	ID string  `iri:"@id"`
}

func Test_outputInline(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want bool
	}{
		{
			name: "struct is CreationInfo",
			v:    reflect.ValueOf(CreationInfo{}),
			want: false,
		},
		{
			name: "struct is an element",
			v:    reflect.ValueOf(Element{}), // Assuming Element is a valid type checked by isElement
			want: false,
		},
		{
			name: "struct is neither CreationInfo nor element",
			v:    reflect.ValueOf(SomeOtherStruct{}), // Assuming SomeOtherStruct is not matched by isElement
			want: true,
		},
		{
			name: "non-struct type slice",
			v:    reflect.ValueOf([]int{1, 2, 3}),
			want: true,
		},
		{
			name: "non-struct type integer",
			v:    reflect.ValueOf(42),
			want: true,
		},
		{
			name: "nested pointer to struct CreationInfo",
			v:    reflect.ValueOf(&CreationInfo{}),
			want: false,
		},
		{
			name: "nested pointer to non-struct",
			v:    reflect.ValueOf(func() *int { i := 7; return &i }()),
			want: true,
		},
		{
			name: "nil reflect.Value",
			v:    reflect.Value{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := outputInline(tt.v); got != tt.want {
				t.Errorf("outputInline() = %v, want %v", got, tt.want)
			}
		})
	}
}
