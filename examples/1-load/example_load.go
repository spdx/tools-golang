// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *tagvalue*, *spdx*

// This example demonstrates loading an SPDX tag-value file from disk into
// memory, and printing some of its contents to standard output.
// Run project:  go run example_load.go ../sample-docs/tv/hello.spdx

package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/spdxlib"
	"github.com/spdx/tools-golang/tagvalue"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX tag-value file <spdx-file-in>, and\n")
		fmt.Printf("  print a portion of its contents.\n")
		return
	}

	// open the SPDX file
	filename := args[1]
	r, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filename, err)
		return
	}
	defer r.Close()

	// try to load the SPDX file's contents as a tag-value file
	doc, err := tagvalue.Read(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", filename, err)
		return
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded %s\n\n", filename)

	// we can now take a look at its contents via the various data
	// structures representing the SPDX document's sections.

	// print the struct containing the SPDX file's Creation Info section data
	fmt.Printf("==============\n")
	fmt.Printf("Creation info:\n")
	fmt.Printf("==============\n")
	fmt.Printf("%#v\n\n", doc.CreationInfo)

	// check whether the SPDX file has at least one package that it describes
	pkgIDs, err := spdxlib.GetDescribedPackageIDs(doc)
	if err != nil {
		fmt.Printf("Unable to get describe packages from SPDX document: %v\n", err)
		return
	}

	if len(pkgIDs) == 0 {
		return
	}

	// it does, so we'll go through each one
	for _, pkg := range doc.Packages {
		var documentDescribesPackage bool
		for _, describedPackageID := range pkgIDs {
			if pkg.PackageSPDXIdentifier == describedPackageID {
				documentDescribesPackage = true
				break
			}
		}

		if !documentDescribesPackage {
			continue
		}

		pkgID := pkg.PackageSPDXIdentifier

		// check whether the package had its files analyzed
		if !pkg.FilesAnalyzed {
			fmt.Printf("Package %s (%s) had FilesAnalyzed: false\n", string(pkgID), pkg.PackageName)
			continue
		}

		// also check whether the package has any files present
		if pkg.Files == nil || len(pkg.Files) < 1 {
			fmt.Printf("Package %s (%s) has no Files\n", string(pkgID), pkg.PackageName)
			continue
		}

		// if we got here, there's at least one file
		// print the filename and license info for the first 50
		fmt.Printf("============================\n")
		fmt.Printf("Package %s (%s)\n", string(pkgID), pkg.PackageName)
		fmt.Printf("File info (up to first 50):\n")
		i := 1
		for _, f := range pkg.Files {
			// note that these will be in random order, since we're pulling
			// from a map. if we care about order, we should first pull the
			// IDs into a slice, sort it, and then print the ordered files.
			fmt.Printf("- File %d: %s\n", i, f.FileName)
			fmt.Printf("    License from file: %v\n", f.LicenseInfoInFiles)
			fmt.Printf("    License concluded: %v\n", f.LicenseConcluded)
			i++
			if i > 50 {
				break
			}
		}
	}
}
