// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v3

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx/v2_3"
)

func (parser *tvParser2_3) parsePairFromOtherLicense2_3(tag string, value string) error {
	switch tag {
	// tag for creating new other license section
	case "LicenseID":
		parser.otherLic = &v2_3.OtherLicense{}
		parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
		parser.otherLic.LicenseIdentifier = value
	case "ExtractedText":
		parser.otherLic.ExtractedText = value
	case "LicenseName":
		parser.otherLic.LicenseName = value
	case "LicenseCrossReference":
		parser.otherLic.LicenseCrossReferences = append(parser.otherLic.LicenseCrossReferences, value)
	case "LicenseComment":
		parser.otherLic.LicenseComment = value
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &v2_3.Relationship{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_3(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_3(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &v2_3.Annotation{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_3(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_3(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_3(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_3(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_3(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_3
		return parser.parsePairFromReview2_3(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in OtherLicense section", tag)
	}

	return nil
}
