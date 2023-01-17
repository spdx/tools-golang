// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import (
	"fmt"
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// ===== Parser package section state change tests =====
func TestParserPackageStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	// create the first package
	pkgOldName := "p1"

	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: pkgOldName, PackageSPDXIdentifier: "p1"},
	}
	pkgOld := parser.pkg
	parser.doc.Packages = append(parser.doc.Packages, pkgOld)
	// the Document's Packages should have this one only
	if parser.doc.Packages[0] != pkgOld {
		t.Errorf("expected package %v, got %v", pkgOld, parser.doc.Packages[0])
	}
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}

	// now add a new package
	pkgName := "p2"
	err := parser.parsePair("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psPackage {
		t.Errorf("expected state to be %v, got %v", psPackage, parser.st)
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
	if parser.doc.Packages[0] != pkgOld {
		t.Errorf("Expected package %v, got %v", pkgOld, parser.doc.Packages[0])
	}
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}
}

func TestParserPackageStartsNewPackageAfterParsingPackageNameTagWhileInUnpackaged(t *testing.T) {
	// pkg is nil, so that Files appearing before the first PackageName tag
	// are added to Files instead of Packages
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psFile,
		pkg: nil,
	}
	// the Document's Packages should be empty
	if len(parser.doc.Packages) != 0 {
		t.Errorf("Expected zero packages, got %d", len(parser.doc.Packages))
	}

	// now add a new package
	pkgName := "p2"
	err := parser.parsePair("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psPackage {
		t.Errorf("expected state to be %v, got %v", psPackage, parser.st)
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

func TestParserPackageMovesToFileAfterParsingFileNameTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	pkgCurrent := parser.pkg

	err := parser.parsePair("FileName", "testFile")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psFile {
		t.Errorf("expected state to be %v, got %v", psFile, parser.st)
	}
	// and current package should remain what it was
	if parser.pkg != pkgCurrent {
		t.Fatalf("expected package to remain %v, got %v", pkgCurrent, parser.pkg)
	}
}

func TestParserPackageMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("expected state to be %v, got %v", psOtherLicense, parser.st)
	}
}

func TestParserPackageMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("expected state to be %v, got %v", psReview, parser.st)
	}
}

func TestParserPackageStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should remain unchanged
	if parser.st != psPackage {
		t.Errorf("expected state to be %v, got %v", psPackage, parser.st)
	}

	err = parser.parsePair("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psPackage {
		t.Errorf("expected state to be %v, got %v", psPackage, parser.st)
	}
}

func TestParserPackageStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psPackage {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage)
	}

	err = parser.parsePair("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psPackage {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage)
	}

	err = parser.parsePair("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psPackage {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage)
	}

	err = parser.parsePair("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psPackage {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage)
	}

	err = parser.parsePair("AnnotationComment", "i guess i had something to say about this package")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psPackage {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage)
	}
}

// ===== Package data section tests =====
func TestParserCanParsePackageTags(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// should not yet be in Packages map, b/c no SPDX identifier
	if len(parser.doc.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(parser.doc.Packages))
	}

	// Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
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
	err = parser.parsePairFromPackage("SPDXID", "SPDXRef-p1")
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
	if parser.doc.Packages[0] != parser.pkg {
		t.Errorf("expected to point to parser.pkg, got %v", parser.doc.Packages[0])
	}

	// Package Version
	err = parser.parsePairFromPackage("PackageVersion", "2.1.1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageVersion != "2.1.1" {
		t.Errorf("got %v for PackageVersion", parser.pkg.PackageVersion)
	}

	// Package File Name
	err = parser.parsePairFromPackage("PackageFileName", "p1-2.1.1.tar.gz")
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
	err = parser.parsePairFromPackage("PackageDownloadLocation", "https://example.com/whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageDownloadLocation != "https://example.com/whatever" {
		t.Errorf("got %v for PackageDownloadLocation", parser.pkg.PackageDownloadLocation)
	}

	// Files Analyzed
	err = parser.parsePairFromPackage("FilesAnalyzed", "false")
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

	testChecksums := map[common.ChecksumAlgorithm]string{
		"MD5":    "624c1abb3664f4b35547e7c73864ad24",
		"SHA1":   "85ed0817af83a24ad8da68c2b5094de69833983c",
		"SHA256": "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd",
		"SHA512": "4ced3267f5ed38df65ceebc43e97aa6c2948cc7ef3288c2e5074e7df7fab544cc93339604513ea5f65616f9ed1c48581465043c8a9b693ef20fd4fddaf25e1b9",
	}

	for algo, tc := range testChecksums {
		if err := parser.parsePairFromPackage(
			"PackageChecksum", fmt.Sprintf("%s: %s", algo, tc)); err != nil {
			t.Errorf("expected error, got %v", err)
		}
	}

	for _, checksum := range parser.pkg.PackageChecksums {
		if checksum.Value != testChecksums[checksum.Algorithm] {
			t.Errorf(
				"expected %s for PackageChecksum%s, got %s",
				testChecksums[checksum.Algorithm], checksum.Algorithm, checksum.Value,
			)
		}
	}

	// Package Home Page
	err = parser.parsePairFromPackage("PackageHomePage", "https://example.com/whatever2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageHomePage != "https://example.com/whatever2" {
		t.Errorf("got %v for PackageHomePage", parser.pkg.PackageHomePage)
	}

	// Package Source Info
	err = parser.parsePairFromPackage("PackageSourceInfo", "random comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSourceInfo != "random comment" {
		t.Errorf("got %v for PackageSourceInfo", parser.pkg.PackageSourceInfo)
	}

	// Package License Concluded
	err = parser.parsePairFromPackage("PackageLicenseConcluded", "Apache-2.0 OR GPL-2.0-or-later")
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
		err = parser.parsePairFromPackage("PackageLicenseInfoFromFiles", lic)
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
	err = parser.parsePairFromPackage("PackageLicenseDeclared", "Apache-2.0 OR GPL-2.0-or-later")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageLicenseDeclared != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("got %v for PackageLicenseDeclared", parser.pkg.PackageLicenseDeclared)
	}

	// Package License Comments
	err = parser.parsePairFromPackage("PackageLicenseComments", "this is a license comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageLicenseComments != "this is a license comment" {
		t.Errorf("got %v for PackageLicenseComments", parser.pkg.PackageLicenseComments)
	}

	// Package Copyright Text
	err = parser.parsePairFromPackage("PackageCopyrightText", "Copyright (c) me myself and i")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageCopyrightText != "Copyright (c) me myself and i" {
		t.Errorf("got %v for PackageCopyrightText", parser.pkg.PackageCopyrightText)
	}

	// Package Summary
	err = parser.parsePairFromPackage("PackageSummary", "i wrote this package")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSummary != "i wrote this package" {
		t.Errorf("got %v for PackageSummary", parser.pkg.PackageSummary)
	}

	// Package Description
	err = parser.parsePairFromPackage("PackageDescription", "i wrote this package a lot")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageDescription != "i wrote this package a lot" {
		t.Errorf("got %v for PackageDescription", parser.pkg.PackageDescription)
	}

	// Package Comment
	err = parser.parsePairFromPackage("PackageComment", "i scanned this package")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageComment != "i scanned this package" {
		t.Errorf("got %v for PackageComment", parser.pkg.PackageComment)
	}

	// Package Attribution Text
	attrs := []string{
		"Include this notice in all advertising materials",
		"This is a \nmulti-line string",
	}
	for _, attr := range attrs {
		err = parser.parsePairFromPackage("PackageAttributionText", attr)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, attrWant := range attrs {
		flagFound := false
		for _, attrCheck := range parser.pkg.PackageAttributionTexts {
			if attrWant == attrCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in PackageAttributionText", attrWant)
		}
	}
	if len(attrs) != len(parser.pkg.PackageAttributionTexts) {
		t.Errorf("expected %d attribution texts in PackageAttributionTexts, got %d", len(attrs),
			len(parser.pkg.PackageAttributionTexts))
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
	err = parser.parsePairFromPackage("ExternalRef", ref1)
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
	err = parser.parsePairFromPackage("ExternalRefComment", ref1Comment)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromPackage("ExternalRef", ref2)
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
	err = parser.parsePairFromPackage("ExternalRefComment", ref2Comment)
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

func TestParserCanParsePackageSupplierPersonTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Supplier: Person
	err := parser.parsePairFromPackage("PackageSupplier", "Person: John Doe")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSupplier.Supplier != "John Doe" {
		t.Errorf("got %v for PackageSupplierPerson", parser.pkg.PackageSupplier.Supplier)
	}
}

func TestParserCanParsePackageSupplierOrganizationTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Supplier: Organization
	err := parser.parsePairFromPackage("PackageSupplier", "Organization: John Doe, Inc.")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSupplier.Supplier != "John Doe, Inc." {
		t.Errorf("got %v for PackageSupplierOrganization", parser.pkg.PackageSupplier.Supplier)
	}
}

func TestParserCanParsePackageSupplierNOASSERTIONTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Supplier: NOASSERTION
	err := parser.parsePairFromPackage("PackageSupplier", "NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSupplier.Supplier != "NOASSERTION" {
		t.Errorf("got value for Supplier, expected NOASSERTION")
	}
}

func TestParserCanParsePackageOriginatorPersonTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Originator: Person
	err := parser.parsePairFromPackage("PackageOriginator", "Person: John Doe")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageOriginator.Originator != "John Doe" {
		t.Errorf("got %v for PackageOriginator", parser.pkg.PackageOriginator.Originator)
	}
}

func TestParserCanParsePackageOriginatorOrganizationTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Originator: Organization
	err := parser.parsePairFromPackage("PackageOriginator", "Organization: John Doe, Inc.")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageOriginator.Originator != "John Doe, Inc." {
		t.Errorf("got %v for PackageOriginator", parser.pkg.PackageOriginator.Originator)
	}
}

func TestParserCanParsePackageOriginatorNOASSERTIONTag(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Originator: NOASSERTION
	err := parser.parsePairFromPackage("PackageOriginator", "NOASSERTION")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageOriginator.Originator != "NOASSERTION" {
		t.Errorf("got false for PackageOriginatorNOASSERTION")
	}
}

func TestParserCanParsePackageVerificationCodeTagWithExcludes(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Verification Code with excludes parenthetical
	code := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	fileName := "./package.spdx"
	fullCodeValue := "d6a770ba38583ed4bb4525bd96e50461655d2758 (excludes: ./package.spdx)"
	err := parser.parsePairFromPackage("PackageVerificationCode", fullCodeValue)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageVerificationCode.Value != code {
		t.Errorf("got %v for PackageVerificationCode", parser.pkg.PackageVerificationCode)
	}
	if len(parser.pkg.PackageVerificationCode.ExcludedFiles) != 1 || parser.pkg.PackageVerificationCode.ExcludedFiles[0] != fileName {
		t.Errorf("got %v for PackageVerificationCodeExcludedFile", parser.pkg.PackageVerificationCode.ExcludedFiles)
	}

}

func TestParserCanParsePackageVerificationCodeTagWithoutExcludes(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Verification Code without excludes parenthetical
	code := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	err := parser.parsePairFromPackage("PackageVerificationCode", code)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageVerificationCode.Value != code {
		t.Errorf("got %v for PackageVerificationCode", parser.pkg.PackageVerificationCode)
	}
	if len(parser.pkg.PackageVerificationCode.ExcludedFiles) != 0 {
		t.Errorf("got %v for PackageVerificationCodeExcludedFile", parser.pkg.PackageVerificationCode.ExcludedFiles)
	}

}

func TestParserPackageExternalRefPointerChangesAfterTags(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	ref1 := "SECURITY cpe23Type cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*"
	err := parser.parsePairFromPackage("ExternalRef", ref1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef == nil {
		t.Errorf("expected non-nil external reference pointer, got nil")
	}

	// now, a comment; pointer should go away
	err = parser.parsePairFromPackage("ExternalRefComment", "whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef != nil {
		t.Errorf("expected nil external reference pointer, got non-nil")
	}

	ref2 := "Other LocationRef-something https://example.com/whatever"
	err = parser.parsePairFromPackage("ExternalRef", ref2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef == nil {
		t.Errorf("expected non-nil external reference pointer, got nil")
	}

	// and some other random tag makes the pointer go away too
	err = parser.parsePairFromPackage("PackageSummary", "whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkgExtRef != nil {
		t.Errorf("expected nil external reference pointer, got non-nil")
	}
}

func TestParserPackageCreatesRelationshipInDocument(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-whatever")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.rln == nil {
		t.Fatalf("parser didn't create and point to Relationship struct")
	}
	if parser.rln != parser.doc.Relationships[0] {
		t.Errorf("pointer to new Relationship doesn't match idx 0 for doc.Relationships[]")
	}
}

func TestParserPackageCreatesAnnotationInDocument(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.ann == nil {
		t.Fatalf("parser didn't create and point to Annotation struct")
	}
	if parser.ann != parser.doc.Annotations[0] {
		t.Errorf("pointer to new Annotation doesn't match idx 0 for doc.Annotations[]")
	}
}

func TestParserPackageUnknownTagFails(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: "p1", PackageSPDXIdentifier: "p1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePairFromPackage("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestParserFailsIfInvalidSPDXIDInPackageSection(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid ID format
	err = parser.parsePairFromPackage("SPDXID", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfInvalidPackageSupplierFormat(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid supplier format
	err = parser.parsePairFromPackage("PackageSupplier", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfUnknownPackageSupplierType(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid supplier type
	err = parser.parsePairFromPackage("PackageSupplier", "whoops: John Doe")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfInvalidPackageOriginatorFormat(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid originator format
	err = parser.parsePairFromPackage("PackageOriginator", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfUnknownPackageOriginatorType(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid originator type
	err = parser.parsePairFromPackage("PackageOriginator", "whoops: John Doe")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserSetsFilesAnalyzedTagsCorrectly(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// set tag
	err = parser.parsePairFromPackage("FilesAnalyzed", "true")
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

func TestParserFailsIfInvalidPackageChecksumFormat(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid checksum format
	err = parser.parsePairFromPackage("PackageChecksum", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfInvalidPackageChecksumType(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid checksum type
	err = parser.parsePairFromPackage("PackageChecksum", "whoops: blah")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfInvalidExternalRefFormat(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid external ref format
	err = parser.parsePairFromPackage("ExternalRef", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsIfExternalRefCommentBeforeExternalRef(t *testing.T) {
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{},
	}

	// start with Package Name
	err := parser.parsePairFromPackage("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// external ref comment before external ref
	err = parser.parsePairFromPackage("ExternalRefComment", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// ===== Helper function tests =====

func TestCanCheckAndExtractExcludesFilenameAndCode(t *testing.T) {
	code := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	fileName := "./package.spdx"
	fullCodeValue := "d6a770ba38583ed4bb4525bd96e50461655d2758 (excludes: ./package.spdx)"

	gotCode := extractCodeAndExcludes(fullCodeValue)
	if gotCode.Value != code {
		t.Errorf("got %v for gotCode", gotCode)
	}
	if len(gotCode.ExcludedFiles) != 1 || gotCode.ExcludedFiles[0] != fileName {
		t.Errorf("got %v for gotFileName", gotCode.ExcludedFiles)
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

func TestParserPackageWithoutSpdxIdentifierThrowsError(t *testing.T) {
	// More than one package, the previous package doesn't contain an SPDX ID
	pkgOldName := "p1"
	parser := tvParser{
		doc: &v2_2.Document{Packages: []*v2_2.Package{}},
		st:  psPackage,
		pkg: &v2_2.Package{PackageName: pkgOldName},
	}
	pkgOld := parser.pkg
	parser.doc.Packages = append(parser.doc.Packages, pkgOld)
	// the Document's Packages should have this one only
	if parser.doc.Packages[0] != pkgOld {
		t.Errorf("expected package %v, got %v", pkgOld, parser.doc.Packages[0])
	}
	if len(parser.doc.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(parser.doc.Packages))
	}

	pkgName := "p2"
	err := parser.parsePair("PackageName", pkgName)
	if err == nil {
		t.Errorf("package without SPDX Identifier getting accepted")
	}
}
