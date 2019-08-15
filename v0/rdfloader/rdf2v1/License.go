package rdf2v1

import "github.com/deltamobile/goraptor"

type License struct {
	LicenseComment                ValueStr
	LicenseName                   ValueStr
	LicenseText                   ValueStr
	StandardLicenseHeader         ValueStr
	LicenseSeeAlso                []ValueStr
	LicenseIsFsLibre              ValueStr
	StandardLicenseTemplate       ValueStr
	StandardLicenseHeaderTemplate ValueStr
	LicenseId                     ValueStr
	LicenseisOsiApproved          ValueStr
}
type DisjunctiveLicenseSet struct {
	Member []ValueStr
}
type ConjunctiveLicenseSet struct {
	License                *License
	ExtractedLicensingInfo *ExtractedLicensingInfo
}

func (p *Parser) requestLicense(node goraptor.Term) (*License, error) {
	obj, err := p.requestElementType(node, typeLicense)
	if err != nil {
		return nil, err
	}
	return obj.(*License), err
}
func (p *Parser) requestDisjunctiveLicenseSet(node goraptor.Term) (*DisjunctiveLicenseSet, error) {
	obj, err := p.requestElementType(node, typeDisjunctiveLicenseSet)
	if err != nil {
		return nil, err
	}
	return obj.(*DisjunctiveLicenseSet), err
}
func (p *Parser) requestConjunctiveLicenseSet(node goraptor.Term) (*ConjunctiveLicenseSet, error) {
	obj, err := p.requestElementType(node, typeConjunctiveLicenseSet)
	if err != nil {
		return nil, err
	}
	return obj.(*ConjunctiveLicenseSet), err
}
func (p *Parser) MapLicense(lic *License) *builder {
	builder := &builder{t: typeLicense, ptr: lic}
	builder.updaters = map[string]updater{
		"rdfs:comment":                  update(&lic.LicenseComment),
		"name":                          update(&lic.LicenseName),
		"licenseText":                   update(&lic.LicenseText),
		"licenseId":                     update(&lic.LicenseId),
		"rdfs:seeAlso":                  updateList(&lic.LicenseSeeAlso),
		"isFsfLibre":                    update(&lic.LicenseIsFsLibre),
		"isOsiApproved":                 update(&lic.LicenseisOsiApproved),
		"standardLicenseHeader":         update(&lic.StandardLicenseHeader),
		"standardLicenseTemplate":       update(&lic.StandardLicenseTemplate),
		"standardLicenseHeaderTemplate": update(&lic.StandardLicenseTemplate),
	}
	return builder
}

func (p *Parser) MapDisjunctiveLicenseSet(dls *DisjunctiveLicenseSet) *builder {
	builder := &builder{t: typeDisjunctiveLicenseSet, ptr: dls}
	builder.updaters = map[string]updater{
		"member": updateList(&dls.Member),
	}
	return builder
}
func (p *Parser) MapConjunctiveLicenseSet(cls *ConjunctiveLicenseSet) *builder {
	builder := &builder{t: typeConjunctiveLicenseSet, ptr: cls}
	builder.updaters = map[string]updater{
		"member": func(obj goraptor.Term) error {

			lic, err := p.requestLicense(obj)
			cls.License = lic
			if err != nil {
				eli, err := p.requestExtractedLicensingInfo(obj)
				cls.ExtractedLicensingInfo = eli
				return err
			}
			return nil
		},
	}
	return builder
}
