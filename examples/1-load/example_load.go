// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *tvloader*, *spdx*

// This example demonstrates loading an SPDX tag-value file from disk into
// memory, and printing some of its contents to standard output.

package main

import (
	"fmt"
	"os"

	"github.com/swinslow/spdx-go/v0/tvloader"
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

	// we can now take a look at its contents via the various data
	// structures representing the SPDX document's sections.

	// print the struct containing the SPDX file's Creation Info section data
	fmt.Printf("==============\n")
	fmt.Printf("Creation info:\n")
	fmt.Printf("==============\n")
	fmt.Printf("%#v\n\n", doc.CreationInfo)

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
	// print the filename and license info for the first 50
	fmt.Printf("============================\n")
	fmt.Printf("Files info (up to first 50):\n")
	fmt.Printf("============================\n")
	i := 1
	for _, f := range pkg.Files {
		fmt.Printf("File %d: %s\n", i, f.FileName)
		fmt.Printf("  License from file: %v\n", f.LicenseInfoInFile)
		fmt.Printf("  License concluded: %v\n", f.LicenseConcluded)
		i++
		if i > 50 {
			break
		}
	}
}
