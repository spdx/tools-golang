// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderCreationInfo2_2(t *testing.T) {
	type args struct {
		ci           *spdx.CreationInfo2_2
		jsondocument map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				ci: &spdx.CreationInfo2_2{
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
				jsondocument: make(map[string]interface{}),
			},
			want: map[string]interface{}{
				"dataLicense":       "CC0-1.0",
				"spdxVersion":       "SPDX-2.2",
				"SPDXID":            "SPDXRef-DOCUMENT",
				"documentNamespace": "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301",
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := renderCreationInfo2_2(tt.args.ci, tt.args.jsondocument); (err != nil) != tt.wantErr {
				t.Errorf("renderCreationInfo2_2() error = %v, wantErr %v", err, tt.wantErr)
			}
			for k, v := range tt.want {
				if !reflect.DeepEqual(tt.args.jsondocument[k], v) {
					t.Errorf("renderCreationInfo2_2() = %v, want %v", tt.args.jsondocument[k], v)
				}
			}
		})
	}
}
