// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
// Run project: go run exampleRDFLoader.go ../../testdata/spdx-examples/rdf/SPDXRdfExample.rdf
package main

import (
	"fmt"
	"github.com/spdx/tools-golang/rdfloader"
	"os"
	"strings"
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
		_ = fmt.Errorf("%v", ok)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		_ = fmt.Errorf("error opening File: %s", err)
		return
	}

	// loading the spdx-document
	doc, err := rdfloader.Load2_2(file)
	if err != nil {
		fmt.Println(fmt.Errorf("error parsing given spdx document: %s", err))
		os.Exit(1)
	}

	// Printing some of the document Information
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Some Attributes of the Document:")
	fmt.Printf("Document Name:         %s\n", doc.CreationInfo.DocumentName)
	fmt.Printf("DataLicense:           %s\n", doc.CreationInfo.DataLicense)
	fmt.Printf("Document NameSpace:    %s\n", doc.CreationInfo.DocumentNamespace)
	fmt.Printf("SPDX Document Version: %s\n", doc.CreationInfo.SPDXVersion)
	fmt.Println(strings.Repeat("=", 80))
}
