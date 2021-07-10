// Package saver2v2 contains functions to render and write a json
// formatted version of an in-memory SPDX document and its sections
// (version 2.2).
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdxlib"
)

// RenderDocument2_2 is the main entry point to take an SPDX in-memory
// Document (version 2.2), and render it to the received *bytes.Buffer.
// It is only exported in order to be available to the jsonsaver package,
// and typically does not need to be called by client code.
func RenderDocument2_2(doc *spdx.Document2_2, buf *bytes.Buffer) error {

	fmt.Fprintln(buf, "{")
	// start to parse the creationInfo
	if doc.CreationInfo == nil {
		return fmt.Errorf("document had nil CreationInfo section")
	}
	renderCreationInfo2_2(doc.CreationInfo, buf)

	// parse otherlicenses from sodx struct to json
	if doc.OtherLicenses != nil {
		renderOtherLicenses2_2(doc.OtherLicenses, buf)
	}

	// parse document level annotations
	if doc.Annotations != nil {
		ann, _ := renderAnnotations2_2(doc.Annotations, spdx.MakeDocElementID("", string(doc.CreationInfo.SPDXIdentifier)))
		annotationjson, _ := json.Marshal(ann)
		fmt.Fprintf(buf, "\"%s\": %s ,", "annotations", annotationjson)
	}

	// parse document namespace
	if doc.CreationInfo.DocumentNamespace != "" {
		fmt.Fprintf(buf, "\"%s\": \"%s\",", "documentNamespace", doc.CreationInfo.DocumentNamespace)
	}

	// parse document describes
	describes, _ := spdxlib.GetDescribedPackageIDs2_2(doc)
	if describes != nil {
		var describesID []string
		for _, v := range describes {
			describesID = append(describesID, spdx.RenderElementID(v))
		}
		describesjson, _ := json.Marshal(describesID)
		fmt.Fprintf(buf, "\"%s\": %s,", "documentDescribes", describesjson)
	}

	// parse packages from spdx to json
	if doc.Packages != nil {
		renderPackage2_2(doc, buf)
	}
	// parse files from spdx to json
	if doc.UnpackagedFiles != nil {
		renderfiles2_2(doc, buf)
	}

	// parse relationships  from spdx to json
	if doc.Relationships != nil {
		rels, _ := renderRelationships2_2(doc.Relationships)
		relsjson, _ := json.Marshal(rels)
		fmt.Fprintf(buf, "\"%s\": %s ,", "relationships", relsjson)
	}

	// parsing ends
	buf.WriteRune('}')
	// remove the pattern ",}" from the json
	final := bytes.ReplaceAll(buf.Bytes(), []byte(",}"), []byte("}"))
	// indent the json
	var b []byte
	newbuf := bytes.NewBuffer(b)
	err := json.Indent(newbuf, final, "", "\t")
	if err != nil {
		return err
	}
	str := newbuf.String()
	logger := log.Default()
	logger.Fatal(str)
	return nil
}
