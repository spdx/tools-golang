package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
	"strings"
)

var FileInfoHeadersByColumn = map[string]string{
	"A": common.FileInfoFileName,
	"B": common.FileInfoSPDXIdentifier,
	"C": common.FileInfoPackageIdentifier,
	"D": common.FileInfoFileTypes,
	"E": common.FileInfoFileChecksums,
	"F": common.FileInfoLicenseConcluded,
	"G": common.FileInfoLicenseInfoInFile,
	"H": common.FileInfoLicenseComments,
	"I": common.FileInfoFileCopyrightText,
	"J": common.FileInfoNoticeText,
	"K": common.FileInfoArtifactOfProject,
	"L": common.FileInfoArtifactOfHomepage,
	"M": common.FileInfoArtifactOfURL,
	"N": common.FileInfoContributors,
	"O": common.FileInfoFileComment,
	"P": common.FileInfoFileDependencies,
	"Q": common.FileInfoAttributionText,
}

func WriteFileInfoRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	rowNum := 2

	// files can appear at the document level, or the package level

	// document-level
	for _, file := range doc.Files {
		err := processFileInfo(file, "", spreadsheet, rowNum)
		if err != nil {
			return fmt.Errorf("failed to process document-level file info: %s", err.Error())
		}

		rowNum += 1
	}

	// package-level
	for _, pkg := range doc.Packages {
		for _, file := range pkg.Files {
			err := processFileInfo(file, pkg.PackageSPDXIdentifier, spreadsheet, rowNum)
			if err != nil {
				return fmt.Errorf("failed to process package-level file info: %s", err.Error())
			}

			rowNum += 1
		}
	}

	return nil
}

func processFileInfo(file *spdx.File2_2, packageIdentifier spdx.ElementID, spreadsheet *excelize.File, rowNum int) error {
	for column, valueType := range FileInfoHeadersByColumn {
		axis := common.PositionToAxis(column, rowNum)

		// set `value` to the value to be written to the spreadsheet cell
		var value interface{}
		// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
		var err error

		switch valueType {
		case common.FileInfoFileName:
			value = file.FileName
		case common.FileInfoSPDXIdentifier:
			value = file.FileSPDXIdentifier
		case common.FileInfoPackageIdentifier:
			// a file can optionally be associated with a package
			value = ""
			if packageIdentifier != "" {
				value = packageIdentifier
			}
		case common.FileInfoFileTypes:
			value = strings.Join(file.FileTypes, "\n")
		case common.FileInfoFileChecksums:
			checksums := make([]string, 0, len(file.Checksums))
			for _, checksum := range file.Checksums {
				if err = checksum.Validate(); err != nil {
					break
				}

				checksums = append(checksums, checksum.String())
			}

			value = strings.Join(checksums, "\n")
		case common.FileInfoLicenseConcluded:
			value = file.LicenseConcluded
		case common.FileInfoLicenseInfoInFile:
			value = strings.Join(file.LicenseInfoInFiles, ", ")
		case common.FileInfoLicenseComments:
			value = file.LicenseComments
		case common.FileInfoFileCopyrightText:
			value = file.FileCopyrightText
		case common.FileInfoNoticeText:
			value = file.FileNotice
		case common.FileInfoArtifactOfProject:
			// ignored
		case common.FileInfoArtifactOfHomepage:
			// ignored
		case common.FileInfoArtifactOfURL:
			// ignored
		case common.FileInfoContributors:
			contributors := make([]string, 0, len(file.FileContributors))
			for _, contributor := range file.FileContributors {
				// these get wrapped in quotes
				contributors = append(contributors, fmt.Sprintf("\"%s\"", contributor))
			}
			value = strings.Join(contributors, ",")
		case common.FileInfoFileComment:
			value = file.FileComment
		case common.FileInfoFileDependencies:
			// ignored
		case common.FileInfoAttributionText:
			texts := make([]string, 0, len(file.FileAttributionTexts))
			for _, text := range file.FileAttributionTexts {
				// these get wrapped in quotes
				texts = append(texts, fmt.Sprintf("\"%s\"", text))
			}
			value = strings.Join(texts, "\n")
		}

		if err != nil {
			return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
		}

		err = spreadsheet.SetCellValue(common.SheetNameFileInfo, axis, value)
		if err != nil {
			return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
		}
	}

	return nil
}
