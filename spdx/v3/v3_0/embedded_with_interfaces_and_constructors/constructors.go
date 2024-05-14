package v3_0

import (
	"time"
)

// Could add convenience creation props like:

func NewPackage(props PackageProps) Package {
	out := &PackageImpl{}
	out.SetSpdxID(props.SpdxID)
	out.SetName(props.Name)
	out.SetBuiltTime(props.BuiltTime)
	out.SetCopyrightText(props.CopyrightText)
	out.SetPackageVersion(props.PackageVersion)
	return out
}

type PackageProps struct {
	// Element properties
	SpdxID ElementID
	Name   string

	// Artifact properties
	BuiltTime time.Time

	// SoftwareArtifact properties
	CopyrightText string

	// Package properties
	PackageVersion string
}

func NewFile(props FileProps) File {
	out := &FileImpl{}
	out.SetSpdxID(props.SpdxID)
	out.SetName(props.Name)
	out.SetBuiltTime(props.BuiltTime)
	out.SetCopyrightText(props.CopyrightText)
	out.SetContentType(props.ContentType)
	return out
}

type FileProps struct {
	// Element properties
	SpdxID ElementID
	Name   string

	// Artifact properties
	BuiltTime time.Time

	// SoftwareArtifact properties
	CopyrightText string

	// File properties
	ContentType string
}

func NewRelationship(props RelationshipProps) Relationship {
	out := &RelationshipImpl{}
	out.SetSpdxID(props.SpdxID)
	out.SetName(props.Name)
	out.SetRelationshipType(props.RelationshipType)
	out.SetFrom(props.From)
	out.SetTo(props.To)
	return out
}

type RelationshipProps struct {
	// Element properties
	SpdxID ElementID
	Name   string

	// Relationship properties
	RelationshipType RelationshipType
	From             Element
	To               []Element
}

func NewSpdxDocument(props SpdxDocumentProps) SpdxDocument {
	out := &SpdxDocumentImpl{}
	out.SetSpdxID(props.SpdxID)
	out.SetName(props.Name)
	out.SetProfileConformance(props.ProfileConformance)
	out.SetRootElement(props.RootElement)
	out.SetElements(props.Element)
	return out
}

type SpdxDocumentProps struct {
	// Element properties
	SpdxID ElementID
	Name   string

	// ElementCollection properties
	ProfileConformance ProfileIdentifierType
	RootElement        []Element
	Element            []Element
}
