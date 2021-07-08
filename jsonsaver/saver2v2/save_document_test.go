// Package saver2v2 contains functions to render and write a json
// formatted version of an in-memory SPDX document and its sections
// (version 2.2).
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package saver2v2

import (
	"bytes"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestRenderDocument2_2(t *testing.T) {

	test1 := &spdx.Document2_2{
		CreationInfo: &spdx.CreationInfo2_2{
			SPDXVersion:       "SPDX-2.2",
			DataLicense:       "CC0-1.0",
			SPDXIdentifier:    spdx.ElementID("DOCUMENT"),
			DocumentName:      "tools-golang-0.0.1.abcdef",
			DocumentNamespace: "https://github.com/spdx/spdx-docs/tools-golang/tools-golang-0.0.1.abcdef.whatever",
			CreatorPersons: []string{
				"John Doe ()",
			},
			CreatorOrganizations: []string{"ExampleCodeInspect"},
			CreatorTools:         []string{"LicenseFind-1.0"},
			DocumentComment:      "This document was created using SPDX 2.0 using licenses from the web site.",
			LicenseListVersion:   "3.8",
			CreatorComment:       "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
			Created:              "2018-10-10T06:20:00Z",
			ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{
				"spdx-tool-1.2": {
					DocumentRefID: "spdx-tool-1.2",
					URI:           "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
					Alg:           "SHA1",
					Checksum:      "d6a770ba38583ed4bb4525bd96e50461655d2759",
				},
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
				t.Errorf("RenderDocument2_2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
