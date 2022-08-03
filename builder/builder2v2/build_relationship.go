// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// BuildRelationshipSection2_2 creates an SPDX Relationship (version 2.2)
// solely for the document "DESCRIBES" package relationship, returning that
// relationship or error if any is encountered. Arguments:
//   - packageName: name of package / directory
func BuildRelationshipSection2_2(packageName string) (*v2_2.Relationship, error) {
	rln := &v2_2.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", fmt.Sprintf("Package-%s", packageName)),
		Relationship: "DESCRIBES",
	}

	return rln, nil
}
