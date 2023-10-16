package marshal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		in       interface{}
		expected string
	}{
		{
			name:     "basic usage",
			in:       "<html>",
			expected: `"<html>"`,
		},
		{
			name: "within MarshalJSON callbacks",
			in: s1{
				s2{
					s3{
						Value: "<html>",
					},
				},
			},
			expected: `{"S2":{"S3":{"Value":"<html>"}}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := JSON(test.in)
			require.NoError(t, err)
			require.Equal(t, test.expected, string(got))
		})
	}
}

type s1 struct {
	S2 s2
}

type s2 struct {
	S3 s3
}

func (s *s2) MarshalJSON() ([]byte, error) {
	return JSON(s.S3)
}

type s3 struct {
	Value string
}

func (s *s3) MarshalJSON() ([]byte, error) {
	return JSON(s.Value)
}
