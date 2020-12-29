// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

func renderReview2_2(rev *spdx.Review2_2, w io.Writer) error {
	if rev.Reviewer != "" && rev.ReviewerType != "" {
		fmt.Fprintf(w, "Reviewer: %s: %s\n", rev.ReviewerType, rev.Reviewer)
	}
	if rev.Date != "" {
		fmt.Fprintf(w, "ReviewDate: %s\n", rev.Date)
	}
	if rev.Comment != "" {
		fmt.Fprintf(w, "ReviewComment: %s\n", rev.Comment)
	}

	fmt.Fprintf(w, "\n")

	return nil
}
