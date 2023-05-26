// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package yaml

import (
	"bytes"
	jsonenc "encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_2/example"
	"github.com/spdx/tools-golang/yaml"
)

func Test_Read(t *testing.T) {
	want := example.Copy()

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

	if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(spdx.Package{}), cmpopts.SortSlices(relationshipLess)); len(diff) > 0 {
		t.Errorf("got incorrect struct after parsing YAML example: %s", diff)
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

	if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(spdx.Package{}), cmpopts.SortSlices(relationshipLess)); len(diff) > 0 {
		t.Errorf("got incorrect struct after writing and re-parsing YAML example: %s", diff)
		return
	}
}

func relationshipLess(a, b *spdx.Relationship) bool {
	aStr, _ := jsonenc.Marshal(a)
	bStr, _ := jsonenc.Marshal(b)
	return string(aStr) < string(bStr)
}
