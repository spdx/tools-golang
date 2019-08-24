// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Document struct {
	SPDXVersion            ValueStr
	DataLicense            ValueStr
	CreationInfo           *CreationInfo
	Review                 []*Review
	DocumentName           ValueStr
	DocumentNamespace      ValueStr
	SPDXID                 ValueStr
	DocumentComment        ValueStr
	ExtractedLicensingInfo []*ExtractedLicensingInfo
	Relationship           []*Relationship
	License                *License
	Annotation             []*Annotation
	ExternalDocumentRef    *ExternalDocumentRef
}

type ExternalDocumentRef struct {
	ExternalDocumentId ValueStr
	Checksum           *Checksum
	SPDXDocument       ValueStr
}

func (p *Parser) requestDocument(node goraptor.Term) (*Document, error) {
	obj, err := p.requestElementType(node, TypeDocument)
	if err != nil {
		return nil, err
	}
	return obj.(*Document), err
}

func (p *Parser) requestExternalDocumentRef(node goraptor.Term) (*ExternalDocumentRef, error) {
	obj, err := p.requestElementType(node, TypeExternalDocumentRef)
	if err != nil {
		return nil, err
	}
	return obj.(*ExternalDocumentRef), err
}

func (p *Parser) MapDocument(doc *Document) *builder {
	builder := &builder{t: TypeDocument, ptr: doc}
	doc.DocumentNamespace = DocumentNamespace
	doc.SPDXID = SPDXID
	builder.updaters = map[string]updater{
		"specVersion": update(&doc.SPDXVersion),
		// Example: gets CC0-1.0 from "http://spdx.org/licenses/CC0-1.0"
		"dataLicense": func(obj goraptor.Term) error {
			lic, err := p.requestLicense(obj)
			doc.License = lic
			return err
		},
		"creationInfo": func(obj goraptor.Term) error {
			ci, err := p.requestCreationInfo(obj)
			doc.CreationInfo = ci
			return err
		},
		"reviewed": func(obj goraptor.Term) error {
			rev, err := p.requestReview(obj)
			if err != nil {
				return err
			}
			doc.Review = append(doc.Review, rev)
			return err
		},
		"name":         update(&doc.DocumentName),
		"rdfs:comment": update(&doc.DocumentComment),
		"hasExtractedLicensingInfo": func(obj goraptor.Term) error {
			eli, err := p.requestExtractedLicensingInfo(obj)
			if err != nil {
				return err
			}
			doc.ExtractedLicensingInfo = append(doc.ExtractedLicensingInfo, eli)
			return nil
		},
		"relationship": func(obj goraptor.Term) error {
			rel, err := p.requestRelationship(obj)
			if err != nil {
				return err
			}
			doc.Relationship = append(doc.Relationship, rel)
			if rel != nil {
				DoctoRel[SPDXID] = append(DoctoRel[SPDXID], rel)
			}
			return nil
		},
		"annotation": func(obj goraptor.Term) error {
			an, err := p.requestAnnotation(obj)
			if err != nil {
				return err
			}
			doc.Annotation = append(doc.Annotation, an)
			if an != nil {
				DoctoAnno[SPDXID] = append(DoctoAnno[SPDXID], an)
			}
			return err
		},
		"externalDocumentRef": func(obj goraptor.Term) error {
			edr, err := p.requestExternalDocumentRef(obj)
			doc.ExternalDocumentRef = edr
			return err
		},
	}
	return builder
}

func (p *Parser) MapExternalDocumentRef(edr *ExternalDocumentRef) *builder {
	builder := &builder{t: TypeExternalDocumentRef, ptr: edr}
	builder.updaters = map[string]updater{
		"externalDocumentId": update(&edr.ExternalDocumentId),
		"checksum": func(obj goraptor.Term) error {
			cksum, err := p.requestChecksum(obj)

			edr.Checksum = cksum
			return err
		},
		"spdxDocument": update(&edr.SPDXDocument),
	}
	return builder

}
