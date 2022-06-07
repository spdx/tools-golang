// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/spdx/tools-golang/utils"
	"sort"
	"strings"
)

type Supplier struct {
	// can be "NOASSERTION"
	Supplier string
	// SupplierType can be one of "Person", "Organization", or empty if Supplier is "NOASSERTION"
	SupplierType string
}

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (s Supplier) Validate() error {
	// SupplierType is allowed to be empty if Supplier is "NOASSERTION"
	if s.Supplier == "" || (s.SupplierType == "" && s.Supplier != "NOASSERTION") {
		return fmt.Errorf("invalid Supplier, missing fields. %+v", s)
	}

	return nil
}

// FromString parses a string into a Supplier.
// These stings take the form: "<SupplierType>: <Supplier>"
func (s *Supplier) FromString(value string) error {
	if value == "NOASSERTION" {
		s.Supplier = value
		return nil
	}

	supplierType, supplier, err := utils.ExtractSubs(value)
	if err != nil {
		return err
	}

	s.SupplierType = supplierType
	s.Supplier = supplier

	return nil
}

// String converts the Supplier to a string in the form "<SupplierType>: <Supplier>"
func (s Supplier) String() string {
	if s.Supplier == "NOASSERTION" {
		return s.Supplier
	}

	return fmt.Sprintf("%s: %s", s.SupplierType, s.Supplier)
}

// UnmarshalJSON takes a supplier in the typical one-line format and parses it into a Supplier struct.
// This function is also used when unmarshalling YAML
func (s *Supplier) UnmarshalJSON(data []byte) error {
	// the value is just a string presented as a slice of bytes
	supplierStr := string(data)
	supplierStr = strings.Trim(supplierStr, "\"")

	return s.FromString(supplierStr)
}

// MarshalJSON converts the receiver into a slice of bytes representing a Supplier in string form.
// This function is also used when marshalling to YAML
func (s Supplier) MarshalJSON() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(s.String())
}

type Originator struct {
	// can be "NOASSERTION"
	Originator string
	// OriginatorType can be one of "Person", "Organization", or empty if Originator is "NOASSERTION"
	OriginatorType string
}

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (o Originator) Validate() error {
	// Originator is allowed to be empty if Originator is "NOASSERTION"
	if o.Originator == "" || (o.OriginatorType == "" && o.Originator != "NOASSERTION") {
		return fmt.Errorf("invalid Originator, missing fields. %+v", o)
	}

	return nil
}

// FromString parses a string into a Originator.
// These stings take the form: "<OriginatorType>: <Originator>"
func (o *Originator) FromString(value string) error {
	if value == "NOASSERTION" {
		o.Originator = value
		return nil
	}

	originatorType, originator, err := utils.ExtractSubs(value)
	if err != nil {
		return err
	}

	o.OriginatorType = originatorType
	o.Originator = originator

	return nil
}

// String converts the Originator to a string in the form "<OriginatorType>: <Originator>"
func (o Originator) String() string {
	if o.Originator == "NOASSERTION" {
		return o.Originator
	}

	return fmt.Sprintf("%s: %s", o.OriginatorType, o.Originator)
}

// UnmarshalJSON takes an originator in the typical one-line format and parses it into an Originator struct.
// This function is also used when unmarshalling YAML
func (o *Originator) UnmarshalJSON(data []byte) error {
	// the value is just a string presented as a slice of bytes
	originatorStr := string(data)
	originatorStr = strings.Trim(originatorStr, "\"")

	return o.FromString(originatorStr)
}

// MarshalJSON converts the receiver into a slice of bytes representing an Originator in string form.
// This function is also used when marshalling to YAML
func (o Originator) MarshalJSON() ([]byte, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(o.String())
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
	PackageVersion string `json:"versionInfo,omitempty"`

	// 3.4: Package File Name
	// Cardinality: optional, one
	PackageFileName string `json:"packageFileName,omitempty"`

	// 3.5: Package Supplier: may have single result for either Person or Organization,
	//                        or NOASSERTION
	// Cardinality: optional, one
	PackageSupplier *Supplier `json:"supplier,omitempty"`

	// 3.6: Package Originator: may have single result for either Person or Organization,
	//                          or NOASSERTION
	// Cardinality: optional, one
	PackageOriginator *Originator `json:"originator,omitempty"`

	// 3.7: Package Download Location
	// Cardinality: mandatory, one
	PackageDownloadLocation string `json:"downloadLocation"`

	// 3.8: FilesAnalyzed
	// Cardinality: optional, one; default value is "true" if omitted
	FilesAnalyzed bool `json:"filesAnalyzed,omitempty"`
	// NOT PART OF SPEC: did FilesAnalyzed tag appear?
	IsFilesAnalyzedTagPresent bool `json:"-"`

	// 3.9: Package Verification Code
	PackageVerificationCode PackageVerificationCode `json:"packageVerificationCode"`

	// 3.10: Package Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: optional, one or many
	PackageChecksums []Checksum `json:"checksums,omitempty"`

	// 3.11: Package Home Page
	// Cardinality: optional, one
	PackageHomePage string `json:"homepage,omitempty"`

	// 3.12: Source Information
	// Cardinality: optional, one
	PackageSourceInfo string `json:"sourceInfo,omitempty"`

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
	PackageLicenseComments string `json:"licenseComments,omitempty"`

	// 3.17: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageCopyrightText string `json:"copyrightText"`

	// 3.18: Package Summary Description
	// Cardinality: optional, one
	PackageSummary string `json:"summary,omitempty"`

	// 3.19: Package Detailed Description
	// Cardinality: optional, one
	PackageDescription string `json:"description,omitempty"`

	// 3.20: Package Comment
	// Cardinality: optional, one
	PackageComment string `json:"comment,omitempty"`

	// 3.21: Package External Reference
	// Cardinality: optional, one or many
	PackageExternalReferences []*PackageExternalReference2_1 `json:"externalRefs,omitempty"`

	// Files contained in this Package
	Files []*File2_1

	Annotations []Annotation2_1 `json:"annotations,omitempty"`
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
	PackageVersion string `json:"versionInfo,omitempty"`

	// 3.4: Package File Name
	// Cardinality: optional, one
	PackageFileName string `json:"packageFileName,omitempty"`

	// 3.5: Package Supplier: may have single result for either Person or Organization,
	//                        or NOASSERTION
	// Cardinality: optional, one
	PackageSupplier *Supplier `json:"supplier,omitempty"`

	// 3.6: Package Originator: may have single result for either Person or Organization,
	//                          or NOASSERTION
	// Cardinality: optional, one
	PackageOriginator *Originator `json:"originator,omitempty"`

	// 3.7: Package Download Location
	// Cardinality: mandatory, one
	PackageDownloadLocation string `json:"downloadLocation"`

	// 3.8: FilesAnalyzed
	// Cardinality: optional, one; default value is "true" if omitted
	FilesAnalyzed bool `json:"filesAnalyzed,omitempty"`
	// NOT PART OF SPEC: did FilesAnalyzed tag appear?
	IsFilesAnalyzedTagPresent bool

	// 3.9: Package Verification Code
	PackageVerificationCode PackageVerificationCode `json:"packageVerificationCode"`

	// 3.10: Package Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: optional, one or many
	PackageChecksums []Checksum `json:"checksums"`

	// 3.11: Package Home Page
	// Cardinality: optional, one
	PackageHomePage string `json:"homepage,omitempty"`

	// 3.12: Source Information
	// Cardinality: optional, one
	PackageSourceInfo string `json:"sourceInfo,omitempty"`

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
	PackageLicenseComments string `json:"licenseComments,omitempty"`

	// 3.17: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	PackageCopyrightText string `json:"copyrightText"`

	// 3.18: Package Summary Description
	// Cardinality: optional, one
	PackageSummary string `json:"summary,omitempty"`

	// 3.19: Package Detailed Description
	// Cardinality: optional, one
	PackageDescription string `json:"description,omitempty"`

	// 3.20: Package Comment
	// Cardinality: optional, one
	PackageComment string `json:"comment,omitempty"`

	// 3.21: Package External Reference
	// Cardinality: optional, one or many
	PackageExternalReferences []*PackageExternalReference2_2 `json:"externalRefs,omitempty"`

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	// contained within PackageExternalReference2_1 struct, if present

	// 3.23: Package Attribution Text
	// Cardinality: optional, one or many
	PackageAttributionTexts []string `json:"attributionTexts,omitempty"`

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

// MakePackageVerificationCode2_1 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.1, section 3.9.4.
func MakePackageVerificationCode2_1(files []*File2_1, excludeFile string) (PackageVerificationCode, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return PackageVerificationCode{}, fmt.Errorf("got nil file for identifier %v", i)
		}
		if f.FileName != excludeFile {
			// find the SHA1 hash, if present
			for _, checksum := range f.Checksums {
				if checksum.Algorithm == SHA1 {
					shas = append(shas, checksum.Value)
				}
			}
		}
	}

	// sort the strings
	sort.Strings(shas)

	// concatenate them into one string, with no trailing separators
	shasConcat := strings.Join(shas, "")

	// and get its SHA1 value
	hsha1 := sha1.New()
	hsha1.Write([]byte(shasConcat))
	bs := hsha1.Sum(nil)

	code := PackageVerificationCode{
		Value:         fmt.Sprintf("%x", bs),
		ExcludedFiles: []string{excludeFile},
	}

	return code, nil
}

// MakePackageVerificationCode2_2 takes a slice of files and an optional filename
// for an "excludes" file, and returns a Package Verification Code calculated
// according to SPDX spec version 2.2, section 3.9.4.
func MakePackageVerificationCode2_2(files []*File2_2, excludeFile string) (PackageVerificationCode, error) {
	// create slice of strings - unsorted SHA1s for all files
	shas := []string{}
	for i, f := range files {
		if f == nil {
			return PackageVerificationCode{}, fmt.Errorf("got nil file for identifier %v", i)
		}
		if f.FileName != excludeFile {
			// find the SHA1 hash, if present
			for _, checksum := range f.Checksums {
				if checksum.Algorithm == SHA1 {
					shas = append(shas, checksum.Value)
				}
			}
		}
	}

	// sort the strings
	sort.Strings(shas)

	// concatenate them into one string, with no trailing separators
	shasConcat := strings.Join(shas, "")

	// and get its SHA1 value
	hsha1 := sha1.New()
	hsha1.Write([]byte(shasConcat))
	bs := hsha1.Sum(nil)

	code := PackageVerificationCode{
		Value:         fmt.Sprintf("%x", bs),
		ExcludedFiles: []string{excludeFile},
	}

	return code, nil
}
