package reader

import (
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
	"testing"
)

func Test_tvParser_parsePairForDocument(t *testing.T) {
	type fields struct {
		doc       *v2_2.Document
		st        tvParserState
		pkg       *v2_2.Package
		pkgExtRef *v2_2.PackageExternalReference
		file      *v2_2.File
		fileAOP   *v2_2.ArtifactOfProject
		snippet   *v2_2.Snippet
		otherLic  *v2_2.OtherLicense
		rln       *v2_2.Relationship
		ann       *v2_2.Annotation
		rev       *v2_2.Review
	}
	type args struct {
		tag   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test parser has no document",

			fields: fields{
				doc: nil,
			},
			args: args{
				tag:   "DocumentComment",
				value: "test comment",
			},
			wantErr: true,
		},
		{
			name: "test tag not equal to document comment",

			fields: fields{
				doc: &v2_2.Document{},
			},
			args: args{
				tag:   "invalid tag",
				value: "test comment",
			},
			wantErr: true,
		},

		{
			name: "good test",

			fields: fields{
				doc: &v2_2.Document{},
			},
			args: args{
				tag:   "DocumentComment",
				value: "test comment",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &tvParser{
				doc:       tt.fields.doc,
				st:        tt.fields.st,
				pkg:       tt.fields.pkg,
				pkgExtRef: tt.fields.pkgExtRef,
				file:      tt.fields.file,
				fileAOP:   tt.fields.fileAOP,
				snippet:   tt.fields.snippet,
				otherLic:  tt.fields.otherLic,
				rln:       tt.fields.rln,
				ann:       tt.fields.ann,
				rev:       tt.fields.rev,
			}
			if err := parser.parsePairForDocument(tt.args.tag, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("parsePairForDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
