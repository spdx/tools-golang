package v3_0

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestParseLicenseExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       AnyLicenseInfo
		wantErr    bool
	}{
		{
			name:       "simple license",
			expression: "MIT",
			want:       &ListedLicense{Name: "MIT"},
		},
		{
			name:       "license with dots",
			expression: "Apache-2.0",
			want:       &ListedLicense{Name: "Apache-2.0"},
		},
		{
			name:       "LicenseRef",
			expression: "LicenseRef-custom-1.0",
			want:       &CustomLicense{ID: "LicenseRef-custom-1.0"},
		},
		{
			name:       "DocumentRef",
			expression: "DocumentRef-ext:LicenseRef-custom",
			want:       &CustomLicense{ID: "DocumentRef-ext:LicenseRef-custom"},
		},
		{
			name:       "NONE",
			expression: "NONE",
			want:       IndividualLicensingInfo_NoneLicense,
		},
		{
			name:       "NOASSERTION",
			expression: "NOASSERTION",
			want:       IndividualLicensingInfo_NoAssertionLicense,
		},
		{
			name:       "NONE within expression",
			expression: "MIT OR NONE",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					IndividualLicensingInfo_NoneLicense,
				},
			},
		},
		{
			name:       "NOASSERTION within expression",
			expression: "MIT AND NOASSERTION",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					IndividualLicensingInfo_NoAssertionLicense,
				},
			},
		},
		{
			name:       "or later",
			expression: "GPL-2.0-only+",
			want: &OrLaterOperator{
				SubjectLicense: &ListedLicense{Name: "GPL-2.0-only"},
			},
		},
		{
			name:       "OR",
			expression: "MIT OR Apache-2.0",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
				},
			},
		},
		{
			name:       "AND",
			expression: "MIT AND Apache-2.0",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
				},
			},
		},
		{
			name:       "WITH",
			expression: "GPL-2.0-only WITH Classpath-exception-2.0",
			want: &WithAdditionOperator{
				SubjectExtendableLicense: &ListedLicense{Name: "GPL-2.0-only"},
				SubjectAddition:          &ListedLicenseException{Name: "Classpath-exception-2.0"},
			},
		},
		{
			name:       "WITH custom addition",
			expression: "MIT WITH AdditionRef-my-exception",
			want: &WithAdditionOperator{
				SubjectExtendableLicense: &ListedLicense{Name: "MIT"},
				SubjectAddition:          &CustomLicenseAddition{ID: "AdditionRef-my-exception"},
			},
		},
		{
			name:       "or later with exception",
			expression: "GPL-2.0-only+ WITH Classpath-exception-2.0",
			want: &WithAdditionOperator{
				SubjectExtendableLicense: &OrLaterOperator{
					SubjectLicense: &ListedLicense{Name: "GPL-2.0-only"},
				},
				SubjectAddition: &ListedLicenseException{Name: "Classpath-exception-2.0"},
			},
		},
		{
			name:       "flattened OR",
			expression: "MIT OR Apache-2.0 OR BSD-3-Clause",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
					&ListedLicense{Name: "BSD-3-Clause"},
				},
			},
		},
		{
			name:       "flattened AND",
			expression: "MIT AND Apache-2.0 AND BSD-3-Clause",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
					&ListedLicense{Name: "BSD-3-Clause"},
				},
			},
		},
		{
			name:       "AND precedence over OR",
			expression: "MIT OR Apache-2.0 AND BSD-3-Clause",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ConjunctiveLicenseSet{
						Members: LicenseInfoList{
							&ListedLicense{Name: "Apache-2.0"},
							&ListedLicense{Name: "BSD-3-Clause"},
						},
					},
				},
			},
		},
		{
			name:       "parentheses",
			expression: "(MIT OR Apache-2.0) AND BSD-3-Clause",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&DisjunctiveLicenseSet{
						Members: LicenseInfoList{
							&ListedLicense{Name: "MIT"},
							&ListedLicense{Name: "Apache-2.0"},
						},
					},
					&ListedLicense{Name: "BSD-3-Clause"},
				},
			},
		},
		{
			name:       "complex nested parentheses",
			expression: "(MIT or (Apache-2.0 or LicenseRef-MIT-modified)) and GPL-2.0-only",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&DisjunctiveLicenseSet{
						Members: LicenseInfoList{
							&ListedLicense{Name: "MIT"},
							&DisjunctiveLicenseSet{
								Members: LicenseInfoList{
									&ListedLicense{Name: "Apache-2.0"},
									&CustomLicense{ID: "LicenseRef-MIT-modified"},
								},
							},
						},
					},
					&ListedLicense{Name: "GPL-2.0-only"},
				},
			},
		},
		{
			name:       "complex expression",
			expression: "MIT AND (GPL-2.0-only WITH Classpath-exception-2.0 OR Apache-2.0) AND LicenseRef-custom",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&DisjunctiveLicenseSet{
						Members: LicenseInfoList{
							&WithAdditionOperator{
								SubjectExtendableLicense: &ListedLicense{Name: "GPL-2.0-only"},
								SubjectAddition:          &ListedLicenseException{Name: "Classpath-exception-2.0"},
							},
							&ListedLicense{Name: "Apache-2.0"},
						},
					},
					&CustomLicense{ID: "LicenseRef-custom"},
				},
			},
		},
		{
			name:       "extra whitespace",
			expression: "  MIT   OR   Apache-2.0  ",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
				},
			},
		},
		{
			name:       "tabs and newlines",
			expression: "\tMIT\n\tOR\r\n\tApache-2.0\t",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&ListedLicense{Name: "Apache-2.0"},
				},
			},
		},
		{
			name:       "or-later combined with conjunction and grouped or-later",
			expression: "MIT+ AND Apache-2.0+ and (MIT OR GPL-2.0-only+) AND BSD-3-Clause",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&OrLaterOperator{SubjectLicense: &ListedLicense{Name: "MIT"}},
					&OrLaterOperator{SubjectLicense: &ListedLicense{Name: "Apache-2.0"}},
					&DisjunctiveLicenseSet{
						Members: LicenseInfoList{
							&ListedLicense{Name: "MIT"},
							&OrLaterOperator{SubjectLicense: &ListedLicense{Name: "GPL-2.0-only"}},
						},
					},
					&ListedLicense{Name: "BSD-3-Clause"},
				},
			},
		},
		{
			// makeAddition's AdditionRef- prefix handling on the addition side,
			// with a LicenseRef as the extendable subject.
			name:       "WITH on a LicenseRef with AdditionRef",
			expression: "LicenseRef-MyLic WITH AdditionRef-MyExc",
			want: &WithAdditionOperator{
				SubjectExtendableLicense: &CustomLicense{ID: "LicenseRef-MyLic"},
				SubjectAddition:          &CustomLicenseAddition{ID: "AdditionRef-MyExc"},
			},
		},
		{
			// the colon-aware ident scanner cooperates with the + suffix.
			name:       "DocumentRef with or-later suffix",
			expression: "DocumentRef-ext:LicenseRef-foo+",
			want: &OrLaterOperator{
				SubjectLicense: &CustomLicense{ID: "DocumentRef-ext:LicenseRef-foo"},
			},
		},
		{
			// WITH binds tighter than OR, so the explicit grouping matches the
			// default precedence.
			name:       "WITH inside paren-grouped OR",
			expression: "MIT OR (GPL-2.0-only WITH Classpath-exception-2.0)",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&WithAdditionOperator{
						SubjectExtendableLicense: &ListedLicense{Name: "GPL-2.0-only"},
						SubjectAddition:          &ListedLicenseException{Name: "Classpath-exception-2.0"},
					},
				},
			},
		},
		// error cases
		{
			name:       "empty",
			expression: "",
			wantErr:    true,
		},
		{
			name:       "just whitespace",
			expression: "   ",
			wantErr:    true,
		},
		{
			name:       "unclosed paren",
			expression: "(MIT",
			wantErr:    true,
		},
		{
			name:       "unexpected close paren",
			expression: "MIT)",
			wantErr:    true,
		},
		{
			name:       "missing operand after OR",
			expression: "MIT OR",
			wantErr:    true,
		},
		{
			name:       "missing operand after AND",
			expression: "MIT AND",
			wantErr:    true,
		},
		{
			name:       "missing exception after WITH",
			expression: "MIT WITH",
			wantErr:    true,
		},
		{
			name:       "double operator",
			expression: "MIT OR OR Apache-2.0",
			wantErr:    true,
		},
		{
			name:       "leading operator",
			expression: "OR MIT",
			wantErr:    true,
		},
		{
			name:       "trailing plus on group",
			expression: "(MIT OR Apache-2.0)+",
			wantErr:    true,
		},
		{
			// the exception side of WITH is a bare identifier; it cannot be parenthesized.
			name:       "parenthesized exception",
			expression: "MIT WITH (Classpath-exception-2.0)",
			wantErr:    true,
		},
		{
			// + is not allowed on the exception side; the trailing + is left as an
			// unexpected token.
			name:       "or-later suffix on exception",
			expression: "MIT WITH GPL-2.0-only+",
			wantErr:    true,
		},
		{
			// WITH requires an extendable license; a license set is not extendable,
			// so this fails the AnyExtendableLicense type assertion.
			name:       "WITH applied to a license set",
			expression: "(MIT OR Apache-2.0) WITH Classpath-exception-2.0",
			wantErr:    true,
		},
		{
			// a space before + detaches it from the identifier, leaving it as an
			// unexpected trailing token.
			name:       "space before or-later suffix",
			expression: "MIT +",
			wantErr:    true,
		},
		{
			name:       "empty parens",
			expression: "()",
			wantErr:    true,
		},
		{
			name:       "empty parens with whitespace",
			expression: "( )",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLicenseExpression(tt.expression)
			if tt.wantErr {
				require.Errorf(t, err, "expected error for input %q, got nil", tt.expression)
				return
			}
			require.NoError(t, err)
			d := cmp.Diff(tt.want, got, diffOpts()...)
			require.Empty(t, d, "ParseLicenseExpression(%q) mismatch (-want +got):\n%s", tt.expression, d)
		})
	}
}

func TestConvert23LicenseExpressionResolvesNestedRefs(t *testing.T) {
	// custom licenses registered by convert23license, as populated when the
	// document's OtherLicenses are converted.
	custom1 := &CustomLicense{ID: "LicenseRef-custom-1", Name: "Custom One", Text: "custom one text"}
	custom2 := &CustomLicense{ID: "LicenseRef-custom-2", Name: "Custom Two", Text: "custom two text"}

	tests := []struct {
		name       string
		expression string
		want       AnyLicenseInfo
	}{
		{
			name:       "ref nested in conjunctive and disjunctive sets",
			expression: "MIT AND (Apache-2.0 OR LicenseRef-custom-1)",
			want: &ConjunctiveLicenseSet{
				Members: LicenseInfoList{
					&ListedLicense{Name: "MIT"},
					&DisjunctiveLicenseSet{
						Members: LicenseInfoList{
							&ListedLicense{Name: "Apache-2.0"},
							custom1,
						},
					},
				},
			},
		},
		{
			name:       "multiple refs nested at different depths",
			expression: "LicenseRef-custom-1 OR (MIT AND LicenseRef-custom-2)",
			want: &DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					custom1,
					&ConjunctiveLicenseSet{
						Members: LicenseInfoList{
							&ListedLicense{Name: "MIT"},
							custom2,
						},
					},
				},
			},
		},
		{
			name:       "ref nested in or-later operator",
			expression: "LicenseRef-custom-1+",
			want: &OrLaterOperator{
				SubjectLicense: custom1,
			},
		},
		{
			// the subject license is resolved; the addition is left untouched.
			name:       "ref nested in with-addition operator",
			expression: "LicenseRef-custom-1 WITH AdditionRef-custom",
			want: &WithAdditionOperator{
				SubjectExtendableLicense: custom1,
				SubjectAddition:          &CustomLicenseAddition{ID: "AdditionRef-custom"},
			},
		},
		{
			name:       "unregistered ref left as placeholder",
			expression: "LicenseRef-unknown",
			want:       &CustomLicense{ID: "LicenseRef-unknown"},
		},
	}

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(
			ListedLicense{},
			CustomLicense{},
			OrLaterOperator{},
			DisjunctiveLicenseSet{},
			ConjunctiveLicenseSet{},
			WithAdditionOperator{},
			ListedLicenseException{},
			CustomLicenseAddition{},
			IndividualLicensingInfo{},
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &documentConverter{
				idMap: map[string]any{
					custom1.ID: custom1,
					custom2.ID: custom2,
				},
			}
			got := c.convert23licenseExpression(tt.expression)
			d := cmp.Diff(tt.want, got, opts...)
			require.Emptyf(t, d, "convert23licenseExpression(%q) mismatch (-want +got):\n%s", tt.expression, d)
		})
	}
}
