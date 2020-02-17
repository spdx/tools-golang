// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *licensediff*, *tvloader*

// This example demonstrates loading two SPDX tag-value files from disk into
// memory, and generating a diff of the concluded licenses for Files in the
// first-listed Packages in each document.
// This is generally only useful when run with two SPDX documents that
// describe licenses for subsequent versions of the same set of files.

package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/licensediff"
	"github.com/spdx/tools-golang/tvloader"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("Usage: %v <spdx-file-first> <spdx-file-second>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 tag-value files <spdx-file-first> and <spdx-file-second>,\n")
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

	// try to load the first SPDX file's contents as a tag-value file, version 2.1
	docFirst, err := tvloader.Load2_1(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", filenameFirst, err)
		return
	}
	// make sure it has at least one Package
	if len(docFirst.Packages) < 1 {
		fmt.Printf("Error, no packages found in SPDX file %s\n", filenameFirst)
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

	// try to load the second SPDX file's contents as a tag-value file, version 2.1
	docSecond, err := tvloader.Load2_1(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", filenameSecond, err)
		return
	}
	// make sure it has at least one Package
	if len(docSecond.Packages) < 1 {
		fmt.Printf("Error, no packages found in SPDX file %s\n", filenameSecond)
		return
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded second SPDX file %s\n\n", filenameSecond)

	// extract the first-listed package from each Document.
	// (there could be more than one Package within a Document. For
	// purposes of this example, we'll just look at the first one.)
	// we've already confirmed above that each has at least one Package.
	p1 := docFirst.Packages[0]
	p2 := docSecond.Packages[0]

	// now, run a diff between the two
	pairs, err := licensediff.MakePairs(p1, p2)
	if err != nil {
		fmt.Printf("Error generating licensediff pairs: %v\n", err)
		return
	}

	// take the pairs and turn them into a more structured results set
	resultSet, err := licensediff.MakeResults(pairs)
	if err != nil {
		fmt.Printf("Error generating licensediff results set: %v\n", err)
		return
	}

	// print some information about the results
	fmt.Printf("Files in first only: %d\n", len(resultSet.InFirstOnly))
	fmt.Printf("Files in second only: %d\n", len(resultSet.InSecondOnly))
	fmt.Printf("Files in both with different licenses: %d\n", len(resultSet.InBothChanged))
	fmt.Printf("Files in both with same licenses: %d\n", len(resultSet.InBothSame))
}
