package v3_0_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx/v3/internal/ld"
	spdx "github.com/spdx/tools-golang/spdx/v3/v3_0"
)

func Test_customSerialization(t *testing.T) {
	d := spdx.NewDocument(spdx.ProfileIdentifierType_Software, "adoc", &spdx.Person{
		Name: "Keith",
		ExternalIdentifiers: spdx.ExternalIdentifierList{
			&spdx.ExternalIdentifier{
				Type:       spdx.ExternalIdentifierType_Email,
				Identifier: "keith@example.com",
			},
		},
	}, nil)

	sbom := &spdx.SBOM{}
	d.RootElements = spdx.ElementList{sbom} // only 1 element appended to document RootElement list

	p := &spdx.Package{}
	p.Name = "a pkg"
	sbom.RootElements = append(sbom.RootElements, p)

	p.ContentIdentifiers = append(p.ContentIdentifiers, &spdx.ContentIdentifier{
		Type:  spdx.ContentIdentifierType_Gitoid,
		Value: "http://example.org",
	})

	a := &spdx.Person{
		Name: "Alice",
	}
	sbom.RootElements = append(sbom.RootElements, a)

	validationErr := d.Validate(false)
	require.Error(t, validationErr) // we did not set creationInfo, should be invalid document
	contents := bytes.Buffer{}
	err := d.ToJSON(&contents)
	require.NoError(t, err)

	structure := map[string]any{}
	err = json.Unmarshal(contents.Bytes(), &structure)
	require.NoError(t, err)

	graph := structure["@graph"].([]any)

	var spdxDoc map[string]any
	for _, e := range graph {
		if e, ok := e.(map[string]any); ok {
			if e["type"] == "SpdxDocument" {
				spdxDoc = e
				break
			}
		}
	}
	require.NotNil(t, spdxDoc)

	rootElements, _ := spdxDoc["rootElement"].([]any)
	require.NotNil(t, rootElements)

	serializedSBOMId := ""
	for _, e := range rootElements {
		if e, ok := e.(string); ok {
			serializedSBOMId = e
		}
	}
	require.NotEmpty(t, serializedSBOMId)

	// sbom is serialized to the graph root, find it by id:
	var serializedSBOM map[string]any
	for _, e := range graph {
		if e, ok := e.(map[string]any); ok {
			if e["type"] == "software_Sbom" && e["spdxId"] == serializedSBOMId {
				serializedSBOM = e
			}
		}
	}
	require.NotNil(t, serializedSBOM)

	rootElements, _ = serializedSBOM["rootElement"].([]any)
	require.NotNil(t, rootElements)
	require.Len(t, rootElements, 2) // we specifically added 2 elements to the SBOM
}

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
	pkg1 := &spdx.Package{
		Name: "the pkg 2",
	}
	file1 := &spdx.File{
		Name:        "a file",
		ContentType: "text", // validation error
		Kind:        spdx.FileKindType{},
	}
	d.RootElements = spdx.ElementList{
		&spdx.SBOM{
			Name:     "My Bom",
			Elements: spdx.ElementList{},
		},
		file1,
		pkg1,
		&spdx.Package{
			ID:          "some ID!",
			Name:        "some name",
			Description: "descr",
			ExternalIdentifiers: spdx.ExternalIdentifierList{
				&spdx.ExternalIdentifier{
					IdentifierLocators: []ld.URI{
						"locator1",
						"locator2",
					},
					Identifier: "CVE-2024-1234",
					Type:       spdx.ExternalIdentifierType_Cve,
				},
			},
			ExternalRefs:  nil,
			Summary:       "",
			VerifiedUsing: nil,
			StandardNames: []string{
				"standard-name1",
				"standard-name2",
			},
			ReleaseTime: time.Now(),
			AdditionalPurposes: []spdx.SoftwarePurpose{
				spdx.SoftwarePurpose_Container,
				spdx.SoftwarePurpose_Library,
			},
			PrimaryPurpose: spdx.SoftwarePurpose_Application,
			CopyrightText:  "",
		},
	}

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
	d2 := newTestDocument()
	err = d2.FromJSON(&buf)
	require.NoError(t, err)
	// these would (correctly) cause a failure, to validate unexported and time fields are being properly checked:
	//d2.RootElements.Packages().Views()[1].PrimaryPurpose = spdx.SoftwarePurpose_Archive
	//_ = spdx.As(d2.CreationInfo, func(info *spdx.CreationInfo) error {
	//	info.Created = info.Created.Add(time.Hour)
	//	return nil
	//})
	diff := cmp.Diff(d.SpdxDocument, d2.SpdxDocument, diffOpts()...)
	require.Empty(t, diff)
}

func diffOpts() []cmp.Option {
	var out []cmp.Option
	for _, t := range []any{
		spdx.Package{},
		spdx.AIPackage{},
		spdx.Relationship{},
		spdx.File{},
		spdx.Snippet{},
		spdx.Annotation{},
		spdx.Tool{},
		spdx.Person{},
		spdx.Organization{},
		spdx.CustomLicense{},
		spdx.ListedLicense{},
		spdx.SpdxDocument{},
		spdx.SBOM{},
	} {
		out = append(out,
			cmpopts.IgnoreUnexported(t),
			cmpopts.IgnoreFields(t, "ID", "CreationInfo"),
		)
	}
	for _, t := range []any{
		spdx.SpdxDocument{},
		spdx.SBOM{},
		spdx.Bundle{},
	} {
		out = append(out,
			cmpopts.IgnoreFields(t, "Elements"),
		)
	}
	out = append(out,
		cmp.Transformer("truncate_time.Time", func(t time.Time) time.Time {
			return t.Truncate(time.Second)
		}),
		cmpopts.IgnoreFields(spdx.Document{}, "LDContext"),
		cmpopts.IgnoreFields(spdx.CreationInfo{}, "CreatedUsing"),
		cmpopts.EquateComparable(
			spdx.ExternalIdentifierType{},
			spdx.HashAlgorithm{},
			spdx.FileKindType{},
			spdx.SoftwarePurpose{},
			spdx.PresenceType{},
			spdx.SafetyRiskAssessmentType{},
			spdx.RelationshipCompleteness{},
			spdx.RelationshipType{},
			spdx.AnnotationType{},
			spdx.ProfileIdentifierType{},
			spdx.ExternalIRI{},
			spdx.EnergyUnitType{},
		),
	)
	return out
}

func Test_reader(t *testing.T) {
	f, err := os.Open("testdata/test.json")
	require.NoError(t, err)
	d := newTestDocument()
	err = d.FromJSON(f)
	require.NoError(t, err)
	sboms := d.RootElements.SBOMs()
	require.Len(t, sboms, 1)

	packages := sboms[0].GetRootElements().Packages()
	require.Len(t, packages, 1)

	rels := d.Elements.Relationships()
	require.Len(t, rels, 1)

	// this is the only reference to the package I see:
	p := rels[0].GetFrom().(*spdx.Package)
	require.NotNil(t, p)
	require.NotEqual(t, time.Time{}, p.BuiltTime)
	require.Equal(t, "my-package", p.Name)
}

func Test_readerExpanded(t *testing.T) {
	f, err := os.Open("testdata/test.expanded.json")
	require.NoError(t, err)
	d := newTestDocument()
	err = d.FromJSON(f)
	require.NoError(t, err)
	for _, fi := range d.Elements.Files() {
		//println("File ID: " + fi.ID)
		println("File Name: " + fi.GetName())
	}

	sboms := d.RootElements.SBOMs()
	require.Len(t, sboms, 1)

	packages := sboms[0].GetRootElements().Packages()
	require.Len(t, packages, 1)

	rels := d.Elements.Relationships()
	require.Len(t, rels, 1)

	// this is the only reference to the package I see:
	p := rels[0].GetFrom().(*spdx.Package)
	require.NotNil(t, p)
	require.NotEqual(t, time.Time{}, p.BuiltTime)
	require.Equal(t, "my-package", p.Name)
}

func Test_reader2(t *testing.T) {
	contents := `{
	  "@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
	  "@graph": [
		{ "type": "SpdxDocument",
		  "rootElement": [ "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization" ]
		}
	  ]
	}`

	d := newTestDocument()
	err := d.FromJSON(strings.NewReader(contents))
	require.NoError(t, err)
	require.Len(t, d.RootElements, 1)
	spdxOrgID, _ := ld.GetID(spdx.Organization_SpdxOrganization)
	gotID, _ := ld.GetID(d.RootElements[0])
	require.Equal(t, spdxOrgID, gotID)
}

func Test_exportImportExport(t *testing.T) {
	// create a document
	doc := spdx.NewDocument(
		spdx.ProfileIdentifierType_Software,
		"My Document",
		&spdx.SoftwareAgent{
			Name:    "tools-golang",
			Summary: "a summary",
		},
		&spdx.Tool{
			ExternalIdentifiers: spdx.ExternalIdentifierList{
				&spdx.ExternalIdentifier{
					Type:       spdx.ExternalIdentifierType_Cpe23,
					Identifier: "cpe:2.3:a:myvendor:my-product:*:*:*:*:*:*:*:*",
				},
			},
			ExternalRefs: nil,
			Name:         "not-tools-golang",
		})

	sbom := &spdx.SBOM{}
	doc.RootElements = append(doc.RootElements, sbom)

	// create a package

	pkg1 := &spdx.Package{
		Name:    "some-package-1",
		Version: "1.2.3",
	}

	// create another package

	pkg2 := &spdx.Package{}
	pkg2.Name = "some-package-2"
	pkg2.Version = "2.4.5"

	// add the packages to the sbom

	sbom.RootElements = append(sbom.RootElements, pkg1, pkg2)

	// add a file

	file1 := &spdx.File{
		Name: "/bin/bash",
	}
	sbom.RootElements = append(sbom.RootElements, file1)

	// add relationships

	sbom.RootElements = append(sbom.RootElements, &spdx.Relationship{
		From: file1,
		Type: spdx.RelationshipType_Contains,
		To: spdx.ElementList{
			pkg1,
			pkg2,
		},
	})

	sbom.RootElements = append(sbom.RootElements, &spdx.Relationship{
		From: pkg1,
		Type: spdx.RelationshipType_DependsOn,
		To: spdx.ElementList{
			pkg2,
		},
	})

	// serialize

	buf := bytes.Buffer{}
	err := doc.ToJSON(&buf)
	require.NoError(t, err)

	json1 := buf.String()

	// deserialize to a new document

	newDoc := newTestDocument()
	err = newDoc.FromJSON(strings.NewReader(json1))
	require.NoError(t, err)

	// compare original to parsed -- this includes Element lists, etc.
	diff := cmp.Diff(doc, newDoc, diffOpts()...)
	require.Empty(t, diff)

	// some basic usage:

	var pkgs []spdx.AnyPackage
	for _, sbom := range doc.RootElements.SBOMs() {
		for _, rel := range sbom.GetRootElements().Relationships() {
			if rel.GetType() != spdx.RelationshipType_Contains {
				continue
			}
			f := rel.GetFrom().(*spdx.File)
			if f.Name == "/bin/bash" {
				for _, pkg := range rel.GetTo().Packages() {
					pkgs = append(pkgs, pkg)
				}
			}
		}
	}
	require.Len(t, pkgs, 2)
}

func Test_aiProfile(t *testing.T) {
	doc := spdx.NewDocument(spdx.ProfileIdentifierType_Ai, "", &spdx.SoftwareAgent{
		Name:    "tools-golang",
		Summary: "a summary",
	}, nil)

	aiPkg := &spdx.AIPackage{
		Name: "some ai package",
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
	require.NoError(t, err)

	json1 := buf.String()

	// deserialize to a new document

	newDoc := newTestDocument()
	err = newDoc.FromJSON(strings.NewReader(json1))
	require.NoError(t, err)

	diff := cmp.Diff(doc, newDoc, diffOpts()...)
	require.Empty(t, diff)
}

func newTestDocument() *spdx.Document {
	return spdx.NewDocument(spdx.ProfileIdentifierType_Software, "test document",
		&spdx.SoftwareAgent{Name: "tools-golang-tests-agent", Summary: "a summary"},
		&spdx.Tool{Name: "tools-golang-tests-tool"})
}
