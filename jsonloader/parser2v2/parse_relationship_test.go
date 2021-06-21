// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonRelationships2_2(t *testing.T) {
	data := []byte(`{
		"relationships" : [ {
			"spdxElementId" : "SPDXRef-DOCUMENT",
			"relatedSpdxElement" : "DocumentRef-spdx-tool-1.2:SPDXRef-ToolsElement",
			"relationshipType" : "COPY_OF"
		  }, {
			"spdxElementId" : "SPDXRef-DOCUMENT",
			"relatedSpdxElement" : "SPDXRef-Package",
			"relationshipType" : "CONTAINS"
		  }, {
			"spdxElementId" : "SPDXRef-DOCUMENT",
			"relatedSpdxElement" : "SPDXRef-File",
			"relationshipType" : "DESCRIBES"
		  }, {
			"spdxElementId" : "SPDXRef-DOCUMENT",
			"relatedSpdxElement" : "SPDXRef-Package",
			"relationshipType" : "DESCRIBES"
		  }, {
			"spdxElementId" : "SPDXRef-Package",
			"relatedSpdxElement" : "SPDXRef-Saxon",
			"relationshipType" : "DYNAMIC_LINK"
		  }, {
			"spdxElementId" : "SPDXRef-Package",
			"relatedSpdxElement" : "SPDXRef-JenaLib",
			"relationshipType" : "CONTAINS"
		  },{
			"spdxElementId" : "SPDXRef-CommonsLangSrc",
			"relatedSpdxElement" : "NOASSERTION",
			"relationshipType" : "GENERATED_FROM"
		  } , {
			"spdxElementId" : "SPDXRef-JenaLib",
			"relatedSpdxElement" : "SPDXRef-Package",
			"relationshipType" : "CONTAINS"
		  }, {
			"spdxElementId" : "SPDXRef-File",
			"relatedSpdxElement" : "SPDXRef-fromDoap-0",
			"relationshipType" : "GENERATED_FROM"
		  } ]
		  }
  `)

	Relationship := []*spdx.Relationship2_2{
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
			RefA:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "CommonsLangSrc"},
			RefB:         spdx.DocElementID{DocumentRefID: "", ElementRefID: "", SpecialID: "NOASSERTION"},
			Relationship: "GENERATED_FROM",
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
	}

	var specs JSONSpdxDocument
	json.Unmarshal(data, &specs)

	type args struct {
		key   string
		value interface{}
		doc   *spdxDocument2_2
	}
	tests := []struct {
		name    string
		spec    JSONSpdxDocument
		args    args
		want    []*spdx.Relationship2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:   "relationships",
				value: specs["relationships"],
				doc:   &spdxDocument2_2{},
			},
			want:    Relationship,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonRelationships2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonRelationships2_2() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i := 0; i < len(tt.want); i++ {
				if !reflect.DeepEqual(tt.args.doc.Relationships[i], tt.want[i]) {
					t.Errorf("Load2_2() = %v, want %v", tt.args.doc.Relationships[i], tt.want[i])
				}
			}

		})
	}
}
