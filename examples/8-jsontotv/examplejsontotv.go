// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *jsonparser2v2*, *tvsaver*

// This example demonstrates loading an SPDX json from disk into memory,
// and then re-saving it to a different file on disk in tag-value format .
// Run project: go run examplejsontotv.go ../sample-docs/json/SPDXJSONExample-v2.2.spdx.json example.spdx
package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/jsonloader"
	"github.com/spdx/tools-golang/tvsaver"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("Usage: %v <json-file-in> <spdx-file-out>\n", args[0])
		fmt.Printf("  Load SPDX 2.2 tag-value file <spdx-file-in>, and\n")
		fmt.Printf("  save it out to <spdx-file-out>.\n")
		return
	}

	// open the SPDX file
	fileIn := args[1]
	r, err := os.Open(fileIn)
	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", fileIn, err)
		return
	}
	defer r.Close()

	// try to load the SPDX file's contents as a json file, version 2.2
	doc, err := jsonloader.Load2_2(r)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v", args[1], err)
		return
	}

	// if we got here, the file is now loaded into memory.
	fmt.Printf("Successfully loaded %s\n", args[1])

	// we can now save it back to disk, using tvsaver.

	// create a new file for writing
	fileOut := args[2]
	w, err := os.Create(fileOut)
	if err != nil {
		fmt.Printf("Error while opening %v for writing: %v", fileOut, err)
		return
	}
	defer w.Close()

	// try to save the document to disk as an SPDX tag-value file, version 2.2
	err = tvsaver.Save2_2(doc, w)
	if err != nil {
		fmt.Printf("Error while saving %v: %v", fileOut, err)
		return
	}

	// it worked
	fmt.Printf("Successfully saved %s\n", fileOut)
}
