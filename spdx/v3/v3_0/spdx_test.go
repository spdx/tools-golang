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
	doc.RootElements.Append(sbom)

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

	sbom.RootElements.Append(pkg1, pkg2)

	// add a file

	file1 := &spdx.SoftwareFile{SoftwareSoftwareArtifact: spdx.SoftwareSoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
		Name: "/bin/bash",
	}}}}
	sbom.RootElements.Append(file1)

	// add relationships

	sbom.RootElements.Append(&spdx.Relationship{
		From:             file1,
		RelationshipType: spdx.RelationshipType_Contains,
		Tos: spdx.ElementList{
			pkg1,
			pkg2,
		},
	})

	sbom.RootElements.Append(&spdx.Relationship{
		From:             pkg1,
		RelationshipType: spdx.RelationshipType_DependsOn,
		Tos: spdx.ElementList{
			pkg2,
		},
	})

	// serialize

	//buf := bytes.Buffer{}
	//err := doc.ToJSON(&buf)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//json1 := buf.String()
	//fmt.Printf("--------- initial JSON: ----------\n%s\n\n", json1)
	//
	//// deserialize to a new document
	//
	//doc = spdx.NewDocument(&spdx.SoftwareAgent{})
	//err = doc.FromJSON(strings.NewReader(json1))
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//// re-serialize
	//
	//buf.Reset()
	//err = doc.ToJSON(&buf)
	//if err != nil {
	//	t.Error(err)
	//}
	//json2 := buf.String()
	//fmt.Printf("--------- reserialized JSON: ----------\n%s\n", json2)
	//
	//// compare original to parsed and re-encoded
	//
	//diff := difflib.UnifiedDiff{
	//	A:        difflib.SplitLines(json1),
	//	B:        difflib.SplitLines(json2),
	//	FromFile: "Original",
	//	ToFile:   "Current",
	//	Context:  3,
	//}
	//text, _ := difflib.GetUnifiedDiffString(diff)
	//if text != "" {
	//	t.Errorf(text)
	//}

	// some basic usage:

	var pkgs []*spdx.SoftwarePackage
	for _, sbom := range doc.RootElements.SoftwareSbomIter() {
		for _, rel := range sbom.RootElements.RelationshipIter() {
			if rel.RelationshipType == spdx.RelationshipType_Contains {
				spdx.As(rel.From, func(f *spdx.SoftwareFile) {
					if f.Name == "/bin/bash" {
						for _, pkg := range rel.Tos.SoftwarePackageIter() {
							pkgs = append(pkgs, pkg)
						}
					}
				})
			}
		}
	}
	if len(pkgs) != 2 {
		t.Error("wrong packages returned")
	}
}

func Test_aiProfile(t *testing.T) {
	doc := spdx.NewDocument(spdx.ProfileIdentifierType_Ai, "", &spdx.SoftwareAgent{Agent: spdx.Agent{Element: spdx.Element{
		Name:    "tools-golang",
		Summary: "a summary",
	}}}, nil)

	aiPkg := &spdx.AiAIPackage{
		SoftwarePackage: spdx.SoftwarePackage{SoftwareSoftwareArtifact: spdx.SoftwareSoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
			Name: "some ai package",
		}}}},
		AiEnergyConsumption: &spdx.AiEnergyConsumption{
			AiFinetuningEnergyConsumptions: spdx.AiEnergyConsumptionDescriptionList{
				&spdx.AiEnergyConsumptionDescription{
					AiEnergyQuantity: 1.2,
					AiEnergyUnit:     spdx.AiEnergyUnitType_KilowattHour,
				},
			},
			AiTrainingEnergyConsumptions: spdx.AiEnergyConsumptionDescriptionList{
				&spdx.AiEnergyConsumptionDescription{
					AiEnergyQuantity: 5032402,
					AiEnergyUnit:     spdx.AiEnergyUnitType_KilowattHour,
				},
			},
		},
		AiTypeOfModels: []string{
			"Llama 3 8B",
		},
	}

	doc.RootElements.Append(aiPkg)

	// serialize

	//buf := bytes.Buffer{}
	//err := doc.ToJSON(&buf)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//json1 := buf.String()
	//fmt.Printf("--------- initial JSON: ----------\n%s\n\n", json1)
	//
	//// deserialize to a new document
	//
	//doc = spdx.NewDocument(&spdx.SoftwareAgent{})
	//err = doc.FromJSON(strings.NewReader(json1))
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//// re-serialize
	//
	//buf.Reset()
	//err = doc.ToJSON(&buf)
	//if err != nil {
	//	t.Error(err)
	//}
	//json2 := buf.String()
	//fmt.Printf("--------- reserialized JSON: ----------\n%s\n", json2)
	//
	//// compare original to parsed and re-encoded
	//
	//diff := difflib.UnifiedDiff{
	//	A:        difflib.SplitLines(json1),
	//	B:        difflib.SplitLines(json2),
	//	FromFile: "Original",
	//	ToFile:   "Current",
	//	Context:  3,
	//}
	//text, _ := difflib.GetUnifiedDiffString(diff)
	//if text != "" {
	//	t.Errorf(text)
	//}
}
