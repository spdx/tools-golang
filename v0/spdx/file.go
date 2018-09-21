// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

// File2_1 is a File section of an SPDX Document for version 2.1 of the spec.
type File2_1 struct {

	// 4.1: File Name
	// Cardinality: mandatory, one
	FileName string

	// 4.2: File SPDX Identifier: "SPDXRef-[idstring]"
	// Cardinality: mandatory, one
	FileSPDXIdentifier string

	// 4.3: File Type
	// Cardinality: optional, multiple
	FileType []string

	// 4.4: File Checksum: may have keys for SHA1, SHA256 and/or MD5
	// Cardinality: mandatory, one SHA1, others may be optionally provided
	FileChecksumSHA1   string
	FileChecksumSHA256 string
	FileChecksumMD5    string

	// 4.5: Concluded License: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	LicenseConcluded string

	// 4.6: License Information in File: SPDX License Expression, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one or many
	LicenseInfoInFile []string

	// 4.7: Comments on License
	// Cardinality: optional, one
	LicenseComments string

	// 4.8: Copyright Text: copyright notice(s) text, "NONE" or "NOASSERTION"
	// Cardinality: mandatory, one
	FileCopyrightText string

	// DEPRECATED in version 2.1 of spec
	// 4.9: Artifact of Project Name
	// Cardinality: optional, one or many
	ArtifactOfProjectName []string

	// DEPRECATED in version 2.1 of spec
	// 4.10: Artifact of Project Homepage: URL or "UNKNOWN"
	// Cardinality: optional, one or many
	ArtifactOfProjectHomePage []string

	// DEPRECATED in version 2.1 of spec
	// 4.11: Artifact of Project Uniform Resource Identifier
	// Cardinality: optional, one or many
	ArtifactOfProjectURI []string

	// 4.12: File Comment
	// Cardinality: optional, one
	FileComment string

	// 4.13: File Notice
	// Cardinality: optional, one
	FileNotice string

	// 4.14: File Contributor
	// Cardinality: optional, one or many
	FileContributor []string

	// DEPRECATED in version 2.0 of spec
	// 4.15: File Dependencies
	// Cardinality: optional, one or many
	FileDependencies []string

	// Snippets contained in this File
	Snippets []*Snippet2_1

	// Relationships applicable to this File
	Relationships []*Relationship2_1

	// Annotations applicable to this File
	Annotations []*Annotation2_1
}
