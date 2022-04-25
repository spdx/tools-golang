package parse

import (
	"encoding/csv"
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"strconv"
	"strings"
)

func ProcessPackageInfoRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2
		newPackage := spdx.Package2_2{}

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.PackageName:
				newPackage.PackageName = value
			case common.PackageSPDXIdentifier:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid Package SPDX Identifier in row %d: %v", rowNum, err.Error())
				}

				newPackage.PackageSPDXIdentifier = id.ElementRefID
			case common.PackageVersion:
				newPackage.PackageVersion = value
			case common.PackageFileName:
				newPackage.PackageFileName = value
			case common.PackageSupplier:
				supplier := spdx.Supplier{}
				err := supplier.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid Package Supplier in row %d: %v", rowNum, err.Error())
				}

				newPackage.PackageSupplier = &supplier
			case common.PackageOriginator:
				originator := spdx.Originator{}
				err := originator.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid Package Originator in row %d: %v", rowNum, err.Error())
				}

				newPackage.PackageOriginator = &originator
			case common.PackageHomePage:
				newPackage.PackageHomePage = value
			case common.PackageDownloadLocation:
				newPackage.PackageDownloadLocation = value
			case common.PackageChecksum:
				checksums := strings.Split(value, "\n")
				for _, checksumStr := range checksums {
					checksum := spdx.Checksum{}
					err := checksum.FromString(checksumStr)
					if err != nil {
						return fmt.Errorf("invalid Package Checksum in row %d: %v", rowNum, err.Error())
					}

					newPackage.PackageChecksums = append(newPackage.PackageChecksums, checksum)
				}
			case common.PackageVerificationCode:
				newPackage.PackageVerificationCode.Value = value
			case common.PackageVerificationCodeExcludedFiles:
				newPackage.PackageVerificationCode.ExcludedFiles = append(newPackage.PackageVerificationCode.ExcludedFiles, value)
			case common.PackageSourceInfo:
				newPackage.PackageSourceInfo = value
			case common.PackageLicenseDeclared:
				newPackage.PackageLicenseDeclared = value
			case common.PackageLicenseConcluded:
				newPackage.PackageLicenseConcluded = value
			case common.PackageLicenseInfoFromFiles:
				files := strings.Split(value, ",")
				newPackage.PackageLicenseInfoFromFiles = append(newPackage.PackageLicenseInfoFromFiles, files...)
			case common.PackageLicenseComments:
				newPackage.PackageLicenseComments = value
			case common.PackageCopyrightText:
				newPackage.PackageCopyrightText = value
			case common.PackageSummary:
				newPackage.PackageSummary = value
			case common.PackageDescription:
				newPackage.PackageDescription = value
			case common.PackageAttributionText:
				attributionTexts, err := csv.NewReader(strings.NewReader(value)).Read()
				if err != nil {
					return fmt.Errorf("invalid Package Attribution Text in row %d: %s", rowNum, err.Error())
				}
				newPackage.PackageAttributionTexts = attributionTexts
			case common.PackageFilesAnalyzed:
				filesAnalyzed, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("invalid boolean for Files Analyzed in row %d (should be 'true' or 'false')", rowNum)
				}

				newPackage.FilesAnalyzed = filesAnalyzed
			case common.PackageComments:
				newPackage.PackageComment = value
			}
		}

		// TODO: validate?
		doc.Packages = append(doc.Packages, &newPackage)
	}

	return nil
}
