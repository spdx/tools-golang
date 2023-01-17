// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func (parser *rdfParser2_3) getExtractedLicensingInfoFromNode(node *gordfParser.Node) (lic ExtractedLicensingInfo, err error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	var restTriples []*gordfParser.Triple
	for _, triple := range associatedTriples {
		switch triple.Predicate.ID {
		case SPDX_EXTRACTED_TEXT:
			lic.extractedText = triple.Object.ID
		default:
			restTriples = append(restTriples, triple)
		}
	}
	lic.SimpleLicensingInfo, err = parser.getSimpleLicensingInfoFromTriples(restTriples)
	if err != nil {
		return lic, fmt.Errorf("error setting simple licensing information of extracted licensing info: %s", err)
	}
	return lic, nil
}

func (parser *rdfParser2_3) extractedLicenseToOtherLicense(extLicense ExtractedLicensingInfo) (othLic spdx.OtherLicense) {
	othLic.LicenseIdentifier = extLicense.licenseID
	othLic.ExtractedText = extLicense.extractedText
	othLic.LicenseComment = extLicense.comment
	othLic.LicenseCrossReferences = extLicense.seeAlso
	othLic.LicenseName = extLicense.name
	return othLic
}
