// Package jsonloader is used to load and parse SPDX JSON documents
// into tools-golang data structures.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

//TODO: json validity check
//TODO: passsing an unrecornized key

func TestLoad2_2(t *testing.T) {

	jsonData, err := ioutil.ReadFile("jsonfiles/test.json") // b has type []byte
	if err != nil {
		log.Fatal(err)
	}

	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *spdxDocument2_2
		wantErr bool
	}{
		{
			name: "True test",
			args: args{content: jsonData},
			want: &spdxDocument2_2{
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
					{
						AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
						AnnotationDate:           "2011-03-13T00:00:00Z",
						AnnotationType:           "REVIEW",
						AnnotatorType:            "Person",
						Annotator:                "Suzanne Reviewer",
						AnnotationComment:        "Another example reviewer.",
					},
					{
						AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
						AnnotationDate:           "2010-01-29T18:30:22Z",
						AnnotationType:           "OTHER",
						AnnotatorType:            "Person",
						Annotator:                "Jane Doe ()",
						AnnotationComment:        "Document level annotation",
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
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
						Relationship: "DESCRIBES",
					},
					{
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
						Relationship: "DESCRIBES",
					},
					{
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "Saxon"},
						Relationship: "DYNAMIC_LINK",
					},
					{
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "JenaLib"},
						Relationship: "CONTAINS",
					},
					{
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "CommonsLangSrc"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "", SpecialID: "NOASSERTION"},
						Relationship: "GENERATED_FROM",
					},
					{
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "JenaLib"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
						Relationship: "CONTAINS",
					},
					{
						RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
						RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "fromDoap-0"},
						Relationship: "GENERATED_FROM",
					},
				},
				Packages: map[spdx.ElementID]*spdx.Package2_2{
					"Package": {
						PackageSPDXIdentifier:   "Package",
						PackageAttributionTexts: []string{"The GNU C Library is free software.  See the file COPYING.LIB for copying conditions, and LICENSES for notices about a few contributions that require these additional notices to be distributed.  License copyright years may be listed using range notation, e.g., 1996-2015, indicating that every year in the range, inclusive, is a copyrightable year that would otherwise be listed individually."},
						PackageChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
							"SHA256": {
								Algorithm: "SHA256",
								Value:     "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
							},
							"SHA1": {
								Algorithm: "SHA1",
								Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
							},
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
							"CommonsLangSrc": {
								FileSPDXIdentifier: "CommonsLangSrc",
								FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
									"SHA1": {
										Algorithm: "SHA1",
										Value:     "c2b4e1c67a2d28fced849ee1bb76e7391b93f125",
									},
								},
								FileComment:       "This file is used by Jena",
								FileCopyrightText: "Copyright 2001-2011 The Apache Software Foundation",
								FileContributor:   []string{"Apache Software Foundation"},
								FileName:          "./lib-source/commons-lang3-3.1-sources.jar",
								FileType:          []string{"ARCHIVE"},
								LicenseConcluded:  "Apache-2.0",
								LicenseInfoInFile: []string{"Apache-2.0"},
								FileNotice:        "Apache Commons Lang\nCopyright 2001-2011 The Apache Software Foundation\n\nThis product includes software developed by\nThe Apache Software Foundation (http://www.apache.org/).\n\nThis product includes software from the Spring Framework,\nunder the Apache License 2.0 (see: StringUtils.containsWhitespace())",
							},
							"JenaLib": {
								FileSPDXIdentifier: "JenaLib",
								FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
									"SHA1": {
										Algorithm: "SHA1",
										Value:     "3ab4e1c67a2d28fced849ee1bb76e7391b93f125",
									},
								},
								FileComment:       "This file belongs to Jena",
								FileCopyrightText: "(c) Copyright 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009 Hewlett-Packard Development Company, LP",
								FileContributor:   []string{"Apache Software Foundation", "Hewlett Packard Inc."},
								FileDependencies:  []string{"SPDXRef-CommonsLangSrc"},
								FileName:          "./lib-source/jena-2.6.3-sources.jar",
								FileType:          []string{"ARCHIVE"},
								LicenseComments:   "This license is used by Jena",
								LicenseConcluded:  "LicenseRef-1",
								LicenseInfoInFile: []string{"LicenseRef-1"},
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
					"fromDoap-1": {
						PackageSPDXIdentifier:     "fromDoap-1",
						PackageComment:            "This package was converted from a DOAP Project by the same name",
						PackageCopyrightText:      "NOASSERTION",
						PackageDownloadLocation:   "NOASSERTION",
						FilesAnalyzed:             false,
						IsFilesAnalyzedTagPresent: true,
						PackageHomePage:           "http://commons.apache.org/proper/commons-lang/",
						PackageLicenseConcluded:   "NOASSERTION",
						PackageLicenseDeclared:    "NOASSERTION",
						PackageName:               "Apache Commons Lang",
					},
					"fromDoap-0": {
						PackageSPDXIdentifier:     "fromDoap-0",
						PackageComment:            "This package was converted from a DOAP Project by the same name",
						PackageCopyrightText:      "NOASSERTION",
						PackageDownloadLocation:   "NOASSERTION",
						FilesAnalyzed:             false,
						IsFilesAnalyzedTagPresent: true,
						PackageHomePage:           "http://www.openjena.org/",
						PackageLicenseConcluded:   "NOASSERTION",
						PackageLicenseDeclared:    "NOASSERTION",
						PackageName:               "Jena",
					},

					"Saxon": {
						PackageSPDXIdentifier: "Saxon",
						PackageChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
							"SHA1": {
								Algorithm: "SHA1",
								Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
							},
						},
						PackageDescription:        "The Saxon package is a collection of tools for processing XML documents.",
						PackageDownloadLocation:   "https://sourceforge.net/projects/saxon/files/Saxon-B/8.8.0.7/saxonb8-8-0-7j.zip/download",
						FilesAnalyzed:             false,
						IsFilesAnalyzedTagPresent: true,
						PackageHomePage:           "http://saxon.sourceforge.net/",
						PackageLicenseComments:    "Other versions available for a commercial license",
						PackageLicenseConcluded:   "MPL-1.0",
						PackageLicenseDeclared:    "MPL-1.0",
						PackageName:               "Saxon",
						PackageFileName:           "saxonB-8.8.zip",
						PackageVersion:            "8.8",
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
				OtherLicenses: []*spdx.OtherLicense2_2{
					{
						ExtractedText:     "\"THE BEER-WARE LICENSE\" (Revision 42):\nphk@FreeBSD.ORG wrote this file. As long as you retain this notice you\ncan do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return Poul-Henning Kamp  </\nLicenseName: Beer-Ware License (Version 42)\nLicenseCrossReference:  http://people.freebsd.org/~phk/\nLicenseComment: \nThe beerware license has a couple of other standard variants.",
						LicenseIdentifier: "LicenseRef-Beerware-4.2",
					},
					{
						ExtractedText:     "/*\n * (c) Copyright 2009 University of Bristol\n * All rights reserved.\n *\n * Redistribution and use in source and binary forms, with or without\n * modification, are permitted provided that the following conditions\n * are met:\n * 1. Redistributions of source code must retain the above copyright\n *    notice, this list of conditions and the following disclaimer.\n * 2. Redistributions in binary form must reproduce the above copyright\n *    notice, this list of conditions and the following disclaimer in the\n *    documentation and/or other materials provided with the distribution.\n * 3. The name of the author may not be used to endorse or promote products\n *    derived from this software without specific prior written permission.\n *\n * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR\n * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\n * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.\n * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,\n * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT\n * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\n * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\n * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF\n * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n*/",
						LicenseIdentifier: "LicenseRef-4",
					},
					{
						LicenseComment:         "This is tye CyperNeko License",
						ExtractedText:          "The CyberNeko Software License, Version 1.0\n\n \n(C) Copyright 2002-2005, Andy Clark.  All rights reserved.\n \nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions\nare met:\n\n1. Redistributions of source code must retain the above copyright\n   notice, this list of conditions and the following disclaimer. \n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in\n   the documentation and/or other materials provided with the\n   distribution.\n\n3. The end-user documentation included with the redistribution,\n   if any, must include the following acknowledgment:  \n     \"This product includes software developed by Andy Clark.\"\n   Alternately, this acknowledgment may appear in the software itself,\n   if and wherever such third-party acknowledgments normally appear.\n\n4. The names \"CyberNeko\" and \"NekoHTML\" must not be used to endorse\n   or promote products derived from this software without prior \n   written permission. For written permission, please contact \n   andyc@cyberneko.net.\n\n5. Products derived from this software may not be called \"CyberNeko\",\n   nor may \"CyberNeko\" appear in their name, without prior written\n   permission of the author.\n\nTHIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESSED OR IMPLIED\nWARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\nOF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR OTHER CONTRIBUTORS\nBE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, \nOR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT \nOF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR \nBUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, \nWHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE \nOR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, \nEVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
						LicenseIdentifier:      "LicenseRef-3",
						LicenseName:            "CyberNeko License",
						LicenseCrossReferences: []string{"http://people.apache.org/~andyc/neko/LICENSE", "http://justasample.url.com"},
					},
					{
						ExtractedText:     "This package includes the GRDDL parser developed by Hewlett Packard under the following license:\n� Copyright 2007 Hewlett-Packard Development Company, LP\n\nRedistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met: \n\nRedistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer. \nRedistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution. \nThe name of the author may not be used to endorse or promote products derived from this software without specific prior written permission. \nTHIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
						LicenseIdentifier: "LicenseRef-2",
					},
					{
						ExtractedText:     "/*\n * (c) Copyright 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009 Hewlett-Packard Development Company, LP\n * All rights reserved.\n *\n * Redistribution and use in source and binary forms, with or without\n * modification, are permitted provided that the following conditions\n * are met:\n * 1. Redistributions of source code must retain the above copyright\n *    notice, this list of conditions and the following disclaimer.\n * 2. Redistributions in binary form must reproduce the above copyright\n *    notice, this list of conditions and the following disclaimer in the\n *    documentation and/or other materials provided with the distribution.\n * 3. The name of the author may not be used to endorse or promote products\n *    derived from this software without specific prior written permission.\n *\n * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR\n * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\n * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.\n * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,\n * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT\n * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\n * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\n * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF\n * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n*/",
						LicenseIdentifier: "LicenseRef-1",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load2_2(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load2_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// check creation info
			if !reflect.DeepEqual(got.CreationInfo, tt.want.CreationInfo) {
				t.Errorf("Load2_2() = %v, want %v", got.CreationInfo, tt.want.CreationInfo)
			}
			// check annotations
			for i := 0; i < len(tt.want.Annotations); i++ {
				if !reflect.DeepEqual(got.Annotations[i], tt.want.Annotations[i]) {
					t.Errorf("Load2_2() = %v, want %v", got.Annotations[i], tt.want.Annotations[i])
				}
			}
			// check relationships
			for i := 0; i < len(got.Relationships); i++ {
				if !reflect.DeepEqual(got.Relationships[i], tt.want.Relationships[i]) {
					t.Errorf("Load2_2() = %v, want %v", got.Relationships[i], tt.want.Relationships[i])
				}
			}
			//check unpackaged files
			for k, v := range tt.want.UnpackagedFiles {
				if !reflect.DeepEqual(got.UnpackagedFiles[k], v) {
					t.Errorf("Load2_2() = %v, want %v", got.UnpackagedFiles[k], v)
				}
			}
			// check packages
			for k, v := range tt.want.Packages {
				if !reflect.DeepEqual(got.Packages[k], v) {
					t.Errorf("Load2_2() = %v, want %v", got.Packages[k], v)
				}
			}
			// check other licenses
			for i := 0; i < len(got.OtherLicenses); i++ {
				if !reflect.DeepEqual(got.OtherLicenses[i], tt.want.OtherLicenses[i]) {
					t.Errorf("Load2_2() = %v, want %v", got.OtherLicenses[i], tt.want.OtherLicenses[i])
				}
			}
			// check reviews
			for i := 0; i < len(got.Reviews); i++ {
				if !reflect.DeepEqual(got.Reviews[i], tt.want.Reviews[i]) {
					t.Errorf("Load2_2() = %v, want %v", got.Reviews[i], tt.want.Reviews[i])
				}
			}

		})
	}
}
