// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderSnippets2_2(t *testing.T) {
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
					UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{
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
									SnippetAttributionTexts:       []string{"text1", "text2 "},
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
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
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
					"snippetFromFile":  "SPDXRef-DoapSource",
					"attributionTexts": []string{"text1", "text2 "},
				},
			},
		},
		{
			name: "success empty",
			args: args{
				doc:          &spdx.Document2_2{},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
				map[string]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderSnippets2_2(tt.args.jsondocument, tt.args.doc.UnpackagedFiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderSnippets2_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got {
				if !reflect.DeepEqual(v, tt.want[k]) {
					t.Errorf("renderSnippets2_2() = %v, want %v", v, tt.want[k])
				}
			}
		})
	}
}
