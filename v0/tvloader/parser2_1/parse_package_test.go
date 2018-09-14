// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"testing"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// ===== Parser package section state change tests =====
func TestParser2_1PackageStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	// create the first package
	pkgOldName := "p1"

	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: pkgOldName},
	}
	pkgOld := parser.pkg
	parser.doc.Packages = append(parser.doc.Packages, pkgOld)
	// the Document's Packages should have this one only
	if parser.doc.Packages[0] != pkgOld {
		t.Errorf("Expected package %v in Packages[0], got %v", pkgOld, parser.doc.Packages[0])
	}
	if parser.doc.Packages[0].PackageName != pkgOldName {
		t.Errorf("expected package name %s in Packages[0], got %s", pkgOldName, parser.doc.Packages[0].PackageName)
	}

	// now add a new package
	pkgName := "p2"
	err := parser.parsePair2_1("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}
	// and a package should be created
	if parser.pkg == nil {
		t.Fatalf("parser didn't create new package")
	}
	// and the package name should be as expected
	if parser.pkg.PackageName != pkgName {
		t.Errorf("expected package name %s, got %s", pkgName, parser.pkg.PackageName)
	}
	// and the package should _not_ be an "unpackaged" placeholder
	if parser.pkg.IsUnpackaged == true {
		t.Errorf("package incorrectly has IsUnpackaged flag set")
	}
	// and the Document's Packages should be of size 2 and have these two
	if parser.doc.Packages[0] != pkgOld {
		t.Errorf("Expected package %v in Packages[0], got %v", pkgOld, parser.doc.Packages[0])
	}
	if parser.doc.Packages[0].PackageName != pkgOldName {
		t.Errorf("expected package name %s in Packages[0], got %s", pkgOldName, parser.doc.Packages[0].PackageName)
	}
	if parser.doc.Packages[1] != parser.pkg {
		t.Errorf("Expected package %v in Packages[1], got %v", parser.pkg, parser.doc.Packages[1])
	}
	if parser.doc.Packages[1].PackageName != pkgName {
		t.Errorf("expected package name %s in Packages[1], got %s", pkgName, parser.doc.Packages[1].PackageName)
	}
}

// func TestParser2_1CIMovesToFileAfterParsingFileNameTag(t *testing.T) {
// 	parser := tvParser2_1{
// 		doc: &spdx.Document2_1{},
// 		st:  psCreationInfo2_1,
// 	}
// 	err := parser.parsePair2_1("FileName", "testFile")
// 	if err != nil {
// 		t.Errorf("got error when calling parsePair2_1: %v", err)
// 	}
// 	// state should be correct
// 	if parser.st != psFile2_1 {
// 		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
// 	}
// 	// and current package should be an "unpackaged" placeholder
// 	if parser.pkg == nil {
// 		t.Fatalf("parser didn't create placeholder package")
// 	}
// 	if !parser.pkg.IsUnpackaged {
// 		t.Errorf("placeholder package is not set as unpackaged")
// 	}
// }

// func TestParser2_1CIMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
// 	parser := tvParser2_1{
// 		doc: &spdx.Document2_1{},
// 		st:  psCreationInfo2_1,
// 	}
// 	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
// 	if err != nil {
// 		t.Errorf("got error when calling parsePair2_1: %v", err)
// 	}
// 	if parser.st != psOtherLicense2_1 {
// 		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense2_1)
// 	}
// }

// func TestParser2_1CIStaysAfterParsingRelationshipTags(t *testing.T) {
// 	parser := tvParser2_1{
// 		doc: &spdx.Document2_1{},
// 		st:  psCreationInfo2_1,
// 	}

// 	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-else")
// 	if err != nil {
// 		t.Errorf("got error when calling parsePair2_1: %v", err)
// 	}
// 	if parser.st != psCreationInfo2_1 {
// 		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
// 	}

// 	err = parser.parsePair2_1("RelationshipComment", "blah")
// 	if err != nil {
// 		t.Errorf("got error when calling parsePair2_1: %v", err)
// 	}
// 	if parser.st != psCreationInfo2_1 {
// 		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
// 	}
// }
