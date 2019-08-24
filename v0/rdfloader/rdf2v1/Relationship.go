// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Relationship struct {
	RelationshipType    ValueStr
	Package             []*Package
	File                []*File
	RelatedSpdxElement  ValueStr
	SpdxElement         *SpdxElement
	RelationshipComment ValueStr
}
type SpdxElement struct {
	SpdxElement ValueStr
}

func (p *Parser) requestRelationship(node goraptor.Term) (*Relationship, error) {
	obj, err := p.requestElementType(node, TypeRelationship)
	if err != nil {
		return nil, err
	}
	return obj.(*Relationship), err
}
func (p *Parser) requestSpdxElement(node goraptor.Term) (*SpdxElement, error) {
	obj, err := p.requestElementType(node, TypeSpdxElement)
	if err != nil {
		return nil, err
	}
	return obj.(*SpdxElement), err
}

func (p *Parser) MapRelationship(rel *Relationship) *builder {
	builder := &builder{t: TypeRelationship, ptr: rel}
	builder.updaters = map[string]updater{
		"relationshipType": update(&rel.RelationshipType),
		"rdfs:comment":     update(&rel.RelationshipComment),
		"relatedSpdxElement": func(obj goraptor.Term) error {
			_, ok := builder.updaters["http://spdx.org/rdf/terms#relatedSpdxElement"]
			if ok {
				builder.updaters = map[string]updater{"relatedSpdxElement": update(&rel.RelatedSpdxElement)}
				return nil
			}
			pkg, err := p.requestPackage(obj)
			rel.Package = append(rel.Package, pkg)

			// Relates Relationship to Package
			if pkg != nil {
				ReltoPackage[SPDXIDRelationship] = append(ReltoPackage[SPDXIDRelationship], pkg)
			}
			if err != nil {
				file, err := p.requestFile(obj)
				rel.File = append(rel.File, file)
				if file != nil {
					ReltoFile[SPDXIDRelationship] = append(ReltoFile[SPDXIDRelationship], file)
				}
				if err != nil {
					se, err := p.requestSpdxElement(obj)
					rel.SpdxElement = se
					return err
				}

			}
			return nil
		},
	}
	return builder
}

func (p *Parser) MapSpdxElement(se *SpdxElement) *builder {
	builder := &builder{t: TypeSpdxElement, ptr: se}
	builder.updaters = map[string]updater{}
	return builder
}
