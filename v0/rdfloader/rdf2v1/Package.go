package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Package struct {
	PackageName                 ValueStr
	PackageVersionInfo          ValueStr
	PackageFileName             ValueStr
	PackageDownloadLocation     ValueStr
	PackageVerificationCode     *PackageVerificationCode
	PackageComment              ValueStr
	PackageChecksum             *Checksum
	PackageLicense              *License
	PackageLicenseComments      ValueStr
	DisjunctiveLicenseSet       *DisjunctiveLicenseSet
	ConjunctiveLicenseSet       *ConjunctiveLicenseSet
	PackageLicenseInfoFromFiles []ValueStr
	PackageLicenseDeclared      ValueStr
	PackageCopyrightText        ValueStr
	File                        []*File
	PackageRelationship         *Relationship
	PackageHomepage             ValueStr
	PackageSupplier             ValueStr
	PackageExternalRef          *ExternalRef
	PackageOriginator           ValueStr
	PackageSourceInfo           ValueStr
	FilesAnalyzed               ValueStr
	PackageSummary              ValueStr
	PackageDescription          ValueStr
	Annotation                  []*Annotation
}
type PackageVerificationCode struct {
	PackageVerificationCode             ValueStr
	PackageVerificationCodeExcludedFile ValueStr
}
type PackageRelationship struct {
	Relationshiptype   ValueStr
	relatedSpdxElement ValueStr
}

func (p *Parser) requestPackage(node goraptor.Term) (*Package, error) {
	obj, err := p.requestElementType(node, typePackage)
	if err != nil {
		return nil, err
	}
	return obj.(*Package), err
}
func (p *Parser) requestPackageRelationship(node goraptor.Term) (*PackageRelationship, error) {
	obj, err := p.requestElementType(node, typeRelationship)
	if err != nil {
		return nil, err
	}
	return obj.(*PackageRelationship), err
}

func (p *Parser) requestPackageVerificationCode(node goraptor.Term) (*PackageVerificationCode, error) {
	obj, err := p.requestElementType(node, typePackageVerificationCode)
	if err != nil {
		return nil, err
	}
	return obj.(*PackageVerificationCode), err
}

func (p *Parser) MapPackage(pkg *Package) *builder {
	builder := &builder{t: typePackage, ptr: pkg}
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
			pkgdls, _ := p.requestDisjunctiveLicenseSet(obj)
			pkg.DisjunctiveLicenseSet = pkgdls
			pkglic, _ := p.requestLicense(obj)
			pkg.PackageLicense = pkglic
			pkgcls, _ := p.requestConjunctiveLicenseSet(obj)
			pkg.ConjunctiveLicenseSet = pkgcls
			return nil
		},
		"licenseDeclared":      update(&pkg.PackageLicenseDeclared),
		"licenseInfoFromFiles": updateList(&pkg.PackageLicenseInfoFromFiles),
		"copyrightText":        update(&pkg.PackageCopyrightText),
		"hasFile": func(obj goraptor.Term) error {
			file, err := p.requestFile(obj)
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
			pkg.PackageExternalRef = er
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
			return err
		},
		"rdfs:comment": update(&pkg.PackageComment)}
	return builder
}

func (p *Parser) MapPackageRelationship(pkgrel *PackageRelationship) *builder {
	builder := &builder{t: typePackageVerificationCode, ptr: pkgrel}
	builder.updaters = map[string]updater{
		"relationshipType":   update(&pkgrel.Relationshiptype),
		"relatedSpdxElement": update(&pkgrel.Relationshiptype),
	}
	return builder
}

func (p *Parser) MapPackageVerificationCode(pkgvc *PackageVerificationCode) *builder {
	builder := &builder{t: typePackageVerificationCode, ptr: pkgvc}
	builder.updaters = map[string]updater{
		"packageVerificationCodeValue":        update(&pkgvc.PackageVerificationCode),
		"packageVerificationCodeExcludedFile": update(&pkgvc.PackageVerificationCodeExcludedFile),
	}
	return builder
}
