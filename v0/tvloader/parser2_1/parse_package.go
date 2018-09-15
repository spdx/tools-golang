// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"fmt"
	"strings"

	"github.com/swinslow/spdx-go/v0/spdx"
)

func (parser *tvParser2_1) parsePairFromPackage2_1(tag string, value string) error {
	switch tag {
	case "PackageName":
		// if package already has a name, create and go on to a new package
		if parser.pkg.PackageName != "" {
			parser.pkg = &spdx.Package2_1{
				IsUnpackaged:              false,
				FilesAnalyzed:             true,
				IsFilesAnalyzedTagPresent: false,
			}
			parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
		}
		parser.pkg.PackageName = value
	// tag for going on to file section
	case "FileName":
		parser.st = psFile2_1
		return parser.parsePairFromFile2_1(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_1
		return parser.parsePairFromOtherLicense2_1(tag, value)
	case "SPDXID":
		parser.pkg.PackageSPDXIdentifier = value
	case "PackageVersion":
		parser.pkg.PackageVersion = value
	case "PackageFileName":
		parser.pkg.PackageFileName = value
	case "PackageSupplier":
		if value == "NOASSERTION" {
			parser.pkg.PackageSupplierNOASSERTION = true
			break
		}
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "Person":
			parser.pkg.PackageSupplierPerson = subvalue
		case "Organization":
			parser.pkg.PackageSupplierOrganization = subvalue
		default:
			return fmt.Errorf("unrecognized PackageSupplier type %v", subkey)
		}
	case "PackageOriginator":
		if value == "NOASSERTION" {
			parser.pkg.PackageOriginatorNOASSERTION = true
			break
		}
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "Person":
			parser.pkg.PackageOriginatorPerson = subvalue
		case "Organization":
			parser.pkg.PackageOriginatorOrganization = subvalue
		default:
			return fmt.Errorf("unrecognized PackageSupplier type %v", subkey)
		}
	case "PackageDownloadLocation":
		parser.pkg.PackageDownloadLocation = value
	case "FilesAnalyzed":
		parser.pkg.IsFilesAnalyzedTagPresent = true
		if value == "false" {
			parser.pkg.FilesAnalyzed = false
		} else if value == "true" {
			parser.pkg.FilesAnalyzed = true
		}
		parser.pkg.PackageDownloadLocation = value
	case "PackageVerificationCode":
		code, excludesFileName := extractCodeAndExcludes(value)
		parser.pkg.PackageVerificationCode = code
		parser.pkg.PackageVerificationCodeExcludedFile = excludesFileName
	}

	return nil
}

// ===== Helper functions =====

func extractCodeAndExcludes(value string) (string, string) {
	// FIXME this should probably be done using regular expressions instead
	// split by paren + word "excludes:"
	sp := strings.SplitN(value, "(excludes:", 2)
	if len(sp) < 2 {
		// not found; return the whole string as just the code
		return value, ""
	}

	// if we're here, code is in first part and excludes filename is in
	// second part, with trailing paren
	code := strings.TrimSpace(sp[0])
	parsedSp := strings.SplitN(sp[1], ")", 2)
	fileName := strings.TrimSpace(parsedSp[0])
	return code, fileName
}
