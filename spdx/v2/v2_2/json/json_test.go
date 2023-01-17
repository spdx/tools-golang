// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_2/example"
)

func TestLoad(t *testing.T) {
	want := example.Copy()

	file, err := os.Open("../../../../examples/sample-docs/json/SPDXJSONExample-v2.2.spdx.json")
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	var got v2_2.Document
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
	var got v2_2.Document
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
