// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
	"strings"
)

func (parser *rdfParser2_2) getExternalLicensingInfoFromNode(node *gordfParser.Node) (*spdx.OtherLicense2_2, error) {
	lic := &spdx.OtherLicense2_2{}
	licensePrefix := "LicenseRef-"
	for _, triple := range parser.nodeToTriples[node.String()] {
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_LICENSE_ID:
			fragment := strings.TrimSpace(getLastPartOfURI(triple.Subject.ID))
			if !strings.HasPrefix(fragment, licensePrefix) {
				return nil, fmt.Errorf("license ID must be of type \"LicenseRef-[idstring]\"; found %s", fragment)
			}
			lic.LicenseIdentifier = strings.TrimSuffix(fragment, licensePrefix)
		case SPDX_EXTRACTED_TEXT:
			lic.ExtractedText = triple.Object.ID
		case SPDX_NAME:
			lic.LicenseName = triple.Object.ID
		case RDFS_SEE_ALSO:
			lic.LicenseCrossReferences = append(lic.LicenseCrossReferences, triple.Object.ID)
		case RDFS_COMMENT:
			lic.LicenseComment = triple.Object.ID
		default:
			return nil, fmt.Errorf("unknown predicate %v while parsing extractedLicensingInfo", triple.Predicate)
		}
	}
	return lic, nil
}

// parses the other license and appends it to the doc if no error is encountered.
func (parser *rdfParser2_2) parseOtherLicenseFromNode(node *gordfParser.Node) error {
	ol := &spdx.OtherLicense2_2{}
	ol.LicenseIdentifier = getLicenseStringFromURI(node.ID) // 6.1
	for _, triple := range parser.nodeToTriples[node.String()] {
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_EXTRACTED_TEXT: // 6.2
			ol.ExtractedText = triple.Object.ID
		case SPDX_NAME: // 6.3
			ol.LicenseName = triple.Object.ID
		case RDFS_SEE_ALSO: // 6.4
			ol.LicenseCrossReferences = append(ol.LicenseCrossReferences, triple.Object.ID)
		case RDFS_COMMENT: // 6.5
			ol.LicenseComment = triple.Object.ID
		case SPDX_LICENSE_ID:
			// override licenseId from the rdf:about tag.
			ol.LicenseIdentifier = getLicenseStringFromURI(triple.Object.ID)
		default:
			return fmt.Errorf("unknown predicate (%s) while parsing other license", triple.Predicate.ID)
		}
	}

	parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, ol)
	return nil
}
