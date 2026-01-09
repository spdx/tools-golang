package ld

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getId(t *testing.T) {
	tests := []struct {
		name     string
		instance func() any
		expected string
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name: "no id string",
			instance: func() any {
				return ""
			},
			expected: "",
		},
		{
			name: "an id string",
			instance: func() any {
				return "an-id"
			},
			expected: "an-id",
		},
		{
			name: "no id struct",
			instance: func() any {
				type ty struct {
					id string `iri:"not_id"`
				}
				return ty{
					id: "a-val",
				}
			},
			expected: "",
			wantErr:  require.Error,
		},
		{
			name: "direct id struct",
			instance: func() any {
				type ty struct {
					id string `iri:"@id"`
				}
				return ty{
					id: "a-val",
				}
			},
			expected: "a-val",
		},
		{
			name: "exported single embedded id struct",
			instance: func() any {
				type Ty struct {
					Id string `iri:"@id"`
				}
				type ta struct {
					Ty
				}
				return ta{
					Ty: Ty{
						Id: "ta-val",
					},
				}
			},
			expected: "ta-val",
		},
		{
			name: "unexported embedded id struct",
			instance: func() any {
				type ty struct {
					id string `iri:"@id"`
				}
				type tb struct {
					ty
				}
				return tb{
					ty: ty{
						id: "tb-val",
					},
				}
			},
			expected: "tb-val",
		},
		{
			name: "multiple embedded id struct",
			instance: func() any {
				type Ta struct {
					Name string `iri:"not_id"`
				}
				type Tb struct {
					ID string `iri:"@id"`
				}
				type Ty struct {
					Ta
					Tb
				}
				return Ty{
					Ta: Ta{
						Name: "no-id",
					},
					Tb: Tb{
						ID: "tb-val",
					},
				}
			},
			expected: "tb-val",
		},
		{
			name: "recursive embedded id struct",
			instance: func() any {
				type Ty struct {
					ID string `iri:"@id"`
				}
				type Tb struct {
					Ty
				}
				type Tc struct {
					Tb
				}
				return Tc{
					Tb: Tb{
						Ty: Ty{
							ID: "tc-val",
						},
					},
				}
			},
			expected: "tc-val",
		},
		{
			name: "invalid data type returns error",
			instance: func() any {
				return 1
			},
			expected: "",
			wantErr:  require.Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel() // this function must be safe to call in parallel
			v := test.instance()
			got, err := GetID(v)
			wantErr := test.wantErr
			if wantErr == nil {
				wantErr = require.NoError
			}
			wantErr(t, err)
			require.Equal(t, test.expected, got)
		})
	}
}
