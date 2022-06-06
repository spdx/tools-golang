// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Parser file section state change tests =====
func TestParser2_1FileStartsNewFileAfterParsingFileNameTag(t *testing.T) {
	// create the first file
	fileOldName := "f1.txt"

	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: fileOldName, FileSPDXIdentifier: "f1"},
	}
	fileOld := parser.file
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, fileOld)
	// the Package's Files should have this one only
	if len(parser.pkg.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(parser.pkg.Files))
	}
	if parser.pkg.Files[0] != fileOld {
		t.Errorf("expected file %v in Files[f1], got %v", fileOld, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != fileOldName {
		t.Errorf("expected file name %s in Files[f1], got %s", fileOldName, parser.pkg.Files[0].FileName)
	}

	// now add a new file
	fileName := "f2.txt"
	err := parser.parsePair2_1("FileName", fileName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psFile2_1 {
		t.Errorf("expected state to be %v, got %v", psFile2_1, parser.st)
	}
	// and a file should be created
	if parser.file == nil {
		t.Fatalf("parser didn't create new file")
	}
	// and the file name should be as expected
	if parser.file.FileName != fileName {
		t.Errorf("expected file name %s, got %s", fileName, parser.file.FileName)
	}
	// and the Package's Files should still be of size 1 and not have this new
	// one yet, since it hasn't seen an SPDX identifier
	if len(parser.pkg.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(parser.pkg.Files))
	}
	if parser.pkg.Files[0] != fileOld {
		t.Errorf("expected file %v in Files[f1], got %v", fileOld, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != fileOldName {
		t.Errorf("expected file name %s in Files[f1], got %s", fileOldName, parser.pkg.Files[0].FileName)
	}

	// now parse an SPDX identifier tag
	err = parser.parsePair2_1("SPDXID", "SPDXRef-f2ID")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// and the Package's Files should now be of size 2 and have this new one
	if len(parser.pkg.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(parser.pkg.Files))
	}
	if parser.pkg.Files[0] != fileOld {
		t.Errorf("expected file %v in Files[f1], got %v", fileOld, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != fileOldName {
		t.Errorf("expected file name %s in Files[f1], got %s", fileOldName, parser.pkg.Files[0].FileName)
	}
	if parser.pkg.Files[1] != parser.file {
		t.Errorf("expected file %v in Files[f2ID], got %v", parser.file, parser.pkg.Files[1])
	}
	if parser.pkg.Files[1].FileName != fileName {
		t.Errorf("expected file name %s in Files[f2ID], got %s", fileName, parser.pkg.Files[1].FileName)
	}
}

func TestParser2_1FileAddsToPackageOrUnpackagedFiles(t *testing.T) {
	// start with no packages
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// add a file and SPDX identifier
	fileName := "f2.txt"
	err := parser.parsePair2_1("FileName", fileName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	err = parser.parsePair2_1("SPDXID", "SPDXRef-f2ID")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	fileOld := parser.file
	// should have been added to Files
	if len(parser.doc.Files) != 1 {
		t.Fatalf("expected 1 file in Files, got %d", len(parser.doc.Files))
	}
	if parser.doc.Files[0] != fileOld {
		t.Errorf("expected file %v in Files[f2ID], got %v", fileOld, parser.doc.Files[0])
	}
	// now create a package and a new file
	err = parser.parsePair2_1("PackageName", "package1")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	err = parser.parsePair2_1("SPDXID", "SPDXRef-pkg1")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	err = parser.parsePair2_1("FileName", "f3.txt")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	err = parser.parsePair2_1("SPDXID", "SPDXRef-f3ID")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// Files should still be size 1 and have old file only
	if len(parser.doc.Files) != 1 {
		t.Fatalf("expected 1 file in Files, got %d", len(parser.doc.Files))
	}
	if parser.doc.Files[0] != fileOld {
		t.Errorf("expected file %v in Files[f2ID], got %v", fileOld, parser.doc.Files[0])
	}
	// and new package should have gotten the new file
	if len(parser.pkg.Files) != 1 {
		t.Fatalf("expected 1 file in Files, got %d", len(parser.pkg.Files))
	}
	if parser.pkg.Files[0] != parser.file {
		t.Errorf("expected file %v in Files[f3ID], got %v", parser.file, parser.pkg.Files[0])
	}
}

func TestParser2_1FileStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	// create the first file and package
	p1Name := "package1"
	f1Name := "f1.txt"

	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: p1Name, PackageSPDXIdentifier: "package1", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: f1Name, FileSPDXIdentifier: "f1"},
	}
	p1 := parser.pkg
	f1 := parser.file
	parser.doc.Packages = append(parser.doc.Packages, p1)
	parser.pkg.Files = append(parser.pkg.Files, f1)

	// now add a new package
	p2Name := "package2"
	err := parser.parsePair2_1("PackageName", p2Name)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should go back to Package
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
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
	if *parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	// and the new Package should have no files
	if len(parser.pkg.Files) != 0 {
		t.Errorf("Expected no files in pkg.Files, got %d", len(parser.pkg.Files))
	}
	// and the Document's Packages should still be of size 1 and not have this
	// new one, because no SPDX identifier has been seen yet
	if len(parser.doc.Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(parser.doc.Packages))
	}
	if parser.doc.Packages[0] != p1 {
		t.Errorf("Expected package %v in Packages[package1], got %v", p1, parser.doc.Packages[0])
	}
	if parser.doc.Packages[0].PackageName != p1Name {
		t.Errorf("expected package name %s in Packages[package1], got %s", p1Name, parser.doc.Packages[0].PackageName)
	}
	// and the first Package's Files should be of size 1 and have f1 only
	if len(parser.doc.Packages[0].Files) != 1 {
		t.Errorf("Expected 1 file in Packages[package1].Files, got %d", len(parser.doc.Packages[0].Files))
	}
	if parser.doc.Packages[0].Files[0] != f1 {
		t.Errorf("Expected file %v in Files[f1], got %v", f1, parser.doc.Packages[0].Files[0])
	}
	if parser.doc.Packages[0].Files[0].FileName != f1Name {
		t.Errorf("expected file name %s in Files[f1], got %s", f1Name, parser.doc.Packages[0].Files[0].FileName)
	}
	// and the current file should be nil
	if parser.file != nil {
		t.Errorf("Expected nil for parser.file, got %v", parser.file)
	}
}

func TestParser2_1FileMovesToSnippetAfterParsingSnippetSPDXIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	fileCurrent := parser.file

	err := parser.parsePair2_1("SnippetSPDXID", "SPDXRef-Test1")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psSnippet2_1 {
		t.Errorf("expected state to be %v, got %v", psSnippet2_1, parser.st)
	}
	// and current file should remain what it was
	if parser.file != fileCurrent {
		t.Fatalf("expected file to remain %v, got %v", fileCurrent, parser.file)
	}
}

func TestParser2_1FileMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psOtherLicense2_1 {
		t.Errorf("expected state to be %v, got %v", psOtherLicense2_1, parser.st)
	}
}

func TestParser2_1FileMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePair2_1("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("expected state to be %v, got %v", psReview2_1, parser.st)
	}
}

func TestParser2_1FileStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePair2_1("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should remain unchanged
	if parser.st != psFile2_1 {
		t.Errorf("expected state to be %v, got %v", psFile2_1, parser.st)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psFile2_1 {
		t.Errorf("expected state to be %v, got %v", psFile2_1, parser.st)
	}
}

func TestParser2_1FileStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePair2_1("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}

	err = parser.parsePair2_1("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}

	err = parser.parsePair2_1("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}

	err = parser.parsePair2_1("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}

	err = parser.parsePair2_1("AnnotationComment", "i guess i had something to say about this particular file")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}
}

// ===== File data section tests =====
func TestParser2_1CanParseFileTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileName != "f1.txt" {
		t.Errorf("got %v for FileName", parser.file.FileName)
	}
	// should not yet be added to the Packages file list, because we haven't
	// seen an SPDX identifier yet
	if len(parser.pkg.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(parser.pkg.Files))
	}

	// File SPDX Identifier
	err = parser.parsePairFromFile2_1("SPDXID", "SPDXRef-f1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileSPDXIdentifier != "f1" {
		t.Errorf("got %v for FileSPDXIdentifier", parser.file.FileSPDXIdentifier)
	}
	// should now be added to the Packages file list
	if len(parser.pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(parser.pkg.Files))
	}
	if parser.pkg.Files[0] != parser.file {
		t.Errorf("expected Files[f1] to be %v, got %v", parser.file, parser.pkg.Files[0])
	}

	// File Type
	fileTypes := []string{
		"TEXT",
		"DOCUMENTATION",
	}
	for _, ty := range fileTypes {
		err = parser.parsePairFromFile2_1("FileType", ty)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, typeWant := range fileTypes {
		flagFound := false
		for _, typeCheck := range parser.file.FileTypes {
			if typeWant == typeCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in FileTypes", typeWant)
		}
	}
	if len(fileTypes) != len(parser.file.FileTypes) {
		t.Errorf("expected %d types in FileTypes, got %d", len(fileTypes),
			len(parser.file.FileTypes))
	}

	// File Checksums
	codeSha1 := "85ed0817af83a24ad8da68c2b5094de69833983c"
	sumSha1 := "SHA1: 85ed0817af83a24ad8da68c2b5094de69833983c"
	codeSha256 := "11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"
	sumSha256 := "SHA256: 11b6d3ee554eedf79299905a98f9b9a04e498210b59f15094c916c91d150efcd"
	codeMd5 := "624c1abb3664f4b35547e7c73864ad24"
	sumMd5 := "MD5: 624c1abb3664f4b35547e7c73864ad24"
	err = parser.parsePairFromFile2_1("FileChecksum", sumSha1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromFile2_1("FileChecksum", sumSha256)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromFile2_1("FileChecksum", sumMd5)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	for _, checksum := range parser.file.Checksums {
		switch checksum.Algorithm {
		case spdx.SHA1:
			if checksum.Value != codeSha1 {
				t.Errorf("expected %s for FileChecksumSHA1, got %s", codeSha1, checksum.Value)
			}
		case spdx.SHA256:
			if checksum.Value != codeSha256 {
				t.Errorf("expected %s for FileChecksumSHA1, got %s", codeSha256, checksum.Value)
			}
		case spdx.MD5:
			if checksum.Value != codeMd5 {
				t.Errorf("expected %s for FileChecksumSHA1, got %s", codeMd5, checksum.Value)
			}
		}
	}
	// Concluded License
	err = parser.parsePairFromFile2_1("LicenseConcluded", "Apache-2.0 OR GPL-2.0-or-later")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.LicenseConcluded != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("got %v for LicenseConcluded", parser.file.LicenseConcluded)
	}

	// License Information in File
	lics := []string{
		"Apache-2.0",
		"GPL-2.0-or-later",
		"CC0-1.0",
	}
	for _, lic := range lics {
		err = parser.parsePairFromFile2_1("LicenseInfoInFile", lic)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, licWant := range lics {
		flagFound := false
		for _, licCheck := range parser.file.LicenseInfoInFiles {
			if licWant == licCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in LicenseInfoInFiles", licWant)
		}
	}
	if len(lics) != len(parser.file.LicenseInfoInFiles) {
		t.Errorf("expected %d licenses in LicenseInfoInFiles, got %d", len(lics),
			len(parser.file.LicenseInfoInFiles))
	}

	// Comments on License
	err = parser.parsePairFromFile2_1("LicenseComments", "this is a comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.LicenseComments != "this is a comment" {
		t.Errorf("got %v for LicenseComments", parser.file.LicenseComments)
	}

	// Copyright Text
	err = parser.parsePairFromFile2_1("FileCopyrightText", "copyright (c) me")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileCopyrightText != "copyright (c) me" {
		t.Errorf("got %v for FileCopyrightText", parser.file.FileCopyrightText)
	}

	// Artifact of Projects: Name, HomePage and URI
	// Artifact set 1
	err = parser.parsePairFromFile2_1("ArtifactOfProjectName", "project1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromFile2_1("ArtifactOfProjectHomePage", "http://example.com/1/")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromFile2_1("ArtifactOfProjectURI", "http://example.com/1/uri.whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// Artifact set 2 -- just name
	err = parser.parsePairFromFile2_1("ArtifactOfProjectName", "project2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// Artifact set 3 -- just name and home page
	err = parser.parsePairFromFile2_1("ArtifactOfProjectName", "project3")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromFile2_1("ArtifactOfProjectHomePage", "http://example.com/3/")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// Artifact set 4 -- just name and URI
	err = parser.parsePairFromFile2_1("ArtifactOfProjectName", "project4")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromFile2_1("ArtifactOfProjectURI", "http://example.com/4/uri.whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if len(parser.file.ArtifactOfProjects) != 4 {
		t.Fatalf("expected len %d, got %d", 4, len(parser.file.ArtifactOfProjects))
	}

	aop := parser.file.ArtifactOfProjects[0]
	if aop.Name != "project1" {
		t.Errorf("expected %v, got %v", "project1", aop.Name)
	}
	if aop.HomePage != "http://example.com/1/" {
		t.Errorf("expected %v, got %v", "http://example.com/1/", aop.HomePage)
	}
	if aop.URI != "http://example.com/1/uri.whatever" {
		t.Errorf("expected %v, got %v", "http://example.com/1/uri.whatever", aop.URI)
	}

	aop = parser.file.ArtifactOfProjects[1]
	if aop.Name != "project2" {
		t.Errorf("expected %v, got %v", "project2", aop.Name)
	}
	if aop.HomePage != "" {
		t.Errorf("expected %v, got %v", "", aop.HomePage)
	}
	if aop.URI != "" {
		t.Errorf("expected %v, got %v", "", aop.URI)
	}

	aop = parser.file.ArtifactOfProjects[2]
	if aop.Name != "project3" {
		t.Errorf("expected %v, got %v", "project3", aop.Name)
	}
	if aop.HomePage != "http://example.com/3/" {
		t.Errorf("expected %v, got %v", "http://example.com/3/", aop.HomePage)
	}
	if aop.URI != "" {
		t.Errorf("expected %v, got %v", "", aop.URI)
	}

	aop = parser.file.ArtifactOfProjects[3]
	if aop.Name != "project4" {
		t.Errorf("expected %v, got %v", "project4", aop.Name)
	}
	if aop.HomePage != "" {
		t.Errorf("expected %v, got %v", "", aop.HomePage)
	}
	if aop.URI != "http://example.com/4/uri.whatever" {
		t.Errorf("expected %v, got %v", "http://example.com/4/uri.whatever", aop.URI)
	}

	// File Comment
	err = parser.parsePairFromFile2_1("FileComment", "this is a comment")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileComment != "this is a comment" {
		t.Errorf("got %v for FileComment", parser.file.FileComment)
	}

	// File Notice
	err = parser.parsePairFromFile2_1("FileNotice", "this is a Notice")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileNotice != "this is a Notice" {
		t.Errorf("got %v for FileNotice", parser.file.FileNotice)
	}

	// File Contributor
	contribs := []string{
		"John Doe jdoe@example.com",
		"EvilCorp",
	}
	for _, contrib := range contribs {
		err = parser.parsePairFromFile2_1("FileContributor", contrib)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, contribWant := range contribs {
		flagFound := false
		for _, contribCheck := range parser.file.FileContributors {
			if contribWant == contribCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in FileContributors", contribWant)
		}
	}
	if len(contribs) != len(parser.file.FileContributors) {
		t.Errorf("expected %d contribenses in FileContributors, got %d", len(contribs),
			len(parser.file.FileContributors))
	}

	// File Dependencies
	deps := []string{
		"f-1.txt",
		"g.txt",
	}
	for _, dep := range deps {
		err = parser.parsePairFromFile2_1("FileDependency", dep)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, depWant := range deps {
		flagFound := false
		for _, depCheck := range parser.file.FileDependencies {
			if depWant == depCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in FileDependency", depWant)
		}
	}
	if len(deps) != len(parser.file.FileDependencies) {
		t.Errorf("expected %d depenses in FileDependency, got %d", len(deps),
			len(parser.file.FileDependencies))
	}

}

func TestParser2_1FileCreatesRelationshipInDocument(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

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

func TestParser2_1FileCreatesAnnotationInDocument(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

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

func TestParser2_1FileUnknownTagFails(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePairFromFile2_1("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestFileAOPPointerChangesAfterTags(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
		file: &spdx.File2_1{FileName: "f1.txt", FileSPDXIdentifier: "f1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePairFromFile2_1("ArtifactOfProjectName", "project1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.fileAOP == nil {
		t.Errorf("expected non-nil AOP pointer, got nil")
	}
	curPtr := parser.fileAOP

	// now, a home page; pointer should stay
	err = parser.parsePairFromFile2_1("ArtifactOfProjectHomePage", "http://example.com/1/")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.fileAOP != curPtr {
		t.Errorf("expected no change in AOP pointer, was %v, got %v", curPtr, parser.fileAOP)
	}

	// a URI; pointer should stay
	err = parser.parsePairFromFile2_1("ArtifactOfProjectURI", "http://example.com/1/uri.whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.fileAOP != curPtr {
		t.Errorf("expected no change in AOP pointer, was %v, got %v", curPtr, parser.fileAOP)
	}

	// now, another artifact name; pointer should change but be non-nil
	// now, a home page; pointer should stay
	err = parser.parsePairFromFile2_1("ArtifactOfProjectName", "project2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.fileAOP == curPtr {
		t.Errorf("expected change in AOP pointer, got no change")
	}

	// finally, an unrelated tag; pointer should go away
	err = parser.parsePairFromFile2_1("FileComment", "whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.fileAOP != nil {
		t.Errorf("expected nil AOP pointer, got %v", parser.fileAOP)
	}
}

func TestParser2_1FailsIfInvalidSPDXIDInFileSection(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// start with File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid SPDX Identifier
	err = parser.parsePairFromFile2_1("SPDXID", "whoops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfInvalidChecksumFormatInFileSection(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// start with File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// invalid format for checksum line, missing colon
	err = parser.parsePairFromFile2_1("FileChecksum", "SHA1 85ed0817af83a24ad8da68c2b5094de69833983c")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfUnknownChecksumTypeInFileSection(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// start with File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// unknown checksum type
	err = parser.parsePairFromFile2_1("FileChecksum", "Special: 85ed0817af83a24ad8da68c2b5094de69833983c")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfArtifactHomePageBeforeArtifactName(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// start with File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// artifact home page appears before artifact name
	err = parser.parsePairFromFile2_1("ArtifactOfProjectHomePage", "https://example.com")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FailsIfArtifactURIBeforeArtifactName(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test", PackageSPDXIdentifier: "test", Files: []*spdx.File2_1{}},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)

	// start with File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// artifact home page appears before artifact name
	err = parser.parsePairFromFile2_1("ArtifactOfProjectURI", "https://example.com")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestParser2_1FilesWithoutSpdxIdThrowError(t *testing.T) {
	// case 1: The previous file (packaged or unpackaged) does not contain spdxID
	parser1 := tvParser2_1{
		doc:  &spdx.Document2_1{Packages: []*spdx.Package2_1{}},
		st:   psFile2_1,
		file: &spdx.File2_1{FileName: "FileName"},
	}

	err := parser1.parsePair2_1("FileName", "f2")
	if err == nil {
		t.Errorf("file without SPDX Identifier getting accepted")
	}

	// case 2: Invalid file with snippet
	// Last unpackaged file before a snippet starts
	sid1 := spdx.ElementID("s1")
	fileName := "f2.txt"
	parser2 := tvParser2_1{
		doc:  &spdx.Document2_1{},
		st:   psCreationInfo2_1,
		file: &spdx.File2_1{FileName: fileName},
	}
	err = parser2.parsePair2_1("SnippetSPDXID", string(sid1))
	if err == nil {
		t.Errorf("file without SPDX Identifier getting accepted")
	}

	// case 3 : Invalid File without snippets
	// Last unpackaged file before a package starts
	// Last file of a package and New package starts
	parser3 := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	fileName = "f3.txt"
	err = parser3.parsePair2_1("FileName", fileName)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = parser3.parsePair2_1("PackageName", "p2")
	if err == nil {
		t.Errorf("files withoutSpdx Identifier getting accepted")
	}
}
