package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Snippet struct {
	SnippetName             ValueStr
	SnippetCopyrightText    ValueStr
	SnippetLicenseComments  ValueStr
	SnippetFromFile         *File
	SnippetStartEndPointer  []*SnippetStartEndPointer
	SnippetLicenseConcluded ValueStr
	SnippetComment          ValueStr
	LicenseInfoInSnippet    []ValueStr
	SnippetSPDXIdentifier   ValueStr
}

type SnippetStartEndPointer struct {
	ByteOffsetPointer []*ByteOffsetPointer
	LineCharPointer   []*LineCharPointer
}

type ByteOffsetPointer struct {
	Reference ValueStr
	Offset    ValueStr
}

type LineCharPointer struct {
	Reference  ValueStr
	LineNumber ValueStr
}

func (p *Parser) requestSnippet(node goraptor.Term) (*Snippet, error) {
	obj, err := p.requestElementType(node, TypeSnippet)
	if err != nil {
		return nil, err
	}
	return obj.(*Snippet), err
}

func (p *Parser) requestReferenceType(node goraptor.Term) (*ReferenceType, error) {
	obj, err := p.requestElementType(node, TypeReferenceType)
	if err != nil {
		return nil, err
	}
	return obj.(*ReferenceType), err
}
func (p *Parser) requestSnippetStartEndPointer(node goraptor.Term) (*SnippetStartEndPointer, error) {
	obj, err := p.requestElementType(node, TypeSnippetStartEndPointer)
	if err != nil {
		return nil, err
	}
	return obj.(*SnippetStartEndPointer), err
}

func (p *Parser) requestByteOffsetPointer(node goraptor.Term) (*ByteOffsetPointer, error) {
	obj, err := p.requestElementType(node, TypeByteOffsetPointer)
	if err != nil {
		return nil, err
	}
	return obj.(*ByteOffsetPointer), err
}
func (p *Parser) requestLineCharPointer(node goraptor.Term) (*LineCharPointer, error) {
	obj, err := p.requestElementType(node, TypeLineCharPointer)
	if err != nil {
		return nil, err
	}
	return obj.(*LineCharPointer), err
}
func (p *Parser) MapSnippet(s *Snippet) *builder {
	builder := &builder{t: TypeSnippet, ptr: s}
	s.SnippetSPDXIdentifier = SPDXIDSnippet
	builder.updaters = map[string]updater{
		"name":            update(&s.SnippetName),
		"copyrightText":   update(&s.SnippetCopyrightText),
		"licenseComments": update(&s.SnippetLicenseComments),
		"snippetFromFile": func(obj goraptor.Term) error {
			file, err := p.requestFile(obj)
			s.SnippetFromFile = file
			return err
		},
		"licenseInfoInSnippet": updateList(&s.LicenseInfoInSnippet),
		"rdfs:comment":         update(&s.SnippetComment),
		"licenseConcluded":     update(&s.SnippetLicenseConcluded),
		"range": func(obj goraptor.Term) error {
			sep, err := p.requestSnippetStartEndPointer(obj)
			s.SnippetStartEndPointer = append(s.SnippetStartEndPointer, sep)
			return err
		},
	}
	return builder
}

func (p *Parser) MapReferenceType(rt *ReferenceType) *builder {
	builder := &builder{t: TypeReferenceType, ptr: rt}
	builder.updaters = map[string]updater{}
	return builder
}

func (p *Parser) MapSnippetStartEndPointer(sep *SnippetStartEndPointer) *builder {
	builder := &builder{t: TypeSnippetStartEndPointer, ptr: sep}
	builder.updaters = map[string]updater{
		"j.0:startPointer": func(obj goraptor.Term) error {
			lc, err := p.requestLineCharPointer(obj)
			sep.LineCharPointer = append(sep.LineCharPointer, lc)
			if err != nil {
				bo, err := p.requestByteOffsetPointer(obj)
				sep.ByteOffsetPointer = append(sep.ByteOffsetPointer, bo)
				return err
			}
			return nil
		},
		"j.0:endPointer": func(obj goraptor.Term) error {
			lc, err := p.requestLineCharPointer(obj)
			sep.LineCharPointer = append(sep.LineCharPointer, lc)
			if err != nil {
				bo, err := p.requestByteOffsetPointer(obj)
				sep.ByteOffsetPointer = append(sep.ByteOffsetPointer, bo)
				return err
			}
			return nil
		},
	}
	return builder
}

func (p *Parser) MapLineCharPointer(lc *LineCharPointer) *builder {
	builder := &builder{t: TypeLineCharPointer, ptr: lc}
	builder.updaters = map[string]updater{
		"j.0:reference":  update(&lc.Reference),
		"j.0:lineNumber": update(&lc.LineNumber),
	}
	return builder
}
func (p *Parser) MapByteOffsetPointer(bo *ByteOffsetPointer) *builder {
	builder := &builder{t: TypeByteOffsetPointer, ptr: bo}
	builder.updaters = map[string]updater{
		"j.0:reference": update(&bo.Reference),
		"j.0:offset":    update(&bo.Offset),
	}
	return builder
}
