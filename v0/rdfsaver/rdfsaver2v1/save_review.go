package rdfsaver2v1

import (
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"

	"github.com/deltamobile/goraptor"
)

func (f *Formatter) Review(r *rdf2v1.Review) (id goraptor.Term, err error) {
	id = f.NodeId("rev")

	if err = f.setNodeType(id, rdf2v1.TypeReview); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"reviewer", r.Reviewer.Val},
		Pair{"reviewDate", r.ReviewDate.Val},
		Pair{"rdfs:comment", r.ReviewComment.Val},
	)

	return id, err
}
func (f *Formatter) Reviews(parent goraptor.Term, element string, rs []*rdf2v1.Review) error {

	if len(rs) == 0 {
		return nil
	}

	for _, r := range rs {
		revId, err := f.Review(r)
		if err != nil {
			return err
		}
		if revId == nil {
			continue
		}
		if err = f.addTerm(parent, element, revId); err != nil {
			return err
		}
	}
	return nil
}
