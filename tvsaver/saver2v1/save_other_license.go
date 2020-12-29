// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

func renderOtherLicense2_1(ol *spdx.OtherLicense2_1, w io.Writer) error {
	if ol.Identifier != "" {
		fmt.Fprintf(w, "LicenseID: %s\n", ol.Identifier)
	}
	if ol.ExtractedText != "" {
		fmt.Fprintf(w, "ExtractedText: %s\n", textify(ol.ExtractedText))
	}
	if ol.Name != "" {
		fmt.Fprintf(w, "LicenseName: %s\n", ol.Name)
	}
	for _, s := range ol.CrossReferences {
		fmt.Fprintf(w, "LicenseCrossReference: %s\n", s)
	}
	if ol.Comment != "" {
		fmt.Fprintf(w, "LicenseComment: %s\n", textify(ol.Comment))
	}

	fmt.Fprintf(w, "\n")

	return nil
}
