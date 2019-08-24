// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) License(lic *rdf2v1.License) (id goraptor.Term, err error) {
	id = f.NodeId("lic")

	if err = f.setNodeType(id, rdf2v1.TypeLicense); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"rdfs:comment", lic.LicenseComment.Val},
		Pair{"name", lic.LicenseName.Val},
		Pair{"licenseText", lic.LicenseText.Val},
		Pair{"standardLicenseHeader", lic.StandardLicenseHeader.V()},
		Pair{"standardLicenseTemplate", lic.StandardLicenseTemplate.V()},
		Pair{"standardLicenseHeaderTemplate", lic.StandardLicenseHeaderTemplate.Val},
		Pair{"isFsfLibre", lic.LicenseIsFsLibre.Val},
		Pair{"licenseId", lic.LicenseId.Val},
		Pair{"licenseOsiApproved", lic.LicenseisOsiApproved.Val},
	)
	for _, sa := range lic.LicenseSeeAlso {
		if err = f.addLiteral(id, "rdfs:seeAlso", sa.Val); err != nil {
			return
		}
	}

	return id, err
}
func (f *Formatter) ConjunctiveLicenseSet(cls *rdf2v1.ConjunctiveLicenseSet) (id goraptor.Term, err error) {
	id = f.NodeId("cls")

	if err = f.setNodeType(id, rdf2v1.TypeConjunctiveLicenseSet); err != nil {
		return
	}

	if id, err := f.License(cls.License); err == nil {
		if err = f.addTerm(id, "member", id); err != nil {
			return id, err
		}
	} else if id, err := f.ExtractedLicInfo(cls.ExtractedLicensingInfo); err == nil {
		if err = f.addTerm(id, "member", id); err != nil {
			return id, err
		}
	} else {
		return id, err
	}

	return id, err
}

func (f *Formatter) DisjunctiveLicenseSet(dls *rdf2v1.DisjunctiveLicenseSet) (id goraptor.Term, err error) {
	id = f.NodeId("dls")

	if err = f.setNodeType(id, rdf2v1.TypeDisjunctiveLicenseSet); err != nil {
		return
	}

	for _, mem := range dls.Member {
		if err = f.addLiteral(id, "member", mem.Val); err != nil {
			return
		}
	}

	return id, err
}
