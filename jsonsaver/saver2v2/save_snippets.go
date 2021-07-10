// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func rendersnippets2_2(doc *spdx.Document2_2, buf *bytes.Buffer) error {

	var snippets []interface{}
	for k, v := range doc.UnpackagedFiles {
		snippet := make(map[string]interface{})
		snippet["SPDXID"] = spdx.RenderElementID(k)
		ann, _ := renderAnnotations2_2(doc.Annotations, spdx.MakeDocElementID("", string(v.FileSPDXIdentifier)))
		if ann != nil {
			snippet["annotations"] = ann
		}
		if v.FileContributor != nil {
			snippet["attributionTexts"] = v.FileContributor
		}
		if v.FileComment != "" {
			snippet["comment"] = v.FileComment
		}

		// parse package checksums
		if v.FileChecksums != nil {
			var checksums []interface{}
			for _, value := range v.FileChecksums {
				checksum := make(map[string]interface{})
				checksum["algorithm"] = value.Algorithm
				checksum["checksumValue"] = value.Value
				checksums = append(checksums, checksum)
			}
			snippet["checksums"] = checksums
		}
		if v.FileCopyrightText != "" {
			snippet["copyrightText"] = v.FileCopyrightText
		}
		if v.FileName != "" {
			snippet["fileName"] = v.FileName
		}
		if v.FileType != nil {
			snippet["fileTypes"] = v.FileType
		}
		if v.LicenseComments != "" {
			snippet["licenseComments"] = v.LicenseComments
		}
		if v.LicenseConcluded != "" {
			snippet["licenseConcluded"] = v.LicenseConcluded
		}
		if v.LicenseInfoInFile != nil {
			snippet["licenseInfoFromFiles"] = v.LicenseInfoInFile
		}
		if v.FileNotice != "" {
			snippet["name"] = v.FileNotice
		}
		if v.FileContributor != nil {
			snippet["fileContributors"] = v.FileContributor
		}
		if v.FileDependencies != nil {
			snippet["fileDependencies"] = v.FileDependencies
		}
		if v.FileAttributionTexts != nil {
			snippet["attributionTexts"] = v.FileAttributionTexts
		}

		snippets = append(snippets, snippet)
	}
	snippetjson, _ := json.Marshal(snippets)
	fmt.Fprintf(buf, "\"%s\": %s ,", "snippets", snippetjson)

	return nil
}
