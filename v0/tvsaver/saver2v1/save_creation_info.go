// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"fmt"
	"io"

	"github.com/swinslow/spdx-go/v0/spdx"
)

func renderCreationInfo2_1(ci *spdx.CreationInfo2_1, w io.Writer) error {
	if ci.SPDXVersion != "" {
		fmt.Fprintf(w, "SPDXVersion: %s\n", textify(ci.SPDXVersion))
	}
	if ci.DataLicense != "" {
		fmt.Fprintf(w, "DataLicense: %s\n", textify(ci.DataLicense))
	}
	if ci.SPDXIdentifier != "" {
		fmt.Fprintf(w, "SPDXID: %s\n", textify(ci.SPDXIdentifier))
	}
	if ci.DocumentName != "" {
		fmt.Fprintf(w, "DocumentName: %s\n", textify(ci.DocumentName))
	}
	if ci.DocumentNamespace != "" {
		fmt.Fprintf(w, "DocumentNamespace: %s\n", textify(ci.DocumentNamespace))
	}
	for _, s := range ci.ExternalDocumentReferences {
		fmt.Fprintf(w, "ExternalDocumentRef: %s\n", textify(s))
	}
	if ci.LicenseListVersion != "" {
		fmt.Fprintf(w, "LicenseListVersion: %s\n", textify(ci.LicenseListVersion))
	}
	for _, s := range ci.CreatorPersons {
		fmt.Fprintf(w, "Creator: Person: %s\n", textify(s))
	}
	for _, s := range ci.CreatorOrganizations {
		fmt.Fprintf(w, "Creator: Organization: %s\n", textify(s))
	}
	for _, s := range ci.CreatorTools {
		fmt.Fprintf(w, "Creator: Tool: %s\n", textify(s))
	}
	if ci.Created != "" {
		fmt.Fprintf(w, "Created: %s\n", textify(ci.Created))
	}
	if ci.CreatorComment != "" {
		fmt.Fprintf(w, "CreatorComment: %s\n", textify(ci.CreatorComment))
	}
	if ci.DocumentComment != "" {
		fmt.Fprintf(w, "DocumentComment: %s\n", textify(ci.DocumentComment))
	}

	// add blank newline b/c end of a main section
	fmt.Fprintf(w, "\n")

	return nil
}
