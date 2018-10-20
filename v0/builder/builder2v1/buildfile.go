// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"fmt"
	"path/filepath"

	"github.com/swinslow/spdx-go/v0/spdx"
)

func buildFileSection2_1(filePath string, prefix string, fileNumber int) (*spdx.File2_1, error) {
	// build the full file path
	p := filepath.Join(prefix, filePath)

	// make sure we can get the file and its hashes
	ssha1, ssha256, smd5, err := getHashesForFilePath(p)
	if err != nil {
		return nil, err
	}

	// build the identifier
	i := fmt.Sprintf("SPDXRef-File%d", fileNumber)

	// now build the File section
	f := &spdx.File2_1{
		FileName:           filePath,
		FileSPDXIdentifier: i,
		FileChecksumSHA1:   ssha1,
		FileChecksumSHA256: ssha256,
		FileChecksumMD5:    smd5,
		LicenseConcluded:   "NOASSERTION",
		LicenseInfoInFile:  []string{},
		FileCopyrightText:  "NOASSERTION",
	}

	return f, nil
}
