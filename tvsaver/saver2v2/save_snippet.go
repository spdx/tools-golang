// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

func renderSnippet2_2(sn *spdx.Snippet2_2, w io.Writer) error {
	if sn.SnippetSPDXIdentifier != "" {
		fmt.Fprintf(w, "SnippetSPDXID: %s\n", spdx.RenderElementID(sn.SnippetSPDXIdentifier))
	}
	snFromFileIDStr := spdx.RenderElementID(sn.SnippetFromFileSPDXIdentifier)
	if snFromFileIDStr != "" {
		fmt.Fprintf(w, "SnippetFromFileSPDXID: %s\n", snFromFileIDStr)
	}

	for _, snippetRange := range sn.Ranges {
		if snippetRange.StartPointer.Offset != 0 && snippetRange.EndPointer.Offset != 0 {
			fmt.Fprintf(w, "SnippetByteRange: %d:%d\n", snippetRange.StartPointer.Offset, snippetRange.EndPointer.Offset)
		}
		if snippetRange.StartPointer.LineNumber != 0 && snippetRange.EndPointer.LineNumber != 0 {
			fmt.Fprintf(w, "SnippetLineRange: %d:%d\n", snippetRange.StartPointer.LineNumber, snippetRange.EndPointer.LineNumber)
		}
	}
	if sn.SnippetLicenseConcluded != "" {
		fmt.Fprintf(w, "SnippetLicenseConcluded: %s\n", sn.SnippetLicenseConcluded)
	}
	for _, s := range sn.LicenseInfoInSnippet {
		fmt.Fprintf(w, "LicenseInfoInSnippet: %s\n", s)
	}
	if sn.SnippetLicenseComments != "" {
		fmt.Fprintf(w, "SnippetLicenseComments: %s\n", textify(sn.SnippetLicenseComments))
	}
	if sn.SnippetCopyrightText != "" {
		fmt.Fprintf(w, "SnippetCopyrightText: %s\n", textify(sn.SnippetCopyrightText))
	}
	if sn.SnippetComment != "" {
		fmt.Fprintf(w, "SnippetComment: %s\n", textify(sn.SnippetComment))
	}
	if sn.SnippetName != "" {
		fmt.Fprintf(w, "SnippetName: %s\n", sn.SnippetName)
	}
	for _, s := range sn.SnippetAttributionTexts {
		fmt.Fprintf(w, "SnippetAttributionText: %s\n", textify(s))
	}

	fmt.Fprintf(w, "\n")

	return nil
}
