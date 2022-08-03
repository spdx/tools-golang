// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
)

type tvParser2_1 struct {
	// document into which data is being parsed
	doc *v2_1.Document

	// current parser state
	st tvParserState2_1

	// current SPDX item being filled in, if any
	pkg       *v2_1.Package
	pkgExtRef *v2_1.PackageExternalReference
	file      *v2_1.File
	fileAOP   *v2_1.ArtifactOfProject
	snippet   *v2_1.Snippet
	otherLic  *v2_1.OtherLicense
	rln       *v2_1.Relationship
	ann       *v2_1.Annotation
	rev       *v2_1.Review
	// don't need creation info pointer b/c only one,
	// and we can get to it via doc.CreationInfo
}

// parser state (SPDX document version 2.1)
type tvParserState2_1 int

const (
	// at beginning of document
	psStart2_1 tvParserState2_1 = iota

	// in document creation info section
	psCreationInfo2_1

	// in package data section
	psPackage2_1

	// in file data section (including "unpackaged" files)
	psFile2_1

	// in snippet data section (including "unpackaged" files)
	psSnippet2_1

	// in other license section
	psOtherLicense2_1

	// in review section
	psReview2_1
)

const nullSpdxElementId2_1 = common.ElementID("")
