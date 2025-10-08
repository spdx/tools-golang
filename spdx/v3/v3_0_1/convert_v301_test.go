package v3_0_1

import (
	"time"

	"github.com/kzantow/go-ld"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func v301doc() *Document {
	ci := v301creationInfo()

	d := NewDocument(
		ProfileIdentifierType_Software,
		"Example Software Package",
		ci.CreatedBy[0],
		ci.CreatedUsing[0],
	)
	d.DataLicense = &ListedLicense{
		License: License{
			ExtendableLicense: ExtendableLicense{
				LicenseInfo: LicenseInfo{
					Element: Element{
						Name: v2_3.DataLicense,
					},
				},
			},
		},
	}
	d.SpdxDocument.ID = "SPDXRef-DOCUMENT"
	d.CreationInfo.(*CreationInfo).Created = parseTime("2023-01-15T10:30:00Z")
	d.SpdxDocument.Element.Comment = "This is a sample SPDX document for testing purposes."
	d.NamespaceMaps = NamespaceMapList{
		&NamespaceMap{
			Namespace: "https://example.com/spdx/example-software-1.0.0",
		},
	}

	sbom := &SBOM{}
	d.RootElements = ElementList{sbom}

	// Add external document references
	d.Imports = ExternalMapList{
		v301externalMap1(),
		v301externalMap2(),
	}

	sbom.Elements = append(sbom.Elements,
		v301customLicense1(),
		v301customLicense2(),
	)

	pkg1, extraElems := v301package1()
	sbom.Elements = append(sbom.Elements, extraElems...)
	sbom.Elements = append(sbom.Elements, pkg1)

	// pkg1 is the only thing referenced directly as a document root element
	sbom.RootElements = append(sbom.RootElements, pkg1)

	pkg2, extraElems := v301package2()
	sbom.Elements = append(sbom.Elements, extraElems...)
	sbom.Elements = append(sbom.Elements, pkg2)

	file1, extraElems := v301file1()
	sbom.Elements = append(sbom.Elements, file1)
	sbom.Elements = append(sbom.Elements, extraElems...)

	file2, extraElems := v301file2()
	sbom.Elements = append(sbom.Elements, file2)
	sbom.Elements = append(sbom.Elements, extraElems...)

	sbom.Elements = append(sbom.Elements,
		v301annotation1(pkg1),
		v301annotation2(pkg2),
	)

	snippet1, extraElems := v301snippet1(file1)
	sbom.Elements = append(sbom.Elements, snippet1)
	sbom.Elements = append(sbom.Elements, extraElems...)

	snippet2, extraElems := v301snippet2(file2)
	sbom.Elements = append(sbom.Elements, snippet2)
	sbom.Elements = append(sbom.Elements, extraElems...)

	sbom.Elements = append(sbom.Elements,
		// SpdxRef-DOCUMENT relationships are handled by object structure
		&Relationship{
			From: pkg1,
			To:   ElementList{pkg2},
			Type: RelationshipType_DependsOn,
			Element: Element{
				Comment: "Main package depends on utility tools",
			},
		},
		&Relationship{
			From: pkg1,
			To:   ElementList{file1},
			Type: RelationshipType_Contains,
			Element: Element{
				Comment: "Package contains main source file",
			},
		},
	)

	return d
}

func v301creationInfo() *CreationInfo {
	return &CreationInfo{

		Created: parseTime("2023-01-15T10:30:00Z"),
		CreatedBy: AgentList{
			&Person{
				Agent: Agent{
					Element: Element{
						Name:                "John Doe",
						ExternalIdentifiers: externalIdentifierListEmail("john@example.com"),
					},
				},
			},
		},
		CreatedUsing: ToolList{
			&Tool{
				Element: Element{
					Name: "tools-golang-v1.1.0",
				},
			},
		},

		//LicenseListVersion: "3.19",
		Comment: "Created during automated build process",
	}
}

func v301package1() (*Package, ElementList) {
	p := &Package{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID:          "SPDXRef-Package-ExampleLib",
					Name:        "example-library",
					Description: "This is a detailed description of the example library package.",
					Summary:     "A sample library for demonstration",
					Comment:     "Package built with standard configuration" + ";License determined from LICENSE file",

					VerifiedUsing: IntegrityMethodList{
						&PackageVerificationCode{
							Algorithm:     HashAlgorithm_Sha1,
							HashValue:     "d6a770ba38583ed4bb4525bd96e50461655d2758",
							ExcludedFiles: []string{"./exclude1.txt", "./exclude2.txt"},
						},
						&Hash{
							Algorithm: HashAlgorithm_Sha1,
							Value:     "aabbccdd",
						},
						&Hash{
							Algorithm: HashAlgorithm_Sha256,
							Value:     "11223344",
						},
					},
					ExternalIdentifiers: ExternalIdentifierList{
						&ExternalIdentifier{
							Type:       ExternalIdentifierType_PackageURL,
							Identifier: "pkg:npm/example-library@1.2.3",
							Comment:    "NPM package reference",
						},
						&ExternalIdentifier{
							Type:       ExternalIdentifierType_Cpe23,
							Identifier: "cpe:2.3:a:example:library:1.2.3:*:*:*:*:*:*:*",
							Comment:    "CPE reference for security scanning",
						},
					},
				},
				SuppliedBy: &Organization{
					Agent: Agent{
						Element: Element{
							Name:                "Example Corp",
							ExternalIdentifiers: externalIdentifierListEmail("support@example.com"),
						},
					},
				},

				OriginatedBy: AgentList{
					&Person{
						Agent: Agent{
							Element: Element{
								Name:                "Jane Smith",
								ExternalIdentifiers: externalIdentifierListEmail("jane@example.com"),
							},
						},
					},
				},

				ReleaseTime:    parseTime("2023-01-10T00:00:00Z"),
				BuiltTime:      parseTime("2023-01-15T08:00:00Z"),
				ValidUntilTime: parseTime("2025-01-10T00:00:00Z"),
			},
			CopyrightText: "Copyright 2023 Example Corp",
			AttributionTexts: []string{
				"This package includes code from Project ABC",
				"Special thanks to the open source community",
			},
			//FilesAnalyzed: true,
			PrimaryPurpose: SoftwarePurpose_Library,
		},
		Version:          "1.2.3",
		HomePage:         "https://example.com/library",
		SourceInfo:       "Built from git tag v1.2.3",
		DownloadLocation: "https://github.com/example/library/archive/v1.2.3.tar.gz",
		PackageURL:       "pkg:npm/example-library-main@1.2.3",
	}

	// comment was appended to pkg comment: "License determined from LICENSE file"

	l := &ListedLicense{}
	l.Name = "MIT"

	l2 := &ListedLicense{}
	l2.Name = "Apache-2.0"

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l},
		From: p,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{l},
		From: p,
	}

	r3 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l2},
		From: p,
	}

	return p, ElementList{l, r1, r2, l2, r3}
}

func v301package2() (*Package, ElementList) {
	p := &Package{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID:          "SPDXRef-Package-UtilityTools",
					Name:        "utility-tools",
					Description: "A comprehensive set of utility tools for developers.",
					Summary:     "Collection of utility tools",

					VerifiedUsing: IntegrityMethodList{
						&Hash{
							Algorithm: HashAlgorithm_Sha256,
							Value:     "ffaabbcc11223344",
						},
					},
				},
				SuppliedBy: &Person{
					Agent: Agent{
						Element: Element{
							Name:                "Bob Johnson",
							ExternalIdentifiers: externalIdentifierListEmail("bob@tools.com"),
						},
					},
				},
			},
			CopyrightText:  "Copyright 2023 Tools Inc",
			PrimaryPurpose: SoftwarePurpose_Application,
		},
		Version:          "2.1.0",
		HomePage:         ld.URI("https://tools.com/utility"),
		DownloadLocation: ld.URI("https://tools.com/download/utility-tools-2.1.0.zip"),
	}

	l := &ListedLicense{}
	l.Name = "Apache-2.0"

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l},
		From: p,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{l},
		From: p,
	}

	return p, ElementList{l, r1, r2}
}

func v301customLicense1() AnyLicense {
	return &CustomLicense{
		License: License{
			ExtendableLicense: ExtendableLicense{
				LicenseInfo: LicenseInfo{
					Element: Element{
						ID:      "LicenseRef-CustomLicense1",
						Name:    "Custom Example License",
						Comment: "Custom license used for internal tools",
					},
				},
			},
			SeeAlsos: []ld.URI{
				"https://example.com/licenses/custom",
			},
			Text: "This is a custom license text for demonstration purposes.\n\nPermission is granted to use this software...",
		},
	}
}

func v301customLicense2() AnyLicense {
	return &CustomLicense{
		License: License{
			ExtendableLicense: ExtendableLicense{
				LicenseInfo: LicenseInfo{
					Element: Element{
						ID:      "LicenseRef-CustomLicense2",
						Name:    "Another Custom License",
						Comment: "License for third-party components",
					},
				},
			},
			SeeAlsos: []ld.URI{
				"https://example.com/licenses/another-custom",
				"https://internal.example.com/legal/licenses",
			},
			Text: "Another custom license text with different terms.\n\nThis software may be used under the following conditions...",
		},
	}
}

func v301file1() (AnyFile, ElementList) {
	f := &File{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID:   "SPDXRef-File-Other",
					Name: "./src/other.c",
					VerifiedUsing: IntegrityMethodList{
						&Hash{
							Value:     "da39a3ee5e6b4b0d3255bfef95601890afd80709",
							Algorithm: HashAlgorithm_Sha1,
						},
						&Hash{
							Value:     "d41d8cd98f00b204e9800998ecf8427e",
							Algorithm: HashAlgorithm_Md5,
						},
					},
					Comment:     "Other application entry point",
					Description: "This file contains the other function",
				},
				OriginatedBy: AgentList{
					&Person{
						Agent: Agent{
							Element: Element{
								Name:                "Other John Doe",
								ExternalIdentifiers: externalIdentifierListEmail("john@doe.com"),
							},
						},
					},
					&Person{
						Agent: Agent{
							Element: Element{
								Name:                "Jane Smith",
								ExternalIdentifiers: externalIdentifierListEmail("jane.smith@example.org"),
							},
						},
					},
				},
			},

			CopyrightText: "Copyright 2023 Example Corp",
			AttributionTexts: []string{
				"Based on other example code from something",
			},
			//FileTypes: []FileType{
			//	FileType_Source,
			//},
			//PrimaryPurpose: SoftwarePurpose_Source,
		},
		Kind: FileKindType_File,
	}

	l := &ListedLicense{}
	l.Name = "MIT"

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l},
		From: f,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{l},
		From: f,
	}

	return f, ElementList{l, r1, r2}
}

func v301file2() (*File, ElementList) {
	f := &File{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID: "SPDXRef-File-Header",

					Name: "./include/header.h",
					VerifiedUsing: IntegrityMethodList{
						&Hash{
							Value:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
							Algorithm: HashAlgorithm_Sha256,
						},
					},
					Comment: "Header file with function declarations",
				},
			},
			CopyrightText: "Copyright 2023 Example Corp",
			//FileTypes: []FileType{
			//	FileType_Header,
			//},
		},
	}

	l := &ListedLicense{}
	l.Name = "MIT"

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l},
		From: f,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{l},
		From: f,
	}

	return f, ElementList{l, r1, r2}
}

func v301externalMap1() *ExternalMap {
	return &ExternalMap{
		ExternalSpdxID: "DocumentRef-External1",
		LocationHint:   "https://external.com/spdx/external-doc-1.0.0",
		VerifiedUsing: IntegrityMethodList{
			&Hash{
				Algorithm: HashAlgorithm_Sha1,
				Value:     "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			},
		},
	}
}

func v301externalMap2() *ExternalMap {
	return &ExternalMap{
		ExternalSpdxID: "DocumentRef-External2",
		LocationHint:   "https://external2.com/spdx/external-doc-2.0.0",
		VerifiedUsing: IntegrityMethodList{
			&Hash{
				Algorithm: HashAlgorithm_Sha256,
				Value:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
	}
}

func v301annotation1(subject AnyElement) *Annotation {
	return &Annotation{
		Type: AnnotationType_Review,
		Element: Element{
			ID: "SPDXRef-Annotation-1",
			CreationInfo: &CreationInfo{
				Created: parseTime("2023-01-20T14:30:00Z"),
				CreatedBy: AgentList{
					&Person{
						Agent: Agent{
							Element: Element{
								Name:                "Security Team",
								ExternalIdentifiers: externalIdentifierListEmail("security@example.com"),
							},
						},
					},
				},
			},
		},
		Subject:   subject,
		Statement: "Security review completed - no vulnerabilities found",
	}
}

func v301annotation2(subject AnyElement) *Annotation {
	return &Annotation{
		Type: AnnotationType_Other,
		Element: Element{
			ID: "SPDXRef-Annotation-2",
			CreationInfo: &CreationInfo{
				Created: parseTime("2023-01-21T09:15:00Z"),
				CreatedBy: AgentList{
					&Person{
						Agent: Agent{
							Element: Element{
								Name: "vulnerability-scanner-v1.5",
							},
						},
					},
				},
			},
			Comment: "Automated scan completed - clean",
		},
		Subject: subject,
	}
}

func v301snippet1(fileRef AnyFile) (AnyElement, ElementList) {
	s := &Snippet{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID:      "SPDXRef-Snippet1",
					Name:    "Core Algorithm",
					Comment: "Key algorithm implementation",
				},
			},
			CopyrightText:    "Copyright 2023 Example Corp",
			AttributionTexts: []string{"Algorithm based on research paper XYZ"},
		},
		FromFile: fileRef,
		ByteRange: &PositiveIntegerRange{
			BeginIntegerRange: 100,
			EndIntegerRange:   200,
		},
		LineRange: &PositiveIntegerRange{
			BeginIntegerRange: 10,
			EndIntegerRange:   15,
		},

		//LicenseComments:  "License applies to this code snippet",
	}

	l := &ListedLicense{}
	l.Name = "MIT"

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l},
		From: s,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{l},
		From: s,
	}

	return s, ElementList{l, r1, r2}
}

func v301snippet2(fileRef AnyFile) (AnySnippet, ElementList) {
	s := &Snippet{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID:      "SPDXRef-Snippet2",
					Name:    "API Declarations",
					Comment: "Function declarations",
				},
			},
			CopyrightText: "Copyright 2023 Example Corp",
		},
		FromFile: fileRef,
		ByteRange: &PositiveIntegerRange{
			EndIntegerRange:   50,
			BeginIntegerRange: 150,
		},
		LineRange: &PositiveIntegerRange{
			EndIntegerRange:   5,
			BeginIntegerRange: 8,
		},
	}

	l := &ListedLicense{}
	l.Name = "MIT"

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{l},
		From: s,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{l},
		From: s,
	}

	return s, ElementList{l, r1, r2}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
