package example

import (
	"fmt"

	converter "github.com/anchore/go-struct-converter"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// Copy provides a deep copy of the example
func Copy() spdx.Document {
	out := spdx.Document{}
	err := converter.Convert(example, &out)
	if err != nil {
		panic(fmt.Errorf("unable to convert example doc: %w", err))
	}
	return out
}

// Example is handwritten translation of an official example SPDX document into a Go struct.
// We expect that the result of parsing the official document should be this value.
// We expect that the result of writing this struct should match the official example document.
var example = spdx.Document{
	DataLicense:       spdx.DataLicense,
	SPDXVersion:       spdx.Version,
	SPDXIdentifier:    "DOCUMENT",
	DocumentName:      "SPDX-Tools-v2.0",
	DocumentNamespace: "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301",
	CreationInfo: &spdx.CreationInfo{
		LicenseListVersion: "3.9",
		Creators: []common.Creator{
			{CreatorType: "Tool", Creator: "LicenseFind-1.0"},
			{CreatorType: "Organization", Creator: "ExampleCodeInspect ()"},
			{CreatorType: "Person", Creator: "Jane Doe ()"},
		},
		Created:        "2010-01-29T18:30:22Z",
		CreatorComment: "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
	},
	DocumentComment: "This document was created using SPDX 2.0 using licenses from the web site.",
	ExternalDocumentReferences: []spdx.ExternalDocumentRef{
		{
			DocumentRefID: "spdx-tool-1.2",
			URI:           "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
			Checksum: common.Checksum{
				Algorithm: common.SHA1,
				Value:     "d6a770ba38583ed4bb4525bd96e50461655d2759",
			},
		},
	},
	OtherLicenses: []*spdx.OtherLicense{
		{
			LicenseIdentifier: "LicenseRef-1",
			ExtractedText:     "/*\n * (c) Copyright 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009 Hewlett-Packard Development Company, LP\n * All rights reserved.\n *\n * Redistribution and use in source and binary forms, with or without\n * modification, are permitted provided that the following conditions\n * are met:\n * 1. Redistributions of source code must retain the above copyright\n *    notice, this list of conditions and the following disclaimer.\n * 2. Redistributions in binary form must reproduce the above copyright\n *    notice, this list of conditions and the following disclaimer in the\n *    documentation and/or other materials provided with the distribution.\n * 3. The name of the author may not be used to endorse or promote products\n *    derived from this software without specific prior written permission.\n *\n * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR\n * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\n * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.\n * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,\n * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT\n * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\n * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\n * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF\n * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n*/",
		},
		{
			LicenseIdentifier: "LicenseRef-2",
			ExtractedText:     "This package includes the GRDDL parser developed by Hewlett Packard under the following license:\n� Copyright 2007 Hewlett-Packard Development Company, LP\n\nRedistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met: \n\nRedistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer. \nRedistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution. \nThe name of the author may not be used to endorse or promote products derived from this software without specific prior written permission. \nTHIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
		},
		{
			LicenseIdentifier: "LicenseRef-4",
			ExtractedText:     "/*\n * (c) Copyright 2009 University of Bristol\n * All rights reserved.\n *\n * Redistribution and use in source and binary forms, with or without\n * modification, are permitted provided that the following conditions\n * are met:\n * 1. Redistributions of source code must retain the above copyright\n *    notice, this list of conditions and the following disclaimer.\n * 2. Redistributions in binary form must reproduce the above copyright\n *    notice, this list of conditions and the following disclaimer in the\n *    documentation and/or other materials provided with the distribution.\n * 3. The name of the author may not be used to endorse or promote products\n *    derived from this software without specific prior written permission.\n *\n * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR\n * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\n * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.\n * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,\n * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT\n * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\n * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\n * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF\n * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n*/",
		},
		{
			LicenseIdentifier:      "LicenseRef-Beerware-4.2",
			ExtractedText:          "\"THE BEER-WARE LICENSE\" (Revision 42):\nphk@FreeBSD.ORG wrote this file. As long as you retain this notice you\ncan do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return Poul-Henning Kamp",
			LicenseComment:         "The beerware license has a couple of other standard variants.",
			LicenseName:            "Beer-Ware License (Version 42)",
			LicenseCrossReferences: []string{"http://people.freebsd.org/~phk/"},
		},
		{
			LicenseIdentifier: "LicenseRef-3",
			ExtractedText:     "The CyberNeko Software License, Version 1.0\n\n \n(C) Copyright 2002-2005, Andy Clark.  All rights reserved.\n \nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions\nare met:\n\n1. Redistributions of source code must retain the above copyright\n   notice, this list of conditions and the following disclaimer. \n\n2. Redistributions in binary form must reproduce the above copyright\n   notice, this list of conditions and the following disclaimer in\n   the documentation and/or other materials provided with the\n   distribution.\n\n3. The end-user documentation included with the redistribution,\n   if any, must include the following acknowledgment:  \n     \"This product includes software developed by Andy Clark.\"\n   Alternately, this acknowledgment may appear in the software itself,\n   if and wherever such third-party acknowledgments normally appear.\n\n4. The names \"CyberNeko\" and \"NekoHTML\" must not be used to endorse\n   or promote products derived from this software without prior \n   written permission. For written permission, please contact \n   andyc@cyberneko.net.\n\n5. Products derived from this software may not be called \"CyberNeko\",\n   nor may \"CyberNeko\" appear in their name, without prior written\n   permission of the author.\n\nTHIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESSED OR IMPLIED\nWARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES\nOF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR OTHER CONTRIBUTORS\nBE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, \nOR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT \nOF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR \nBUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, \nWHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE \nOR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, \nEVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
			LicenseName:       "CyberNeko License",
			LicenseCrossReferences: []string{
				"http://people.apache.org/~andyc/neko/LICENSE",
				"http://justasample.url.com",
			},
			LicenseComment: "This is tye CyperNeko License",
		},
	},
	Annotations: []*spdx.Annotation{
		{
			Annotator: common.Annotator{
				Annotator:     "Jane Doe ()",
				AnnotatorType: "Person",
			},
			AnnotationDate:    "2010-01-29T18:30:22Z",
			AnnotationType:    "OTHER",
			AnnotationComment: "Document level annotation",
		},
		{
			Annotator: common.Annotator{
				Annotator:     "Joe Reviewer",
				AnnotatorType: "Person",
			},
			AnnotationDate:    "2010-02-10T00:00:00Z",
			AnnotationType:    "REVIEW",
			AnnotationComment: "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
		},
		{
			Annotator: common.Annotator{
				Annotator:     "Suzanne Reviewer",
				AnnotatorType: "Person",
			},
			AnnotationDate:    "2011-03-13T00:00:00Z",
			AnnotationType:    "REVIEW",
			AnnotationComment: "Another example reviewer.",
		},
	},
	Packages: []*spdx.Package{
		{
			PackageName:           "glibc",
			PackageSPDXIdentifier: "Package",
			PackageVersion:        "2.11.1",
			PackageFileName:       "glibc-2.11.1.tar.gz",
			PackageSupplier: &common.Supplier{
				Supplier:     "Jane Doe (jane.doe@example.com)",
				SupplierType: "Person",
			},
			PackageOriginator: &common.Originator{
				Originator:     "ExampleCodeInspect (contact@example.com)",
				OriginatorType: "Organization",
			},
			PackageDownloadLocation:   "http://ftp.gnu.org/gnu/glibc/glibc-ports-2.15.tar.gz",
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: true,
			PackageVerificationCode: common.PackageVerificationCode{
				Value:         "d6a770ba38583ed4bb4525bd96e50461655d2758",
				ExcludedFiles: []string{"./package.spdx"},
			},
			PackageChecksums: []common.Checksum{
				{
					Algorithm: "MD5",
					Value:     "624c1abb3664f4b35547e7c73864ad24",
				},
				{
					Algorithm: "SHA1",
					Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
				},
				{
					Algorithm: "SHA256",
					Value:     "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
				},
			},
			PackageHomePage:         "http://ftp.gnu.org/gnu/glibc",
			PackageSourceInfo:       "uses glibc-2_11-branch from git://sourceware.org/git/glibc.git.",
			PackageLicenseConcluded: "(LGPL-2.0-only OR LicenseRef-3)",
			PackageLicenseInfoFromFiles: []string{
				"GPL-2.0-only",
				"LicenseRef-2",
				"LicenseRef-1",
			},
			PackageLicenseDeclared: "(LGPL-2.0-only AND LicenseRef-3)",
			PackageLicenseComments: "The license for this project changed with the release of version x.y.  The version of the project included here post-dates the license change.",
			PackageCopyrightText:   "Copyright 2008-2010 John Smith",
			PackageSummary:         "GNU C library.",
			PackageDescription:     "The GNU C Library defines functions that are specified by the ISO C standard, as well as additional features specific to POSIX and other derivatives of the Unix operating system, and extensions specific to GNU systems.",
			PackageComment:         "",
			PackageExternalReferences: []*spdx.PackageExternalReference{
				{
					Category: "SECURITY",
					RefType:  "cpe23Type",
					Locator:  "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
				},
				{
					Category:           "OTHER",
					RefType:            "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#LocationRef-acmeforge",
					Locator:            "acmecorp/acmenator/4.1.3-alpha",
					ExternalRefComment: "This is the external ref for Acme",
				},
			},
			PackageAttributionTexts: []string{
				"The GNU C Library is free software.  See the file COPYING.LIB for copying conditions, and LICENSES for notices about a few contributions that require these additional notices to be distributed.  License copyright years may be listed using range notation, e.g., 1996-2015, indicating that every year in the range, inclusive, is a copyrightable year that would otherwise be listed individually.",
			},
			Files: nil,
			Annotations: []spdx.Annotation{
				{
					Annotator: common.Annotator{
						Annotator:     "Package Commenter",
						AnnotatorType: "Person",
					},
					AnnotationDate:    "2011-01-29T18:30:22Z",
					AnnotationType:    "OTHER",
					AnnotationComment: "Package level annotation",
				},
			},
		},
		{
			PackageSPDXIdentifier:     "fromDoap-1",
			PackageCopyrightText:      "NOASSERTION",
			PackageDownloadLocation:   "NOASSERTION",
			FilesAnalyzed:             false,
			IsFilesAnalyzedTagPresent: true,
			PackageHomePage:           "http://commons.apache.org/proper/commons-lang/",
			PackageLicenseConcluded:   "NOASSERTION",
			PackageLicenseDeclared:    "NOASSERTION",
			PackageName:               "Apache Commons Lang",
		},
		{
			PackageName:             "Jena",
			PackageSPDXIdentifier:   "fromDoap-0",
			PackageCopyrightText:    "NOASSERTION",
			PackageDownloadLocation: "https://search.maven.org/remotecontent?filepath=org/apache/jena/apache-jena/3.12.0/apache-jena-3.12.0.tar.gz",
			PackageExternalReferences: []*spdx.PackageExternalReference{
				{
					Category: "PACKAGE-MANAGER",
					RefType:  "purl",
					Locator:  "pkg:maven/org.apache.jena/apache-jena@3.12.0",
				},
			},
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: false,
			PackageHomePage:           "http://www.openjena.org/",
			PackageLicenseConcluded:   "NOASSERTION",
			PackageLicenseDeclared:    "NOASSERTION",
			PackageVersion:            "3.12.0",
		},
		{
			PackageSPDXIdentifier: "Saxon",
			PackageChecksums: []common.Checksum{
				{
					Algorithm: "SHA1",
					Value:     "85ed0817af83a24ad8da68c2b5094de69833983c",
				},
			},
			PackageCopyrightText:      "Copyright Saxonica Ltd",
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
	Files: []*spdx.File{
		{
			FileName:           "./src/org/spdx/parser/DOAPProject.java",
			FileSPDXIdentifier: "DoapSource",
			FileTypes: []string{
				"SOURCE",
			},
			Checksums: []common.Checksum{
				{
					Algorithm: "SHA1",
					Value:     "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
				},
			},
			LicenseConcluded: "Apache-2.0",
			LicenseInfoInFiles: []string{
				"Apache-2.0",
			},
			FileCopyrightText: "Copyright 2010, 2011 Source Auditor Inc.",
			FileContributors: []string{
				"Protecode Inc.",
				"SPDX Technical Team Members",
				"Open Logic Inc.",
				"Source Auditor Inc.",
				"Black Duck Software In.c",
			},
		},
		{
			FileSPDXIdentifier: "CommonsLangSrc",
			Checksums: []common.Checksum{
				{
					Algorithm: "SHA1",
					Value:     "c2b4e1c67a2d28fced849ee1bb76e7391b93f125",
				},
			},
			FileComment:        "This file is used by Jena",
			FileCopyrightText:  "Copyright 2001-2011 The Apache Software Foundation",
			FileContributors:   []string{"Apache Software Foundation"},
			FileName:           "./lib-source/commons-lang3-3.1-sources.jar",
			FileTypes:          []string{"ARCHIVE"},
			LicenseConcluded:   "Apache-2.0",
			LicenseInfoInFiles: []string{"Apache-2.0"},
			FileNotice:         "Apache Commons Lang\nCopyright 2001-2011 The Apache Software Foundation\n\nThis product includes software developed by\nThe Apache Software Foundation (http://www.apache.org/).\n\nThis product includes software from the Spring Framework,\nunder the Apache License 2.0 (see: StringUtils.containsWhitespace())",
		},
		{
			FileSPDXIdentifier: "JenaLib",
			Checksums: []common.Checksum{
				{
					Algorithm: "SHA1",
					Value:     "3ab4e1c67a2d28fced849ee1bb76e7391b93f125",
				},
			},
			FileComment:        "This file belongs to Jena",
			FileCopyrightText:  "(c) Copyright 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009 Hewlett-Packard Development Company, LP",
			FileContributors:   []string{"Apache Software Foundation", "Hewlett Packard Inc."},
			FileName:           "./lib-source/jena-2.6.3-sources.jar",
			FileTypes:          []string{"ARCHIVE"},
			LicenseComments:    "This license is used by Jena",
			LicenseConcluded:   "LicenseRef-1",
			LicenseInfoInFiles: []string{"LicenseRef-1"},
		},
		{
			FileSPDXIdentifier: "File",
			Annotations: []spdx.Annotation{
				{
					Annotator: common.Annotator{
						Annotator:     "File Commenter",
						AnnotatorType: "Person",
					},
					AnnotationDate:    "2011-01-29T18:30:22Z",
					AnnotationType:    "OTHER",
					AnnotationComment: "File level annotation",
				},
			},
			Checksums: []common.Checksum{
				{
					Algorithm: "SHA1",
					Value:     "d6a770ba38583ed4bb4525bd96e50461655d2758",
				},
				{
					Algorithm: "MD5",
					Value:     "624c1abb3664f4b35547e7c73864ad24",
				},
			},
			FileComment:        "The concluded license was taken from the package level that the file was included in.\nThis information was found in the COPYING.txt file in the xyz directory.",
			FileCopyrightText:  "Copyright 2008-2010 John Smith",
			FileContributors:   []string{"The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation"},
			FileName:           "./package/foo.c",
			FileTypes:          []string{"SOURCE"},
			LicenseComments:    "The concluded license was taken from the package level that the file was included in.",
			LicenseConcluded:   "(LGPL-2.0-only OR LicenseRef-2)",
			LicenseInfoInFiles: []string{"GPL-2.0-only", "LicenseRef-2"},
			FileNotice:         "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.com\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the �Software�), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: \nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED �AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.",
		},
	},
	Snippets: []spdx.Snippet{
		{
			SnippetSPDXIdentifier:         "Snippet",
			SnippetFromFileSPDXIdentifier: "DoapSource",
			Ranges: []common.SnippetRange{
				{
					StartPointer: common.SnippetRangePointer{
						Offset:             310,
						FileSPDXIdentifier: "DoapSource",
					},
					EndPointer: common.SnippetRangePointer{
						Offset:             420,
						FileSPDXIdentifier: "DoapSource",
					},
				},
				{
					StartPointer: common.SnippetRangePointer{
						LineNumber:         5,
						FileSPDXIdentifier: "DoapSource",
					},
					EndPointer: common.SnippetRangePointer{
						LineNumber:         23,
						FileSPDXIdentifier: "DoapSource",
					},
				},
			},
			SnippetLicenseConcluded: "GPL-2.0-only",
			LicenseInfoInSnippet:    []string{"GPL-2.0-only"},
			SnippetLicenseComments:  "The concluded license was taken from package xyz, from which the snippet was copied into the current file. The concluded license information was found in the COPYING.txt file in package xyz.",
			SnippetCopyrightText:    "Copyright 2008-2010 John Smith",
			SnippetComment:          "This snippet was identified as significant and highlighted in this Apache-2.0 file, when a commercial scanner identified it as being derived from file foo.c in package xyz which is licensed under GPL-2.0.",
			SnippetName:             "from linux kernel",
		},
	},
	Relationships: []*spdx.Relationship{
		{
			RefA:         common.MakeDocElementID("", "DOCUMENT"),
			RefB:         common.MakeDocElementID("", "Package"),
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.MakeDocElementID("", "Package"),
			RefB:         common.MakeDocElementID("", "CommonsLangSrc"),
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.MakeDocElementID("", "Package"),
			RefB:         common.MakeDocElementID("", "DoapSource"),
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.MakeDocElementID("", "DOCUMENT"),
			RefB:         common.MakeDocElementID("spdx-tool-1.2", "ToolsElement"),
			Relationship: "COPY_OF",
		},
		{
			RefA:         common.MakeDocElementID("", "DOCUMENT"),
			RefB:         common.MakeDocElementID("", "File"),
			Relationship: "DESCRIBES",
		},
		{
			RefA:         common.MakeDocElementID("", "DOCUMENT"),
			RefB:         common.MakeDocElementID("", "Package"),
			Relationship: "DESCRIBES",
		},
		{
			RefA:         common.MakeDocElementID("", "Package"),
			RefB:         common.MakeDocElementID("", "JenaLib"),
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.MakeDocElementID("", "Package"),
			RefB:         common.MakeDocElementID("", "Saxon"),
			Relationship: "DYNAMIC_LINK",
		},
		{
			RefA:         common.MakeDocElementID("", "CommonsLangSrc"),
			RefB:         common.MakeDocElementSpecial("NOASSERTION"),
			Relationship: "GENERATED_FROM",
		},
		{
			RefA:         common.MakeDocElementID("", "JenaLib"),
			RefB:         common.MakeDocElementID("", "Package"),
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.MakeDocElementID("", "File"),
			RefB:         common.MakeDocElementID("", "fromDoap-0"),
			Relationship: "GENERATED_FROM",
		},
	},
}
