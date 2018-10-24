// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *reporter*, *tvloader*

// This example demonstrates loading an SPDX tag-value file from disk into memory,
// generating a basic report listing counts of the concluded licenses for its
// files, and printing the report to standard output.

package main

import (
	"fmt"
	"os"

	"github.com/swinslow/spdx-go/v0/reporter"
	"github.com/swinslow/spdx-go/v0/tvloader"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 tag-value file <spdx-file-in>, and\n")
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

	// try to load the SPDX file's contents as a tag-value file, version 2.1
	doc, err := tvloader.Load2_1(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", filename, err)
		return
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded %s\n\n", filename)

	// check whether the SPDX file has at least one package
	if doc.Packages == nil || len(doc.Packages) < 1 {
		fmt.Printf("No packages found in SPDX document\n")
		return
	}

	// it does, so we'll choose the first one
	pkg := doc.Packages[0]

	// check whether the package had its files analyzed
	if !pkg.FilesAnalyzed {
		fmt.Printf("First Package (%s) had FilesAnalyzed: false\n", pkg.PackageName)
		return
	}

	// also check whether the package has any files present
	if pkg.Files == nil || len(pkg.Files) < 1 {
		fmt.Printf("No Files found in first Package (%s)\n", pkg.PackageName)
		return
	}

	// if we got here, there's at least one file
	// generate and print a report of the Package's Files' LicenseConcluded
	// values, sorted by # of occurrences
	err = reporter.Generate(pkg, os.Stdout)
	if err != nil {
		fmt.Printf("Error while generating report: %v\n", err)
	}
}
