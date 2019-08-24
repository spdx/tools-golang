// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"github.com/spdx/tools-golang/v0/spdx"

	"github.com/deltamobile/goraptor"
)

type CreationInfo struct {
	SPDXIdentifier     ValueStr
	LicenseListVersion ValueStr
	Creator            []ValueStr
	Create             ValueStr
	Comment            ValueStr
}

func (p *Parser) requestCreationInfo(node goraptor.Term) (*CreationInfo, error) {

	obj, err := p.requestElementType(node, TypeCreationInfo)
	if err != nil {
		return nil, err
	}
	return obj.(*CreationInfo), err
}

func (p *Parser) MapCreationInfo(ci *CreationInfo) *builder {
	builder := &builder{t: TypeCreationInfo, ptr: ci}
	ci.SPDXIdentifier = SPDXID
	builder.updaters = map[string]updater{
		"licenseListVersion": update(&ci.LicenseListVersion),
		"creator":            updateList(&ci.Creator),
		"created":            update(&ci.Create),
		"rdfs:comment":       update(&ci.Comment),
	}
	return builder
}

func ExtractCreator(ci *CreationInfo, creator string) []string {

	var val []string
	for _, c := range ValueList(ci.Creator) {
		subkey, subvalue, _ := extractSubs(c)
		if subkey == creator {
			val = append(val, subvalue)
		}
	}
	return val
}
func InsertCreator(ci *spdx.CreationInfo2_1) []ValueStr {

	var val []string
	if len(ci.CreatorPersons) != 0 {
		for _, person := range ci.CreatorPersons {
			val = append(val, person)
		}
	}
	if len(ci.CreatorOrganizations) != 0 {
		for _, org := range ci.CreatorOrganizations {
			val = append(val, org)
		}
	}
	if len(ci.CreatorOrganizations) != 0 {
		for _, org := range ci.CreatorOrganizations {
			val = append(val, org)
		}
	}
	valstr := ValueStrList(val)
	return valstr
}
