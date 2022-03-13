// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func Test_renderReviews2_2(t *testing.T) {
	type args struct {
		reviews      []*spdx.Review2_2
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
				reviews: []*spdx.Review2_2{
					{
						ReviewDate:    "2010-02-10T00:00:00Z",
						ReviewerType:  "Person",
						Reviewer:      "Joe Reviewer",
						ReviewComment: "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
					},
				},
				jsondocument: make(map[string]interface{}),
			},
			want: []interface{}{
				map[string]interface{}{
					"reviewDate": "2010-02-10T00:00:00Z",
					"reviewer":   "Person: Joe Reviewer",
					"comment":    "This is just an example.  Some of the non-standard licenses look like they are actually BSD 3 clause licenses",
				},
			},
		},
		{
			name: "success empty",
			args: args{
				reviews: []*spdx.Review2_2{
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
			got, err := renderReviews2_2(tt.args.reviews, tt.args.jsondocument)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderReviews2_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got {
				if !reflect.DeepEqual(v, tt.want[k]) {
					t.Errorf("renderReviews2_2() = %v, want %v", v, tt.want[k])
				}
			}
		})
	}
}
