// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx/v2_1"
)

func renderCreationInfo2_1(ci *v2_1.CreationInfo, w io.Writer) error {
	if ci.LicenseListVersion != "" {
		fmt.Fprintf(w, "LicenseListVersion: %s\n", ci.LicenseListVersion)
	}
	for _, creator := range ci.Creators {
		fmt.Fprintf(w, "Creator: %s: %s\n", creator.CreatorType, creator.Creator)
	}
	if ci.Created != "" {
		fmt.Fprintf(w, "Created: %s\n", ci.Created)
	}
	if ci.CreatorComment != "" {
		fmt.Fprintf(w, "CreatorComment: %s\n", textify(ci.CreatorComment))
	}

	// add blank newline b/c end of a main section
	fmt.Fprintf(w, "\n")

	return nil
}
