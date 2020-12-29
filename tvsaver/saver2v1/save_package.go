// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"fmt"
	"io"
	"sort"

	"github.com/spdx/tools-golang/spdx"
)

func renderPackage2_1(pkg *spdx.Package2_1, w io.Writer) error {
	if pkg.Name != "" {
		fmt.Fprintf(w, "PackageName: %s\n", pkg.Name)
	}
	if pkg.SPDXIdentifier != "" {
		fmt.Fprintf(w, "SPDXID: %s\n", spdx.RenderElementID(pkg.SPDXIdentifier))
	}
	if pkg.Version != "" {
		fmt.Fprintf(w, "PackageVersion: %s\n", pkg.Version)
	}
	if pkg.FileName != "" {
		fmt.Fprintf(w, "PackageFileName: %s\n", pkg.FileName)
	}
	if pkg.SupplierPerson != "" {
		fmt.Fprintf(w, "PackageSupplier: Person: %s\n", pkg.SupplierPerson)
	}
	if pkg.SupplierOrganization != "" {
		fmt.Fprintf(w, "PackageSupplier: Organization: %s\n", pkg.SupplierOrganization)
	}
	if pkg.SupplierNOASSERTION == true {
		fmt.Fprintf(w, "PackageSupplier: NOASSERTION\n")
	}
	if pkg.OriginatorPerson != "" {
		fmt.Fprintf(w, "PackageOriginator: Person: %s\n", pkg.OriginatorPerson)
	}
	if pkg.OriginatorOrganization != "" {
		fmt.Fprintf(w, "PackageOriginator: Organization: %s\n", pkg.OriginatorOrganization)
	}
	if pkg.OriginatorNOASSERTION == true {
		fmt.Fprintf(w, "PackageOriginator: NOASSERTION\n")
	}
	if pkg.DownloadLocation != "" {
		fmt.Fprintf(w, "PackageDownloadLocation: %s\n", pkg.DownloadLocation)
	}
	if pkg.FilesAnalyzed == true {
		if pkg.IsFilesAnalyzedTagPresent == true {
			fmt.Fprintf(w, "FilesAnalyzed: true\n")
		}
	} else {
		fmt.Fprintf(w, "FilesAnalyzed: false\n")
	}
	if pkg.VerificationCode != "" && pkg.FilesAnalyzed == true {
		if pkg.VerificationCodeExcludedFile == "" {
			fmt.Fprintf(w, "PackageVerificationCode: %s\n", pkg.VerificationCode)
		} else {
			fmt.Fprintf(w, "PackageVerificationCode: %s (excludes %s)\n", pkg.VerificationCode, pkg.VerificationCodeExcludedFile)
		}
	}
	if pkg.ChecksumSHA1 != "" {
		fmt.Fprintf(w, "PackageChecksum: SHA1: %s\n", pkg.ChecksumSHA1)
	}
	if pkg.ChecksumSHA256 != "" {
		fmt.Fprintf(w, "PackageChecksum: SHA256: %s\n", pkg.ChecksumSHA256)
	}
	if pkg.ChecksumMD5 != "" {
		fmt.Fprintf(w, "PackageChecksum: MD5: %s\n", pkg.ChecksumMD5)
	}
	if pkg.HomePage != "" {
		fmt.Fprintf(w, "PackageHomePage: %s\n", pkg.HomePage)
	}
	if pkg.SourceInfo != "" {
		fmt.Fprintf(w, "PackageSourceInfo: %s\n", textify(pkg.SourceInfo))
	}
	if pkg.LicenseConcluded != "" {
		fmt.Fprintf(w, "PackageLicenseConcluded: %s\n", pkg.LicenseConcluded)
	}
	if pkg.FilesAnalyzed == true {
		for _, s := range pkg.LicenseInfoFromFiles {
			fmt.Fprintf(w, "PackageLicenseInfoFromFiles: %s\n", s)
		}
	}
	if pkg.LicenseDeclared != "" {
		fmt.Fprintf(w, "PackageLicenseDeclared: %s\n", pkg.LicenseDeclared)
	}
	if pkg.LicenseComments != "" {
		fmt.Fprintf(w, "PackageLicenseComments: %s\n", textify(pkg.LicenseComments))
	}
	if pkg.CopyrightText != "" {
		fmt.Fprintf(w, "PackageCopyrightText: %s\n", pkg.CopyrightText)
	}
	if pkg.Summary != "" {
		fmt.Fprintf(w, "PackageSummary: %s\n", textify(pkg.Summary))
	}
	if pkg.Description != "" {
		fmt.Fprintf(w, "PackageDescription: %s\n", textify(pkg.Description))
	}
	if pkg.Comment != "" {
		fmt.Fprintf(w, "PackageComment: %s\n", textify(pkg.Comment))
	}
	for _, s := range pkg.ExternalReferences {
		fmt.Fprintf(w, "ExternalRef: %s %s %s\n", s.Category, s.RefType, s.Locator)
		if s.ExternalRefComment != "" {
			fmt.Fprintf(w, "ExternalRefComment: %s\n", s.ExternalRefComment)
		}
	}

	fmt.Fprintf(w, "\n")

	// also render any files for this package
	// get slice of File identifiers so we can sort them
	fileKeys := []string{}
	for k := range pkg.Files {
		fileKeys = append(fileKeys, string(k))
	}
	sort.Strings(fileKeys)
	for _, fiID := range fileKeys {
		fi := pkg.Files[spdx.ElementID(fiID)]
		renderFile2_1(fi, w)
	}

	return nil
}
