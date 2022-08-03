// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

type tvParser2_2 struct {
	// document into which data is being parsed
	doc *v2_2.Document

	// current parser state
	st tvParserState2_2

	// current SPDX item being filled in, if any
	pkg       *v2_2.Package
	pkgExtRef *v2_2.PackageExternalReference
	file      *v2_2.File
	fileAOP   *v2_2.ArtifactOfProject
	snippet   *v2_2.Snippet
	otherLic  *v2_2.OtherLicense
	rln       *v2_2.Relationship
	ann       *v2_2.Annotation
	rev       *v2_2.Review
	// don't need creation info pointer b/c only one,
	// and we can get to it via doc.CreationInfo
}

// parser state (SPDX document version 2.2)
type tvParserState2_2 int

const (
	// at beginning of document
	psStart2_2 tvParserState2_2 = iota

	// in document creation info section
	psCreationInfo2_2

	// in package data section
	psPackage2_2

	// in file data section (including "unpackaged" files)
	psFile2_2

	// in snippet data section (including "unpackaged" files)
	psSnippet2_2

	// in other license section
	psOtherLicense2_2

	// in review section
	psReview2_2
)

const nullSpdxElementId2_2 = common.ElementID("")
