package v3_0

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// ==============================
// Conversion Context
// ==============================

type ConversionContext struct {
	idMap    map[string]common.ElementID
	Counter  int
	Warnings []string
}

func StartConversion() *ConversionContext {
	return &ConversionContext{
		idMap: make(map[string]common.ElementID),
	}
}

// ==============================
// Relationship Mapping
// ==============================

type RelationshipMapping struct {
	V2Type  string
	Reverse bool
	ToField string
}

var relMap = map[string]RelationshipMapping{
	// keys must match rel.Type.GetID()
	"https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/dependsOn": {
		V2Type: string(common.TypeRelationshipDependsOn),
	},
	"https://spdx.org/rdf/3.0.1/terms/Core/RelationshipType/contains": {
		V2Type: string(common.TypeRelationshipContains),
	},
	// placeholder for N:1
	"https://spdx.org/rdf/3.0.1/terms/Software/hasConcludedLicense": {
		ToField: "PackageLicenseConcluded",
	},
}

// ==============================
// Package Conversion
// ==============================

func ConvertPackageNameToV23(p AnyPackage) *v2_3.Package {
	if p == nil {
		return nil
	}

	return &v2_3.Package{
		PackageName: p.GetName(),
	}
}

func ConvertPackageToV23(p AnyPackage, ctx *ConversionContext) *v2_3.Package {
	if p == nil || ctx == nil {
		return nil
	}

	v3ID := p.GetID()

	if v3ID == "" {
		ctx.Warnings = append(ctx.Warnings, "package missing ID")
		return nil
	}

	if existing, ok := ctx.idMap[v3ID]; ok {
		return &v2_3.Package{
			PackageName:           p.GetName(),
			PackageSPDXIdentifier: existing,
		}
	}

	ctx.Counter++
	spdxID := fmt.Sprintf("SPDXRef-Package-%d", ctx.Counter)

	ctx.idMap[v3ID] = common.ElementID(spdxID)

	return &v2_3.Package{
		PackageName:           p.GetName(),
		PackageSPDXIdentifier: common.ElementID(spdxID),
	}
}

// ==============================
// Relationship Conversion
// ==============================

func ConvertRelationshipToV23(rel *Relationship, ctx *ConversionContext) []*v2_3.Relationship {
	if rel == nil || ctx == nil {
		return nil
	}

	typeID := rel.Type.GetID()

	mapping, ok := relMap[typeID]
	if !ok {
		ctx.Warnings = append(ctx.Warnings, "unsupported relationship: "+typeID)
		return nil
	}

	fromID := rel.From.GetID()
	if fromID == "" {
		ctx.Warnings = append(ctx.Warnings, "relationship source missing ID")
		return nil
	}

	v2From, ok := ctx.idMap[fromID]
	if !ok {
		ctx.Warnings = append(ctx.Warnings, "missing ID mapping for source: "+fromID)
		return nil
	}

	var results []*v2_3.Relationship

	for _, to := range rel.To {
		toID := to.GetID()

		if toID == "" {
			ctx.Warnings = append(ctx.Warnings, "relationship target missing ID")
			continue
		}

		v2To, ok := ctx.idMap[toID]
		if !ok {
			ctx.Warnings = append(ctx.Warnings, "missing ID mapping for target: "+toID)
			continue
		}

		// N:1 placeholder
		if mapping.ToField != "" {
			ctx.Warnings = append(ctx.Warnings, "field mapping not implemented: "+typeID)
			continue
		}

		src := v2From
		dst := v2To

		if mapping.Reverse {
			src, dst = dst, src
		}

		results = append(results, &v2_3.Relationship{
			RefA: common.DocElementID{
				ElementRefID: src,
			},
			RefB: common.DocElementID{
				ElementRefID: dst,
			},
			Relationship: mapping.V2Type,
		})
	}

	return results
}

// ==============================
// Warnings → Annotations
// ==============================

func ApplyWarningsAsAnnotations(doc *v2_3.Document, ctx *ConversionContext) {
	if doc == nil || ctx == nil {
		return
	}

	for _, w := range ctx.Warnings {
		doc.Annotations = append(doc.Annotations, &v2_3.Annotation{
			Annotator: common.Annotator{
				Annotator:     "SPDX-Converter",
				AnnotatorType: "Tool",
			},
			AnnotationType:    "OTHER",
			AnnotationComment: "LOSSY: " + w,
		})
	}
}