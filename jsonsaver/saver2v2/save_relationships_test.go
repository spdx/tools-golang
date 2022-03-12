// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderRelationships2_2(t *testing.T) {
	type args struct {
		relationships []*spdx.Relationship2_2
		jsondocument  map[string]interface{}
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
				relationships: []*spdx.Relationship2_2{
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
						RefA:                spdx.DocElementID{DocumentRefID: "", ElementRefID: "DOCUMENT"},
						RefB:                spdx.DocElementID{DocumentRefID: "", ElementRefID: "File"},
						Relationship:        "DESCRIBES",
						RelationshipComment: "This is a comment.",
					},
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
				map[string]interface{}{
					"spdxElementId":      "SPDXRef-DOCUMENT",
					"relatedSpdxElement": "DocumentRef-spdx-tool-1.2:SPDXRef-ToolsElement",
					"relationshipType":   "COPY_OF",
				},
				map[string]interface{}{
					"spdxElementId":      "SPDXRef-DOCUMENT",
					"relatedSpdxElement": "SPDXRef-Package",
					"relationshipType":   "CONTAINS",
				},
				map[string]interface{}{
					"spdxElementId":      "SPDXRef-DOCUMENT",
					"relatedSpdxElement": "SPDXRef-File",
					"relationshipType":   "DESCRIBES",
					"comment":            "This is a comment.",
				},
			},
		},
		{
			name: "success empty",
			args: args{
				relationships: []*spdx.Relationship2_2{
					{},
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
			got, err := renderRelationships2_2(tt.args.relationships, tt.args.jsondocument)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderRelationships2_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got {
				if !reflect.DeepEqual(v, tt.want[k]) {
					t.Errorf("renderRelationships2_2() = %v, want %v", v, tt.want[k])
				}
			}

		})
	}
}
