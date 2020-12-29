// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v1

import (
	"fmt"
)

func (parser *tvParser2_1) parsePairForAnnotation2_1(tag string, value string) error {
	if parser.ann == nil {
		return fmt.Errorf("no annotation struct created in parser ann pointer")
	}

	switch tag {
	case "Annotator":
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		if subkey == "Person" || subkey == "Organization" || subkey == "Tool" {
			parser.ann.AnnotatorType = subkey
			parser.ann.Annotator = subvalue
			return nil
		}
		return fmt.Errorf("unrecognized Annotator type %v", subkey)
	case "AnnotationDate":
		parser.ann.Date = value
	case "AnnotationType":
		parser.ann.Type = value
	case "SPDXREF":
		deID, err := extractDocElementID(value)
		if err != nil {
			return err
		}
		parser.ann.SPDXIdentifier = deID
	case "AnnotationComment":
		parser.ann.Comment = value
	default:
		return fmt.Errorf("received unknown tag %v in Annotation section", tag)
	}

	return nil
}
