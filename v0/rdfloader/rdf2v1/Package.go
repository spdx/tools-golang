// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Package struct {
	PackageName                  ValueStr
	PackageVersionInfo           ValueStr
	PackageFileName              ValueStr
	PackageSPDXIdentifier        ValueStr
	PackageDownloadLocation      ValueStr
	PackageVerificationCode      *PackageVerificationCode
	PackageComment               ValueStr
	PackageChecksum              *Checksum
	PackageLicense               *License
	PackageLicenseComments       ValueStr
	DisjunctiveLicenseSet        *DisjunctiveLicenseSet
	ConjunctiveLicenseSet        *ConjunctiveLicenseSet
	PackageLicenseInfoFromFiles  []ValueStr
	PackageLicenseDeclared       ValueStr
	PackageCopyrightText         ValueStr
	File                         []*File
	PackageRelationship          *Relationship
	PackageHomepage              ValueStr
	PackageSupplier              ValueStr
	PackageExternalRef           []*ExternalRef
	PackageOriginator            ValueStr
	PackageSourceInfo            ValueStr
	FilesAnalyzed                ValueStr
	PackageSummary               ValueStr
	PackageDescription           ValueStr
	Annotation                   []*Annotation
	PackageLicenseSPDXIdentifier ValueStr
}
type PackageVerificationCode struct {
	PackageVerificationCode             ValueStr
	PackageVerificationCodeExcludedFile ValueStr
}
type ExternalRef struct {
	ReferenceLocator  ValueStr
	ReferenceType     *ReferenceType
	ReferenceCategory ValueStr
	ReferenceComment  ValueStr
}
type ReferenceType struct {
	ReferenceType ValueStr
}

func (p *Parser) requestPackage(node goraptor.Term) (*Package, error) {
	obj, err := p.requestElementType(node, TypePackage)
	if err != nil {
		return nil, err
	}
	return obj.(*Package), err
}

func (p *Parser) requestPackageVerificationCode(node goraptor.Term) (*PackageVerificationCode, error) {
	obj, err := p.requestElementType(node, TypePackageVerificationCode)
	if err != nil {
		return nil, err
	}
	return obj.(*PackageVerificationCode), err
}

func (p *Parser) requestExternalRef(node goraptor.Term) (*ExternalRef, error) {
	obj, err := p.requestElementType(node, TypeExternalRef)
	if err != nil {
		return nil, err
	}
	return obj.(*ExternalRef), err
}
func (p *Parser) MapPackage(pkg *Package) *builder {
	builder := &builder{t: TypePackage, ptr: pkg}
	pkg.PackageSPDXIdentifier = SPDXIDPackage
	pkg.PackageLicenseSPDXIdentifier = SPDXIDLicense
	builder.updaters = map[string]updater{
		"name":             update(&pkg.PackageName),
		"versionInfo":      update(&pkg.PackageVersionInfo),
		"packageFileName":  update(&pkg.PackageFileName),
		"downloadLocation": update(&pkg.PackageDownloadLocation),
		"packageVerificationCode": func(obj goraptor.Term) error {
			pkgvc, err := p.requestPackageVerificationCode(obj)
			pkg.PackageVerificationCode = pkgvc
			return err
		},
		"checksum": func(obj goraptor.Term) error {
			pkgcksum, err := p.requestChecksum(obj)
			pkg.PackageChecksum = pkgcksum
			return err
		},
		"licenseComments": update(&pkg.PackageLicenseComments),
		"licenseConcluded": func(obj goraptor.Term) error {
			pkgdls, err := p.requestDisjunctiveLicenseSet(obj)
			pkg.DisjunctiveLicenseSet = pkgdls
			if err != nil {
				pkglic, err := p.requestLicense(obj)
				pkg.PackageLicense = pkglic
				if err != nil {
					pkgcls, err := p.requestConjunctiveLicenseSet(obj)
					pkg.ConjunctiveLicenseSet = pkgcls
					return err
				}
			}
			return nil
		},
		"licenseDeclared": func(obj goraptor.Term) error {
			_, ok := builder.updaters["http://spdx.org/rdf/terms#licenseDeclared"]
			if ok {
				builder.updaters = map[string]updater{"licenseDeclared": update(&pkg.PackageLicenseDeclared)}
			}
			pkgdls, _ := p.requestDisjunctiveLicenseSet(obj)
			pkg.DisjunctiveLicenseSet = pkgdls

			pkglic, _ := p.requestLicense(obj)
			pkg.PackageLicense = pkglic

			pkgcls, _ := p.requestConjunctiveLicenseSet(obj)
			pkg.ConjunctiveLicenseSet = pkgcls

			return nil
		},
		"licenseInfoFromFiles": updateList(&pkg.PackageLicenseInfoFromFiles),
		"copyrightText":        update(&pkg.PackageCopyrightText),
		"hasFile": func(obj goraptor.Term) error {
			file, err := p.requestFile(obj)

			// Relates File to Package
			if file != nil {
				PackagetoFile[SPDXIDPackage] = append(PackagetoFile[SPDXIDPackage], file)
			}
			if err != nil {
				return err
			}
			pkg.File = append(pkg.File, file)
			return nil
		},
		"relationship": func(obj goraptor.Term) error {
			rel, err := p.requestRelationship(obj)
			pkg.PackageRelationship = rel
			return err
		},
		"doap:homepage": update(&pkg.PackageHomepage),
		"supplier":      update(&pkg.PackageSupplier),
		"externalRef": func(obj goraptor.Term) error {
			er, err := p.requestExternalRef(obj)
			pkg.PackageExternalRef = append(pkg.PackageExternalRef, er)
			return err
		},
		"originator":    update(&pkg.PackageOriginator),
		"sourceInfo":    update((&pkg.PackageSourceInfo)),
		"summary":       update((&pkg.PackageSummary)),
		"filesAnalyzed": update((&pkg.FilesAnalyzed)),
		"description":   update((&pkg.PackageDescription)),
		"annotation": func(obj goraptor.Term) error {
			an, err := p.requestAnnotation(obj)
			pkg.Annotation = append(pkg.Annotation, an)
			if an != nil {
				PackagetoAnno[SPDXIDPackage] = append(PackagetoAnno[SPDXIDPackage], an)
			}
			return err
		},
		"rdfs:comment": update(&pkg.PackageComment)}
	return builder
}

func (p *Parser) MapPackageVerificationCode(pkgvc *PackageVerificationCode) *builder {
	builder := &builder{t: TypePackageVerificationCode, ptr: pkgvc}
	builder.updaters = map[string]updater{
		"packageVerificationCodeValue":        update(&pkgvc.PackageVerificationCode),
		"packageVerificationCodeExcludedFile": update(&pkgvc.PackageVerificationCodeExcludedFile),
	}
	return builder
}

func (p *Parser) MapExternalRef(er *ExternalRef) *builder {
	builder := &builder{t: TypeExternalRef, ptr: er}
	builder.updaters = map[string]updater{
		"referenceLocator":  update(&er.ReferenceLocator),
		"referenceCategory": update(&er.ReferenceCategory),
		"rdfs:comment":      update(&er.ReferenceComment),
		"referenceType": func(obj goraptor.Term) error {
			rt, err := p.requestReferenceType(obj)
			er.ReferenceType = rt
			return err
		},
	}
	return builder
}
