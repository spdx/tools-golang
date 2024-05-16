package v3_0_test

import (
	"fmt"
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v3/v3_0/embedded_second"
)

func Test_makeAnSpdxDocument(t *testing.T) {
	// creating new documents: 2 packages found from 1 file with 2 relationships

	// embedded nesting is not ideal when instantiating:
	pkg1 := &spdx.PackageData{
		SoftwareArtifactData: spdx.SoftwareArtifactData{
			ArtifactData: spdx.ArtifactData{
				ElementData: spdx.ElementData{
					SpdxID: "pkg-1",
					Name:   "pkg-1",
				},
			},
		},
		PackageVersion: "1.0.0",
	}

	// can just reference properties directly, once an object exists:
	pkg2 := &spdx.PackageData{}
	pkg2.SpdxID = "package-2"
	pkg2.Name = "package-2"
	pkg2.PackageVersion = "2.0.0"

	file1 := &spdx.FileData{}
	file1.SpdxID = "file-1"
	file1.Name = "file-1"
	file1.ContentType = "text/plain"

	file1containsPkg1 := &spdx.RelationshipData{
		RelationshipType: "CONTAINS",
		From:             file1,
		To:               []spdx.Element{pkg1},
	}

	pkg1dependsOnFile1 := &spdx.RelationshipData{
		RelationshipType: "DEPENDS_ON",
		From:             pkg1,
		To:               []spdx.Element{pkg2},
	}

	doc := spdx.SpdxDocumentData{
		ElementCollectionData: spdx.ElementCollectionData{
			ElementData: spdx.ElementData{
				SpdxID: "spdx-document",
			},
			Elements: []spdx.Element{
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
		if e, ok := e.(spdx.Package); ok {
			e.AsPackage().Name = "updated-name"
		}
	}
	fmt.Printf("%#v\n", doc)
}
