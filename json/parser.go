// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx_json

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

// Load2_2 takes in an io.Reader and returns an SPDX document.
func Load2_2(content io.Reader) (*spdx.Document2_2, error) {
	// convert io.Reader to a slice of bytes and call the parser
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(content)
	if err != nil {
		return nil, err
	}

	var doc spdx.Document2_2
	err = json.Unmarshal(buf.Bytes(), &doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
