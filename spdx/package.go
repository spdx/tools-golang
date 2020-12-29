// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

// Package2_1 is a Package section of an SPDX Document for version 2.1 of the spec.
type Package2_1 struct {

	// 3.1: Package Name
	// Cardinality: mandatory, one
	Name string

	// 3.2: Package SPDX Identifier: "SPDXRef-[idstring]"
	// Cardinality: mandatory, one
	SPDXIdentifier ElementID

	// 3.3: Package Version
	// Cardinality: optional, one
	Version string

	// 3.4: Package File Name
	// Cardinality: optional, one
	FileName string

	// 3.5: Package Supplier: may have single result for either Person or Organization,
	//                        or NOASSERTION
	// Cardinality: optional, one
	SupplierPerson       string
	SupplierOrganization string
	SupplierNOASSERTION  bool

	// 3.6: Package Originator: may have single result for either Person or Organization,
	//                          or NOASSERTION
	// Cardinality: optional, one
	OriginatorPerson       string
	OriginatorOrganization string
	OriginatorNOASSERTION  bool

	// 3.7: Package Download Location
	// Cardinality: mandatory, one
	DownloadLocation string

	// 3.8: FilesAnalyzed
	// Cardinality: optional, one; default value is "true" if omitted
	FilesAnalyzed bool
	// NOT PART OF SPEC: did FilesAnalyzed tag appear?
	IsFilesAnalyzedTagPresent bool

	// 3.9: Package Verification Code
	// Cardinality: mandatory, one if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	VerificationCode string
	// Spec also allows specifying a single file to exclude from the
	// verification code algorithm; intended to enable exclusion of
	// the SPDX document file itself.
	VerificationCodeExcludedFile string

	// 3.10: Package Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: optional, one or many
	ChecksumSHA1   string
	ChecksumSHA256 string
	ChecksumMD5    string

	// 3.11: Package Home Page
	// Cardinality: optional, one
	HomePage string

	// 3.12: Source Information
	// Cardinality: optional, one
	SourceInfo string

	// 3.13: Concluded License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	LicenseConcluded string

	// 3.14: All Licenses Info from Files: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one or many if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	LicenseInfoFromFiles []string

	// 3.15: Declared License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	LicenseDeclared string

	// 3.16: Comments on License
	// Cardinality: optional, one
	LicenseComments string

	// 3.17: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	CopyrightText string

	// 3.18: Package Summary Description
	// Cardinality: optional, one
	Summary string

	// 3.19: Package Detailed Description
	// Cardinality: optional, one
	Description string

	// 3.20: Package Comment
	// Cardinality: optional, one
	Comment string

	// 3.21: Package External Reference
	// Cardinality: optional, one or many
	ExternalReferences []*PackageExternalReference2_1

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	// contained within PackageExternalReference2_1 struct, if present

	// Files contained in this Package
	Files map[ElementID]*File2_1
}

// PackageExternalReference2_1 is an External Reference to additional info
// about a Package, as defined in section 3.21 in version 2.1 of the spec.
type PackageExternalReference2_1 struct {

	// category is "SECURITY", "PACKAGE-MANAGER" or "OTHER"
	Category string

	// type is an [idstring] as defined in Appendix VI;
	// called RefType here due to "type" being a Golang keyword
	RefType string

	// locator is a unique string to access the package-specific
	// info, metadata or content within the target location
	Locator string

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	ExternalRefComment string
}

// Package2_2 is a Package section of an SPDX Document for version 2.2 of the spec.
type Package2_2 struct {

	// NOT PART OF SPEC
	// flag: does this "package" contain files that were in fact "unpackaged",
	// e.g. included directly in the Document without being in a Package?
	IsUnpackaged bool

	// 3.1: Package Name
	// Cardinality: mandatory, one
	Name string

	// 3.2: Package SPDX Identifier: "SPDXRef-[idstring]"
	// Cardinality: mandatory, one
	SPDXIdentifier ElementID

	// 3.3: Package Version
	// Cardinality: optional, one
	Version string

	// 3.4: Package File Name
	// Cardinality: optional, one
	FileName string

	// 3.5: Package Supplier: may have single result for either Person or Organization,
	//                        or NOASSERTION
	// Cardinality: optional, one
	SupplierPerson       string
	SupplierOrganization string
	SupplierNOASSERTION  bool

	// 3.6: Package Originator: may have single result for either Person or Organization,
	//                          or NOASSERTION
	// Cardinality: optional, one
	OriginatorPerson       string
	OriginatorOrganization string
	OriginatorNOASSERTION  bool

	// 3.7: Package Download Location
	// Cardinality: mandatory, one
	DownloadLocation string

	// 3.8: FilesAnalyzed
	// Cardinality: optional, one; default value is "true" if omitted
	FilesAnalyzed bool
	// NOT PART OF SPEC: did FilesAnalyzed tag appear?
	IsFilesAnalyzedTagPresent bool

	// 3.9: Package Verification Code
	// Cardinality: mandatory, one if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	VerificationCode string
	// Spec also allows specifying a single file to exclude from the
	// verification code algorithm; intended to enable exclusion of
	// the SPDX document file itself.
	VerificationCodeExcludedFile string

	// 3.10: Package Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: optional, one or many
	ChecksumSHA1   string
	ChecksumSHA256 string
	ChecksumMD5    string

	// 3.11: Package Home Page
	// Cardinality: optional, one
	HomePage string

	// 3.12: Source Information
	// Cardinality: optional, one
	SourceInfo string

	// 3.13: Concluded License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	LicenseConcluded string

	// 3.14: All Licenses Info from Files: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one or many if filesAnalyzed is true / omitted;
	//              zero (must be omitted) if filesAnalyzed is false
	LicenseInfoFromFiles []string

	// 3.15: Declared License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	LicenseDeclared string

	// 3.16: Comments on License
	// Cardinality: optional, one
	LicenseComments string

	// 3.17: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	CopyrightText string

	// 3.18: Package Summary Description
	// Cardinality: optional, one
	Summary string

	// 3.19: Package Detailed Description
	// Cardinality: optional, one
	Description string

	// 3.20: Package Comment
	// Cardinality: optional, one
	Comment string

	// 3.21: Package External Reference
	// Cardinality: optional, one or many
	ExternalReferences []*PackageExternalReference2_2

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	// contained within PackageExternalReference2_1 struct, if present

	// 3.23: Package Attribution Text
	// Cardinality: optional, one or many
	AttributionTexts []string

	// Files contained in this Package
	Files map[ElementID]*File2_2
}

// PackageExternalReference2_2 is an External Reference to additional info
// about a Package, as defined in section 3.21 in version 2.2 of the spec.
type PackageExternalReference2_2 struct {

	// category is "SECURITY", "PACKAGE-MANAGER", "PERSISTENT-ID" or "OTHER"
	Category string

	// type is an [idstring] as defined in Appendix VI;
	// called RefType here due to "type" being a Golang keyword
	RefType string

	// locator is a unique string to access the package-specific
	// info, metadata or content within the target location
	Locator string

	// 3.22: Package External Reference Comment
	// Cardinality: conditional (optional, one) for each External Reference
	ExternalRefComment string
}
