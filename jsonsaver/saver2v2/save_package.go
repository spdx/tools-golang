// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func renderPackage2_2(doc *spdx.Document2_2, buf *bytes.Buffer) error {

	var packages []interface{}
	for k, v := range doc.Packages {
		pkg := make(map[string]interface{})
		pkg["SPDXID"] = spdx.RenderElementID(k)
		ann, _ := renderAnnotations2_2(doc.Annotations, spdx.MakeDocElementID("", string(v.PackageSPDXIdentifier)))
		if ann != nil {
			pkg["annotations"] = ann
		}
		if v.PackageAttributionTexts != nil {
			pkg["attributionTexts"] = v.PackageAttributionTexts
		}
		// parse package checksums
		if v.PackageChecksums != nil {
			var checksums []interface{}
			for _, value := range v.PackageChecksums {
				checksum := make(map[string]interface{})
				checksum["algorithm"] = value.Algorithm
				checksum["checksumValue"] = value.Value
				checksums = append(checksums, checksum)
			}
			pkg["checksums"] = checksums
		}
		if v.PackageCopyrightText != "" {
			pkg["copyrightText"] = v.PackageCopyrightText
		}
		if v.PackageDescription != "" {
			pkg["description"] = v.PackageDescription
		}
		if v.PackageDownloadLocation != "" {
			pkg["downloadLocation"] = v.PackageDownloadLocation
		}

		//parse document external refereneces
		if v.PackageExternalReferences != nil {
			var externalrefs []interface{}
			for _, value := range v.PackageExternalReferences {
				ref := make(map[string]interface{})
				ref["referenceCategory"] = value.Category
				ref["referenceLocator"] = value.Locator
				ref["referenceType"] = value.RefType
				if value.ExternalRefComment != "" {
					ref["comment"] = value.ExternalRefComment
				}
				externalrefs = append(externalrefs, ref)
			}
			pkg["externalRefs"] = externalrefs
		}

		pkg["filesAnalyzed"] = v.FilesAnalyzed

		// parse package hass files
		if v.Files != nil {
			var fileIds []string
			for k, v := range v.Files {
				doc.UnpackagedFiles[k] = v
				fileIds = append(fileIds, spdx.RenderElementID(k))
			}
			pkg["hasFiles"] = fileIds
		}

		if v.PackageHomePage != "" {
			pkg["homepage"] = v.PackageHomePage
		}
		if v.PackageLicenseComments != "" {
			pkg["licenseComments"] = v.PackageLicenseComments
		}
		if v.PackageLicenseConcluded != "" {
			pkg["licenseConcluded"] = v.PackageLicenseConcluded
		}
		if v.PackageLicenseDeclared != "" {
			pkg["licenseDeclared"] = v.PackageLicenseDeclared
		}
		if v.PackageLicenseInfoFromFiles != nil {
			pkg["licenseInfoFromFiles"] = v.PackageLicenseInfoFromFiles
		}
		if v.PackageName != "" {
			pkg["name"] = v.PackageName
		}
		if v.PackageSourceInfo != "" {
			pkg["sourceInfo"] = v.PackageSourceInfo
		}
		if v.PackageSummary != "" {
			pkg["summary"] = v.PackageSummary
		}
		if v.PackageVersion != "" {
			pkg["versionInfo"] = v.PackageVersion
		}
		if v.PackageFileName != "" {
			pkg["packageFileName"] = v.PackageFileName
		}

		//parse package originator
		if v.PackageOriginatorPerson != "" {
			pkg["originator"] = fmt.Sprintf("Person: %s", v.PackageOriginatorPerson)
		}
		if v.PackageOriginatorOrganization != "" {
			pkg["originator"] = fmt.Sprintf("Organization: %s", v.PackageOriginatorOrganization)
		}
		if v.PackageOriginatorNOASSERTION {
			pkg["originator"] = "NOASSERTION"
		}

		//parse package verification code
		if v.PackageVerificationCode != "" {
			verification := make(map[string]interface{})
			verification["packageVerificationCodeExcludedFiles"] = []string{v.PackageVerificationCodeExcludedFile}
			verification["packageVerificationCodeValue"] = v.PackageVerificationCode
			pkg["packageVerificationCode"] = verification
		}

		//parse package supplier
		if v.PackageSupplierPerson != "" {
			pkg["supplier"] = fmt.Sprintf("Person: %s", v.PackageSupplierPerson)
		}
		if v.PackageSupplierOrganization != "" {
			pkg["supplier"] = fmt.Sprintf("Organization: %s", v.PackageSupplierOrganization)
		}
		if v.PackageSupplierNOASSERTION {
			pkg["supplier"] = "NOASSERTION"
		}

		packages = append(packages, pkg)
	}
	packagejson, _ := json.Marshal(packages)
	fmt.Fprintf(buf, "\"%s\": %s ,", "packages", packagejson)

	return nil
}
