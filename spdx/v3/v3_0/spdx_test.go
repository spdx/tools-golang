package v3_0_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kzantow/go-ld"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"

	spdx "github.com/spdx/tools-golang/spdx/v3/v3_0"
)

func Test_validateMinList(t *testing.T) {
	a := &spdx.Person{}
	a.Name = "me"
	e := &spdx.CreationInfo{
		SpecVersion:  "3.0.1",
		Created:      time.Now(),
		CreatedUsing: nil,
		Comment:      "",
		CreatedBy:    spdx.AgentList{},
	}
	err := e.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "must have")
}

func Test_writer(t *testing.T) {
	d := newTestDocument()
	pkg1 := &spdx.Package{SoftwareArtifact: spdx.SoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
		Name: "the pkg 2",
	}}}}
	file1 := &spdx.File{
		SoftwareArtifact: spdx.SoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
			Name: "a file",
		}}},
		ContentType: "text", // validation error
		FileKind:    spdx.FileKindType{},
	}
	d.Append(
		&spdx.Sbom{Bom: spdx.Bom{Bundle: spdx.Bundle{ElementCollection: spdx.ElementCollection{
			Element: spdx.Element{
				Name: "My Bom",
			},
			Elements: spdx.ElementList{},
		}}},
			SbomTypes: nil,
		},
		file1,
		pkg1,
		&spdx.Package{
			SoftwareArtifact: spdx.SoftwareArtifact{
				Artifact: spdx.Artifact{
					Element: spdx.Element{
						ID:          "some ID!",
						Name:        "some name",
						Description: "descr",
						ExternalIdentifiers: spdx.ExternalIdentifierList{
							&spdx.ExternalIdentifier{
								IdentifierLocators: []ld.URI{
									"locator1",
									"locator2",
								},
								Identifier:             "CVE-2024-1234",
								ExternalIdentifierType: spdx.ExternalIdentifierType_Cve,
							},
						},
						ExternalRefs:   nil,
						Summary:        "",
						VerifiedUsings: nil,
					},
					StandardNames: []string{
						"standard-name1",
						"standard-name2",
					},
					ReleaseTime: time.Now(),
				},
				AdditionalPurposes: []spdx.SoftwarePurpose{
					spdx.SoftwarePurpose_Container,
					spdx.SoftwarePurpose_Library,
				},
				PrimaryPurpose: spdx.SoftwarePurpose_Application,
				CopyrightText:  "",
			},
		},
	)

	// many validation issues
	err := d.Validate(false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CreationInfo")

	err = d.Validate(true)
	require.Error(t, err)
	require.NotContains(t, err.Error(), "CreationInfo")

	// fix validation issue
	file1.ContentType = "text/plain"
	err = d.Validate(false) // already use creationInfo
	require.NoError(t, err)

	buf := bytes.Buffer{}
	err = d.ToJSON(&buf)
	require.NoError(t, err)
	fmt.Printf("%#v\n", buf.String())

	d2 := newTestDocument()
	err = d2.FromJSON(&buf)
	require.NoError(t, err)
	// these would (correctly) cause a failure, to validate unexported and time fields are being properly checked:
	//d2.RootElements.Packages().Views()[1].PrimaryPurpose = spdx.SoftwarePurpose_Archive
	//_ = spdx.As(d2.CreationInfo, func(info *spdx.CreationInfo) error {
	//	info.Created = info.Created.Add(time.Hour)
	//	return nil
	//})
	diff := cmp.Diff(d.SpdxDocument, d2.SpdxDocument, testOpts...)
	if diff != "" {
		t.Fatal(diff)
	}
}

var testOpts = []cmp.Option{
	cmp.Transformer("truncate_time.Time", func(t time.Time) time.Time {
		return t.Truncate(time.Second)
	}),
	// export and compare unexported fields
	cmp.Exporter(func(r reflect.Type) bool {
		return true
	}),
}

func Test_reader(t *testing.T) {
	f, err := os.Open("test.json")
	require.NoError(t, err)
	d := newTestDocument()
	err = d.FromJSON(f)
	require.NoError(t, err)
	fmt.Printf("%#v\n", d)

	require.Equal(t, d.Elements.Files().Len(), 1)
	for _, fi := range d.Elements.Files() {
		if fi.PrimaryPurpose == spdx.SoftwarePurpose_Executable {
			println("Got Executable File ID: " + fi.ID)
		}
	}

	// the example is incorrect, it doesn't include the package in the element root
	pkgs := d.Elements.Packages().Views()
	//require.Len(t, pkgs, 1)
	//require.NotEqual(t, time.Time{}, pkgs[0].BuiltTime)
	require.Empty(t, pkgs) // FIXME this shouldn't be true, but the example is wrong

	rels := d.Elements.Relationships().Views()
	require.Len(t, rels, 1)

	if p, ok := rels[0].From.(*spdx.Package); ok {
		require.NotEqual(t, time.Time{}, p.BuiltTime)
	}
	// this is the only reference to the package I see:
	_ = spdx.As(rels[0].From, func(p *spdx.Package) error {
		require.NotEqual(t, time.Time{}, p.BuiltTime)
		return nil
	})
}

func Test_readerExpanded(t *testing.T) {
	f, err := os.Open("test.expanded.json")
	require.NoError(t, err)
	d := newTestDocument()
	err = d.FromJSON(f)
	require.NoError(t, err)
	fmt.Printf("%#v\n", d)
	for _, fi := range d.Elements.Files() {
		println("File ID: " + fi.ID)
	}

	pkgs := d.Elements.Packages().Views()
	require.Len(t, pkgs, 1)
	require.NotEqual(t, time.Time{}, pkgs[0].BuiltTime)
}

func Test_reader2(t *testing.T) {
	contents := `
		{
  "@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
  "@graph": [
	{ "spdxId": "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization" },
	{ "spdxId": "SpdxOrganization" },
{
      "type": "CreationInfo",
      "@id": "_:creationinfo",
      "createdBy": [
        "http://spdx.example.com/Agent/JoshuaWatt"
      ],
      "specVersion": "3.0.1",
      "created": "2024-03-06T00:00:00Z"
    }
  ]
}
	`
	graph, err := spdx.LDContext().FromJSON(strings.NewReader(contents))
	require.NoError(t, err)
	for _, fi := range graph {
		println("Elem" + fmt.Sprintf("%#v", fi))
	}
}

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

	sbom := &spdx.Sbom{}
	doc.RootElements = append(doc.RootElements, sbom)

	// create a package

	pkg1 := &spdx.Package{
		SoftwareArtifact: spdx.SoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
			Name: "some-package-1",
		}}},
		PackageVersion: "1.2.3",
	}

	// create another package

	pkg2 := &spdx.AIPackage{}
	pkg2.Name = "some-package-2"
	pkg2.PackageVersion = "2.4.5"

	// add the packages to the sbom

	sbom.RootElements = append(sbom.RootElements, pkg1, pkg2)

	// add a file

	file1 := &spdx.File{SoftwareArtifact: spdx.SoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
		Name: "/bin/bash",
	}}}}
	sbom.RootElements = append(sbom.RootElements, file1)

	// add relationships

	sbom.RootElements = append(sbom.RootElements, &spdx.Relationship{
		From:             file1,
		RelationshipType: spdx.RelationshipType_Contains,
		To: spdx.ElementList{
			pkg1,
			pkg2,
		},
	})

	sbom.RootElements = append(sbom.RootElements, &spdx.Relationship{
		From:             pkg1,
		RelationshipType: spdx.RelationshipType_DependsOn,
		To: spdx.ElementList{
			pkg2,
		},
	})

	// serialize

	buf := bytes.Buffer{}
	err := doc.ToJSON(&buf)
	if err != nil {
		t.Error(err)
	}

	json1 := buf.String()
	fmt.Printf("--------- initial JSON: ----------\n%s\n\n", json1)

	// deserialize to a new document

	newDoc := newTestDocument()
	err = newDoc.FromJSON(strings.NewReader(json1))
	if err != nil {
		t.Error(err)
	}

	// re-serialize

	buf.Reset()
	err = newDoc.ToJSON(&buf)
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

	// some basic usage:

	var pkgs []*spdx.Package
	for _, sbom := range doc.RootElements.Sboms() {
		for _, rel := range sbom.RootElements.Relationships() {
			if rel.RelationshipType != spdx.RelationshipType_Contains {
				continue
			}
			_ = spdx.As(rel.From, func(f *spdx.File) any {
				if f.Name == "/bin/bash" {
					for _, pkg := range rel.To.Packages() {
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

func Test_aiProfile(t *testing.T) {
	doc := spdx.NewDocument(spdx.ProfileIdentifierType_Ai, "", &spdx.SoftwareAgent{Agent: spdx.Agent{Element: spdx.Element{
		Name:    "tools-golang",
		Summary: "a summary",
	}}}, nil)

	aiPkg := &spdx.AIPackage{
		Package: spdx.Package{SoftwareArtifact: spdx.SoftwareArtifact{Artifact: spdx.Artifact{Element: spdx.Element{
			Name: "some ai package",
		}}}},
		EnergyConsumption: &spdx.EnergyConsumption{
			FinetuningEnergyConsumptions: spdx.EnergyConsumptionDescriptionList{
				&spdx.EnergyConsumptionDescription{
					EnergyQuantity: 1.2,
					EnergyUnit:     spdx.EnergyUnitType_KilowattHour,
				},
			},
			TrainingEnergyConsumptions: spdx.EnergyConsumptionDescriptionList{
				&spdx.EnergyConsumptionDescription{
					EnergyQuantity: 5032402,
					EnergyUnit:     spdx.EnergyUnitType_KilowattHour,
				},
			},
		},
		TypeOfModels: []string{
			"Llama 3 8B",
		},
	}

	doc.RootElements = append(doc.RootElements, aiPkg)

	// serialize

	buf := bytes.Buffer{}
	err := doc.ToJSON(&buf)
	if err != nil {
		t.Error(err)
	}

	json1 := buf.String()
	fmt.Printf("--------- initial JSON: ----------\n%s\n\n", json1)

	// deserialize to a new document

	doc = newTestDocument()
	//doc.RootElements.Append(&spdx.Agent{})
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

func newTestDocument() *spdx.Document {
	//return spdx.NewDocument(spdx.ProfileIdentifierType_Lite, "test document",
	return spdx.NewDocument(spdx.ProfileIdentifierType_Software, "test document",
		&spdx.SoftwareAgent{Agent: spdx.Agent{Element: spdx.Element{Name: "tools-golang-tests-agent", Summary: "a summary"}}},
		&spdx.Tool{Element: spdx.Element{Name: "tools-golang-tests-tool"}})
}
