// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"github.com/spdx/tools-golang/spdx"
)

type tvParser2_1 struct {
	// document into which data is being parsed
	doc *spdx.Document2_1

	// current parser state
	st tvParserState2_1

	// current SPDX item being filled in, if any
	pkg       *spdx.Package2_1
	pkgExtRef *spdx.PackageExternalReference2_1
	file      *spdx.File2_1
	fileAOP   *spdx.ArtifactOfProject2_1
	snippet   *spdx.Snippet2_1
	otherLic  *spdx.OtherLicense2_1
	rln       *spdx.Relationship2_1
	ann       *spdx.Annotation2_1
	rev       *spdx.Review2_1
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
