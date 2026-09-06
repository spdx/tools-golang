// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/v2_1"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func TestReadVersion(t *testing.T) {
	for _, version := range []string{v2_1.Version, v2_2.Version, v2_3.Version} {
		t.Run(version, func(t *testing.T) {
			// Put the version after nested data to exercise inspection of the
			// whole document, including fields unknown to the SPDX model.
			input := fmt.Sprintf(`{"name":"example","SPDXID":"SPDXRef-DOCUMENT","unknown":{"spdxVersion":"ignored","nested":[1,true,null]},"packages":[{"name":"dependency","SPDXID":"SPDXRef-dependency"}],"spdxVersion":%q}`, version)
			doc, err := Read(strings.NewReader(input))
			require.NoError(t, err)
			require.Equal(t, v2_3.Version, doc.SPDXVersion)
			require.Equal(t, "example", doc.DocumentName)
			require.Len(t, doc.Packages, 1)
			require.Equal(t, "dependency", doc.Packages[0].PackageName)

			var original any
			switch version {
			case v2_1.Version:
				original = &v2_1.Document{}
			case v2_2.Version:
				original = &v2_2.Document{}
			case v2_3.Version:
				original = &v2_3.Document{}
			}
			require.NoError(t, ReadInto(strings.NewReader(input), original))
		})
	}
}

func TestReadVersionErrors(t *testing.T) {
	for _, tc := range []struct{ name, input, message string }{
		{"null", `null`, "not a valid SPDX JSON document"},
		{"array", `[]`, "not a valid SPDX JSON document"},
		{"string", `"document"`, "not a valid SPDX JSON document"},
		{"number", `42`, "not a valid SPDX JSON document"},
		{"boolean", `true`, "not a valid SPDX JSON document"},
		{"missing", `{}`, "does not contain spdxVersion"},
		{"case sensitive", `{"SPDXVersion":"SPDX-2.3"}`, "does not contain spdxVersion"},
		{"unsupported", `{"spdxVersion":"SPDX-9.9"}`, "unsupported SDPX version"},
		{"null version", `{"spdxVersion":null}`, "does not contain spdxVersion"},
		{"numeric version", `{"spdxVersion":2.3}`, "does not contain spdxVersion"},
		{"object version", `{"spdxVersion":{}}`, "does not contain spdxVersion"},
		{"array version", `{"spdxVersion":[]}`, "does not contain spdxVersion"},
		{"malformed", `{"spdxVersion":"SPDX-2.3","packages":[}`, "invalid character"},
		{"trailing data", `{"spdxVersion":"SPDX-2.3"} {}`, "invalid character"},
		{"last duplicate wins", `{"spdxVersion":"SPDX-2.3","spdxVersion":"SPDX-9.9"}`, "unsupported SDPX version"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := spdx.Document{DocumentName: "unchanged"}
			err := ReadInto(strings.NewReader(tc.input), &doc)
			require.ErrorContains(t, err, tc.message)
			require.Equal(t, "unchanged", doc.DocumentName)
		})
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestReadIOError(t *testing.T) {
	want := errors.New("reader failed")
	_, err := Read(failingReader{want})
	require.ErrorIs(t, err, want)
}
