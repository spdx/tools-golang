package v3_0

import "time"

type ElementID string

// ----------------------- SpdxDocumentData -----------------------

type SpdxDocumentData struct {
	ElementCollectionData

	NamespaceMap map[string]string
}

type SpdxDocument interface {
	ElementCollection
}

// ----------------------- ElementCollectionData -----------------------

type ElementCollectionData struct {
	ElementData

	RootElements []Element
	Elements     []Element
}

type ElementCollection interface {
	Element

	AsElementCollection() *ElementCollectionData
}

// ----------------------- ElementData -----------------------

type ElementData struct {
	SpdxID ElementID
	Name   string
}

func (e *ElementData) AsElement() *ElementData {
	return e
}

type Element interface {
	AsElement() *ElementData
}

var _ interface {
	Element
} = (*ElementData)(nil)

// ----------------------- ArtifactData -----------------------

type ArtifactData struct {
	ElementData

	BuiltTime time.Time
}

func (a *ArtifactData) AsArtifact() *ArtifactData {
	return a
}

type Artifact interface {
	Element

	AsArtifact() *ArtifactData
}

var _ interface {
	Artifact
} = (*ArtifactData)(nil)

// ----------------------- SoftwareArtifactData -----------------------

type SoftwareArtifactData struct {
	ArtifactData

	CopyrightText string
}

type ISoftwareArtifact interface {
	Artifact

	AsSoftwareArtifact() *SoftwareArtifactData
}

func (a *SoftwareArtifactData) AsSoftwareArtifact() *SoftwareArtifactData {
	return a
}

var _ interface {
	ISoftwareArtifact
} = (*SoftwareArtifactData)(nil)

// ----------------------- PackageData -----------------------

type Package interface {
	ISoftwareArtifact

	AsPackage() *PackageData
}

type PackageData struct {
	SoftwareArtifactData

	PackageVersion string
}

func (p *PackageData) AsPackage() *PackageData {
	return p
}

var _ interface {
	Package
} = (*PackageData)(nil)

// ----------------------- FileData -----------------------

type File interface {
	ISoftwareArtifact

	AsFile() *FileData
}

type FileData struct {
	SoftwareArtifactData

	ContentType string
}

func (f *FileData) AsFile() *FileData {
	return f
}

var _ Package = (*PackageData)(nil)

// ----------------------- RelationshipData -----------------------

type RelationshipType string

type RelationshipData struct {
	ElementData

	RelationshipType RelationshipType
	From             Element
	To               []Element
}

type Relationship interface {
	Element

	AsRelationship() *RelationshipData
}

func (r *RelationshipData) AsRelationship() *RelationshipData {
	return r
}

var _ Relationship = (*RelationshipData)(nil)
