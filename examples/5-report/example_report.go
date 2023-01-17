// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *reporter*, *tagvalue*

// This example demonstrates loading an SPDX tag-value file from disk into memory,
// generating a basic report listing counts of the concluded licenses for its
// files, and printing the report to standard output.
// Run project: go run example_report.go ../sample-docs/tv/hello.spdx
package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/reporter"
	"github.com/spdx/tools-golang/spdxlib"
	"github.com/spdx/tools-golang/tagvalue"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX tag-value file <spdx-file-in>, and\n")
		fmt.Printf("  generate and print a report of its concluded licenses.\n")
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
			return
		}

		// also check whether the package has any files present
		if pkg.Files == nil || len(pkg.Files) < 1 {
			fmt.Printf("Package %s (%s) has no Files\n", string(pkgID), pkg.PackageName)
			return
		}

		// if we got here, there's at least one file
		// generate and print a report of the Package's Files' LicenseConcluded
		// values, sorted by # of occurrences
		fmt.Printf("============================\n")
		fmt.Printf("Package %s (%s)\n", string(pkgID), pkg.PackageName)
		err = reporter.Generate(pkg, os.Stdout)
		if err != nil {
			fmt.Printf("Error while generating report: %v\n", err)
		}
	}
}
