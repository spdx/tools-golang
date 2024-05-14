package v3_0

import "time"

type ElementID string

type SpdxDocument interface {
	ElementCollection
}

type SpdxDocumentImpl struct {
	ElementCollectionImpl
}

var _ SpdxDocument = (*SpdxDocumentImpl)(nil)

// ----------------------- Element -----------------------

type Element interface {
	SpdxID() ElementID
	SetSpdxID(spdxID ElementID)

	Name() string
	SetName(name string)
}

type ElementImpl struct {
	spdxID ElementID
	name   string
}

var _ Element = (*ElementImpl)(nil)

func (p *ElementImpl) SpdxID() ElementID {
	return p.spdxID
}

func (p *ElementImpl) SetSpdxID(id ElementID) {
	p.spdxID = id
}

func (p *ElementImpl) Name() string {
	return p.name
}

func (p *ElementImpl) SetName(name string) {
	p.name = name
}

// ----------------------- Artifact -----------------------

type Artifact interface {
	Element

	BuiltTime() time.Time
	SetBuiltTime(builtTime time.Time)
}

type ArtifactImpl struct {
	ElementImpl

	builtTime time.Time
}

var _ Artifact = (*ArtifactImpl)(nil)

func (p *ArtifactImpl) BuiltTime() time.Time {
	return p.builtTime
}

func (p *ArtifactImpl) SetBuiltTime(builtTime time.Time) {
	p.builtTime = builtTime
}

// ----------------------- SoftwareArtifact -----------------------

type SoftwareArtifact interface {
	Artifact

	CopyrightText() string
	SetCopyrightText(copyrightText string)
}

type SoftwareArtifactImpl struct {
	ArtifactImpl

	copyrightText string
}

var _ SoftwareArtifact = (*SoftwareArtifactImpl)(nil)

func (p *SoftwareArtifactImpl) CopyrightText() string {
	return p.copyrightText
}

func (p *SoftwareArtifactImpl) SetCopyrightText(copyrightText string) {
	p.copyrightText = copyrightText
}

// ----------------------- Package -----------------------

type Package interface {
	SoftwareArtifact

	PackageVersion() string
	SetPackageVersion(packageVersion string)
}

var _ Element = (Package)(nil)

type PackageImpl struct {
	SoftwareArtifactImpl

	packageVersion string
}

var _ Package = (*PackageImpl)(nil)

func (p *PackageImpl) PackageVersion() string {
	return p.packageVersion
}

func (p *PackageImpl) SetPackageVersion(packageVersion string) {
	p.packageVersion = packageVersion
}

// ----------------------- File -----------------------

type File interface {
	SoftwareArtifact

	ContentType() string
	SetContentType(contentType string)
}

type FileImpl struct {
	SoftwareArtifactImpl

	// File properties
	contentType string
}

var _ File = (*FileImpl)(nil)

func (p *FileImpl) ContentType() string {
	return p.contentType
}

func (p *FileImpl) SetContentType(contentType string) {
	p.contentType = contentType
}

// ----------------------- ElementCollection -----------------------

type ProfileIdentifierType string

type ElementCollection interface {
	Element

	ProfileConformance() ProfileIdentifierType
	SetProfileConformance(profileConformance ProfileIdentifierType)

	RootElement() []Element
	SetRootElement(element []Element)

	Elements() []Element
	SetElements(elements []Element)
}

type ElementCollectionImpl struct {
	ElementImpl

	// ElementCollection properties
	profileConformance ProfileIdentifierType
	rootElement        []Element
	elements           []Element
}

var _ ElementCollection = (*ElementCollectionImpl)(nil)
var _ Element = (*ElementCollectionImpl)(nil)

func (p *ElementCollectionImpl) ProfileConformance() ProfileIdentifierType {
	return p.profileConformance
}

func (p *ElementCollectionImpl) SetProfileConformance(profileConformance ProfileIdentifierType) {
	p.profileConformance = profileConformance
}

func (e *ElementCollectionImpl) RootElement() []Element {
	return e.rootElement
}

func (e *ElementCollectionImpl) SetRootElement(rootElement []Element) {
	e.rootElement = rootElement
}

func (e *ElementCollectionImpl) Elements() []Element {
	return e.elements
}

func (e *ElementCollectionImpl) SetElements(elements []Element) {
	e.elements = elements
}

// ----------------------- Relationship -----------------------

type RelationshipType string

type Relationship interface {
	Element

	RelationshipType() RelationshipType
	SetRelationshipType(relationshipType RelationshipType)

	From() Element
	SetFrom(element Element)

	To() []Element
	SetTo(element []Element)
}

type RelationshipImpl struct {
	*ElementImpl

	relationshipType RelationshipType
	from             Element
	to               []Element
}

var _ Relationship = (*RelationshipImpl)(nil)
var _ Element = (*RelationshipImpl)(nil)

func (e *RelationshipImpl) RelationshipType() RelationshipType {
	return e.relationshipType
}

func (e *RelationshipImpl) SetRelationshipType(relationshipType RelationshipType) {
	e.relationshipType = relationshipType
}

func (e *RelationshipImpl) From() Element {
	return e.from
}

func (e *RelationshipImpl) SetFrom(from Element) {
	e.from = from
}

func (e *RelationshipImpl) To() []Element {
	return e.to
}

func (e *RelationshipImpl) SetTo(to []Element) {
	e.to = to
}
