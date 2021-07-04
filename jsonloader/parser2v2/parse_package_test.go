// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonPackages2_2(t *testing.T) {

	data := []byte(`{
		"packages" : [ {
			"SPDXID" : "SPDXRef-Package",
			"annotations" : [ {
			  "annotationDate" : "2011-01-29T18:30:22Z",
			  "annotationType" : "OTHER",
			  "annotator" : "Person: Package Commenter",
			  "comment" : "Package level annotation"
			} ],
			"attributionTexts" : [ "The GNU C Library is free software.  See the file COPYING.LIB for copying conditions, and LICENSES for notices about a few contributions that require these additional notices to be distributed.  License copyright years may be listed using range notation, e.g., 1996-2015, indicating that every year in the range, inclusive, is a copyrightable year that would otherwise be listed individually." ],
			"checksums" : [ {
			  "algorithm" : "SHA256",
			  "checksumValue" : "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"
			}, {
			  "algorithm" : "SHA1",
			  "checksumValue" : "85ed0817af83a24ad8da68c2b5094de69833983c"
			}, {
			  "algorithm" : "MD5",
			  "checksumValue" : "624c1abb3664f4b35547e7c73864ad24"
			} ],
			"copyrightText" : "Copyright 2008-2010 John Smith",
			"description" : "The GNU C Library defines functions that are specified by the ISO C standard, as well as additional features specific to POSIX and other derivatives of the Unix operating system, and extensions specific to GNU systems.",
			"downloadLocation" : "http://ftp.gnu.org/gnu/glibc/glibc-ports-2.15.tar.gz",
			"externalRefs" : [ {
			  "comment" : "This is the external ref for Acme",
			  "referenceCategory" : "OTHER",
			  "referenceLocator" : "acmecorp/acmenator/4.1.3-alpha",
			  "referenceType" : "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#LocationRef-acmeforge"
			}, {
			  "referenceCategory" : "SECURITY",
			  "referenceLocator" : "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
			  "referenceType" : "http://spdx.org/rdf/references/cpe23Type"
			} ],
			"filesAnalyzed" : true,
			"hasFiles" : [ "SPDXRef-JenaLib", "SPDXRef-DoapSource", "SPDXRef-CommonsLangSrc" ],
			"homepage" : "http://ftp.gnu.org/gnu/glibc",
			"licenseComments" : "The license for this project changed with the release of version x.y.  The version of the project included here post-dates the license change.",
			"licenseConcluded" : "(LGPL-2.0-only OR LicenseRef-3)",
			"licenseDeclared" : "(LGPL-2.0-only AND LicenseRef-3)",
			"licenseInfoFromFiles" : [ "GPL-2.0-only", "LicenseRef-2", "LicenseRef-1" ],
			"name" : "glibc",
			"originator" : "Organization: ExampleCodeInspect (contact@example.com)",
			"packageFileName" : "glibc-2.11.1.tar.gz",
			"packageVerificationCode" : {
			  "packageVerificationCodeExcludedFiles" : [ "excludes: ./package.spdx" ],
			  "packageVerificationCodeValue" : "d6a770ba38583ed4bb4525bd96e50461655d2758"
			},
			"sourceInfo" : "uses glibc-2_11-branch from git://sourceware.org/git/glibc.git.",
			"summary" : "GNU C library.",
			"supplier" : "Person: Jane Doe (jane.doe@example.com)",
			"versionInfo" : "2.11.1"
		  }, {
			"SPDXID" : "SPDXRef-fromDoap-1",
			"comment" : "This package was converted from a DOAP Project by the same name",
			"copyrightText" : "NOASSERTION",
			"downloadLocation" : "NOASSERTION",
			"filesAnalyzed" : false,
			"homepage" : "http://commons.apache.org/proper/commons-lang/",
			"licenseConcluded" : "NOASSERTION",
			"licenseDeclared" : "NOASSERTION",
			"name" : "Apache Commons Lang"
		  }, {
			"SPDXID" : "SPDXRef-fromDoap-0",
			"comment" : "This package was converted from a DOAP Project by the same name",
			"copyrightText" : "NOASSERTION",
			"downloadLocation" : "NOASSERTION",
			"filesAnalyzed" : false,
			"homepage" : "http://www.openjena.org/",
			"licenseConcluded" : "NOASSERTION",
			"licenseDeclared" : "NOASSERTION",
			"name" : "Jena"
		  }, {
			"SPDXID" : "SPDXRef-Saxon",
			"checksums" : [ {
			  "algorithm" : "SHA1",
			  "checksumValue" : "85ed0817af83a24ad8da68c2b5094de69833983c"
			} ],
			"description" : "The Saxon package is a collection of tools for processing XML documents.",
			"downloadLocation" : "https://sourceforge.net/projects/saxon/files/Saxon-B/8.8.0.7/saxonb8-8-0-7j.zip/download",
			"filesAnalyzed" : false,
			"homepage" : "http://saxon.sourceforge.net/",
			"licenseComments" : "Other versions available for a commercial license",
			"licenseConcluded" : "MPL-1.0",
			"licenseDeclared" : "MPL-1.0",
			"name" : "Saxon",
			"packageFileName" : "saxonB-8.8.zip",
			"versionInfo" : "8.8"
		  } ]		
		}
  `)

	document := spdxDocument2_2{
		UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{
			"DoapSource": {
				FileSPDXIdentifier: "DoapSource",
				FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
					"SHA1": {
						Algorithm: "SHA1",
						Value:     "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
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
	}

	packagetest1 := map[spdx.ElementID]*spdx.Package2_2{
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
	}
	var specs JSONSpdxDocument
	json.Unmarshal(data, &specs)

	type args struct {
		key   string
		value interface{}
		doc   *spdxDocument2_2
	}
	tests := []struct {
		name    string
		spec    JSONSpdxDocument
		args    args
		want    map[spdx.ElementID]*spdx.Package2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:   "packages",
				value: specs["packages"],
				doc:   &document,
			},
			want:    packagetest1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonPackages2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonPackages2_2() error = %v, wantErr %v", err, tt.wantErr)
			}

			for k, v := range tt.want {
				if !reflect.DeepEqual(tt.args.doc.Packages[k], v) {
					t.Errorf("Load2_2() = %v, want %v", tt.args.doc.Packages[k], v)
				}
			}

		})
	}
}
