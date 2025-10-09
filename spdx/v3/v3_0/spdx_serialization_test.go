package v3_0

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/spdx/tools-golang/spdx/v3/internal"
)

func Test_spdxExportImportExport(t *testing.T) {
	doc := SpdxDocument{
		DataLicense: nil,
		Imports:     nil,
	}

	doc.ID = "new-id"

	agent := &SoftwareAgent{
		Name:    "some-agent",
		Summary: "summary",
	}
	c := &CreationInfo{
		Comment: "some-comment",
		Created: time.Now(),
		CreatedBy: AgentList{
			agent,
		},
		CreatedUsing: []AnyTool{
			&Tool{
				ExternalIdentifiers: ExternalIdentifierList{
					&ExternalIdentifier{
						Type:       ExternalIdentifierType_Cpe23,
						Identifier: "cpe:2.3:a:myvendor:my-product:*:*:*:*:*:*:*:*",
					},
				},
				Name: "not-tools-golang",
			},
		},
		SpecVersion: "",
	}
	agent.CreationInfo = c

	// add a package

	pkg1 := &Package{
		Name:         "some-package-1",
		CreationInfo: c,
		Version:      "1.2.3",
	}
	pkg2 := &Package{
		Name:         "some-package-2",
		CreationInfo: c,
		Version:      "2.4.5",
	}
	doc.RootElements = append(doc.RootElements, pkg2)

	file1 := &File{
		Name:         "/bin/bash",
		CreationInfo: c,
	}
	doc.RootElements = append(doc.RootElements, file1)

	// add relationships

	doc.RootElements = append(doc.RootElements,
		&Relationship{
			CreationInfo: c,
			From:         file1,
			Type:         RelationshipType_Contains,
			To: ElementList{
				pkg1,
				pkg2,
			},
		},
	)

	doc.RootElements = append(doc.RootElements,
		&Relationship{
			CreationInfo: c,
			From:         pkg1,
			Type:         RelationshipType_DependsOn,
			To: ElementList{
				pkg2,
			},
		},
	)

	doc.RootElements = append(doc.RootElements,
		&AIPackage{
			CreationInfo: c,
			TypeOfModels: []string{"a model"},
		},
	)

	got := encodeDecode(t, &doc)

	// some basic verification:

	var pkgs PackageList
	for _, rel := range got.RootElements.Relationships() {
		if rel.GetType() != RelationshipType_Contains {
			continue
		}
		if from, ok := rel.GetFrom().(AnyFile); ok {
			if from.GetName() == "/bin/bash" {
				for _, pkg := range rel.GetTo().Packages() {
					pkgs = append(pkgs, pkg)
				}
			}
		}
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
	doc := &SpdxDocument{
		ProfileConformances: []ProfileIdentifierType{
			ProfileIdentifierType_Software,
		},
	}
	encodeDecode(t, doc)
}

func Test_externalID(t *testing.T) {
	doc := &SpdxDocument{
		Elements: ElementList{
			// FIXME update the ExtenralIRI for flat struct generation
			//NewExternalIRI("http://someplace.org/ac7b643f0b2d"),
		},
	}
	encodeDecode(t, doc)
}

// encodeDecode encodes to JSON, decodes from the JSON and compares the decoded struct against the input
func encodeDecode[T comparable](t *testing.T, obj T) T {
	// serialization:
	buf := bytes.Buffer{}
	err := internal.ToJSON("https://spdx.org/rdf/3.0.1/spdx-context.jsonld", context(), obj, &buf)
	//err := context().ToJSON(&buf, obj)
	if err != nil {
		t.Fatal(err)
	}

	json1 := buf.String()

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

	diff := cmp.Diff(obj, got, diffOpts()...)
	if diff != "" {
		t.Fatal(diff)
	}

	return got
}
