package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Relationship struct {
	RelationshipType    ValueStr
	Package             []*Package
	File                []*File
	relatedSpdxElement  ValueStr
	SpdxElement         *SpdxElement
	RelationshipComment ValueStr
}
type SpdxElement struct {
	SpdxElement ValueStr
}

func (p *Parser) requestRelationship(node goraptor.Term) (*Relationship, error) {
	obj, err := p.requestElementType(node, typeRelationship)
	if err != nil {
		return nil, err
	}
	return obj.(*Relationship), err
}
func (p *Parser) requestSpdxElement(node goraptor.Term) (*SpdxElement, error) {
	obj, err := p.requestElementType(node, typeSpdxElement)
	if err != nil {
		return nil, err
	}
	return obj.(*SpdxElement), err
}

func (p *Parser) MapRelationship(rel *Relationship) *builder {
	builder := &builder{t: typeRelationship, ptr: rel}
	builder.updaters = map[string]updater{
		"relationshipType": update(&rel.RelationshipType),
		"rdfs:comment":     update(&rel.RelationshipComment),
		"relatedSpdxElement": func(obj goraptor.Term) error {
			_, ok := builder.updaters["http://spdx.org/rdf/terms#relatedSpdxElement"]
			if ok {
				builder.updaters = map[string]updater{"relatedSpdxElement": update(&rel.relatedSpdxElement)}
				return nil
			}
			pkg, err := p.requestPackage(obj)
			rel.Package = append(rel.Package, pkg)
			if err != nil {
				file, err := p.requestFile(obj)
				rel.File = append(rel.File, file)
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
	builder := &builder{t: typeSpdxElement, ptr: se}
	builder.updaters = map[string]updater{}
	return builder
}
