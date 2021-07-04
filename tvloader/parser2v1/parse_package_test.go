// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Parser package section state change tests =====
func TestParser2_1PackageStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	// create the first package
	pkgOldName := "p1"

	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: pkgOldName, PackageSPDXIdentifier: "p1"},
	}
	pkgOld := parser.pkg
	parser.doc.Packages["p1"] = pkgOld
	// the Document's Packages should have this one only
	if parser.doc.Packages["p1"] != pkgOld {
		t.Errorf("expected package %v, got %v", pkgOld, parser.doc.Packages["p1"])
	}
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}

	// now add a new package
	pkgName := "p2"
	err := parser.parsePair2_1("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}
	// and a package should be created
	if parser.pkg == nil {
		t.Fatalf("parser didn't create new package")
	}
	// and it should not be pkgOld
	if parser.pkg == pkgOld {
		t.Errorf("expected new package, got pkgOld")
	}
	// and the package name should be as expected
	if parser.pkg.PackageName != pkgName {
		t.Errorf("expected package name %s, got %s", pkgName, parser.pkg.PackageName)
	}
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the Document's Packages should still be of size 1 and have pkgOld only
	if parser.doc.Packages["p1"] != pkgOld {
		t.Errorf("Expected package %v, got %v", pkgOld, parser.doc.Packages["p1"])
	}
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}
}

func TestParser2_1PackageStartsNewPackageAfterParsingPackageNameTagWhileInUnpackaged(t *testing.T) {
	// pkg is nil, so that Files appearing before the first PackageName tag
	// are added to UnpackagedFiles instead of Packages
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: nil,
	}
	// the Document's Packages should be empty
	if len(parser.doc.Packages) != 0 {
		t.Errorf("Expected zero packages, got %d", len(parser.doc.Packages))
	}

	// now add a new package
	pkgName := "p2"
	err := parser.parsePair2_1("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}
	// and a package should be created
	if parser.pkg == nil {
		t.Fatalf("parser didn't create new package")
	}
	// and the package name should be as expected
	if parser.pkg.PackageName != pkgName {
		t.Errorf("expected package name %s, got %s", pkgName, parser.pkg.PackageName)
	}
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the Document's Packages should be of size 0, because the prior was
	// unpackaged files and this one won't be added until an SPDXID is seen
	if len(parser.doc.Packages) != 0 {
		t.Errorf("Expected %v packages in doc, got %v", 0, len(parser.doc.Packages))
	}
}

func TestParser2_1PackageMovesToFileAfterParsingFileNameTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg
	pkgCurrent := parser.pkg

	err := parser.parsePair2_1("FileName", "testFile")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psFile2_1 {
		t.Errorf("expected state to be %v, got %v", psFile2_1, parser.st)
	}
	// and current package should remain what it was
	if parser.pkg != pkgCurrent {
		t.Fatalf("expected package to remain %v, got %v", pkgCurrent, parser.pkg)
	}
}

func TestParser2_1PackageMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psOtherLicense2_1 {
		t.Errorf("expected state to be %v, got %v", psOtherLicense2_1, parser.st)
	}
}

func TestParser2_1PackageMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePair2_1("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("expected state to be %v, got %v", psReview2_1, parser.st)
	}
}

func TestParser2_1PackageStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePair2_1("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should remain unchanged
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}
}

func TestParser2_1PackageStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePair2_1("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}

	err = parser.parsePair2_1("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}

	err = parser.parsePair2_1("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}

	err = parser.parsePair2_1("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}

	err = parser.parsePair2_1("AnnotationComment", "i guess i had something to say about this package")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psPackage2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_1)
	}
}

// ===== Package data section tests =====
func TestParser2_1CanParsePackageTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// should not yet be in Packages map, b/c no SPDX identifier
	if len(parser.doc.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(parser.doc.Packages))
	}

	// Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageName != "p1" {
		t.Errorf("got %v for PackageName", parser.pkg.PackageName)
	}
	// still should not yet be in Packages map, b/c no SPDX identifier
	if len(parser.doc.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(parser.doc.Packages))
	}

	// Package SPDX Identifier
	err = parser.parsePairFromPackage2_1("SPDXID", "SPDXRef-p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// "SPDXRef-" prefix should be removed from the item
	if parser.pkg.PackageSPDXIdentifier != "p1" {
		t.Errorf("got %v for PackageSPDXIdentifier", parser.pkg.PackageSPDXIdentifier)
	}
	// and it should now be added to the Packages map
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}
	if parser.doc.Packages["p1"] != parser.pkg {
		t.Errorf("expected to point to parser.pkg, got %v", parser.doc.Packages["p1"])
	}

	// Package Version
	err = parser.parsePairFromPackage2_1("PackageVersion", "2.1.1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageVersion != "2.1.1" {
		t.Errorf("got %v for PackageVersion", parser.pkg.PackageVersion)
	}

	// Package File Name
	err = parser.parsePairFromPackage2_1("PackageFileName", "p1-2.1.1.tar.gz")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageFileName != "p1-2.1.1.tar.gz" {
		t.Errorf("got %v for PackageFileName", parser.pkg.PackageFileName)
	}

	// Package Supplier
	// SKIP -- separate tests for subvalues below

	// Package Originator
	// SKIP -- separate tests for subvalues below

	// Package Download Location
	err = parser.parsePairFromPackage2_1("PackageDownloadLocation", "https://example.com/whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageDownloadLocation != "https://example.com/whatever" {
		t.Errorf("got %v for PackageDownloadLocation", parser.pkg.PackageDownloadLocation)
	}

	// Files Analyzed
	err = parser.parsePairFromPackage2_1("FilesAnalyzed", "false")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.FilesAnalyzed != false {
		t.Errorf("got %v for FilesAnalyzed", parser.pkg.FilesAnalyzed)
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != true {
		t.Errorf("got %v for IsFilesAnalyzedTagPresent", parser.pkg.IsFilesAnalyzedTagPresent)
	}

	// Package Verification Code
	// SKIP -- separate tests for "excludes", or not, below

	// Package Checksums
	codeSha1 := "85ed0817af83a24ad8da68c2b5094de69833983c"
	sumSha1 := "SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c"
	codeSha256 := "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"
	sumSha256 := "SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"
	codeMd5 := "624c1abb3664f4b35547e7c73864ad24"
	sumMd5 := "MD5: 624c1abb3664f4b35547e7c73864ad24"
	err = parser.parsePairFromPackage2_1("PackageChecksum", sumSha1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromPackage2_1("PackageChecksum", sumSha256)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromPackage2_1("PackageChecksum", sumMd5)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageChecksumSHA1 != codeSha1 {
		t.Errorf("expected %s for PackageChecksumSHA1, got %s", codeSha1, parser.pkg.PackageChecksumSHA1)
	}
	if parser.pkg.PackageChecksumSHA256 != codeSha256 {
		t.Errorf("expected %s for PackageChecksumSHA256, got %s", codeSha256, parser.pkg.PackageChecksumSHA256)
	}
	if parser.pkg.PackageChecksumMD5 != codeMd5 {
		t.Errorf("expected %s for PackageChecksumMD5, got %s", codeMd5, parser.pkg.PackageChecksumMD5)
	}

	// Package Home Page
	err = parser.parsePairFromPackage2_1("PackageHomePage", "https://example.com/whatever2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageHomePage != "https://example.com/whatever2" {
		t.Errorf("got %v for PackageHomePage", parser.pkg.PackageHomePage)
	}

	// Package Source Info
	err = parser.parsePairFromPackage2_1("PackageSourceInfo", "random comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSourceInfo != "random comment" {
		t.Errorf("got %v for PackageSourceInfo", parser.pkg.PackageSourceInfo)
	}

	// Package License Concluded
	err = parser.parsePairFromPackage2_1("PackageLicenseConcluded", "Apache-2.0 OR GPL-2.0-or-later")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageLicenseConcluded != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("got %v for PackageLicenseConcluded", parser.pkg.PackageLicenseConcluded)
	}

	// All Licenses Info From Files
	lics := []string{
		"Apache-2.0",
		"GPL-2.0-or-later",
		"CC0-1.0",
	}
	for _, lic := range lics {
		err = parser.parsePairFromPackage2_1("PackageLicenseInfoFromFiles", lic)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, licWant := range lics {
		flagFound := false
		for _, licCheck := range parser.pkg.PackageLicenseInfoFromFiles {
			if licWant == licCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in PackageLicenseInfoFromFiles", licWant)
		}
	}
	if len(lics) != len(parser.pkg.PackageLicenseInfoFromFiles) {
		t.Errorf("expected %d licenses in PackageLicenseInfoFromFiles, got %d", len(lics),
			len(parser.pkg.PackageLicenseInfoFromFiles))
	}

	// Package License Declared
	err = parser.parsePairFromPackage2_1("PackageLicenseDeclared", "Apache-2.0 OR GPL-2.0-or-later")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageLicenseDeclared != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("got %v for PackageLicenseDeclared", parser.pkg.PackageLicenseDeclared)
	}

	// Package License Comments
	err = parser.parsePairFromPackage2_1("PackageLicenseComments", "this is a license comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageLicenseComments != "this is a license comment" {
		t.Errorf("got %v for PackageLicenseComments", parser.pkg.PackageLicenseComments)
	}

	// Package Copyright Text
	err = parser.parsePairFromPackage2_1("PackageCopyrightText", "Copyright (c) me myself and i")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageCopyrightText != "Copyright (c) me myself and i" {
		t.Errorf("got %v for PackageCopyrightText", parser.pkg.PackageCopyrightText)
	}

	// Package Summary
	err = parser.parsePairFromPackage2_1("PackageSummary", "i wrote this package")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSummary != "i wrote this package" {
		t.Errorf("got %v for PackageSummary", parser.pkg.PackageSummary)
	}

	// Package Description
	err = parser.parsePairFromPackage2_1("PackageDescription", "i wrote this package a lot")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageDescription != "i wrote this package a lot" {
		t.Errorf("got %v for PackageDescription", parser.pkg.PackageDescription)
	}

	// Package Comment
	err = parser.parsePairFromPackage2_1("PackageComment", "i scanned this package")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageComment != "i scanned this package" {
		t.Errorf("got %v for PackageComment", parser.pkg.PackageComment)
	}

	// Package External References and Comments
	ref1 := "SECURITY cpe23Type cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"
	ref1Category := "SECURITY"
	ref1Type := "cpe23Type"
	ref1Locator := "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"
	ref1Comment := "this is comment #1"
	ref2 := "OTHER LocationRef-acmeforge acmecorp/acmenator/4.1.3alpha"
	ref2Category := "OTHER"
	ref2Type := "LocationRef-acmeforge"
	ref2Locator := "acmecorp/acmenator/4.1.3alpha"
	ref2Comment := "this is comment #2"
	err = parser.parsePairFromPackage2_1("ExternalRef", ref1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.pkg.PackageExternalReferences) != 1 {
		t.Errorf("expected 1 external reference, got %d", len(parser.pkg.PackageExternalReferences))
	}
	if parser.pkgExtRef == nil {
		t.Errorf("expected non-nil pkgExtRef, got nil")
	}
	if parser.pkg.PackageExternalReferences[0] == nil {
		t.Errorf("expected non-nil PackageExternalReferences[0], got nil")
	}
	if parser.pkgExtRef != parser.pkg.PackageExternalReferences[0] {
		t.Errorf("expected pkgExtRef to match PackageExternalReferences[0], got no match")
	}
	err = parser.parsePairFromPackage2_1("ExternalRefComment", ref1Comment)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromPackage2_1("ExternalRef", ref2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.pkg.PackageExternalReferences) != 2 {
		t.Errorf("expected 2 external references, got %d", len(parser.pkg.PackageExternalReferences))
	}
	if parser.pkgExtRef == nil {
		t.Errorf("expected non-nil pkgExtRef, got nil")
	}
	if parser.pkg.PackageExternalReferences[1] == nil {
		t.Errorf("expected non-nil PackageExternalReferences[1], got nil")
	}
	if parser.pkgExtRef != parser.pkg.PackageExternalReferences[1] {
		t.Errorf("expected pkgExtRef to match PackageExternalReferences[1], got no match")
	}
	err = parser.parsePairFromPackage2_1("ExternalRefComment", ref2Comment)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// finally, check these values
	gotRef1 := parser.pkg.PackageExternalReferences[0]
	if gotRef1.Category != ref1Category {
		t.Errorf("expected ref1 category to be %s, got %s", gotRef1.Category, ref1Category)
	}
	if gotRef1.RefType != ref1Type {
		t.Errorf("expected ref1 type to be %s, got %s", gotRef1.RefType, ref1Type)
	}
	if gotRef1.Locator != ref1Locator {
		t.Errorf("expected ref1 locator to be %s, got %s", gotRef1.Locator, ref1Locator)
	}
	if gotRef1.ExternalRefComment != ref1Comment {
		t.Errorf("expected ref1 comment to be %s, got %s", gotRef1.ExternalRefComment, ref1Comment)
	}
	gotRef2 := parser.pkg.PackageExternalReferences[1]
	if gotRef2.Category != ref2Category {
		t.Errorf("expected ref2 category to be %s, got %s", gotRef2.Category, ref2Category)
	}
	if gotRef2.RefType != ref2Type {
		t.Errorf("expected ref2 type to be %s, got %s", gotRef2.RefType, ref2Type)
	}
	if gotRef2.Locator != ref2Locator {
		t.Errorf("expected ref2 locator to be %s, got %s", gotRef2.Locator, ref2Locator)
	}
	if gotRef2.ExternalRefComment != ref2Comment {
		t.Errorf("expected ref2 comment to be %s, got %s", gotRef2.ExternalRefComment, ref2Comment)
	}

}

func TestParser2_1CanParsePackageSupplierPersonTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Supplier: Person
	err := parser.parsePairFromPackage2_1("PackageSupplier", "Person: John Doe")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSupplierPerson != "John Doe" {
		t.Errorf("got %v for PackageSupplierPerson", parser.pkg.PackageSupplierPerson)
	}
}

func TestParser2_1CanParsePackageSupplierOrganizationTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Supplier: Organization
	err := parser.parsePairFromPackage2_1("PackageSupplier", "Organization: John Doe, Inc.")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSupplierOrganization != "John Doe, Inc." {
		t.Errorf("got %v for PackageSupplierOrganization", parser.pkg.PackageSupplierOrganization)
	}
}

func TestParser2_1CanParsePackageSupplierNOASSERTIONTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Supplier: NOASSERTION
	err := parser.parsePairFromPackage2_1("PackageSupplier", "NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSupplierNOASSERTION != true {
		t.Errorf("got false for PackageSupplierNOASSERTION")
	}
}

func TestParser2_1CanParsePackageOriginatorPersonTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Originator: Person
	err := parser.parsePairFromPackage2_1("PackageOriginator", "Person: John Doe")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageOriginatorPerson != "John Doe" {
		t.Errorf("got %v for PackageOriginatorPerson", parser.pkg.PackageOriginatorPerson)
	}
}

func TestParser2_1CanParsePackageOriginatorOrganizationTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Originator: Organization
	err := parser.parsePairFromPackage2_1("PackageOriginator", "Organization: John Doe, Inc.")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageOriginatorOrganization != "John Doe, Inc." {
		t.Errorf("got %v for PackageOriginatorOrganization", parser.pkg.PackageOriginatorOrganization)
	}
}

func TestParser2_1CanParsePackageOriginatorNOASSERTIONTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Originator: NOASSERTION
	err := parser.parsePairFromPackage2_1("PackageOriginator", "NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageOriginatorNOASSERTION != true {
		t.Errorf("got false for PackageOriginatorNOASSERTION")
	}
}

func TestParser2_1CanParsePackageVerificationCodeTagWithExcludes(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Verification Code with excludes parenthetical
	code := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	fileName := "./package.spdx"
	fullCodeValue := "d6a770ba38583ed4bb4525bd96e50461655d2758 (excludes: ./package.spdx)"
	err := parser.parsePairFromPackage2_1("PackageVerificationCode", fullCodeValue)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageVerificationCode != code {
		t.Errorf("got %v for PackageVerificationCode", parser.pkg.PackageVerificationCode)
	}
	if parser.pkg.PackageVerificationCodeExcludedFile != fileName {
		t.Errorf("got %v for PackageVerificationCodeExcludedFile", parser.pkg.PackageVerificationCodeExcludedFile)
	}

}

func TestParser2_1CanParsePackageVerificationCodeTagWithoutExcludes(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	// Package Verification Code without excludes parenthetical
	code := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	err := parser.parsePairFromPackage2_1("PackageVerificationCode", code)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageVerificationCode != code {
		t.Errorf("got %v for PackageVerificationCode", parser.pkg.PackageVerificationCode)
	}
	if parser.pkg.PackageVerificationCodeExcludedFile != "" {
		t.Errorf("got %v for PackageVerificationCodeExcludedFile", parser.pkg.PackageVerificationCodeExcludedFile)
	}

}

func TestPackageExternalRefPointerChangesAfterTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	ref1 := "SECURITY cpe23Type cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"
	err := parser.parsePairFromPackage2_1("ExternalRef", ref1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef == nil {
		t.Errorf("expected non-nil external reference pointer, got nil")
	}

	// now, a comment; pointer should go away
	err = parser.parsePairFromPackage2_1("ExternalRefComment", "whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef != nil {
		t.Errorf("expected nil external reference pointer, got non-nil")
	}

	ref2 := "Other LocationRef-something https://example.com/whatever"
	err = parser.parsePairFromPackage2_1("ExternalRef", ref2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef == nil {
		t.Errorf("expected non-nil external reference pointer, got nil")
	}

	// and some other random tag makes the pointer go away too
	err = parser.parsePairFromPackage2_1("PackageSummary", "whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef != nil {
		t.Errorf("expected nil external reference pointer, got non-nil")
	}
}

func TestParser2_1PackageCreatesRelationshipInDocument(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePair2_1("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-whatever")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.rln == nil {
		t.Fatalf("parser didn't create and point to Relationship struct")
	}
	if parser.rln != parser.doc.Relationships[0] {
		t.Errorf("pointer to new Relationship doesn't match idx 0 for doc.Relationships[]")
	}
}

func TestParser2_1PackageCreatesAnnotationInDocument(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePair2_1("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.ann == nil {
		t.Fatalf("parser didn't create and point to Annotation struct")
	}
	if parser.ann != parser.doc.Annotations[0] {
		t.Errorf("pointer to new Annotation doesn't match idx 0 for doc.Annotations[]")
	}
}

func TestParser2_1PackageUnknownTagFails(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages["p1"] = parser.pkg

	err := parser.parsePairFromPackage2_1("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestParser2_1FailsIfInvalidSPDXIDInPackageSection(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid ID format
	err = parser.parsePairFromPackage2_1("SPDXID", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfInvalidPackageSupplierFormat(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid supplier format
	err = parser.parsePairFromPackage2_1("PackageSupplier", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfUnknownPackageSupplierType(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid supplier type
	err = parser.parsePairFromPackage2_1("PackageSupplier", "whoops: John Doe")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfInvalidPackageOriginatorFormat(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid originator format
	err = parser.parsePairFromPackage2_1("PackageOriginator", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfUnknownPackageOriginatorType(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid originator type
	err = parser.parsePairFromPackage2_1("PackageOriginator", "whoops: John Doe")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1SetsFilesAnalyzedTagsCorrectly(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// set tag
	err = parser.parsePairFromPackage2_1("FilesAnalyzed", "true")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected %v, got %v", true, parser.pkg.FilesAnalyzed)
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != true {
		t.Errorf("expected %v, got %v", true, parser.pkg.IsFilesAnalyzedTagPresent)
	}
}

func TestParser2_1FailsIfInvalidPackageChecksumFormat(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid checksum format
	err = parser.parsePairFromPackage2_1("PackageChecksum", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfInvalidPackageChecksumType(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid checksum type
	err = parser.parsePairFromPackage2_1("PackageChecksum", "whoops: blah")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfInvalidExternalRefFormat(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid external ref format
	err = parser.parsePairFromPackage2_1("ExternalRef", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfExternalRefCommentBeforeExternalRef(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// external ref comment before external ref
	err = parser.parsePairFromPackage2_1("ExternalRefComment", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// ===== Helper function tests =====

func TestCanCheckAndExtractExcludesFilenameAndCode(t *testing.T) {
	code := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	fileName := "./package.spdx"
	fullCodeValue := "d6a770ba38583ed4bb4525bd96e50461655d2758 (excludes: ./package.spdx)"

	gotCode, gotFileName := extractCodeAndExcludes(fullCodeValue)
	if gotCode != code {
		t.Errorf("got %v for gotCode", gotCode)
	}
	if gotFileName != fileName {
		t.Errorf("got %v for gotFileName", gotFileName)
	}
}

func TestCanExtractPackageExternalReference(t *testing.T) {
	ref1 := "SECURITY cpe23Type cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"
	category := "SECURITY"
	refType := "cpe23Type"
	location := "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"

	gotCategory, gotRefType, gotLocation, err := extractPackageExternalReference(ref1)
	if err != nil {
		t.Errorf("got non-nil error: %v", err)
	}
	if gotCategory != category {
		t.Errorf("expected category %s, got %s", category, gotCategory)
	}
	if gotRefType != refType {
		t.Errorf("expected refType %s, got %s", refType, gotRefType)
	}
	if gotLocation != location {
		t.Errorf("expected location %s, got %s", location, gotLocation)
	}
}

func TestCanExtractPackageExternalReferenceWithExtraWhitespace(t *testing.T) {
	ref1 := "  SECURITY    \t cpe23Type   cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:* \t "
	category := "SECURITY"
	refType := "cpe23Type"
	location := "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"

	gotCategory, gotRefType, gotLocation, err := extractPackageExternalReference(ref1)
	if err != nil {
		t.Errorf("got non-nil error: %v", err)
	}
	if gotCategory != category {
		t.Errorf("expected category %s, got %s", category, gotCategory)
	}
	if gotRefType != refType {
		t.Errorf("expected refType %s, got %s", refType, gotRefType)
	}
	if gotLocation != location {
		t.Errorf("expected location %s, got %s", location, gotLocation)
	}
}

func TestFailsPackageExternalRefWithInvalidFormat(t *testing.T) {
	_, _, _, err := extractPackageExternalReference("whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1PackageWithoutSpdxIdentifierThrowsError(t *testing.T) {
	// More than one package, the previous package doesn't contain the SPDXID
	pkgOldName := "p1"
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: map[spdx.ElementID]*spdx.Package2_1{}},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: pkgOldName},
	}
	pkgOld := parser.pkg
	parser.doc.Packages["p1"] = pkgOld
	// the Document's Packages should have this one only
	if parser.doc.Packages["p1"] != pkgOld {
		t.Errorf("expected package %v, got %v", pkgOld, parser.doc.Packages["p1"])
	}
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}

	pkgName := "p2"
	err := parser.parsePair2_1("PackageName", pkgName)
	if err == nil {
		t.Errorf("package without SPDX Identifier getting accepted")
	}
}
