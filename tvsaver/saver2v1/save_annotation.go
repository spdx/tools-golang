// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v1

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/spdx"
)

func renderAnnotation2_1(ann *spdx.Annotation2_1, w io.Writer) error {
	if ann.Annotator != "" && ann.AnnotatorType != "" {
		fmt.Fprintf(w, "Annotator: %s: %s\n", ann.AnnotatorType, ann.Annotator)
	}
	if ann.Date != "" {
		fmt.Fprintf(w, "AnnotationDate: %s\n", ann.Date)
	}
	if ann.Type != "" {
		fmt.Fprintf(w, "AnnotationType: %s\n", ann.Type)
	}
	annIDStr := spdx.RenderDocElementID(ann.SPDXIdentifier)
	if annIDStr != "SPDXRef-" {
		fmt.Fprintf(w, "SPDXREF: %s\n", annIDStr)
	}
	if ann.Comment != "" {
		fmt.Fprintf(w, "AnnotationComment: %s\n", textify(ann.Comment))
	}

	return nil
}
