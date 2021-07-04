// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jsonloader

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

func TestLoad2_2(t *testing.T) {

	file, err := os.Open("./parser2v2/jsonfiles/jsonloadertest.json")
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	type args struct {
		content io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *spdx.Document2_2
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success test",
			args: args{
				content: file,
			},
			want: &spdx.Document2_2{
				CreationInfo: &spdx.CreationInfo2_2{
					DataLicense:                "CC0-1.0",
					SPDXVersion:                "SPDX-2.2",
					SPDXIdentifier:             "DOCUMENT",
					DocumentName:               "SPDX-Tools-v2.0",
					ExternalDocumentReferences: make(map[string]spdx.ExternalDocumentRef2_2),
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
		})
	}
}
