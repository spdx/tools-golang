package v3_0

import (
	"fmt"
	"testing"

	v3_0 "github.com/spdx/tools-golang/spdx/v3/v3_0/interfaces_only"
)

func Test_makeAnSpdxDocument(t *testing.T) {
	// creating new documents: 2 packages found from 1 file with 2 relationships

	// must call setters
	pkg1, _ := NewPackage(PackageProps{
		SpdxID:         "package-1",
		Name:           "package-1",
		PackageVersion: "1.0.0",
	})

	pkg2, _ := NewPackage(PackageProps{
		SpdxID:         "package-2",
		Name:           "package-2",
		PackageVersion: "2.0.0",
	})

	file1, _ := NewFile(FileProps{
		SpdxID:      "file-1",
		Name:        "file-1",
		ContentType: "text/plain",
	})

	file1containsPkg1, _ := NewRelationship(RelationshipProps{
		RelationshipType: "CONTAINS",
		From:             file1,
		To:               []v3_0.Element{pkg1},
	})

	pkg1dependsOnFile1, _ := NewRelationship(RelationshipProps{
		RelationshipType: "DEPENDS_ON",
		From:             pkg1,
		To:               []v3_0.Element{pkg2},
	})

	doc, _ := NewSpdxDocument(SpdxDocumentProps{
		SpdxID: "spdx-document",
		Element: []v3_0.Element{
			pkg1,
			pkg2,
			pkg1dependsOnFile1,
			file1containsPkg1,
		},
	})
	fmt.Printf("%#v\n", doc)

	// working with existing documents

	for _, e := range doc.Elements() {
		if e, ok := e.(v3_0.Package); ok {
			_ = e.SetName("updated-name")
		}
	}
	fmt.Printf("%#v\n", doc)
}
