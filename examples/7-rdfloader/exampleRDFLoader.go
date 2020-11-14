// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"fmt"
	"github.com/spdx/tools-golang/rdfloader"
	"os"
	"strings"
)

func getFilePathFromUser() string {
	if len(os.Args) == 1 {
		// user hasn't specified the rdf file path
		panic("kindly provide path of the rdf file to be loaded as a spdx-document while running this file")
	}
	return os.Args[1]
}

func main() {
	// example to use the rdfLoader.
	filePath := getFilePathFromUser()
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
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
