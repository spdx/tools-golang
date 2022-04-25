package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
)

var RelationshipsHeadersByColumn = map[string]string{
	"A": common.RelationshipsRefA,
	"B": common.RelationshipsRelationship,
	"C": common.RelationshipsRefB,
	"D": common.RelationshipsComment,
}

func WriteRelationshipsRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	for ii, relationship := range doc.Relationships {
		// get correct row number. first row is headers (+1) and Go slices are zero-indexed (+1)
		rowNum := ii + 2

		for column, valueType := range RelationshipsHeadersByColumn {
			axis := common.PositionToAxis(column, rowNum)

			// set `value` to the value to be written to the spreadsheet cell
			var value interface{}
			// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
			var err error

			switch valueType {
			case common.RelationshipsRefA:
				value = relationship.RefA
			case common.RelationshipsRelationship:
				value = relationship.Relationship
			case common.RelationshipsRefB:
				value = relationship.RefB
			case common.RelationshipsComment:
				value = relationship.RelationshipComment
			}

			if err != nil {
				return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
			}

			err = spreadsheet.SetCellValue(common.SheetNameRelationships, axis, value)
			if err != nil {
				return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
			}
		}
	}

	return nil
}
