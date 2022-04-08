// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Supplier struct {
	// can be "NOASSERTION"
	Supplier string
	// SupplierType can be one of "Person", "Organization", or empty if Supplier is "NOASSERTION"
	SupplierType string
}

// UnmarshalJSON takes a supplier in the typical one-line format and parses it into a Supplier struct.
// This function is also used when unmarshalling YAML
func (s *Supplier) UnmarshalJSON(data []byte) error {
	// the value is just a string presented as a slice of bytes
	supplierStr := string(data)
	supplierStr = strings.Trim(supplierStr, "\"")

	if supplierStr == "NOASSERTION" {
		s.Supplier = supplierStr
		return nil
	}

	supplierFields := strings.SplitN(supplierStr, ": ", 2)

	if len(supplierFields) != 2 {
		return fmt.Errorf("failed to parse Supplier '%s'", supplierStr)
	}

	s.SupplierType = supplierFields[0]
	s.Supplier = supplierFields[1]

	return nil
}

// MarshalJSON converts the receiver into a slice of bytes representing a Supplier in string form.
// This function is also used when marshalling to YAML
func (s Supplier) MarshalJSON() ([]byte, error) {
	if s.Supplier == "NOASSERTION" {
		return json.Marshal(s.Supplier)
	} else if s.Supplier != "" {
		return json.Marshal(fmt.Sprintf("%s: %s", s.SupplierType, s.Supplier))
	}

	return []byte{}, nil
}

type Originator struct {
	// can be "NOASSERTION"
	Originator string
	// OriginatorType can be one of "Person", "Organization", or empty if Originator is "NOASSERTION"
	OriginatorType string
}

// UnmarshalJSON takes an originator in the typical one-line format and parses it into an Originator struct.
// This function is also used when unmarshalling YAML
func (o *Originator) UnmarshalJSON(data []byte) error {
	// the value is just a string presented as a slice of bytes
	originatorStr := string(data)
	originatorStr = strings.Trim(originatorStr, "\"")

	if originatorStr == "NOASSERTION" {
		o.Originator = originatorStr
		return nil
	}

	originatorFields := strings.SplitN(originatorStr, ": ", 2)

	if len(originatorFields) != 2 {
		return fmt.Errorf("failed to parse Originator '%s'", originatorStr)
	}

	o.OriginatorType = originatorFields[0]
	o.Originator = originatorFields[1]

	return nil
}

// MarshalJSON converts the receiver into a slice of bytes representing an Originator in string form.
// This function is also used when marshalling to YAML
func (o Originator) MarshalJSON() ([]byte, error) {
	if o.Originator == "NOASSERTION" {
		return json.Marshal(o.Originator)
	} else if o.Originator != "" {
		return json.Marshal(fmt.Sprintf("%s: %s", o.OriginatorType, o.Originator))
	}

	return []byte{}, nil
}

type PackageVerificationCode struct {
	// Cardinality: mandatory, one if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	Value string `json:"packageVerificationCodeValue"`
	// Spec also allows specifying files to exclude from the
	// verification code algorithm; intended to enable exclusion of
	// the SPDX document file itself.
	ExcludedFiles []string `json:"packageVerificationCodeExcludedFiles"`
}

// Package2_1 is a Package section of an SPDX Document for version 2.1 of the spec.
type Package2_1 struct {
	// 3.1: Package Name
	// Cardinality: mandatory, one
	PackageName string `json:"name"`

	// 3.2: Package SPDX Identifier: "SPDXRef-[idstring]"
	// Cardinality: mandatory, one
	PackageSPDXIdentifier ElementID `json:"SPDXID"`

	// 3.3: Package Version
	// Cardinality: optional, one
	PackageVersion string `json:"versionInfo"`

	// 3.4: Package File Name
	// Cardinality: optional, one
	PackageFileName string `json:"packageFileName"`

	// 3.5: Package Supplier: may have single result for either Person or Organization,
	//                        or NOASSERTION
	// Cardinality: optional, one
	PackageSupplier Supplier `json:"supplier"`

	// 3.6: Package Originator: may have single result for either Person or Organization,
	//                          or NOASSERTION
	// Cardinality: optional, one
	PackageOriginator Originator `json:"originator"`

	// 3.7: Package Download Location
	// Cardinality: mandatory, one
	PackageDownloadLocation string `json:"downloadLocation"`

	// 3.8: FilesAnalyzed
	// Cardinality: optional, one; default value is "true" if omitted
	FilesAnalyzed bool `json:"filesAnalyzed"`
	// NOT PART OF SPEC: did FilesAnalyzed tag appear?
	IsFilesAnalyzedTagPresent bool

	// 3.9: Package Verification Code
	PackageVerificationCode PackageVerificationCode `json:"packageVerificationCode"`

	// 3.10: Package Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: optional, one or many
	PackageChecksums []Checksum `json:"checksums"`

	// 3.11: Package Home Page
	// Cardinality: optional, one
	PackageHomePage string `json:"homepage"`

	// 3.12: Source Information
	// Cardinality: optional, one
	PackageSourceInfo string `json:"sourceInfo"`

	// 3.13: Concluded License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageLicenseConcluded string `json:"licenseConcluded"`

	// 3.14: All Licenses Info from Files: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one or many if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	PackageLicenseInfoFromFiles []string `json:"licenseInfoFromFiles"`

	// 3.15: Declared License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageLicenseDeclared string `json:"licenseDeclared"`

	// 3.16: Comments on License
	// Cardinality: optional, one
	PackageLicenseComments string `json:"licenseComments"`

	// 3.17: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageCopyrightText string `json:"copyrightText"`

	// 3.18: Package Summary Description
	// Cardinality: optional, one
	PackageSummary string `json:"summary"`

	// 3.19: Package Detailed Description
	// Cardinality: optional, one
	PackageDescription string `json:"description"`

	// 3.20: Package Comment
	// Cardinality: optional, one
	PackageComment string `json:"comment"`

	// 3.21: Package External Reference
	// Cardinality: optional, one or many
	PackageExternalReferences []*PackageExternalReference2_1 `json:"externalRefs"`

	// Files contained in this Package
	Files []*File2_1

	Annotations []Annotation2_1 `json:"annotations"`
}

// PackageExternalReference2_1 is an External Reference to additional info
// about a Package, as defined in section 3.21 in version 2.1 of the spec.
type PackageExternalReference2_1 struct {
	// category is "SECURITY", "PACKAGE-MANAGER" or "OTHER"
	Category string `json:"referenceCategory"`

	// type is an [idstring] as defined in Appendix VI;
	// called RefType here due to "type" being a Golang keyword
	RefType string `json:"referenceType"`

	// locator is a unique string to access the package-specific
	// info, metadata or content within the target location
	Locator string `json:"referenceLocator"`

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	ExternalRefComment string `json:"comment"`
}

// Package2_2 is a Package section of an SPDX Document for version 2.2 of the spec.
type Package2_2 struct {
	// NOT PART OF SPEC
	// flag: does this "package" contain files that were in fact "unpackaged",
	// e.g. included directly in the Document without being in a Package?
	IsUnpackaged bool

	// 3.1: Package Name
	// Cardinality: mandatory, one
	PackageName string `json:"name"`

	// 3.2: Package SPDX Identifier: "SPDXRef-[idstring]"
	// Cardinality: mandatory, one
	PackageSPDXIdentifier ElementID `json:"SPDXID"`

	// 3.3: Package Version
	// Cardinality: optional, one
	PackageVersion string `json:"versionInfo"`

	// 3.4: Package File Name
	// Cardinality: optional, one
	PackageFileName string `json:"packageFileName"`

	// 3.5: Package Supplier: may have single result for either Person or Organization,
	//                        or NOASSERTION
	// Cardinality: optional, one
	PackageSupplier Supplier `json:"supplier"`

	// 3.6: Package Originator: may have single result for either Person or Organization,
	//                          or NOASSERTION
	// Cardinality: optional, one
	PackageOriginator Originator `json:"originator"`

	// 3.7: Package Download Location
	// Cardinality: mandatory, one
	PackageDownloadLocation string `json:"downloadLocation"`

	// 3.8: FilesAnalyzed
	// Cardinality: optional, one; default value is "true" if omitted
	FilesAnalyzed bool `json:"filesAnalyzed"`
	// NOT PART OF SPEC: did FilesAnalyzed tag appear?
	IsFilesAnalyzedTagPresent bool

	// 3.9: Package Verification Code
	PackageVerificationCode PackageVerificationCode `json:"packageVerificationCode"`

	// 3.10: Package Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: optional, one or many
	PackageChecksums []Checksum `json:"checksums"`

	// 3.11: Package Home Page
	// Cardinality: optional, one
	PackageHomePage string `json:"homepage"`

	// 3.12: Source Information
	// Cardinality: optional, one
	PackageSourceInfo string `json:"sourceInfo"`

	// 3.13: Concluded License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageLicenseConcluded string `json:"licenseConcluded"`

	// 3.14: All Licenses Info from Files: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one or many if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	PackageLicenseInfoFromFiles []string `json:"licenseInfoFromFiles"`

	// 3.15: Declared License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageLicenseDeclared string `json:"licenseDeclared"`

	// 3.16: Comments on License
	// Cardinality: optional, one
	PackageLicenseComments string `json:"licenseComments"`

	// 3.17: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageCopyrightText string `json:"copyrightText"`

	// 3.18: Package Summary Description
	// Cardinality: optional, one
	PackageSummary string `json:"summary"`

	// 3.19: Package Detailed Description
	// Cardinality: optional, one
	PackageDescription string `json:"description"`

	// 3.20: Package Comment
	// Cardinality: optional, one
	PackageComment string `json:"comment"`

	// 3.21: Package External Reference
	// Cardinality: optional, one or many
	PackageExternalReferences []*PackageExternalReference2_2 `json:"externalRefs"`

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	// contained within PackageExternalReference2_1 struct, if present

	// 3.23: Package Attribution Text
	// Cardinality: optional, one or many
	PackageAttributionTexts []string `json:"attributionTexts"`

	// Files contained in this Package
	Files []*File2_2

	Annotations []Annotation2_2 `json:"annotations"`
}

// PackageExternalReference2_2 is an External Reference to additional info
// about a Package, as defined in section 3.21 in version 2.2 of the spec.
type PackageExternalReference2_2 struct {
	// category is "SECURITY", "PACKAGE-MANAGER" or "OTHER"
	Category string `json:"referenceCategory"`

	// type is an [idstring] as defined in Appendix VI;
	// called RefType here due to "type" being a Golang keyword
	RefType string `json:"referenceType"`

	// locator is a unique string to access the package-specific
	// info, metadata or content within the target location
	Locator string `json:"referenceLocator"`

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	ExternalRefComment string `json:"comment"`
}
