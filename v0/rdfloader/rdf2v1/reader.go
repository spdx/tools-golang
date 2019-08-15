package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

var (
	URInsType = uri("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	typeDocument                = prefix("SpdxDocument")
	typeCreationInfo            = prefix("CreationInfo")
	typeExtractedLicensingInfo  = prefix("ExtractedLicensingInfo")
	typeRelationship            = prefix("Relationship")
	typePackage                 = prefix("Package")
	typePackageVerificationCode = prefix("PackageVerificationCode")
	typeChecksum                = prefix("Checksum")
	typeDisjunctiveLicenseSet   = prefix("DisjunctiveLicenseSet")
	typeConjunctiveLicenseSet   = prefix("ConjunctiveLicenseSet")
	typeFile                    = prefix("File")
	typeSpdxElement             = prefix("SpdxElement")
	typeSnippet                 = prefix("Snippet")
	typeLicenseConcluded        = prefix("licenseConcluded")
	typeReview                  = prefix("Review")
	typeAnnotation              = prefix("Annotation")
	typeLicense                 = prefix("License")
	typeExternalDocumentRef     = prefix("ExternalDocumentRef")
	typeExternalRef             = prefix("ExternalRef")
	typeProject                 = prefix("doap:Project")
	typeReferenceType           = prefix("ReferenceType")
	typeSnippetStartEndPointer  = prefix("j.0:StartEndPointer")
	typeByteOffsetPointer       = prefix("j.0:ByteOffsetPointer")
	typeLineCharPointer         = prefix("j.0:LineCharPointer")
)
var DocumentNamespace ValueStr

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

// Process parsed goraptor statements
func (p *Parser) ProcessTriple(stm *goraptor.Statement) error {
	node := termStr(stm.Subject)
	ns, id, _ := ExtractNs(node)
	if id == "SPDXRef-DOCUMENT" {
		if DocumentNamespace.Val == "" {
			DocumentNamespace = Str(ns)
		}
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

// checks compatiblity of terms, returns builder by type, runs buffer
func (p *Parser) setNodeType(node, t goraptor.Term) (interface{}, error) {
	nodeStr := termStr(node)
	builder, ok := p.Index[nodeStr]

	if ok {
		if !checkRaptorTypes(builder.t, t) && builder.checkPredicate("ns:type") {

			if err := builder.apply(uri("ns:type"), t); err != nil {
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
	case t.Equals(typeDocument):
		p.Doc = new(Document)
		builder = p.MapDocument(p.Doc)

	case t.Equals(typeCreationInfo):
		builder = p.MapCreationInfo(new(CreationInfo))

	case t.Equals(typeExtractedLicensingInfo):
		builder = p.MapExtractedLicensingInfo(new(ExtractedLicensingInfo))

	case t.Equals(typeRelationship):
		builder = p.MapRelationship(new(Relationship))

	case t.Equals(typePackage):
		builder = p.MapPackage(new(Package))

	case t.Equals(typePackageVerificationCode):
		builder = p.MapPackageVerificationCode(new(PackageVerificationCode))

	case t.Equals(typeChecksum):
		builder = p.MapChecksum(new(Checksum))

	case t.Equals(typeDisjunctiveLicenseSet):
		builder = p.MapDisjunctiveLicenseSet(new(DisjunctiveLicenseSet))

	case t.Equals(typeFile):
		builder = p.MapFile(new(File))

	case t.Equals(typeReview):
		builder = p.MapReview(new(Review))

	case t.Equals(typeLicense):
		builder = p.MapLicense(new(License))

	case t.Equals(typeAnnotation):
		builder = p.MapAnnotation(new(Annotation))

	case t.Equals(typeExternalRef):
		builder = p.MapExternalRef(new(ExternalRef))

	case t.Equals(typeReferenceType):
		builder = p.MapReferenceType(new(ReferenceType))

	case t.Equals(typeExternalDocumentRef):
		builder = p.MapExternalDocumentRef(new(ExternalDocumentRef))

	case t.Equals(typeProject):
		builder = p.MapProject(new(Project))

	case t.Equals(typeSnippet):
		p.Snip = new(Snippet)
		builder = p.MapSnippet(p.Snip)

	case t.Equals(typeSpdxElement):
		builder = p.MapSpdxElement(new(SpdxElement))

	case t.Equals(typeConjunctiveLicenseSet):
		builder = p.MapConjunctiveLicenseSet(new(ConjunctiveLicenseSet))

	case t.Equals(typeSnippetStartEndPointer):
		builder = p.MapSnippetStartEndPointer(new(SnippetStartEndPointer))

	case t.Equals(typeLineCharPointer):
		builder = p.MapLineCharPointer(new(LineCharPointer))

	case t.Equals(typeByteOffsetPointer):
		builder = p.MapByteOffsetPointer(new(ByteOffsetPointer))
	default:
		return nil, fmt.Errorf("New Builder: Types does not match.")
	}

	p.Index[nodeStr] = builder

	// run buffer
	buf := p.Buffer[nodeStr]
	for _, stm := range buf {
		if err := builder.apply(stm.Predicate, stm.Object); err != nil {
			return nil, err
		}
	}
	delete(p.Buffer, nodeStr)

	return builder.ptr, nil
}

// compares goraptor terms
func checkRaptorTypes(found goraptor.Term, need ...goraptor.Term) bool {
	for _, b := range need {
		if found == b || found.Equals(b) {
			return true
		}
	}
	return false
}

// check compatiblity of goraptor term with input
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
	property := shortPrefix(pred)
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
