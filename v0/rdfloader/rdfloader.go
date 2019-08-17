package rdfloader

import (
	"fmt"

	"github.com/spdx/tools-golang/v0/spdx"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
)

func Reader2_1(input string) *spdx.Document2_1 {

	spdxdoc, sp, err := Parse(input)
	if err != nil {
		fmt.Errorf("Parsing Error")
	}
	doc2v1 := rdf2v1.TransferDocument(spdxdoc, sp)
	if doc2v1 == nil {
		fmt.Errorf("Translation Error")
	}
	return doc2v1
}

func Parse(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
