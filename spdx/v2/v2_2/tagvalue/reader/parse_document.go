package reader

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func (parser *tvParser) parsePairForDocument(tag string, value string) error {
	if parser.doc == nil {
		return fmt.Errorf("no document struct created in parser ann pointer")
	}

	switch tag {
	case "DocumentComment":
		parser.doc = &v2_2.Document{}
		parser.doc.DocumentComment = value
	default:
		return fmt.Errorf("received unknown tag %v in Document section", tag)
	}

	return nil
}
