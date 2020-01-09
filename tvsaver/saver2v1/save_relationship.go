// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

func renderRelationship2_1(rln *spdx.Relationship2_1, w io.Writer) error {
	if rln.RefA != "" && rln.RefB != "" && rln.Relationship != "" {
		fmt.Fprintf(w, "Relationship: %s %s %s\n", rln.RefA, rln.Relationship, rln.RefB)
	}
	if rln.RelationshipComment != "" {
		fmt.Fprintf(w, "RelationshipComment: %s\n", rln.RelationshipComment)
	}

	return nil
}
