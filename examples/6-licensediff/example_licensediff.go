// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *licensediff*, *tvloader*

// This example demonstrates loading two SPDX tag-value files from disk into
// memory, and generating a diff of the concluded licenses for Files in
// Packages with matching IDs in each document.
// This is generally only useful when run with two SPDX documents that
// describe licenses for subsequent versions of the same set of files, AND if
// they have the same identifier in both documents.
// Run project: go run example_licensediff.go ../sample-docs/tv/hello.spdx  ../sample-docs/tv/hello-modified.spdx
package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/licensediff"
	"github.com/spdx/tools-golang/spdxlib"
	"github.com/spdx/tools-golang/tvloader"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("Usage: %v <spdx-file-first> <spdx-file-second>\n", args[0])
		fmt.Printf("  Load SPDX 2.2 tag-value files <spdx-file-first> and <spdx-file-second>,\n")
		fmt.Printf("  run a diff between their concluded licenses, and print basic results.\n")
		return
	}

	// open the first SPDX file
	filenameFirst := args[1]
	r, err := os.Open(filenameFirst)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filenameFirst, err)
		return
	}
	defer r.Close()

	// try to load the first SPDX file's contents as a tag-value file, version 2.2
	docFirst, err := tvloader.Load2_2(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", filenameFirst, err)
		return
	}
	// check whether the SPDX file has at least one package that it describes
	pkgIDsFirst, err := spdxlib.GetDescribedPackageIDs2_2(docFirst)
	if err != nil {
		fmt.Printf("Unable to get describe packages from first SPDX document: %v\n", err)
		return
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded first SPDX file %s\n", filenameFirst)

	// open the second SPDX file
	filenameSecond := args[2]
	r, err = os.Open(filenameSecond)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filenameSecond, err)
		return
	}
	defer r.Close()

	// try to load the second SPDX file's contents as a tag-value file, version 2.2
	docSecond, err := tvloader.Load2_2(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", filenameSecond, err)
		return
	}
	// check whether the SPDX file has at least one package that it describes
	pkgIDsSecond, err := spdxlib.GetDescribedPackageIDs2_2(docSecond)
	if err != nil {
		fmt.Printf("Unable to get describe packages from second SPDX document: %v\n", err)
		return
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded second SPDX file %s\n\n", filenameSecond)

	// compare the described packages from each Document, by SPDX ID
	// go through the first set first, report if they aren't in the second set
	for _, pkgID := range pkgIDsFirst {
		fmt.Printf("================================\n")
		p1, okFirst := docFirst.Packages[pkgID]
		if !okFirst {
			fmt.Printf("Package %s has described relationship in first document but ID not found\n", string(pkgID))
			continue
		}
		fmt.Printf("Package %s (%s)\n", string(pkgID), p1.PackageName)
		p2, okSecond := docSecond.Packages[pkgID]
		if !okSecond {
			fmt.Printf("  Found in first document, not found in second\n")
			continue
		}

		// now, run a diff between the two
		pairs, err := licensediff.MakePairs2_2(p1, p2)
		if err != nil {
			fmt.Printf("  Error generating licensediff pairs: %v\n", err)
			continue
		}

		// take the pairs and turn them into a more structured results set
		resultSet, err := licensediff.MakeResults(pairs)
		if err != nil {
			fmt.Printf("  Error generating licensediff results set: %v\n", err)
			continue
		}

		// print some information about the results
		fmt.Printf("  Files in first only: %d\n", len(resultSet.InFirstOnly))
		fmt.Printf("  Files in second only: %d\n", len(resultSet.InSecondOnly))
		fmt.Printf("  Files in both with different licenses: %d\n", len(resultSet.InBothChanged))
		fmt.Printf("  Files in both with same licenses: %d\n", len(resultSet.InBothSame))
	}

	// now report if there are any package IDs in the second set that aren't
	// in the first
	for _, pkgID := range pkgIDsSecond {
		p2, okSecond := docSecond.Packages[pkgID]
		if !okSecond {
			fmt.Printf("================================\n")
			fmt.Printf("Package %s has described relationship in second document but ID not found\n", string(pkgID))
			continue
		}
		_, okFirst := docFirst.Packages[pkgID]
		if !okFirst {
			fmt.Printf("================================\n")
			fmt.Printf("Package %s (%s)\n", string(pkgID), p2.PackageName)
			fmt.Printf("  Found in second document, not found in first\n")
		}
	}
}
