// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestJSONSpdxDocument_parseJsonReviews2_2(t *testing.T) {

	data := []byte(`{
		"revieweds" : [ {
		"reviewDate" : "2010-02-10T00:00:00Z",
		"reviewer" : "Person: Joe Reviewer",
		"comment" : "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses"
	  }, {
		"reviewDate" : "2011-03-13T00:00:00Z",
		"reviewer" : "Person: Suzanne Reviewer",
		"comment" : "Another example reviewer."
	  }]
	}
  `)

	reviewstest1 := []*spdx.Review2_2{
		{
			ReviewDate:    "2010-02-10T00:00:00Z",
			ReviewerType:  "Person",
			Reviewer:      "Joe Reviewer",
			ReviewComment: "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
		},
		{
			ReviewDate:    "2011-03-13T00:00:00Z",
			ReviewerType:  "Person",
			Reviewer:      "Suzanne Reviewer",
			ReviewComment: "Another example reviewer.",
		},
	}

	var specs JSONSpdxDocument
	json.Unmarshal(data, &specs)

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
		want    []*spdx.Review2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successTest",
			spec: specs,
			args: args{
				key:   "revieweds",
				value: specs["revieweds"],
				doc:   &spdxDocument2_2{},
			},
			want:    reviewstest1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.spec.parseJsonReviews2_2(tt.args.key, tt.args.value, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("JSONSpdxDocument.parseJsonAnnotations2_2() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i := 0; i < len(tt.want); i++ {
				if !reflect.DeepEqual(tt.args.doc.Reviews[i], tt.want[i]) {
					t.Errorf("Load2_2() = %v, want %v", tt.args.doc.Reviews[i], tt.want[i])
				}
			}

		})
	}
}
