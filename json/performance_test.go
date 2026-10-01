// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func BenchmarkRead(b *testing.B) {
	for _, size := range []int{10, 1000} {
		b.Run(fmt.Sprintf("packages=%d", size), func(b *testing.B) {
			doc := spdx.Document{SPDXVersion: "SPDX-2.3", SPDXIdentifier: "DOCUMENT"}
			for i := 0; i < size; i++ {
				doc.Packages = append(doc.Packages, &spdx.Package{
					PackageName:             fmt.Sprintf("package-%d", i),
					PackageSPDXIdentifier:   common.ElementID(fmt.Sprintf("package-%d", i)),
					PackageVersion:          "1.2.3",
					PackageDownloadLocation: "https://example.com/package.tar.gz",
					PackageLicenseConcluded: "Apache-2.0",
					PackageDescription:      strings.Repeat("Package description. ", 10),
				})
			}
			data, err := json.Marshal(doc)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := Read(bytes.NewReader(data)); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkMarshalPackageV22(b *testing.B) {
	for _, withCode := range []bool{false, true} {
		b.Run(fmt.Sprintf("verificationCode=%t", withCode), func(b *testing.B) {
			pkg := v2_2.Package{
				PackageName:             "example",
				PackageSPDXIdentifier:   "example",
				PackageDownloadLocation: "https://example.com/package.tar.gz",
				PackageLicenseConcluded: "Apache-2.0",
				PackageDescription:      strings.Repeat("Package description. ", 10),
			}
			if withCode {
				pkg.PackageVerificationCode.Value = strings.Repeat("a", 40)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := json.Marshal(&pkg); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
