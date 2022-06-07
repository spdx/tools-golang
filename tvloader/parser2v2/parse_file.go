// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	"github.com/spdx/tools-golang/utils"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_2) parsePairFromFile2_2(tag string, value string) error {
	// expire fileAOP for anything other than an AOPHomePage or AOPURI
	// (we'll actually handle the HomePage and URI further below)
	if tag != "ArtifactOfProjectHomePage" && tag != "ArtifactOfProjectURI" {
		parser.fileAOP = nil
	}

	switch tag {
	// tag for creating new file section
	case "FileName":
		// check if the previous file contained an spdx Id or not
		if parser.file != nil && parser.file.FileSPDXIdentifier == nullSpdxElementId2_2 {
			return fmt.Errorf("file with FileName %s does not have SPDX identifier", parser.file.FileName)
		}
		parser.file = &spdx.File2_2{}
		parser.file.FileName = value
	// tag for creating new package section and going back to parsing Package
	case "PackageName":
		parser.st = psPackage2_2
		// check if the previous file contained an spdx Id or not
		if parser.file != nil && parser.file.FileSPDXIdentifier == nullSpdxElementId2_2 {
			return fmt.Errorf("file with FileName %s does not have SPDX identifier", parser.file.FileName)
		}
		parser.file = nil
		return parser.parsePairFromPackage2_2(tag, value)
	// tag for going on to snippet section
	case "SnippetSPDXID":
		parser.st = psSnippet2_2
		return parser.parsePairFromSnippet2_2(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_2
		return parser.parsePairFromOtherLicense2_2(tag, value)
	// tags for file data
	case "SPDXID":
		err := parser.file.FileSPDXIdentifier.FromString(value)
		if err != nil {
			return err
		}
		if parser.pkg == nil {
			if parser.doc.Files == nil {
				parser.doc.Files = []*spdx.File2_2{}
			}
			parser.doc.Files = append(parser.doc.Files, parser.file)
		} else {
			if parser.pkg.Files == nil {
				parser.pkg.Files = []*spdx.File2_2{}
			}
			parser.pkg.Files = append(parser.pkg.Files, parser.file)
		}
	case "FileType":
		parser.file.FileTypes = append(parser.file.FileTypes, value)
	case "FileChecksum":
		subkey, subvalue, err := utils.ExtractSubs(value)
		if err != nil {
			return err
		}
		if parser.file.Checksums == nil {
			parser.file.Checksums = []spdx.Checksum{}
		}
		switch spdx.ChecksumAlgorithm(subkey) {
		case spdx.SHA1, spdx.SHA256, spdx.MD5:
			algorithm := spdx.ChecksumAlgorithm(subkey)
			parser.file.Checksums = append(parser.file.Checksums, spdx.Checksum{Algorithm: algorithm, Value: subvalue})
		default:
			return fmt.Errorf("got unknown checksum type %s", subkey)
		}
	case "LicenseConcluded":
		parser.file.LicenseConcluded = value
	case "LicenseInfoInFile":
		parser.file.LicenseInfoInFiles = append(parser.file.LicenseInfoInFiles, value)
	case "LicenseComments":
		parser.file.LicenseComments = value
	case "FileCopyrightText":
		parser.file.FileCopyrightText = value
	case "ArtifactOfProjectName":
		parser.fileAOP = &spdx.ArtifactOfProject2_2{}
		parser.file.ArtifactOfProjects = append(parser.file.ArtifactOfProjects, parser.fileAOP)
		parser.fileAOP.Name = value
	case "ArtifactOfProjectHomePage":
		if parser.fileAOP == nil {
			return fmt.Errorf("no current ArtifactOfProject found")
		}
		parser.fileAOP.HomePage = value
	case "ArtifactOfProjectURI":
		if parser.fileAOP == nil {
			return fmt.Errorf("no current ArtifactOfProject found")
		}
		parser.fileAOP.URI = value
	case "FileComment":
		parser.file.FileComment = value
	case "FileNotice":
		parser.file.FileNotice = value
	case "FileContributor":
		parser.file.FileContributors = append(parser.file.FileContributors, value)
	case "FileDependency":
		parser.file.FileDependencies = append(parser.file.FileDependencies, value)
	case "FileAttributionText":
		parser.file.FileAttributionTexts = append(parser.file.FileAttributionTexts, value)
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
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_2
		return parser.parsePairFromReview2_2(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in File section", tag)
	}

	return nil
}
