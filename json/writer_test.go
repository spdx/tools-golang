package json_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func Test_Write(t *testing.T) {
	tests := []struct {
		name   string
		doc    common.AnyDocument
		option []json.WriteOption
		want   string
	}{
		{
			name: "happy path",
			doc: spdx.Document{
				SPDXVersion:  "2.3",
				DocumentName: "test_doc",
			},
			want: `{"spdxVersion":"2.3","dataLicense":"","SPDXID":"SPDXRef-","name":"test_doc","documentNamespace":"","creationInfo":null}
`,
		},
		{
			name: "happy path with Indent option",
			doc: spdx.Document{
				SPDXVersion:  "2.3",
				DocumentName: "test_doc",
			},
			option: []json.WriteOption{json.Indent(" ")},
			want: `{
 "spdxVersion": "2.3",
 "dataLicense": "",
 "SPDXID": "SPDXRef-",
 "name": "test_doc",
 "documentNamespace": "",
 "creationInfo": null
}
`,
		},
		{
			name: "happy path with EscapeHTML==true option",
			doc: spdx.Document{
				SPDXVersion:  "2.3",
				DocumentName: "test_doc_>",
			},
			option: []json.WriteOption{json.EscapeHTML(true)},
			want:   "{\"spdxVersion\":\"2.3\",\"dataLicense\":\"\",\"SPDXID\":\"SPDXRef-\",\"name\":\"test_doc_\\u003e\",\"documentNamespace\":\"\",\"creationInfo\":null}\n",
		},
		{
			name: "happy path with EscapeHTML==false option",
			doc: spdx.Document{
				SPDXVersion:  "2.3",
				DocumentName: "test_doc_>",
			},
			option: []json.WriteOption{json.EscapeHTML(false)},
			want:   "{\"spdxVersion\":\"2.3\",\"dataLicense\":\"\",\"SPDXID\":\"SPDXRef-\",\"name\":\"test_doc_>\",\"documentNamespace\":\"\",\"creationInfo\":null}\n",
		},
		{
			name: "happy path with EscapeHTML==false option",
			doc: spdx.Document{
				SPDXVersion:  "2.3",
				DocumentName: "test_doc_>",
			},
			option: []json.WriteOption{json.EscapeHTML(false)},
			want:   "{\"spdxVersion\":\"2.3\",\"dataLicense\":\"\",\"SPDXID\":\"SPDXRef-\",\"name\":\"test_doc_>\",\"documentNamespace\":\"\",\"creationInfo\":null}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := json.Write(tt.doc, buf, tt.option...)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
