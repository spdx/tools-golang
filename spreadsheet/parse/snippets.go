package parse

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"strings"
)

func ProcessSnippetsRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2
		newSnippet := spdx.Snippet2_2{}

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.SnippetsID:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.SnippetsID, rowNum, err.Error())
				}

				newSnippet.SnippetSPDXIdentifier = id.ElementRefID
			case common.SnippetsName:
				newSnippet.SnippetName = value
			case common.SnippetsFromFileID:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.SnippetsFromFileID, rowNum, err.Error())
				}

				newSnippet.SnippetFromFileSPDXIdentifier = id.ElementRefID
			case common.SnippetsByteRange:
				snippetRange := spdx.SnippetRange{}
				err := snippetRange.FromString(value, true)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.SnippetsByteRange, rowNum, err.Error())
				}

				newSnippet.Ranges = append(newSnippet.Ranges, snippetRange)
			case common.SnippetsLineRange:
				snippetRange := spdx.SnippetRange{}
				err := snippetRange.FromString(value, false)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.SnippetsLineRange, rowNum, err.Error())
				}

				newSnippet.Ranges = append(newSnippet.Ranges, snippetRange)
			case common.SnippetsLicenseConcluded:
				newSnippet.SnippetLicenseConcluded = value
			case common.SnippetsLicenseInfoInSnippet:
				newSnippet.LicenseInfoInSnippet = strings.Split(value, ", ")
			case common.SnippetsLicenseComments:
				newSnippet.SnippetLicenseComments = value
			case common.SnippetsCopyrightText:
				newSnippet.SnippetCopyrightText = value
			case common.SnippetsComment:
				newSnippet.SnippetComment = value
			}
		}

		// TODO: validate?
		doc.Snippets = append(doc.Snippets, newSnippet)
	}

	return nil
}
