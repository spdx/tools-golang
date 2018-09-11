// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"fmt"
	"strings"

	"github.com/swinslow/spdx-go/v0/spdx"
)

type tvParser2_1 struct {
	// document into which data is being parsed
	doc *spdx.Document2_1

	// current parser state
	st tvParserState2_1

	// current SPDX item being filled in, if any
	pkg      *spdx.Package2_1
	file     *spdx.File2_1
	snippet  *spdx.Snippet2_1
	otherLic *spdx.OtherLicense2_1
	rln      *spdx.Relationship2_1
	ann      *spdx.Annotation2_1
	rev      *spdx.Review2_1
	// don't need creation info pointer b/c only one,
	// and we can get to it via doc.CreationInfo
}

// parser state (SPDX document version 2.1)
type tvParserState2_1 int

const (
	// at beginning of document
	psStart2_1 tvParserState2_1 = iota

	// in document creation info section
	psCreationInfo2_1

	// in package data section
	psPackage2_1

	// in file data section (including "unpackaged" files)
	psFile2_1

	// in snippet data section (including "unpackaged" files)
	psSnippet2_1

	// in other license section
	psOtherLicense2_1

	// at end of document
	psEnd2_1

	// error has been encountered
	psError2_1
)

// used to extract key / value from embedded substrings
// returns subkey, subvalue, nil if no error, or "", "", error otherwise
func extractSubs(value string) (string, string, error) {
	// parse the value to see if it's a valid subvalue format
	sp := strings.SplitN(value, ":", 2)
	if len(sp) == 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (no colon found)", value)
	}

	subkey := strings.TrimSpace(sp[0])
	subvalue := strings.TrimSpace(sp[1])

	// fail if there's another colon in the subvalue
	colon := strings.SplitN(subvalue, ":", 2)
	if len(colon) != 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (>1 colon found)", value)
	}

	return subkey, subvalue, nil
}

func (parser *tvParser2_1) parsePair2_1(tag string, value string) error {
	switch parser.st {
	case psStart2_1:
		return parser.parsePairFromStart2_1(tag, value)
	case psCreationInfo2_1:
		return parser.parsePairFromCreationInfo2_1(tag, value)
	// case psPackage2_1:
	// 	return parser.parsePairFromPackage2_1(tag, value)
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

func (parser *tvParser2_1) parsePairFromCreationInfo2_1(tag string, value string) error {
	// fail if not in Creation Info parser state
	if parser.st != psCreationInfo2_1 {
		return fmt.Errorf("Got invalid state %v in parsePairFromCreationInfo2_1", parser.st)
	}

	// create an SPDX Creation Info data struct if we don't have one already
	if parser.doc.CreationInfo == nil {
		parser.doc.CreationInfo = &spdx.CreationInfo2_1{}
	}

	ci := parser.doc.CreationInfo
	switch tag {
	// FIXME check for tags that go on to next section and change state
	// FIXME check for tags that add relationship / annotation, keeping state
	case "SPDXVersion":
		ci.SPDXVersion = value
	case "DataLicense":
		ci.DataLicense = value
	case "SPDXID":
		ci.SPDXIdentifier = value
	case "DocumentName":
		ci.DocumentName = value
	case "DocumentNamespace":
		ci.DocumentNamespace = value
	case "ExternalDocumentRef":
		ci.ExternalDocumentReferences = append(ci.ExternalDocumentReferences, value)
	case "LicenseListVersion":
		ci.LicenseListVersion = value
	case "Creator":
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "Person":
			ci.CreatorPersons = append(ci.CreatorPersons, subvalue)
		case "Organization":
			ci.CreatorOrganizations = append(ci.CreatorOrganizations, subvalue)
		case "Tool":
			ci.CreatorTools = append(ci.CreatorTools, subvalue)
		}
	case "Created":
		ci.Created = value
	case "CreatorComment":
		ci.CreatorComment = value
	case "DocumentComment":
		ci.DocumentComment = value
	}
	// FIXME complete and add default

	return nil
}
