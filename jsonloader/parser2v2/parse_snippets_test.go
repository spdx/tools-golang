// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonSnippets2_2(t *testing.T) {

	data := []byte(`{
		"snippets" : [ {
			"SPDXID" : "SPDXRef-Snippet",
			"comment" : "This snippet was identified as significant and highlighted in this Apache-2.0 file, when a commercial scanner identified it as being derived from file foo.c in package xyz which is licensed under GPL-2.0.",
			"copyrightText" : "Copyright 2008-2010 John Smith",
			"licenseComments" : "The concluded license was taken from package xyz, from which the snippet was copied into the current file. The concluded license information was found in the COPYING.txt file in package xyz.",
			"licenseConcluded" : "GPL-2.0-only",
			"licenseInfoInSnippets" : [ "GPL-2.0-only" ],
			"attributionTexts":["text1"],
			"name" : "from linux kernel",
			"ranges" : [ {
			  "endPointer" : {
				"lineNumber" : 23,
				"reference" : "SPDXRef-DoapSource"
			  },
			  "startPointer" : {
				"lineNumber" : 5,
				"reference" : "SPDXRef-DoapSource"
			  }
			}, {
			  "endPointer" : {
				"offset" : 420,
				"reference" : "SPDXRef-DoapSource"
			  },
			  "startPointer" : {
				"offset" : 310,
				"reference" : "SPDXRef-DoapSource"
			  }
			} ],
			"snippetFromFile" : "SPDXRef-DoapSource"
		  } ]
		}
  	`)

	filetest1 := spdx.File2_2{
		FileSPDXIdentifier: "DoapSource",
		FileChecksums:      map[spdx.ChecksumAlgorithm]spdx.Checksum{},
		Snippets: map[spdx.ElementID]*spdx.Snippet2_2{
			"Snippet": {
				SnippetSPDXIdentifier:         "Snippet",
				SnippetAttributionTexts:       []string{"text1"},
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
	}

	var specs JSONSpdxDocument
	err := json.Unmarshal(data, &specs)

	if err != nil {
		log.Fatal(err)
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
		want    *spdx.File2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:   "snippets",
				value: specs["snippets"],
				doc: &spdxDocument2_2{UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{
					"DoapSource": {
						FileSPDXIdentifier: "DoapSource",
						FileChecksums:      map[spdx.ChecksumAlgorithm]spdx.Checksum{},
						Snippets:           map[spdx.ElementID]*spdx.Snippet2_2{},
					},
				}},
			},
			want:    &filetest1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonSnippets2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonSnippets2_2() error = %v, wantErr %v", err, tt.wantErr)
			}

			for k, v := range tt.want.Snippets {
				if !reflect.DeepEqual(tt.args.doc.UnpackagedFiles["DoapSource"].Snippets[k], v) {
					t.Errorf("Load2_2() = %v, want %v", tt.args.doc.UnpackagedFiles["DoapSource"].Snippets[k], v)
				}
			}

		})
	}
}
