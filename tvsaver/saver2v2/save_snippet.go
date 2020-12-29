// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

func renderSnippet2_2(sn *spdx.Snippet2_2, w io.Writer) error {
	if sn.SPDXIdentifier != "" {
		fmt.Fprintf(w, "SnippetSPDXIdentifier: %s\n", spdx.RenderElementID(sn.SPDXIdentifier))
	}
	snFromFileIDStr := spdx.RenderDocElementID(sn.SnippetFromFileSPDXIdentifier)
	if snFromFileIDStr != "" {
		fmt.Fprintf(w, "SnippetFromFileSPDXID: %s\n", snFromFileIDStr)
	}
	if sn.ByteRangeStart != 0 && sn.ByteRangeEnd != 0 {
		fmt.Fprintf(w, "SnippetByteRange: %d:%d\n", sn.ByteRangeStart, sn.ByteRangeEnd)
	}
	if sn.LineRangeStart != 0 && sn.LineRangeEnd != 0 {
		fmt.Fprintf(w, "SnippetLineRange: %d:%d\n", sn.LineRangeStart, sn.LineRangeEnd)
	}
	if sn.LicenseConcluded != "" {
		fmt.Fprintf(w, "SnippetLicenseConcluded: %s\n", sn.LicenseConcluded)
	}
	for _, s := range sn.LicenseInfoInSnippet {
		fmt.Fprintf(w, "LicenseInfoInSnippet: %s\n", s)
	}
	if sn.LicenseComments != "" {
		fmt.Fprintf(w, "SnippetLicenseComments: %s\n", textify(sn.LicenseComments))
	}
	if sn.CopyrightText != "" {
		fmt.Fprintf(w, "SnippetCopyrightText: %s\n", sn.CopyrightText)
	}
	if sn.Comment != "" {
		fmt.Fprintf(w, "SnippetComment: %s\n", textify(sn.Comment))
	}
	if sn.Name != "" {
		fmt.Fprintf(w, "SnippetName: %s\n", sn.Name)
	}
	for _, s := range sn.AttributionTexts {
		fmt.Fprintf(w, "SnippetAttributionText: %s\n", textify(s))
	}

	fmt.Fprintf(w, "\n")

	return nil
}
