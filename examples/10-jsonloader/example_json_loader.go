// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *jsonparser2v2*

// This example demonstrates loading an SPDX json from disk into memory,
// and then logging out some attributes to the console .
// Run project: go run example_json_loader.go ../sample-docs/json/SPDXJSONExample-v2.2.spdx.json example.spdx
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spdx/tools-golang/jsonloader"
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

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Some Attributes of the Document:")
	fmt.Printf("Document Name:         %s\n", doc.CreationInfo.DocumentName)
	fmt.Printf("DataLicense:           %s\n", doc.CreationInfo.DataLicense)
	fmt.Printf("Document NameSpace:    %s\n", doc.CreationInfo.DocumentNamespace)
	fmt.Printf("SPDX Document Version: %s\n", doc.CreationInfo.SPDXVersion)
	fmt.Println(strings.Repeat("=", 80))
}
