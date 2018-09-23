// Package tvloader is used to load and parse SPDX tag-value documents
// into spdx-go data structures.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"github.com/swinslow/spdx-go/v0/tvloader/parser2v1"
	"github.com/swinslow/spdx-go/v0/tvloader/reader"
)

type TVLoader struct {
	version   string
	reader    *reader.tvReader
	parser2_1 *parser2v1.parser2_1
}
