// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"reflect"
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func Test_setPackageSupplier(t *testing.T) {
	var err error

	// TestCase 1: no assertion must set PackageSupplierNOASSERTION field to true
	pkg := &spdx.Package{}
	err = setPackageSupplier(pkg, "NOASSERTION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pkg.PackageSupplier.Supplier != "NOASSERTION" {
		t.Errorf("PackageSupplier must've been set to NOASSERTION")
	}

	// TestCase 2: lower-case noassertion must also set the
	// PackageSupplierNOASSERTION to true.
	pkg = &spdx.Package{}
	err = setPackageSupplier(pkg, "noassertion")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pkg.PackageSupplier.Supplier != "NOASSERTION" {
		t.Errorf("PackageSupplier must've been set to NOASSERTION")
	}

	// TestCase 3: invalid input without colon separator. must raise an error
	pkg = &spdx.Package{}
	input := "string without colon separator"
	err = setPackageSupplier(pkg, input)
	if err == nil {
		t.Errorf("invalid input \"%s\" didn't raise an error", input)
	}

	// TestCase 4: Valid Person
	pkg = &spdx.Package{}
	personName := "Rishabh Bhatnagar"
	input = "Person: " + personName
	err = setPackageSupplier(pkg, input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pkg.PackageSupplier.Supplier != personName {
		t.Errorf("PackageSupplierPerson should be %s. found %s", personName, pkg.PackageSupplier.Supplier)
	}

	// TestCase 5: Valid Organization
	pkg = &spdx.Package{}
	orgName := "SPDX"
	input = "Organization: " + orgName
	err = setPackageSupplier(pkg, input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pkg.PackageSupplier.Supplier != orgName {
		t.Errorf("PackageSupplierPerson should be %s. found %s", orgName, pkg.PackageSupplier.Supplier)
	}

	// TestCase 6: Invalid EntityType
	pkg = &spdx.Package{}
	input = "InvalidEntity: entity"
	err = setPackageSupplier(pkg, input)
	if err == nil {
		t.Errorf("invalid entity should've raised an error")
	}
}

func Test_setPackageOriginator(t *testing.T) {
	var err error

	// TestCase 1: no assertion must set PackageSupplierNOASSERTION field to true
	pkg := &spdx.Package{}
	err = setPackageOriginator(pkg, "NOASSERTION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pkg.PackageOriginator.Originator != "NOASSERTION" {
		t.Errorf("PackageOriginator must've been set to NOASSERTION")
	}

	// TestCase 2: lower-case noassertion must also set the
	// PackageOriginatorNOASSERTION to true.
	pkg = &spdx.Package{}
	err = setPackageOriginator(pkg, "noassertion")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pkg.PackageOriginator.Originator != "NOASSERTION" {
		t.Errorf("PackageOriginator must've been set to NOASSERTION")
	}

	// TestCase 3: invalid input without colon separator. must raise an error
	pkg = &spdx.Package{}
	input := "string without colon separator"
	err = setPackageOriginator(pkg, input)
	if err == nil {
		t.Errorf("invalid input \"%s\" didn't raise an error", input)
	}

	// TestCase 4: Valid Person
	pkg = &spdx.Package{}
	personName := "Rishabh Bhatnagar"
	input = "Person: " + personName
	err = setPackageOriginator(pkg, input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pkg.PackageOriginator.Originator != personName {
		t.Errorf("PackageOriginatorPerson should be %s. found %s", personName, pkg.PackageOriginator.Originator)
	}

	// TestCase 5: Valid Organization
	pkg = &spdx.Package{}
	orgName := "SPDX"
	input = "Organization: " + orgName
	err = setPackageOriginator(pkg, input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if pkg.PackageOriginator.Originator != orgName {
		t.Errorf("PackageOriginatorOrganization should be %s. found %s", orgName, pkg.PackageOriginator.Originator)
	}

	// TestCase 6: Invalid EntityType
	pkg = &spdx.Package{}
	input = "InvalidEntity: entity"
	err = setPackageOriginator(pkg, input)
	if err == nil {
		t.Errorf("invalid entity should've raised an error")
	}
}

func Test_rdfParser2_3_setPackageVerificationCode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var pkg *spdx.Package
	var err error

	// TestCase 1: invalid predicate must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx.PackageVerificationCode>
			<spdx:invalidPredicate />
			<spdx:packageVerificationCodeValue>cbceb8b5689b75a584efe35587b5d41bd48820ce</spdx:packageVerificationCodeValue>
			<spdx:packageVerificationCodeExcludedFile>./package.spdx</spdx:packageVerificationCodeExcludedFile>
		</spdx.PackageVerificationCode>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	pkg = &spdx.Package{}
	err = parser.setPackageVerificationCode(pkg, node)
	if err == nil {
		t.Errorf("expected an error due to invalid predicate, got <nil>")
	}

	// TestCase 2: valid input
	parser, _ = parserFromBodyContent(`
		<spdx.PackageVerificationCode>
			<spdx:packageVerificationCodeValue>cbceb8b5689b75a584efe35587b5d41bd48820ce</spdx:packageVerificationCodeValue>
			<spdx:packageVerificationCodeExcludedFile>./package.spdx</spdx:packageVerificationCodeExcludedFile>
		</spdx.PackageVerificationCode>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	pkg = &spdx.Package{}
	err = parser.setPackageVerificationCode(pkg, node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedValue := "cbceb8b5689b75a584efe35587b5d41bd48820ce"
	if pkg.PackageVerificationCode.Value != expectedValue {
		t.Errorf("expected %v, got %v", expectedValue, pkg.PackageVerificationCode)
	}
	expectedExcludedFile := "./package.spdx"
	if pkg.PackageVerificationCode.ExcludedFiles[0] != expectedExcludedFile {
		t.Errorf("expected %v, got %v", expectedExcludedFile, pkg.PackageVerificationCode.ExcludedFiles)
	}
}

func Test_rdfParser2_3_getPackageExternalRef(t *testing.T) {
	var extRef *spdx.PackageExternalReference
	var err error
	var parser *rdfParser2_3
	var node *gordfParser.Node

	// TestCase 1: invalid reference category
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalRef>
			<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
			<spdx:referenceType>
				<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
			</spdx:referenceType>
			<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_invalid"/>
		</spdx:ExternalRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	extRef, err = parser.getPackageExternalRef(node)
	if err == nil {
		t.Errorf("expected an error due to invalid referenceCategory, got <nil>")
	}

	// TestCase 2: invalid predicate
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalRef>
			<spdx:unknownPredicate />
			<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
			<spdx:referenceType>
				<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
			</spdx:referenceType>
			<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_security"/>
		</spdx:ExternalRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	extRef, err = parser.getPackageExternalRef(node)
	if err == nil {
		t.Errorf("expected an error due to invalid referenceCategory, got <nil>")
	}

	// TestCase 3: valid example (referenceCategory_security)
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalRef>
			<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
			<spdx:referenceType>
				<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
			</spdx:referenceType>
			<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_security"/>
			<rdfs:comment>comment</rdfs:comment>
		</spdx:ExternalRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	extRef, err = parser.getPackageExternalRef(node)
	if err != nil {
		t.Fatalf("unexpected error parsing a valid example: %v", err)
	}
	expectedExtRef := &spdx.PackageExternalReference{
		Locator:            "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
		RefType:            "http://spdx.org/rdf/references/cpe23Type",
		Category:           "SECURITY",
		ExternalRefComment: "comment",
	}
	if !reflect.DeepEqual(extRef, expectedExtRef) {
		t.Errorf("expected: \n%+v\ngot: \n%+v", expectedExtRef, extRef)
	}

	// TestCase 4: valid example (referenceCategory_packageManager)
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalRef>
			<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
			<spdx:referenceType>
				<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
			</spdx:referenceType>
			<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_packageManager"/>
			<rdfs:comment>comment</rdfs:comment>
		</spdx:ExternalRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	extRef, err = parser.getPackageExternalRef(node)
	if err != nil {
		t.Fatalf("unexpected error parsing a valid example: %v", err)
	}
	expectedExtRef = &spdx.PackageExternalReference{
		Locator:            "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
		RefType:            "http://spdx.org/rdf/references/cpe23Type",
		Category:           "PACKAGE-MANAGER",
		ExternalRefComment: "comment",
	}
	if !reflect.DeepEqual(extRef, expectedExtRef) {
		t.Errorf("expected: \n%+v\ngot: \n%+v", expectedExtRef, extRef)
	}

	// TestCase 5: valid example (referenceCategory_other)
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalRef>
			<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
			<spdx:referenceType>
				<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
			</spdx:referenceType>
			<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_other"/>
			<rdfs:comment>comment</rdfs:comment>
		</spdx:ExternalRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	extRef, err = parser.getPackageExternalRef(node)
	if err != nil {
		t.Fatalf("unexpected error parsing a valid example: %v", err)
	}
	expectedExtRef = &spdx.PackageExternalReference{
		Locator:            "cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*",
		RefType:            "http://spdx.org/rdf/references/cpe23Type",
		Category:           "OTHER",
		ExternalRefComment: "comment",
	}
	if !reflect.DeepEqual(extRef, expectedExtRef) {
		t.Errorf("expected: \n%+v\ngot: \n%+v", expectedExtRef, extRef)
	}
}

func Test_rdfParser2_3_getPrimaryPackagePurpose(t *testing.T) {
	// TestCase 1: basic purpose
	value := getPrimaryPackagePurpose("packagePurpose_container")
	if value != "CONTAINER" {
		t.Errorf("expected primary package purpose to be CONTAINER. got: '%s'", value)
	}

	// TestCase 2: purpose with underscore-to-dash
	value = getPrimaryPackagePurpose("packagePurpose_operating_system")
	if value != "OPERATING-SYSTEM" {
		t.Errorf("expected primary package purpose to be OPERATING-SYSTEM. got: '%s'", value)
	}

	// TestCase 3: invalid purpose
	value = getPrimaryPackagePurpose("packagePurpose_invalid")
	if value != "" {
		t.Errorf("expected invalid primary package purpose to be empty. got: '%s'", value)
	}
}

func Test_rdfParser2_3_getPackageFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var err error

	// TestCase 1: invalid elementId
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#upload2">
            <spdx:name>time-1.9.tar.gz</spdx:name>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(missing SPDXRef- prefix), found %v", err)
	}

	// TestCase 2: Invalid License Concluded must raise an error:
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
            <spdx:licenseConcluded rdf:resource="http://spdx.org/licenses/IPL-3.0"/>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid license), found %v", err)
	}

	// TestCase 2: Invalid License Declared must raise an error:
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
            <spdx:licenseDeclared rdf:resource="http://spdx.org/licenses/IPL-3.0"/>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid license), found %v", err)
	}

	// TestCase 3: Invalid ExternalRef
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:externalRef>			
				<spdx:ExternalRef>
					<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
					<spdx:referenceType>
						<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
					</spdx:referenceType>
					<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_invalid"/>
				</spdx:ExternalRef>
			</spdx:externalRef>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid externalRef), found %v", err)
	}

	// TestCase 4: invalid file must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:hasFile>
              <spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#item8"/>
            </spdx:hasFile>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid file), found %v", err)
	}

	// TestCase 5: invalid predicate must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:hasFiles>
              <spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#item8"/>
            </spdx:hasFiles>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid predicate), found %v", err)
	}

	// TestCase 6: invalid annotation must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:annotation>
				<spdx:Annotation>
					<spdx:unknownAttribute />
				</spdx:Annotation>
			</spdx:annotation>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid annotation), found %v", err)
	}

	// TestCase 6: invalid homepage must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<doap:homepage>u r i</doap:homepage>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err == nil {
		t.Errorf("expected an error(invalid homepage uri), found %v", err)
	}

	// TestCase 7: Package tag declared more than once should be parsed into a single object's definition
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:name>Test Package</spdx:name>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid package: %v", err)
	}
	yetAnotherPkgTriple := gordfParser.Triple{
		Subject: node,
		Predicate: &gordfParser.Node{
			NodeType: gordfParser.IRI,
			ID:       SPDX_PACKAGE_FILE_NAME,
		},
		Object: &gordfParser.Node{
			NodeType: gordfParser.LITERAL,
			ID:       "packageFileName",
		},
	}
	parser.nodeStringToTriples[node.String()] = append(parser.nodeStringToTriples[node.String()], &yetAnotherPkgTriple)
	pkg, err := parser.getPackageFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid package: %v", err)
	}
	// validating if all the attributes that spanned over two tags are included in the parsed package.
	expectedID := "upload2"
	if string(pkg.PackageSPDXIdentifier) != expectedID {
		t.Errorf("expected package id: %s, got %s", expectedID, pkg.PackageSPDXIdentifier)
	}
	expectedPkgFileName := "packageFileName"
	if expectedPkgFileName != pkg.PackageFileName {
		t.Errorf("expected package file name: %s, got %s", expectedPkgFileName, pkg.PackageFileName)
	}
	expectedName := "Test Package"
	if pkg.PackageName != expectedName {
		t.Errorf("expected package name: %s, got %s", expectedPkgFileName, pkg.PackageName)
	}

	// TestCase 8: Checking if packages can handle cyclic dependencies:
	// Simulating a smallest possible cycle: package related to itself.
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:name>Test Package</spdx:name>
			<spdx:relationship>
			    <spdx:Relationship>
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_describes" />
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
							<spdx:versionInfo>1.1.1</spdx:versionInfo>
						</spdx:Package>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	pkg, err = parser.getPackageFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid package: %v", err)
	}
	// checking if both the attributes of the packages are set.
	expectedVersionInfo := "1.1.1"
	expectedPackageName := "Test Package"
	if pkg.PackageVersion != expectedVersionInfo {
		t.Errorf("Expected %s, found %s", expectedVersionInfo, pkg.PackageVersion)
	}
	if pkg.PackageName != expectedPackageName {
		t.Errorf("Expected %s, found %s", expectedPackageName, pkg.PackageName)
	}

	// TestCase 9: everything valid
	parser, _ = parserFromBodyContent(`
		<spdx:Package rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2">
			<spdx:name>Test Package</spdx:name>
			<spdx:versionInfo>1.1.1</spdx:versionInfo>
			<spdx:packageFileName>time-1.9.tar.gz</spdx:packageFileName>
			<spdx:supplier>Person: Jane Doe (jane.doe@example.com)</spdx:supplier>
			<spdx:originator>Organization: SPDX</spdx:originator>
			<spdx:downloadLocation rdf:resource="http://spdx.org/rdf/terms#noassertion" />
			<spdx:filesAnalyzed>true</spdx:filesAnalyzed>
			<spdx:packageVerificationCode>
                <spdx.PackageVerificationCode>
                    <spdx:packageVerificationCodeValue>cbceb8b5689b75a584efe35587b5d41bd48820ce</spdx:packageVerificationCodeValue>
					<spdx:packageVerificationCodeExcludedFile>./package.spdx</spdx:packageVerificationCodeExcludedFile>
                </spdx.PackageVerificationCode>
            </spdx:packageVerificationCode>
			<spdx:checksum>
                <spdx:Checksum>
					<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1" />
					<spdx:checksumValue>75068c26abbed3ad3980685bae21d7202d288317</spdx:checksumValue>
                </spdx:Checksum>
            </spdx:checksum>
			<doap:homepage>http://www.openjena.org/</doap:homepage>
			<spdx:sourceInfo>uses glibc-2_11-branch from git://sourceware.org/git/glibc.git.</spdx:sourceInfo>
			<spdx:licenseConcluded>
                <spdx:DisjunctiveLicenseSet>
					<spdx:member rdf:resource="http://spdx.org/licenses/Nokia"/>
					<spdx:member rdf:resource="http://spdx.org/licenses/LGPL-2.0"/>
                </spdx:DisjunctiveLicenseSet>
            </spdx:licenseConcluded>
			<spdx:licenseInfoFromFiles rdf:resource="http://spdx.org/rdf/terms#noassertion" />
			<spdx:licenseDeclared rdf:resource="http://spdx.org/rdf/terms#noassertion" />
			<spdx:licenseComments>Other versions available for a commercial license</spdx:licenseComments>
			<spdx:copyrightText rdf:resource="http://spdx.org/rdf/terms#noassertion" />
			<spdx:summary> Package for Testing </spdx:summary>
			<spdx:description> Some tags are taken from other spdx autogenerated files </spdx:description>
			<rdfs:comment>no comments</rdfs:comment>
			<spdx:externalRef>
				<spdx:ExternalRef>
					<spdx:referenceLocator>cpe:2.3:a:pivotal_software:spring_framework:4.1.0:*:*:*:*:*:*:*</spdx:referenceLocator>
					<spdx:referenceType>
						<spdx:ReferenceType rdf:about="http://spdx.org/rdf/references/cpe23Type"/>
					</spdx:referenceType>
					<spdx:referenceCategory rdf:resource="http://spdx.org/rdf/terms#referenceCategory_security"/>
				</spdx:ExternalRef>
			</spdx:externalRef>
			<spdx:hasFile rdf:resource="http://spdx.org/documents/spdx-toolsv2.1.7-SNAPSHOT#SPDXRef-129" />
			<spdx:relationship>
			    <spdx:Relationship>
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_describes" />
					<spdx:relatedSpdxElement rdf:resource="http://anupam-VirtualBox/repo/SPDX2_time-1.9#SPDXRef-upload2" />
				</spdx:Relationship>
			</spdx:relationship>
			<spdx:attributionText>attribution text</spdx:attributionText>
			<spdx:annotation>
				<spdx:Annotation>
					<spdx:annotationDate>2011-01-29T18:30:22Z</spdx:annotationDate>
					<rdfs:comment>Package level annotation</rdfs:comment>
					<spdx:annotator>Person: Package Commenter</spdx:annotator>
					<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_other"/>
				</spdx:Annotation>
			</spdx:annotation>
		</spdx:Package>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getPackageFromNode(node)
	if err != nil {
		t.Errorf("error parsing a valid package: %v", err)
	}
}

func Test_rdfParser2_3_setFileToPackage(t *testing.T) {
	var pkg *spdx.Package
	var file *spdx.File
	var parser *rdfParser2_3

	// TestCase 1: setting to a nil files attribute shouldn't panic.
	parser, _ = parserFromBodyContent(``)
	pkg = &spdx.Package{}
	file = &spdx.File{}
	parser.setFileToPackage(pkg, file)
	if len(pkg.Files) != 1 {
		t.Errorf("expected given package to have one file after setting, got %d", len(pkg.Files))
	}
	if parser.assocWithPackage[file.FileSPDXIdentifier] != true {
		t.Errorf("given file should've been associated with a package, assocWithPackage is false")
	}
}

func Test_rdfParser2_3_setPackageChecksum(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var pkg *spdx.Package
	var expectedChecksumValue string
	var err error

	// TestCase 1: invalid checksum algorithm
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha999"/>
		</spdx:Checksum>
	`)
	pkg = &spdx.Package{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setPackageChecksum(pkg, node)
	if err == nil {
		t.Error("expected an error due to invalid checksum node, got <nil>")
	}

	// TestCase 1: valid checksum algorithm which is invalid for package
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha2000"/>
		</spdx:Checksum>
	`)
	pkg = &spdx.Package{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setPackageChecksum(pkg, node)
	if err == nil {
		t.Error("expected an error due to invalid checksum for package, got <nil>")
	}

	// TestCase 2: valid checksum (sha1)
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1"/>
		</spdx:Checksum>
	`)
	pkg = &spdx.Package{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setPackageChecksum(pkg, node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedChecksumValue = "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"

	for _, checksum := range pkg.PackageChecksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != expectedChecksumValue {
				t.Errorf("expected %v, got: %v", expectedChecksumValue, checksum.Value)
			}
		}
	}

	// TestCase 3: valid checksum (sha256)
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha256"/>
		</spdx:Checksum>
	`)
	pkg = &spdx.Package{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setPackageChecksum(pkg, node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedChecksumValue = "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"
	for _, checksum := range pkg.PackageChecksums {
		switch checksum.Algorithm {
		case common.SHA256:
			if checksum.Value != expectedChecksumValue {
				t.Errorf("expected %v, got: %v", expectedChecksumValue, checksum.Value)
			}
		}
	}

	// TestCase 4: valid checksum (md5)
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_md5"/>
		</spdx:Checksum>
	`)
	pkg = &spdx.Package{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setPackageChecksum(pkg, node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedChecksumValue = "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"
	for _, checksum := range pkg.PackageChecksums {
		switch checksum.Algorithm {
		case common.MD5:
			if checksum.Value != expectedChecksumValue {
				t.Errorf("expected %v, got: %v", expectedChecksumValue, checksum.Value)
			}
		}
	}
}

func Test_setDocumentLocationFromURI(t *testing.T) {
	var pkg *spdx.Package
	var expectedDocumentLocation, gotDocumentLocation string
	var inputURI string
	var err error

	// TestCase 1: NOASSERTION
	inputURI = SPDX_NOASSERTION_SMALL
	pkg = &spdx.Package{}
	err = setDocumentLocationFromURI(pkg, inputURI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDocumentLocation = "NOASSERTION"
	gotDocumentLocation = pkg.PackageDownloadLocation
	if expectedDocumentLocation != gotDocumentLocation {
		t.Errorf("expected: %v, got: %v", expectedDocumentLocation, gotDocumentLocation)
	}

	// TestCase 2: NONE
	inputURI = SPDX_NONE_CAPS
	pkg = &spdx.Package{}
	err = setDocumentLocationFromURI(pkg, inputURI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDocumentLocation = "NONE"
	gotDocumentLocation = pkg.PackageDownloadLocation
	if expectedDocumentLocation != gotDocumentLocation {
		t.Errorf("expected: %v, got: %v", expectedDocumentLocation, gotDocumentLocation)
	}

	// TestCase 3: valid uri
	inputURI = "https://www.gnu.org/software/texinfo/"
	pkg = &spdx.Package{}
	err = setDocumentLocationFromURI(pkg, inputURI)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedDocumentLocation = "https://www.gnu.org/software/texinfo/"
	gotDocumentLocation = pkg.PackageDownloadLocation
	if expectedDocumentLocation != gotDocumentLocation {
		t.Errorf("expected: %v, got: %v", expectedDocumentLocation, gotDocumentLocation)
	}

	// TestCase 3: invalid uri
	inputURI = " "
	pkg = &spdx.Package{}
	err = setDocumentLocationFromURI(pkg, inputURI)
	if err == nil {
		t.Fatalf("expected an error due to invalid uri, got %v", err)
	}
}

func Test_setFilesAnalyzed(t *testing.T) {
	var pkg *spdx.Package
	var err error

	// TestCase 1: not a valid bool value:
	pkg = &spdx.Package{}
	err = setFilesAnalyzed(pkg, "no")
	if err == nil {
		t.Errorf("expected an error due to invalid bool input, got %v", err)
	}

	// TestCase 2: valid input
	pkg = &spdx.Package{}
	err = setFilesAnalyzed(pkg, "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pkg.IsFilesAnalyzedTagPresent {
		t.Errorf("should've set IsFilesAnalyzedTagPresent, got: %t", pkg.IsFilesAnalyzedTagPresent)
	}
	if !pkg.FilesAnalyzed {
		t.Errorf("expected: %t, got: %t", true, pkg.FilesAnalyzed)
	}
}
