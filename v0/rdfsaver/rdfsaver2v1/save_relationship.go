package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) Relationship(rel *rdf2v1.Relationship) (id goraptor.Term, err error) {
	id = f.NodeId("rel")

	if err = f.setNodeType(id, rdf2v1.TypeRelationship); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"rdfs:comment", rel.RelationshipComment.Val},
	)
	if err != nil {
		return
	}

	if rel.RelationshipType.Val != "" {
		if err = f.addTerm(id, "relationshipType", rdf2v1.Prefix(rel.RelationshipType.Val)); err != nil {
			return
		}
	}
	if rel.RelatedSpdxElement.Val != "" {
		if err = f.addTerm(id, "relatedSpdxElement", rdf2v1.Prefix(rel.RelatedSpdxElement.Val)); err != nil {
			return
		}
	}
	if rel.SpdxElement != nil {
		seId, err := f.SpdxElement(rel.SpdxElement)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "relatedSpdxElement", seId); err != nil {
			return id, err
		}
	}

	if err = f.Files(id, "relatedSpdxElement", rel.File); err != nil {
		if err = f.Packages(id, "relatedSpdxElement", rel.Package); err != nil {
			return
		}
	}

	return id, err
}

func (f *Formatter) Relationships(parent goraptor.Term, element string, rels []*rdf2v1.Relationship) error {
	if len(rels) == 0 {
		return nil
	}
	for _, rel := range rels {
		relId, err := f.Relationship(rel)
		if err != nil {
			return err
		}
		if relId == nil {
			continue
		}
		if err = f.addTerm(parent, element, relId); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) SpdxElement(se *rdf2v1.SpdxElement) (id goraptor.Term, err error) {
	id = f.NodeId("se")

	if err = f.setNodeType(id, rdf2v1.TypeSpdxElement); err != nil {
		return
	}
	if se.SpdxElement.Val != "" {
		if err = f.addTerm(id, "SpdxElement", rdf2v1.Prefix(se.SpdxElement.Val)); err != nil {
			return
		}
	}
	return id, err
}
