package parse

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
)

func ProcessAnnotationsRows(rows [][]string, doc *spdx.Document2_2) error {
	// the first row is column headers, keep track of which order they appear in
	columnsByIndex := make(map[int]string)
	for index, header := range rows[0] {
		columnsByIndex[index] = header
	}

	for rowNum, row := range rows[1:] {
		// set rowNum to the correct value, Go slices are zero-indexed (+1), and we started iterating on the second element (+1)
		rowNum = rowNum + 2

		newAnnotation := spdx.Annotation2_2{}

		for columnIndex, value := range row {
			if value == "" {
				continue
			}

			switch columnsByIndex[columnIndex] {
			case common.AnnotationsSPDXIdentifier:
				id := spdx.DocElementID{}
				err := id.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.AnnotationsSPDXIdentifier, rowNum, err.Error())
				}

				newAnnotation.AnnotationSPDXIdentifier = id
			case common.AnnotationsComment:
				newAnnotation.AnnotationComment = value
			case common.AnnotationsDate:
				newAnnotation.AnnotationDate = value
			case common.AnnotationsAnnotator:
				annotator := spdx.Annotator{}
				err := annotator.FromString(value)
				if err != nil {
					return fmt.Errorf("invalid %s in row %d: %v", common.AnnotationsAnnotator, rowNum, err.Error())
				}

				newAnnotation.Annotator = annotator
			case common.AnnotationsType:
				newAnnotation.AnnotationType = value
			}
		}

		// TODO: validate?

		// an annotation can be at the Document level, File level, or Package level
		if newAnnotation.AnnotationSPDXIdentifier.DocumentRefID == "" && newAnnotation.AnnotationSPDXIdentifier.ElementRefID != doc.SPDXIdentifier {
			var found bool
			for ii, pkg := range doc.Packages {
				if newAnnotation.AnnotationSPDXIdentifier.ElementRefID == pkg.PackageSPDXIdentifier {
					// package level
					found = true
					doc.Packages[ii].Annotations = append(doc.Packages[ii].Annotations, newAnnotation)
					break
				}
			}

			if !found {
				for ii, file := range doc.Files {
					if newAnnotation.AnnotationSPDXIdentifier.ElementRefID == file.FileSPDXIdentifier {
						// file level
						found = true
						doc.Files[ii].Annotations = append(doc.Files[ii].Annotations, newAnnotation)
						break
					}
				}
			}

			if !found {
				return fmt.Errorf("annotation SPDX Identifier from row %d not found in document: %s", rowNum, newAnnotation.AnnotationSPDXIdentifier)
			}
		} else {
			// document level
			doc.Annotations = append(doc.Annotations, &newAnnotation)
		}
	}

	return nil
}
