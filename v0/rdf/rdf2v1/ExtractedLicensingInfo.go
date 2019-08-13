package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type ExtractedLicensingInfo struct {
	LicenseIdentifier ValueStr
	LicenseName       []ValueStr
	ExtractedText     ValueStr
	LicenseComment    ValueStr
	LicenseSeeAlso    []ValueStr
}

func (p *Parser) requestExtractedLicensingInfo(node goraptor.Term) (*ExtractedLicensingInfo, error) {
	obj, err := p.requestElementType(node, typeExtractedLicensingInfo)
	if err != nil {
		return nil, err
	}
	return obj.(*ExtractedLicensingInfo), err
}

func (p *Parser) MapExtractedLicensingInfo(lic *ExtractedLicensingInfo) *builder {
	builder := &builder{t: typeExtractedLicensingInfo, ptr: lic}
	builder.updaters = map[string]updater{
		"licenseId":     update(&lic.LicenseIdentifier),
		"name":          updateList(&lic.LicenseName),
		"extractedText": update(&lic.ExtractedText),
		"rdfs:comment":  update(&lic.LicenseComment),
		"rdfs:seeAlso":  updateList(&lic.LicenseSeeAlso),
	}
	return builder
}
