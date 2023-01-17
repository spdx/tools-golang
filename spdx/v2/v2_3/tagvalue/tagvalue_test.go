// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package tagvalue

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/spdx/tools-golang/spdx/v2/v2_3/example"
	"github.com/spdx/tools-golang/tagvalue"
)

var update = *flag.Bool("update-snapshots", false, "update the example snapshot")

func Test_Read(t *testing.T) {
	fileName := "../../../../examples/sample-docs/tv/SPDXTagExample-v2.3.spdx"

	want := example.Copy()

	if update {
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			t.Errorf("unable to open file to write: %v", err)
		}
		err = tagvalue.Write(want, f)
		if err != nil {
			t.Errorf("unable to call tagvalue.Write: %v", err)
		}
	} else {
		// tagvalue.Write sorts a few items in the in-memory SPDX document, run this
		// so the same sorting happens:
		err := tagvalue.Write(want, &bytes.Buffer{})
		if err != nil {
			t.Errorf("unable to call tagvalue.Write: %v", err)
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	var got spdx.Document
	err = tagvalue.ReadInto(file, &got)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}

	if !cmp.Equal(want, got, ignores...) {
		t.Errorf("got incorrect struct after parsing example: %s", cmp.Diff(want, got, ignores...))
		return
	}
}

func Test_ReadWrite(t *testing.T) {
	want := example.Copy()

	w := &bytes.Buffer{}
	// get a copy of the handwritten struct so we don't mutate it on accident
	if err := tagvalue.Write(want, w); err != nil {
		t.Errorf("Write() error = %v", err.Error())
		return
	}

	// we should be able to parse what the writer wrote, and it should be identical to the original struct we wrote
	var got spdx.Document
	err := tagvalue.ReadInto(bytes.NewReader(w.Bytes()), &got)
	if err != nil {
		t.Errorf("failed to parse written document: %v", err.Error())
		return
	}

	if !cmp.Equal(want, got, ignores...) {
		t.Errorf("got incorrect struct after writing and re-parsing YAML example: %s", cmp.Diff(want, got, ignores...))
		return
	}
}

var ignores = []cmp.Option{
	cmpopts.IgnoreFields(spdx.Document{}, "Snippets"),
	cmpopts.IgnoreFields(spdx.File{}, "Annotations"),
	cmpopts.IgnoreFields(spdx.Package{}, "IsFilesAnalyzedTagPresent", "Annotations"),
}
