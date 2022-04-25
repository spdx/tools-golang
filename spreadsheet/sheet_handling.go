package spdx_xls

import (
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spreadsheet/common"
	"github.com/spdx/tools-golang/spreadsheet/parse"
	"github.com/spdx/tools-golang/spreadsheet/write"
	"github.com/xuri/excelize/v2"
)

// sheetParserFunc is a function that takes in the data from a sheet as a slice of rows and iterates through them to
// fill in information in the given spdx.Document2_2.
// Returns an error if any occurred.
type sheetParserFunc func(rows [][]string, doc *spdx.Document2_2) error

// sheetWriterFunc is a function that takes in a spdx.Document2_2 and a spreadsheet as a *excelize.File and iterates
// through particular section of the Document spdx.Document2_2 in order to write out data to the spreadsheet.
// Returns an error if any occurred.
type sheetWriterFunc func(doc *spdx.Document2_2, spreadsheet *excelize.File) error

// sheetHandlingInformation defines info that is needed for parsing individual sheets in a workbook.
type sheetHandlingInformation struct {
	// SheetName is the name of the sheet
	SheetName string

	// HeadersByColumn is a map of header names to which column the header should go in.
	// This is used only when writing/exporting a spreadsheet.
	// During parsing/imports, the header positions are parsed dynamically.
	HeadersByColumn map[string]string

	// ParserFunc is the function that should be used to parse a particular sheet
	ParserFunc sheetParserFunc

	// WriterFunc is the function that should be used to write a particular sheet
	WriterFunc sheetWriterFunc

	// SheetIsRequired denotes whether the sheet is required to be present in the workbook, or if it is optional (false)
	SheetIsRequired bool
}

// sheetHandlers contains handling information for each sheet in the workbook.
// The order of this slice determines the order in which the sheets are processed.
var sheetHandlers = []sheetHandlingInformation{
	{
		SheetName:       common.SheetNameDocumentInfo,
		HeadersByColumn: write.DocumentInfoHeadersByColumn,
		ParserFunc:      parse.ProcessDocumentInfoRows,
		WriterFunc:      write.WriteDocumentInfoRows,
		SheetIsRequired: true,
	},
	{
		SheetName:       common.SheetNamePackageInfo,
		HeadersByColumn: write.PackageInfoHeadersByColumn,
		ParserFunc:      parse.ProcessPackageInfoRows,
		WriterFunc:      write.WritePackageInfoRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameExternalRefs,
		HeadersByColumn: write.ExternalRefsHeadersByColumn,
		ParserFunc:      parse.ProcessPackageExternalRefsRows,
		WriterFunc:      write.WriteExternalRefsRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameExtractedLicenseInfo,
		HeadersByColumn: write.ExtractedLicenseInfoHeadersByColumn,
		ParserFunc:      parse.ProcessExtractedLicenseInfoRows,
		WriterFunc:      write.WriteExtractedLicenseInfoRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameFileInfo,
		HeadersByColumn: write.FileInfoHeadersByColumn,
		ParserFunc:      parse.ProcessPerFileInfoRows,
		WriterFunc:      write.WriteFileInfoRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameRelationships,
		HeadersByColumn: write.RelationshipsHeadersByColumn,
		ParserFunc:      parse.ProcessRelationshipsRows,
		WriterFunc:      write.WriteRelationshipsRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameAnnotations,
		HeadersByColumn: write.AnnotationsHeadersByColumn,
		ParserFunc:      parse.ProcessAnnotationsRows,
		WriterFunc:      write.WriteAnnotationsRows,
		SheetIsRequired: false,
	},
	{
		SheetName:       common.SheetNameSnippets,
		HeadersByColumn: write.SnippetsHeadersByColumn,
		ParserFunc:      parse.ProcessSnippetsRows,
		WriterFunc:      write.WriteSnippetsRows,
		SheetIsRequired: false,
	},
}
