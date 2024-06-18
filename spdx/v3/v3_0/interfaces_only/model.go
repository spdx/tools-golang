package v3_0

import "time"

type ElementID string

type SpdxDocument interface {
	ElementCollection
}

type spdxDocument struct {
	// Element properties
	spdxID ElementID
	name   string

	// ElementCollection properties
	profileConformance ProfileIdentifierType
	rootElement        []Element
	elements           []Element
}

func (p *spdxDocument) SpdxID() ElementID {
	return p.spdxID
}

func (p *spdxDocument) SetSpdxID(id ElementID) error {
	p.spdxID = id
	return nil
}

func (p *spdxDocument) Name() string {
	return p.name
}

func (p *spdxDocument) SetName(name string) error {
	p.name = name
	return nil
}

func (p *spdxDocument) ProfileConformance() ProfileIdentifierType {
	return p.profileConformance
}

func (p *spdxDocument) SetProfileConformance(profileConformance ProfileIdentifierType) error {
	p.profileConformance = profileConformance
	return nil
}

func (e *spdxDocument) RootElement() []Element {
	return e.rootElement
}

func (e *spdxDocument) SetRootElement(rootElement []Element) error {
	e.rootElement = rootElement
	return nil
}

func (e *spdxDocument) Elements() []Element {
	return e.elements
}

func (e *spdxDocument) SetElements(elements []Element) error {
	e.elements = elements
	return nil
}

func NewSpdxDocument() SpdxDocument {
	return &spdxDocument{}
}

// ----------------------- Element -----------------------

type Element interface {
	SpdxID() ElementID
	SetSpdxID(spdxID ElementID) error

	Name() string
	SetName(name string) error
}

// ----------------------- Artifact -----------------------

type Artifact interface {
	Element

	BuiltTime() time.Time
	SetBuiltTime(builtTime time.Time) error
}

// ----------------------- SoftwareArtifact -----------------------

type SoftwareArtifact interface {
	Artifact

	CopyrightText() string
	SetCopyrightText(copyrightText string) error
}

// ----------------------- Package -----------------------

type Package interface {
	SoftwareArtifact

	PackageVersion() string
	SetPackageVersion(packageVersion string) error
}

// "package" is a reserved word...
type packageImpl struct {
	// Element properties
	spdxID ElementID
	name   string

	// Artifact properties
	builtTime time.Time

	// SoftwareArtifact properties
	copyrightText string

	// Package properties
	packageVersion string
}

func (p *packageImpl) SpdxID() ElementID {
	return p.spdxID
}

func (p *packageImpl) SetSpdxID(id ElementID) error {
	p.spdxID = id
	return nil
}

func (p *packageImpl) Name() string {
	return p.name
}

func (p *packageImpl) SetName(name string) error {
	p.name = name
	return nil
}

func (p *packageImpl) BuiltTime() time.Time {
	return p.builtTime
}

func (p *packageImpl) SetBuiltTime(builtTime time.Time) error {
	p.builtTime = builtTime
	return nil
}

func (p *packageImpl) CopyrightText() string {
	return p.copyrightText
}

func (p *packageImpl) SetCopyrightText(copyrightText string) error {
	p.copyrightText = copyrightText
	return nil
}

func (p *packageImpl) PackageVersion() string {
	return p.packageVersion
}

func (p *packageImpl) SetPackageVersion(packageVersion string) error {
	p.packageVersion = packageVersion
	return nil
}

func NewPackage() Package {
	return &packageImpl{}
}

var _ Package = (*packageImpl)(nil)

// ----------------------- File -----------------------

type File interface {
	SoftwareArtifact

	ContentType() string
	SetContentType(contentType string) error
}

type file struct {
	// Element properties
	spdxID ElementID
	name   string

	// Artifact properties
	builtTime time.Time

	// SoftwareArtifact properties
	copyrightText string

	// File properties
	contentType string
}

func (p *file) SpdxID() ElementID {
	return p.spdxID
}

func (p *file) SetSpdxID(id ElementID) error {
	p.spdxID = id
	return nil
}

func (p *file) Name() string {
	return p.name
}

func (p *file) SetName(name string) error {
	p.name = name
	return nil
}

func (p *file) BuiltTime() time.Time {
	return p.builtTime
}

func (p *file) SetBuiltTime(builtTime time.Time) error {
	p.builtTime = builtTime
	return nil
}

func (p *file) CopyrightText() string {
	return p.copyrightText
}

func (p *file) SetCopyrightText(copyrightText string) error {
	p.copyrightText = copyrightText
	return nil
}

func (p *file) ContentType() string {
	return p.contentType
}

func (p *file) SetContentType(contentType string) error {
	p.contentType = contentType
	return nil
}

func NewFile() File {
	return &file{}
}

var _ File = (*file)(nil)

// ----------------------- ElementCollection -----------------------

type ProfileIdentifierType string

type ElementCollection interface {
	Element

	ProfileConformance() ProfileIdentifierType
	SetProfileConformance(profileConformance ProfileIdentifierType) error

	RootElement() []Element
	SetRootElement(element []Element) error

	Elements() []Element
	SetElements(elements []Element) error
}

// ----------------------- Relationship -----------------------

type RelationshipType string

type Relationship interface {
	Element

	RelationshipType() RelationshipType
	SetRelationshipType(relationshipType RelationshipType) error

	From() Element
	SetFrom(element Element) error

	To() []Element
	SetTo(element []Element) error
}

func NewRelationship() Relationship {
	return &relationship{}
}

type relationship struct {
	// Element properties
	spdxID ElementID
	name   string

	// Relationship properties
	relationshipType RelationshipType
	from             Element
	to               []Element
}

func (p *relationship) SpdxID() ElementID {
	return p.spdxID
}

func (p *relationship) SetSpdxID(id ElementID) error {
	p.spdxID = id
	return nil
}

func (p *relationship) Name() string {
	return p.name
}

func (p *relationship) SetName(name string) error {
	p.name = name
	return nil
}

func (e *relationship) RelationshipType() RelationshipType {
	return e.relationshipType
}

func (e *relationship) SetRelationshipType(relationshipType RelationshipType) error {
	e.relationshipType = relationshipType
	return nil
}

func (e *relationship) From() Element {
	return e.from
}

func (e *relationship) SetFrom(from Element) error {
	e.from = from
	return nil
}

func (e *relationship) To() []Element {
	return e.to
}

func (e *relationship) SetTo(to []Element) error {
	e.to = to
	return nil
}
