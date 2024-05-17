// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/utils"
)

// BuildPackageSection creates an SPDX Package, returning
// that package or error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - dirRoot: path to directory to be analyzed
//   - pathsIgnore: slice of strings for filepaths to ignore
func BuildPackageSection(packageName string, dirRoot string, pathsIgnore []string) (*spdx.Package, error) {
	// build the file section first, so we'll have it available
	// for calculating the package verification code
	relativePaths, err := utils.GetAllFilePaths(dirRoot, pathsIgnore)

	if err != nil {
		return nil, err
	}

	files := []*spdx.File{}
	fileNumber := 0
	for _, filePath := range relativePaths {
		// SPDX spec says file names should generally start with ./
		// see: https://spdx.github.io/spdx-spec/v2.3/file-information/#81-file-name-field
		relativePath := "." + filePath
		newFile, err := BuildFileSection(relativePath, dirRoot, fileNumber)
		if err != nil {
			return nil, err
		}
		files = append(files, newFile)
		fileNumber++
	}

	// get the verification code
	code, err := utils.GetVerificationCode(files, "")
	if err != nil {
		return nil, err
	}

	// now build the package section
	pkg := &spdx.Package{
		PackageName:                 packageName,
		PackageSPDXIdentifier:       common.ElementID(fmt.Sprintf("Package-%s", packageName)),
		PackageDownloadLocation:     "NOASSERTION",
		FilesAnalyzed:               true,
		IsFilesAnalyzedTagPresent:   true,
		PackageVerificationCode:     &code,
		PackageLicenseConcluded:     "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{},
		PackageLicenseDeclared:      "NOASSERTION",
		PackageCopyrightText:        "NOASSERTION",
		Files:                       files,
	}

	return pkg, nil
}
