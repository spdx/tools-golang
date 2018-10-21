// Package idsearcher is used to search for short-form IDs in files
// within a directory, and to build an SPDX Document containing those
// license findings.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package idsearcher

import (
	"bufio"
	"os"
	"strings"
)

// ===== Utility functions =====
func searchFileIDs(filePath string) ([]string, error) {
	ids := []string{}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "SPDX-License-Identifier:") {
			strs := strings.SplitAfterN(scanner.Text(), "SPDX-License-Identifier:", 2)
			lid := strings.TrimSpace(strs[1])
			ids = append(ids, lid)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
