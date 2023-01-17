// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// BuildRelationshipSection creates an SPDX Relationship
// solely for the document "DESCRIBES" package relationship, returning that
// relationship or error if any is encountered. Arguments:
//   - packageName: name of package / directory
func BuildRelationshipSection(packageName string) (*spdx.Relationship, error) {
	rln := &spdx.Relationship{
		RefA:         common.MakeDocElementID("", "DOCUMENT"),
		RefB:         common.MakeDocElementID("", fmt.Sprintf("Package-%s", packageName)),
		Relationship: "DESCRIBES",
	}

	return rln, nil
}
