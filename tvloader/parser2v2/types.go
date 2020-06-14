// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"github.com/spdx/tools-golang/spdx"
)

type tvParser2_2 struct {
	// document into which data is being parsed
	doc *spdx.Document2_2

	// current parser state
	st tvParserState2_2

	// current SPDX item being filled in, if any
	pkg       *spdx.Package2_2
	pkgExtRef *spdx.PackageExternalReference2_2
	file      *spdx.File2_2
	fileAOP   *spdx.ArtifactOfProject2_2
	snippet   *spdx.Snippet2_2
	otherLic  *spdx.OtherLicense2_2
	rln       *spdx.Relationship2_2
	ann       *spdx.Annotation2_2
	rev       *spdx.Review2_2
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
