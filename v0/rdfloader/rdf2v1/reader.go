// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"fmt"
	"strings"

	"github.com/deltamobile/goraptor"
)

var (
	URInsType = Uri("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	TypeDocument                = Prefix("SpdxDocument")
	TypeCreationInfo            = Prefix("CreationInfo")
	TypeExtractedLicensingInfo  = Prefix("ExtractedLicensingInfo")
	TypeRelationship            = Prefix("Relationship")
	TypePackage                 = Prefix("Package")
	TypePackageVerificationCode = Prefix("PackageVerificationCode")
	TypeChecksum                = Prefix("Checksum")
	TypeDisjunctiveLicenseSet   = Prefix("DisjunctiveLicenseSet")
	TypeConjunctiveLicenseSet   = Prefix("ConjunctiveLicenseSet")
	TypeFile                    = Prefix("File")
	TypeSpdxElement             = Prefix("SpdxElement")
	TypeSnippet                 = Prefix("Snippet")
	TypeLicenseConcluded        = Prefix("licenseConcluded")
	TypeReview                  = Prefix("Review")
	TypeAnnotation              = Prefix("Annotation")
	TypeLicense                 = Prefix("License")
	TypeExternalDocumentRef     = Prefix("ExternalDocumentRef")
	TypeExternalRef             = Prefix("ExternalRef")
	TypeProject                 = Prefix("doap:Project")
	TypeReferenceType           = Prefix("ReferenceType")
	TypeSnippetStartEndPointer  = Prefix("j.0:StartEndPointer")
	TypeByteOffsetPointer       = Prefix("j.0:ByteOffsetPointer")
	TypeLineCharPointer         = Prefix("j.0:LineCharPointer")
)
var (
	DocumentNamespace  ValueStr
	ProjectURI         ValueStr
	SPDXID             ValueStr
	SPDXIDFile         ValueStr
	SPDXIDRelationship ValueStr
	SPDXIDSnippet      ValueStr
	SPDXIDPackage      ValueStr
	SPDXIDLicense      ValueStr
	SPDXIDCLicense     ValueStr
)

var (
	counter       int
	PackagetoFile = make(map[ValueStr][]*File)
	ReltoPackage  = make(map[ValueStr][]*Package)
	ReltoFile     = make(map[ValueStr][]*File)
	DoctoRel      = make(map[ValueStr][]*Relationship)
	SniptoFile    = make(map[ValueStr]*File)
	DoctoAnno     = make(map[ValueStr][]*Annotation)
	FiletoAnno    = make(map[ValueStr][]*Annotation)
	PackagetoAnno = make(map[ValueStr][]*Annotation)
)

// Parser Struct and associated methods
type Parser struct {
	Rdfparser *goraptor.Parser
	Input     string
	Index     map[string]*builder
	Buffer    map[string][]*goraptor.Statement
	Doc       *Document
	Snip      *Snippet
}

// NewParser uses goraptor.NewParser to initialse a new parser interface
func NewParser(input string) *Parser {

	return &Parser{
		Rdfparser: goraptor.NewParser("guess"),
		Input:     input,
		Index:     make(map[string]*builder),
		Buffer:    make(map[string][]*goraptor.Statement),
	}
}

func (p *Parser) Parse() (*Document, *Snippet, error) {
	// PARSE FILE method - Takes the file location as an input
	ch := p.Rdfparser.ParseFile(p.Input, "")
	var err error

	for statement := range ch {
		if err = p.ProcessTriple(statement); err != nil {
			break
		}
	}
	return p.Doc, p.Snip, err
}

// Free the goraptor parser.
func (p *Parser) Free() {
	p.Rdfparser.Free()
	p.Snip = nil
	p.Doc = nil
}

func (p *Parser) ProcessTriple(stm *goraptor.Statement) error {
	node := termStr(stm.Subject)
	ns, id, _ := ExtractNs(node)
	if id == "SPDXRef-DOCUMENT" {
		SPDXID = Str(id)
		if DocumentNamespace.Val == "" {
			DocumentNamespace = Str(ns)
		}
	}

	if ExtractId(termStr(stm.Predicate)) == "member" {
		SPDXIDCLicense = Str(ExtractId(termStr(stm.Object)))
		if SPDXIDLicense == Str("") {
			SPDXIDCLicense = Str(strings.Replace(termStr(stm.Object), "http://spdx.org/licenses/", "", 1))
		}
	}
	if ExtractId(termStr(stm.Predicate)) == "relationshipType" {
		SPDXIDRelationship = Str(strings.Replace(termStr(stm.Object), "http://spdx.org/rdf/terms#relationshipType_", "", 1))
		counter++
	}

	if stm.Predicate.Equals(URInsType) {
		_, err := p.setNodeType(stm.Subject, stm.Object)
		return err
	}

	// apply function if it's a builder
	builder, ok := p.Index[node]
	if ok {
		return builder.apply(stm.Predicate, stm.Object)
	}

	// buffer statement
	if _, ok := p.Buffer[node]; !ok {
		p.Buffer[node] = make([]*goraptor.Statement, 0)
	}
	p.Buffer[node] = append(p.Buffer[node], stm)
	return nil
}

func (p *Parser) setNodeType(node, t goraptor.Term) (interface{}, error) {
	nodeStr := termStr(node)
	builder, ok := p.Index[nodeStr]
	if ExtractId(termStr(t)) == "File" {
		SPDXIDFile = Str(ExtractId(termStr(node)))
	}
	if ExtractId(termStr(t)) == "Package" {
		SPDXIDPackage = Str(ExtractId(termStr(node)))
	}
	if ExtractId(termStr(t)) == "Snippet" {
		SPDXIDSnippet = Str(ExtractId(termStr(node)))
	}
	if ExtractId(termStr(t)) == "License" {
		SPDXIDLicense = Str(ExtractId(termStr(node)))
	}
	if ExtractId(termStr(t)) == "Project" {
		ProjectURI = Str(termStr(node))
	}

	if ok {
		if !checkRaptorTypes(builder.t, t) && builder.checkPredicate("ns:type") {
			if err := builder.apply(Uri("ns:type"), t); err != nil {
				return nil, err
			}
			return builder.ptr, nil
		}
		if !checkCompatibleTypes(builder.t, t) {
			return nil, fmt.Errorf("IncompatibleType")
		}
		return builder.ptr, nil
	}

	// new builder by type
	switch {
	// t is goraptor Object
	case t.Equals(TypeDocument):
		p.Doc = new(Document)
		builder = p.MapDocument(p.Doc)

	case t.Equals(TypeCreationInfo):
		builder = p.MapCreationInfo(new(CreationInfo))

	case t.Equals(TypeExtractedLicensingInfo):
		builder = p.MapExtractedLicensingInfo(new(ExtractedLicensingInfo))

	case t.Equals(TypeRelationship):
		builder = p.MapRelationship(new(Relationship))

	case t.Equals(TypePackage):
		builder = p.MapPackage(new(Package))

	case t.Equals(TypePackageVerificationCode):
		builder = p.MapPackageVerificationCode(new(PackageVerificationCode))

	case t.Equals(TypeChecksum):
		builder = p.MapChecksum(new(Checksum))

	case t.Equals(TypeDisjunctiveLicenseSet):
		builder = p.MapDisjunctiveLicenseSet(new(DisjunctiveLicenseSet))

	case t.Equals(TypeFile):
		builder = p.MapFile(new(File))

	case t.Equals(TypeReview):
		builder = p.MapReview(new(Review))

	case t.Equals(TypeLicense):
		builder = p.MapLicense(new(License))

	case t.Equals(TypeAnnotation):
		builder = p.MapAnnotation(new(Annotation))

	case t.Equals(TypeExternalRef):
		builder = p.MapExternalRef(new(ExternalRef))

	case t.Equals(TypeReferenceType):
		builder = p.MapReferenceType(new(ReferenceType))

	case t.Equals(TypeExternalDocumentRef):
		builder = p.MapExternalDocumentRef(new(ExternalDocumentRef))

	case t.Equals(TypeProject):
		builder = p.MapProject(new(Project))

	case t.Equals(TypeSnippet):
		p.Snip = new(Snippet)
		builder = p.MapSnippet(p.Snip)

	case t.Equals(TypeSpdxElement):
		builder = p.MapSpdxElement(new(SpdxElement))

	case t.Equals(TypeConjunctiveLicenseSet):
		builder = p.MapConjunctiveLicenseSet(new(ConjunctiveLicenseSet))

	case t.Equals(TypeSnippetStartEndPointer):
		builder = p.MapSnippetStartEndPointer(new(SnippetStartEndPointer))

	case t.Equals(TypeLineCharPointer):
		builder = p.MapLineCharPointer(new(LineCharPointer))

	case t.Equals(TypeByteOffsetPointer):
		builder = p.MapByteOffsetPointer(new(ByteOffsetPointer))
	default:
		return nil, fmt.Errorf("New Builder: Types does not match.")
	}

	p.Index[nodeStr] = builder

	buf := p.Buffer[nodeStr]
	for _, stm := range buf {
		if err := builder.apply(stm.Predicate, stm.Object); err != nil {
			return nil, err
		}
	}
	delete(p.Buffer, nodeStr)

	return builder.ptr, nil
}

func checkRaptorTypes(found goraptor.Term, need ...goraptor.Term) bool {
	for _, b := range need {
		if found == b || found.Equals(b) {
			return true
		}
	}
	return false
}

func checkCompatibleTypes(input, required goraptor.Term) bool {
	if checkRaptorTypes(input, required) {
		return true
	}
	return false
}

func (p *Parser) requestElementType(node, t goraptor.Term) (interface{}, error) {
	builder, ok := p.Index[termStr(node)]
	if ok {
		if !checkCompatibleTypes(builder.t, t) {
			return nil, fmt.Errorf("%v and %v are Incompatible Type", builder.t, t)
		}
		return builder.ptr, nil
	}
	return p.setNodeType(node, t)
}

// Builder Struct and associated methods
type builder struct {
	t        goraptor.Term // type of element this builder represents
	ptr      interface{}   // the spdx element that this builder builds
	updaters map[string]updater
}

func (b *builder) apply(pred, obj goraptor.Term) error {
	property := ShortPrefix(pred)
	f, ok := b.updaters[property]

	if !ok {
		return fmt.Errorf("Property %s is not supported for %s.", property, b.t)
	}
	return f(obj)
}

// to check if builder contains a predicate
func (b *builder) checkPredicate(pred string) bool {
	_, ok := b.updaters[pred]
	return ok
}
