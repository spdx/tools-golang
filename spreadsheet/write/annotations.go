package write

import (
	"fmt"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/xuri/excelize/v2"
)

var AnnotationsHeadersByColumn = map[string]string{
	"A": common.AnnotationsSPDXIdentifier,
	"B": common.AnnotationsComment,
	"C": common.AnnotationsDate,
	"D": common.AnnotationsAnnotator,
	"E": common.AnnotationsType,
}

func WriteAnnotationsRows(doc *spdx.Document2_2, spreadsheet *excelize.File) error {
	rowNum := 2

	// annotations can be document-level, package-level, or file-level

	// document-level
	for _, annotation := range doc.Annotations {
		err := processAnnotation(annotation, spdx.MakeDocElementID("", string(doc.SPDXIdentifier)), spreadsheet, rowNum)
		if err != nil {
			return fmt.Errorf("failed to process document-level annotation: %s", err.Error())
		}

		rowNum += 1
	}

	// package-level
	for _, pkg := range doc.Packages {
		for _, annotation := range pkg.Annotations {
			err := processAnnotation(&annotation, spdx.MakeDocElementID("", string(pkg.PackageSPDXIdentifier)), spreadsheet, rowNum)
			if err != nil {
				return fmt.Errorf("failed to process package-level annotation: %s", err.Error())
			}

			rowNum += 1
		}
	}

	// file-level
	for _, file := range doc.Files {
		for _, annotation := range file.Annotations {
			err := processAnnotation(&annotation, spdx.MakeDocElementID("", string(file.FileSPDXIdentifier)), spreadsheet, rowNum)
			if err != nil {
				return fmt.Errorf("failed to process file-level annotation: %s", err.Error())
			}

			rowNum += 1
		}
	}

	return nil
}

func processAnnotation(annotation *spdx.Annotation2_2, spdxID spdx.DocElementID, spreadsheet *excelize.File, rowNum int) error {
	for column, valueType := range AnnotationsHeadersByColumn {
		axis := common.PositionToAxis(column, rowNum)

		// set `value` to the value to be written to the spreadsheet cell
		var value interface{}
		// if there was a problem determining `value`, set err to something non-nil and processing will be aborted
		var err error

		switch valueType {
		case common.AnnotationsSPDXIdentifier:
			value = spdxID
		case common.AnnotationsComment:
			value = annotation.AnnotationComment
		case common.AnnotationsDate:
			value = annotation.AnnotationDate
		case common.AnnotationsAnnotator:
			err = annotation.Annotator.Validate()
			value = annotation.Annotator.String()
		case common.AnnotationsType:
			value = annotation.AnnotationType
		}

		if err != nil {
			return fmt.Errorf("failed to translate %s for row %d: %s", valueType, rowNum, err.Error())
		}

		err = spreadsheet.SetCellValue(common.SheetNameAnnotations, axis, value)
		if err != nil {
			return fmt.Errorf("failed to set cell %s to %+v: %s", axis, value, err.Error())
		}
	}

	return nil
}
