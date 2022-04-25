// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx_xls

import (
	"errors"
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/xuri/excelize/v2"
	"io"
)

// Load2_2 takes in an io.Reader and returns an SPDX document.
func Load2_2(content io.Reader) (*spdx.Document2_2, error) {
	workbook, err := excelize.OpenReader(content)
	if err != nil {
		return nil, err
	}

	doc, err := parseWorkbook(workbook)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func parseWorkbook(workbook *excelize.File) (*spdx.Document2_2, error) {
	doc := spdx.Document2_2{
		// ensure this pointer is not nil
		CreationInfo: &spdx.CreationInfo2_2{},
	}

	for _, sheetHandlingInfo := range sheetHandlers {
		rows, err := workbook.GetRows(sheetHandlingInfo.SheetName)
		if err != nil {
			// if the sheet doesn't exist and is required, that's a problem
			if errors.As(err, &excelize.ErrSheetNotExist{}) {
				if sheetHandlingInfo.SheetIsRequired {
					return nil, fmt.Errorf("sheet '%s' is required but is not present", sheetHandlingInfo.SheetName)
				} else {
					// if it is not required, skip it
					continue
				}
			} else {
				// some other error happened
				return nil, err
			}
		}

		// the first row is column headers, and the next row would contain actual data.
		// if there are less than 2 rows present, there is no actual data in the sheet.
		if len(rows) < 2 {
			if sheetHandlingInfo.SheetIsRequired {
				return nil, fmt.Errorf("sheet '%s' is required but contains no data", sheetHandlingInfo.SheetName)
			}

			continue
		}

		err = sheetHandlingInfo.ParserFunc(rows, &doc)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sheet '%s': %w", sheetHandlingInfo.SheetName, err)
		}
	}

	return &doc, nil
}
