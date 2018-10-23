// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// ===== Parser file section state change tests =====
func TestParser2_1FileStartsNewFileAfterParsingFileNameTag(t *testing.T) {
	// create the first file
	fileOldName := "f1.txt"

	parser := tvParser2_1{
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: fileOldName},
	}
	fileOld := parser.file
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, fileOld)
	// the Package's Files should have this one only
	if parser.pkg.Files[0] != fileOld {
		t.Errorf("Expected file %v in Files[0], got %v", fileOld, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != fileOldName {
		t.Errorf("expected file name %s in Files[0], got %s", fileOldName, parser.pkg.Files[0].FileName)
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
	// and the Package's Files should be of size 2 and have these two
	if parser.pkg.Files[0] != fileOld {
		t.Errorf("Expected file %v in Files[0], got %v", fileOld, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != fileOldName {
		t.Errorf("expected file name %s in Files[0], got %s", fileOldName, parser.pkg.Files[0].FileName)
	}
	if parser.pkg.Files[1] != parser.file {
		t.Errorf("Expected file %v in Files[1], got %v", parser.file, parser.pkg.Files[1])
	}
	if parser.pkg.Files[1].FileName != fileName {
		t.Errorf("expected file name %s in Files[1], got %s", fileName, parser.pkg.Files[1].FileName)
	}
}

func TestParser2_1FileStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	// create the first file and package
	p1Name := "package1"
	f1Name := "f1.txt"

	parser := tvParser2_1{
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: p1Name},
		file: &spdx.File2_1{FileName: f1Name},
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
	if parser.doc.Packages[0] != p1 {
		t.Errorf("Expected package %v in Packages[0], got %v", p1, parser.doc.Packages[0])
	}
	if parser.doc.Packages[0].PackageName != p1Name {
		t.Errorf("expected package name %s in Packages[0], got %s", p1Name, parser.doc.Packages[0].PackageName)
	}
	if parser.doc.Packages[1] != parser.pkg {
		t.Errorf("Expected package %v in Packages[1], got %v", parser.pkg, parser.doc.Packages[1])
	}
	if parser.doc.Packages[1].PackageName != p2Name {
		t.Errorf("expected package name %s in Packages[1], got %s", p2Name, parser.doc.Packages[1].PackageName)
	}
	// and the first Package's Files should be of size 1 and have f1 only
	if len(parser.doc.Packages[0].Files) != 1 {
		t.Errorf("Expected 1 file in Packages[0].Files, got %d", len(parser.doc.Packages[0].Files))
	}
	if parser.doc.Packages[0].Files[0] != f1 {
		t.Errorf("Expected file %v in Files[0], got %v", f1, parser.doc.Packages[0].Files[0])
	}
	if parser.doc.Packages[0].Files[0].FileName != f1Name {
		t.Errorf("expected file name %s in Files[0], got %s", f1Name, parser.doc.Packages[0].Files[0].FileName)
	}
	// and the second Package should have no files
	if len(parser.doc.Packages[1].Files) != 0 {
		t.Errorf("Expected no files in Packages[1].Files, got %d", len(parser.doc.Packages[1].Files))
	}
	// and the current file should be nil
	if parser.file != nil {
		t.Errorf("Expected nil for parser.file, got %v", parser.file)
	}
}

func TestParser2_1FileMovesToSnippetAfterParsingSnippetSPDXIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-else")
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
		doc: &spdx.Document2_1{},
		st:  psFile2_1,
		pkg: &spdx.Package2_1{PackageName: "test"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	// File Name
	err := parser.parsePairFromFile2_1("FileName", "f1.txt")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileName != "f1.txt" {
		t.Errorf("got %v for FileName", parser.file.FileName)
	}

	// File SPDX Identifier
	err = parser.parsePairFromFile2_1("SPDXID", "SPDXRef-f1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.file.FileSPDXIdentifier != "SPDXRef-f1" {
		t.Errorf("got %v for FileSPDXIdentifier", parser.file.FileSPDXIdentifier)
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
		for _, typeCheck := range parser.file.FileType {
			if typeWant == typeCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in FileType", typeWant)
		}
	}
	if len(fileTypes) != len(parser.file.FileType) {
		t.Errorf("expected %d types in FileType, got %d", len(fileTypes),
			len(parser.file.FileType))
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
	if parser.file.FileChecksumSHA1 != codeSha1 {
		t.Errorf("expected %s for FileChecksumSHA1, got %s", codeSha1, parser.file.FileChecksumSHA1)
	}
	if parser.file.FileChecksumSHA256 != codeSha256 {
		t.Errorf("expected %s for FileChecksumSHA256, got %s", codeSha256, parser.file.FileChecksumSHA256)
	}
	if parser.file.FileChecksumMD5 != codeMd5 {
		t.Errorf("expected %s for FileChecksumMD5, got %s", codeMd5, parser.file.FileChecksumMD5)
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
		for _, licCheck := range parser.file.LicenseInfoInFile {
			if licWant == licCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in LicenseInfoInFile", licWant)
		}
	}
	if len(lics) != len(parser.file.LicenseInfoInFile) {
		t.Errorf("expected %d licenses in LicenseInfoInFile, got %d", len(lics),
			len(parser.file.LicenseInfoInFile))
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
		for _, contribCheck := range parser.file.FileContributor {
			if contribWant == contribCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in FileContributor", contribWant)
		}
	}
	if len(contribs) != len(parser.file.FileContributor) {
		t.Errorf("expected %d contribenses in FileContributor, got %d", len(contribs),
			len(parser.file.FileContributor))
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)

	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-whatever")
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
		doc:  &spdx.Document2_1{},
		st:   psFile2_1,
		pkg:  &spdx.Package2_1{PackageName: "test"},
		file: &spdx.File2_1{FileName: "f1.txt"},
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
