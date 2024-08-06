package v3_0_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"

	spdx "github.com/spdx/tools-golang/spdx/v3/v3_0"
)

func Test_exportImportExport(t *testing.T) {
	doc := spdx.NewDocument(&spdx.SoftwareAgent{
		Name:    "tools-golang",
		Summary: "a summary",
	})
	doc.Document().SetProfileConformance(spdx.ProfileIdentifierType_Software)

	doc.CreationInfo().SetCreatedUsing(
		&spdx.Tool{
			ExternalIdentifiers: []spdx.IExternalIdentifier{
				&spdx.ExternalIdentifier{
					ExternalIdentifierType: spdx.ExternalIdentifierType_Cpe23,
					Identifier:             "cpe:2.3:a:myvendor:my-product:*:*:*:*:*:*:*:*",
				},
			},
			ExternalRefs: nil,
			Name:         "not-tools-golang",
		},
	)

	doc.Document().SetName("My Document")

	// add a package

	pkg1 := &spdx.Package{
		Name:           "some-package-1",
		PackageVersion: "1.2.3",
	}
	pkg2 := &spdx.Package{
		Name:           "some-package-2",
		PackageVersion: "2.4.5",
	}
	doc.AddElement(pkg1, pkg2)

	file1 := &spdx.File{
		Name: "/bin/bash",
	}
	doc.AddElement(file1)

	// add relationships

	doc.AddElement(
		&spdx.Relationship{
			From:             file1,
			RelationshipType: spdx.RelationshipType_Contains,
			To: []spdx.IElement{
				pkg1,
				pkg2,
			},
		},
	)

	doc.AddElement(
		&spdx.Relationship{
			From:             pkg1,
			RelationshipType: spdx.RelationshipType_DependsOn,
			To: []spdx.IElement{
				pkg2,
			},
		},
	)

	// serialize

	buf := bytes.Buffer{}
	err := doc.ToJSON(&buf)
	if err != nil {
		t.Error(err)
	}

	json1 := buf.String()
	fmt.Printf("--------- initial JSON: ----------\n%s\n\n", json1)

	// deserialize to a new document

	doc = spdx.NewDocument(&spdx.SoftwareAgent{})
	err = doc.FromJSON(strings.NewReader(json1))
	if err != nil {
		t.Error(err)
	}

	// re-serialize

	buf.Reset()
	err = doc.ToJSON(&buf)
	if err != nil {
		t.Error(err)
	}
	json2 := buf.String()
	fmt.Printf("--------- reserialized JSON: ----------\n%s\n", json2)

	// compare original to parsed and re-encoded

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(json1),
		B:        difflib.SplitLines(json2),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	if text != "" {
		t.Errorf(text)
	}

	// some basic verification:

	var pkgs []spdx.IPackage
	for _, e := range doc.GetElements() {
		if rel, ok := e.(spdx.IRelationship); ok && rel.GetRelationshipType() == spdx.RelationshipType_Contains {
			if from, ok := rel.GetFrom().(spdx.IFile); ok && from.GetName() == "/bin/bash" {
				for _, el := range rel.GetTo() {
					if pkg, ok := el.(spdx.IPackage); ok {
						pkgs = append(pkgs, pkg)
					}
				}

			}
		}
	}
	if len(pkgs) != 2 {
		t.Error("wrong packages returned")
	}
}

func Test_aiProfile(t *testing.T) {
	doc := spdx.NewDocument(&spdx.SoftwareAgent{
		Name:    "tools-golang",
		Summary: "a summary",
	})
	doc.Document().SetProfileConformance(spdx.ProfileIdentifierType_Ai)

	aiPkg := &spdx.AIPackage{
		Name: "some ai package",
		EnergyConsumption: &spdx.EnergyConsumption{
			FinetuningEnergyConsumption: []spdx.IEnergyConsumptionDescription{
				&spdx.EnergyConsumptionDescription{
					EnergyQuantity: 1.2,
					EnergyUnit:     spdx.EnergyUnitType_KilowattHour,
				},
			},
			TrainingEnergyConsumption: []spdx.IEnergyConsumptionDescription{
				&spdx.EnergyConsumptionDescription{
					EnergyQuantity: 5032402,
					EnergyUnit:     spdx.EnergyUnitType_KilowattHour,
				},
			},
		},
		TypeOfModel: []string{
			"Llama 3 8B",
		},
	}

	doc.AddElement(aiPkg)

	// serialize

	buf := bytes.Buffer{}
	err := doc.ToJSON(&buf)
	if err != nil {
		t.Error(err)
	}

	json1 := buf.String()
	fmt.Printf("--------- initial JSON: ----------\n%s\n\n", json1)

	// deserialize to a new document

	doc = spdx.NewDocument(&spdx.SoftwareAgent{})
	err = doc.FromJSON(strings.NewReader(json1))
	if err != nil {
		t.Error(err)
	}

	// re-serialize

	buf.Reset()
	err = doc.ToJSON(&buf)
	if err != nil {
		t.Error(err)
	}
	json2 := buf.String()
	fmt.Printf("--------- reserialized JSON: ----------\n%s\n", json2)

	// compare original to parsed and re-encoded

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(json1),
		B:        difflib.SplitLines(json2),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	if text != "" {
		t.Errorf(text)
	}
}
