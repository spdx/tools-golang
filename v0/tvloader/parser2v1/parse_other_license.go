// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"

	"github.com/spdx/tools-golang/v0/spdx"
)

func (parser *tvParser2_1) parsePairFromOtherLicense2_1(tag string, value string) error {
	switch tag {
	// tag for creating new other license section
	case "LicenseID":
		parser.otherLic = &spdx.OtherLicense2_1{}
		parser.doc.OtherLicenses = append(parser.doc.OtherLicenses, parser.otherLic)
		parser.otherLic.LicenseIdentifier = value
	case "ExtractedText":
		parser.otherLic.ExtractedText = value
	case "LicenseName":
		parser.otherLic.LicenseName = value
	case "LicenseCrossReference":
		parser.otherLic.LicenseCrossReferences = append(parser.otherLic.LicenseCrossReferences, value)
	case "LicenseComment":
		parser.otherLic.LicenseComment = value
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_1
		return parser.parsePairFromReview2_1(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in OtherLicense section", tag)
	}

	return nil
}
