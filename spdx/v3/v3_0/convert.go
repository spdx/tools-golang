package v3_0

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/kzantow/go-ld"

	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func From_v2_3(doc v2_3.Document, d *Document) {
	if d.CreationInfo == nil {
		d.CreationInfo = &CreationInfo{}
	}
	if c, ok := d.CreationInfo.(*CreationInfo); ok {
		c.SpecVersion = Version
	}

	idMap := map[string]any{}

	d.ID = string(doc.SPDXIdentifier)
	d.Comment = doc.DocumentComment
	d.CreationInfo = convert23CreationInfo(doc.CreationInfo)
	d.Name = doc.DocumentName
	d.DataLicense = &LicenseInfo{Element: Element{
		Name: doc.DataLicense,
	}}

	for _, pkg := range doc.Packages {
		newPkg := convert23Package(idMap, pkg)
		if newPkg != nil {
			d.RootElements = append(d.RootElements, newPkg)
			idMap[string(pkg.PackageSPDXIdentifier)] = newPkg
		}
	}

	for _, file := range doc.Files {
		newFile := convert23File(idMap, file)
		if newFile != nil {
			d.RootElements = append(d.RootElements, newFile)
			idMap[string(file.FileSPDXIdentifier)] = newFile
		}
	}

	rels := relMap{}
	for _, rel := range doc.Relationships {
		newRel := convert23Relationship(idMap, rel)
		if newRel != nil {
			rels.add(newRel)
		}
	}
	for _, relTypes := range rels {
		for _, to := range relTypes {
			d.RootElements = append(d.RootElements, to)
		}
	}
}

type relMap map[AnyElement]map[RelationshipType]*Relationship

func (r relMap) add(relationship *Relationship) {
	relTypes := r[relationship.From]
	if relTypes == nil {
		relTypes = map[RelationshipType]*Relationship{}
		r[relationship.From] = relTypes
	}
	existing := relTypes[relationship.RelationshipType]
	if existing == nil {
		relTypes[relationship.RelationshipType] = relationship
		return
	}
	existing.To = appendUnique(existing.To, relationship.To...)
}

func appendUnique[T comparable](existing []T, adding ...T) []T {
	for _, add := range adding {
		if slices.Contains(existing, add) {
			continue
		}
		existing = append(existing, add)
	}
	return existing
}

func convert23Relationship(idMap map[string]any, rel *v2_3.Relationship) *Relationship {
	if rel == nil {
		return nil
	}
	from, _ := idMap[string(rel.RefA.ElementRefID)].(AnyElement)
	to, _ := idMap[string(rel.RefB.ElementRefID)].(AnyElement)
	if from == nil || to == nil {
		return nil
	}
	return &Relationship{
		Element: Element{
			Comment: rel.RelationshipComment,
		},
		From:             from,
		RelationshipType: convert23RelationshipType(rel.Relationship),
		To:               []AnyElement{to},
	}
}

func convert23RelationshipType(relationship string) RelationshipType {
	switch strings.ToLower(relationship) {
	case "contains", "contained_by":
		return RelationshipType_Contains
	case "depends_on", "dependency_of":
		return RelationshipType_DependsOn
	}
	return RelationshipType{}
}

func convert23CreationInfo(info *v2_3.CreationInfo) *CreationInfo {
	return &CreationInfo{
		Comment:     info.CreatorComment,
		Created:     convert23Time(info.Created),
		CreatedBy:   convert23Creators(info.Creators),
		SpecVersion: info.LicenseListVersion,
	}
}

func convert23Creators(creators []common.Creator) []AnyAgent {
	var out []AnyAgent
	for _, c := range creators {
		out = append(out, convert23Agent(c.CreatorType, c.Creator))
	}
	return out
}

func convert23Agent(agentType string, agent string) AnyAgent {
	switch strings.ToLower(agentType) {
	case "person":
		return &Person{Agent: Agent{Element: Element{
			Name: agent,
		}}}
	case "organization", "org":
		return &Organization{Agent: Agent{Element: Element{
			Name: agent,
		}}}
	}
	return nil
}

func convert23File(idMap map[string]any, file *v2_3.File) AnyFile {
	if file == nil {
		return nil
	}
	return &File{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{Element: Element{
				ID:      string(file.FileSPDXIdentifier),
				Comment: file.FileComment,
				Name:    file.FileName,
			}},
			CopyrightText:    file.FileCopyrightText,
			AttributionTexts: file.FileAttributionTexts,
		},
		FileKind: FileKindType_File,
	}
}

func convert23Package(idMap map[string]any, pkg *v2_3.Package) *Package {
	if pkg == nil {
		return nil
	}
	return &Package{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					ID:                  string(pkg.PackageSPDXIdentifier),
					Name:                pkg.PackageName,
					Summary:             pkg.PackageSummary,
					Comment:             pkg.PackageComment,
					Description:         pkg.PackageDescription,
					ExternalIdentifiers: convert23ExternalIdentifiers(pkg.PackageExternalReferences),
					VerifiedUsing:       convert23VerifiedUsing(pkg.PackageVerificationCode),
				},
				BuiltTime:      convert23Time(pkg.BuiltDate),
				OriginatedBy:   convert23PackageOriginator(pkg.PackageOriginator),
				ReleaseTime:    convert23Time(pkg.ReleaseDate),
				SuppliedBy:     convert23Supplier(pkg.PackageSupplier),
				ValidUntilTime: convert23Time(pkg.ValidUntilDate),
			},
			AttributionTexts: pkg.PackageAttributionTexts,
			CopyrightText:    pkg.PackageCopyrightText,
			PrimaryPurpose:   convert23PrimaryPurpose(pkg.PrimaryPackagePurpose),
		},

		DownloadLocation: convert23URI(pkg.PackageDownloadLocation),
		HomePage:         convert23URI(pkg.PackageHomePage),
		PackageUrl:       convert23PackageUrl(pkg.PackageExternalReferences),
		PackageVersion:   pkg.PackageVersion,
		SourceInfo:       pkg.PackageSourceInfo,
	}
}

func convert23Time(date string) time.Time {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		logDropped(err)
		return time.Time{}
	}
	return t
}

func convert23URI(uri string) ld.URI {
	return ld.URI(uri)
}

func logDropped(value any) {
	_, _ = fmt.Fprintf(os.Stderr, "dropped: %v", value)
}

func convert23PackageUrl(references []*v2_3.PackageExternalReference) ld.URI {
	for _, ref := range references {
		if ref.RefType == common.TypePackageManagerPURL {
			return convert23URI(ref.Locator)
		}
	}
	return ""
}

func convert23PrimaryPurpose(purpose string) SoftwarePurpose {
	switch purpose {
	case "container":
		return SoftwarePurpose_Container
	case "library":
		return SoftwarePurpose_Library
	case "application":
		return SoftwarePurpose_Application
	}
	return SoftwarePurpose{}
}

func convert23Supplier(supplier *common.Supplier) AnyAgent {
	if supplier == nil {
		return nil
	}
	return convert23Agent(supplier.SupplierType, supplier.Supplier)
}

func convert23PackageOriginator(originator *common.Originator) []AnyAgent {
	if originator == nil {
		return nil
	}
	return []AnyAgent{convert23Agent(originator.OriginatorType, originator.Originator)}
}

func convert23VerifiedUsing(verificationCode *common.PackageVerificationCode) []AnyIntegrityMethod {
	return nil // TODO
}

func convert23ExternalIdentifiers(references []*v2_3.PackageExternalReference) []AnyExternalIdentifier {
	var out []AnyExternalIdentifier
	for _, r := range references {
		typ := ExternalIdentifierType{}
		switch r.RefType {
		case common.TypeSecurityCPE22Type:
			typ = ExternalIdentifierType_Cpe22
		case common.TypeSecurityCPE23Type:
			typ = ExternalIdentifierType_Cpe23
		case common.TypePackageManagerPURL:
			typ = ExternalIdentifierType_PackageUrl
		default:
			continue // unknown
		}
		out = append(out, &ExternalIdentifier{
			Comment:                r.ExternalRefComment,
			ExternalIdentifierType: typ,
			Identifier:             r.Locator,
		})
	}
	return out
}
