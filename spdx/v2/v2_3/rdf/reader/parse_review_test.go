// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"testing"
)

func Test_rdfParser2_3_setReviewFromNode(t *testing.T) {
	// TestCase 1: unknown predicate must raise an error
	parser, _ := parserFromBodyContent(`
		<spdx:Review>
			<rdfs:comment>Another example reviewer.</rdfs:comment>
			<spdx:reviewDate>2011-03-13T00:00:00Z</spdx:reviewDate>
			<spdx:reviewer>Person: Suzanne Reviewer</spdx:reviewer>
			<spdx:unknown />
		</spdx:Review>
	`)
	reviewNode := parser.gordfParserObj.Triples[0].Subject
	err := parser.setReviewFromNode(reviewNode)
	if err == nil {
		t.Errorf("unknown predicate should've elicit an error")
	}

	// TestCase 2: wrong reviewer format must raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:Review>
			<rdfs:comment>Another example reviewer.</rdfs:comment>
			<spdx:reviewDate>2011-03-13T00:00:00Z</spdx:reviewDate>
			<spdx:reviewer>Suzanne Reviewer</spdx:reviewer>
		</spdx:Review>
	`)
	reviewNode = parser.gordfParserObj.Triples[0].Subject
	err = parser.setReviewFromNode(reviewNode)
	if err == nil {
		t.Errorf("incorrect should've elicit an error")
	}

	// TestCase 3: valid input
	parser, _ = parserFromBodyContent(`
		<spdx:Review>
			<rdfs:comment>Another example reviewer.</rdfs:comment>
			<spdx:reviewDate>2011-03-13T00:00:00Z</spdx:reviewDate>
			<spdx:reviewer>Person: Suzanne</spdx:reviewer>
		</spdx:Review>
	`)
	reviewNode = parser.gordfParserObj.Triples[0].Subject
	err = parser.setReviewFromNode(reviewNode)
	if err != nil {
		t.Errorf("error parsing a valid node")
	}
	n := len(parser.doc.Reviews)
	if n != 1 {
		t.Errorf("expected doc to have 1 review, found %d", n)
	}
	review := parser.doc.Reviews[0]
	expectedComment := "Another example reviewer."
	if review.ReviewComment != expectedComment {
		t.Errorf("expected: %v, found: %s", expectedComment, review.ReviewComment)
	}
	expectedDate := "2011-03-13T00:00:00Z"
	if review.ReviewDate != expectedDate {
		t.Errorf("expected %s, found %s", expectedDate, review.ReviewDate)
	}
	expectedReviewer := "Suzanne"
	if review.Reviewer != expectedReviewer {
		t.Errorf("expected %s, found %s", expectedReviewer, review.Reviewer)
	}
	expectedReviewerType := "Person"
	if review.ReviewerType != expectedReviewerType {
		t.Errorf("expected %s, found %s", expectedReviewerType, review.ReviewerType)
	}
}
