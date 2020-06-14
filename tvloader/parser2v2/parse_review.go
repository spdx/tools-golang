// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_2) parsePairFromReview2_2(tag string, value string) error {
	switch tag {
	// tag for creating new review section
	case "Reviewer":
		parser.rev = &spdx.Review2_2{}
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
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_2{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_2(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_2(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_2{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_2(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_2(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_2(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_2(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_2(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in Review section", tag)
	}

	return nil
}
