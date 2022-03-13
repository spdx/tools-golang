// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderPackage2_2(t *testing.T) {
	type args struct {
		doc          *spdx.Document2_2
		jsondocument map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				doc: &spdx.Document2_2{
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
								"SHA1": {
									Algorithm: "SHA1",
									Value:     "d6a770ba38583ed4bb4525bd96e50461655d2758",
								},
								"MD5": {
									Algorithm: "MD5",
									Value:     "624c1abb3664f4b35547e7c73864ad24",
								},
							},
							FileComment:       "The concluded license was taken from the package level that the file was included in.\nThis information was found in the COPYING.txt file in the xyz directory.",
							FileCopyrightText: "Copyright 2008-2010 John Smith",
							FileContributor:   []string{"The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation"},
							FileName:          "./package/foo.c",
							FileType:          []string{"SOURCE"},
							LicenseComments:   "The concluded license was taken from the package level that the file was included in.",
							LicenseConcluded:  "(LGPL-2.0-only OR LicenseRef-2)",
							LicenseInfoInFile: []string{"GPL-2.0-only", "LicenseRef-2"},
							FileNotice:        "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.com\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the �Software�), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: \nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED �AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.",
						},
					},
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
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
		},
		{
			name: "success empty",
			args: args{
				doc: &spdx.Document2_2{
					Annotations: []*spdx.Annotation2_2{{}},
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
				map[string]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderPackage2_2(tt.args.doc, tt.args.jsondocument, make(map[spdx.ElementID]*spdx.File2_2))
			if (err != nil) != tt.wantErr {
				t.Errorf("renderPackage2_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got {
				if !reflect.DeepEqual(v, tt.want[k]) {
					t.Errorf("renderPackage2_2() = %v, want %v", v, tt.want[k])
				}
			}
		})
	}
}
