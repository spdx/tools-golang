// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) CreationInfo(ci *rdf2v1.CreationInfo) (id goraptor.Term, err error) {
	id = f.NodeId("cri")

	if err = f.setNodeType(id, rdf2v1.TypeCreationInfo); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"created", ci.Create.Val},
		Pair{"rdfs:comment", ci.Comment.Val},
		Pair{"licenseListVersion", ci.LicenseListVersion.Val},
	)

	if err != nil {
		return
	}

	for _, creator := range ci.Creator {
		if err = f.addLiteral(id, "creator", creator.Val); err != nil {
			return
		}
	}

	return id, nil
}
