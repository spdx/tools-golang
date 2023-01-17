// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
)

func Test_rdfParser2_3_getExternalDocumentRefFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var err error

	// TestCase 1: invalid checksum
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalDocumentRef>
			<spdx:externalDocumentId>DocumentRef-spdx-tool-1.2</spdx:externalDocumentId>
			<spdx:checksum>
				<spdx:Checksum>
					<spdx:checksumValue>d6a770ba38583ed4bb4525bd96e50461655d2759</spdx:checksumValue>
					<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha999"/>
				</spdx:Checksum>
			</spdx:checksum>
			<spdx:spdxDocument rdf:resource="http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"/>
		</spdx:ExternalDocumentRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getExternalDocumentRefFromNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid checksum, found %v", err)
	}

	// TestCase 2: unknown predicate
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalDocumentRef>
			<spdx:unknownTag />
		</spdx:ExternalDocumentRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getExternalDocumentRefFromNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid predicate, found %v", err)
	}

	// TestCase 3: valid example
	parser, _ = parserFromBodyContent(`
		<spdx:ExternalDocumentRef>
			<spdx:externalDocumentId>DocumentRef-spdx-tool-1.2</spdx:externalDocumentId>
			<spdx:checksum>
				<spdx:Checksum>
					<spdx:checksumValue>d6a770ba38583ed4bb4525bd96e50461655d2759</spdx:checksumValue>
					<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha256"/>
				</spdx:Checksum>
			</spdx:checksum>
			<spdx:spdxDocument rdf:resource="http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"/>
		</spdx:ExternalDocumentRef>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getExternalDocumentRefFromNode(node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_rdfParser2_3_parseSpdxDocumentNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var err error

	// TestCase 1: invalid spdx id of the document
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301/Document"/>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid document id, got %v", err)
	}

	// TestCase 2: erroneous dataLicense
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:dataLicense rdf:resource="http://spdx.org/rdf/terms#Unknown" />
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid dataLicense, got %v", err)
	}

	// TestCase 3: invalid external document ref
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:externalDocumentRef>
				<spdx:ExternalDocumentRef>
					<spdx:externalDocumentId>DocumentRef-spdx-tool-1.2</spdx:externalDocumentId>
					<spdx:checksum>
						<spdx:Checksum>
							<spdx:checksumValue>d6a770ba38583ed4bb4525bd96e50461655d2759</spdx:checksumValue>
							<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha999"/>
						</spdx:Checksum>
					</spdx:checksum>
					<spdx:spdxDocument rdf:resource="http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"/>
				</spdx:ExternalDocumentRef>
        </spdx:externalDocumentRef>
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid externalDocumentRef, got %v", err)
	}

	// TestCase 4: invalid package
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:describesPackage>
				<spdx:Package rdf:about="http://www.spdx.org/spdxdocs/8f141b09-1138-4fc5-aecb-fc10d9ac1eed#SPDX-1"/>
			</spdx:describesPackage>
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid externalDocumentRef, got %v", err)
	}

	// TestCase 5: error in extractedLicensingInfo
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:hasExtractedLicensingInfo>
				<spdx:ExtractedLicensingInfo rdf:about="#LicenseRef-Freeware">
					<spdx:invalidTag />
					<spdx:licenseId>LicenseRef-Freeware</spdx:licenseId>
					<spdx:name>freeware</spdx:name>
					<spdx:extractedText>
<![CDATA[Software classified as freeware is licensed at no cost and is either fully functional for an unlimited time; or has only basic functions enabled with a fully functional version available commercially or as shareware.[8] In contrast to free software, the author usually restricts one or more rights of the user, including the rights to use, copy, distribute, modify and make derivative works of the software or extract the source code.[1][2][9][10] The software license may impose various additional restrictions on the type of use, e.g. only for personal use, private use, individual use, non-profit use, non-commercial use, academic use, educational use, use in charity or humanitarian organizations, non-military use, use by public authorities or various other combinations of these type of restrictions.[11] For instance, the license may be "free for private, non-commercial use". The software license may also impose various other restrictions, such as restricted use over a network, restricted use on a server, restricted use in a combination with some types of other software or with some hardware devices, prohibited distribution over the Internet other than linking to author's website, restricted distribution without author's consent, restricted number of copies, etc.]]>
					</spdx:extractedText>
				</spdx:ExtractedLicensingInfo>
			</spdx:hasExtractedLicensingInfo>
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid extractedLicensingInfo, got %v", err)
	}

	// TestCase 6: error in annotation
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:annotation>
				<spdx:Annotation>
					<spdx:unknownAttribute />
				</spdx:Annotation>
			</spdx:annotation>
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to invalid extractedLicensingInfo, got %v", err)
	}

	// TestCase 7: invalid predicate
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:unknownTag />
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err == nil {
		t.Errorf("expected an error due to unknown predicate, got %v", err)
	}

	// TestCase 7: everything valid
	parser, _ = parserFromBodyContent(`
		<spdx:SpdxDocument rdf:about="#SPDXRef-Document">
			<spdx:specVersion>SPDX-2.1</spdx:specVersion>
    		<spdx:dataLicense rdf:resource="http://spdx.org/licenses/CC0-1.0" />
			<spdx:name>/test/example</spdx:name>
			<spdx:externalDocumentRef>
				<spdx:ExternalDocumentRef>
					<spdx:externalDocumentId>DocumentRef-spdx-tool-1.2</spdx:externalDocumentId>
					<spdx:checksum>
						<spdx:Checksum>
							<spdx:checksumValue>d6a770ba38583ed4bb4525bd96e50461655d2759</spdx:checksumValue>
							<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1"/>
						</spdx:Checksum>
					</spdx:checksum>
					<spdx:spdxDocument rdf:resource="http://spdx.org/spdxdocs/spdx-tools-v1.2-3F2504E0-4F89-41D3-9A0C-0305E82C3301"/>
				</spdx:ExternalDocumentRef>
			</spdx:externalDocumentRef>	
			<spdx:creationInfo>
			    <spdx:CreationInfo>
					<spdx:licenseListVersion>2.6</spdx:licenseListVersion>
					<spdx:creator>Person: spdx (y)</spdx:creator>
					<spdx:creator>Organization: </spdx:creator>
					<spdx:creator>Tool: spdx2</spdx:creator>
					<spdx:created>2018-08-24T19:55:34Z</spdx:created>
			    </spdx:CreationInfo>
			</spdx:creationInfo>
    		<rdfs:comment>test</rdfs:comment>
			<spdx:reviewed>
				<spdx:Review>
					<rdfs:comment>Another example reviewer.</rdfs:comment>
					<spdx:reviewDate>2011-03-13T00:00:00Z</spdx:reviewDate>
					<spdx:reviewer>Person: Suzanne Reviewer</spdx:reviewer>
				</spdx:Review>
	        </spdx:reviewed>
			<spdx:describesPackage>
				<spdx:Package rdf:about="#SPDXRef-1"/>
			</spdx:describesPackage>
			<spdx:hasExtractedLicensingInfo rdf:resource="http://spdx.org/licenses/CC0-1.0"/>
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_containedBy"/>
					<spdx:relatedSpdxElement rdf:resource="http://spdx.org/documents/spdx-toolsv2.1.7-SNAPSHOT#SPDXRef-1"/>
					<rdfs:comment></rdfs:comment>
				</spdx:Relationship>
			</spdx:relationship>
			<spdx:annotation>
				<spdx:Annotation>
					<spdx:annotationDate>2011-01-29T18:30:22Z</spdx:annotationDate>
					<rdfs:comment>test annotation</rdfs:comment>
					<spdx:annotator>Person: Rishabh Bhatnagar</spdx:annotator>
					<spdx:annotationType rdf:resource="http://spdx.org/rdf/terms#annotationType_other"/>
				</spdx:Annotation>
			</spdx:annotation>
		</spdx:SpdxDocument>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseSpdxDocumentNode(node)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
