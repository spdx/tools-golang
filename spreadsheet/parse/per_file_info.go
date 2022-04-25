package parse

import (
	"encoding/csv"
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"strings"
)

func ProcessPerFileInfoRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2
		newFile := spdx.File2_2{}
		var associatedPackageSPDXID spdx.ElementID

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.FileInfoFileName:
				newFile.FileName = value
			case common.FileInfoSPDXIdentifier:
				var id spdx.DocElementID
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid SPDX Identifier in row %d: %v", rowNum, err.Error())
				}

				newFile.FileSPDXIdentifier = id.ElementRefID
			case common.FileInfoPackageIdentifier:
				// in spreadsheet formats, file<->package relationships are dictated by this column.
				// if there is no value in this column, the file is not associated with a particular package
				var id spdx.DocElementID
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid Package SPDX Identifier in row %d: %v", rowNum, err.Error())
				}

				associatedPackageSPDXID = id.ElementRefID
			case common.FileInfoFileTypes:
				newFile.FileTypes = strings.Split(value, ", ")
			case common.FileInfoFileChecksums:
				checksums := strings.Split(value, "\n")
				for _, checksumStr := range checksums {
					checksum := spdx.Checksum{}
					err := checksum.FromString(checksumStr)
					if err != nil {
						return fmt.Errorf("invalid File Checksum in row %d: %v", rowNum, err.Error())
					}

					newFile.Checksums = append(newFile.Checksums, checksum)
				}
			case common.FileInfoLicenseConcluded:
				newFile.LicenseConcluded = value
			case common.FileInfoLicenseInfoInFile:
				newFile.LicenseInfoInFiles = strings.Split(value, ", ")
			case common.FileInfoLicenseComments:
				newFile.LicenseComments = value
			case common.FileInfoFileCopyrightText:
				newFile.FileCopyrightText = value
			case common.FileInfoNoticeText:
				newFile.FileNotice = value
			case common.FileInfoArtifactOfProject:
				// ignored
			case common.FileInfoArtifactOfHomepage:
				// ignored
			case common.FileInfoArtifactOfURL:
				// ignored
			case common.FileInfoContributors:
				contributors, err := csv.NewReader(strings.NewReader(value)).Read()
				if err != nil {
					return fmt.Errorf("invalid File Contributors in row %d: %s", rowNum, err.Error())
				}
				newFile.FileContributors = contributors
			case common.FileInfoFileComment:
				newFile.FileComment = value
			case common.FileInfoFileDependencies:
				newFile.FileDependencies = strings.Split(value, ", ")
			case common.FileInfoAttributionText:
				newFile.FileAttributionTexts = strings.Split(value, ", ")
			}
		}

		// TODO: validate?
		doc.Files = append(doc.Files, &newFile)

		// add this file to the associated package, if it is associated with a package
		if associatedPackageSPDXID != "" {
			for ii, pkg := range doc.Packages {
				if pkg.PackageSPDXIdentifier == associatedPackageSPDXID {
					doc.Packages[ii].Files = append(doc.Packages[ii].Files, &newFile)
					break
				}
			}
		}
	}

	return nil
}
