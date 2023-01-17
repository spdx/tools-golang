// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"bufio"
	"strings"
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	rdfloader2 "github.com/spdx/gordf/rdfloader/xmlreader"
	gordfWriter "github.com/spdx/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// content is the tags within the rdf:RDF tag
// pads the content with the enclosing rdf:RDF tag
func wrapIntoTemplate(content string) string {
	header := `<rdf:RDF
        xmlns:spdx="http://spdx.org/rdf/terms#"
        xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
        xmlns="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#"
        xmlns:doap="http://usefulinc.com/ns/doap#"
        xmlns:j.0="http://www.w3.org/2009/pointers#"
        xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#">`
	footer := `</rdf:RDF>`
	return header + content + footer
}

func parserFromBodyContent(content string) (*rdfParser2_2, error) {
	rdfContent := wrapIntoTemplate(content)
	xmlreader := rdfloader2.XMLReaderFromFileObject(bufio.NewReader(strings.NewReader(rdfContent)))
	rootBlock, err := xmlreader.Read()
	if err != nil {
		return nil, err
	}
	parser := gordfParser.New()
	err = parser.Parse(rootBlock)
	if err != nil {
		return nil, err
	}
	nodeToTriples := gordfWriter.GetNodeToTriples(parser.Triples)
	rdfParser := NewParser2_2(parser, nodeToTriples)
	return rdfParser, err
}

func Test_rdfParser2_2_getArtifactFromNode(t *testing.T) {
	// TestCase 1: artifactOf without project URI
	rdfParser, err := parserFromBodyContent(
		`<spdx:File>
			<spdx:artifactOf>
				<doap:Project>
					<doap:homepage>http://www.openjena.org/</doap:homepage>
					<doap:name>Jena</doap:name>
				</doap:Project>
			</spdx:artifactOf>
		</spdx:File>`)
	if err != nil {
		t.Errorf("unexpected error while parsing a valid example: %v", err)
	}
	artifactOfNode := gordfWriter.FilterTriples(rdfParser.gordfParserObj.Triples, nil, &SPDX_ARTIFACT_OF, nil)[0].Object
	artifact, err := rdfParser.getArtifactFromNode(artifactOfNode)
	if err != nil {
		t.Errorf("error parsing a valid artifactOf node: %v", err)
	}
	if artifact.Name != "Jena" {
		t.Errorf("expected name of artifact: %s, found: %s", "Jena", artifact.Name)
	}
	expectedHomePage := "http://www.openjena.org/"
	if artifact.HomePage != expectedHomePage {
		t.Errorf("wrong artifact homepage. Expected: %s, found: %s", expectedHomePage, artifact.HomePage)
	}
	if artifact.URI != "" {
		t.Errorf("wrong artifact URI. Expected: %s, found: %s", "", artifact.URI)
	}

	// TestCase 2: artifactOf with a Project URI
	rdfParser, err = parserFromBodyContent(
		`<spdx:File>
			<spdx:artifactOf>
				<doap:Project rdf:about="http://subversion.apache.org/doap.rdf">
					<doap:homepage>http://www.openjena.org/</doap:homepage>
					<doap:name>Jena</doap:name>
				</doap:Project>
			</spdx:artifactOf>
		</spdx:File>`)
	if err != nil {
		t.Errorf("unexpected error while parsing a valid example: %v", err)
	}
	artifactOfNode = gordfWriter.FilterTriples(rdfParser.gordfParserObj.Triples, nil, &SPDX_ARTIFACT_OF, nil)[0].Object
	artifact, err = rdfParser.getArtifactFromNode(artifactOfNode)
	if err != nil {
		t.Errorf("error parsing a valid artifactOf node: %v", err)
	}
	expectedURI := "http://subversion.apache.org/doap.rdf"
	if artifact.URI != expectedURI {
		t.Errorf("wrong artifact URI. Expected: %s, found: %s", expectedURI, artifact.URI)
	}

	// TestCase 3: artifactOf with unknown predicate
	rdfParser, err = parserFromBodyContent(
		`<spdx:File>
			<spdx:artifactOf>
				<doap:Project rdf:about="http://subversion.apache.org/doap.rdf">
					<doap:homepage>http://www.openjena.org/</doap:homepage>
					<doap:name>Jena</doap:name>
					<doap:invalidTag rdf:ID="invalid"/>
				</doap:Project>
			</spdx:artifactOf>
		</spdx:File>`)
	if err != nil {
		t.Errorf("unexpected error while parsing a valid example: %v", err)
	}
	artifactOfNode = gordfWriter.FilterTriples(rdfParser.gordfParserObj.Triples, nil, &SPDX_ARTIFACT_OF, nil)[0].Object
	_, err = rdfParser.getArtifactFromNode(artifactOfNode)
	if err == nil {
		t.Errorf("must've raised an error for an invalid predicate")
	}
}

func Test_rdfParser2_2_getFileTypeFromUri(t *testing.T) {
	rdfParser, _ := parserFromBodyContent(``)

	// TestCase 1: Valid fileType URI:
	fileTypeURI := "http://spdx.org/rdf/terms#fileType_source"
	fileType, err := rdfParser.getFileTypeFromUri(fileTypeURI)
	if err != nil {
		t.Errorf("error in a valid example: %v", err)
	}
	if fileType != "source" {
		t.Errorf("wrong fileType. expected: %s, found: %s", "source", fileType)
	}

	// TestCase 2: Invalid fileType URI format.
	fileTypeURI = "http://spdx.org/rdf/terms#source"
	fileType, err = rdfParser.getFileTypeFromUri(fileTypeURI)
	if err == nil {
		t.Error("should've raised an error for invalid fileType")
	}
}

func Test_rdfParser2_2_setUnpackagedFiles(t *testing.T) {
	// unpackaged files are the files which are not associated with any package
	// file associated with a package sets parser.assocWithPackage[fileID] to true.
	rdfParser, _ := parserFromBodyContent(``)
	file1 := &v2_2.File{FileSPDXIdentifier: common.ElementID("file1")}
	file2 := &v2_2.File{FileSPDXIdentifier: common.ElementID("file2")}
	file3 := &v2_2.File{FileSPDXIdentifier: common.ElementID("file3")}

	// setting files to the document as if it were to be set when it was parsed using triples.
	rdfParser.files[file1.FileSPDXIdentifier] = file1
	rdfParser.files[file2.FileSPDXIdentifier] = file2
	rdfParser.files[file3.FileSPDXIdentifier] = file3

	// assuming file1 is associated with a package
	rdfParser.assocWithPackage[file1.FileSPDXIdentifier] = true

	rdfParser.setUnpackagedFiles()

	// after setting unpackaged files, parser.doc.Files must've file2 and file3
	if n := len(rdfParser.doc.Files); n != 2 {
		t.Errorf("unpackage files should've had 2 files, found %d files", n)
	}

	// checking if the unpackagedFiles contain only file2 & file3.
	for _, file := range rdfParser.doc.Files {
		switch string(file.FileSPDXIdentifier) {
		case "file2", "file3":
			continue
		default:
			t.Errorf("unexpected file with id %s found in unpackaged files", file.FileSPDXIdentifier)
		}
	}
}

func Test_setFileIdentifier(t *testing.T) {
	file := &v2_2.File{}

	// TestCase 1: valid example
	err := setFileIdentifier("http://spdx.org/documents/spdx-toolsv2.1.7-SNAPSHOT#SPDXRef-129", file)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if file.FileSPDXIdentifier != "129" {
		t.Errorf("expected %s, found: %s", "129", file.FileSPDXIdentifier)
	}

	// TestCase 2: invalid example
	err = setFileIdentifier("http://spdx.org/documents/spdx-toolsv2.1.7-SNAPSHOT#129", file)
	if err == nil {
		t.Errorf("should've raised an error for an invalid example")
	}
}

func Test_rdfParser2_2_setFileChecksumFromNode(t *testing.T) {
	// TestCase 1: md5 checksum
	parser, _ := parserFromBodyContent(` 
		<spdx:Checksum>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_md5" />
		    <spdx:checksumValue>d2356e0fe1c0b85285d83c6b2ad51b5f</spdx:checksumValue>
		</spdx:Checksum>
    `)
	checksumNode := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_CHECKSUM_CAPITALIZED)[0].Subject
	file := &v2_2.File{}
	err := parser.setFileChecksumFromNode(file, checksumNode)
	if err != nil {
		t.Errorf("error parsing a valid checksum node")
	}
	checksumValue := "d2356e0fe1c0b85285d83c6b2ad51b5f"
	for _, checksum := range file.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != "" {
				t.Errorf("incorrectly set sha1, should've been empty")
			}
		case common.SHA256:
			if checksum.Value != "" {
				t.Errorf("incorrectly set sha256, should've been empty")
			}
		case common.MD5:
			if checksum.Value != checksumValue {
				t.Errorf("wrong checksum value for md5. Expected: %s, found: %s", checksumValue, checksum.Value)
			}
		}
	}

	// TestCase 2: valid sha1 checksum
	parser, _ = parserFromBodyContent(` 
		<spdx:Checksum>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1" />
		    <spdx:checksumValue>d2356e0fe1c0b85285d83c6b2ad51b5f</spdx:checksumValue>
		</spdx:Checksum>
    `)
	checksumNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_CHECKSUM_CAPITALIZED)[0].Subject
	file = &v2_2.File{}
	err = parser.setFileChecksumFromNode(file, checksumNode)
	if err != nil {
		t.Errorf("error parsing a valid checksum node")
	}
	for _, checksum := range file.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != checksumValue {
				t.Errorf("wrong checksum value for sha1. Expected: %s, found: %s", checksumValue, checksum.Value)
			}
		case common.SHA256:
			if checksum.Value != "" {
				t.Errorf("incorrectly set sha256, should've been empty")
			}
		case common.MD5:
			if checksum.Value != checksumValue {
				t.Errorf("incorrectly set md5, should've been empty")
			}
		}
	}

	// TestCase 3: valid sha256 checksum
	parser, _ = parserFromBodyContent(` 
		<spdx:Checksum>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha256" />
		    <spdx:checksumValue>d2356e0fe1c0b85285d83c6b2ad51b5f</spdx:checksumValue>
		</spdx:Checksum>
    `)
	checksumNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_CHECKSUM_CAPITALIZED)[0].Subject
	file = &v2_2.File{}
	err = parser.setFileChecksumFromNode(file, checksumNode)
	if err != nil {
		t.Errorf("error parsing a valid checksum node")
	}
	for _, checksum := range file.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != checksumValue {
				t.Errorf("incorrectly set sha1, should've been empty")
			}
		case common.SHA256:
			if checksum.Value != checksumValue {
				t.Errorf("wrong checksum value for sha256. Expected: %s, found: %s", checksumValue, checksum.Value)
			}
		case common.MD5:
			if checksum.Value != checksumValue {
				t.Errorf("incorrectly set md5, should've been empty")
			}
		}
	}

	// TestCase 4: checksum node without one of the mandatory attributes
	parser, _ = parserFromBodyContent(` 
		<spdx:Checksum>
		    <spdx:checksumValue>d2356e0fe1c0b85285d83c6b2ad51b5f</spdx:checksumValue>
		</spdx:Checksum>
    `)
	checksumNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_CHECKSUM_CAPITALIZED)[0].Subject
	file = &v2_2.File{}
	err = parser.setFileChecksumFromNode(file, checksumNode)
	if err == nil {
		t.Errorf("should've raised an error parsing an invalid checksum node")
	}

	// TestCase 5: invalid checksum algorithm
	parser, _ = parserFromBodyContent(` 
		<spdx:Checksum>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_md43" />
		    <spdx:checksumValue>d2356e0fe1c0b85285d83c6b2ad51b5f</spdx:checksumValue>
		</spdx:Checksum>
    `)
	checksumNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_CHECKSUM_CAPITALIZED)[0].Subject
	file = &v2_2.File{}
	err = parser.setFileChecksumFromNode(file, checksumNode)
	if err == nil {
		t.Errorf("should've raised an error parsing an invalid checksum node")
	}

	// TestCase 6: valid checksum algorithm which is invalid for file (like md4, md6, sha384, etc.)
	parser, _ = parserFromBodyContent(` 
		<spdx:Checksum>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha2000" />
		    <spdx:checksumValue>d2356e0fe1c0b85285d83c6b2ad51b5f</spdx:checksumValue>
		</spdx:Checksum>
    `)
	checksumNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_CHECKSUM_CAPITALIZED)[0].Subject
	file = &v2_2.File{}
	err = parser.setFileChecksumFromNode(file, checksumNode)
	if err == nil {
		t.Errorf("should've raised an error parsing an invalid checksum algorithm for a file")
	}
}

func Test_rdfParser2_2_getFileFromNode(t *testing.T) {
	// TestCase 1: file with invalid id
	parser, _ := parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gzspdx.rdf#item177"/>
	`)
	fileNode := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err := parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid file ID")
	}

	// TestCase 2: invalid fileType
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:fileType rdf:resource="http://spdx.org/rdf/terms#source"/>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid fileType")
	}

	// TestCase 3: invalid file checksum
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:checksum>
				<spdx:Checksum>
					<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha2000" />
					<spdx:checksumValue>0a3a0e1ab72b7c132f5021c538a7a3ea6d539bcd</spdx:checksumValue>
				</spdx:Checksum>
			</spdx:checksum>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid checksum")
	}

	// TestCase 4: invalid license concluded
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:licenseConcluded rdf:resource="http://spdx.org/rdf/terms#invalid_license" />
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid license Concluded")
	}

	// TestCase 5: invalid artifactOf attribute
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:artifactOf>
				<doap:Project>
					<doap:unknown_tag />
					<doap:name>Jena</doap:name>
				</doap:Project>
			</spdx:artifactOf>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid artifactOf predicate")
	}

	// TestCase 6: invalid file dependency
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:fileDependency rdf:resource="http://spdx.org/spdxdocs/spdx-example#CommonsLangSrc"/>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid fileDependency")
	}

	// TestCase 7: invalid annotation with unknown predicate
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:annotation>
				<spdx:Annotation>
					<spdx:unknownAttribute />
				</spdx:Annotation>
			</spdx:annotation>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid annotation predicate")
	}

	// TestCase 8: invalid relationship
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#dynamicLink"/>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Errorf("should've raised an error stating invalid relationship Type")
	}

	// TestCase 8: unknown predicate
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:unknown />
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Error("should've raised an error stating invalid predicate for a file")
	}

	// TestCase 9: invalid licenseInfoInFile.
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:licenseInfoInFile rdf:resource="http://spdx.org/licenses/DC0-1.0" />
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	_, err = parser.getFileFromNode(fileNode)
	if err == nil {
		t.Error("should've raised an error stating invalid licenseInfoInFile for a file")
	}

	// TestCase 10: Splitting of File definition into parents of different tags mustn't create new file objects.
	fileDefinitions := []string{
		`<spdx:Package rdf:about="#SPDXRef-Package1">
			<spdx:hasFile>
				<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
					<spdx:fileName>time-1.9/ChangeLog</spdx:fileName>
					<spdx:fileType rdf:resource="http://spdx.org/rdf/terms#fileType_source"/>
				</spdx:File>
			</spdx:hasFile>
		</spdx:Package>`,
		`<spdx:Package rdf:about="#SPDXRef-Package2">
			<spdx:hasFile>
				<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
					<spdx:licenseConcluded rdf:resource="http://spdx.org/rdf/terms#noassertion" />
					<spdx:licenseInfoInFile rdf:resource="http://spdx.org/rdf/terms#NOASSERTION" />
				</spdx:File>
			</spdx:hasFile>
		</spdx:Package>`,
	}
	parser, _ = parserFromBodyContent(strings.Join(fileDefinitions, ""))

	var file *v2_2.File
	packageTypeTriples := gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_PACKAGE)
	for _, typeTriple := range packageTypeTriples {
		pkg, err := parser.getPackageFromNode(typeTriple.Subject)
		if err != nil {
			t.Errorf("unexpected error parsing a valid package: %v", err)
		}
		if n := len(pkg.Files); n != 1 {
			t.Errorf("expected package to contain exactly 1 file. Found %d files", n)
		}
		for _, file = range pkg.Files {
		}
	}

	// checking if all the attributes that spanned over a several tags are set in the same variable.
	expectedFileName := "time-1.9/ChangeLog"
	if file.FileName != expectedFileName {
		t.Errorf("expected %s, found %s", expectedFileName, file.FileName)
	}
	expectedLicenseConcluded := "NOASSERTION"
	if file.LicenseConcluded != expectedLicenseConcluded {
		t.Errorf("expected %s, found %s", expectedLicenseConcluded, file.LicenseConcluded)
	}
	expectedFileType := "source"
	if file.FileTypes[0] != expectedFileType {
		t.Errorf("expected %s, found %s", expectedFileType, file.FileTypes)
	}
	expectedLicenseInfoInFile := "NOASSERTION"
	if file.LicenseInfoInFiles[0] != expectedLicenseInfoInFile {
		t.Errorf("expected %s, found %s", expectedLicenseInfoInFile, file.LicenseInfoInFiles[0])
	}

	// TestCase 12: checking if recursive dependencies are resolved.
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="#SPDXRef-ParentFile">
			<spdx:fileType rdf:resource="http://spdx.org/rdf/terms#fileType_source"/>
			<spdx:fileDependency>
				<spdx:File rdf:about="#SPDXRef-ChildFile">
					<spdx:fileDependency>
						<spdx:File rdf:about="#SPDXRef-ParentFile">
							<spdx:fileName>ParentFile</spdx:fileName>
						</spdx:File>
					</spdx:fileDependency>
				</spdx:File>
			</spdx:fileDependency>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	file, err = parser.getFileFromNode(fileNode)

	// TestCase 11: all valid attribute and it's values.
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://anupam-VirtualBox/repo/SPDX2_time-1.9.tar.gz_1535120734-spdx.rdf#SPDXRef-item177">
			<spdx:fileName>time-1.9/ChangeLog</spdx:fileName>
			<spdx:name/>
			<spdx:fileType rdf:resource="http://spdx.org/rdf/terms#fileType_source"/>
			<spdx:checksum>
				<spdx:Checksum>
					<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1" />
					<spdx:checksumValue>0a3a0e1ab72b7c132f5021c538a7a3ea6d539bcd</spdx:checksumValue>
				</spdx:Checksum>
			</spdx:checksum>
			<spdx:licenseConcluded rdf:resource="http://spdx.org/rdf/terms#noassertion" />
			<spdx:licenseInfoInFile rdf:resource="http://spdx.org/rdf/terms#NOASSERTION" />
			<spdx:licenseComments>no comments</spdx:licenseComments>
			<spdx:copyrightText>from spdx file</spdx:copyrightText>
			<spdx:artifactOf>
				<doap:Project>
					<doap:homepage>http://www.openjena.org/</doap:homepage>
					<doap:name>Jena</doap:name>
				</doap:Project>
			</spdx:artifactOf>
			<rdfs:comment>no comments</rdfs:comment>
			<spdx:noticeText rdf:resource="http://spdx.org/rdf/terms#noassertion"/>
			<spdx:fileContributor>Some Organization</spdx:fileContributor>
			<spdx:fileDependency rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-CommonsLangSrc"/>
			<spdx:attributionText>attribution text</spdx:attributionText>
			<spdx:annotation>
				<spdx:Annotation>
					<spdx:annotationDate>2011-01-29T18:30:22Z</spdx:annotationDate>
					<rdfs:comment>File level annotation copied from a spdx document</rdfs:comment>
					<spdx:annotator>Person: File Commenter</spdx:annotator>
					<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_other"/>
				</spdx:Annotation>
			</spdx:annotation>
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_contains"/>
					<spdx:relatedSpdxElement rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Package"/>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	fileNode = gordfWriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0].Subject
	file, err = parser.getFileFromNode(fileNode)
	if err != nil {
		t.Errorf("unexpected error parsing a valid file: %v", err)
	}

	// checking each and every attribute of the obtained file.

	expectedFileName = "time-1.9/ChangeLog"
	if file.FileName != expectedFileName {
		t.Errorf("expected %s, found %s", expectedFileName, file.FileName)
	}

	if len(file.FileTypes) != 1 {
		t.Errorf("given file should have 1 fileType attribute. found %d", len(file.FileTypes))
	}
	expectedFileType = "source"
	if file.FileTypes[0] != expectedFileType {
		t.Errorf("expected %s, found %s", expectedFileType, file.FileTypes)
	}

	expectedChecksum := "0a3a0e1ab72b7c132f5021c538a7a3ea6d539bcd"

	for _, checksum := range file.Checksums {
		switch checksum.Algorithm {
		case common.SHA1:
			if checksum.Value != expectedChecksum {
				t.Errorf("expected %s, found %s", expectedChecksum, checksum.Value)
			}
		}
	}

	expectedLicenseConcluded = "NOASSERTION"
	if file.LicenseConcluded != expectedLicenseConcluded {
		t.Errorf("expected %s, found %s", expectedLicenseConcluded, file.LicenseConcluded)
	}

	if len(file.LicenseInfoInFiles) != 1 {
		t.Errorf("given file should have 1 licenseInfoInFile attribute. found %d", len(file.LicenseInfoInFiles))
	}
	expectedLicenseInfoInFile = "NOASSERTION"
	if file.LicenseInfoInFiles[0] != expectedLicenseInfoInFile {
		t.Errorf("expected %s, found %s", expectedLicenseInfoInFile, file.LicenseInfoInFiles[0])
	}

	expectedLicenseComments := "no comments"
	if file.LicenseComments != expectedLicenseComments {
		t.Errorf("expected %s, found %s", expectedLicenseComments, file.LicenseComments)
	}

	expectedCopyrightText := "from spdx file"
	if file.FileCopyrightText != expectedCopyrightText {
		t.Errorf("expected %s, found %s", expectedCopyrightText, file.FileCopyrightText)
	}

	if n := len(file.ArtifactOfProjects); n != 1 {
		t.Errorf("given file should have 1 artifactOfProjects attribute. found %d", n)
	}
	artifactOf := file.ArtifactOfProjects[0]
	expectedHomePage := "http://www.openjena.org/"
	if artifactOf.HomePage != expectedHomePage {
		t.Errorf("expected %s, found %s", expectedHomePage, artifactOf.HomePage)
	}
	if artifactOf.Name != "Jena" {
		t.Errorf("expected %s, found %s", "Jena", artifactOf.Name)
	}
	if artifactOf.URI != "" {
		t.Errorf("expected artifactOf uri to be empty, found %s", artifactOf.URI)
	}

	expectedFileComment := "no comments"
	if file.FileComment != expectedFileComment {
		t.Errorf("expected %s, found %s", expectedFileName, file.FileComment)
	}

	expectedNoticeText := "NOASSERTION"
	if file.FileNotice != expectedNoticeText {
		t.Errorf("expected %s, found %s", expectedNoticeText, file.FileNotice)
	}

	if n := len(file.FileContributors); n != 1 {
		t.Errorf("given file should have 1 fileContributor. found %d", n)
	}
	expectedFileContributor := "Some Organization"
	if file.FileContributors[0] != expectedFileContributor {
		t.Errorf("expected %s, found %s", expectedFileContributor, file.FileContributors)
	}

	if n := len(file.FileDependencies); n != 1 {
		t.Errorf("given file should have 1 fileDependencies. found %d", n)
	}
	expectedFileDependency := "CommonsLangSrc"
	if file.FileDependencies[0] != expectedFileDependency {
		t.Errorf("expected %s, found %s", expectedFileDependency, file.FileDependencies[0])
	}

	if n := len(file.FileAttributionTexts); n != 1 {
		t.Errorf("given file should have 1 attributionText. found %d", n)
	}
	expectedAttributionText := "attribution text"
	if file.FileAttributionTexts[0] != expectedAttributionText {
		t.Errorf("expected %s, found %s", expectedAttributionText, file.FileAttributionTexts[0])
	}

	if n := len(parser.doc.Annotations); n != 1 {
		t.Errorf("doc should've had 1 annotation. found %d", n)
	}
	ann := parser.doc.Annotations[0]
	expectedAnnDate := "2011-01-29T18:30:22Z"
	if ann.AnnotationDate != expectedAnnDate {
		t.Errorf("expected %s, found %s", expectedAnnDate, ann.AnnotationDate)
	}
	expectedAnnComment := "File level annotation copied from a spdx document"
	if ann.AnnotationComment != expectedAnnComment {
		t.Errorf("expected %s, found %s", expectedAnnComment, ann.AnnotationComment)
	}
	expectedAnnotationType := "OTHER"
	if ann.AnnotationType != expectedAnnotationType {
		t.Errorf("expected %s, found %s", expectedAnnotationType, ann.AnnotationType)
	}
	expectedAnnotator := "File Commenter"
	if ann.Annotator.Annotator != expectedAnnotator {
		t.Errorf("expected %s, found %s", expectedAnnotator, ann.Annotator)
	}
	expectedAnnotatorType := "Person"
	if ann.AnnotationType != expectedAnnotationType {
		t.Errorf("expected %s, found %s", expectedAnnotatorType, ann.Annotator.AnnotatorType)
	}

	if n := len(parser.doc.Relationships); n != 1 {
		t.Errorf("doc should've had 1 relation. found %d", n)
	}
	reln := parser.doc.Relationships[0]
	expectedRefAEID := "item177"
	if reln.RefA.DocumentRefID != "" {
		t.Errorf("expected refA.DocumentRefID to be empty, found %s", reln.RefA.DocumentRefID)
	}
	if string(reln.RefA.ElementRefID) != expectedRefAEID {
		t.Errorf("expected %s, found %s", expectedRefAEID, reln.RefA.ElementRefID)
	}
	expectedRefBEID := "Package"
	if reln.RefB.DocumentRefID != "" {
		t.Errorf("expected refB.DocumentRefID to be empty, found %s", reln.RefB.DocumentRefID)
	}
	if string(reln.RefB.ElementRefID) != expectedRefBEID {
		t.Errorf("expected %s, found %s", expectedRefBEID, reln.RefB.ElementRefID)
	}
	expectedRelationType := "contains"
	if reln.Relationship != expectedRelationType {
		t.Errorf("expected %s, found %s", expectedRefBEID, reln.RefB.ElementRefID)
	}
	if reln.RelationshipComment != "" {
		t.Errorf("expected relationship comment to be empty, found %s", reln.RelationshipComment)
	}
}

func Test_getNoticeTextFromNode(t *testing.T) {
	// TestCase 1: SPDX_NOASSERTION_SMALL must return NOASSERTION
	output := getNoticeTextFromNode(&gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       SPDX_NOASSERTION_SMALL,
	})
	if strings.ToUpper(output) != "NOASSERTION" {
		t.Errorf("expected NOASSERTION, found %s", strings.ToUpper(output))
	}

	// TestCase 2: SPDX_NOASSERTION_CAPS must return NOASSERTION
	output = getNoticeTextFromNode(&gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       SPDX_NOASSERTION_CAPS,
	})
	if strings.ToUpper(output) != "NOASSERTION" {
		t.Errorf("expected NOASSERTION, found %s", strings.ToUpper(output))
	}

	// TestCase 3: not a NOASSERTION must return the field verbatim
	// TestCase 1: SPDX_NOASSERTION_SMALL must return NOASSERTION
	output = getNoticeTextFromNode(&gordfParser.Node{
		NodeType: gordfParser.IRI,
		ID:       "text",
	})
	if output != "text" {
		t.Errorf("expected text, found %s", output)
	}
}
