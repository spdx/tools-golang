package rdf

import (
	"io"

	"github.com/spdx/gordf/rdfloader"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
	"github.com/spdx/tools-golang/spdx/v2/v2_2/rdf/reader"
)

// Takes in a file Reader and returns the pertaining spdx document
// or the error if any is encountered while setting the doc.
func Read(content io.Reader) (*v2_2.Document, error) {
	var rdfParserObj, err = rdfloader.LoadFromReaderObject(content)
	if err != nil {
		return nil, err
	}

	doc, err := reader.LoadFromGoRDFParser(rdfParserObj)
	return doc, err
}
