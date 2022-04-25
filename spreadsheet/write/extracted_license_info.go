package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
	"strings"
)

var ExtractedLicenseInfoHeadersByColumn = map[string]string{
	"A": common.LicenseInfoIdentifier,
	"B": common.LicenseInfoExtractedText,
	"C": common.LicenseInfoLicenseName,
	"D": common.LicenseInfoCrossReferenceURLs,
	"E": common.LicenseInfoComment,
}

func WriteExtractedLicenseInfoRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	for ii, license := range doc.OtherLicenses {
		// get correct row number. first row is headers (+1) and Go slices are zero-indexed (+1)
		rowNum := ii + 2

		for column, valueType := range ExtractedLicenseInfoHeadersByColumn {
			axis := common.PositionToAxis(column, rowNum)

			// set `value` to the value to be written to the spreadsheet cell
			var value interface{}
			// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
			var err error

			switch valueType {
			case common.LicenseInfoIdentifier:
				value = license.LicenseIdentifier
			case common.LicenseInfoExtractedText:
				value = license.ExtractedText
			case common.LicenseInfoLicenseName:
				value = license.LicenseName
			case common.LicenseInfoCrossReferenceURLs:
				value = strings.Join(license.LicenseCrossReferences, ", ")
			case common.LicenseInfoComment:
				value = license.LicenseComment
			}

			if err != nil {
				return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
			}

			err = spreadsheet.SetCellValue(common.SheetNameExtractedLicenseInfo, axis, value)
			if err != nil {
				return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
			}
		}
	}

	return nil
}
