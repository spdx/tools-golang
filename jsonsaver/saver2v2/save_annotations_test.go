// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderAnnotations2_2(t *testing.T) {
	type args struct {
		annotations []*spdx.Annotation2_2
		eID         spdx.DocElementID
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
				annotations: []*spdx.Annotation2_2{
					{
						AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
						AnnotationDate:           "2011-01-29T18:30:22Z",
						AnnotationType:           "OTHER",
						AnnotatorType:            "Person",
						Annotator:                "File Commenter",
						AnnotationComment:        "File level annotation",
					},
					{
						AnnotationSPDXIdentifier: spdx.DocElementID{DocumentRefID: "", ElementRefID: "Package"},
						AnnotationDate:           "2011-01-29T18:30:22Z",
						AnnotationType:           "OTHER",
						AnnotatorType:            "Person",
						Annotator:                "Package Commenter",
						AnnotationComment:        "Package level annotation",
					},
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
				eID: spdx.MakeDocElementID("", "File"),
			},
			want: []interface{}{
				map[string]interface{}{
					"annotationDate": "2011-01-29T18:30:22Z",
					"annotationType": "OTHER",
					"annotator":      "Person: File Commenter",
					"comment":        "File level annotation",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderAnnotations2_2(tt.args.annotations, tt.args.eID)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderAnnotations2_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got {
				if !reflect.DeepEqual(v, tt.want[k]) {
					t.Errorf("renderAnnotations2_2() = %v, want %v", v, tt.want[k])
				}
			}
		})
	}
}
