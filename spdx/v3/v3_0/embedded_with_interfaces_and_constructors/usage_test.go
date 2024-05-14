package v3_0_test

import (
	"fmt"
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v3/v3_0/embedded_with_interfaces_and_constructors"
)

func Test_makeAnSpdxDocument(t *testing.T) {
	// creating new documents: 2 packages found from 1 file with 2 relationships

	pkg1 := spdx.NewPackage(spdx.PackageProps{
		SpdxID:         "package-1",
		Name:           "package-1",
		PackageVersion: "1.0.0",
	})

	pkg2 := spdx.NewPackage(spdx.PackageProps{
		SpdxID:         "package-2",
		Name:           "package-2",
		PackageVersion: "2.0.0",
	})

	file1 := spdx.NewFile(spdx.FileProps{
		SpdxID:      "file-1",
		Name:        "file-1",
		ContentType: "text/plain",
	})

	file1containsPkg1 := spdx.NewRelationship(spdx.RelationshipProps{
		RelationshipType: "CONTAINS",
		From:             file1,
		To:               []spdx.Element{pkg1},
	})

	pkg1dependsOnFile1 := spdx.NewRelationship(spdx.RelationshipProps{
		RelationshipType: "DEPENDS_ON",
		From:             pkg1,
		To:               []spdx.Element{pkg2},
	})

	doc := spdx.NewSpdxDocument(spdx.SpdxDocumentProps{
		SpdxID: "spdx-document",
		Element: []spdx.Element{
			pkg1,
			pkg2,
			pkg1dependsOnFile1,
			file1containsPkg1,
		},
	})
	fmt.Printf("%#v\n", doc)

	// working with existing documents

	for _, e := range doc.Elements() {
		if e, ok := e.(spdx.Package); ok {
			e.SetName("updated-name")
		}
	}
	fmt.Printf("%#v\n", doc)
}
