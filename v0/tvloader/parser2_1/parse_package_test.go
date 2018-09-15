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

func TestParser2_1PackageMovesToFileAfterParsingFileNameTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
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
	if parser.pkg.IsUnpackaged {
		t.Errorf("expected IsUnpackaged to be false, got true")
	}
}

func TestParser2_1PackageMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psOtherLicense2_1 {
		t.Errorf("expected state to be %v, got %v", psOtherLicense2_1, parser.st)
	}
}

func TestParser2_1PackageStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should remain unchanges
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should still remain unchanges
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}
}

// ===== Package data section tests =====
func TestParser2_1CanParsePackageTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// Package Name
	err := parser.parsePairFromPackage2_1("PackageName", "p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageName != "p1" {
		t.Errorf("got %v for PackageName", parser.pkg.PackageName)
	}

	// Package SPDX Identifier
	err = parser.parsePairFromPackage2_1("SPDXID", "SPDXRef-p1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.pkg.PackageSPDXIdentifier != "SPDXRef-p1" {
		t.Errorf("got %v for PackageSPDXIdentifier", parser.pkg.PackageSPDXIdentifier)
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
}

func TestParser2_1CanParsePackageSupplierPersonTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
		pkg: &spdx.Package2_1{PackageName: "p1", IsUnpackaged: false},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

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
