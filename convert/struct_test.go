// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FromPtr(t *testing.T) {
	type t1 struct{}
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "struct",
			input:    t1{},
			expected: t1{},
		},
		{
			name:     "ptr",
			input:    &t1{},
			expected: t1{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := FromPtr(test.input)
			assert.Equal(t, test.expected, out)
		})
	}
}

func Test_Describe(t *testing.T) {
	type t1 struct {
		text string
	}
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "struct",
			input: t1{
				text: "some-text",
			},
			expected: "github.com/spdx/tools-golang/convert.t1: {text:some-text}",
		},
		{
			name: "ptr",
			input: &t1{
				text: "some-text",
			},
			expected: "github.com/spdx/tools-golang/convert.*t1: {text:some-text}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := Describe(test.input)
			assert.Equal(t, test.expected, out)
		})
	}
}

func Test_limit(t *testing.T) {
	tests := []struct {
		expected string
		input    string
		length   int
	}{
		{
			expected: "abc",
			input:    "abc",
			length:   3,
		},
		{
			expected: "abc...",
			input:    "abcdefg",
			length:   3,
		},
		{
			expected: "abcdef",
			input:    "abcdef",
			length:   3,
		},
		{
			expected: "abcd",
			input:    "abcd",
			length:   -1,
		},
		{
			expected: "",
			input:    "",
			length:   100,
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			out := limit(test.input, test.length)
			assert.Equal(t, test.expected, out)
		})
	}
}
