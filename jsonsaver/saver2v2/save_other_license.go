// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func renderOtherLicenses2_2(otherlicenses []*spdx.OtherLicense2_2, buf *bytes.Buffer) error {

	var licenses []interface{}
	for _, v := range otherlicenses {
		lic := make(map[string]interface{})
		lic["licenseId"] = v.LicenseIdentifier
		lic["extractedText"] = v.ExtractedText
		if v.LicenseComment != "" {
			lic["comment"] = v.LicenseComment
		}
		if v.LicenseName != "" {
			lic["name"] = v.LicenseName
		}
		if v.LicenseCrossReferences != nil {
			lic["seeAlsos"] = v.LicenseCrossReferences
		}
		licenses = append(licenses, lic)
	}
	licensesjson, _ := json.Marshal(licenses)
	fmt.Fprintf(buf, "\"%s\": %s ,", "hasExtractedLicensingInfos", licensesjson)

	return nil
}
