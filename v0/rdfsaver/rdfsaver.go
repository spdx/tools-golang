// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfsaver

import (
	"fmt"
	"os"

	"github.com/spdx/tools-golang/v0/spdx"

	"github.com/spdx/tools-golang/v0/rdfloader/rdf2v1"
	"github.com/spdx/tools-golang/v0/rdfloader/rdfsaver2v1"
)

func Saver2_1(doc2v1 *spdx.Document2_1) error {

	newdoc2v1 := rdf2v1.CollectDocument(doc2v1)
	newsn2v1 := rdf2v1.CollectSnippets(doc2v1)
	output := os.Stdout

	errdoc := rdfsaver2v1.Write(output, newdoc2v1, newsn2v1)

	if errdoc != nil {
		fmt.Errorf("Parsing Error")
	}
	return errdoc
}
