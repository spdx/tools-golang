package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
)

var DocumentInfoHeadersByColumn = map[string]string{
	"A": common.DocumentInfoSPDXVersion,
	"B": common.DocumentInfoDataLicense,
	"C": common.DocumentInfoSPDXIdentifier,
	"D": common.DocumentInfoLicenseListVersion,
	"E": common.DocumentInfoDocumentName,
	"F": common.DocumentInfoDocumentNamespace,
	"G": common.DocumentInfoExternalDocumentReferences,
	"H": common.DocumentInfoDocumentComment,
	"I": common.DocumentInfoCreator,
	"J": common.DocumentInfoCreated,
	"K": common.DocumentInfoCreatorComment,
}

func WriteDocumentInfoRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	if doc.CreationInfo == nil {
		return fmt.Errorf("document is missing CreationInfo")
	}

	// some data in this sheet gets split across rows, instead of being split up by newlines or commas.
	// the two columns where this happens are Creators and External Document Refs.
	// figure out how many rows we're going to need
	numCreators := len(doc.CreationInfo.Creators)
	numExternalDocRefs := len(doc.ExternalDocumentReferences)
	rowsNeeded := 1
	if numCreators > numExternalDocRefs {
		rowsNeeded = numCreators
	} else if numExternalDocRefs > 1 {
		rowsNeeded = numExternalDocRefs
	}

	for rowNum := 2; rowNum-2 < rowsNeeded; rowNum++ {
		for column, valueType := range DocumentInfoHeadersByColumn {
			// only certain columns are used past the first data row
			if rowNum > 2 && valueType != common.DocumentInfoCreator && valueType != common.DocumentInfoExternalDocumentReferences {
				continue
			}

			axis := common.PositionToAxis(column, rowNum)

			// set `value` to the value to be written to the spreadsheet cell
			var value interface{}
			// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
			var err error

			switch valueType {
			case common.DocumentInfoSPDXVersion:
				value = doc.SPDXVersion
			case common.DocumentInfoDataLicense:
				value = doc.DataLicense
			case common.DocumentInfoSPDXIdentifier:
				value = doc.SPDXIdentifier
			case common.DocumentInfoLicenseListVersion:
				value = doc.CreationInfo.LicenseListVersion
			case common.DocumentInfoDocumentName:
				value = doc.DocumentName
			case common.DocumentInfoDocumentNamespace:
				value = doc.DocumentNamespace
			case common.DocumentInfoExternalDocumentReferences:
				if rowNum-2 > numExternalDocRefs-1 {
					continue
				}

				ref := doc.ExternalDocumentReferences[rowNum-2]
				if err = ref.Validate(); err != nil {
					break
				}

				value = ref.String()
			case common.DocumentInfoDocumentComment:
				value = doc.DocumentComment
			case common.DocumentInfoCreator:
				if rowNum-2 > numCreators-1 {
					continue
				}

				creator := doc.CreationInfo.Creators[rowNum-2]
				if err = creator.Validate(); err != nil {
					break
				}

				value = creator.String()
			case common.DocumentInfoCreated:
				value = doc.CreationInfo.Created
			case common.DocumentInfoCreatorComment:
				value = doc.CreationInfo.CreatorComment
			}

			if err != nil {
				return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
			}

			err = spreadsheet.SetCellValue(common.SheetNameDocumentInfo, axis, value)
			if err != nil {
				return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
			}
		}
	}

	return nil
}
