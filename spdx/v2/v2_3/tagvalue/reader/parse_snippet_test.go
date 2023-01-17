// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import (
	"testing"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// ===== Parser snippet section state change tests =====
func TestParserSnippetStartsNewSnippetAfterParsingSnippetSPDXIDTag(t *testing.T) {
	// create the first snippet
	sid1 := common.ElementID("s1")
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: sid1},
	}
	s1 := parser.snippet
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets[sid1] = parser.snippet

	// the File's Snippets should have this one only
	if len(parser.file.Snippets) != 1 {
		t.Errorf("Expected len(Snippets) to be 1, got %d", len(parser.file.Snippets))
	}
	if parser.file.Snippets["s1"] != s1 {
		t.Errorf("Expected snippet %v in Snippets[s1], got %v", s1, parser.file.Snippets["s1"])
	}
	if parser.file.Snippets["s1"].SnippetSPDXIdentifier != sid1 {
		t.Errorf("expected snippet ID %s in Snippets[s1], got %s", sid1, parser.file.Snippets["s1"].SnippetSPDXIdentifier)
	}

	// now add a new snippet
	err := parser.parsePair("SnippetSPDXID", "SPDXRef-s2")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psSnippet {
		t.Errorf("expected state to be %v, got %v", psSnippet, parser.st)
	}
	// and a snippet should be created
	if parser.snippet == nil {
		t.Fatalf("parser didn't create new snippet")
	}
	// and the snippet ID should be as expected
	if parser.snippet.SnippetSPDXIdentifier != "s2" {
		t.Errorf("expected snippet ID %s, got %s", "s2", parser.snippet.SnippetSPDXIdentifier)
	}
	// and the File's Snippets should be of size 2 and have these two
	if len(parser.file.Snippets) != 2 {
		t.Errorf("Expected len(Snippets) to be 2, got %d", len(parser.file.Snippets))
	}
	if parser.file.Snippets["s1"] != s1 {
		t.Errorf("Expected snippet %v in Snippets[s1], got %v", s1, parser.file.Snippets["s1"])
	}
	if parser.file.Snippets["s1"].SnippetSPDXIdentifier != sid1 {
		t.Errorf("expected snippet ID %s in Snippets[s1], got %s", sid1, parser.file.Snippets["s1"].SnippetSPDXIdentifier)
	}
	if parser.file.Snippets["s2"] != parser.snippet {
		t.Errorf("Expected snippet %v in Snippets[s2], got %v", parser.snippet, parser.file.Snippets["s2"])
	}
	if parser.file.Snippets["s2"].SnippetSPDXIdentifier != "s2" {
		t.Errorf("expected snippet ID %s in Snippets[s2], got %s", "s2", parser.file.Snippets["s2"].SnippetSPDXIdentifier)
	}
}

func TestParserSnippetStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	p1 := parser.pkg
	f1 := parser.file
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets["s1"] = parser.snippet

	// now add a new package
	p2Name := "package2"
	err := parser.parsePair("PackageName", p2Name)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should go back to Package
	if parser.st != psPackage {
		t.Errorf("expected state to be %v, got %v", psPackage, parser.st)
	}
	// and a package should be created
	if parser.pkg == nil {
		t.Fatalf("parser didn't create new pkg")
	}
	// and the package name should be as expected
	if parser.pkg.PackageName != p2Name {
		t.Errorf("expected package name %s, got %s", p2Name, parser.pkg.PackageName)
	}
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the Document's Packages should still be of size 1 b/c no SPDX
	// identifier has been seen yet
	if len(parser.doc.Packages) != 1 {
		t.Errorf("Expected len(Packages) to be 1, got %d", len(parser.doc.Packages))
	}
	if parser.doc.Packages[0] != p1 {
		t.Errorf("Expected package %v in Packages[package1], got %v", p1, parser.doc.Packages[0])
	}
	if parser.doc.Packages[0].PackageName != "package1" {
		t.Errorf("expected package name %s in Packages[package1], got %s", "package1", parser.doc.Packages[0].PackageName)
	}
	// and the first Package's Files should be of size 1 and have f1 only
	if len(parser.doc.Packages[0].Files) != 1 {
		t.Errorf("Expected 1 file in Packages[package1].Files, got %d", len(parser.doc.Packages[0].Files))
	}
	if parser.doc.Packages[0].Files[0] != f1 {
		t.Errorf("Expected file %v in Files[f1], got %v", f1, parser.doc.Packages[0].Files[0])
	}
	if parser.doc.Packages[0].Files[0].FileName != "f1.txt" {
		t.Errorf("expected file name %s in Files[f1], got %s", "f1.txt", parser.doc.Packages[0].Files[0].FileName)
	}
	// and the new Package should have no files
	if len(parser.pkg.Files) != 0 {
		t.Errorf("Expected no files in Packages[1].Files, got %d", len(parser.pkg.Files))
	}
	// and the current file should be nil
	if parser.file != nil {
		t.Errorf("Expected nil for parser.file, got %v", parser.file)
	}
	// and the current snippet should be nil
	if parser.snippet != nil {
		t.Errorf("Expected nil for parser.snippet, got %v", parser.snippet)
	}
}

func TestParserSnippetMovesToFileAfterParsingFileNameTag(t *testing.T) {
	f1Name := "f1.txt"
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	p1 := parser.pkg
	f1 := parser.file
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets["s1"] = parser.snippet

	f2Name := "f2.txt"
	err := parser.parsePair("FileName", f2Name)
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should be correct
	if parser.st != psFile {
		t.Errorf("expected state to be %v, got %v", psSnippet, parser.st)
	}
	// and current package should remain what it was
	if parser.pkg != p1 {
		t.Fatalf("expected package to remain %v, got %v", p1, parser.pkg)
	}
	// and a file should be created
	if parser.file == nil {
		t.Fatalf("parser didn't create new file")
	}
	// and the file name should be as expected
	if parser.file.FileName != f2Name {
		t.Errorf("expected file name %s, got %s", f2Name, parser.file.FileName)
	}
	// and the Package's Files should still be of size 1 since we haven't seen
	// an SPDX identifier yet for this new file
	if len(parser.pkg.Files) != 1 {
		t.Errorf("Expected len(Files) to be 1, got %d", len(parser.pkg.Files))
	}
	if parser.pkg.Files[0] != f1 {
		t.Errorf("Expected file %v in Files[f1], got %v", f1, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != f1Name {
		t.Errorf("expected file name %s in Files[f1], got %s", f1Name, parser.pkg.Files[0].FileName)
	}
	// and the current snippet should be nil
	if parser.snippet != nil {
		t.Errorf("Expected nil for parser.snippet, got %v", parser.snippet)
	}
}

func TestParserSnippetMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets["s1"] = parser.snippet

	err := parser.parsePair("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psOtherLicense {
		t.Errorf("expected state to be %v, got %v", psOtherLicense, parser.st)
	}
}

func TestParserSnippetMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets["s1"] = parser.snippet

	err := parser.parsePair("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psReview {
		t.Errorf("expected state to be %v, got %v", psReview, parser.st)
	}
}

func TestParserSnippetStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets["s1"] = parser.snippet

	err := parser.parsePair("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should remain unchanged
	if parser.st != psSnippet {
		t.Errorf("expected state to be %v, got %v", psSnippet, parser.st)
	}
	// and the relationship should be in the Document's Relationships
	if len(parser.doc.Relationships) != 1 {
		t.Fatalf("expected doc.Relationships to have len 1, got %d", len(parser.doc.Relationships))
	}
	deID := parser.doc.Relationships[0].RefA
	if deID.DocumentRefID != "" || deID.ElementRefID != "blah" {
		t.Errorf("expected RefA to be %s, got %s", "blah", parser.doc.Relationships[0].RefA)
	}

	err = parser.parsePair("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psSnippet {
		t.Errorf("expected state to be %v, got %v", psSnippet, parser.st)
	}
}

func TestParserSnippetStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets["s1"] = parser.snippet

	err := parser.parsePair("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psSnippet {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet)
	}

	err = parser.parsePair("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psSnippet {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet)
	}

	err = parser.parsePair("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psSnippet {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet)
	}

	err = parser.parsePair("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psSnippet {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet)
	}

	err = parser.parsePair("AnnotationComment", "i guess i had something to say about this particular file")
	if err != nil {
		t.Errorf("got error when calling parsePair: %v", err)
	}
	if parser.st != psSnippet {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet)
	}

	// and the annotation should be in the Document's Annotations
	if len(parser.doc.Annotations) != 1 {
		t.Fatalf("expected doc.Annotations to have len 1, got %d", len(parser.doc.Annotations))
	}
	if parser.doc.Annotations[0].Annotator.Annotator != "John Doe ()" {
		t.Errorf("expected Annotator to be %s, got %s", "John Doe ()", parser.doc.Annotations[0].Annotator)
	}
}

// ===== Snippet data section tests =====
func TestParserCanParseSnippetTags(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	// Snippet SPDX Identifier
	err := parser.parsePairFromSnippet("SnippetSPDXID", "SPDXRef-s1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetSPDXIdentifier != "s1" {
		t.Errorf("got %v for SnippetSPDXIdentifier", parser.snippet.SnippetSPDXIdentifier)
	}

	// Snippet from File SPDX Identifier
	err = parser.parsePairFromSnippet("SnippetFromFileSPDXID", "SPDXRef-f1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wantDeID := common.DocElementID{DocumentRefID: "", ElementRefID: common.ElementID("f1")}
	if parser.snippet.SnippetFromFileSPDXIdentifier != wantDeID.ElementRefID {
		t.Errorf("got %v for SnippetFromFileSPDXIdentifier", parser.snippet.SnippetFromFileSPDXIdentifier)
	}

	// Snippet Byte Range
	err = parser.parsePairFromSnippet("SnippetByteRange", "20:320")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.Ranges[0].StartPointer.Offset != 20 {
		t.Errorf("got %v for SnippetByteRangeStart", parser.snippet.Ranges[0].StartPointer.Offset)
	}
	if parser.snippet.Ranges[0].EndPointer.Offset != 320 {
		t.Errorf("got %v for SnippetByteRangeEnd", parser.snippet.Ranges[0].EndPointer.Offset)
	}

	// Snippet Line Range
	err = parser.parsePairFromSnippet("SnippetLineRange", "5:12")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.Ranges[1].StartPointer.LineNumber != 5 {
		t.Errorf("got %v for SnippetLineRangeStart", parser.snippet.Ranges[1].StartPointer.LineNumber)
	}
	if parser.snippet.Ranges[1].EndPointer.LineNumber != 12 {
		t.Errorf("got %v for SnippetLineRangeEnd", parser.snippet.Ranges[1].EndPointer.LineNumber)
	}

	// Snippet Concluded License
	err = parser.parsePairFromSnippet("SnippetLicenseConcluded", "BSD-3-Clause")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetLicenseConcluded != "BSD-3-Clause" {
		t.Errorf("got %v for SnippetLicenseConcluded", parser.snippet.SnippetLicenseConcluded)
	}

	// License Information in Snippet
	lics := []string{
		"Apache-2.0",
		"GPL-2.0-or-later",
		"CC0-1.0",
	}
	for _, lic := range lics {
		err = parser.parsePairFromSnippet("LicenseInfoInSnippet", lic)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, licWant := range lics {
		flagFound := false
		for _, licCheck := range parser.snippet.LicenseInfoInSnippet {
			if licWant == licCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in LicenseInfoInSnippet", licWant)
		}
	}
	if len(lics) != len(parser.snippet.LicenseInfoInSnippet) {
		t.Errorf("expected %d licenses in LicenseInfoInSnippet, got %d", len(lics),
			len(parser.snippet.LicenseInfoInSnippet))
	}

	// Snippet Comments on License
	err = parser.parsePairFromSnippet("SnippetLicenseComments", "this is a comment about the licenses")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetLicenseComments != "this is a comment about the licenses" {
		t.Errorf("got %v for SnippetLicenseComments", parser.snippet.SnippetLicenseComments)
	}

	// Snippet Copyright Text
	err = parser.parsePairFromSnippet("SnippetCopyrightText", "copyright (c) John Doe and friends")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetCopyrightText != "copyright (c) John Doe and friends" {
		t.Errorf("got %v for SnippetCopyrightText", parser.snippet.SnippetCopyrightText)
	}

	// Snippet Comment
	err = parser.parsePairFromSnippet("SnippetComment", "this is a comment about the snippet")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetComment != "this is a comment about the snippet" {
		t.Errorf("got %v for SnippetComment", parser.snippet.SnippetComment)
	}

	// Snippet Name
	err = parser.parsePairFromSnippet("SnippetName", "from some other package called abc")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetName != "from some other package called abc" {
		t.Errorf("got %v for SnippetName", parser.snippet.SnippetName)
	}

	// Snippet Attribution Texts
	attrs := []string{
		"Include this notice in all advertising materials",
		"This is a \nmulti-line string",
	}
	for _, attr := range attrs {
		err = parser.parsePairFromSnippet("SnippetAttributionText", attr)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, attrWant := range attrs {
		flagFound := false
		for _, attrCheck := range parser.snippet.SnippetAttributionTexts {
			if attrWant == attrCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in SnippetAttributionText", attrWant)
		}
	}
	if len(attrs) != len(parser.snippet.SnippetAttributionTexts) {
		t.Errorf("expected %d attribution texts in SnippetAttributionTexts, got %d", len(attrs),
			len(parser.snippet.SnippetAttributionTexts))
	}

}

func TestParserSnippetUnknownTagFails(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{SnippetSPDXIdentifier: "s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePairFromSnippet("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestParserFailsForInvalidSnippetSPDXID(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	// invalid Snippet SPDX Identifier
	err := parser.parsePairFromSnippet("SnippetSPDXID", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsForInvalidSnippetFromFileSPDXID(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	// start with Snippet SPDX Identifier
	err := parser.parsePairFromSnippet("SnippetSPDXID", "SPDXRef-s1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid From File identifier
	err = parser.parsePairFromSnippet("SnippetFromFileSPDXID", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsForInvalidSnippetByteValues(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	// start with Snippet SPDX Identifier
	err := parser.parsePairFromSnippet("SnippetSPDXID", "SPDXRef-s1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid byte formats and values
	err = parser.parsePairFromSnippet("SnippetByteRange", "200 210")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
	err = parser.parsePairFromSnippet("SnippetByteRange", "a:210")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
	err = parser.parsePairFromSnippet("SnippetByteRange", "200:a")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFailsForInvalidSnippetLineValues(t *testing.T) {
	parser := tvParser{
		doc:     &spdx.Document{Packages: []*spdx.Package{}},
		st:      psSnippet,
		pkg:     &spdx.Package{PackageName: "package1", PackageSPDXIdentifier: "package1", Files: []*spdx.File{}},
		file:    &spdx.File{FileName: "f1.txt", FileSPDXIdentifier: "f1", Snippets: map[common.ElementID]*spdx.Snippet{}},
		snippet: &spdx.Snippet{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	// start with Snippet SPDX Identifier
	err := parser.parsePairFromSnippet("SnippetSPDXID", "SPDXRef-s1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid byte formats and values
	err = parser.parsePairFromSnippet("SnippetLineRange", "200 210")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
	err = parser.parsePairFromSnippet("SnippetLineRange", "a:210")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
	err = parser.parsePairFromSnippet("SnippetLineRange", "200:a")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParserFilesWithoutSpdxIdThrowErrorWithSnippets(t *testing.T) {
	// Invalid file with snippet
	// Last unpackaged file before the snippet starts
	// Last file of a package and New package starts
	fileName := "f2.txt"
	sid1 := common.ElementID("s1")
	parser2 := tvParser{
		doc:  &spdx.Document{},
		st:   psCreationInfo,
		file: &spdx.File{FileName: fileName},
	}
	err := parser2.parsePair("SnippetSPDXID", string(sid1))
	if err == nil {
		t.Errorf("file without SPDX Identifier getting accepted")
	}

}
