// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func renderfiles2_2(doc *spdx.Document2_2, buf *bytes.Buffer) error {

	var files []interface{}
	for k, v := range doc.UnpackagedFiles {
		file := make(map[string]interface{})
		file["SPDXID"] = spdx.RenderElementID(k)
		ann, _ := renderAnnotations2_2(doc.Annotations, spdx.MakeDocElementID("", string(v.FileSPDXIdentifier)))
		if ann != nil {
			file["annotations"] = ann
		}
		if v.FileContributor != nil {
			file["attributionTexts"] = v.FileContributor
		}
		if v.FileComment != "" {
			file["comment"] = v.FileComment
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
			file["checksums"] = checksums
		}
		if v.FileCopyrightText != "" {
			file["copyrightText"] = v.FileCopyrightText
		}
		if v.FileName != "" {
			file["fileName"] = v.FileName
		}
		if v.FileType != nil {
			file["fileTypes"] = v.FileType
		}
		if v.LicenseComments != "" {
			file["licenseComments"] = v.LicenseComments
		}
		if v.LicenseConcluded != "" {
			file["licenseConcluded"] = v.LicenseConcluded
		}
		if v.LicenseInfoInFile != nil {
			file["licenseInfoFromFiles"] = v.LicenseInfoInFile
		}
		if v.FileNotice != "" {
			file["noticeText"] = v.FileNotice
		}
		if v.FileContributor != nil {
			file["fileContributors"] = v.FileContributor
		}
		if v.FileDependencies != nil {
			file["fileDependencies"] = v.FileDependencies
		}
		if v.FileAttributionTexts != nil {
			file["attributionTexts"] = v.FileAttributionTexts
		}

		files = append(files, file)
	}
	filesjson, _ := json.Marshal(files)
	fmt.Fprintf(buf, "\"%s\": %s ,", "files", filesjson)

	return nil
}
