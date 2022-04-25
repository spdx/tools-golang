package parse

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
)

func ProcessPackageExternalRefsRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2
		// each external ref is related to a package, make sure we figure out which package
		var packageSPDXID spdx.ElementID
		newExternalRef := spdx.PackageExternalReference2_2{}

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.ExternalRefPackageID:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid Package SPDX Identifier for External Ref in row %d: %v", rowNum, err.Error())
				}

				packageSPDXID = id.ElementRefID
			case common.ExternalRefCategory:
				newExternalRef.Category = value
			case common.ExternalRefType:
				newExternalRef.RefType = value
			case common.ExternalRefLocator:
				newExternalRef.Locator = value
			case common.ExternalRefComment:
				newExternalRef.ExternalRefComment = value
			}
		}

		if packageSPDXID == "" {
			return fmt.Errorf("no SPDX ID given for package external ref in row %d", rowNum)
		}

		// find the package this external ref is related to
		var packageFound bool
		for ii, pkg := range doc.Packages {
			if pkg.PackageSPDXIdentifier == packageSPDXID {
				packageFound = true
				doc.Packages[ii].PackageExternalReferences = append(doc.Packages[ii].PackageExternalReferences, &newExternalRef)
				break
			}
		}

		if !packageFound {
			return fmt.Errorf("package external ref assigned to non-existent package %s in row %d", packageSPDXID, rowNum)
		}
	}

	return nil
}
