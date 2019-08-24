// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/v0/rdfloader"
	"github.com/spdx/tools-golang/v0/rdfsaver"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
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

	// storing the input
	input := args[1]

	// try to load the SPDX file's contents as a rdf file, version 2.1
	doc2v1 := rdfloader.Reader2_1(input)

	err := rdfsaver.Saver2_1(doc2v1)
	if err != nil {
		fmt.Printf("Error while saving the document.")
		return
	}

}

func Parse2_1(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
