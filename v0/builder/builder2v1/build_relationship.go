// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"fmt"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// BuildRelationshipSection2_1 creates an SPDX Relationship (version 2.1)
// solely for the document "DESCRIBES" package relationship, returning that
// relationship or error if any is encountered. Arguments:
//   - packageName: name of package / directory
func BuildRelationshipSection2_1(packageName string) (*spdx.Relationship2_1, error) {
	rln := &spdx.Relationship2_1{
		RefA:         "SPDXRef-DOCUMENT",
		RefB:         fmt.Sprintf("SPDXRef-Package-%s", packageName),
		Relationship: "DESCRIBES",
	}

	return rln, nil
}
