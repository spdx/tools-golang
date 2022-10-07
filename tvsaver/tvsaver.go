// Package tvsaver is used to save tools-golang data structures
// as SPDX tag-value documents.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvsaver

import (
	"io"

	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
	"github.com/spdx/tools-golang/spdx/v2_3"
	"github.com/spdx/tools-golang/tvsaver/saver2v1"
	"github.com/spdx/tools-golang/tvsaver/saver2v2"
	"github.com/spdx/tools-golang/tvsaver/saver2v3"
)

// Save2_1 takes an io.Writer and an SPDX Document (version 2.1),
// and writes it to the writer in tag-value format. It returns error
// if any error is encountered.
func Save2_1(doc *v2_1.Document, w io.Writer) error {
	return saver2v1.RenderDocument2_1(doc, w)
}

// Save2_2 takes an io.Writer and an SPDX Document (version 2.2),
// and writes it to the writer in tag-value format. It returns error
// if any error is encountered.
func Save2_2(doc *v2_2.Document, w io.Writer) error {
	return saver2v2.RenderDocument2_2(doc, w)
}

// Save2_3 takes an io.Writer and an SPDX Document (version 2.3),
// and writes it to the writer in tag-value format. It returns error
// if any error is encountered.
func Save2_3(doc *v2_3.Document, w io.Writer) error {
	return saver2v3.RenderDocument2_3(doc, w)
}
