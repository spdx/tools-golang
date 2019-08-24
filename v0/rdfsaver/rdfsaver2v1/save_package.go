// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) Package(pkg *rdf2v1.Package) (id goraptor.Term, err error) {
	id = f.NodeId("pkg")

	if err = f.setNodeType(id, rdf2v1.TypePackage); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"name", pkg.PackageName.Val},
		Pair{"versionInfo", pkg.PackageVersionInfo.Val},
		Pair{"packageFileName", pkg.PackageFileName.Val},
		Pair{"downloadLocation", pkg.PackageDownloadLocation.Val},
		Pair{"rdfs:comment", pkg.PackageComment.Val},
		Pair{"licenseComments", pkg.PackageLicenseComments.Val},
		Pair{"copyrightText", pkg.PackageCopyrightText.Val},
		Pair{"doap:homepage", pkg.PackageHomepage.Val},
		Pair{"supplier", pkg.PackageSupplier.Val},
		Pair{"originator", pkg.PackageOriginator.V()},
		Pair{"sourceInfo", pkg.PackageSourceInfo.Val},
		Pair{"filesAnalyzed", pkg.FilesAnalyzed.Val},
		Pair{"summary", pkg.PackageSummary.Val},
		Pair{"description", pkg.PackageDescription.Val},
	)
	if err != nil {
		return
	}
	if pkg.PackageVerificationCode != nil {
		pkgid, err := f.PackageVerificationCode(pkg.PackageVerificationCode)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "packageVerificationCode", pkgid); err != nil {
			return id, err
		}
	}

	if pkg.PackageChecksum != nil {
		cksumId, err := f.Checksum(pkg.PackageChecksum)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "checksum", cksumId); err != nil {
			return id, err
		}
	}

	if err = f.Annotations(id, "annotation", pkg.Annotation); err != nil {
		return
	}

	if err = f.Files(id, "hasFile", pkg.File); err != nil {
		return id, err
	}

	for _, per := range pkg.PackageExternalRef {
		if pkg.PackageExternalRef != nil {
			pkgErId, err := f.ExternalRef(per)
			if err != nil {
				return id, err
			}
			if err = f.addTerm(id, "externalRef", pkgErId); err != nil {
				return id, err
			}
		}
	}

	if pkg.PackageRelationship != nil {
		pkgRelId, err := f.Relationship(pkg.PackageRelationship)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "relationship", pkgRelId); err != nil {
			return id, err
		}
	}
	for _, lif := range pkg.PackageLicenseInfoFromFiles {
		if err = f.addTerm(id, "licenseInfoFromFiles", rdf2v1.Prefix(lif.Val)); err != nil {
			return
		}
	}

	if pkg.PackageLicense != nil {
		pkglicId, err := f.License(pkg.PackageLicense)
		if err != nil {
			pkglicId, err = f.DisjunctiveLicenseSet(pkg.DisjunctiveLicenseSet)
			if err != nil {
				pkglicId, err = f.ConjunctiveLicenseSet(pkg.ConjunctiveLicenseSet)
				if err != nil {
					return id, err
				}
			}
		}
		if err = f.addTerm(id, "licenseConcluded", pkglicId); err != nil {
			return id, err

		}
	}

	if pkg.PackageLicenseDeclared.Val != "" {
		if err = f.addTerm(id, "licenseDeclared", rdf2v1.Prefix(pkg.PackageLicenseDeclared.Val)); err != nil {
			pkglicId, err := f.License(pkg.PackageLicense)
			if err != nil {
				pkglicId, err = f.DisjunctiveLicenseSet(pkg.DisjunctiveLicenseSet)
				if err != nil {
					pkglicId, err = f.ConjunctiveLicenseSet(pkg.ConjunctiveLicenseSet)
					if err != nil {
						return id, err
					}
				}
			}
			if err = f.addTerm(id, "licenseDeclared", pkglicId); err != nil {
				return id, err
			}
		}

	}
	return id, err
}

func (f *Formatter) Packages(parent goraptor.Term, element string, pkgs []*rdf2v1.Package) error {
	if len(pkgs) == 0 {
		return nil
	}
	for _, pkg := range pkgs {
		pkgid, err := f.Package(pkg)
		if err != nil {
			return err
		}
		if err = f.addTerm(parent, element, pkgid); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) PackageVerificationCode(pvc *rdf2v1.PackageVerificationCode) (id goraptor.Term, err error) {
	id = f.NodeId("pvc")

	if err = f.setNodeType(id, rdf2v1.TypePackageVerificationCode); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"packageVerificationCodeValue", pvc.PackageVerificationCode.Val},
		Pair{"packageVerificationCodeExcludedFile", pvc.PackageVerificationCodeExcludedFile.Val},
	)

	return id, err
}
func (f *Formatter) ExternalRef(er *rdf2v1.ExternalRef) (id goraptor.Term, err error) {
	id = f.NodeId("er")

	if err = f.setNodeType(id, rdf2v1.TypeExternalRef); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"referenceLocator", er.ReferenceLocator.Val},
		Pair{"rdfs:comment", er.ReferenceComment.Val},
	)
	if id, err := f.ReferenceType(er.ReferenceType); err == nil {
		if err = f.addTerm(id, "referenceType", id); err != nil {
			return id, err
		}
	}
	if er.ReferenceCategory.Val != "" {
		if err = f.addTerm(id, "referenceCategory", rdf2v1.Prefix(er.ReferenceCategory.Val)); err != nil {
			return
		}
	}
	return id, err
}

func (f *Formatter) ReferenceType(rt *rdf2v1.ReferenceType) (id goraptor.Term, err error) {
	id = f.NodeId("rt")

	if err = f.setNodeType(id, rdf2v1.TypeReferenceType); err != nil {
		return
	}

	if rt.ReferenceType.Val != "" {
		if err = f.addTerm(id, "referenceType", rdf2v1.Prefix(rt.ReferenceType.Val)); err != nil {
			return
		}
	}
	return id, err
}
