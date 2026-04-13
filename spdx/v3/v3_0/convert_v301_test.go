package v3_0

import (
	"slices"
	"time"

	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/spdx/tools-golang/spdx/v3/internal"
	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
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
		Name: v2_3.DataLicense,
	}
	d.SpdxDocument.ID = "DOCUMENT"
	d.CreationInfo.(*CreationInfo).Created = parseTime("2023-01-15T10:30:00Z")
	d.SpdxDocument.Comment = "This is a sample SPDX document for testing purposes."
	d.NamespaceMaps = NamespaceMapList{
		&NamespaceMap{
			Prefix:    internal.DefaultSpdxNamespace,
			Namespace: "https://example.com/spdx/example-software-1.0.0#",
		},
	}

	sbom := &SBOM{}
	d.RootElements = ElementList{sbom}

	// Add external document references
	d.Imports = ExternalMapList{
		v301externalMap1(),
		v301externalMap2(),
	}

	addToSBOM := func(elements ...AnyElement) {
		for _, e := range elements {
			if !slices.Contains(sbom.Elements, e) {
				sbom.Elements = append(sbom.Elements, e)
			}
		}
	}

	pkg1, pkg1Elems := v301package1()
	pkg2, pkg2Elems := v301package2()
	file1, file1Elems := v301file1()
	file2, file2Elems := v301file2()
	snippet1, snippet1Elems := v301snippet1(file1)
	snippet2, snippet2Elems := v301snippet2(file2)
	ann1 := v301annotation1(pkg1)
	ann2 := v301annotation2(pkg2)

	// pkg1 is the only thing referenced directly as a document root element
	sbom.RootElements = append(sbom.RootElements, pkg1)

	// top-level converted elements go directly in sbom.Elements
	// (matching converter behavior: packages, files, annotations, snippets, custom licenses)
	addToSBOM(
		v301customLicense1(),
		v301customLicense2(),
		pkg1, pkg2,
		AnyElement(file1), AnyElement(file2),
		ann1, ann2,
		AnyElement(snippet1), AnyElement(snippet2),
	)

	// all relationships go in sbom.Elements
	for _, elems := range []ElementList{pkg1Elems, pkg2Elems, file1Elems, file2Elems, snippet1Elems, snippet2Elems} {
		for _, r := range elems.Relationships() {
			addToSBOM(r)
		}
	}
	addToSBOM(
		&Relationship{
			Type: RelationshipType_Describes,
			From: ann1,
			To:   ElementList{pkg1},
		},
		&Relationship{
			Type: RelationshipType_Describes,
			From: ann2,
			To:   ElementList{pkg2},
		},
		// SpdxRef-DOCUMENT relationships are handled by object structure
		&LifecycleScopedRelationship{
			From:    pkg2,
			To:      ElementList{pkg1},
			Type:    RelationshipType_UsesTool,
			Scope:   LifecycleScopeType_Build,
			Comment: "Main package depends on utility tools",
		},
		&Relationship{
			From:    pkg1,
			To:      ElementList{file1},
			Type:    RelationshipType_Contains,
			Comment: "Package contains main source file",
		},
	)

	// collect all reachable elements into the document Elements list
	for _, e := range collectAllElements(&d.SpdxDocument) {
		d.SpdxDocument.Elements = append(d.SpdxDocument.Elements, e)
	}

	return d
}

func v301creationInfo() *CreationInfo {
	return &CreationInfo{
		SpecVersion: Version,
		Created:     parseTime("2023-01-15T10:30:00Z"),
		CreatedBy: AgentList{
			&Person{
				Name:                "John Doe",
				ExternalIdentifiers: externalIdentifierListEmail("john@example.com"),
			},
		},
		CreatedUsing: ToolList{
			&Tool{
				Name: "tools-golang-v1.1.0",
			},
		},

		Comment: "Created during automated build process" +
			"; LicenseListVersion: 3.19",
	}
}

func v301package1() (*Package, ElementList) {
	p := &Package{
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
		SuppliedBy: &Organization{
			Name:                "Example Corp",
			ExternalIdentifiers: externalIdentifierListEmail("support@example.com"),
		},

		OriginatedBy: AgentList{
			&Person{
				Name:                "Jane Smith",
				ExternalIdentifiers: externalIdentifierListEmail("jane@example.com"),
			},
		},

		ReleaseTime:    parseTime("2023-01-10T00:00:00Z"),
		BuiltTime:      parseTime("2023-01-15T08:00:00Z"),
		ValidUntilTime: parseTime("2025-01-10T00:00:00Z"),
		CopyrightText:  "Copyright 2023 Example Corp",
		AttributionTexts: []string{
			"This package includes code from Project ABC",
			"Special thanks to the open source community",
		},
		//FilesAnalyzed: true,
		PrimaryPurpose:   SoftwarePurpose_Library,
		Version:          "1.2.3",
		HomePage:         "https://example.com/library",
		SourceInfo:       "Built from git tag v1.2.3",
		DownloadLocation: "https://github.com/example/library/archive/v1.2.3.tar.gz",
		PackageURL:       "pkg:npm/example-library-main@1.2.3",
	}

	// comment was appended to pkg comment: "License determined from LICENSE file"

	// converter creates separate license objects for each expression
	lInfoMIT := &ListedLicense{Name: "MIT"}           // from LicenseInfoFromFiles
	lInfoApache := &ListedLicense{Name: "Apache-2.0"} // from LicenseInfoFromFiles
	lConcluded := &ListedLicense{Name: "MIT"}         // from LicenseConcluded
	lDeclared := &ListedLicense{Name: "MIT"}          // from LicenseDeclared

	r1 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{lInfoMIT, lInfoApache, lDeclared},
		From: p,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{lConcluded},
		From: p,
	}

	// file associated with this package (mirrors v23package1 embedded file)
	file3, fileElems := v301file3()

	// annotation associated with this package (mirrors v23package1 embedded annotation)
	annotation3 := v301annotation3(p)

	containsFile := &Relationship{
		Type: RelationshipType_Contains,
		From: p,
		To:   ElementList{file3},
	}
	describesPkg := &Relationship{
		Type: RelationshipType_Describes,
		From: annotation3,
		To:   ElementList{p},
	}

	extras := ElementList{lInfoMIT, lInfoApache, lConcluded, lDeclared, r1, r2, file3, annotation3, containsFile, describesPkg}
	extras = append(extras, fileElems...)

	return p, extras
}

func v301package2() (*Package, ElementList) {
	p := &Package{
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
		SuppliedBy: &Person{
			Name:                "Bob Johnson",
			ExternalIdentifiers: externalIdentifierListEmail("bob@tools.com"),
		},
		CopyrightText:    "Copyright 2023 Tools Inc",
		PrimaryPurpose:   SoftwarePurpose_Application,
		Version:          "2.1.0",
		HomePage:         ld.URI("https://tools.com/utility"),
		DownloadLocation: ld.URI("https://tools.com/download/utility-tools-2.1.0.zip"),
	}

	concludedLicense := &ListedLicense{Name: "Apache-2.0"}

	declaredLicense := &DisjunctiveLicenseSet{
		Members: LicenseInfoList{
			&ListedLicense{Name: "AGPL"},
			&DisjunctiveLicenseSet{
				Members: LicenseInfoList{
					&WithAdditionOperator{
						SubjectExtendableLicense: &ListedLicense{Name: "GPL-2.0-only"},
						SubjectAddition:          &ListedLicenseException{Name: "Classpath-Exception"},
					},
					&ListedLicense{Name: "LicenceRef-CustomLicense1"},
				},
			},
		},
	}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{concludedLicense},
		From: p,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{declaredLicense},
		From: p,
	}

	return p, ElementList{concludedLicense, r1, declaredLicense, r2}
}

func v301customLicense1() AnyLicense {
	return &CustomLicense{
		ID:      "LicenseRef-CustomLicense1",
		Name:    "Custom Example License",
		Comment: "Custom license used for internal tools",
		SeeAlsos: []ld.URI{
			"https://example.com/licenses/custom",
		},
		Text: "This is a custom license text for demonstration purposes.\n\nPermission is granted to use this software...",
	}
}

func v301customLicense2() AnyLicense {
	return &CustomLicense{
		ID:      "LicenseRef-CustomLicense2",
		Name:    "Another Custom License",
		Comment: "License for third-party components",
		SeeAlsos: []ld.URI{
			"https://example.com/licenses/another-custom",
			"https://internal.example.com/legal/licenses",
		},
		Text: "Another custom license text with different terms.\n\nThis software may be used under the following conditions...",
	}
}

func v301file1() (AnyFile, ElementList) {
	f := &File{
		ID:   "SPDXRef-File-Main",
		Name: "./src/main.c",
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
		Comment: "Other application entry point",
		OriginatedBy: AgentList{
			&Person{
				Name:                "Other John Doe",
				ExternalIdentifiers: externalIdentifierListEmail("john@doe.com"),
			},
			&Person{
				Name:                "Jane Smith",
				ExternalIdentifiers: externalIdentifierListEmail("jane.smith@example.org"),
			},
		},

		CopyrightText: "Copyright 2023 Example Corp",
		AttributionTexts: []string{
			"Based on other example code from something",
		},
		Kind: FileKindType_File,
	}

	lConcluded := &ListedLicense{Name: "MIT", Comment: "License applies to this file"}
	lDeclared := &ListedLicense{Name: "MIT"}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{lConcluded},
		From: f,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{lDeclared},
		From: f,
	}

	// snippet associated with this file (mirrors v23file1 embedded snippet)
	snippet3, snippetElems := v301snippet3(f)

	// annotation associated with this file (mirrors v23file1 embedded annotation)
	annotation4 := v301annotation4(f)

	containsSnippet := &Relationship{
		Type: RelationshipType_Contains,
		From: f,
		To:   ElementList{snippet3},
	}
	describesFile := &Relationship{
		Type: RelationshipType_Describes,
		From: annotation4,
		To:   ElementList{f},
	}

	// from FileNotice
	fileNotice := &Annotation{
		Type:      AnnotationType_Other,
		Subject:   f,
		Statement: "This file contains the other function",
	}
	fileNoticeRel := &Relationship{
		Type: RelationshipType_Describes,
		From: fileNotice,
		To:   ElementList{f},
	}

	extras := ElementList{lConcluded, lDeclared, r1, r2, snippet3, annotation4, containsSnippet, describesFile, fileNotice, fileNoticeRel}
	extras = append(extras, snippetElems...)

	return f, extras
}

func v301file2() (*File, ElementList) {
	f := &File{
		ID: "SPDXRef-File-Header",

		Name: "./include/header.h",
		VerifiedUsing: IntegrityMethodList{
			&Hash{
				Value:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Algorithm: HashAlgorithm_Sha256,
			},
		},
		Comment:       "Header file with function declarations",
		CopyrightText: "Copyright 2023 Example Corp",
		//FileTypes: []FileType{
		//	FileType_Header,
		//},
		Kind: FileKindType_File,
	}

	concludedLicense := &ConjunctiveLicenseSet{
		Members: LicenseInfoList{
			&ListedLicense{Name: "MIT"},
			&ListedLicense{Name: "Apache-2.0"},
		},
	}

	declaredLicense := &ListedLicense{Name: "MIT"}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{concludedLicense},
		From: f,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{declaredLicense},
		From: f,
	}

	return f, ElementList{concludedLicense, declaredLicense, r1, r2}
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
		ID:   "SPDXRef-Annotation-1",
		CreationInfo: &CreationInfo{
			Created: parseTime("2023-01-20T14:30:00Z"),
			CreatedBy: AgentList{
				&Person{
					Name:                "Security Team",
					ExternalIdentifiers: externalIdentifierListEmail("security@example.com"),
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
		ID:   "SPDXRef-Annotation-2",
		CreationInfo: &CreationInfo{
			Created: parseTime("2023-01-21T09:15:00Z"),
			CreatedBy: AgentList{
				&SoftwareAgent{
					Name: "vulnerability-scanner-v1.5",
				},
			},
		},
		Statement: "Automated scan completed - clean",
		Subject:   subject,
	}
}

func v301file3() (AnyFile, ElementList) {
	f := &File{
		ID:             "SPDXRef-File-Utils",
		Name:           "./src/utils.c",
		PrimaryPurpose: SoftwarePurpose_Source,
		VerifiedUsing: IntegrityMethodList{
			&Hash{
				Value:     "aabb1234aabb1234aabb1234aabb1234aabb1234aabb1234aabb1234aabb1234",
				Algorithm: HashAlgorithm_Sha256,
			},
		},
		Comment:       "Utility functions for the library",
		CopyrightText: "Copyright 2023 Example Corp",
		Kind:          FileKindType_File,
	}

	lConcluded := &ListedLicense{Name: "Apache-2.0"}
	lDeclared := &ListedLicense{Name: "Apache-2.0"}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{lConcluded},
		From: f,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{lDeclared},
		From: f,
	}

	return f, ElementList{lConcluded, lDeclared, r1, r2}
}

func v301annotation3(subject AnyElement) *Annotation {
	return &Annotation{
		Type: AnnotationType_Other,
		ID:   "SPDXRef-Annotation-3",
		CreationInfo: &CreationInfo{
			Created: parseTime("2023-02-01T10:00:00Z"),
			CreatedBy: AgentList{
				&Person{
					Name:                "Code Review Bot",
					ExternalIdentifiers: externalIdentifierListEmail("bot@example.com"),
				},
			},
		},
		Subject:   subject,
		Statement: "Package structure review completed",
	}
}

func v301annotation4(subject AnyElement) *Annotation {
	return &Annotation{
		Type: AnnotationType_Review,
		ID:   "SPDXRef-Annotation-4",
		CreationInfo: &CreationInfo{
			Created: parseTime("2023-02-05T16:00:00Z"),
			CreatedBy: AgentList{
				&Person{
					Name:                "Code Reviewer",
					ExternalIdentifiers: externalIdentifierListEmail("reviewer@example.com"),
				},
			},
		},
		Subject:   subject,
		Statement: "File review completed - code quality approved",
	}
}

func v301snippet3(fileRef AnyFile) (AnyElement, ElementList) {
	s := &Snippet{
		ID:               "SPDXRef-Snippet3",
		Name:             "Helper Routine",
		Comment:          "Helper routine implementation",
		CopyrightText:    "Copyright 2023 Example Corp",
		AttributionTexts: []string{"Inspired by open source utilities"},
		FromFile:         fileRef,
		ByteRange: &PositiveIntegerRange{
			BeginIntegerRange: 300,
			EndIntegerRange:   400,
		},
		LineRange: &PositiveIntegerRange{
			BeginIntegerRange: 20,
			EndIntegerRange:   30,
		},
	}

	lConcluded := &ListedLicense{Name: "Apache-2.0", Comment: "License for helper routine"}
	lDeclared := &ListedLicense{Name: "Apache-2.0", Comment: ""}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{lConcluded},
		From: s,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{lDeclared},
		From: s,
	}

	return s, ElementList{lConcluded, lDeclared, r1, r2}
}

func v301snippet1(fileRef AnyFile) (AnyElement, ElementList) {
	s := &Snippet{
		ID:               "SPDXRef-Snippet1",
		Name:             "Core Algorithm",
		Comment:          "Key algorithm implementation",
		CopyrightText:    "Copyright 2023 Example Corp",
		AttributionTexts: []string{"Algorithm based on research paper XYZ"},
		FromFile:         fileRef,
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

	lConcluded := &ListedLicense{Name: "MIT", Comment: "License applies to this code snippet"}
	lDeclared := &ListedLicense{Name: "MIT", Comment: ""}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{lConcluded},
		From: s,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{lDeclared},
		From: s,
	}

	return s, ElementList{lConcluded, lDeclared, r1, r2}
}

func v301snippet2(fileRef AnyFile) (AnySnippet, ElementList) {
	s := &Snippet{
		ID:            "SPDXRef-Snippet2",
		Name:          "API Declarations",
		Comment:       "Function declarations",
		CopyrightText: "Copyright 2023 Example Corp",
		FromFile:      fileRef,
		ByteRange: &PositiveIntegerRange{
			BeginIntegerRange: 50,
			EndIntegerRange:   150,
		},
		LineRange: &PositiveIntegerRange{
			BeginIntegerRange: 5,
			EndIntegerRange:   8,
		},
	}

	concludedLicense := &OrLaterOperator{
		SubjectLicense: &ListedLicense{Name: "MIT"},
	}

	declaredLicense := &ListedLicense{Name: "MIT"}

	r1 := &Relationship{
		Type: RelationshipType_HasConcludedLicense,
		To:   ElementList{concludedLicense},
		From: s,
	}

	r2 := &Relationship{
		Type: RelationshipType_HasDeclaredLicense,
		To:   ElementList{declaredLicense},
		From: s,
	}

	return s, ElementList{concludedLicense, declaredLicense, r1, r2}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
