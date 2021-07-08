// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func renderCreationInfo2_2(ci *spdx.CreationInfo2_2, buf *bytes.Buffer) error {
	if ci.SPDXIdentifier != "" {
		fmt.Fprintf(buf, "\"%s\": \"%s\",", "SPDXID", spdx.RenderElementID(ci.SPDXIdentifier))
	}
	if ci.SPDXVersion != "" {
		fmt.Fprintf(buf, "\"%s\": \"%s\",", "spdxVersion", ci.SPDXVersion)
	}
	if ci.CreatorComment != "" || ci.Created != "" || ci.CreatorPersons != nil || ci.CreatorOrganizations != nil || ci.CreatorTools != nil || ci.LicenseListVersion != "" {
		fmt.Fprintf(buf, "\"%s\": %s", "creationInfo", "{")
		if ci.CreatorComment != "" {
			commentjson, _ := json.Marshal(ci.CreatorComment)
			fmt.Fprintf(buf, "\"%s\": %s,", "comment", commentjson)
		}
		if ci.Created != "" {
			fmt.Fprintf(buf, "\"%s\": \"%s\",", "created", ci.Created)
		}
		if ci.CreatorPersons != nil || ci.CreatorOrganizations != nil || ci.CreatorTools != nil {
			var creators []string
			for _, v := range ci.CreatorPersons {
				creators = append(creators, fmt.Sprintf("Person: %s", v))
			}
			for _, v := range ci.CreatorOrganizations {
				creators = append(creators, fmt.Sprintf("Organization: %s", v))
			}
			for _, v := range ci.CreatorTools {
				creators = append(creators, fmt.Sprintf("Tool: %s", v))
			}
			creatorsjson, _ := json.Marshal(creators)
			fmt.Fprintf(buf, "\"%s\": %s ,", "creators", creatorsjson)
		}
		if ci.LicenseListVersion != "" {
			fmt.Fprintf(buf, "\"%s\": \"%s\",", "licenseListVersion", ci.LicenseListVersion)
		}
		fmt.Fprintf(buf, "%s", "},")
	}
	if ci.DocumentName != "" {
		fmt.Fprintf(buf, "\"%s\": \"%s\",", "name", ci.DocumentName)
	}
	if ci.DataLicense != "" {
		fmt.Fprintf(buf, "\"%s\": \"%s\",", "dataLicense", ci.DataLicense)
	}
	if ci.DocumentComment != "" {
		fmt.Fprintf(buf, "\"%s\": \"%s\",", "comment", ci.DocumentComment)
	}
	if ci.ExternalDocumentReferences != nil {
		var refs []interface{}
		for _, v := range ci.ExternalDocumentReferences {
			aa := make(map[string]interface{})
			aa["externalDocumentId"] = v.DocumentRefID
			aa["checksum"] = map[string]string{
				"algorithm":     v.Alg,
				"checksumValue": v.Checksum,
			}
			aa["spdxDocument"] = v.URI
			refs = append(refs, aa)
		}
		refsjson, _ := json.Marshal(refs)
		fmt.Fprintf(buf, "\"%s\": %s ,", "externalDocumentRefs", refsjson)
	}

	return nil
}
