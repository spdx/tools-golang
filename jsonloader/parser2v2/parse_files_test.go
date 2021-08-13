// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

//TODO: tests for annotations parsed from files

func TestJSONSpdxDocument_parseJsonFiles2_2(t *testing.T) {

	data := []byte(`{
		"files" : [ {
			"SPDXID" : "SPDXRef-DoapSource",
			"checksums" : [ {
			  "algorithm" : "SHA1",
			  "checksumValue" : "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"
			} ],
			"copyrightText" : "Copyright 2010, 2011 Source Auditor Inc.",
			"fileContributors" : [ "Protecode Inc.", "SPDX Technical Team Members", "Open Logic Inc.", "Source Auditor Inc.", "Black Duck Software In.c" ],
			"fileDependencies" : [ "SPDXRef-JenaLib", "SPDXRef-CommonsLangSrc" ],
			"fileName" : "./src/org/spdx/parser/DOAPProject.java",
			"fileTypes" : [ "SOURCE" ],
			"licenseConcluded" : "Apache-2.0",
			"attributionTexts":["text1"],
			"licenseInfoInFiles" : [ "Apache-2.0" ]
		  }, {
			"SPDXID" : "SPDXRef-CommonsLangSrc",
			"checksums" : [ {
			  "algorithm" : "SHA1",
			  "checksumValue" : "c2b4e1c67a2d28fced849ee1bb76e7391b93f125"
			} ],
			"comment" : "This file is used by Jena",
			"copyrightText" : "Copyright 2001-2011 The Apache Software Foundation",
			"fileContributors" : [ "Apache Software Foundation" ],
			"fileName" : "./lib-source/commons-lang3-3.1-sources.jar",
			"fileTypes" : [ "ARCHIVE" ],
			"licenseConcluded" : "Apache-2.0",
			"licenseInfoInFiles" : [ "Apache-2.0" ],
			"noticeText" : "Apache Commons Lang\nCopyright 2001-2011 The Apache Software Foundation\n\nThis product includes software developed by\nThe Apache Software Foundation (http://www.apache.org/).\n\nThis product includes software from the Spring Framework,\nunder the Apache License 2.0 (see: StringUtils.containsWhitespace())"
		  }, {
			"SPDXID" : "SPDXRef-JenaLib",
			"checksums" : [ {
			  "algorithm" : "SHA1",
			  "checksumValue" : "3ab4e1c67a2d28fced849ee1bb76e7391b93f125"
			} ],
			"comment" : "This file belongs to Jena",
			"copyrightText" : "(c) Copyright 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009 Hewlett-Packard Development Company, LP",
			"fileContributors" : [ "Apache Software Foundation", "Hewlett Packard Inc." ],
			"fileDependencies" : [ "SPDXRef-CommonsLangSrc" ],
			"fileName" : "./lib-source/jena-2.6.3-sources.jar",
			"fileTypes" : [ "ARCHIVE" ],
			"licenseComments" : "This license is used by Jena",
			"licenseConcluded" : "LicenseRef-1",
			"licenseInfoInFiles" : [ "LicenseRef-1" ]
		  }, {
			"SPDXID" : "SPDXRef-File",
			"annotations" : [ {
			  "annotationDate" : "2011-01-29T18:30:22Z",
			  "annotationType" : "OTHER",
			  "annotator" : "Person: File Commenter",
			  "comment" : "File level annotation"
			} ],
			"checksums" : [ {
			  "algorithm" : "SHA1",
			  "checksumValue" : "d6a770ba38583ed4bb4525bd96e50461655d2758"
			}, {
			  "algorithm" : "MD5",
			  "checksumValue" : "624c1abb3664f4b35547e7c73864ad24"
			} ],
			"comment" : "The concluded license was taken from the package level that the file was included in.\nThis information was found in the COPYING.txt file in the xyz directory.",
			"copyrightText" : "Copyright 2008-2010 John Smith",
			"fileContributors" : [ "The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation" ],
			"fileName" : "./package/foo.c",
			"fileTypes" : [ "SOURCE" ],
			"licenseComments" : "The concluded license was taken from the package level that the file was included in.",
			"licenseConcluded" : "(LGPL-2.0-only OR LicenseRef-2)",
			"licenseInfoInFiles" : [ "GPL-2.0-only", "LicenseRef-2" ],
			"noticeText" : "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.com\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the �Software�), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: \nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED �AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE."
		  } ]	
		}
  `)

	var specs JSONSpdxDocument
	json.Unmarshal(data, &specs)

	filestest1 := map[spdx.ElementID]*spdx.File2_2{
		"DoapSource": {
			FileSPDXIdentifier: "DoapSource",
			FileChecksums: map[spdx.ChecksumAlgorithm]spdx.Checksum{
				"SHA1": {
					Algorithm: "SHA1",
					Value:     "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
				},
			},
			FileCopyrightText:    "Copyright 2010, 2011 Source Auditor Inc.",
			FileContributor:      []string{"Protecode Inc.", "SPDX Technical Team Members", "Open Logic Inc.", "Source Auditor Inc.", "Black Duck Software In.c"},
			FileDependencies:     []string{"SPDXRef-JenaLib", "SPDXRef-CommonsLangSrc"},
			FileName:             "./src/org/spdx/parser/DOAPProject.java",
			FileType:             []string{"SOURCE"},
			LicenseConcluded:     "Apache-2.0",
			FileAttributionTexts: []string{"text1"},
			LicenseInfoInFile:    []string{"Apache-2.0"},
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
	}

	type args struct {
		key   string
		value interface{}
		doc   *spdxDocument2_2
	}
	tests := []struct {
		name    string
		spec    JSONSpdxDocument
		args    args
		want    map[spdx.ElementID]*spdx.File2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:   "files",
				value: specs["files"],
				doc:   &spdxDocument2_2{},
			},
			want:    filestest1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonFiles2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonFiles2_2() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				for k, v := range tt.want {
					if !reflect.DeepEqual(tt.args.doc.UnpackagedFiles[k], v) {
						t.Errorf("Load2_2() = %v, want %v", tt.args.doc.UnpackagedFiles[k], v)
					}
				}
			}

		})
	}
}
