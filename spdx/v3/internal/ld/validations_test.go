package ld

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name  string
		graph any
		errs  int
	}{
		{
			name:  "no validators",
			graph: 1,
			errs:  0,
		},
		{
			name: "no validators slice struct",
			graph: []any{
				time.Time{},
			},
			errs: 0,
		},
		{
			name: "one invalid",
			graph: []any{
				PositiveInt(1),
				notValid("invalid"),
			},
			errs: 1,
		},
		{
			name: "multiple invalid",
			graph: []any{
				notValid("invalid1"),
				notValid("invalid2"),
			},
			errs: 2,
		},
		{
			name: "positive int valid",
			graph: []any{
				ptr(PositiveInt(1)),
			},
			errs: 0,
		},
		{
			name: "positive int invalid",
			graph: []any{
				ptr(PositiveInt(-1)),
			},
			errs: 1,
		},
		{
			name: "non negative int valid",
			graph: []any{
				ptr(NonNegativeInt(1)),
			},
			errs: 0,
		},
		{
			name: "non negative int invalid",
			graph: []any{
				ptr(NonNegativeInt(-1)),
			},
			errs: 1,
		},
		{
			name: "single invalid ptr",
			graph: []any{
				ptr(notValidP("invalid")),
			},
			errs: 1,
		},
		{
			name: "nested valid",
			graph: []any{
				ptr(notValidP("invalid")),
			},
			errs: 1,
		},
		{
			name: "nested single invalid",
			graph: []any{
				ptr(nested{ // valid
					PositiveVal:    -1, // invalid
					NonNegativeVal: 8,  // valid
				}),
			},
			errs: 1,
		},
		{
			name: "nested multiple invalid",
			graph: []any{
				ptr(nested{ // invalid
					PositiveVal:    11, // valid
					NonNegativeVal: -1, // invalid
				}),
			},
			errs: 2,
		},
		{
			name:  "time valid",
			graph: time.Now(),
			errs:  0,
		},
		{
			name:  "DateTime valid",
			graph: DateTime(time.Now()),
			errs:  0,
		},
		{
			name: "Pattern valid",
			graph: validateExpression{
				value:   "text/plain",
				pattern: "^[^\\/]+\\/[^\\/]+$",
			},
			errs: 0,
		},
		{
			name: "Pattern invalid",
			graph: validateExpression{
				value:   "text-plain",
				pattern: "^[^\\/]+\\/[^\\/]+$",
			},
			errs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGraph(tt.graph)
			require.Len(t, flattenErrors(err), tt.errs)
		})
	}
}

type notValid string

func (n notValid) Validate() error {
	return fmt.Errorf(string(n))
}

type notValidP string

func (n *notValidP) Validate() error {
	return fmt.Errorf(string(*n))
}

func ptr[T any](v T) *T {
	return &v
}

type nested struct {
	PositiveVal    PositiveInt
	NonNegativeVal NonNegativeInt
}

func (v *nested) Validate() error {
	if v.PositiveVal >= 10 || v.NonNegativeVal >= 10 {
		return fmt.Errorf("values should be < 10")
	}
	return nil
}

type validateExpression struct {
	value   string
	pattern string
}

func (v validateExpression) Validate() error {
	return ValidateExpression(v.pattern)(v.value)
}
