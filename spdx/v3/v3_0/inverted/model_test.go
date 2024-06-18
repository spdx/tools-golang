package v3_0

import (
	"fmt"
	"testing"
)

func Test_model(t *testing.T) {
	// creating new documents: 2 packages found from 1 file with 2 relationships

	pkg1 := &Element{
		SpdxID:           "package-1",
		Name:             "package-1",
		Artifact:         &Artifact{}, // how would a user know to do this?
		SoftwareArtifact: &SoftwareArtifact{},
		Package: &Package{
			PackageVersion: "1.0.0",
		},
	}

	pkg2 := NewPackage(PackageProps{
		SpdxID:         "package-2",
		Name:           "package-2",
		PackageVersion: "2.0.0",
	})

	file1 := NewFile(FileProps{
		SpdxID:      "file-2",
		Name:        "file-2",
		ContentType: "text/plain",
	})

	file1containsPkg1 := NewRelationship(RelationshipProps{
		RelationshipType: "CONTAINS",
		From:             file1,
		To:               []*Element{pkg1},
	})

	pkg1dependsOnFile1 := NewRelationship(RelationshipProps{
		RelationshipType: "DEPENDS_ON",
		From:             pkg1,
		To:               []*Element{pkg2},
	})

	doc := NewSpdxDocument(SpdxDocumentProps{
		SpdxID: "spdx-document",
		Elements: []*Element{
			pkg1,
			pkg2,
			pkg1dependsOnFile1,
			file1containsPkg1,
		},
	})
	fmt.Printf("%#v\n", doc)

	// working with existing documents

	for _, e := range doc.Elements {
		if p := e.Package; p != nil {
			e.Name = "updated-name"
		}
	}
	fmt.Printf("%#v\n", doc)
}
