// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/spdx/tools-golang/convert"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_1"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/spdx/tools-golang/spdx/v3/v3_0"
)

// Read takes an io.Reader and returns a fully-parsed current model SPDX Document
// or an error if any error is encountered.
func Read(content io.Reader) (*spdx.Document, error) {
	doc := spdx.Document{}
	err := ReadInto(content, &doc)
	return &doc, err
}

// ReadInto takes an io.Reader, reads in the SPDX document at the version provided
// and converts to the doc version
func ReadInto(content io.Reader, doc common.AnyDocument) error {
	if !convert.IsPtr(doc) {
		return fmt.Errorf("doc to read into must be a pointer")
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(content)
	if err != nil {
		return err
	}

	// Keep the document fields encoded while inspecting the version. Decoding
	// every package and file into interface values creates a second object tree
	// that is discarded as soon as the typed document is decoded below.
	var fields map[string]json.RawMessage
	err = json.Unmarshal(buf.Bytes(), &fields)
	if err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			return fmt.Errorf("not a valid SPDX JSON document")
		}
		return err
	}

	if fields == nil {
		return fmt.Errorf("not a valid SPDX JSON document")
	}

	version := stringField(fields, "spdxVersion")
	if version == "" {
		version = stringField(fields, "@context")
		if version != "" {
			extract := regexp.MustCompile(`https://spdx.org/rdf/(\d+(?:\.\d+)+)/spdx-context\.jsonld`)
			matches := extract.FindStringSubmatch(version)
			if len(matches) == 2 {
				version = matches[1]
			}
		}
	}

	if version == "" {
		return fmt.Errorf("JSON document does not contain spdxVersion field")
	}

	var data any
	switch version {
	case v2_1.Version:
		var doc v2_1.Document
		err = json.Unmarshal(buf.Bytes(), &doc)
		if err != nil {
			return err
		}
		data = doc
	case v2_2.Version:
		var doc v2_2.Document
		err = json.Unmarshal(buf.Bytes(), &doc)
		if err != nil {
			return err
		}
		data = doc
	case v2_3.Version:
		var doc v2_3.Document
		err = json.Unmarshal(buf.Bytes(), &doc)
		if err != nil {
			return err
		}
		data = doc
	case "3.0.0":
		fallthrough
	case v3_0.Version:
		// support older 3.0.x versions
		contents := buf.Bytes()
		if version != v3_0.Version {
			// The JSON-LD loader only knows the 3.0.1 context, so point @context at it.
			// Rewrite the parsed field, not the raw bytes, so escaped URLs and "3.0.0" in other values are handled correctly.
			fields["@context"], err = json.Marshal(fmt.Sprintf("https://spdx.org/rdf/%s/spdx-context.jsonld", v3_0.Version))
			if err != nil {
				return err
			}
			contents, err = json.Marshal(fields)
			if err != nil {
				return err
			}
		}
		var in v3_0.Document
		err = json.Unmarshal(contents, &in)
		if err != nil {
			return err
		}
		data = in
	default:
		return fmt.Errorf("unsupported SDPX version: %s", version)
	}

	return convert.Document(data, doc)
}

// stringField returns the named field as a string, or "" if it is missing or not a string.
func stringField(fields map[string]json.RawMessage, name string) string {
	var s string
	if raw, ok := fields[name]; ok {
		_ = json.Unmarshal(raw, &s)
	}
	return s
}
