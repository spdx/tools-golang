// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/utils"
)

func (parser *tvParser2_2) parsePairFromCreationInfo2_2(tag string, value string) error {
	// fail if not in Creation Info parser state
	if parser.st != psCreationInfo2_2 {
		return fmt.Errorf("got invalid state %v in parsePairFromCreationInfo2_2", parser.st)
	}

	// create an SPDX Creation Info data struct if we don't have one already
	if parser.doc.CreationInfo == nil {
		parser.doc.CreationInfo = &spdx.CreationInfo2_2{}
	}

	ci := parser.doc.CreationInfo
	switch tag {
	case "LicenseListVersion":
		ci.LicenseListVersion = value
	case "Creator":
		subkey, subvalue, err := utils.ExtractSubs(value)
		if err != nil {
			return err
		}

		creator := spdx.Creator{Creator: subvalue}
		switch subkey {
		case "Person", "Organization", "Tool":
			creator.CreatorType = subkey
		default:
			return fmt.Errorf("unrecognized Creator type %v", subkey)
		}

		ci.Creators = append(ci.Creators, creator)
	case "Created":
		ci.Created = value
	case "CreatorComment":
		ci.CreatorComment = value

	// tag for going on to package section
	case "PackageName":
		// error if last file does not have an identifier
		// this may be a null case: can we ever have a "last file" in
		// the "creation info" state? should go on to "file" state
		// even when parsing unpackaged files.
		if parser.file != nil && parser.file.FileSPDXIdentifier == nullSpdxElementId2_2 {
			return fmt.Errorf("file with FileName %s does not have SPDX identifier", parser.file.FileName)
		}
		parser.st = psPackage2_2
		parser.pkg = &spdx.Package2_2{
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: false,
		}
		return parser.parsePairFromPackage2_2(tag, value)
	// tag for going on to _unpackaged_ file section
	case "FileName":
		// leave pkg as nil, so that packages will be placed in Files
		parser.st = psFile2_2
		parser.pkg = nil
		return parser.parsePairFromFile2_2(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_2
		return parser.parsePairFromOtherLicense2_2(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_2
		return parser.parsePairFromReview2_2(tag, value)
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
		return fmt.Errorf("received unknown tag %v in CreationInfo section", tag)
	}

	return nil
}
