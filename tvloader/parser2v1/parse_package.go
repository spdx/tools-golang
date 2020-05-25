// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"
	"strings"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_1) parsePairFromPackage2_1(tag string, value string) error {
	// expire pkgExtRef for anything other than a comment
	// (we'll actually handle the comment further below)
	if tag != "ExternalRefComment" {
		parser.pkgExtRef = nil
	}

	switch tag {
	case "PackageName":
		// if package already has a name, create and go on to a new package
		if parser.pkg == nil || parser.pkg.PackageName != "" {
			parser.pkg = &spdx.Package2_1{
				FilesAnalyzed:             true,
				IsFilesAnalyzedTagPresent: false,
			}
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
		eID, err := extractElementID(value)
		if err != nil {
			return err
		}
		parser.pkg.PackageSPDXIdentifier = eID
		if parser.doc.Packages == nil {
			parser.doc.Packages = map[spdx.ElementID]*spdx.Package2_1{}
		}
		parser.doc.Packages[eID] = parser.pkg
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
			return fmt.Errorf("unrecognized PackageOriginator type %v", subkey)
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
	case "PackageVerificationCode":
		code, excludesFileName := extractCodeAndExcludes(value)
		parser.pkg.PackageVerificationCode = code
		parser.pkg.PackageVerificationCodeExcludedFile = excludesFileName
	case "PackageChecksum":
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "SHA1":
			parser.pkg.PackageChecksumSHA1 = subvalue
		case "SHA256":
			parser.pkg.PackageChecksumSHA256 = subvalue
		case "MD5":
			parser.pkg.PackageChecksumMD5 = subvalue
		default:
			return fmt.Errorf("got unknown checksum type %s", subkey)
		}
	case "PackageHomePage":
		parser.pkg.PackageHomePage = value
	case "PackageSourceInfo":
		parser.pkg.PackageSourceInfo = value
	case "PackageLicenseConcluded":
		parser.pkg.PackageLicenseConcluded = value
	case "PackageLicenseInfoFromFiles":
		parser.pkg.PackageLicenseInfoFromFiles = append(parser.pkg.PackageLicenseInfoFromFiles, value)
	case "PackageLicenseDeclared":
		parser.pkg.PackageLicenseDeclared = value
	case "PackageLicenseComments":
		parser.pkg.PackageLicenseComments = value
	case "PackageCopyrightText":
		parser.pkg.PackageCopyrightText = value
	case "PackageSummary":
		parser.pkg.PackageSummary = value
	case "PackageDescription":
		parser.pkg.PackageDescription = value
	case "PackageComment":
		parser.pkg.PackageComment = value
	case "ExternalRef":
		parser.pkgExtRef = &spdx.PackageExternalReference2_1{}
		parser.pkg.PackageExternalReferences = append(parser.pkg.PackageExternalReferences, parser.pkgExtRef)
		category, refType, locator, err := extractPackageExternalReference(value)
		if err != nil {
			return err
		}
		parser.pkgExtRef.Category = category
		parser.pkgExtRef.RefType = refType
		parser.pkgExtRef.Locator = locator
	case "ExternalRefComment":
		if parser.pkgExtRef == nil {
			return fmt.Errorf("no current ExternalRef found")
		}
		parser.pkgExtRef.ExternalRefComment = value
		// now, expire pkgExtRef anyway because it can have at most one comment
		parser.pkgExtRef = nil
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_1{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_1(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_1(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_1{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_1(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_1
		return parser.parsePairFromReview2_1(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in Package section", tag)
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

func extractPackageExternalReference(value string) (string, string, string, error) {
	sp := strings.Split(value, " ")
	// remove any that are just whitespace
	keepSp := []string{}
	for _, s := range sp {
		ss := strings.TrimSpace(s)
		if ss != "" {
			keepSp = append(keepSp, ss)
		}
	}
	// now, should have 3 items and should be able to map them
	if len(keepSp) != 3 {
		return "", "", "", fmt.Errorf("expected 3 elements, got %d", len(keepSp))
	}
	return keepSp[0], keepSp[1], keepSp[2], nil
}
