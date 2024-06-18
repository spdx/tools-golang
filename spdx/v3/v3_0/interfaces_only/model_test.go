package v3_0

import (
	"fmt"
	"testing"
)

func Test_makeAnSpdxDocument(t *testing.T) {
	// creating new documents: 2 packages found from 1 file with 2 relationships

	// must call setters
	pkg1 := NewPackage()
	_ = pkg1.SetSpdxID("package-1")
	_ = pkg1.SetName("package-1")
	_ = pkg1.SetPackageVersion("1.0.0")

	pkg2 := NewPackage()
	_ = pkg2.SetSpdxID("package-2")
	_ = pkg2.SetName("package-2")
	_ = pkg2.SetPackageVersion("2.0.0")

	file1 := NewFile()
	_ = file1.SetSpdxID("file-1")
	_ = file1.SetName("file-1")
	_ = file1.SetContentType("text/plain")

	file1containsPkg1 := NewRelationship()
	_ = file1containsPkg1.SetRelationshipType("CONTAINS")
	_ = file1containsPkg1.SetFrom(file1)
	_ = file1containsPkg1.SetTo([]Element{pkg1})

	pkg1dependsOnFile1 := NewRelationship()
	_ = pkg1dependsOnFile1.SetRelationshipType("DEPENDS_ON")
	_ = pkg1dependsOnFile1.SetFrom(pkg1)
	_ = pkg1dependsOnFile1.SetTo([]Element{pkg2})

	doc := NewSpdxDocument()
	_ = doc.SetSpdxID("spdx-document")
	_ = doc.SetElements([]Element{
		pkg1,
		pkg2,
		pkg1dependsOnFile1,
		file1containsPkg1,
	})
	fmt.Printf("%#v\n", doc)

	// working with existing documents

	for _, e := range doc.Elements() {
		if e, ok := e.(Package); ok {
			_ = e.SetName("updated-name")
		}
	}
	fmt.Printf("%#v\n", doc)
}
