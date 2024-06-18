package v3_0

import (
	"errors"
	"time"

	v3_0 "github.com/spdx/tools-golang/spdx/v3/v3_0/interfaces_only"
)

// Could add convenience creation props like:

func NewPackage(props PackageProps) (v3_0.Package, error) {
	out := v3_0.NewPackage()
	return out, errors.Join(
		out.SetSpdxID(props.SpdxID),
		out.SetName(props.Name),
		out.SetBuiltTime(props.BuiltTime),
		out.SetCopyrightText(props.CopyrightText),
		out.SetPackageVersion(props.PackageVersion),
	)
}

type PackageProps struct {
	// Element properties
	SpdxID v3_0.ElementID
	Name   string

	// Artifact properties
	BuiltTime time.Time

	// SoftwareArtifact properties
	CopyrightText string

	// Package properties
	PackageVersion string
}

func NewFile(props FileProps) (v3_0.File, error) {
	out := v3_0.NewFile()
	return out, errors.Join(
		out.SetSpdxID(props.SpdxID),
		out.SetName(props.Name),
		out.SetBuiltTime(props.BuiltTime),
		out.SetCopyrightText(props.CopyrightText),
		out.SetContentType(props.ContentType),
	)
}

type FileProps struct {
	// Element properties
	SpdxID v3_0.ElementID
	Name   string

	// Artifact properties
	BuiltTime time.Time

	// SoftwareArtifact properties
	CopyrightText string

	// File properties
	ContentType string
}

func NewRelationship(props RelationshipProps) (v3_0.Relationship, error) {
	out := v3_0.NewRelationship()
	return out, errors.Join(
		out.SetSpdxID(props.SpdxID),
		out.SetName(props.Name),
		out.SetRelationshipType(props.RelationshipType),
		out.SetFrom(props.From),
		out.SetTo(props.To),
	)
}

type RelationshipProps struct {
	// Element properties
	SpdxID v3_0.ElementID
	Name   string

	// Relationship properties
	RelationshipType v3_0.RelationshipType
	From             v3_0.Element
	To               []v3_0.Element
}

func NewSpdxDocument(props SpdxDocumentProps) (v3_0.SpdxDocument, error) {
	out := v3_0.NewSpdxDocument()
	return out, errors.Join(
		out.SetSpdxID(props.SpdxID),
		out.SetName(props.Name),
		out.SetProfileConformance(props.ProfileConformance),
		out.SetRootElement(props.RootElement),
		out.SetElements(props.Element),
	)
}

type SpdxDocumentProps struct {
	// Element properties
	SpdxID v3_0.ElementID
	Name   string

	// ElementCollection properties
	ProfileConformance v3_0.ProfileIdentifierType
	RootElement        []v3_0.Element
	Element            []v3_0.Element
}
