// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func (parser *tvParser) parsePairFromOtherLicense(tag string, value string) error {
	switch tag {
	// tag for creating new other license section
	case "LicenseID":
		parser.otherLic = &v2_2.OtherLicense{}
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
		parser.rln = &v2_2.Relationship{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &v2_2.Annotation{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation(tag, value)
	// tags for starting a new package, file or snippet section. Other License sections are not
	// required to be the last sections in a document: producers such as Yocto / OpenEmbedded emit
	// each package's LicenseRef extracted texts directly after the package, followed by the next
	// package. Mirror the File and Snippet states, which already allow these transitions.
	case "PackageName":
		parser.st = psPackage
		parser.file = nil
		parser.snippet = nil
		return parser.parsePairFromPackage(tag, value)
	case "FileName":
		parser.st = psFile
		parser.snippet = nil
		return parser.parsePairFromFile(tag, value)
	case "SnippetSPDXID":
		parser.st = psSnippet
		return parser.parsePairFromSnippet(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview
		return parser.parsePairFromReview(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in OtherLicense section", tag)
	}

	return nil
}
