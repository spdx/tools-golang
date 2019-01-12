// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"

	"github.com/spdx/tools-golang/v0/spdx"
)

func (parser *tvParser2_1) parsePairFromReview2_1(tag string, value string) error {
	switch tag {
	// tag for creating new review section
	case "Reviewer":
		parser.rev = &spdx.Review2_1{}
		parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "Person":
			parser.rev.Reviewer = subvalue
			parser.rev.ReviewerType = "Person"
		case "Organization":
			parser.rev.Reviewer = subvalue
			parser.rev.ReviewerType = "Organization"
		case "Tool":
			parser.rev.Reviewer = subvalue
			parser.rev.ReviewerType = "Tool"
		default:
			return fmt.Errorf("unrecognized Reviewer type %v", subkey)
		}
	case "ReviewDate":
		parser.rev.ReviewDate = value
	case "ReviewComment":
		parser.rev.ReviewComment = value
	default:
		return fmt.Errorf("received unknown tag %v in Review section", tag)
	}

	return nil
}
