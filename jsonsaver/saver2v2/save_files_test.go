// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderFiles2_2(t *testing.T) {
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
					Annotations: []*spdx.Annotation2_2{
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
					UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{
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
							FileComment:          "The concluded license was taken from the package level that the file was .",
							FileCopyrightText:    "Copyright 2008-2010 John Smith",
							FileContributor:      []string{"The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation"},
							FileName:             "./package/foo.c",
							FileType:             []string{"SOURCE"},
							LicenseComments:      "The concluded license was taken from the package level that the file was included in.",
							LicenseConcluded:     "(LGPL-2.0-only OR LicenseRef-2)",
							LicenseInfoInFile:    []string{"GPL-2.0-only", "LicenseRef-2"},
							FileNotice:           "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.",
							FileAttributionTexts: []string{"text1", "text2 "},
						},
					},
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
				map[string]interface{}{
					"SPDXID": "SPDXRef-File",
					"annotations": []interface{}{
						map[string]interface{}{
							"annotationDate": "2011-01-29T18:30:22Z",
							"annotationType": "OTHER",
							"annotator":      "Person: File Commenter",
							"comment":        "File level annotation",
						},
					},
					"checksums": []interface{}{
						map[string]interface{}{
							"algorithm":     "MD5",
							"checksumValue": "624c1abb3664f4b35547e7c73864ad24",
						},
						map[string]interface{}{
							"algorithm":     "SHA1",
							"checksumValue": "d6a770ba38583ed4bb4525bd96e50461655d2758",
						},
					},
					"comment":            "The concluded license was taken from the package level that the file was .",
					"copyrightText":      "Copyright 2008-2010 John Smith",
					"fileContributors":   []string{"The Regents of the University of California", "Modified by Paul Mundt lethal@linux-sh.org", "IBM Corporation"},
					"fileName":           "./package/foo.c",
					"fileTypes":          []string{"SOURCE"},
					"licenseComments":    "The concluded license was taken from the package level that the file was included in.",
					"licenseConcluded":   "(LGPL-2.0-only OR LicenseRef-2)",
					"licenseInfoInFiles": []string{"GPL-2.0-only", "LicenseRef-2"},
					"noticeText":         "Copyright (c) 2001 Aaron Lehmann aaroni@vitelus.",
					"attributionTexts":   []string{"text1", "text2 "},
				},
			},
		},
		{
			name: "success empty",
			args: args{
				doc: &spdx.Document2_2{
					Annotations:     []*spdx.Annotation2_2{},
					UnpackagedFiles: map[spdx.ElementID]*spdx.File2_2{},
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
				map[string]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderFiles2_2(tt.args.doc, tt.args.jsondocument, tt.args.doc.UnpackagedFiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderFiles2_2() error = %v, wantErr %v", err, tt.wantErr)
			}
			for k, v := range got {
				if !reflect.DeepEqual(v, tt.want[k]) {
					t.Errorf("renderFiles2_2() error = %v, want %v", v, tt.want[k])
				}
			}

		})
	}
}
