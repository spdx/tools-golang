// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"
	"strings"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func (parser *rdfParser2_2) getPackageFromNode(packageNode *gordfParser.Node) (pkg *v2_2.Package, err error) {
	pkg = &v2_2.Package{} // new package which will be returned

	currState := parser.cache[packageNode.ID]
	if currState == nil {
		// there is no entry about the state of current package node.
		// this is the first time we're seeing this node.
		parser.cache[packageNode.ID] = &nodeState{
			object: pkg,
			Color:  WHITE,
		}
	} else if currState.Color == GREY {
		// we have already started parsing this package node and we needn't parse it again.
		return currState.object.(*v2_2.Package), nil
	}

	// setting color of the state to grey to indicate that we've started to
	// parse this node once.
	parser.cache[packageNode.ID].Color = GREY

	// setting state color to black to indicate when we're done parsing this node.
	defer func() { parser.cache[packageNode.ID].Color = BLACK }()

	// setting the SPDXIdentifier for the package.
	eId, err := ExtractElementID(getLastPartOfURI(packageNode.ID))
	if err != nil {
		return nil, fmt.Errorf("error extracting elementID of a package identifier: %v", err)
	}
	pkg.PackageSPDXIdentifier = eId // 3.2

	// check if we already have a package initialized for this ID
	existingPackageIndex := -1
	for ii, existingPkg := range parser.doc.Packages {
		if existingPkg != nil && existingPkg.PackageSPDXIdentifier == eId {
			existingPackageIndex = ii
			pkg = existingPkg
			break
		}
	}

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
		case SPDX_ORIGINATOR: // 3.6
			// cardinality: max 1
			err = setPackageOriginator(pkg, subTriple.Object.ID)
		case SPDX_DOWNLOAD_LOCATION: // 3.7
			// cardinality: exactly 1
			err = setDocumentLocationFromURI(pkg, subTriple.Object.ID)
		case SPDX_FILES_ANALYZED: // 3.8
			// cardinality: max 1
			err = setFilesAnalyzed(pkg, subTriple.Object.ID)
		case SPDX_PACKAGE_VERIFICATION_CODE: // 3.9
			// cardinality: max 1
			err = parser.setPackageVerificationCode(pkg, subTriple.Object)
		case SPDX_CHECKSUM: // 3.10
			// cardinality: min 0
			err = parser.setPackageChecksum(pkg, subTriple.Object)
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
			externalDocRef, err := parser.getPackageExternalRef(subTriple.Object)
			if err != nil {
				return nil, fmt.Errorf("error parsing externalRef of a package: %v", err)
			}
			pkg.PackageExternalReferences = append(pkg.PackageExternalReferences, externalDocRef)
		case SPDX_HAS_FILE: // 3.22
			// cardinality: min 0
			file, err := parser.getFileFromNode(subTriple.Object)
			if err != nil {
				return nil, fmt.Errorf("error setting file inside a package: %v", err)
			}
			parser.setFileToPackage(pkg, file)
		case SPDX_RELATIONSHIP:
			// cardinality: min 0
			err = parser.parseRelationship(subTriple)
		case SPDX_ATTRIBUTION_TEXT:
			// cardinality: min 0
			pkg.PackageAttributionTexts = append(pkg.PackageAttributionTexts, subTriple.Object.ID)
		case SPDX_ANNOTATION:
			// cardinality: min 0
			err = parser.parseAnnotationFromNode(subTriple.Object)
		default:
			return nil, fmt.Errorf("unknown predicate id %s while parsing a package", subTriple.Predicate.ID)
		}
		if err != nil {
			return nil, err
		}
	}

	if existingPackageIndex != -1 {
		parser.doc.Packages[existingPackageIndex] = pkg
	} else {
		parser.doc.Packages = append(parser.doc.Packages, pkg)
	}

	return pkg, nil
}

// parses externalReference found in the package by the associated triple.
func (parser *rdfParser2_2) getPackageExternalRef(node *gordfParser.Node) (externalDocRef *v2_2.PackageExternalReference, err error) {
	externalDocRef = &v2_2.PackageExternalReference{}
	for _, triple := range parser.nodeToTriples(node) {
		switch triple.Predicate.ID {
		case SPDX_REFERENCE_CATEGORY:
			// cardinality: exactly 1
			switch triple.Object.ID {
			case SPDX_REFERENCE_CATEGORY_SECURITY:
				externalDocRef.Category = "SECURITY"
			case SPDX_REFERENCE_CATEGORY_PACKAGE_MANAGER:
				externalDocRef.Category = "PACKAGE-MANAGER"
			case SPDX_REFERENCE_CATEGORY_OTHER:
				externalDocRef.Category = "OTHER"
			default:
				return nil, fmt.Errorf("unknown packageManager uri %s", triple.Predicate.ID)
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
			externalDocRef.RefType = triple.Object.ID
		case SPDX_REFERENCE_LOCATOR:
			// cardinality: exactly 1
			externalDocRef.Locator = triple.Object.ID
		case RDFS_COMMENT:
			// cardinality: max 1
			externalDocRef.ExternalRefComment = triple.Object.ID
		default:
			return nil, fmt.Errorf("unknown package external reference predicate id %s", triple.Predicate.ID)
		}
	}
	return
}

func (parser *rdfParser2_2) setPackageVerificationCode(pkg *v2_2.Package, node *gordfParser.Node) error {
	for _, subTriple := range parser.nodeToTriples(node) {
		switch subTriple.Predicate.ID {
		case SPDX_PACKAGE_VERIFICATION_CODE_VALUE:
			// cardinality: exactly 1
			pkg.PackageVerificationCode.Value = subTriple.Object.ID
		case SPDX_PACKAGE_VERIFICATION_CODE_EXCLUDED_FILE:
			// cardinality: min 0
			pkg.PackageVerificationCode.ExcludedFiles = append(pkg.PackageVerificationCode.ExcludedFiles, subTriple.Object.ID)
		case RDF_TYPE:
			// cardinality: exactly 1
			continue
		default:
			return fmt.Errorf("unparsed predicate %s", subTriple.Predicate.ID)
		}
	}
	return nil
}

// appends the file to the package and also sets the assocWithPackage for the
// file to indicate the file is associated with a package
func (parser *rdfParser2_2) setFileToPackage(pkg *v2_2.Package, file *v2_2.File) {
	if pkg.Files == nil {
		pkg.Files = []*v2_2.File{}
	}
	pkg.Files = append(pkg.Files, file)
	parser.assocWithPackage[file.FileSPDXIdentifier] = true
}

// given a supplierObject, sets the PackageSupplier attribute of the pkg.
// Args:
//
//	value: [NOASSERTION | [Person | Organization]: string]
func setPackageSupplier(pkg *v2_2.Package, value string) error {
	value = strings.TrimSpace(value)
	supplier := &common.Supplier{}
	if strings.ToUpper(value) == "NOASSERTION" {
		supplier.Supplier = "NOASSERTION"
		pkg.PackageSupplier = supplier
		return nil
	}

	subKey, subValue, err := ExtractSubs(value, ":")
	if err != nil {
		return fmt.Errorf("package supplier must be of the form NOASSERTION or [Person|Organization]: string. found: %s", value)
	}
	switch subKey {
	case "Person", "Organization":
		supplier.Supplier = subValue
		supplier.SupplierType = subKey
	default:
		return fmt.Errorf("unknown supplier %s", subKey)
	}

	pkg.PackageSupplier = supplier

	return nil
}

// given a OriginatorObject, sets the PackageOriginator attribute of the pkg.
// Args:
//
//	value: [NOASSERTION | [Person | Organization]: string]
func setPackageOriginator(pkg *v2_2.Package, value string) error {
	value = strings.TrimSpace(value)
	originator := &common.Originator{}
	if strings.ToUpper(value) == "NOASSERTION" {
		originator.Originator = "NOASSERTION"
		pkg.PackageOriginator = originator
		return nil
	}

	subKey, subValue, err := ExtractSubs(value, ":")
	if err != nil {
		return fmt.Errorf("package Originator must be of the form NOASSERTION or [Person|Organization]: string. found: %s", value)
	}
	switch subKey {
	case "Person", "Organization":
		originator.Originator = subValue
		originator.OriginatorType = subKey
	default:
		return fmt.Errorf("unknown Originator %s", subKey)
	}

	pkg.PackageOriginator = originator

	return nil
}

// validates the uri and sets the location if it is valid
func setDocumentLocationFromURI(pkg *v2_2.Package, locationURI string) error {
	switch locationURI {
	case SPDX_NOASSERTION_CAPS, SPDX_NOASSERTION_SMALL:
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
func setFilesAnalyzed(pkg *v2_2.Package, boolValue string) (err error) {
	pkg.IsFilesAnalyzedTagPresent = true
	pkg.FilesAnalyzed, err = boolFromString(boolValue)
	return err
}

func (parser *rdfParser2_2) setPackageChecksum(pkg *v2_2.Package, node *gordfParser.Node) error {
	checksumAlgorithm, checksumValue, err := parser.getChecksumFromNode(node)
	if err != nil {
		return fmt.Errorf("error getting checksum algorithm and value from %v", node)
	}
	if pkg.PackageChecksums == nil {
		pkg.PackageChecksums = make([]common.Checksum, 0, 1)
	}
	switch checksumAlgorithm {
	case common.SHA1,
		common.SHA224,
		common.SHA256,
		common.SHA384,
		common.SHA512,
		common.MD2,
		common.MD4,
		common.MD5,
		common.MD6:
		pkg.PackageChecksums = append(pkg.PackageChecksums, common.Checksum{Algorithm: checksumAlgorithm, Value: checksumValue})
	default:
		return fmt.Errorf("unknown checksumAlgorithm %s while parsing a package", checksumAlgorithm)
	}
	return nil
}
