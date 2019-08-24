// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Annotation struct {
	Annotator                ValueStr
	AnnotationType           ValueStr
	AnnotationDate           ValueStr
	AnnotationComment        ValueStr
	AnnotationSPDXIdentifier ValueStr
}

func (p *Parser) requestAnnotation(node goraptor.Term) (*Annotation, error) {
	obj, err := p.requestElementType(node, TypeAnnotation)
	if err != nil {
		return nil, err
	}
	return obj.(*Annotation), err
}
func (p *Parser) MapAnnotation(an *Annotation) *builder {
	builder := &builder{t: TypeAnnotation, ptr: an}
	builder.updaters = map[string]updater{
		"annotationDate": update(&an.AnnotationDate),
		"rdfs:comment":   update(&an.AnnotationComment),
		"annotator":      update(&an.Annotator),
		"annotationType": updateTrimPrefix(BaseUri, &an.AnnotationType),
	}
	return builder

}
