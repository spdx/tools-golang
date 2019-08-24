// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Review struct {
	ReviewComment ValueStr
	ReviewDate    ValueStr
	Reviewer      ValueStr
}

func (p *Parser) requestReview(node goraptor.Term) (*Review, error) {
	obj, err := p.requestElementType(node, TypeReview)
	if err != nil {
		return nil, err
	}
	return obj.(*Review), err
}
func (p *Parser) MapReview(rev *Review) *builder {
	builder := &builder{t: TypeReview, ptr: rev}
	builder.updaters = map[string]updater{
		"rdfs:comment": update(&rev.ReviewComment),
		"reviewDate":   update(&rev.ReviewDate),
		"reviewer":     update(&rev.Reviewer),
	}
	return builder

}
