package v3_0

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx"
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
			name: "licenseExpression",
			expected: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
				},
			},
			convert: func(c *documentConverter) any {
				return c.convert23licenseExpression("MIT OR Apache-2.0")
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
			name: "small document with packages",
			expected: func() *Document {
				p1 := &Package{
					ID:   "pkg-1",
					Name: "pkg-1",
				}
				p2 := &Package{
					ID:   "pkg-2",
					Name: "pkg-2",
				}
				r1 := &Relationship{
					Type: RelationshipType_Contains,
					From: p1,
					To:   ElementList{p2},
				}
				sbom := &SBOM{
					Elements: ElementList{
						p1,
						p2,
						r1,
					},
					RootElements: ElementList{
						p1,
					},
				}
				return &Document{
					SpdxDocument: SpdxDocument{
						NamespaceMaps:       NamespaceMapList{&NamespaceMap{Namespace: "https://spdx.org/spdxdocs/#", Prefix: "SPDXRef"}},
						ProfileConformances: []ProfileIdentifierType{ProfileIdentifierType_Core, ProfileIdentifierType_Software},
						RootElements: ElementList{
							sbom,
						},
						// all collected elements:
						Elements: ElementList{
							sbom,
							p1,
							p2,
							r1,
						},
					},
				}
			}(),
			convert: func(c *documentConverter) any {
				d := &Document{}
				From_v2_3(spdx.Document{
					Packages: []*spdx.Package{
						{
							PackageSPDXIdentifier: "pkg-1",
							PackageName:           "pkg-1",
						},
						{
							PackageSPDXIdentifier: "pkg-2",
							PackageName:           "pkg-2",
						},
					},
					Relationships: []*spdx.Relationship{
						{
							Relationship: spdx.RelationshipDescribes,
							RefA:         spdx.DocElementID{ElementRefID: "DOCUMENT"},
							RefB:         spdx.DocElementID{ElementRefID: "pkg-1"},
						},
						{
							Relationship: spdx.RelationshipContains,
							RefA:         spdx.DocElementID{ElementRefID: "pkg-1"},
							RefB:         spdx.DocElementID{ElementRefID: "pkg-2"},
						},
					},
				}, d)
				return d
			},
		},
		{
			name: "small document with annotations",
			expected: func() *Document {
				a1 := &Annotation{
					ID:        "annotation-1",
					Type:      AnnotationType_Review,
					Statement: "annotation-comment-1",
				}
				sbom := &SBOM{
					Elements: ElementList{
						a1,
					},
					RootElements: nil,
				}
				return &Document{
					SpdxDocument: SpdxDocument{
						NamespaceMaps:       NamespaceMapList{&NamespaceMap{Namespace: "https://spdx.org/spdxdocs/#", Prefix: "SPDXRef"}},
						ProfileConformances: []ProfileIdentifierType{ProfileIdentifierType_Core, ProfileIdentifierType_Software},
						RootElements: ElementList{
							sbom,
						},
						// all collected elements:
						Elements: ElementList{
							sbom,
							a1,
						},
					},
				}
			}(),
			convert: func(c *documentConverter) any {
				d := &Document{}
				From_v2_3(spdx.Document{
					Annotations: []*spdx.Annotation{
						{
							AnnotationSPDXIdentifier: spdx.DocElementID{ElementRefID: "annotation-1"},
							AnnotationType:           "REVIEW",
							AnnotationComment:        "annotation-comment-1",
						},
					},
					Relationships: []*spdx.Relationship{},
				}, d)
				return d
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

func Test_documentConversion(t *testing.T) {
	expected := v301doc()

	converted := &Document{}
	From_v2_3(*v23doc(), converted)

	opts := diffOpts()

	convertedElements := ElementList(slices.Collect(maps.Values(collectAllElements(&converted.SpdxDocument))))
	expectedElements := ElementList(slices.Collect(maps.Values(collectAllElements(&expected.SpdxDocument))))
	tests := []struct {
		name     string
		expected any
		got      any
	}{
		{"Packages", expectedElements.Packages(), convertedElements.Packages()},
		{"Files", expectedElements.Files(), convertedElements.Files()},
		{"Snippets", expectedElements.Snippets(), convertedElements.Snippets()},
		{"Annotations", expectedElements.Annotations(), convertedElements.Annotations()},
		{"Relationships", expectedElements.Relationships(), convertedElements.Relationships()},
		{"CustomLicenses", expectedElements.CustomLicenses(), convertedElements.CustomLicenses()},
		{"ListedLicenses", expectedElements.ListedLicenses(), convertedElements.ListedLicenses()},
		{"DisjunctiveLicenseSets", expectedElements.DisjunctiveLicenseSets(), convertedElements.DisjunctiveLicenseSets()},
		{"ConjunctiveLicenseSets", expectedElements.ConjunctiveLicenseSets(), convertedElements.ConjunctiveLicenseSets()},
		{"WithAdditionOperators", expectedElements.WithAdditionOperators(), convertedElements.WithAdditionOperators()},
		{"ListedLicenseExceptions", expectedElements.ListedLicenseExceptions(), convertedElements.ListedLicenseExceptions()},
		{"OrLaterOperators", expectedElements.OrLaterOperators(), convertedElements.OrLaterOperators()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.expected)
			require.NotEmpty(t, tt.got)
			if diff := cmp.Diff(tt.expected, tt.got, opts...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func diffOpts() []cmp.Option {
	var out []cmp.Option
	for _, t := range []any{
		Package{},
		AIPackage{},
		Relationship{},
		File{},
		Snippet{},
		Annotation{},
		Tool{},
		Person{},
		Organization{},
		SoftwareAgent{},
		CustomLicense{},
		ListedLicense{},
		OrLaterOperator{},
		DisjunctiveLicenseSet{},
		ConjunctiveLicenseSet{},
		WithAdditionOperator{},
		ListedLicenseException{},
		CustomLicenseAddition{},
		SpdxDocument{},
		SBOM{},
	} {
		out = append(out,
			cmpopts.IgnoreUnexported(t),
			cmpopts.IgnoreFields(t, "ID", "CreationInfo"),
		)
	}
	out = append(out,
		cmpopts.SortSlices(func(a, b any) int {
			return strings.Compare(sortKey(a), sortKey(b))
		}),
		cmpopts.IgnoreFields(Document{}, "LDContext"),
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
			LifecycleScopeType{},
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

func sortKey(v any) string {
	return elementSortKey(v, 0)
}

func elementSortKey(v any, depth int) string {
	if depth > 4 {
		return reflect.TypeOf(v).String()
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", v)
	}
	typeName := rv.Type().Name()
	var parts []string
	for i := 0; i < rv.NumField(); i++ {
		ft := rv.Type().Field(i)
		if !ft.IsExported() || ft.Name == "ID" || ft.Name == "CreationInfo" {
			continue
		}
		f := rv.Field(i)
		if f.IsZero() {
			continue
		}
		switch f.Kind() {
		case reflect.String:
			parts = append(parts, ft.Name+"="+f.String())
		case reflect.Struct:
			parts = append(parts, ft.Name+"="+fmt.Sprintf("%v", f.Interface()))
		case reflect.Interface:
			if !f.IsNil() {
				parts = append(parts, ft.Name+"="+elementSortKey(f.Interface(), depth+1))
			}
		case reflect.Slice:
			if f.Len() > 0 {
				var elems []string
				for j := 0; j < f.Len(); j++ {
					elems = append(elems, elementSortKey(f.Index(j).Interface(), depth+1))
				}
				parts = append(parts, ft.Name+"=["+strings.Join(elems, ",")+"]")
			}
		case reflect.Ptr:
			if !f.IsNil() {
				parts = append(parts, ft.Name+"="+elementSortKey(f.Interface(), depth+1))
			}
		default:
			parts = append(parts, fmt.Sprintf("%s=%v", ft.Name, f.Interface()))
		}
	}
	return typeName + "{" + strings.Join(parts, ",") + "}"
}

func first[T1, T2 any](v1 T1, _ T2) T1 {
	return v1
}
