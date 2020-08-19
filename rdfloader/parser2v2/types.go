// copied from tvloader/parser2v2/types.go
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v2

import (
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
)

type rdfParser2_2 struct {
	// fields associated with gordf project which
	// will be required by rdfloader
	gordfParserObj *gordfParser.Parser
	nodeToTriples  map[string][]*gordfParser.Triple

	// document into which data is being parsed
	doc *spdx.Document2_2

	// map of packages and files.
	files            map[spdx.ElementID]*spdx.File2_2
	assocWithPackage map[spdx.ElementID]bool
	packages         map[spdx.ElementID]*spdx.Package2_2

	// mapping of nodeStrings to parsed object to save double computation.
	cache map[string]interface{}
}
