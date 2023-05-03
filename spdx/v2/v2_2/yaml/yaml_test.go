// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package yaml

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_2/example"
	"github.com/spdx/tools-golang/yaml"
)

func Test_Read(t *testing.T) {
	want := example.Copy()

	want.Relationships = append(want.Relationships, []*spdx.Relationship{
		{
			RefA:         common.DocElementID{ElementRefID: "DOCUMENT"},
			RefB:         common.DocElementID{ElementRefID: "File"},
			Relationship: "DESCRIBES",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "DOCUMENT"},
			RefB:         common.DocElementID{ElementRefID: "Package"},
			Relationship: "DESCRIBES",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "Package"},
			RefB:         common.DocElementID{ElementRefID: "CommonsLangSrc"},
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "Package"},
			RefB:         common.DocElementID{ElementRefID: "JenaLib"},
			Relationship: "CONTAINS",
		},
		{
			RefA:         common.DocElementID{ElementRefID: "Package"},
			RefB:         common.DocElementID{ElementRefID: "DoapSource"},
			Relationship: "CONTAINS",
		},
	}...)

	file, err := os.Open("../../../../examples/sample-docs/yaml/SPDXYAMLExample-2.2.spdx.yaml")
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	var got spdx.Document
	err = yaml.ReadInto(file, &got)
	if err != nil {
		t.Errorf("yaml.Read() error = %v", err)
		return
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(spdx.Package{})) {
		t.Errorf("got incorrect struct after parsing YAML example: %s", cmp.Diff(want, got, cmpopts.IgnoreUnexported(spdx.Package{})))
		return
	}
}

func Test_Write(t *testing.T) {
	want := example.Copy()

	// we always output FilesAnalyzed, even though we handle reading files where it is omitted
	for _, p := range want.Packages {
		p.IsFilesAnalyzedTagPresent = true
	}

	w := &bytes.Buffer{}
	if err := yaml.Write(want, w); err != nil {
		t.Errorf("Save() error = %v", err.Error())
		return
	}

	// we should be able to parse what the writer wrote, and it should be identical to the original handwritten struct
	var got spdx.Document
	err := yaml.ReadInto(bytes.NewReader(w.Bytes()), &got)
	if err != nil {
		t.Errorf("failed to parse written document: %v", err.Error())
		return
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(spdx.Package{})) {
		t.Errorf("got incorrect struct after writing and re-parsing YAML example: %s", cmp.Diff(want, got, cmpopts.IgnoreUnexported(spdx.Package{})))
		return
	}
}
