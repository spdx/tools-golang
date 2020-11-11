// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v2

import (
	"testing"

	"github.com/spdx/tools-golang/spdx"
)

// ===== Parser creation info state change tests =====
func TestParser2_2CIMovesToPackageAfterParsingPackageNameTag(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	pkgName := "testPkg"
	err := parser.parsePair2_2("PackageName", pkgName)
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	// state should be correct
	if parser.st != psPackage2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psPackage2_2)
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

func TestParser2_2CIMovesToFileAfterParsingFileNameTagWithNoPackages(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("FileName", "testFile")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	// state should be correct
	if parser.st != psFile2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psFile2_2)
	}
	// and current package should be nil, meaning Files are placed in the
	// UnpackagedFiles map instead of in a Package
	if parser.pkg != nil {
		t.Fatalf("expected pkg to be nil, got non-nil pkg")
	}
}

func TestParser2_2CIMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psOtherLicense2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psOtherLicense2_2)
	}
}

func TestParser2_2CIMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePair2_2("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psReview2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psReview2_2)
	}
}

func TestParser2_2CIStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePair2_2("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}

	err = parser.parsePair2_2("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}
}

func TestParser2_2CIStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePair2_2("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}

	err = parser.parsePair2_2("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}

	err = parser.parsePair2_2("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}

	err = parser.parsePair2_2("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}

	err = parser.parsePair2_2("AnnotationComment", "i guess i had something to say about this spdx file")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.st != psCreationInfo2_2 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_2)
	}
}

func TestParser2_2FailsParsingCreationInfoWithInvalidState(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psPackage2_2,
	}
	err := parser.parsePairFromCreationInfo2_2("SPDXVersion", "SPDX-2.2")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// ===== Creation Info section tests =====
func TestParser2_2HasCreationInfoAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}
	err := parser.parsePairFromCreationInfo2_2("SPDXVersion", "SPDX-2.2")
	if err != nil {
		t.Errorf("got error when calling parsePairFromCreationInfo2_2: %v", err)
	}
	if parser.doc.CreationInfo == nil {
		t.Errorf("doc.CreationInfo is still nil after parsing first pair")
	}
}

func TestParser2_2CanParseCreationInfoTags(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	// SPDX Version
	err := parser.parsePairFromCreationInfo2_2("SPDXVersion", "SPDX-2.2")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.SPDXVersion != "SPDX-2.2" {
		t.Errorf("got %v for SPDXVersion", parser.doc.CreationInfo.SPDXVersion)
	}

	// Data License
	err = parser.parsePairFromCreationInfo2_2("DataLicense", "CC0-1.0")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DataLicense != "CC0-1.0" {
		t.Errorf("got %v for DataLicense", parser.doc.CreationInfo.DataLicense)
	}

	// SPDX Identifier
	err = parser.parsePairFromCreationInfo2_2("SPDXID", "SPDXRef-DOCUMENT")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.SPDXIdentifier != "DOCUMENT" {
		t.Errorf("got %v for SPDXIdentifier", parser.doc.CreationInfo.SPDXIdentifier)
	}

	// Document Name
	err = parser.parsePairFromCreationInfo2_2("DocumentName", "xyz-2.1.5")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DocumentName != "xyz-2.1.5" {
		t.Errorf("got %v for DocumentName", parser.doc.CreationInfo.DocumentName)
	}

	// Document Namespace
	err = parser.parsePairFromCreationInfo2_2("DocumentNamespace", "http://example.com/xyz-2.1.5.spdx")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DocumentNamespace != "http://example.com/xyz-2.1.5.spdx" {
		t.Errorf("got %v for DocumentNamespace", parser.doc.CreationInfo.DocumentNamespace)
	}

	// External Document Reference
	refs := []string{
		"DocumentRef-spdx-tool-1.2 http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301 SHA1: d6a770ba38583ed4bb4525bd96e50461655d2759",
		"DocumentRef-xyz-2.1.2 http://example.com/xyz-2.1.2 SHA1:d6a770ba38583ed4bb4525bd96e50461655d2760",
	}
	wantRef0 := spdx.ExternalDocumentRef2_2{
		DocumentRefID: "spdx-tool-1.2",
		URI:           "http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301",
		Alg:           "SHA1",
		Checksum:      "d6a770ba38583ed4bb4525bd96e50461655d2759",
	}
	wantRef1 := spdx.ExternalDocumentRef2_2{
		DocumentRefID: "xyz-2.1.2",
		URI:           "http://example.com/xyz-2.1.2",
		Alg:           "SHA1",
		Checksum:      "d6a770ba38583ed4bb4525bd96e50461655d2760",
	}
	err = parser.parsePairFromCreationInfo2_2("ExternalDocumentRef", refs[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_2("ExternalDocumentRef", refs[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.ExternalDocumentReferences) != 2 {
		t.Errorf("got %d ExternalDocumentReferences, expected %d", len(parser.doc.CreationInfo.ExternalDocumentReferences), 2)
	}
	gotRef0 := parser.doc.CreationInfo.ExternalDocumentReferences["spdx-tool-1.2"]
	if gotRef0.DocumentRefID != wantRef0.DocumentRefID ||
		gotRef0.URI != wantRef0.URI ||
		gotRef0.Alg != wantRef0.Alg ||
		gotRef0.Checksum != wantRef0.Checksum {
		t.Errorf("got %#v for ExternalDocumentReferences[0], wanted %#v", gotRef0, wantRef0)
	}
	gotRef1 := parser.doc.CreationInfo.ExternalDocumentReferences["xyz-2.1.2"]
	if gotRef1.DocumentRefID != wantRef1.DocumentRefID ||
		gotRef1.URI != wantRef1.URI ||
		gotRef1.Alg != wantRef1.Alg ||
		gotRef1.Checksum != wantRef1.Checksum {
		t.Errorf("got %#v for ExternalDocumentReferences[1], wanted %#v", gotRef1, wantRef1)
	}

	// License List Version
	err = parser.parsePairFromCreationInfo2_2("LicenseListVersion", "2.2")
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
	err = parser.parsePairFromCreationInfo2_2("Creator", refPersons[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_2("Creator", refPersons[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.CreatorPersons) != 2 ||
		parser.doc.CreationInfo.CreatorPersons[0] != "Person A" ||
		parser.doc.CreationInfo.CreatorPersons[1] != "Person B" {
		t.Errorf("got %v for CreatorPersons", parser.doc.CreationInfo.CreatorPersons)
	}

	// Creators: Organizations
	refOrgs := []string{
		"Organization: Organization A",
		"Organization: Organization B",
	}
	err = parser.parsePairFromCreationInfo2_2("Creator", refOrgs[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_2("Creator", refOrgs[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.CreatorOrganizations) != 2 ||
		parser.doc.CreationInfo.CreatorOrganizations[0] != "Organization A" ||
		parser.doc.CreationInfo.CreatorOrganizations[1] != "Organization B" {
		t.Errorf("got %v for CreatorOrganizations", parser.doc.CreationInfo.CreatorOrganizations)
	}

	// Creators: Tools
	refTools := []string{
		"Tool: Tool A",
		"Tool: Tool B",
	}
	err = parser.parsePairFromCreationInfo2_2("Creator", refTools[0])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = parser.parsePairFromCreationInfo2_2("Creator", refTools[1])
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(parser.doc.CreationInfo.CreatorTools) != 2 ||
		parser.doc.CreationInfo.CreatorTools[0] != "Tool A" ||
		parser.doc.CreationInfo.CreatorTools[1] != "Tool B" {
		t.Errorf("got %v for CreatorTools", parser.doc.CreationInfo.CreatorTools)
	}

	// Created date
	err = parser.parsePairFromCreationInfo2_2("Created", "2018-09-10T11:46:00Z")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.Created != "2018-09-10T11:46:00Z" {
		t.Errorf("got %v for Created", parser.doc.CreationInfo.Created)
	}

	// Creator Comment
	err = parser.parsePairFromCreationInfo2_2("CreatorComment", "Blah whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.CreatorComment != "Blah whatever" {
		t.Errorf("got %v for CreatorComment", parser.doc.CreationInfo.CreatorComment)
	}

	// Document Comment
	err = parser.parsePairFromCreationInfo2_2("DocumentComment", "Blah whatever")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.doc.CreationInfo.DocumentComment != "Blah whatever" {
		t.Errorf("got %v for DocumentComment", parser.doc.CreationInfo.DocumentComment)
	}

}

func TestParser2_2InvalidCreatorTagsFail(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePairFromCreationInfo2_2("Creator", "blah: somebody")
	if err == nil {
		t.Errorf("expected error from parsing invalid Creator format, got nil")
	}

	err = parser.parsePairFromCreationInfo2_2("Creator", "Tool with no colons")
	if err == nil {
		t.Errorf("expected error from parsing invalid Creator format, got nil")
	}
}

func TestParser2_2CreatorTagWithMultipleColonsPasses(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePairFromCreationInfo2_2("Creator", "Tool: tool1:2:3")
	if err != nil {
		t.Errorf("unexpected error from parsing valid Creator format")
	}
}

func TestParser2_2CIUnknownTagFails(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePairFromCreationInfo2_2("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}

func TestParser2_2CICreatesRelationship(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePair2_2("Relationship", "SPDXRef-blah CONTAINS SPDXRef-blah-whatever")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
	}
	if parser.rln == nil {
		t.Fatalf("parser didn't create and point to Relationship struct")
	}
	if parser.rln != parser.doc.Relationships[0] {
		t.Errorf("pointer to new Relationship doesn't match idx 0 for doc.Relationships[]")
	}
}

func TestParser2_2CICreatesAnnotation(t *testing.T) {
	parser := tvParser2_2{
		doc: &spdx.Document2_2{},
		st:  psCreationInfo2_2,
	}

	err := parser.parsePair2_2("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_2: %v", err)
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
