package v3_0_1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kzantow/go-ld"
	"github.com/stretchr/testify/require"
)

func Test_convertElements(t *testing.T) {
	tests := []struct {
		name     string
		expected any
		convert  func(*documentConverter) any
	}{
		{
			name:     "creationInfo",
			expected: v301creationInfo(),
			convert: func(c *documentConverter) any {
				return c.convert23creationInfo(v23creationInfo())
			},
		},
		{
			name:     "externalDocumentRef",
			expected: v301externalMap1(),
			convert: func(c *documentConverter) any {
				return c.convert23externalDocumentRef(*v23externalDocumentRef1())
			},
		},
		{
			name:     "snippet",
			expected: first(v301snippet1(nil)),
			convert: func(c *documentConverter) any {
				return c.convert23snippet(*v23snippet1())
			},
		},
		{
			name:     "annotation",
			expected: v301annotation1(nil),
			convert: func(c *documentConverter) any {
				return c.convert23annotation(v23annotation1())
			},
		},
		{
			name:     "otherLicense",
			expected: v301customLicense1(),
			convert: func(c *documentConverter) any {
				return c.convert23license(v23customLicense1())
			},
		},
		{
			name:     "file",
			expected: first(v301file1()),
			convert: func(c *documentConverter) any {
				return c.convert23file(v23file1())
			},
		},
		{
			name:     "package",
			expected: first(v301package1()),
			convert: func(c *documentConverter) any {
				return c.convert23package(v23package1())
			},
		},
		{
			name:     "full document",
			expected: v301doc(),
			convert: func(c *documentConverter) any {
				d := &Document{}
				From_v2_3(*v23doc(), d)
				return d
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestConverter()
			got := tt.convert(c)
			if tt.expected == nil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			diff := cmp.Diff(tt.expected, got, diffOpts()...)
			if diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func each[Element, View any](s ld.TypeSeq[Element, View]) []View {
	var out []View
	for _, v := range s {
		out = append(out, v)
	}
	return out
}

func Test_documentConversion(t *testing.T) {
	expected := v301doc()

	converted := &Document{}
	From_v2_3(*v23doc(), converted)

	startPkgs := each(expected.Elements.Packages())
	gotPkgs := each(converted.Elements.Packages())
	if diff := cmp.Diff(startPkgs, gotPkgs, diffOpts()...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	startRels := each(expected.Elements.Relationships())
	gotRels := each(converted.Elements.Relationships())
	if diff := cmp.Diff(startRels, gotRels, diffOpts()...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func diffOpts() []cmp.Option {
	var out []cmp.Option
	out = append(out,
		cmpopts.IgnoreFields(Document{}, "LDContext"),
		cmpopts.IgnoreFields(Element{}, "ID", "CreationInfo"),
		cmpopts.IgnoreFields(ElementCollection{}, "Elements"), // we only want to compare RootElements
		cmpopts.IgnoreFields(CreationInfo{}, "CreatedUsing"),
		cmpopts.EquateComparable(
			ExternalIdentifierType{},
			HashAlgorithm{},
			FileKindType{},
			SoftwarePurpose{},
			PresenceType{},
			SafetyRiskAssessmentType{},
			RelationshipCompleteness{},
			RelationshipType{},
			AnnotationType{},
			ProfileIdentifierType{},
			ExternalIRI{},
		),
	)
	return out
}

func newTestConverter() *documentConverter {
	d := Document{
		LDContext: context(),
	}
	return newDocumentConverter(&d)
}

func first[T1, T2 any](v1 T1, _ T2) T1 {
	return v1
}
