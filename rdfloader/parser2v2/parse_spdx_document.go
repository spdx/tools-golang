// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
)

func (parser *rdfParser2_2) parseSpdxDocumentNode(spdxDocNode *gordfParser.Node) (err error) {
	// create a new creation info
	ci := parser.doc.CreationInfo

	// parse the document header information (SPDXID and document namespace)
	// the Subject.ID is of type baseURI#spdxID
	baseUri, offset, err := ExtractSubs(spdxDocNode.ID, "#")
	if err != nil {
		return err
	}
	ci.DocumentNamespace = baseUri             // 2.5
	ci.SPDXIdentifier = spdx.ElementID(offset) // 2.3

	// parse other associated triples.
	for _, subTriple := range parser.nodeToTriples(spdxDocNode) {
		objectValue := subTriple.Object.ID
		switch subTriple.Predicate.ID {
		case SPDX_SPEC_VERSION: // 2.1: specVersion
			// cardinality: exactly 1
			ci.SPDXVersion = objectValue
		case SPDX_DATA_LICENSE: // 2.2: dataLicense
			// cardinality: exactly 1
			dataLicense, err := parser.getAnyLicenseFromNode(subTriple.Object)
			if err != nil {
				return err
			}
			ci.DataLicense = dataLicense.ToLicenseString()
		case SPDX_NAME: // 2.4: DocumentName
			// cardinality: exactly 1
			ci.DocumentName = objectValue
		case SPDX_EXTERNAL_DOCUMENT_REF: // 2.6: externalDocumentReferences
			// cardinality: min 0
			var extRef string
			extRef, err = parser.getExternalDocumentRefFromTriples(parser.nodeToTriples(subTriple.Object))
			ci.ExternalDocumentReferences = append(ci.ExternalDocumentReferences, extRef)
		case SPDX_CREATION_INFO: // 2.7 - 2.10:
			// cardinality: exactly 1
			err = parser.parseCreationInfoFromNode(ci, subTriple.Object)
		case RDFS_COMMENT: // 2.11: Document Comment
			// cardinality: max 1
			ci.DocumentComment = objectValue
		case SPDX_REVIEWED: // reviewed:
			// cardinality: min 0
			err = parser.setReviewFromNode(subTriple.Object)
		case SPDX_DESCRIBES_PACKAGE: // describes Package
			// cardinality: min 0
			var pkg *spdx.Package2_2
			pkg, err = parser.getPackageFromNode(subTriple.Object)
			if err != nil {
				return err
			}
			parser.doc.Packages[pkg.PackageSPDXIdentifier] = pkg
		case SPDX_HAS_EXTRACTED_LICENSING_INFO: // hasExtractedLicensingInfo
			// cardinality: min 0
			extractedLicensingInfo, err := parser.getExtractedLicensingInfoFromNode(subTriple.Object)
			if err != nil {
				return fmt.Errorf("error setting extractedLicensingInfo in spdxDocument: %v", err)
			}
			othLicense := parser.extractedLicenseToOtherLicense(extractedLicensingInfo)
			parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, &othLicense)
		case SPDX_RELATIONSHIP: // relationship
			// cardinality: min 0
			err = parser.parseRelationship(subTriple)
		case SPDX_ANNOTATION: // annotations
			// cardinality: min 0
			err = parser.parseAnnotationFromNode(subTriple.Object)
		}
		if err != nil {
			return err
		}
	}

	// control reaches here iff no error is encountered
	// set the ci if no error is encountered while parsing triples.
	parser.doc.CreationInfo = ci
	return nil
}

func (parser *rdfParser2_2) getExternalDocumentRefFromTriples(triples []*gordfParser.Triple) (string, error) {
	var docID, checksumValue, checksumAlgorithm, spdxDocument string
	var err error
	for _, triple := range triples {
		switch triple.Predicate.ID {
		case SPDX_EXTERNAL_DOCUMENT_ID:
			// cardinality: exactly 1
			docID = triple.Object.ID
		case SPDX_SPDX_DOCUMENT:
			// cardinality: exactly 1
			// assumption: "spdxDocument" property of an external document
			// reference is just a uri which doesn't follow a spdxDocument definition
			spdxDocument = triple.Object.ID
		case SPDX_CHECKSUM:
			// cardinality: exactly 1
			checksumAlgorithm, checksumValue, err = parser.getChecksumFromNode(triple.Object)
			if err != nil {
				return "", err
			}
		case RDF_TYPE:
			continue
		default:
			return "", fmt.Errorf("unknown predicate ID (%s) while parsing externalDocumentReference", triple.Predicate.ID)
		}
	}
	// transform the variables into string form (same as that of tag-value).
	return fmt.Sprintf("%s %s %s: %s", docID, spdxDocument, checksumAlgorithm, checksumValue), nil
}
