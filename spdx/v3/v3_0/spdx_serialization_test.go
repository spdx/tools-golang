package v3_0

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_spdxExportImportExport(t *testing.T) {
	doc := SpdxDocument{
		DataLicense: nil,
		Imports:     nil,
	}

	doc.ID = "new-id"

	agent := &SoftwareAgent{Agent: Agent{Element: Element{
		Name:    "some-agent",
		Summary: "summary",
	}}}
	c := &CreationInfo{
		Comment: "some-comment",
		Created: time.Now(),
		CreatedBy: AgentList{
			agent,
		},
		CreatedUsing: []AnyTool{
			&Tool{Element: Element{
				ExternalIdentifiers: ExternalIdentifierList{
					&ExternalIdentifier{
						ExternalIdentifierType: ExternalIdentifierType_Cpe23,
						Identifier:             "cpe23:a:myvendor:my-product:*:*:*:*:*:*:*",
					},
				},
				Name: "not-tools-golang",
			}},
		},
		SpecVersion: "",
	}
	agent.CreationInfo = c

	// add a package

	pkg1 := &Package{SoftwareArtifact: SoftwareArtifact{Artifact: Artifact{Element: Element{
		Name:         "some-package-1",
		CreationInfo: c,
	}}},
		PackageVersion: "1.2.3",
	}
	pkg2 := &Package{SoftwareArtifact: SoftwareArtifact{Artifact: Artifact{Element: Element{
		Name:         "some-package-2",
		CreationInfo: c,
	}}},
		PackageVersion: "2.4.5",
	}
	doc.Elements = append(doc.Elements, pkg2)

	file1 := &File{SoftwareArtifact: SoftwareArtifact{Artifact: Artifact{Element: Element{
		Name:         "/bin/bash",
		CreationInfo: c,
	}}}}
	doc.Elements = append(doc.Elements, file1)

	// add relationships

	doc.Elements = append(doc.Elements,
		&Relationship{Element: Element{
			CreationInfo: c,
		},
			From:             file1,
			RelationshipType: RelationshipType_Contains,
			To: ElementList{
				pkg1,
				pkg2,
			},
		},
	)

	doc.Elements = append(doc.Elements,
		&Relationship{Element: Element{
			CreationInfo: c,
		},
			From:             pkg1,
			RelationshipType: RelationshipType_DependsOn,
			To: ElementList{
				pkg2,
			},
		},
	)

	doc.Elements = append(doc.Elements,
		&AIPackage{Package: Package{SoftwareArtifact: SoftwareArtifact{Artifact: Artifact{Element: Element{
			CreationInfo: c,
		}}}},
			TypeOfModels: []string{"a model"},
		},
	)

	got := encodeDecode(t, &doc)

	// some basic verification:

	var pkgs PackageList
	for _, rel := range got.RootElements.Relationships() {
		if rel.RelationshipType != RelationshipType_Contains {
			continue
		}
		_ = As(rel.From, func(from *File) any {
			if from.Name == "/bin/bash" {
				for _, pkg := range rel.To.Packages() {
					pkgs = append(pkgs, pkg)
				}
			}
			return nil
		})
	}
	if len(pkgs) != 2 {
		t.Error("wrong packages returned")
	}
}

func Test_stringSlice(t *testing.T) {
	p := &AIPackage{
		TypeOfModels: []string{"a model"},
	}
	encodeDecode(t, p)
}

func Test_profileConformance(t *testing.T) {
	doc := &SpdxDocument{ElementCollection: ElementCollection{
		ProfileConformances: []ProfileIdentifierType{
			ProfileIdentifierType_Software,
		},
	}}
	encodeDecode(t, doc)
}

func Test_externalID(t *testing.T) {
	doc := &SpdxDocument{ElementCollection: ElementCollection{
		Elements: ElementList{
			NewExternalIRI("http://someplace.org/ac7b643f0b2d"),
		},
	}}
	encodeDecode(t, doc)
}

// encodeDecode encodes to JSON, decodes from the JSON, and re-encodes in JSON to validate nothing is lost
func encodeDecode[T comparable](t *testing.T, obj T) T {
	// serialization:
	buf := bytes.Buffer{}
	err := context().ToJSON(&buf, obj)
	if err != nil {
		t.Fatal(err)
	}

	json1 := buf.String()
	t.Logf("--------- initial JSON: ----------\n%s\n\n", json1)

	// deserialization:
	graph, err := context().FromJSON(strings.NewReader(json1))
	var got T
	for _, entry := range graph {
		if e, ok := entry.(T); ok {
			got = e
			break
		}
	}

	var empty T
	if got == empty {
		t.Fatalf("did not find object in graph, json: %s", json1)
	}

	diff := cmp.Diff(obj, got, cmpopts.IgnoreUnexported((*T)(nil)))
	if diff != "" {
		t.Fatal(diff)
	}

	return got
}
