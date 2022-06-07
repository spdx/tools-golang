// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"
	"github.com/spdx/tools-golang/utils"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_1) parsePairFromReview2_1(tag string, value string) error {
	switch tag {
	// tag for creating new review section
	case "Reviewer":
		parser.rev = &spdx.Review2_1{}
		parser.doc.Reviews = append(parser.doc.Reviews, parser.rev)
		subkey, subvalue, err := utils.ExtractSubs(value)
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
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_1{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_1(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_1(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_1{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_1(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_1(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in Review section", tag)
	}

	return nil
}
