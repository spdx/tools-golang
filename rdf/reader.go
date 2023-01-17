// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf

import (
	"errors"
	"fmt"
	"io"

	"github.com/spdx/gordf/rdfloader"
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/convert"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
	v2_2_reader "github.com/spdx/tools-golang/spdx/v2/v2_2/rdf/reader"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	v2_3_reader "github.com/spdx/tools-golang/spdx/v2/v2_3/rdf/reader"
)

// Read takes an io.Reader and returns a fully-parsed current model SPDX Document
// or an error if any error is encountered.
func Read(content io.Reader) (*spdx.Document, error) {
	doc := spdx.Document{}
	err := ReadInto(content, &doc)
	return &doc, err
}

// ReadInto takes an io.Reader, reads in the SPDX document at the version provided
// and converts to the doc version
func ReadInto(content io.Reader, doc common.AnyDocument) error {
	if !convert.IsPtr(doc) {
		return fmt.Errorf("doc to read into must be a pointer")
	}
	var rdfParserObj, err = rdfloader.LoadFromReaderObject(content)
	if err != nil {
		return err
	}

	version, err := getSpdxVersion(rdfParserObj)
	if err != nil {
		return err
	}

	var data interface{}
	switch version {
	case v2_2.Version:
		data, err = v2_2_reader.LoadFromGoRDFParser(rdfParserObj)
	case v2_3.Version:
		data, err = v2_3_reader.LoadFromGoRDFParser(rdfParserObj)
	default:
		return fmt.Errorf("unsupported SPDX version: '%v'", version)
	}

	if err != nil {
		return err
	}

	return convert.Document(data.(common.AnyDocument), doc)
}

func getSpdxVersion(parser *gordfParser.Parser) (string, error) {
	version := ""
	for _, node := range parser.Triples {
		if node.Predicate.ID == "http://spdx.org/rdf/terms#specVersion" {
			version = node.Object.ID
			break
		}
	}
	if version == "" {
		return "", errors.New("unable to determine version from RDF document")
	}
	return version, nil
}
