// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func (parser *rdfParser2_3) parseSpdxDocumentNode(spdxDocNode *gordfParser.Node) (err error) {
	// shorthand for document's creation info.
	ci := parser.doc.CreationInfo

	// parse the document header information (SPDXID and document namespace)
	// the Subject.ID is of type baseURI#spdxID
	baseUri, offset, err := ExtractSubs(spdxDocNode.ID, "#")
	if err != nil {
		return err
	}
	parser.doc.DocumentNamespace = baseUri               // 2.5
	parser.doc.SPDXIdentifier = common.ElementID(offset) // 2.3

	// parse other associated triples.
	for _, subTriple := range parser.nodeToTriples(spdxDocNode) {
		objectValue := subTriple.Object.ID
		switch subTriple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_SPEC_VERSION: // 2.1: specVersion
			// cardinality: exactly 1
			parser.doc.SPDXVersion = objectValue
		case SPDX_DATA_LICENSE: // 2.2: dataLicense
			// cardinality: exactly 1
			dataLicense, err := parser.getAnyLicenseFromNode(subTriple.Object)
			if err != nil {
				return err
			}
			parser.doc.DataLicense = dataLicense.ToLicenseString()
		case SPDX_NAME: // 2.4: DocumentName
			// cardinality: exactly 1
			parser.doc.DocumentName = objectValue
		case SPDX_EXTERNAL_DOCUMENT_REF: // 2.6: externalDocumentReferences
			// cardinality: min 0
			var extRef spdx.ExternalDocumentRef
			extRef, err = parser.getExternalDocumentRefFromNode(subTriple.Object)
			if err != nil {
				return err
			}
			parser.doc.ExternalDocumentReferences = append(parser.doc.ExternalDocumentReferences, extRef)
		case SPDX_CREATION_INFO: // 2.7 - 2.10:
			// cardinality: exactly 1
			err = parser.parseCreationInfoFromNode(ci, subTriple.Object)
		case RDFS_COMMENT: // 2.11: Document Comment
			// cardinality: max 1
			parser.doc.DocumentComment = objectValue
		case SPDX_REVIEWED: // reviewed:
			// cardinality: min 0
			err = parser.setReviewFromNode(subTriple.Object)
		case SPDX_DESCRIBES_PACKAGE: // describes Package
			// cardinality: min 0
			var pkg *spdx.Package
			pkg, err = parser.getPackageFromNode(subTriple.Object)
			if err != nil {
				return err
			}
			parser.doc.Packages = append(parser.doc.Packages, pkg)
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
		default:
			return fmt.Errorf("invalid predicate while parsing SpdxDocument: %v", subTriple.Predicate)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (parser *rdfParser2_3) getExternalDocumentRefFromNode(node *gordfParser.Node) (edr spdx.ExternalDocumentRef, err error) {
	for _, triple := range parser.nodeToTriples(node) {
		switch triple.Predicate.ID {
		case SPDX_EXTERNAL_DOCUMENT_ID:
			// cardinality: exactly 1
			edr.DocumentRefID = common.DocumentID(triple.Object.ID)
		case SPDX_SPDX_DOCUMENT:
			// cardinality: exactly 1
			// assumption: "spdxDocument" property of an external document
			// reference is just a uri which doesn't follow a spdxDocument definition
			edr.URI = triple.Object.ID
		case SPDX_CHECKSUM:
			// cardinality: exactly 1
			alg, checksum, err := parser.getChecksumFromNode(triple.Object)
			if err != nil {
				return edr, err
			}
			edr.Checksum.Value = checksum
			edr.Checksum.Algorithm = alg
		case RDF_TYPE:
			continue
		default:
			return edr, fmt.Errorf("unknown predicate ID (%s) while parsing externalDocumentReference", triple.Predicate.ID)
		}
	}
	return edr, nil
}
