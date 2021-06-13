// Package jsonloader is used to load and parse SPDX JSON documents
// into tools-golang data structures.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jsonloader2v2

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
		// TODO: Add test cases.
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
					CreatorPersons:       []string{"Jane Doe"},
					CreatorOrganizations: []string{"ExampleCodeInspect"},
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
			if !reflect.DeepEqual(got.CreationInfo, tt.want.CreationInfo) {
				t.Errorf("Load2_2() = %v, want %v", got.CreationInfo, tt.want.CreationInfo)
			}
			for i := 0; i < len(got.Annotations); i++ {
				if !reflect.DeepEqual(got.Annotations[i], tt.want.Annotations[i]) {
					t.Errorf("Load2_2() = %v, want %v", got.Annotations[i], tt.want.Annotations[i])
				}
			}
			for i := 0; i < len(got.Relationships); i++ {
				if !reflect.DeepEqual(got.Relationships[i], tt.want.Relationships[i]) {
					t.Errorf("Load2_2() = %v, want %v", got.Relationships[i], tt.want.Relationships[i])
				}
			}
		})
	}
}
