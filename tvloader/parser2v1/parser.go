// Package parser2v1 contains functions to read, load and parse
// SPDX tag-value files.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/tvloader/reader"
)

// ParseTagValues takes a list of (tag, value) pairs, parses it and returns
// a pointer to a parsed SPDX Document.
func ParseTagValues(tvs []reader.TagValuePair) (*spdx.Document2_1, error) {
	parser := tvParser2_1{}
	for _, tv := range tvs {
		err := parser.parsePair2_1(tv.Tag, tv.Value)
		if err != nil {
			return nil, err
		}
	}

	return parser.doc, nil
}

func (parser *tvParser2_1) parsePair2_1(tag string, value string) error {
	switch parser.st {
	case psStart2_1:
		return parser.parsePairFromStart2_1(tag, value)
	case psCreationInfo2_1:
		return parser.parsePairFromCreationInfo2_1(tag, value)
	case psPackage2_1:
		return parser.parsePairFromPackage2_1(tag, value)
	case psFile2_1:
		return parser.parsePairFromFile2_1(tag, value)
	case psSnippet2_1:
		return parser.parsePairFromSnippet2_1(tag, value)
	case psOtherLicense2_1:
		return parser.parsePairFromOtherLicense2_1(tag, value)
	case psReview2_1:
		return parser.parsePairFromReview2_1(tag, value)
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
