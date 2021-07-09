// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func renderAnnotations2_2(annotations []*spdx.Annotation2_2, buf *bytes.Buffer) error {

	var ann []interface{}
	for _, v := range annotations {
		annotation := make(map[string]interface{})
		annotation["annotationDate"] = v.AnnotationDate
		annotation["annotationType"] = v.AnnotationType
		annotation["annotator"] = fmt.Sprintf("%s: %s", v.AnnotatorType, v.Annotator)
		if v.AnnotationComment != "" {
			annotation["comment"] = v.AnnotationComment
		}
		ann = append(ann, annotation)
	}
	annotationjson, _ := json.Marshal(ann)
	fmt.Fprintf(buf, "\"%s\": %s ,", "annotations", annotationjson)

	return nil
}
