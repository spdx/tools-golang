package rdf

import (
	"io"

	"github.com/spdx/gordf/rdfloader"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/spdx/tools-golang/spdx/v2/v2_3/rdf/reader"
)

// Takes in a file Reader and returns the pertaining spdx document
// or the error if any is encountered while setting the doc.
func Read(content io.Reader) (*spdx.Document, error) {
	var rdfParserObj, err = rdfloader.LoadFromReaderObject(content)
	if err != nil {
		return nil, err
	}

	doc, err := reader.LoadFromGoRDFParser(rdfParserObj)
	return doc, err
}
