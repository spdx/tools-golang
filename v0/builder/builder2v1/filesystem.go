// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"os"
	"path/filepath"
	"strings"
)

func getAllFilePaths(dirRoot string) ([]string, error) {
	// paths is a _pointer_ to a slice -- not just a slice.
	// this is so that it can be appropriately modified by append
	// in the sub-function.
	paths := &[]string{}
	prefix := strings.TrimSuffix(dirRoot, "/")

	err := filepath.Walk(dirRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// only include path if it's not a directory
		if !fi.IsDir() {
			*paths = append(*paths, strings.TrimPrefix(path, prefix))
		}
		return nil
	})

	return *paths, err
}
