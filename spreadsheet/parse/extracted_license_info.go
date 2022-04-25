package parse

import (
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"strings"
)

func ProcessExtractedLicenseInfoRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for _, row := range rows[1:] {
		newLicense := spdx.OtherLicense2_2{}

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.LicenseInfoIdentifier:
				newLicense.LicenseIdentifier = value
			case common.LicenseInfoExtractedText:
				newLicense.ExtractedText = value
			case common.LicenseInfoLicenseName:
				newLicense.LicenseName = value
			case common.LicenseInfoCrossReferenceURLs:
				newLicense.LicenseCrossReferences = strings.Split(value, ", ")
			case common.LicenseInfoComment:
				newLicense.LicenseComment = value
			}
		}

		// TODO: validate?
		doc.OtherLicenses = append(doc.OtherLicenses, &newLicense)
	}

	return nil
}
