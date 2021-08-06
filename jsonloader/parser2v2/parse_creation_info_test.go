// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonCreationInfo2_2(t *testing.T) {

	var specs JSONSpdxDocument

	type args struct {
		key   string
		value interface{}
		doc   *spdxDocument2_2
	}
	tests := []struct {
		name    string
		spec    JSONSpdxDocument
		args    args
		want    *spdx.CreationInfo2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		// check whether DataLicense is being parsed
		{
			name: "DataLicense",
			spec: specs,
			args: args{
				key:   "dataLicense",
				value: "CC0-1.0",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{DataLicense: "CC0-1.0", ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: false,
		},
		// check whether SPDXID is being parsed
		{
			name: "SPDXID",
			spec: specs,
			args: args{
				key:   "SPDXID",
				value: "SPDXRef-DOCUMENT",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{SPDXIdentifier: "DOCUMENT", ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: false,
		},
		// check whether DocumentName is being parsed
		{
			name: "DocumentName",
			spec: specs,
			args: args{
				key:   "name",
				value: "SPDX-Tools-v2.0",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{DocumentName: "SPDX-Tools-v2.0", ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: false,
		},
		// check whether SPDXVersion is being parsed
		{
			name: "spdxVersion",
			spec: specs,
			args: args{
				key:   "spdxVersion",
				value: "SPDX-2.2",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{SPDXVersion: "SPDX-2.2", ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: false,
		},
		// check whether DocumentComment is being parsed
		{
			name: "comment",
			spec: specs,
			args: args{
				key:   "comment",
				value: "This document was created using SPDX 2.0 using licenses from the web site.",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{DocumentComment: "This document was created using SPDX 2.0 using licenses from the web site.", ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: false,
		},
		// check whether DocumentNamespace is being parsed
		{
			name: "documentNamespace",
			spec: specs,
			args: args{
				key:   "documentNamespace",
				value: "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{DocumentNamespace: "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301", ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: false,
		},
		// check whether CreationInfo(Creators , CreatorComment and Licence List Version) is being parsed
		{
			name: "creationInfo",
			spec: specs,
			args: args{
				key: "creationInfo",
				value: map[string]interface{}{
					"comment":            "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
					"created":            "2010-01-29T18:30:22Z",
					"creators":           []string{"Tool: LicenseFind-1.0", "Organization: ExampleCodeInspect ()", "Person: Jane Doe ()"},
					"licenseListVersion": "3.8",
				},
				doc: &spdxDocument2_2{},
			},
			want: &spdx.CreationInfo2_2{
				CreatorComment:             "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
				Created:                    "2010-01-29T18:30:22Z",
				CreatorPersons:             []string{"Jane Doe ()"},
				CreatorOrganizations:       []string{"ExampleCodeInspect ()"},
				CreatorTools:               []string{"LicenseFind-1.0"},
				LicenseListVersion:         "3.8",
				ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{},
			},
			wantErr: false,
		},
		// check whether  ExternalDocumentReferences is being parsed
		{
			name: "externalDocumentRefs",
			spec: specs,
			args: args{
				key: "externalDocumentRefs",
				value: []map[string]interface{}{
					{
						"externalDocumentId": "DocumentRef-spdx-tool-1.2",
						"spdxDocument":       "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
						"checksum": map[string]interface{}{
							"algorithm":     "SHA1",
							"checksumValue": "d6a770ba38583ed4bb4525bd96e50461655d2759",
						},
					},
				},
				doc: &spdxDocument2_2{},
			},
			want: &spdx.CreationInfo2_2{
				ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{
					"spdx-tool-1.2": {
						DocumentRefID: "spdx-tool-1.2",
						URI:           "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
						Alg:           "SHA1",
						Checksum:      "d6a770ba38583ed4bb4525bd96e50461655d2759",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failure : Invalid tag ",
			spec: specs,
			args: args{
				key:   "invalid",
				value: "This document was created using SPDX 2.0 using licenses from the web site.",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: true,
		},
		{
			name: "failure : DocRef missing in ExternalRefs",
			spec: specs,
			args: args{
				key: "externalDocumentRefs",
				value: []map[string]interface{}{
					{
						"externalDocumentId": "spdx-tool-1.2",
						"spdxDocument":       "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
						"checksum": map[string]interface{}{
							"algorithm":     "SHA1",
							"checksumValue": "d6a770ba38583ed4bb4525bd96e50461655d2759",
						},
					},
				},
				doc: &spdxDocument2_2{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failure : invalid SPDXID",
			spec: specs,
			args: args{
				key:   "SPDXID",
				value: "DOCUMENT",
				doc:   &spdxDocument2_2{},
			},
			want:    &spdx.CreationInfo2_2{ExternalDocumentReferences: map[string]spdx.ExternalDocumentRef2_2{}},
			wantErr: true,
		},
		{
			name: "failure - invalid creator type",
			spec: specs,
			args: args{
				key: "creationInfo",
				value: map[string]interface{}{
					"comment":            "This package has been shipped in source and binary form.\nThe binaries were created with gcc 4.5.1 and expect to link to\ncompatible system run time libraries.",
					"created":            "2010-01-29T18:30:22Z",
					"creators":           []string{"Invalid: LicenseFind-1.0", "Organization: ExampleCodeInspect ()", "Person: Jane Doe ()"},
					"licenseListVersion": "3.8",
				},
				doc: &spdxDocument2_2{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.doc = &spdxDocument2_2{}
			if err := tt.spec.parseJsonCreationInfo2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonCreationInfo2_2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.doc.CreationInfo, tt.want) {
				t.Errorf("Load2_2() = %v, want %v", tt.args.doc.CreationInfo, tt.want)
			}

		})
	}
}
