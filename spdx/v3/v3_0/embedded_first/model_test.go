package v3_0

import (
	"fmt"
	"testing"
)

func Test_makeAnSpdxDocument(t *testing.T) {
	// creating new documents: 2 packages found from 1 file with 2 relationships

	// embedded nesting is not ideal when instantiating:
	pkg1 := &Package{
		SoftwareArtifact: SoftwareArtifact{
			Artifact: Artifact{
				Element: Element{
					SpdxID: "pkg-1",
					Name:   "pkg-1",
				},
			},
		},
		PackageVersion: "1.0.0",
	}

	// can just reference properties directly, once an object exists:
	pkg2 := &Package{}
	pkg2.SpdxID = "package-2"
	pkg2.Name = "package-2"
	pkg2.PackageVersion = "2.0.0"

	file1 := &File{}
	file1.SpdxID = "file-1"
	file1.Name = "file-1"
	file1.ContentType = "text/plain"

	file1containsPkg1 := &Relationship{
		RelationshipType: "CONTAINS",
		From:             file1,
		To:               []IElement{pkg1},
	}

	pkg1dependsOnFile1 := &Relationship{
		RelationshipType: "DEPENDS_ON",
		From:             pkg1,
		To:               []IElement{pkg2},
	}

	doc := SpdxDocument{
		ElementCollection: ElementCollection{
			Element: Element{
				SpdxID: "spdx-document",
			},
			Elements: []IElement{
				pkg1,
				pkg2,
				pkg1dependsOnFile1,
				file1containsPkg1,
			},
		},
	}
	fmt.Printf("%#v\n", doc)

	// working with existing documents

	for _, e := range doc.Elements {
		if e, ok := e.(IPackage); ok {
			e.AsPackage().Name = "updated-name"
		}
	}
	fmt.Printf("%#v\n", doc)
}
