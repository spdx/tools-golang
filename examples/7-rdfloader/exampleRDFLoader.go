// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
// Run project: go run exampleRDFLoader.go ../sample-docs/rdf/SPDXRdfExample-v2.2.spdx.rdf
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spdx/tools-golang/rdf"
)

func getFilePathFromUser() (string, error) {
	if len(os.Args) == 1 {
		// user hasn't specified the rdf file path
		return "", fmt.Errorf("kindly provide path of the rdf file to be loaded as a spdx-document while running this file")
	}
	return os.Args[1], nil
}

func main() {
	// example to use the rdfLoader.
	filePath, ok := getFilePathFromUser()
	if ok != nil {
		fmt.Println(fmt.Errorf("%v", ok))
		os.Exit(1)
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("error opening File: %s", err))
		os.Exit(1)
	}

	// loading the spdx-document
	doc, err := rdf.Read(file)
	if err != nil {
		fmt.Println(fmt.Errorf("error parsing given spdx document: %s", err))
		os.Exit(1)
	}

	// Printing some of the document Information
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Some Attributes of the Document:")
	fmt.Printf("Document Name:         %s\n", doc.DocumentName)
	fmt.Printf("DataLicense:           %s\n", doc.DataLicense)
	fmt.Printf("Document Namespace:    %s\n", doc.DocumentNamespace)
	fmt.Printf("SPDX Version:          %s\n", doc.SPDXVersion)
	fmt.Println(strings.Repeat("=", 80))
}
