package convert

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/common"
)

func TestDocument(t *testing.T) {
	type testStructa struct {
		foo string
	}
	type testStructb struct {
		foo string
	}
	type args struct {
		from common.AnyDocument
		to   common.AnyDocument
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not a pointer",
			args: args{
				from: struct{}{},
				to:   struct{}{},
			},
			wantErr: true,
		},
		{
			name: "both pointers same type",
			args: args{
				from: &testStructa{"a"},
				to:   &testStructa{"a"},
			},
			wantErr: false,
		},
		{
			name: "both pointers different type",
			args: args{
				from: &testStructa{"a"},
				to:   &testStructb{"a"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Document(tt.args.from, tt.args.to)
			if (got != nil) != tt.wantErr {
				t.Errorf("Document() error = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}
