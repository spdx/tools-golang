package v3_0

import (
	"errors"
	"time"

	v3_0_e "github.com/spdx/tools-golang/spdx/v3/v3_0/embedded_with_interfaces"
)

// Could add convenience creation props like:

func NewPackage(props PackageProps) (v3_0_e.Package, error) {
	out := &v3_0_e.PackageImpl{}
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
	SpdxID v3_0_e.ElementID
	Name   string

	// Artifact properties
	BuiltTime time.Time

	// SoftwareArtifact properties
	CopyrightText string

	// Package properties
	PackageVersion string
}

func NewFile(props FileProps) (v3_0_e.File, error) {
	out := &v3_0_e.FileImpl{}
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
	SpdxID v3_0_e.ElementID
	Name   string

	// Artifact properties
	BuiltTime time.Time

	// SoftwareArtifact properties
	CopyrightText string

	// File properties
	ContentType string
}

func NewRelationship(props RelationshipProps) (v3_0_e.Relationship, error) {
	out := &v3_0_e.RelationshipImpl{}
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
	SpdxID v3_0_e.ElementID
	Name   string

	// Relationship properties
	RelationshipType v3_0_e.RelationshipType
	From             v3_0_e.Element
	To               []v3_0_e.Element
}

func NewSpdxDocument(props SpdxDocumentProps) (v3_0_e.SpdxDocument, error) {
	out := &v3_0_e.SpdxDocumentImpl{}
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
	SpdxID v3_0_e.ElementID
	Name   string

	// ElementCollection properties
	ProfileConformance v3_0_e.ProfileIdentifierType
	RootElement        []v3_0_e.Element
	Element            []v3_0_e.Element
}
