// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v3

import (
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

type tvParser2_3 struct {
	// document into which data is being parsed
	doc *v2_3.Document

	// current parser state
	st tvParserState2_3

	// current SPDX item being filled in, if any
	pkg       *v2_3.Package
	pkgExtRef *v2_3.PackageExternalReference
	file      *v2_3.File
	fileAOP   *v2_3.ArtifactOfProject
	snippet   *v2_3.Snippet
	otherLic  *v2_3.OtherLicense
	rln       *v2_3.Relationship
	ann       *v2_3.Annotation
	rev       *v2_3.Review
	// don't need creation info pointer b/c only one,
	// and we can get to it via doc.CreationInfo
}

// parser state (SPDX document version 2.3)
type tvParserState2_3 int

const (
	// at beginning of document
	psStart2_3 tvParserState2_3 = iota

	// in document creation info section
	psCreationInfo2_3

	// in package data section
	psPackage2_3

	// in file data section (including "unpackaged" files)
	psFile2_3

	// in snippet data section (including "unpackaged" files)
	psSnippet2_3

	// in other license section
	psOtherLicense2_3

	// in review section
	psReview2_3
)

const nullSpdxElementId2_3 = common.ElementID("")
