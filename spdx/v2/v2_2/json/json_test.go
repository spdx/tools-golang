// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_2/example"
)

func TestLoad(t *testing.T) {
	want := example.Copy()

	file, err := os.Open("../../../../examples/sample-docs/json/SPDXJSONExample-v2.2.spdx.json")
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	var got spdx.Document
	err = json.ReadInto(file, &got)
	if err != nil {
		t.Errorf("json.parser.Load() error = %v", err)
		return
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got incorrect struct after parsing YAML example: %s", cmp.Diff(want, got))
		return
	}
}

func Test_Write(t *testing.T) {
	want := example.Copy()

	w := &bytes.Buffer{}

	if err := json.Write(&want, w); err != nil {
		t.Errorf("Write() error = %v", err.Error())
		return
	}

	// we should be able to parse what the writer wrote, and it should be identical to the original struct we wrote
	var got spdx.Document
	err := json.ReadInto(bytes.NewReader(w.Bytes()), &got)
	if err != nil {
		t.Errorf("failed to parse written document: %v", err.Error())
		return
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got incorrect struct after writing and re-parsing JSON example: %s", cmp.Diff(want, got))
		return
	}
}

func Test_ShorthandFields(t *testing.T) {
	contents := `{
		"spdxVersion": "SPDX-2.3",
		"dataLicense": "CC0-1.0",
		"SPDXID": "SPDXRef-DOCUMENT",
		"name": "SPDX-Tools-v2.0",
		"documentDescribes": [
			"SPDXRef-Container"
		],
		"packages": [
			{
				"name": "Container",
				"SPDXID": "SPDXRef-Container"
			},
			{
				"name": "Package-1",
				"SPDXID": "SPDXRef-Package-1",
				"versionInfo": "1.1.1",
				"hasFiles": [
					"SPDXRef-File-1",
					"SPDXRef-File-2"
				]
			},
			{
				"name": "Package-2",
				"SPDXID": "SPDXRef-Package-2",
				"versionInfo": "2.2.2"
			}
		],
		"files": [
			{
				"fileName": "./f1",
				"SPDXID": "SPDXRef-File-1"
			},
			{
				"fileName": "./f2",
				"SPDXID": "SPDXRef-File-2"
			}
		]
	}`

	doc := spdx.Document{}
	err := json.ReadInto(strings.NewReader(contents), &doc)

	require.NoError(t, err)

	id := func(s string) common.DocElementID {
		return common.DocElementID{
			ElementRefID: common.ElementID(s),
		}
	}

	require.Equal(t, spdx.Document{
		SPDXVersion:    spdx.Version,
		DataLicense:    spdx.DataLicense,
		SPDXIdentifier: "DOCUMENT",
		DocumentName:   "SPDX-Tools-v2.0",
		Packages: []*spdx.Package{
			{
				PackageName:           "Container",
				PackageSPDXIdentifier: "Container",
			},
			{
				PackageName:           "Package-1",
				PackageSPDXIdentifier: "Package-1",
				PackageVersion:        "1.1.1",
			},
			{
				PackageName:           "Package-2",
				PackageSPDXIdentifier: "Package-2",
				PackageVersion:        "2.2.2",
			},
		},
		Files: []*spdx.File{
			{
				FileName:           "./f1",
				FileSPDXIdentifier: "File-1",
			},
			{
				FileName:           "./f2",
				FileSPDXIdentifier: "File-2",
			},
		},
		Relationships: []*spdx.Relationship{
			{
				RefA:         id("DOCUMENT"),
				RefB:         id("Container"),
				Relationship: common.TypeRelationshipDescribe,
			},
			{
				RefA:         id("Package-1"),
				RefB:         id("File-1"),
				Relationship: common.TypeRelationshipContains,
			},
			{
				RefA:         id("Package-1"),
				RefB:         id("File-2"),
				Relationship: common.TypeRelationshipContains,
			},
		},
	}, doc)
}
