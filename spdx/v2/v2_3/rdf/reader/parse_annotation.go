// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"errors"
	"fmt"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

// creates a new instance of annotation and sets the annotation attributes
// associated with the given node.
// The newly created annotation is appended to the doc.
func (parser *rdfParser2_3) parseAnnotationFromNode(node *gordfParser.Node) (err error) {
	ann := &spdx.Annotation{}
	for _, subTriple := range parser.nodeToTriples(node) {
		switch subTriple.Predicate.ID {
		case SPDX_ANNOTATOR:
			// cardinality: exactly 1
			err = setAnnotatorFromString(subTriple.Object.ID, ann)
		case SPDX_ANNOTATION_DATE:
			// cardinality: exactly 1
			ann.AnnotationDate = subTriple.Object.ID
		case RDFS_COMMENT:
			// cardinality: exactly 1
			ann.AnnotationComment = subTriple.Object.ID
		case SPDX_ANNOTATION_TYPE:
			// cardinality: exactly 1
			err = setAnnotationType(subTriple.Object.ID, ann)
		case RDF_TYPE:
			// cardinality: exactly 1
			continue
		default:
			err = fmt.Errorf("unknown predicate %s while parsing annotation", subTriple.Predicate.ID)
		}
		if err != nil {
			return err
		}
	}
	return setAnnotationToParser(parser, ann)
}

func setAnnotationToParser(parser *rdfParser2_3, annotation *spdx.Annotation) error {
	if parser.doc == nil {
		return errors.New("uninitialized spdx document")
	}
	if parser.doc.Annotations == nil {
		parser.doc.Annotations = []*spdx.Annotation{}
	}
	parser.doc.Annotations = append(parser.doc.Annotations, annotation)
	return nil
}

// annotator is of type [Person|Organization|Tool]:String
func setAnnotatorFromString(annotatorString string, ann *spdx.Annotation) error {
	subkey, subvalue, err := ExtractSubs(annotatorString, ":")
	if err != nil {
		return err
	}
	if subkey == "Person" || subkey == "Organization" || subkey == "Tool" {
		ann.Annotator.AnnotatorType = subkey
		ann.Annotator.Annotator = subvalue
		return nil
	}
	return fmt.Errorf("unrecognized Annotator type %v while parsing annotation", subkey)
}

// it can be NS_SPDX+annotationType_[review|other]
func setAnnotationType(annType string, ann *spdx.Annotation) error {
	switch annType {
	case SPDX_ANNOTATION_TYPE_OTHER:
		ann.AnnotationType = "OTHER"
	case SPDX_ANNOTATION_TYPE_REVIEW:
		ann.AnnotationType = "REVIEW"
	default:
		return fmt.Errorf("unknown annotation type %s", annType)
	}
	return nil
}
