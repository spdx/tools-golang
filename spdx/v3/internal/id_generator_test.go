package internal

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
)

type Element struct {
	_  ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/Element" node-kind:"http://www.w3.org/ns/shacl#IRI"`
	ID string  `iri:"@id"`
}

type Parent struct {
	Element
}

type NotElement struct {
	_  ld.Type `iri:"https://spdx.org/rdf/3.0.1/terms/Core/IndividualElement" node-kind:"http://www.w3.org/ns/shacl#BlankNodeOrIRI"`
	ID string  `iri:"@id"`
}

func Test_prefixedIdGenerator(t *testing.T) {
	tests := []struct {
		name      string
		iriPrefix string
		prefixes  map[string]string
		id        string
		value     reflect.Value
		expected  string
	}{
		{
			name:      "alternate namespace prefix to a URI matching prefix map",
			iriPrefix: "ex",
			prefixes: map[string]string{
				"http://example.org/": "pfx2",
			},
			id:       "http://example.org/resource",
			value:    reflect.ValueOf(&Element{}),
			expected: "pfx2:resource",
		},
		{
			name:      "same namespace prefix to a URI matching prefix map",
			iriPrefix: "ex",
			prefixes: map[string]string{
				"http://example.org/": "ex",
			},
			id:       "http://example.org/resource",
			value:    reflect.ValueOf(&Element{}),
			expected: "ex:resource",
		},
		{
			name:      "blank node if allowed",
			iriPrefix: "ex",
			prefixes:  map[string]string{},
			id:        "",
			value:     reflect.ValueOf(&NotElement{}),
			expected:  "_:NotElement-1",
		},
		{
			name:      "URI prefix if required",
			iriPrefix: "ex",
			prefixes:  map[string]string{},
			id:        "",
			value:     reflect.ValueOf(&Element{}),
			expected:  "ex:Element-1",
		},
		{
			name:      "URI if required parent",
			iriPrefix: "ex",
			prefixes:  map[string]string{},
			id:        "",
			value:     reflect.ValueOf(&Parent{}),
			expected:  "ex:Parent-1",
		},
		{
			name:      "blank node when not a URI",
			iriPrefix: "ex",
			prefixes:  map[string]string{},
			id:        "id123",
			value:     reflect.ValueOf(&NotElement{}),
			expected:  "_:id123",
		},
		{
			name:      "prefixed URI if ID is non-blank",
			iriPrefix: "ex",
			prefixes:  map[string]string{},
			id:        "id123",
			value:     reflect.ValueOf(&Element{}),
			expected:  "ex:id123",
		},
		{
			name:      "keep URI if valid",
			iriPrefix: "ex",
			prefixes: map[string]string{
				"http://example.org/": "ex",
			},
			id:       "http://another.org/resource",
			value:    reflect.ValueOf(&NotElement{}),
			expected: "http://another.org/resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := PrefixedIdGenerator(tt.iriPrefix, tt.prefixes)
			result := generator(tt.id, tt.value)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func Test_uuidGenerator(t *testing.T) {
	tests := []struct {
		name       string
		existingId string
	}{
		{
			name:       "empty existing ID",
			existingId: "",
		},
		{
			name:       "with existing ID",
			existingId: "existing-id",
		},
		{
			name:       "with UUID existing ID",
			existingId: "urn:uuid:12345678-1234-1234-1234-123456789012",
		},
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := uuidGenerator()
			result := gen(tt.existingId, reflect.Value{})

			// Verify format
			if !strings.HasPrefix(result, "urn:uuid:") {
				t.Errorf("expected UUID to start with 'urn:uuid:', got %s", result)
			}

			// Verify UUID format (8-4-4-4-12 hex digits)
			uuidPart := strings.TrimPrefix(result, "urn:uuid:")
			if !uuidRegex.MatchString(uuidPart) {
				t.Errorf("expected valid UUID format, got %s", uuidPart)
			}
		})
	}

	// Test uniqueness separately
	t.Run("generates unique UUIDs", func(t *testing.T) {
		gen := uuidGenerator()
		result1 := gen("", reflect.Value{})
		result2 := gen("", reflect.Value{})

		if result1 == result2 {
			t.Errorf("expected unique UUIDs, got duplicate: %s", result1)
		}
	})
}
