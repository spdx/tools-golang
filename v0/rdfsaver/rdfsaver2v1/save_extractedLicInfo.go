// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) ExtractedLicInfo(lic *rdf2v1.ExtractedLicensingInfo) (id goraptor.Term, err error) {
	id = f.NodeId("lic")

	if err = f.setNodeType(id, rdf2v1.TypeExtractedLicensingInfo); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"licenseId", lic.LicenseIdentifier.Val},
		Pair{"extractedText", lic.ExtractedText.Val},
		Pair{"rdfs:comment", lic.LicenseComment.Val},
	)

	if err != nil {
		return
	}

	for _, name := range lic.LicenseName {
		if err = f.addLiteral(id, "name", name.Val); err != nil {
			return
		}
	}

	for _, seealso := range lic.LicenseSeeAlso {
		if err = f.addLiteral(id, "rdfs:seeAlso", seealso.Val); err != nil {
			return
		}
	}

	return id, err
}

func (f *Formatter) ExtractedLicInfos(parent goraptor.Term, element string, lics []*rdf2v1.ExtractedLicensingInfo) error {

	if len(lics) == 0 {
		return nil
	}

	for _, lic := range lics {
		licId, err := f.ExtractedLicInfo(lic)
		if err != nil {
			return err
		}
		if licId == nil {
			continue
		}
		if err = f.addTerm(parent, element, licId); err != nil {
			return err
		}
	}
	return nil
}
