package v3_0_e

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
	SetSpdxID(spdxID ElementID) error

	Name() string
	SetName(name string) error
}

type ElementImpl struct {
	spdxID ElementID
	name   string
}

var _ Element = (*ElementImpl)(nil)

func (p *ElementImpl) SpdxID() ElementID {
	return p.spdxID
}

func (p *ElementImpl) SetSpdxID(id ElementID) error {
	p.spdxID = id
	return nil
}

func (p *ElementImpl) Name() string {
	return p.name
}

func (p *ElementImpl) SetName(name string) error {
	p.name = name
	return nil
}

// ----------------------- Artifact -----------------------

type Artifact interface {
	Element

	BuiltTime() time.Time
	SetBuiltTime(builtTime time.Time) error
}

type ArtifactImpl struct {
	ElementImpl

	builtTime time.Time
}

var _ Artifact = (*ArtifactImpl)(nil)
var _ Element = (*ArtifactImpl)(nil)

func (p *ArtifactImpl) BuiltTime() time.Time {
	return p.builtTime
}

func (p *ArtifactImpl) SetBuiltTime(builtTime time.Time) error {
	p.builtTime = builtTime
	return nil
}

// ----------------------- SoftwareArtifact -----------------------

type SoftwareArtifact interface {
	Artifact

	CopyrightText() string
	SetCopyrightText(copyrightText string) error
}

type SoftwareArtifactImpl struct {
	ArtifactImpl

	copyrightText string
}

var _ SoftwareArtifact = (*SoftwareArtifactImpl)(nil)
var _ Artifact = (*SoftwareArtifactImpl)(nil)
var _ Element = (*SoftwareArtifactImpl)(nil)

func (p *SoftwareArtifactImpl) CopyrightText() string {
	return p.copyrightText
}

func (p *SoftwareArtifactImpl) SetCopyrightText(copyrightText string) error {
	p.copyrightText = copyrightText
	return nil
}

// ----------------------- Package -----------------------

type Package interface {
	SoftwareArtifact

	PackageVersion() string
	SetPackageVersion(packageVersion string) error
}

var _ Element = (Package)(nil)

type PackageImpl struct {
	SoftwareArtifactImpl

	packageVersion string
}

var _ Package = (*PackageImpl)(nil)
var _ SoftwareArtifact = (*PackageImpl)(nil)
var _ Artifact = (*PackageImpl)(nil)
var _ Element = (*PackageImpl)(nil)

func (p *PackageImpl) PackageVersion() string {
	return p.packageVersion
}

func (p *PackageImpl) SetPackageVersion(packageVersion string) error {
	p.packageVersion = packageVersion
	return nil
}

// ----------------------- File -----------------------

type File interface {
	SoftwareArtifact

	ContentType() string
	SetContentType(contentType string) error
}

type FileImpl struct {
	SoftwareArtifactImpl

	// File properties
	contentType string
}

var _ File = (*FileImpl)(nil)
var _ SoftwareArtifact = (*FileImpl)(nil)
var _ Artifact = (*FileImpl)(nil)
var _ Element = (*FileImpl)(nil)

func (p *FileImpl) ContentType() string {
	return p.contentType
}

func (p *FileImpl) SetContentType(contentType string) error {
	p.contentType = contentType
	return nil
}

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

func (p *ElementCollectionImpl) SetProfileConformance(profileConformance ProfileIdentifierType) error {
	p.profileConformance = profileConformance
	return nil
}

func (e *ElementCollectionImpl) RootElement() []Element {
	return e.rootElement
}

func (e *ElementCollectionImpl) SetRootElement(rootElement []Element) error {
	e.rootElement = rootElement
	return nil
}

func (e *ElementCollectionImpl) Elements() []Element {
	return e.elements
}

func (e *ElementCollectionImpl) SetElements(elements []Element) error {
	e.elements = elements
	return nil
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

type RelationshipImpl struct {
	ElementImpl

	relationshipType RelationshipType
	from             Element
	to               []Element
}

var _ Relationship = (*RelationshipImpl)(nil)
var _ Element = (*RelationshipImpl)(nil)

func (e *RelationshipImpl) RelationshipType() RelationshipType {
	return e.relationshipType
}

func (e *RelationshipImpl) SetRelationshipType(relationshipType RelationshipType) error {
	e.relationshipType = relationshipType
	return nil
}

func (e *RelationshipImpl) From() Element {
	return e.from
}

func (e *RelationshipImpl) SetFrom(from Element) error {
	e.from = from
	return nil
}

func (e *RelationshipImpl) To() []Element {
	return e.to
}

func (e *RelationshipImpl) SetTo(to []Element) error {
	e.to = to
	return nil
}
