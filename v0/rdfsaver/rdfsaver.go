package rdfsaver

import (
	"os"

	"github.com/spdx/tools-golang/v0/rdfsaver/rdfsaver2v1"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
)

func Saver2_1(spdxdoc *rdf2v1.Document, sp *rdf2v1.Snippet) error {

	output := os.Stdout
	err := rdfsaver2v1.Write(output, spdxdoc, sp)
	return err
}
