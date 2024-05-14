package v3_0

import "time"

type ElementID string

// ----------------------- SpdxDocument -----------------------

type SpdxDocument struct {
	ElementCollection

	NamespaceMap map[string]string
}

type ISpdxDocument interface {
	IElementCollection
}

// ----------------------- ElementCollection -----------------------

type ElementCollection struct {
	Element

	RootElements []IElement
	Elements     []IElement
}

type IElementCollection interface {
	IElement

	AsElementCollection() *ElementCollection
}

// ----------------------- Element -----------------------

type Element struct {
	SpdxID ElementID
	Name   string
}

func (e *Element) AsElement() *Element {
	return e
}

type IElement interface {
	AsElement() *Element
}

var _ interface {
	IElement
} = (*Element)(nil)

// ----------------------- Artifact -----------------------

type Artifact struct {
	Element

	BuiltTime time.Time
}

func (a *Artifact) AsArtifact() *Artifact {
	return a
}

type IArtifact interface {
	IElement

	AsArtifact() *Artifact
}

var _ interface {
	IArtifact
} = (*Artifact)(nil)

// ----------------------- SoftwareArtifact -----------------------

type SoftwareArtifact struct {
	Artifact

	CopyrightText string
}

type ISoftwareArtifact interface {
	IArtifact

	AsSoftwareArtifact() *SoftwareArtifact
}

func (a *SoftwareArtifact) AsSoftwareArtifact() *SoftwareArtifact {
	return a
}

var _ interface {
	ISoftwareArtifact
} = (*SoftwareArtifact)(nil)

// ----------------------- Package -----------------------

type IPackage interface {
	ISoftwareArtifact

	AsPackage() *Package
}

type Package struct {
	SoftwareArtifact

	PackageVersion string
}

func (p *Package) AsPackage() *Package {
	return p
}

var _ interface {
	IPackage
} = (*Package)(nil)

// ----------------------- File -----------------------

type IFile interface {
	ISoftwareArtifact

	AsFile() *File
}

type File struct {
	SoftwareArtifact

	ContentType string
}

func (f *File) AsFile() *File {
	return f
}

var _ IPackage = (*Package)(nil)

// ----------------------- Relationship -----------------------

type RelationshipType string

type Relationship struct {
	Element

	RelationshipType RelationshipType
	From             IElement
	To               []IElement
}

type IRelationship interface {
	IElement

	AsRelationship() *Relationship
}

func (r *Relationship) AsRelationship() *Relationship {
	return r
}

var _ IRelationship = (*Relationship)(nil)
