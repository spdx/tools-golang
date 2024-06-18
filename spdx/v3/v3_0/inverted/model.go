package v3_0

import "time"

type ElementID string

type SpdxDocument struct {
}

func NewSpdxDocument(p SpdxDocumentProps) *Element {
	return &Element{
		SpdxID:       p.SpdxID,
		Name:         p.Name,
		SpdxDocument: &SpdxDocument{},
		ElementCollection: &ElementCollection{
			RootElement: p.RootElement,
			Elements:    p.Elements,
		},
	}
}

type SpdxDocumentProps struct {
	SpdxID      ElementID
	Name        string
	RootElement []*Element
	Elements    []*Element
}

// ----------------------- Element -----------------------

type Element struct {
	*ElementCollection
	*SpdxDocument
	*Artifact
	*SoftwareArtifact
	*Package
	*File
	*Relationship

	SpdxID ElementID
	Name   string
}

// ----------------------- Artifact -----------------------

type Artifact struct {
	BuiltTime time.Time
}

// ----------------------- SoftwareArtifact -----------------------

type SoftwareArtifact struct {
	CopyrightText string
}

// ----------------------- Package -----------------------

type Package struct {
	PackageVersion string
}

func NewPackage(p PackageProps) *Element {
	return &Element{
		SpdxID:           p.SpdxID,
		Name:             p.Name,
		Artifact:         &Artifact{},
		SoftwareArtifact: &SoftwareArtifact{},
		Package: &Package{
			PackageVersion: p.PackageVersion,
		},
	}
}

type PackageProps struct {
	SpdxID         ElementID
	Name           string
	PackageVersion string
}

// ----------------------- File -----------------------

type File struct {
	ContentType string
}

func NewFile(p FileProps) *Element {
	return &Element{
		SpdxID:           p.SpdxID,
		Name:             p.Name,
		Artifact:         &Artifact{},
		SoftwareArtifact: &SoftwareArtifact{},
		File: &File{
			ContentType: p.ContentType,
		},
	}
}

type FileProps struct {
	SpdxID      ElementID
	Name        string
	ContentType string
}

// ----------------------- Relationship -----------------------

type Relationship struct {
	RelationshipType string
	From             *Element
	To               []*Element
}

func NewRelationship(p RelationshipProps) *Element {
	return &Element{
		SpdxID: p.SpdxID,
		Name:   p.Name,
		Relationship: &Relationship{
			RelationshipType: p.RelationshipType,
			From:             p.From,
			To:               p.To,
		},
	}
}

type RelationshipProps struct {
	SpdxID           ElementID
	Name             string
	RelationshipType string
	From             *Element
	To               []*Element
}

// ----------------------- ElementCollection -----------------------

type ElementCollection struct {
	RootElement []*Element
	Elements    []*Element
}
