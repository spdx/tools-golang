// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spdx/tools-golang/json"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/spdx/tools-golang/spdx/v2/v2_3/example"
)

var update = *flag.Bool("update-snapshots", false, "update the example snapshot")

func Test_Read(t *testing.T) {
	fileName := "../../../../examples/sample-docs/json/SPDXJSONExample-v2.3.spdx.json"

	want := example.Copy()

	if update {
		w := &bytes.Buffer{}

		err := json.Write(want, w)
		if err != nil {
			t.Errorf("unable to serialize SPDX 2.3 example to JSON: %v", err)
		}
		err = os.WriteFile(fileName, w.Bytes(), 0644)
		if err != nil {
			t.Errorf("unable to write SPDX 2.3 example to JSON: %v", err)
		}
	}

	file, err := os.Open(fileName)
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
