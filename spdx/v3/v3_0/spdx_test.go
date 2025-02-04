package v3_0_test

import (
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v3/v3_0"
)

func Test_exportImportExport(t *testing.T) {
	// create a document
	doc := spdx.NewDocument(
		spdx.ProfileIdentifierType_Software,
		"My Document",
		&spdx.SoftwareAgent{Agent: spdx.Agent{Element: spdx.Element{
			Name:    "tools-golang",
			Summary: "a summary",
		}}},
		&spdx.Tool{Element: spdx.Element{
			ExternalIdentifiers: spdx.ExternalIdentifierList{
				&spdx.ExternalIdentifier{
					ExternalIdentifierType: spdx.ExternalIdentifierType_Cpe23,
					Identifier:             "cpe:2.3:a:myvendor:my-product:*:*:*:*:*:*:*:*",
				},
			},
			ExternalRefs: nil,
			Name:         "not-tools-golang",
		}})

	sbom := &spdx.SoftwareSbom{}
	doc.RootElements = append(doc.RootElements, sbom)

	// create a package

	pkg1 := &spdx.SoftwarePackage{
		SoftwareSoftwareArtifact: spdx.SoftwareSoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
			Name: "some-package-1",
		}}},
		SoftwarePackageVersion: "1.2.3",
	}

	// create another package

	pkg2 := &spdx.AiAIPackage{}
	pkg2.Name = "some-package-2"
	pkg2.SoftwarePackageVersion = "2.4.5"

	// add the packages to the sbom

	sbom.RootElements = append(sbom.RootElements, pkg1, pkg2)

	// add a file

	file1 := &spdx.SoftwareFile{SoftwareSoftwareArtifact: spdx.SoftwareSoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
		Name: "/bin/bash",
	}}}}
	sbom.RootElements = append(sbom.RootElements, file1)

	// add relationships

	sbom.RootElements = append(sbom.RootElements, &spdx.Relationship{
		From:             file1,
		RelationshipType: spdx.RelationshipType_Contains,
		Tos: spdx.ElementList{
			pkg1,
			pkg2,
		},
	})

	sbom.RootElements = append(sbom.RootElements, &spdx.Relationship{
		From:             pkg1,
		RelationshipType: spdx.RelationshipType_DependsOn,
		Tos: spdx.ElementList{
			pkg2,
		},
	})

	// some basic usage:

	var pkgs []*spdx.SoftwarePackage
	for _, sbom := range doc.RootElements.SoftwareSbomIter() {
		for _, rel := range sbom.RootElements.RelationshipIter() {
			if rel.RelationshipType != spdx.RelationshipType_Contains {
				continue
			}
			_ = spdx.As(rel.From, func(f *spdx.SoftwareFile) any {
				if f.Name == "/bin/bash" {
					for _, pkg := range rel.Tos.SoftwarePackageIter() {
						pkgs = append(pkgs, pkg)
					}
				}
				return nil
			})

		}
	}
	if len(pkgs) != 2 {
		t.Error("wrong packages returned")
	}
}

func newTestDocument() *spdx.Document {
	return spdx.NewDocument(spdx.ProfileIdentifierType_Lite, "test document",
		&spdx.SoftwareAgent{Agent: spdx.Agent{Element: spdx.Element{Name: "tools-golang-tests-agent", Summary: "a summary"}}},
		&spdx.Tool{Element: spdx.Element{Name: "tools-golang-tests-tool"}})
}
