// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spdx/tools-golang/json/marshal"
)

func Test_DocElementIDEncoding(t *testing.T) {
	tests := []struct {
		name     string
		value    DocElementID
		expected string
		err      bool
	}{
		{
			name: "ElementRefID",
			value: DocElementID{
				ElementRefID: "some-id",
			},
			expected: "SPDXRef-some-id",
		},
		{
			name: "DocumentRefID:ElementRefID",
			value: DocElementID{
				DocumentRefID: "a-doc",
				ElementRefID:  "some-id",
			},
			expected: "DocumentRef-a-doc:SPDXRef-some-id",
		},
		{
			name: "DocumentRefID no ElementRefID",
			value: DocElementID{
				DocumentRefID: "a-doc",
			},
			err: true,
		},
		{
			name: "SpecialID",
			value: DocElementID{
				SpecialID: "special-id",
			},
			expected: "special-id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := marshal.JSON(test.value)
			switch {
			case !test.err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.err && err == nil:
				t.Fatalf("expected error but got none")
			case test.err:
				return
			}
			s := string(result)
			if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
				t.Fatalf("string was not returned: %s", s)
			}
			s = strings.Trim(s, `"`)
			if test.expected != s {
				t.Fatalf("%s != %s", test.expected, s)
			}
		})
	}
}

func Test_DocElementIDDecoding(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected DocElementID
		err      bool
	}{
		{
			name:  "ElementRefID",
			value: "SPDXRef-some-id",
			expected: DocElementID{
				ElementRefID: "some-id",
			},
		},
		{
			name:  "DocumentRefID:ElementRefID",
			value: "DocumentRef-a-doc:SPDXRef-some-id",
			expected: DocElementID{
				DocumentRefID: "a-doc",
				ElementRefID:  "some-id",
			},
		},
		{
			name:  "DocumentRefID no ElementRefID",
			value: "DocumentRef-a-doc",
			expected: DocElementID{
				DocumentRefID: "a-doc",
			},
		},
		{
			name:  "DocumentRefID:ElementRefID without spdxref prefix",
			value: "DocumentRef-a-doc:some-id",
			expected: DocElementID{
				DocumentRefID: "a-doc",
				ElementRefID:  "some-id",
			},
		},
		{
			name:  "without spdxref prefix",
			value: "some-id-without-spdxref",
			expected: DocElementID{
				ElementRefID: "some-id-without-spdxref",
			},
		},
		{
			name:  "SpecialID NONE",
			value: "NONE",
			expected: DocElementID{
				SpecialID: "NONE",
			},
		},
		{
			name:  "SpecialID NOASSERTION",
			value: "NOASSERTION",
			expected: DocElementID{
				SpecialID: "NOASSERTION",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := DocElementID{}
			s := fmt.Sprintf(`"%s"`, test.value)
			err := json.Unmarshal([]byte(s), &out)
			switch {
			case !test.err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.err && err == nil:
				t.Fatalf("expected error but got none")
			case test.err:
				return
			}
			if !reflect.DeepEqual(test.expected, out) {
				t.Fatalf("unexpected value: %v != %v", test.expected, out)
			}
		})
	}
}

func Test_ElementIDEncoding(t *testing.T) {
	tests := []struct {
		name     string
		value    ElementID
		expected string
		err      bool
	}{
		{
			name:     "appends spdxref",
			value:    ElementID("some-id"),
			expected: "SPDXRef-some-id",
		},
		{
			name:     "appends spdxref",
			value:    ElementID("SPDXRef-some-id"),
			expected: "SPDXRef-some-id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := marshal.JSON(test.value)
			switch {
			case !test.err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.err && err == nil:
				t.Fatalf("expected error but got none")
			case test.err:
				return
			}
			s := string(result)
			if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
				t.Fatalf("string was not returned: %s", s)
			}
			s = strings.Trim(s, `"`)
			if test.expected != s {
				t.Fatalf("%s != %s", test.expected, s)
			}
		})
	}
}

func Test_ElementIDDecoding(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected ElementID
		err      bool
	}{
		{
			name:     "valid id",
			value:    "SPDXRef-some-id",
			expected: ElementID("some-id"),
		},
		{
			name:  "without prefix",
			value: "some-id-without-spdxref",
			expected: ElementID("some-id-without-spdxref"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var out ElementID
			s := fmt.Sprintf(`"%s"`, test.value)
			err := json.Unmarshal([]byte(s), &out)
			switch {
			case !test.err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.err && err == nil:
				t.Fatalf("expected error but got none")
			case test.err:
				return
			}
			if !reflect.DeepEqual(test.expected, out) {
				t.Fatalf("unexpected value: %v != %v", test.expected, out)
			}
		})
	}
}

func Test_ElementIDStructEncoding(t *testing.T) {
	type typ struct {
		Id ElementID `json:"id"`
	}
	tests := []struct {
		name     string
		value    typ
		expected string
		err      bool
	}{
		{
			name: "appends spdxref",
			value: typ{
				Id: ElementID("some-id"),
			},
			expected: `{"id":"SPDXRef-some-id"}`,
		},
		{
			name: "appends spdxref",
			value: typ{
				Id: ElementID("SPDXRef-some-id"),
			},
			expected: `{"id":"SPDXRef-some-id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := marshal.JSON(test.value)
			switch {
			case !test.err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.err && err == nil:
				t.Fatalf("expected error but got none")
			case test.err:
				return
			}
			s := string(result)
			if test.expected != s {
				t.Fatalf("%s != %s", test.expected, s)
			}
		})
	}
}

func Test_ElementIDStructDecoding(t *testing.T) {
	type typ struct {
		Id ElementID `json:"id"`
	}
	tests := []struct {
		name     string
		value    string
		expected typ
		err      bool
	}{
		{
			name: "valid id",
			expected: typ{
				Id: ElementID("some-id"),
			},
			value: `{"id":"SPDXRef-some-id"}`,
		},
		{
			name:  "without prefix",
			expected: typ{
				Id: ElementID("some-id"),
			},
			value: `{"id":"some-id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := typ{}
			err := json.Unmarshal([]byte(test.value), &out)
			switch {
			case !test.err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case test.err && err == nil:
				t.Fatalf("expected error but got none")
			case test.err:
				return
			}
			if !reflect.DeepEqual(test.expected, out) {
				t.Fatalf("unexpected value: %v != %v", test.expected, out)
			}
		})
	}
}
