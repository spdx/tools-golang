package parse

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
)

func ProcessRelationshipsRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2
		newRelationship := spdx.Relationship2_2{}

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.RelationshipsRefA:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.RelationshipsRefA, rowNum, err.Error())
				}

				newRelationship.RefA = id
			case common.RelationshipsRelationship:
				newRelationship.Relationship = value
			case common.RelationshipsRefB:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.RelationshipsRefB, rowNum, err.Error())
				}

				newRelationship.RefB = id
			case common.RelationshipsComment:
				newRelationship.RelationshipComment = value
			}
		}

		// TODO: validate?
		doc.Relationships = append(doc.Relationships, &newRelationship)
	}

	return nil
}
