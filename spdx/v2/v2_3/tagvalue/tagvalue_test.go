// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package tagvalue

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	spdxjson "github.com/spdx/tools-golang/json"
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

// Test_WriteNullElements covers the cross-format path: the JSON reader accepts a
// document with a null element and leaves it in the slice, so writing the result
// back out as tag-value hands the writer a nil. Both ends are the exported API a
// caller converting between formats would use.
func Test_WriteNullElements(t *testing.T) {
	for _, tc := range []struct {
		name string
		// a JSON member to splice into an otherwise minimal document
		member string
		// if set, output that must survive the element beside the null being skipped
		want string
	}{
		{name: "packages", member: `"packages": [null]`},
		{name: "files", member: `"files": [null]`},
		{name: "hasExtractedLicensingInfos", member: `"hasExtractedLicensingInfos": [null]`},
		{name: "annotations", member: `"annotations": [null]`},
		{
			// nested in a package, and with a second element so the sort
			// comparator runs before the loop does
			name: "package files",
			member: `"packages": [{"SPDXID": "SPDXRef-Package123", "name": "p1", "downloadLocation": "NOASSERTION",
				"files": [null, {"SPDXID": "SPDXRef-File123", "fileName": "/tmp/whatever.txt"}]}]`,
			want: "FileName: /tmp/whatever.txt",
		},
		{
			name: "package externalRefs",
			member: `"packages": [{"SPDXID": "SPDXRef-Package123", "name": "p1", "downloadLocation": "NOASSERTION",
				"externalRefs": [null, {"referenceCategory": "OTHER", "referenceType": "anything", "referenceLocator": "whatever"}]}]`,
			want: "ExternalRef: OTHER anything whatever",
		},
		{
			name: "file artifactOfs",
			member: `"files": [{"SPDXID": "SPDXRef-File123", "fileName": "/tmp/whatever.txt",
				"artifactOfs": [null, {"name": "project1"}]}]`,
			want: "ArtifactOfProjectName: project1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := fmt.Sprintf(`{
				"SPDXID": "SPDXRef-DOCUMENT",
				"spdxVersion": "SPDX-2.3",
				"creationInfo": {"created": "2018-10-10T06:20:00Z", "creators": ["Person: John Doe"]},
				"name": "null-elements",
				"dataLicense": "CC0-1.0",
				"documentNamespace": "https://example.com/null-elements",
				%s
			}`, tc.member)

			doc, err := spdxjson.Read(strings.NewReader(in))
			if err != nil {
				t.Fatalf("unable to read JSON document: %v", err)
			}

			var got bytes.Buffer
			if err := tagvalue.Write(doc, &got); err != nil {
				t.Errorf("unable to write tag-value document: %v", err)
			}
			if tc.want != "" && !strings.Contains(got.String(), tc.want) {
				t.Errorf("expected rendered document to contain %q, got:\n%s", tc.want, got.String())
			}
		})
	}
}

var ignores = []cmp.Option{
	cmpopts.IgnoreUnexported(spdx.Package{}),
	cmpopts.IgnoreFields(spdx.Document{}, "Snippets"),
	cmpopts.IgnoreFields(spdx.File{}, "Annotations"),
	cmpopts.IgnoreFields(spdx.Package{}, "IsFilesAnalyzedTagPresent", "Annotations"),
}
