package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
)

var ExternalRefsHeadersByColumn = map[string]string{
	"A": common.ExternalRefPackageID,
	"B": common.ExternalRefCategory,
	"C": common.ExternalRefType,
	"D": common.ExternalRefLocator,
	"E": common.ExternalRefComment,
}

func WriteExternalRefsRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	rowNum := 2

	for _, pkg := range doc.Packages {
		for _, externalRef := range pkg.PackageExternalReferences {
			for column, valueType := range ExternalRefsHeadersByColumn {
				axis := common.PositionToAxis(column, rowNum)

				// set `value` to the value to be written to the spreadsheet cell
				var value interface{}
				// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
				var err error

				switch valueType {
				case common.ExternalRefPackageID:
					value = pkg.PackageSPDXIdentifier
				case common.ExternalRefCategory:
					value = externalRef.Category
				case common.ExternalRefType:
					value = externalRef.RefType
				case common.ExternalRefLocator:
					value = externalRef.Locator
				case common.ExternalRefComment:
					value = externalRef.ExternalRefComment
				}

				if err != nil {
					return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
				}

				err = spreadsheet.SetCellValue(common.SheetNameExternalRefs, axis, value)
				if err != nil {
					return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
				}
			}

			rowNum += 1
		}
	}

	return nil
}
