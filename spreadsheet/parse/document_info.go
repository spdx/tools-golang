package parse

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
)

func ProcessDocumentInfoRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2
		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.SPDXVersion:
				doc.SPDXVersion = value
			case common.DataLicense:
				doc.DataLicense = value
			case common.SPDXIdentifier:
				var id spdx.DocElementID
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid SPDX Identifier in row %d: %v", rowNum, err.Error())
				}

				doc.SPDXIdentifier = id.ElementRefID
			case common.LicenseListVersion:
				doc.CreationInfo.LicenseListVersion = value
			case common.DocumentName:
				doc.DocumentName = value
			case common.DocumentNamespace:
				doc.DocumentNamespace = value
			case common.DocumentComment:
				doc.DocumentComment = value
			case common.ExternalDocumentReferences:
				externalDocRef := spdx.ExternalDocumentRef2_2{}
				err := externalDocRef.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid External Document Ref in row %d: %v", rowNum, err.Error())
				}

				doc.ExternalDocumentReferences = append(doc.ExternalDocumentReferences, externalDocRef)
			case common.Created:
				doc.CreationInfo.Created = value
			case common.CreatorComment:
				doc.CreationInfo.CreatorComment = value
			case common.Creator:
				creator := spdx.Creator{}
				err := creator.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid Creator in row %d: %v", rowNum, err.Error())
				}

				doc.CreationInfo.Creators = append(doc.CreationInfo.Creators, creator)
			}
		}
	}

	return nil
}
