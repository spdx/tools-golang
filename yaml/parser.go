// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx_yaml

import (
	"bytes"
	"io"

	"github.com/spdx/tools-golang/spdx/v2_2"
	"sigs.k8s.io/yaml"
)

// Load2_2 takes in an io.Reader and returns an SPDX document.
func Load2_2(content io.Reader) (*v2_2.Document, error) {
	// convert io.Reader to a slice of bytes and call the parser
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(content)
	if err != nil {
		return nil, err
	}

	var doc v2_2.Document
	err = yaml.Unmarshal(buf.Bytes(), &doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
