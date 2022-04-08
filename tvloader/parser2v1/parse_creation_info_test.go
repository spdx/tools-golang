// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Parser creation info state change tests =====
func TestParser2_1CIMovesToPackageAfterParsingPackageNameTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	pkgName := "testPkg"
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
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the package should NOT be in the SPDX Document's map of packages,
	// because it doesn't have an SPDX identifier yet
	if len(parser.doc.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(parser.doc.Packages))
	}
}

func TestParser2_1CIMovesToFileAfterParsingFileNameTagWithNoPackages(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("FileName", "testFile")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psFile2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_1)
	}
	// and current package should be nil, meaning Files are placed in the
	// Files map instead of in a Package
	if parser.pkg != nil {
		t.Fatalf("expected pkg to be nil, got non-nil pkg")
	}
}

func TestParser2_1CIMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psOtherLicense2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense2_1)
	}
}

func TestParser2_1CIMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePair2_1("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_1)
	}
}

func TestParser2_1CIStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePair2_1("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}
}

func TestParser2_1CIStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePair2_1("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}

	err = parser.parsePair2_1("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}

	err = parser.parsePair2_1("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}

	err = parser.parsePair2_1("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}

	err = parser.parsePair2_1("AnnotationComment", "i guess i had something to say about this spdx file")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}
}

func TestParser2_1FailsParsingCreationInfoWithInvalidState(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psPackage2_1,
	}
	err := parser.parsePairFromCreationInfo2_1("SPDXVersion", "SPDX-2.1")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// ===== Creation Info section tests =====
func TestParser2_1HasCreationInfoAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}
	err := parser.parsePairFromCreationInfo2_1("LicenseListVersion", "3.9")
	if err != nil {
		t.Errorf("got error when calling parsePairFromCreationInfo2_1: %v", err)
	}
	if parser.doc.CreationInfo == nil {
		t.Errorf("doc.CreationInfo is still nil after parsing first pair")
	}
}

func TestParser2_1CanParseCreationInfoTags(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	// License List Version
	err := parser.parsePairFromCreationInfo2_1("LicenseListVersion", "2.2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.LicenseListVersion != "2.2" {
		t.Errorf("got %v for LicenseListVersion", parser.doc.CreationInfo.LicenseListVersion)
	}

	// Creators: Persons
	refPersons := []string{
		"Person: Person A",
		"Person: Person B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refPersons[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refPersons[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.Creators) != 2 ||
		parser.doc.CreationInfo.Creators[0].Creator != "Person A" ||
		parser.doc.CreationInfo.Creators[1].Creator != "Person B" {
		t.Errorf("got %+v for Creators", parser.doc.CreationInfo.Creators)
	}

	// Creators: Organizations
	refOrgs := []string{
		"Organization: Organization A",
		"Organization: Organization B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refOrgs[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refOrgs[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.Creators) != 4 ||
		parser.doc.CreationInfo.Creators[2].Creator != "Organization A" ||
		parser.doc.CreationInfo.Creators[3].Creator != "Organization B" {
		t.Errorf("got %+v for CreatorOrganizations", parser.doc.CreationInfo.Creators)
	}

	// Creators: Tools
	refTools := []string{
		"Tool: Tool A",
		"Tool: Tool B",
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refTools[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_1("Creator", refTools[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.Creators) != 6 ||
		parser.doc.CreationInfo.Creators[4].Creator != "Tool A" ||
		parser.doc.CreationInfo.Creators[5].Creator != "Tool B" {
		t.Errorf("got %v for CreatorTools", parser.doc.CreationInfo.Creators)
	}

	// Created date
	err = parser.parsePairFromCreationInfo2_1("Created", "2018-09-10T11:46:00Z")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.Created != "2018-09-10T11:46:00Z" {
		t.Errorf("got %v for Created", parser.doc.CreationInfo.Created)
	}

	// Creator Comment
	err = parser.parsePairFromCreationInfo2_1("CreatorComment", "Blah whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.CreatorComment != "Blah whatever" {
		t.Errorf("got %v for CreatorComment", parser.doc.CreationInfo.CreatorComment)
	}
}

func TestParser2_1InvalidCreatorTagsFail(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePairFromCreationInfo2_1("Creator", "blah: somebody")
	if err == nil {
		t.Errorf("expected error from parsing invalid Creator format, got nil")
	}

	err = parser.parsePairFromCreationInfo2_1("Creator", "Tool with no colons")
	if err == nil {
		t.Errorf("expected error from parsing invalid Creator format, got nil")
	}
}

func TestParser2_1CreatorTagWithMultipleColonsPasses(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePairFromCreationInfo2_1("Creator", "Tool: tool1:2:3")
	if err != nil {
		t.Errorf("unexpected error from parsing valid Creator format")
	}
}

func TestParser2_1CIUnknownTagFails(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

	err := parser.parsePairFromCreationInfo2_1("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestParser2_1CICreatesRelationship(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

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

func TestParser2_1CICreatesAnnotation(t *testing.T) {
	parser := tvParser2_1{
		doc: &spdx.Document2_1{},
		st:  psCreationInfo2_1,
	}

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

// ===== Helper function tests =====

func TestCanExtractExternalDocumentReference(t *testing.T) {
	refstring := "DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 SHA1:d6a770ba38583ed4bb4525bd96e50461655d2759"
	wantDocumentRefID := "spdx-tool-1.2"
	wantURI := "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"
	wantAlg := "SHA1"
	wantChecksum := "d6a770ba38583ed4bb4525bd96e50461655d2759"

	gotDocumentRefID, gotURI, gotAlg, gotChecksum, err := extractExternalDocumentReference(refstring)
	if err != nil {
		t.Errorf("got non-nil error: %v", err)
	}
	if wantDocumentRefID != gotDocumentRefID {
		t.Errorf("wanted document ref ID %s, got %s", wantDocumentRefID, gotDocumentRefID)
	}
	if wantURI != gotURI {
		t.Errorf("wanted URI %s, got %s", wantURI, gotURI)
	}
	if wantAlg != gotAlg {
		t.Errorf("wanted alg %s, got %s", wantAlg, gotAlg)
	}
	if wantChecksum != gotChecksum {
		t.Errorf("wanted checksum %s, got %s", wantChecksum, gotChecksum)
	}
}

func TestCanExtractExternalDocumentReferenceWithExtraWhitespace(t *testing.T) {
	refstring := "   DocumentRef-spdx-tool-1.2    \t http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 \t SHA1:  \t   d6a770ba38583ed4bb4525bd96e50461655d2759"
	wantDocumentRefID := "spdx-tool-1.2"
	wantURI := "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"
	wantAlg := "SHA1"
	wantChecksum := "d6a770ba38583ed4bb4525bd96e50461655d2759"

	gotDocumentRefID, gotURI, gotAlg, gotChecksum, err := extractExternalDocumentReference(refstring)
	if err != nil {
		t.Errorf("got non-nil error: %v", err)
	}
	if wantDocumentRefID != gotDocumentRefID {
		t.Errorf("wanted document ref ID %s, got %s", wantDocumentRefID, gotDocumentRefID)
	}
	if wantURI != gotURI {
		t.Errorf("wanted URI %s, got %s", wantURI, gotURI)
	}
	if wantAlg != gotAlg {
		t.Errorf("wanted alg %s, got %s", wantAlg, gotAlg)
	}
	if wantChecksum != gotChecksum {
		t.Errorf("wanted checksum %s, got %s", wantChecksum, gotChecksum)
	}
}

func TestFailsExternalDocumentReferenceWithInvalidFormats(t *testing.T) {
	invalidRefs := []string{
		"whoops",
		"DocumentRef-",
		"DocumentRef-   ",
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 d6a770ba38583ed4bb4525bd96e50461655d2759",
		"DocumentRef-spdx-tool-1.2",
		"spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 SHA1:d6a770ba38583ed4bb4525bd96e50461655d2759",
	}
	for _, refstring := range invalidRefs {
		_, _, _, _, err := extractExternalDocumentReference(refstring)
		if err == nil {
			t.Errorf("expected non-nil error for %s, got nil", refstring)
		}
	}
}
