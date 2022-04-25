package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
	"strings"
)

var PackageInfoHeadersByColumn = map[string]string{
	"A": common.PackageName,
	"B": common.PackageSPDXIdentifier,
	"C": common.PackageVersion,
	"D": common.PackageFileName,
	"E": common.PackageSupplier,
	"F": common.PackageOriginator,
	"G": common.PackageHomePage,
	"H": common.PackageDownloadLocation,
	"I": common.PackageChecksum,
	"J": common.PackageVerificationCode,
	"K": common.PackageVerificationCodeExcludedFiles,
	"L": common.PackageSourceInfo,
	"M": common.PackageLicenseDeclared,
	"N": common.PackageLicenseConcluded,
	"O": common.PackageLicenseInfoFromFiles,
	"P": common.PackageLicenseComments,
	"Q": common.PackageCopyrightText,
	"R": common.PackageSummary,
	"S": common.PackageDescription,
	"T": common.PackageAttributionText,
	"U": common.PackageFilesAnalyzed,
	"V": common.PackageComments,
}

func WritePackageInfoRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	for ii, pkg := range doc.Packages {
		// get correct row number. first row is headers (+1) and Go slices are zero-indexed (+1)
		rowNum := ii + 2

		for column, valueType := range PackageInfoHeadersByColumn {
			axis := common.PositionToAxis(column, rowNum)

			// set `value` to the value to be written to the spreadsheet cell
			var value interface{}
			// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
			var err error

			switch valueType {
			case common.PackageName:
				value = pkg.PackageName
			case common.PackageSPDXIdentifier:
				value = pkg.PackageSPDXIdentifier
			case common.PackageVersion:
				value = pkg.PackageVersion
			case common.PackageFileName:
				value = pkg.PackageFileName
			case common.PackageSupplier:
				if pkg.PackageSupplier == nil {
					continue
				}

				err = pkg.PackageSupplier.Validate()
				value = pkg.PackageSupplier.String()
			case common.PackageOriginator:
				if pkg.PackageOriginator == nil {
					continue
				}

				err = pkg.PackageOriginator.Validate()
				value = pkg.PackageOriginator.String()
			case common.PackageHomePage:
				value = pkg.PackageHomePage
			case common.PackageDownloadLocation:
				value = pkg.PackageDownloadLocation
			case common.PackageChecksum:
				checksums := make([]string, 0, len(pkg.PackageChecksums))
				for _, checksum := range pkg.PackageChecksums {
					if err = checksum.Validate(); err != nil {
						break
					}

					checksums = append(checksums, checksum.String())
				}

				value = strings.Join(checksums, "\n")
			case common.PackageVerificationCode:
				value = pkg.PackageVerificationCode.Value
			case common.PackageVerificationCodeExcludedFiles:
				value = strings.Join(pkg.PackageVerificationCode.ExcludedFiles, "\n")
			case common.PackageSourceInfo:
				value = pkg.PackageSourceInfo
			case common.PackageLicenseDeclared:
				value = pkg.PackageLicenseDeclared
			case common.PackageLicenseConcluded:
				value = pkg.PackageLicenseConcluded
			case common.PackageLicenseInfoFromFiles:
				value = strings.Join(pkg.PackageLicenseInfoFromFiles, ",")
			case common.PackageLicenseComments:
				value = pkg.PackageLicenseComments
			case common.PackageCopyrightText:
				value = pkg.PackageCopyrightText
			case common.PackageSummary:
				value = pkg.PackageSummary
			case common.PackageDescription:
				value = pkg.PackageDescription
			case common.PackageAttributionText:
				texts := make([]string, 0, len(pkg.PackageAttributionTexts))
				for _, text := range pkg.PackageAttributionTexts {
					// these get wrapped in quotes
					texts = append(texts, fmt.Sprintf("\"%s\"", text))
				}
				value = strings.Join(texts, "\n")
			case common.PackageFilesAnalyzed:
				value = pkg.FilesAnalyzed
			case common.PackageComments:
				value = pkg.PackageComment
			}

			if err != nil {
				return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
			}

			err = spreadsheet.SetCellValue(common.SheetNamePackageInfo, axis, value)
			if err != nil {
				return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
			}
		}
	}

	return nil
}
