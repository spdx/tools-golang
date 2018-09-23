// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"

	"github.com/swinslow/spdx-go/v0/spdx"
)

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
		default:
			return fmt.Errorf("unrecognized Creator type %v", subkey)
		}
	case "Created":
		ci.Created = value
	case "CreatorComment":
		ci.CreatorComment = value
	case "DocumentComment":
		ci.DocumentComment = value

	// tag for going on to package section
	case "PackageName":
		parser.st = psPackage2_1
		parser.pkg = &spdx.Package2_1{
			IsUnpackaged:              false,
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: false,
		}
		parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
		return parser.parsePairFromPackage2_1(tag, value)
	// tag for going on to _unpackaged_ file section
	case "FileName":
		// create an "unpackaged" Package structure
		parser.st = psFile2_1
		parser.pkg = &spdx.Package2_1{
			IsUnpackaged:              true,
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: false,
		}
		parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
		return parser.parsePairFromFile2_1(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_1
		return parser.parsePairFromOtherLicense2_1(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_1
		return parser.parsePairFromReview2_1(tag, value)
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
		return fmt.Errorf("received unknown tag %v in CreationInfo section", tag)
	}

	return nil
}
