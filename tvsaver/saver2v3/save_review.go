// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v3

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx/v2_3"
)

func renderReview2_3(rev *v2_3.Review, w io.Writer) error {
	if rev.Reviewer != "" && rev.ReviewerType != "" {
		fmt.Fprintf(w, "Reviewer: %s: %s\n", rev.ReviewerType, rev.Reviewer)
	}
	if rev.ReviewDate != "" {
		fmt.Fprintf(w, "ReviewDate: %s\n", rev.ReviewDate)
	}
	if rev.ReviewComment != "" {
		fmt.Fprintf(w, "ReviewComment: %s\n", textify(rev.ReviewComment))
	}

	fmt.Fprintf(w, "\n")

	return nil
}
