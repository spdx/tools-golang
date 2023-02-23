// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package yaml

import (
	"io"

	"sigs.k8s.io/yaml"

	"github.com/spdx/tools-golang/convert"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdxlib"
)

// Write takes an SPDX Document and an io.Writer, and writes the document to the writer in YAML format.
func Write(doc common.AnyDocument, w io.Writer) error {
	var targetDoc spdx.Document
	err := convert.Document(doc, &targetDoc)
	if err != nil {
		return err
	}

	err = spdxlib.PopulateJsonSchemaFields(&targetDoc)
	if err != nil {
		return err
	}

	err = convert.Document(&targetDoc, doc)
	if err != nil {
		return err
	}
	buf, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
