// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"strings"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) Checksum(cksum *rdf2v1.Checksum) (id goraptor.Term, err error) {
	id = f.NodeId("cksum")

	if err = f.setNodeType(id, rdf2v1.TypeChecksum); err != nil {
		return
	}

	err = f.addLiteral(id, "checksumValue", cksum.ChecksumValue.Val)
	if err != nil {
		return
	}

	algo := strings.ToLower(cksum.Algorithm.Val)
	if algo == "sha1" {
		err = f.addTerm(id, "algorithm", rdf2v1.Prefix("checksumAlgorithm_sha1"))

	} else if algo == "md5" {
		err = f.addTerm(id, "algorithm", rdf2v1.Prefix("checksumAlgorithm_md5"))
	} else if algo == "sha256" {
		err = f.addTerm(id, "algorithm", rdf2v1.Prefix("checksumAlgorithm_sha256"))
	} else {
		err = f.addLiteral(id, "algorithm", algo)
	}

	return id, err
}
