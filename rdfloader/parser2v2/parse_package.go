// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
	"strings"
)

func (parser *rdfParser2_2) getPackageFromNode(packageNode *gordfParser.Node) (pkg *spdx.Package2_2, err error) {
	pkg = &spdx.Package2_2{} // new package which will be returned

	// setting the SPDXIdentifier for the package.
	eId, err := ExtractElementID(getLastPartOfURI(packageNode.ID))
	if err != nil {
		return nil, fmt.Errorf("error extracting elementID of a package identifier: %v", err)
	}
	pkg.PackageSPDXIdentifier = eId // 3.2

	// iterate over all the triples associated with the provided package packageNode.
	for _, subTriple := range parser.nodeToTriples(packageNode) {
		switch subTriple.Predicate.ID {
		case RDF_TYPE:
			// cardinality: exactly 1
			continue
		case SPDX_NAME: // 3.1
			// cardinality: exactly 1
			pkg.PackageName = subTriple.Object.ID
		case SPDX_VERSION_INFO: // 3.3
			// cardinality: max 1
			pkg.PackageVersion = subTriple.Object.ID
		case SPDX_PACKAGE_FILE_NAME: // 3.4
			// cardinality: max 1
			pkg.PackageFileName = subTriple.Object.ID
		case SPDX_SUPPLIER: // 3.5
			// cardinality: max 1
			err = setPackageSupplier(pkg, subTriple.Object.ID)
			if err != nil {
				return nil, err
			}
		case SPDX_ORIGINATOR: // 3.6
			// cardinality: max 1
			err = setPackageOriginator(pkg, subTriple.Object.ID)
			if err != nil {
				return nil, err
			}
		case SPDX_DOWNLOAD_LOCATION: // 3.7
			// cardinality: exactly 1
			err = setDocumentLocationFromURI(subTriple.Object.ID, pkg)
			if err != nil {
				return nil, err
			}
		case SPDX_FILES_ANALYZED: // 3.8
			// cardinality: max 1
			err = setFilesAnalyzed(subTriple.Object.ID, pkg)
			if err != nil {
				err = fmt.Errorf("error setting a filesAnalyzed attribute of a package: %v", err)
			}
		case SPDX_PACKAGE_VERIFICATION_CODE: // 3.9
			// cardinality: max 1
			err = parser.setPackageVerificationCode(pkg, subTriple.Object)
		case SPDX_CHECKSUM: // 3.10
			// cardinality: min 0
			err = setPackageChecksum(parser, pkg, subTriple.Object)
		case DOAP_HOMEPAGE: // 3.11
			// cardinality: max 1
			// homepage must be a valid Uri
			if !isUriValid(subTriple.Object.ID) {
				return nil, fmt.Errorf("invalid uri %s while parsing doap_homepage in a package", subTriple.Object.ID)
			}
			pkg.PackageHomePage = subTriple.Object.ID
		case SPDX_SOURCE_INFO: // 3.12
			// cardinality: max 1
			pkg.PackageSourceInfo = subTriple.Object.ID
		case SPDX_LICENSE_CONCLUDED: // 3.13
			// cardinality: exactly 1
			anyLicenseInfo, err := parser.getAnyLicenseFromNode(subTriple.Object)
			if err != nil {
				return nil, err
			}
			pkg.PackageLicenseConcluded = anyLicenseInfo.ToLicenseString()
		case SPDX_LICENSE_INFO_FROM_FILES: // 3.14
			// cardinality: min 0
			pkg.PackageLicenseInfoFromFiles = append(pkg.PackageLicenseInfoFromFiles, getLicenseStringFromURI(subTriple.Object.ID))
		case SPDX_LICENSE_DECLARED: // 3.15
			// cardinality: exactly 1
			anyLicenseInfo, err := parser.getAnyLicenseFromNode(subTriple.Object)
			if err != nil {
				return nil, err
			}
			pkg.PackageLicenseDeclared = anyLicenseInfo.ToLicenseString()
		case SPDX_LICENSE_COMMENTS: // 3.16
			// cardinality: max 1
			pkg.PackageLicenseComments = subTriple.Object.ID
		case SPDX_COPYRIGHT_TEXT: // 3.17
			// cardinality: exactly 1
			pkg.PackageCopyrightText = subTriple.Object.ID
		case SPDX_SUMMARY: // 3.18
			// cardinality: max 1
			pkg.PackageSummary = subTriple.Object.ID
		case SPDX_DESCRIPTION: // 3.19
			// cardinality: max 1
			pkg.PackageDescription = subTriple.Object.ID
		case RDFS_COMMENT: // 3.20
			// cardinality: max 1
			pkg.PackageComment = subTriple.Object.ID
		case SPDX_EXTERNAL_REF: // 3.21
			// cardinality: min 0
			externalDocRef, err := parser.getPackageExternalRef(subTriple)
			if err != nil {
				return nil, err
			}
			pkg.PackageExternalReferences = append(pkg.PackageExternalReferences, externalDocRef)
		case SPDX_HAS_FILE: // 3.22
			// cardinality: min 0
			file, err := parser.getFileFromNode(subTriple.Object)
			if err != nil {
				return nil, err
			}
			parser.setFileToPackage(pkg, file)
		case SPDX_RELATIONSHIP:
			// cardinality: min 0
			err := parser.parseRelationship(subTriple)
			if err != nil {
				return nil, err
			}
		case SPDX_ATTRIBUTION_TEXT:
			// cardinality: min 0
			pkg.PackageAttributionTexts = append(pkg.PackageAttributionTexts, subTriple.Object.ID)
		case SPDX_ANNOTATION:
			// cardinality: min 0
			err := parser.parseAnnotationFromNode(subTriple.Object)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown predicate id %s while parsing a package", subTriple.Predicate.ID)
		}
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

// parses externalReference found in the package by the associated triple.
func (parser *rdfParser2_2) getPackageExternalRef(triple *gordfParser.Triple) (externalDocRef *spdx.PackageExternalReference2_2, err error) {
	externalDocRef = &spdx.PackageExternalReference2_2{}
	for _, subTriple := range parser.nodeToTriples(triple.Object) {
		switch subTriple.Predicate.ID {
		case SPDX_REFERENCE_CATEGORY:
			// cardinality: exactly 1
			switch subTriple.Object.ID {
			case SPDX_REFERENCE_CATEGORY_SECURITY:
				externalDocRef.Category = "SECURITY"
			case SPDX_REFERENCE_CATEGORY_PACKAGE_MANAGER:
				externalDocRef.Category = "PACKAGE-MANAGER"
			case SPDX_REFERENCE_CATEGORY_OTHER:
				externalDocRef.Category = "OTHER"
			default:
				return nil, fmt.Errorf("unknown packageManager predicate uri %s", subTriple.Predicate.ID)
			}
		case RDF_TYPE:
			continue
		case SPDX_REFERENCE_TYPE:
			// assumes: the reference type is associated with just the uri and
			// 			other associated fields are ignored.
			// other fields include:
			//		1. contextualExample,
			//		2. documentation and,
			//		3. externalReferenceSite
			externalDocRef.RefType = subTriple.Object.ID
		case SPDX_REFERENCE_LOCATOR:
			// cardinality: exactly 1
			externalDocRef.Locator = subTriple.Object.ID
		case RDFS_COMMENT:
			// cardinality: max 1
			externalDocRef.ExternalRefComment = subTriple.Object.ID
		default:
			return nil, fmt.Errorf("unknown package external reference predicate id %s", subTriple.Predicate.ID)
		}
	}
	return
}

func (parser *rdfParser2_2) setPackageVerificationCode(pkg *spdx.Package2_2, node *gordfParser.Node) error {
	for _, subTriple := range parser.nodeToTriples(node) {
		switch subTriple.Predicate.ID {
		case SPDX_PACKAGE_VERIFICATION_CODE_VALUE:
			// cardinality: exactly 1
			pkg.PackageVerificationCode = subTriple.Object.ID
		case SPDX_PACKAGE_VERIFICATION_CODE_EXCLUDED_FILE:
			// cardinality: min 0
			pkg.PackageVerificationCodeExcludedFile = subTriple.Object.ID
		case RDF_TYPE:
			// cardinality: exactly 1
			continue
		default:
			return fmt.Errorf("unparsed predicate %s", subTriple.Predicate.ID)
		}
	}
	return nil
}

func (parser *rdfParser2_2) setFileToPackage(pkg *spdx.Package2_2, file *spdx.File2_2) {
	if pkg.Files == nil {
		pkg.Files = map[spdx.ElementID]*spdx.File2_2{}
	}
	pkg.Files[file.FileSPDXIdentifier] = file
	parser.assocWithPackage[file.FileSPDXIdentifier] = true
}

// given a supplierObject, sets the PackageSupplier attribute of the pkg.
// Args:
//    value: [NOASSERTION | [Person | Organization]: string]
func setPackageSupplier(pkg *spdx.Package2_2, value string) error {
	value = strings.TrimSpace(value)
	if strings.ToUpper(value) == "NOASSERTION" {
		pkg.PackageSupplierNOASSERTION = true
		return nil
	}
	subKey, subValue, err := ExtractSubs(value, ":")
	if err != nil {
		return fmt.Errorf("package supplier must be of the form NOASSERTION or [Person|Organization]: string. found: %s", value)
	}
	switch subKey {
	case "Person":
		pkg.PackageSupplierPerson = subValue
	case "Organization":
		pkg.PackageSupplierOrganization = subValue
	default:
		return fmt.Errorf("unknown supplier %s", subKey)
	}
	return nil
}

// given a OriginatorObject, sets the PackageOriginator attribute of the pkg.
// Args:
//    value: [NOASSERTION | [Person | Organization]: string]
func setPackageOriginator(pkg *spdx.Package2_2, value string) error {
	value = strings.TrimSpace(value)
	if strings.ToUpper(value) == "NOASSERTION" {
		pkg.PackageOriginatorNOASSERTION = true
		return nil
	}
	subKey, subValue, err := ExtractSubs(value, ":")
	if err != nil {
		return fmt.Errorf("package originator must be of the form NOASSERTION or [Person|Organization]: string. found: %s", value)
	}

	switch subKey {
	case "Person":
		pkg.PackageOriginatorPerson = subValue
	case "Organization":
		pkg.PackageOriginatorOrganization = subValue
	default:
		return fmt.Errorf("originator can be either a Person or Organization. found %s", subKey)
	}
	return nil
}

// validates the uri and sets the location if it is valid
func setDocumentLocationFromURI(locationURI string, pkg *spdx.Package2_2) error {
	switch locationURI {
	case SPDX_NOASSERTION_CAPS, SPDX_NOASSERTION_CAPS:
		pkg.PackageDownloadLocation = "NOASSERTION"
	case SPDX_NONE_CAPS, SPDX_NONE_SMALL:
		pkg.PackageDownloadLocation = "NONE"
	default:
		if !isUriValid(locationURI) {
			return fmt.Errorf("%s is not a valid uri", locationURI)
		}
		pkg.PackageDownloadLocation = locationURI
	}
	return nil
}

// sets the FilesAnalyzed attribute to the given package
// boolValue is a string of type "true" or "false"
func setFilesAnalyzed(boolValue string, pkg *spdx.Package2_2) error {
	pkg.IsFilesAnalyzedTagPresent = true
	switch strings.TrimSpace(boolValue) {
	case "true":
		pkg.FilesAnalyzed = true
	case "false":
		pkg.FilesAnalyzed = false
	default:
		return fmt.Errorf("filesAnalyzed can be either true/false. found %s", boolValue)
	}
	return nil
}

func setPackageChecksum(parser *rdfParser2_2, pkg *spdx.Package2_2, node *gordfParser.Node) error {
	checksumAlgorithm, checksumValue, err := parser.getChecksumFromNode(node)
	if err != nil {
		return err
	}
	switch checksumAlgorithm {
	case "MD5":
		pkg.PackageChecksumMD5 = checksumValue
	case "SHA1":
		pkg.PackageChecksumSHA1 = checksumValue
	case "SHA256":
		pkg.PackageChecksumSHA256 = checksumValue
	default:
		return fmt.Errorf("unknown checksumAlgorithm %s while parsing a package", checksumAlgorithm)
	}
	return nil
}
