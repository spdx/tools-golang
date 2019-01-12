// Package tvsaver is used to save tools-golang data structures
// as SPDX tag-value documents.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvsaver

import (
	"io"

	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/spdx/tools-golang/v0/tvsaver/saver2v1"
)

// Save2_1 takes an io.Writer and an SPDX Document (version 2.1),
// and writes it to the writer in tag-value format. It returns error
// if any error is encountered.
func Save2_1(doc *spdx.Document2_1, w io.Writer) error {
	return saver2v1.RenderDocument2_1(doc, w)
}
