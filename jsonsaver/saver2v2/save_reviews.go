// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func renderReviews2_2(reviews []*spdx.Review2_2, buf *bytes.Buffer) error {

	var review []interface{}
	for _, v := range reviews {
		rev := make(map[string]interface{})
		rev["reviewDate"] = v.ReviewDate
		rev["reviewer"] = fmt.Sprintf("%s: %s", v.ReviewerType, v.Reviewer)
		if v.ReviewComment != "" {
			rev["comment"] = v.ReviewComment
		}
		review = append(review, rev)
	}
	reviewjson, _ := json.Marshal(review)
	fmt.Fprintf(buf, "\"%s\": %s ,", "revieweds", reviewjson)

	return nil
}
