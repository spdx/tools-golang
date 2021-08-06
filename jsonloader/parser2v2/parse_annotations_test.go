// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonAnnotations2_2(t *testing.T) {

	data := []byte(`{
		"annotations" : [ {
		"annotationDate" : "2010-02-10T00:00:00Z",
		"annotationType" : "REVIEW",
		"annotator" : "Person: Joe Reviewer",
		"comment" : "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses"
	  }, {
		"annotationDate" : "2011-03-13T00:00:00Z",
		"annotationType" : "REVIEW",
		"annotator" : "Person: Suzanne Reviewer",
		"comment" : "Another example reviewer."
	  }, {
		"annotationDate" : "2010-01-29T18:30:22Z",
		"annotationType" : "OTHER",
		"annotator" : "Person: Jane Doe ()",
		"comment" : "Document level annotation"
	  } ]
	}
  `)
	data2 := []byte(`{
	"annotations" : [ {
	"annotationDate" : "2010-02-10T00:00:00Z",
	"annotationType" : "REVIEW",
	"annotator" : "Person: Joe Reviewer",
	"comment" : "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
	"Hello":"hellp"
  }]
}
`)
	data3 := []byte(`{
	"annotations" : [ {
	"annotationDate" : "2010-02-10T00:00:00Z",
	"annotationType" : "REVIEW",
	"annotator" : "Fasle: Joe Reviewer",
	"comment" : "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
	"Hello":"hellp"
  }]
}
`)

	annotationstest1 := []*spdx.Annotation2_2{
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
	}

	var specs JSONSpdxDocument
	var specs2 JSONSpdxDocument
	var specs3 JSONSpdxDocument

	json.Unmarshal(data, &specs)
	json.Unmarshal(data2, &specs2)
	json.Unmarshal(data3, &specs3)

	type args struct {
		key           string
		value         interface{}
		doc           *spdxDocument2_2
		SPDXElementID spdx.DocElementID
	}
	tests := []struct {
		name    string
		spec    JSONSpdxDocument
		args    args
		want    []*spdx.Annotation2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:           "annotations",
				value:         specs["annotations"],
				doc:           &spdxDocument2_2{},
				SPDXElementID: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
			},
			want:    annotationstest1,
			wantErr: false,
		},
		{
			name: "failure test - invaid creator type",
			spec: specs2,
			args: args{
				key:           "annotations",
				value:         specs2["annotations"],
				doc:           &spdxDocument2_2{},
				SPDXElementID: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failure test - invalid tag",
			spec: specs3,
			args: args{
				key:           "annotations",
				value:         specs3["annotations"],
				doc:           &spdxDocument2_2{},
				SPDXElementID: spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonAnnotations2_2(tt.args.key, tt.args.value, tt.args.doc, tt.args.SPDXElementID); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonAnnotations2_2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				for i := 0; i < len(tt.want); i++ {
					if !reflect.DeepEqual(tt.args.doc.Annotations[i], tt.want[i]) {
						t.Errorf("Load2_2() = %v, want %v", tt.args.doc.Annotations[i], tt.want[i])
					}
				}
			}

		})
	}
}
