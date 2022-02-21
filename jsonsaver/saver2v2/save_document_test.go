// Package saver2v2 contains functions to render and write a json
// formatted version of an in-memory SPDX document and its sections
// (version 2.2).
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package saver2v2

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestRenderDocument2_2(t *testing.T) {

	test1 := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			DataLicense:          "CC0-1.0",
			SPDXVersion:          "SPDX-2.2",
			SPDXIdentifier:       "DOCUMENT",
			DocumentComment:      "This document was created using SPDX 2.0 using licenses from the web site.",
			LicenseListVersion:   "3.8",
			Created:              "2010-01-29T18:30:22Z",
			CreatorPersons:       []string{"Jane Doe ()"},
			CreatorOrganizations: []string{"ExampleCodeInspect ()"},
			CreatorTools:         []string{"LicenseFind-1.0"},
			DocumentName:         "SPDX-Tools-v2.0",
			DocumentNamespace:    "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301",
			CreatorComment:       "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
			ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{
				"spdx-tool-1.2": {
					DocumentRefID: "spdx-tool-1.2",
					URI:           "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
					Alg:           "SHA1",
					Checksum:      "d6a770ba38583ed4bb4525bd96e50461655d2759",
				},
			},
		},
		OtherLicenses: []*spdx.OtherLicense2_2{
			{
				ExtractedText:     "\"THE BEER-WARE LICENSE\" (Revision 42):\nphk@FreeBSD.ORG wrote this file. As long as you retain this notice you\ncan do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return Poul-Henning Kamp  </\nLicenseName: Beer-Ware License (Version 42)\nLicenseCrossReference:  http://people.freebsd.org/~phk/\nLicenseComment: \nThe beerware license has a couple of other standard variants.",
				LicenseIdentifier: "LicenseRef-Beerware-4.2",
			},
			{
				LicenseComment:         "This is tye CyperNeko License",
				ExtractedText:          "The CyberNeko Software License, Version 1.0\n\n \n(C) Copyright 2002-2005, Andy Clark.  All rights reserved.\n \nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions\nare met:\n\n1. Redistributions of source code must retain the above copyright\n   notice, this list of conditions and the following disclaimer. \n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in\n   the documentation and/or other materials provided with the\n   distribution.\n\n3. The end-user documentation included with the redistribution,\n   if any, must include the following acknowledgment:  \n     \"This product includes software developed by Andy Clark.\"\n   Alternately, this acknowledgment may appear in the software itself,\n   if and wherever such third-party acknowledgments normally appear.\n\n4. The names \"CyberNeko\" and \"NekoHTML\" must not be used to endorse\n   or promote products derived from this software without prior \n   written permission. For written permission, please contact \n   andyc@cyberneko.net.\n\n5. Products derived from this software may not be called \"CyberNeko\",\n   nor may \"CyberNeko\" appear in their name, without prior written\n   permission of the author.\n\nTHIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESSED OR IMPLIED\nWARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\nOF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR OTHER CONTRIBUTORS\nBE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, \nOR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT \nOF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR \nBUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, \nWHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE \nOR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, \nEVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
				LicenseIdentifier:      "LicenseRef-3",
				LicenseName:            "CyberNeko License",
				LicenseCrossReferences: []string{"http://people.apache.org/~andyc/neko/LICENSE", "http://justasample.url.com"},
			},
		},
		Annotations: []*spdx.Annotation2_2{
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
				AnnotationDate:           "2011-01-29T18:30:22Z",
				AnnotationType:           "OTHER",
				AnnotatorType:            "Person",
				Annotator:                "File Commenter",
				AnnotationComment:        "File level annotation",
			},
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
				AnnotationDate:           "2011-01-29T18:30:22Z",
				AnnotationType:           "OTHER",
				AnnotatorType:            "Person",
				Annotator:                "Package Commenter",
				AnnotationComment:        "Package level annotation",
			},
			{
				AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				AnnotationDate:           "2010-02-10T00:00:00Z",
				AnnotationType:           "REVIEW",
				AnnotatorType:            "Person",
				Annotator:                "Joe Reviewer",
				AnnotationComment:        "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
			},
		},
		Relationships: []*spdx.Relationship2_2{
			{
				RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				RefB:         spdx.DocElementID{DocumentRefID: "spdx-tool-1.2", ElementRefID: "ToolsElement"},
				Relationship: "COPY_OF",
			},
			{
				RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
				Relationship: "CONTAINS",
			},
			{
				RefA:                spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
				RefB:                spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
				Relationship:        "DESCRIBES",
				RelationshipComment: "This is a comment.",
			},
		},
		Packages: map[spdx.ElementID]*spdx.Package2_2{
			"Package": {
				PackageSPDXIdentifier:   "Package",
				PackageAttributionTexts: []string{"The GNU C Library is free software.  See the file COPYING.LIB for copying conditions, and LICENSES for notices about a few contributions that require these additional notices to be distributed.  License copyright years may be listed using range notation, e.g., 1996-2015, indicating that every year in the range, inclusive, is a copyrightable year that would otherwise be listed individually."},
				PackageChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
					"MD5": {
						Algorithm: "MD5",
						Value:     "624c1abb3664f4b35547e7c73864ad24",
					},
				},
				PackageCopyrightText:    "Copyright 2008-2010 John Smith",
				PackageDescription:      "The GNU C Library defines functions that are specified by the ISO C standard, as well as additional features specific to POSIX and other derivatives of the Unix operating system, and extensions specific to GNU systems.",
				PackageDownloadLocation: "http://ftp.gnu.org/gnu/glibc/glibc-ports-2.15.tar.gz",
				PackageExternalReferences: []*spdx.PackageExternalReference2_2{
					{
						RefType:            "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#LocationRef-acmeforge",
						ExternalRefComment: "This is the external ref for Acme",
						Category:           "OTHER",
						Locator:            "acmecorp/acmenator/4.1.3-alpha",
					},
					{
						RefType:  "http://spdx.org/rdf/references/cpe23Type",
						Category: "SECURITY",
						Locator:  "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
					},
				},
				FilesAnalyzed:             true,
				IsFilesAnalyzedTagPresent: true,
				Files: map[spdx.ElementID]*spdx.File2_2{
					"DoapSource": {
						FileSPDXIdentifier: "DoapSource",
						FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
							"SHA1": {
								Algorithm: "SHA1",
								Value:     "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
							},
						},
						Snippets: map[spdx.ElementID]*spdx.Snippet2_2{
							"Snippet": {
								SnippetSPDXIdentifier:         "Snippet",
								SnippetFromFileSPDXIdentifier: spdx.DocElementID{ElementRefID: "DoapSource"},
								SnippetComment:                "This snippet was identified as significant and highlighted in this Apache-2.0 file, when a commercial scanner identified it as being derived from file foo.c in package xyz which is licensed under GPL-2.0.",
								SnippetCopyrightText:          "Copyright 2008-2010 John Smith",
								SnippetLicenseComments:        "The concluded license was taken from package xyz, from which the snippet was copied into the current file. The concluded license information was found in the COPYING.txt file in package xyz.",
								SnippetLicenseConcluded:       "GPL-2.0-only",
								LicenseInfoInSnippet:          []string{"GPL-2.0-only"},
								SnippetName:                   "from linux kernel",
								SnippetByteRangeStart:         310,
								SnippetByteRangeEnd:           420,
								SnippetLineRangeStart:         5,
								SnippetLineRangeEnd:           23,
							},
						},
						FileCopyrightText: "Copyright 2010, 2011 Source Auditor Inc.",
						FileContributor:   []string{"Protecode Inc.", "SPDX Technical Team Members", "Open Logic Inc.", "Source Auditor Inc.", "Black Duck Software In.c"},
						FileDependencies:  []string{"SPDXRef-JenaLib", "SPDXRef-CommonsLangSrc"},
						FileName:          "./src/org/spdx/parser/DOAPProject.java",
						FileType:          []string{"SOURCE"},
						LicenseConcluded:  "Apache-2.0",
						LicenseInfoInFile: []string{"Apache-2.0"},
					},
				},
				PackageHomePage:                     "http://ftp.gnu.org/gnu/glibc",
				PackageLicenseComments:              "The license for this project changed with the release of version x.y.  The version of the project included here post-dates the license change.",
				PackageLicenseConcluded:             "(LGPL-2.0-only OR LicenseRef-3)",
				PackageLicenseDeclared:              "(LGPL-2.0-only AND LicenseRef-3)",
				PackageLicenseInfoFromFiles:         []string{"GPL-2.0-only", "LicenseRef-2", "LicenseRef-1"},
				PackageName:                         "glibc",
				PackageOriginatorOrganization:       "ExampleCodeInspect (contact@example.com)",
				PackageFileName:                     "glibc-2.11.1.tar.gz",
				PackageVerificationCodeExcludedFile: "./package.spdx",
				PackageVerificationCode:             "d6a770ba38583ed4bb4525bd96e50461655d2758",
				PackageSourceInfo:                   "uses glibc-2_11-branch from git://sourceware.org/git/glibc.git.",
				PackageSummary:                      "GNU C library.",
				PackageSupplierPerson:               "Jane Doe (jane.doe@example.com)",
				PackageVersion:                      "2.11.1",
			},
		},
		UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{
			"File": {
				FileSPDXIdentifier: "File",
				FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
					"MD5": {
						Algorithm: "MD5",
						Value:     "624c1abb3664f4b35547e7c73864ad24",
					},
					"SHA1": {
						Algorithm: "SHA1",
						Value:     "d6a770ba38583ed4bb4525bd96e50461655d2758",
					},
				},
				FileComment:       "The concluded license was taken from the package level that the file was .\nThis information was found in the COPYING.txt file in the xyz directory.",
				FileCopyrightText: "Copyright 2008-2010 John Smith",
				FileContributor:   []string{"The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation"},
				FileName:          "./package/foo.c",
				FileType:          []string{"SOURCE"},
				LicenseComments:   "The concluded license was taken from the package level that the file was .\nThis information was found in the COPYING.txt file in the xyz directory.",
				LicenseConcluded:  "(LGPL-2.0-only OR LicenseRef-2)",
				LicenseInfoInFile: []string{"GPL-2.0-only", "LicenseRef-2"},
				FileNotice:        "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.",
			},
		},
	}

	want := map[string]interface{}{
		"dataLicense":       "CC0-1.0",
		"spdxVersion":       "SPDX-2.2",
		"SPDXID":            "SPDXRef-DOCUMENT",
		"documentNamespace": "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301",
		"documentDescribes": []string{"SPDXRef-Package"},
		"name":              "SPDX-Tools-v2.0",
		"comment":           "This document was created using SPDX 2.0 using licenses from the web site.",
		"creationInfo": map[string]interface{}{
			"comment":            "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
			"created":            "2010-01-29T18:30:22Z",
			"creators":           []string{"Tool: LicenseFind-1.0", "Organization: ExampleCodeInspect ()", "Person: Jane Doe ()"},
			"licenseListVersion": "3.8",
		},
		"externalDocumentRefs": []interface{}{
			map[string]interface{}{
				"externalDocumentId": "DocumentRef-spdx-tool-1.2",
				"spdxDocument":       "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
				"checksum": map[string]string{
					"algorithm":     "SHA1",
					"checksumValue": "d6a770ba38583ed4bb4525bd96e50461655d2759",
				},
			},
		},
		"hasExtractedLicensingInfos": []interface{}{
			map[string]interface{}{
				"licenseId":     "LicenseRef-Beerware-4.2",
				"extractedText": "\"THE BEER-WARE LICENSE\" (Revision 42):\nphk@FreeBSD.ORG wrote this file. As long as you retain this notice you\ncan do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return Poul-Henning Kamp  </\nLicenseName: Beer-Ware License (Version 42)\nLicenseCrossReference:  http://people.freebsd.org/~phk/\nLicenseComment: \nThe beerware license has a couple of other standard variants.",
			},
			map[string]interface{}{
				"licenseId":     "LicenseRef-3",
				"comment":       "This is tye CyperNeko License",
				"extractedText": "The CyberNeko Software License, Version 1.0\n\n \n(C) Copyright 2002-2005, Andy Clark.  All rights reserved.\n \nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions\nare met:\n\n1. Redistributions of source code must retain the above copyright\n   notice, this list of conditions and the following disclaimer. \n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in\n   the documentation and/or other materials provided with the\n   distribution.\n\n3. The end-user documentation included with the redistribution,\n   if any, must include the following acknowledgment:  \n     \"This product includes software developed by Andy Clark.\"\n   Alternately, this acknowledgment may appear in the software itself,\n   if and wherever such third-party acknowledgments normally appear.\n\n4. The names \"CyberNeko\" and \"NekoHTML\" must not be used to endorse\n   or promote products derived from this software without prior \n   written permission. For written permission, please contact \n   andyc@cyberneko.net.\n\n5. Products derived from this software may not be called \"CyberNeko\",\n   nor may \"CyberNeko\" appear in their name, without prior written\n   permission of the author.\n\nTHIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESSED OR IMPLIED\nWARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\nOF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR OTHER CONTRIBUTORS\nBE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, \nOR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT \nOF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR \nBUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, \nWHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE \nOR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, \nEVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
				"name":          "CyberNeko License",
				"seeAlsos":      []string{"http://people.apache.org/~andyc/neko/LICENSE", "http://justasample.url.com"},
			},
		},
		"annotations": []interface{}{
			map[string]interface{}{
				"annotationDate": "2010-02-10T00:00:00Z",
				"annotationType": "REVIEW",
				"annotator":      "Person: Joe Reviewer",
				"comment":        "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
			},
		},
		"relationships": []interface{}{
			map[string]interface{}{
				"spdxElementId":      "SPDXRef-DOCUMENT",
				"relatedSpdxElement": "DocumentRef-spdx-tool-1.2:SPDXRef-ToolsElement",
				"relationshipType":   "COPY_OF",
			},
			map[string]interface{}{
				"spdxElementId":      "SPDXRef-DOCUMENT",
				"relatedSpdxElement": "SPDXRef-Package",
				"relationshipType":   "CONTAINS",
			},
			map[string]interface{}{
				"spdxElementId":      "SPDXRef-DOCUMENT",
				"relatedSpdxElement": "SPDXRef-File",
				"relationshipType":   "DESCRIBES",
				"comment":            "This is a comment.",
			},
		},
		"packages": []interface{}{
			map[string]interface{}{
				"SPDXID": "SPDXRef-Package",
				"annotations": []interface{}{
					map[string]interface{}{
						"annotationDate": "2011-01-29T18:30:22Z",
						"annotationType": "OTHER",
						"annotator":      "Person: Package Commenter",
						"comment":        "Package level annotation",
					},
				},
				"attributionTexts": []string{"The GNU C Library is free software.  See the file COPYING.LIB for copying conditions, and LICENSES for notices about a few contributions that require these additional notices to be distributed.  License copyright years may be listed using range notation, e.g., 1996-2015, indicating that every year in the range, inclusive, is a copyrightable year that would otherwise be listed individually."},
				"checksums": []interface{}{
					map[string]interface{}{
						"algorithm":     "MD5",
						"checksumValue": "624c1abb3664f4b35547e7c73864ad24",
					},
				},
				"copyrightText":    "Copyright 2008-2010 John Smith",
				"description":      "The GNU C Library defines functions that are specified by the ISO C standard, as well as additional features specific to POSIX and other derivatives of the Unix operating system, and extensions specific to GNU systems.",
				"downloadLocation": "http://ftp.gnu.org/gnu/glibc/glibc-ports-2.15.tar.gz",
				"externalRefs": []interface{}{
					map[string]interface{}{
						"comment":           "This is the external ref for Acme",
						"referenceCategory": "OTHER",
						"referenceLocator":  "acmecorp/acmenator/4.1.3-alpha",
						"referenceType":     "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#LocationRef-acmeforge",
					},
					map[string]interface{}{
						"referenceCategory": "SECURITY",
						"referenceLocator":  "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
						"referenceType":     "http://spdx.org/rdf/references/cpe23Type",
					},
				},
				"filesAnalyzed":        true,
				"hasFiles":             []string{"SPDXRef-DoapSource"},
				"homepage":             "http://ftp.gnu.org/gnu/glibc",
				"licenseComments":      "The license for this project changed with the release of version x.y.  The version of the project included here post-dates the license change.",
				"licenseConcluded":     "(LGPL-2.0-only OR LicenseRef-3)",
				"licenseDeclared":      "(LGPL-2.0-only AND LicenseRef-3)",
				"licenseInfoFromFiles": []string{"GPL-2.0-only", "LicenseRef-2", "LicenseRef-1"},
				"name":                 "glibc",
				"originator":           "Organization: ExampleCodeInspect (contact@example.com)",
				"packageFileName":      "glibc-2.11.1.tar.gz",
				"packageVerificationCode": map[string]interface{}{
					"packageVerificationCodeExcludedFiles": []string{"./package.spdx"},
					"packageVerificationCodeValue":         "d6a770ba38583ed4bb4525bd96e50461655d2758",
				},
				"sourceInfo":  "uses glibc-2_11-branch from git://sourceware.org/git/glibc.git.",
				"summary":     "GNU C library.",
				"supplier":    "Person: Jane Doe (jane.doe@example.com)",
				"versionInfo": "2.11.1",
			},
		},
		"files": []interface{}{
			map[string]interface{}{
				"SPDXID": "SPDXRef-DoapSource",
				"checksums": []interface{}{
					map[string]interface{}{
						"algorithm":     "SHA1",
						"checksumValue": "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
					},
				},
				"copyrightText":      "Copyright 2010, 2011 Source Auditor Inc.",
				"fileContributors":   []string{"Protecode Inc.", "SPDX Technical Team Members", "Open Logic Inc.", "Source Auditor Inc.", "Black Duck Software In.c"},
				"fileDependencies":   []string{"SPDXRef-JenaLib", "SPDXRef-CommonsLangSrc"},
				"fileName":           "./src/org/spdx/parser/DOAPProject.java",
				"fileTypes":          []string{"SOURCE"},
				"licenseConcluded":   "Apache-2.0",
				"licenseInfoInFiles": []string{"Apache-2.0"},
			},
			map[string]interface{}{
				"SPDXID": "SPDXRef-File",
				"annotations": []interface{}{
					map[string]interface{}{
						"annotationDate": "2011-01-29T18:30:22Z",
						"annotationType": "OTHER",
						"annotator":      "Person: File Commenter",
						"comment":        "File level annotation",
					},
				},
				"checksums": []interface{}{
					map[string]interface{}{
						"algorithm":     "MD5",
						"checksumValue": "624c1abb3664f4b35547e7c73864ad24",
					},

					map[string]interface{}{
						"algorithm":     "SHA1",
						"checksumValue": "d6a770ba38583ed4bb4525bd96e50461655d2758",
					},
				},
				"comment":            "The concluded license was taken from the package level that the file was .\nThis information was found in the COPYING.txt file in the xyz directory.",
				"copyrightText":      "Copyright 2008-2010 John Smith",
				"fileContributors":   []string{"The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation"},
				"fileName":           "./package/foo.c",
				"fileTypes":          []string{"SOURCE"},
				"licenseComments":    "The concluded license was taken from the package level that the file was .\nThis information was found in the COPYING.txt file in the xyz directory.",
				"licenseConcluded":   "(LGPL-2.0-only OR LicenseRef-2)",
				"licenseInfoInFiles": []string{"GPL-2.0-only", "LicenseRef-2"},
				"noticeText":         "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.",
			},
		},
		"snippets": []interface{}{
			map[string]interface{}{
				"SPDXID":                "SPDXRef-Snippet",
				"comment":               "This snippet was identified as significant and highlighted in this Apache-2.0 file, when a commercial scanner identified it as being derived from file foo.c in package xyz which is licensed under GPL-2.0.",
				"copyrightText":         "Copyright 2008-2010 John Smith",
				"licenseComments":       "The concluded license was taken from package xyz, from which the snippet was copied into the current file. The concluded license information was found in the COPYING.txt file in package xyz.",
				"licenseConcluded":      "GPL-2.0-only",
				"licenseInfoInSnippets": []string{"GPL-2.0-only"},
				"name":                  "from linux kernel",
				"ranges": []interface{}{
					map[string]interface{}{
						"endPointer": map[string]interface{}{
							"offset":    420,
							"reference": "SPDXRef-DoapSource",
						},
						"startPointer": map[string]interface{}{
							"offset":    310,
							"reference": "SPDXRef-DoapSource",
						},
					},
					map[string]interface{}{
						"endPointer": map[string]interface{}{
							"lineNumber": 23,
							"reference":  "SPDXRef-DoapSource",
						},
						"startPointer": map[string]interface{}{
							"lineNumber": 5,
							"reference":  "SPDXRef-DoapSource",
						},
					},
				},
				"snippetFromFile": "SPDXRef-DoapSource",
			},
		},
	}
	var b []byte

	type args struct {
		doc *spdx.Document2_2
		buf *bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success test ",
			args: args{
				doc: test1,
				buf: bytes.NewBuffer(b),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RenderDocument2_2(tt.args.doc, tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("RenderDocument2_2() error = %v, wantErr %v", err, want)
			}
			jsonspec, _ := json.MarshalIndent(want, "", "\t")
			result := tt.args.buf.Bytes()
			if !reflect.DeepEqual(jsonspec, result) {
				t.Errorf("RenderDocument2_2() error = %v, wantErr %v", string(jsonspec), string(result))
			}
		})
	}
}
