// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx_xls

import (
	"errors"
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/spdx/tools-golang/spreadsheet/parse"
	"github.com/xuri/excelize/v2"
	"io"
)

// sheetHandlerFunc is a func that takes in the data from a sheet as a slice of rows and iterates through them to
// fill in information in the given spdx.Document2_2.
// Returns an error if any occurred.
type sheetHandlerFunc func(rows [][]string, doc *spdx.Document2_2) error

// sheetHandlingInformation defines info that is needed for parsing individual sheets in a workbook.
type sheetHandlingInformation struct {
	// SheetName is the name of the sheet
	SheetName string
	// HandlerFunc is the function that should be used to parse a particular sheet
	HandlerFunc sheetHandlerFunc
	// SheetIsRequired denotes whether the sheet is required to be present in the workbook, or if it is optional (false)
	SheetIsRequired bool
}

// sheetHandlers contains handling information for each sheet in the workbook.
// The order of this slice determines the order in which the sheets are processed.
var sheetHandlers = []sheetHandlingInformation{
	{
		SheetName:       common.SheetNameDocumentInfo,
		HandlerFunc:     parse.ProcessDocumentInfoRows,
		SheetIsRequired: true,
	},
	{
		SheetName:       common.SheetNamePackageInfo,
		HandlerFunc:     parse.ProcessPackageInfoRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameExternalRefs,
		HandlerFunc:     parse.ProcessPackageExternalRefsRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameExtractedLicenseInfo,
		HandlerFunc:     parse.ProcessExtractedLicenseInfoRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameFileInfo,
		HandlerFunc:     parse.ProcessPerFileInfoRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameRelationships,
		HandlerFunc:     parse.ProcessRelationshipsRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameAnnotations,
		HandlerFunc:     parse.ProcessAnnotationsRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameSnippets,
		HandlerFunc:     parse.ProcessSnippetsRows,
		SheetIsRequired: false,
	},
}

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

		err = sheetHandlingInfo.HandlerFunc(rows, &doc)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sheet '%s': %w", sheetHandlingInfo.SheetName, err)
		}
	}

	return &doc, nil
}
