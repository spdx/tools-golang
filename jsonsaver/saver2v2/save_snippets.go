// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"github.com/spdx/tools-golang/spdx"
)

func renderSnippets2_2(doc *spdx.Document2_2, jsondocument map[string]interface{}) error {

	var snippets []interface{}
	for _, value := range doc.UnpackagedFiles {
		snippet := make(map[string]interface{})
		for _, v := range value.Snippets {
			snippet["SPDXID"] = spdx.RenderElementID(v.SnippetSPDXIdentifier)
			if v.SnippetComment != "" {
				snippet["comment"] = v.SnippetComment
			}
			if v.SnippetCopyrightText != "" {
				snippet["copyrightText"] = v.SnippetCopyrightText
			}
			if v.SnippetLicenseComments != "" {
				snippet["licenseComments"] = v.SnippetLicenseComments
			}
			if v.SnippetLicenseConcluded != "" {
				snippet["licenseConcluded"] = v.SnippetLicenseConcluded
			}
			if v.LicenseInfoInSnippet != nil {
				snippet["licenseInfoInSnippets"] = v.LicenseInfoInSnippet
			}
			if v.SnippetName != "" {
				snippet["name"] = v.SnippetName
			}
			if v.SnippetName != "" {
				snippet["snippetFromFile"] = spdx.RenderDocElementID(v.SnippetFromFileSPDXIdentifier)
			}
			if v.SnippetAttributionTexts != nil {
				snippet["attributionTexts"] = v.SnippetAttributionTexts
			}

			// save  snippet ranges
			var ranges []interface{}
			byterange := map[string]interface{}{
				"endPointer": map[string]interface{}{
					"offset":    v.SnippetByteRangeEnd,
					"reference": spdx.RenderDocElementID(v.SnippetFromFileSPDXIdentifier),
				},
				"startPointer": map[string]interface{}{
					"offset":    v.SnippetByteRangeStart,
					"reference": spdx.RenderDocElementID(v.SnippetFromFileSPDXIdentifier),
				},
			}
			linerange := map[string]interface{}{
				"endPointer": map[string]interface{}{
					"lineNumber": v.SnippetLineRangeEnd,
					"reference":  spdx.RenderDocElementID(v.SnippetFromFileSPDXIdentifier),
				},
				"startPointer": map[string]interface{}{
					"lineNumber": v.SnippetLineRangeStart,
					"reference":  spdx.RenderDocElementID(v.SnippetFromFileSPDXIdentifier),
				},
			}
			ranges = append(ranges, byterange, linerange)
			snippet["ranges"] = ranges
			snippets = append(snippets, snippet)
		}
	}
	jsondocument["snippets"] = snippets
	return nil
}
