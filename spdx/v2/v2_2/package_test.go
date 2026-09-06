// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package v2_2

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx/v2/common"
)

func TestPackageMarshalJSONVerificationCode(t *testing.T) {
	for _, tc := range []struct {
		name string
		code common.PackageVerificationCode
		omit bool
	}{
		{"absent", common.PackageVerificationCode{}, true},
		{"value", common.PackageVerificationCode{Value: "abcd"}, false},
		{"empty exclusions", common.PackageVerificationCode{ExcludedFiles: []string{}}, false},
		{"exclusions", common.PackageVerificationCode{ExcludedFiles: []string{"./sbom.json"}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pkg := Package{
				PackageName:               "example <&>",
				PackageSPDXIdentifier:     "example",
				PackageVerificationCode:   tc.code,
				PackageSupplier:           &common.Supplier{Supplier: "NOASSERTION"},
				PackageExternalReferences: []*PackageExternalReference{{Category: "PACKAGE-MANAGER", RefType: "purl", Locator: "pkg:generic/example"}},
				Files:                     []*File{{FileName: "./example", FileSPDXIdentifier: "file"}},
			}
			// Compare every field against ordinary struct encoding, with only
			// the verification field removed when absent.
			type plainPackage Package
			encoded, err := json.Marshal(plainPackage(pkg))
			require.NoError(t, err)
			var want map[string]any
			require.NoError(t, json.Unmarshal(encoded, &want))
			if tc.omit {
				delete(want, "packageVerificationCode")
			}
			for _, value := range []any{pkg, &pkg} {
				data, err := json.Marshal(value)
				require.NoError(t, err)
				var got map[string]any
				require.NoError(t, json.Unmarshal(data, &got))
				require.Equal(t, want, got)
			}
			data, err := pkg.MarshalJSON()
			require.NoError(t, err)
			require.Contains(t, string(data), "example <&>")
		})
	}
}

func TestPackageMarshalJSONInvalidSupplier(t *testing.T) {
	for _, code := range []string{"", "abcd"} {
		pkg := Package{PackageSupplier: &common.Supplier{}, PackageVerificationCode: common.PackageVerificationCode{Value: code}}
		_, err := json.Marshal(pkg)
		require.ErrorContains(t, err, "failed to marshal invalid Supplier")
	}
}
