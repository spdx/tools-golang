package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
	"strings"
)

var SnippetsHeadersByColumn = map[string]string{
	"A": common.SnippetsID,
	"B": common.SnippetsName,
	"C": common.SnippetsFromFileID,
	"D": common.SnippetsByteRange,
	"E": common.SnippetsLineRange,
	"F": common.SnippetsLicenseConcluded,
	"G": common.SnippetsLicenseInfoInSnippet,
	"H": common.SnippetsLicenseComments,
	"I": common.SnippetsCopyrightText,
	"J": common.SnippetsComment,
}

func WriteSnippetsRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	for ii, snippet := range doc.Snippets {
		// get correct row number. first row is headers (+1) and Go slices are zero-indexed (+1)
		rowNum := ii + 2

		for column, valueType := range SnippetsHeadersByColumn {
			axis := common.PositionToAxis(column, rowNum)

			// set `value` to the value to be written to the spreadsheet cell
			var value interface{}
			// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
			var err error

			switch valueType {
			case common.SnippetsID:
				value = snippet.SnippetSPDXIdentifier
			case common.SnippetsName:
				value = snippet.SnippetName
			case common.SnippetsFromFileID:
				value = snippet.SnippetFromFileSPDXIdentifier
			case common.SnippetsByteRange:
				// find a byte range, if there is one
				value = ""
				for _, snippetRange := range snippet.Ranges {
					if snippetRange.EndPointer.Offset != 0 {
						value = snippetRange.String()
						break
					}
				}
			case common.SnippetsLineRange:
				// find a line range, if there is one
				value = ""
				for _, snippetRange := range snippet.Ranges {
					if snippetRange.EndPointer.LineNumber != 0 {
						value = snippetRange.String()
						break
					}
				}
			case common.SnippetsLicenseConcluded:
				value = snippet.SnippetLicenseConcluded
			case common.SnippetsLicenseInfoInSnippet:
				value = strings.Join(snippet.LicenseInfoInSnippet, ", ")
			case common.SnippetsLicenseComments:
				value = snippet.SnippetLicenseComments
			case common.SnippetsCopyrightText:
				value = snippet.SnippetCopyrightText
			case common.SnippetsComment:
				value = snippet.SnippetComment
			}

			if err != nil {
				return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
			}

			err = spreadsheet.SetCellValue(common.SheetNameSnippets, axis, value)
			if err != nil {
				return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
			}
		}
	}

	return nil
}
