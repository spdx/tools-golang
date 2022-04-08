// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"fmt"
	"path/filepath"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/utils"
)

// BuildFileSection2_1 creates an SPDX File (version 2.1), returning that
// file or error if any is encountered. Arguments:
//   - filePath: path to file, relative to prefix
//   - prefix: relative directory for filePath
//   - fileNumber: integer index (unique within package) to use in identifier
func BuildFileSection2_1(filePath string, prefix string, fileNumber int) (*spdx.File2_1, error) {
	// build the full file path
	p := filepath.Join(prefix, filePath)

	// make sure we can get the file and its hashes
	ssha1, ssha256, smd5, err := utils.GetHashesForFilePath(p)
	if err != nil {
		return nil, err
	}

	// build the identifier
	i := fmt.Sprintf("File%d", fileNumber)

	// now build the File section
	f := &spdx.File2_1{
		FileName:           filePath,
		FileSPDXIdentifier: spdx.ElementID(i),
		Checksums: []spdx.Checksum{
			{
				Algorithm: spdx.SHA1,
				Value:     ssha1,
			},
			{
				Algorithm: spdx.SHA256,
				Value:     ssha256,
			},
			{
				Algorithm: spdx.MD5,
				Value:     smd5,
			},
		},
		LicenseConcluded:   "NOASSERTION",
		LicenseInfoInFiles: []string{"NOASSERTION"},
		FileCopyrightText:  "NOASSERTION",
	}

	return f, nil
}
