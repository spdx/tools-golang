// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package tvloader

import (
	"fmt"

	"github.com/swinslow/spdx-go/v0/spdx"
)

func (parser *tvParser2_1) parsePairFromFile2_1(tag string, value string) error {
	switch tag {
	// tag for creating new file section
	case "FileName":
		parser.file = &spdx.File2_1{}
		parser.pkg.Files = append(parser.pkg.Files, parser.file)
		parser.file.FileName = value
	// tag for creating new package section and going back to parsing Package
	case "PackageName":
		parser.st = psPackage2_1
		parser.file = nil
		return parser.parsePairFromPackage2_1(tag, value)
	// tag for going on to snippet section
	case "SnippetSPDXID":
		parser.st = psSnippet2_1
		return parser.parsePairFromSnippet2_1(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_1
		return parser.parsePairFromOtherLicense2_1(tag, value)
	// tags for package data
	case "SPDXID":
		parser.file.FileSPDXIdentifier = value
	case "FileType":
		parser.file.FileType = append(parser.file.FileType, value)
	case "FileChecksum":
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "SHA1":
			parser.file.FileChecksumSHA1 = subvalue
		case "SHA256":
			parser.file.FileChecksumSHA256 = subvalue
		case "MD5":
			parser.file.FileChecksumMD5 = subvalue
		default:
			return fmt.Errorf("got unknown checksum type %s", subkey)
		}
	case "LicenseConcluded":
		parser.file.LicenseConcluded = value
	case "LicenseInfoInFile":
		parser.file.LicenseInfoInFile = append(parser.file.LicenseInfoInFile, value)
	case "LicenseComments":
		parser.file.LicenseComments = value
	case "FileCopyrightText":
		parser.file.FileCopyrightText = value
	case "ArtifactOfProjectName":
		parser.file.ArtifactOfProjectName = append(parser.file.ArtifactOfProjectName, value)
	case "ArtifactOfProjectHomePage":
		parser.file.ArtifactOfProjectHomePage = append(parser.file.ArtifactOfProjectHomePage, value)
	case "ArtifactOfProjectURI":
		parser.file.ArtifactOfProjectURI = append(parser.file.ArtifactOfProjectURI, value)
	case "FileComment":
		parser.file.FileComment = value
	case "FileNotice":
		parser.file.FileNotice = value
	case "FileContributor":
		parser.file.FileContributor = append(parser.file.FileContributor, value)
	case "FileDependency":
		parser.file.FileDependencies = append(parser.file.FileDependencies, value)
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_1{}
		parser.file.Relationships = append(parser.file.Relationships, parser.rln)
		return parser.parsePairForRelationship2_1(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_1(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_1{}
		parser.file.Annotations = append(parser.file.Annotations, parser.ann)
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
		return fmt.Errorf("received unknown tag %v in File section", tag)
	}

	return nil
}
