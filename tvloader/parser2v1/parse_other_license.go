// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_1) parsePairFromOtherLicense2_1(tag string, value string) error {
	switch tag {
	// tag for creating new other license section
	case "LicenseID":
		parser.otherLic = &spdx.OtherLicense2_1{}
		parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
		parser.otherLic.Identifier = value
	case "ExtractedText":
		parser.otherLic.ExtractedText = value
	case "LicenseName":
		parser.otherLic.Name = value
	case "LicenseCrossReference":
		parser.otherLic.CrossReferences = append(parser.otherLic.CrossReferences, value)
	case "LicenseComment":
		parser.otherLic.Comment = value
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
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_1
		return parser.parsePairFromReview2_1(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in OtherLicense section", tag)
	}

	return nil
}
