package v3_0

import (
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func v23doc() *v2_3.Document {
	return &v2_3.Document{
		SPDXVersion:       v2_3.Version,
		DataLicense:       v2_3.DataLicense,
		SPDXIdentifier:    "SPDXRef-DOCUMENT",
		DocumentName:      "Example Software Package",
		DocumentComment:   "This is a sample SPDX document for testing purposes.",
		CreationInfo:      v23creationInfo(),
		DocumentNamespace: "https://example.com/spdx/example-software-1.0.0",
		ExternalDocumentReferences: []v2_3.ExternalDocumentRef{
			*v23externalDocumentRef1(),
			*v23externalDocumentRef2(),
		},
		Packages: []*v2_3.Package{
			v23package1(),
			v23package2(),
		},
		Files: []*v2_3.File{
			v23file1(),
			v23file2(),
		},
		OtherLicenses: []*v2_3.OtherLicense{
			v23customLicense1(),
			v23customLicense2(),
		},
		Annotations: []*v2_3.Annotation{
			v23annotation1(),
			v23annotation2(),
		},
		Snippets: []v2_3.Snippet{
			*v23snippet1(),
			*v23snippet2(),
		},
		Relationships: []*v2_3.Relationship{
			{
				RefA:                common.DocElementID{ElementRefID: "SPDXRef-DOCUMENT"},
				RefB:                common.DocElementID{ElementRefID: "SPDXRef-Package-ExampleLib"}, // pkg1
				Relationship:        common.TypeRelationshipDescribe,
				RelationshipComment: "Document describes the main package",
			},
			{
				RefA:                common.DocElementID{ElementRefID: "SPDXRef-Package-ExampleLib"},
				RefB:                common.DocElementID{ElementRefID: "SPDXRef-Package-UtilityTools"},
				Relationship:        common.TypeRelationshipBuildToolOf,
				RelationshipComment: "Main package depends on utility tools",
			},
			{
				RefA:                common.DocElementID{ElementRefID: "SPDXRef-Package-ExampleLib"},
				RefB:                common.DocElementID{ElementRefID: "SPDXRef-File-Main"},
				Relationship:        common.TypeRelationshipContains,
				RelationshipComment: "Package contains main source file",
			},
		},
		// Reviews are dropped
	}
}

func v23creationInfo() *v2_3.CreationInfo {
	return &v2_3.CreationInfo{
		Created: "2023-01-15T10:30:00Z",
		Creators: []common.Creator{
			{CreatorType: "Person", Creator: "John Doe (john@example.com)"},
			{CreatorType: "Tool", Creator: "tools-golang-v1.1.0"},
		},
		LicenseListVersion: "3.19",
		CreatorComment:     "Created during automated build process",
	}
}

func v23package1() *v2_3.Package {
	return &v2_3.Package{
		Annotations: []v2_3.Annotation{
			*v23annotation3(),
		},
		PackageName:           "example-library",
		PackageSPDXIdentifier: "SPDXRef-Package-ExampleLib",
		PackageDescription:    "This is a detailed description of the example library package.",
		PackageVersion:        "1.2.3",
		PackageSupplier: &common.Supplier{
			Supplier:     "Example Corp (support@example.com)",
			SupplierType: "Organization",
		},
		PackageOriginator: &common.Originator{
			Originator:     "Jane Smith (jane@example.com)",
			OriginatorType: "Person",
		},
		PackageDownloadLocation: "https://github.com/example/library/archive/v1.2.3.tar.gz",
		FilesAnalyzed:           true,
		Files: []*v2_3.File{
			v23file3(),
		},
		IsFilesAnalyzedTagPresent: true,
		PackageVerificationCode: &common.PackageVerificationCode{
			Value:         "d6a770ba38583ed4bb4525bd96e50461655d2758",
			ExcludedFiles: []string{"./exclude1.txt", "./exclude2.txt"},
		},
		PackageChecksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "aabbccdd"},
			{Algorithm: common.SHA256, Value: "11223344"},
		},
		PackageHomePage:             "https://example.com/library",
		PackageSourceInfo:           "Built from git tag v1.2.3",
		PackageLicenseConcluded:     "MIT",
		PackageLicenseInfoFromFiles: []string{"MIT", "Apache-2.0"},
		PackageLicenseDeclared:      "MIT",
		PackageLicenseComments:      "License determined from LICENSE file",
		PackageCopyrightText:        "Copyright 2023 Example Corp",
		PackageSummary:              "A sample library for demonstration",
		PackageComment:              "Package built with standard configuration",
		PackageExternalReferences: []*v2_3.PackageExternalReference{
			{
				Category:           common.CategoryPackageManager,
				RefType:            common.TypePackageManagerPURL,
				Locator:            "pkg:npm/example-library@1.2.3",
				ExternalRefComment: "NPM package reference",
			},
			{
				Category: common.CategoryPackageManager,
				RefType:  common.TypePackageManagerPURL,
				Locator:  "pkg:npm/example-library-main@1.2.3",
			},
			{
				Category:           common.CategorySecurity,
				RefType:            common.TypeSecurityCPE23Type,
				Locator:            "cpe:2.3:a:example:library:1.2.3:*:*:*:*:*:*:*",
				ExternalRefComment: "CPE reference for security scanning",
			},
		},
		PackageAttributionTexts: []string{
			"This package includes code from Project ABC",
			"Special thanks to the open source community",
		},
		PrimaryPackagePurpose: "Library", // TODO constant?
		ReleaseDate:           "2023-01-10T00:00:00Z",
		BuiltDate:             "2023-01-15T08:00:00Z",
		ValidUntilDate:        "2025-01-10T00:00:00Z",
	}
}

func v23package2() *v2_3.Package {
	return &v2_3.Package{
		PackageName:           "utility-tools",
		PackageSPDXIdentifier: "SPDXRef-Package-UtilityTools",
		PackageVersion:        "2.1.0",
		PackageSupplier: &common.Supplier{
			Supplier:     "Bob Johnson (bob@tools.com)",
			SupplierType: "Person",
		},
		PackageDownloadLocation:   "https://tools.com/download/utility-tools-2.1.0.zip",
		FilesAnalyzed:             true,
		IsFilesAnalyzedTagPresent: true,
		PackageChecksums: []common.Checksum{
			{Algorithm: common.SHA256, Value: "ffaabbcc11223344"},
		},
		PackageHomePage:         "https://tools.com/utility",
		PackageLicenseConcluded: "Apache-2.0",
		PackageLicenseDeclared:  "AGPL OR (GPL-2.0-only with Classpath-Exception OR LicenseRef-CustomLicense1)",
		PackageCopyrightText:    "Copyright 2023 Tools Inc",
		PackageSummary:          "Collection of utility tools",
		PackageDescription:      "A comprehensive set of utility tools for developers.",
		PrimaryPackagePurpose:   "Application", // TODO constant?
	}
}
func v23customLicense1() *v2_3.OtherLicense {
	return &v2_3.OtherLicense{
		LicenseIdentifier: "LicenseRef-CustomLicense1",
		ExtractedText:     "This is a custom license text for demonstration purposes.\n\nPermission is granted to use this software...",
		LicenseName:       "Custom Example License",
		LicenseCrossReferences: []string{
			"https://example.com/licenses/custom",
		},
		LicenseComment: "Custom license used for internal tools",
	}
}

func v23customLicense2() *v2_3.OtherLicense {
	return &v2_3.OtherLicense{
		LicenseIdentifier: "LicenseRef-CustomLicense2",
		ExtractedText:     "Another custom license text with different terms.\n\nThis software may be used under the following conditions...",
		LicenseName:       "Another Custom License",
		LicenseCrossReferences: []string{
			"https://example.com/licenses/another-custom",
			"https://internal.example.com/legal/licenses",
		},
		LicenseComment: "License for third-party components",
	}
}

func v23file1() *v2_3.File {
	return &v2_3.File{
		FileName:           "./src/main.c",
		FileSPDXIdentifier: common.ElementID("SPDXRef-File-Main"),
		FileTypes:          []string{"FILE"},
		Checksums: []common.Checksum{
			{Algorithm: common.SHA1, Value: "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
			{Algorithm: common.MD5, Value: "d41d8cd98f00b204e9800998ecf8427e"},
		},
		LicenseConcluded:   "MIT",
		LicenseInfoInFiles: []string{"MIT"},
		LicenseComments:    "License applies to this file",
		FileCopyrightText:  "Copyright 2023 Example Corp",
		ArtifactOfProjects: []*v2_3.ArtifactOfProject{},
		FileComment:        "Other application entry point",
		FileNotice:         "This file contains the other function",
		FileContributors:   []string{"Other John Doe (john@doe.com)", "Jane Smith <jane.smith@example.org>"},
		FileAttributionTexts: []string{
			"Based on other example code from something",
		},
		//FileDependencies: nil, // skipped
		Snippets: map[common.ElementID]*v2_3.Snippet{
			"SPDXRef-Snippet3": v23snippet3(),
		},
		Annotations: []v2_3.Annotation{
			*v23annotation4(),
		},
		FileDependencies: []string{"dep1", "dep2"},
	}
}

func v23file2() *v2_3.File {
	return &v2_3.File{
		FileName:           "./include/header.h",
		FileSPDXIdentifier: common.ElementID("SPDXRef-File-Header"),
		Checksums: []common.Checksum{
			{Algorithm: common.SHA256, Value: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		},
		LicenseConcluded:   "MIT AND Apache-2.0",
		LicenseInfoInFiles: []string{"MIT"},
		FileCopyrightText:  "Copyright 2023 Example Corp",
		FileComment:        "Header file with function declarations",
		FileTypes:          []string{"HEADER"},
	}
}

func v23externalDocumentRef1() *v2_3.ExternalDocumentRef {
	return &v2_3.ExternalDocumentRef{
		DocumentRefID: "DocumentRef-External1",
		URI:           "https://external.com/spdx/external-doc-1.0.0",
		Checksum: common.Checksum{
			Algorithm: common.SHA1,
			Value:     "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
	}
}

func v23externalDocumentRef2() *v2_3.ExternalDocumentRef {
	return &v2_3.ExternalDocumentRef{
		DocumentRefID: "DocumentRef-External2",
		URI:           "https://external2.com/spdx/external-doc-2.0.0",
		Checksum: common.Checksum{
			Algorithm: common.SHA256,
			Value:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}
}

func v23annotation1() *v2_3.Annotation {
	return &v2_3.Annotation{
		Annotator: common.Annotator{
			Annotator:     "Security Team (security@example.com)",
			AnnotatorType: "Person",
		},
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.DocElementID{ElementRefID: "SPDXRef-Package-ExampleLib"},
		AnnotationDate:           "2023-01-20T14:30:00Z",
		AnnotationComment:        "Security review completed - no vulnerabilities found",
	}
}

func v23annotation2() *v2_3.Annotation {
	return &v2_3.Annotation{
		Annotator: common.Annotator{
			Annotator:     "vulnerability-scanner-v1.5",
			AnnotatorType: "Tool",
		},
		AnnotationType:           "OTHER",
		AnnotationSPDXIdentifier: common.DocElementID{ElementRefID: "SPDXRef-Package-UtilityTools"},
		AnnotationDate:           "2023-01-21T09:15:00Z",
		AnnotationComment:        "Automated scan completed - clean",
	}
}

func v23snippet1() *v2_3.Snippet {
	return &v2_3.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("SPDXRef-Snippet1"),
		SnippetFromFileSPDXIdentifier: common.ElementID("SPDXRef-File-Main"),
		Ranges: []common.SnippetRange{
			{
				StartPointer: common.SnippetRangePointer{
					Offset:     100,
					LineNumber: 10,
				},
				EndPointer: common.SnippetRangePointer{
					Offset:     200,
					LineNumber: 15,
				},
			},
		},
		SnippetLicenseConcluded: "MIT",
		LicenseInfoInSnippet:    []string{"MIT"},
		SnippetLicenseComments:  "License applies to this code snippet",
		SnippetCopyrightText:    "Copyright 2023 Example Corp",
		SnippetComment:          "Key algorithm implementation",
		SnippetName:             "Core Algorithm",
		SnippetAttributionTexts: []string{"Algorithm based on research paper XYZ"},
	}
}

func v23snippet2() *v2_3.Snippet {
	return &v2_3.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("SPDXRef-Snippet2"),
		SnippetFromFileSPDXIdentifier: common.ElementID("SPDXRef-File-Header"),
		Ranges: []common.SnippetRange{
			{
				StartPointer: common.SnippetRangePointer{Offset: 50},
				EndPointer:   common.SnippetRangePointer{Offset: 150},
			},
			{
				StartPointer: common.SnippetRangePointer{LineNumber: 5},
				EndPointer:   common.SnippetRangePointer{LineNumber: 8},
			},
		},
		SnippetLicenseConcluded: "MIT+",
		LicenseInfoInSnippet:    []string{"MIT"},
		SnippetCopyrightText:    "Copyright 2023 Example Corp",
		SnippetComment:          "Function declarations",
		SnippetName:             "API Declarations",
	}
}

func v23file3() *v2_3.File {
	return &v2_3.File{
		FileName:           "./src/utils.c",
		FileSPDXIdentifier: common.ElementID("SPDXRef-File-Utils"),
		FileTypes:          []string{"SOURCE"},
		Checksums: []common.Checksum{
			{Algorithm: common.SHA256, Value: "aabb1234aabb1234aabb1234aabb1234aabb1234aabb1234aabb1234aabb1234"},
		},
		LicenseConcluded:   "Apache-2.0",
		LicenseInfoInFiles: []string{"Apache-2.0"},
		FileCopyrightText:  "Copyright 2023 Example Corp",
		FileComment:        "Utility functions for the library",
	}
}

func v23annotation3() *v2_3.Annotation {
	return &v2_3.Annotation{
		Annotator: common.Annotator{
			Annotator:     "Code Review Bot (bot@example.com)",
			AnnotatorType: "Person",
		},
		AnnotationType:           "OTHER",
		AnnotationSPDXIdentifier: common.DocElementID{ElementRefID: "SPDXRef-Package-ExampleLib"},
		AnnotationDate:           "2023-02-01T10:00:00Z",
		AnnotationComment:        "Package structure review completed",
	}
}

func v23annotation4() *v2_3.Annotation {
	return &v2_3.Annotation{
		Annotator: common.Annotator{
			Annotator:     "Code Reviewer (reviewer@example.com)",
			AnnotatorType: "Person",
		},
		AnnotationType:           "REVIEW",
		AnnotationSPDXIdentifier: common.DocElementID{ElementRefID: "SPDXRef-File-Main"},
		AnnotationDate:           "2023-02-05T16:00:00Z",
		AnnotationComment:        "File review completed - code quality approved",
	}
}

func v23snippet3() *v2_3.Snippet {
	return &v2_3.Snippet{
		SnippetSPDXIdentifier:         common.ElementID("SPDXRef-Snippet3"),
		SnippetFromFileSPDXIdentifier: common.ElementID("SPDXRef-File-Main"),
		Ranges: []common.SnippetRange{
			{
				StartPointer: common.SnippetRangePointer{
					Offset:     300,
					LineNumber: 20,
				},
				EndPointer: common.SnippetRangePointer{
					Offset:     400,
					LineNumber: 30,
				},
			},
		},
		SnippetLicenseConcluded: "Apache-2.0",
		LicenseInfoInSnippet:    []string{"Apache-2.0"},
		SnippetLicenseComments:  "License for helper routine",
		SnippetCopyrightText:    "Copyright 2023 Example Corp",
		SnippetComment:          "Helper routine implementation",
		SnippetName:             "Helper Routine",
		SnippetAttributionTexts: []string{"Inspired by open source utilities"},
	}
}
