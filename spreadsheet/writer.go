// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx_xls

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
	"io"
)

// Save2_2 takes an SPDX Document (version 2.2) and an io.Writer, and writes the document to the writer as an XLSX file.
func Save2_2(doc *spdx.Document2_2, w io.Writer) error {
	spreadsheet := excelize.NewFile()

	for _, sheetHandlingInfo := range sheetHandlers {
		spreadsheet.NewSheet(sheetHandlingInfo.SheetName)

		err := writeHeaders(spreadsheet, sheetHandlingInfo.SheetName, sheetHandlingInfo.HeadersByColumn)

		err = sheetHandlingInfo.WriterFunc(doc, spreadsheet)
		if err != nil {
			return fmt.Errorf("failed to write data for sheet %s: %s", sheetHandlingInfo.SheetName, err.Error())
		}
	}

	err := spreadsheet.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func writeHeaders(spreadsheet *excelize.File, sheetName string, headersByColumn map[string]string) error {
	for column, header := range headersByColumn {
		err := spreadsheet.SetCellValue(sheetName, common.PositionToAxis(column, 1), header)
		if err != nil {
			return err
		}
	}

	return nil
}
