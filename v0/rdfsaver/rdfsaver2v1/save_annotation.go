// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver2v1

import (
	"github.com/deltamobile/goraptor"
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
)

func (f *Formatter) Annotation(an *rdf2v1.Annotation) (id goraptor.Term, err error) {
	id = f.NodeId("an")

	if err = f.setNodeType(id, rdf2v1.TypeAnnotation); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"annotationDate", an.AnnotationDate.Val},
		Pair{"rdfs:comment", an.AnnotationComment.Val},
		Pair{"annotator", an.Annotator.Val},
	)
	if err != nil {
		return
	}

	if an.AnnotationType.Val != "" {
		if err = f.addTerm(id, "annotationType", rdf2v1.Prefix(an.AnnotationType.Val)); err != nil {
			return
		}
	}
	return id, err
}
func (f *Formatter) Annotations(parent goraptor.Term, element string, ans []*rdf2v1.Annotation) error {

	if len(ans) == 0 {
		return nil
	}

	for _, an := range ans {
		annId, err := f.Annotation(an)
		if err != nil {
			return err
		}
		if annId == nil {
			continue
		}
		if err = f.addTerm(parent, element, annId); err != nil {
			return err
		}
	}
	return nil
}
