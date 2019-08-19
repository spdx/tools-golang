package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/v0/rdfloader"
	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
	"github.com/spdx/tools-golang/v0/spdx"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 tag-value file <spdx-file-in>, and\n")
		fmt.Printf("  print a portion of its contents.\n")
		return
	}

	// storing the input
	input := args[1]

	// try to load the SPDX file's contents as a rdf file, version 2.1
	doc2v1 := rdfloader.Reader2_1(input)

	fmt.Printf("\n\n\n%#v\n\n\n", doc2v1)

	// if we got here, the file is now loaded into memory.
	// we can now take a look at its contents via the various data
	// structures representing the SPDX document's sections.

	// print the struct containing the SPDX file's Creation Info section data
	fmt.Printf("==============\n")
	fmt.Printf("Creation info:\n")
	fmt.Printf("==============\n")
	fmt.Printf("%#v\n\n", doc2v1.CreationInfo)

	// check whether the SPDX file has at least one package
	fmt.Printf("==============\n")
	fmt.Printf("Packages:\n")
	fmt.Printf("==============\n")

	var packages []*spdx.Package2_1

	for _, pkg := range doc2v1.Packages {
		if pkg != nil {
			packages = append(packages, pkg)
			fmt.Printf("Package %s found.\n", pkg.PackageName)
		} else {
			fmt.Printf("No packages found in SPDX document\n")
			return
		}
	}

	// take a package and check whether the package has any files present
	pkg := packages[0]

	if pkg.Files == nil || len(pkg.Files) < 1 {
		fmt.Printf("No Files found in first Package (%s)\n", pkg.PackageName)
		return
	}

	// if we got here, there's at least one file
	// print the filename and license info for the first 50
	fmt.Printf("============================\n")
	fmt.Printf("Files info (up to first 50):\n")
	fmt.Printf("============================\n")
	i := 1
	for _, f := range pkg.Files {
		fmt.Printf("File %d: %s\n", i, f.FileName)
		fmt.Printf("  License from file: %v\n", f.LicenseInfoInFile)
		i++
		if i > 50 {
			break
		}
	}
}
func Parse2_1(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
