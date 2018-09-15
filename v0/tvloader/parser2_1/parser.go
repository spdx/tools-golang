// Package tvloader contains functions to read, load and parse
// SPDX tag-value files.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"fmt"

	"github.com/swinslow/spdx-go/v0/spdx"
)

func (parser *tvParser2_1) parsePair2_1(tag string, value string) error {
	switch parser.st {
	case psStart2_1:
		return parser.parsePairFromStart2_1(tag, value)
	case psCreationInfo2_1:
		return parser.parsePairFromCreationInfo2_1(tag, value)
	case psPackage2_1:
		return parser.parsePairFromPackage2_1(tag, value)
	// case psFile2_1:
	// 	return parser.parsePairFromFile2_1(tag, value)
	// case psSnippet2_1:
	// 	return parser.parsePairFromSnippet2_1(tag, value)
	// case psOtherLicense2_1:
	// 	return parser.parsePairFromOtherLicense2_1(tag, value)
	case psEnd2_1:
		return fmt.Errorf("Already in parser end state when parsing (%s, %s)", tag, value)
	case psError2_1:
		return fmt.Errorf("In parser error state when parsing (%s, %s)", tag, value)
	default:
		return fmt.Errorf("Parser state %v not recognized when parsing (%s, %s)", parser.st, tag, value)
	}
}

func (parser *tvParser2_1) parsePairFromStart2_1(tag string, value string) error {
	// fail if not in Start parser state
	if parser.st != psStart2_1 {
		return fmt.Errorf("Got invalid state %v in parsePairFromStart2_1", parser.st)
	}

	// create an SPDX Document data struct if we don't have one already
	if parser.doc == nil {
		parser.doc = &spdx.Document2_1{}
	}

	// move to Creation Info parser state
	parser.st = psCreationInfo2_1

	// and ask Creation Info subfunc to parse
	return parser.parsePairFromCreationInfo2_1(tag, value)
}

func (parser *tvParser2_1) parsePairFromFile2_1(tag string, value string) error {
	return nil
}

func (parser *tvParser2_1) parsePairFromOtherLicense2_1(tag string, value string) error {
	return nil
}
