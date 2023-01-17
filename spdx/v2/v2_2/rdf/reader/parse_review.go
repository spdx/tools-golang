// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func (parser *rdfParser2_2) setReviewFromNode(reviewedNode *gordfParser.Node) error {
	review := v2_2.Review{}
	for _, triple := range parser.nodeToTriples(reviewedNode) {
		switch triple.Predicate.ID {
		case RDF_TYPE:
			// cardinality: exactly 1
			continue
		case RDFS_COMMENT:
			// cardinality: max 1
			review.ReviewComment = triple.Object.ID
		case SPDX_REVIEW_DATE:
			// cardinality: exactly 1
			review.ReviewDate = triple.Object.ID
		case SPDX_REVIEWER:
			// cardinality: max 1
			var err error
			review.ReviewerType, review.Reviewer, err = ExtractSubs(triple.Object.ID, ":")
			if err != nil {
				return fmt.Errorf("error parsing reviewer: %v", err)
			}
		default:
			return fmt.Errorf("unknown predicate %v for review triples", triple.Predicate)
		}
	}
	parser.doc.Reviews = append(parser.doc.Reviews, &review)
	return nil
}
