// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v3

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

// BuildRelationshipSection2_3 creates an SPDX Relationship (version 2.3)
// solely for the document "DESCRIBES" package relationship, returning that
// relationship or error if any is encountered. Arguments:
//   - packageName: name of package / directory
func BuildRelationshipSection2_3(packageName string) (*v2_3.Relationship, error) {
	rln := &v2_3.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", fmt.Sprintf("Package-%s", packageName)),
		Relationship: "DESCRIBES",
	}

	return rln, nil
}
