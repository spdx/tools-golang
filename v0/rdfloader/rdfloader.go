package rdfloader

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
)

func Reader2_1() (rdf2v1.Document, rdf2v1.Snippet, error) {

	args := os.Args
	if len(args) != 2 {
		fmt.Errorf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Errorf("  Load SPDX 2.1 RDF file <spdx-file-in>, and\n")
		fmt.Errorf("  print its contents.\n")

	}
	var spdxdoc *rdf2v1.Document
	var sp *rdf2v1.Snippet
	var err error

	input := args[1]
	spdxdoc, sp, err = Parse(input)

	_, _ = spdxdoc, sp
	// var doc2v1 *spdx.Document2_1
	// doc2v1 = TransferDocument(spdxdoc, sp)
	if err != nil {
		fmt.Errorf("Parsing Error")
	}
	return *spdxdoc, *sp, err
}

func Parse(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
