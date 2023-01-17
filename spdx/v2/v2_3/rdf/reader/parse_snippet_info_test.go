// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"testing"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
)

func Test_rdfParser2_3_getSnippetInformationFromTriple2_3(t *testing.T) {
	var err error
	var parser *rdfParser2_3
	var node *gordfParser.Node

	// TestCase 1: invalid snippet id:
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#Snippet">
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err == nil {
		t.Errorf("expected an error due to invalid, got %v", err)
	}

	// TestCase 2: Invalid LicenseInfoInSnippet
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:licenseInfoInSnippet rdf:resource="http://spdx.org/licenses/Unknown"/>
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err == nil {
		t.Errorf("expected an error due to invalid licenseInfoInSnippet, got %v", err)
	}

	// TestCase 3: Invalid range.
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:range>
				<spdx:StartEndPointer>
					<spdx:unknownTag />
				</spdx:StartEndPointer>
			</spdx:range>
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err == nil {
		t.Errorf("expected an error due to invalid range, got %v", err)
	}

	// TestCase 3: invalid file in snippetFromFile
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:snippetFromFile>
				<spdx:File rdf:resource="http://anupam-VirtualBox/spdx.rdf#item8" />
			</spdx:snippetFromFile>
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err == nil {
		t.Errorf("expected an error due to invalid snippetFromFile, got %v", err)
	}

	// TestCase 4: unknown predicate
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:unknownPredicate />
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err == nil {
		t.Errorf("expected an error due to invalid predicate, got %v", err)
	}

	// TestCase 5: invalid license concluded:
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:licenseConcluded rdf:resource="http://spdx.org/licenses/Unknown"/>
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err == nil {
		t.Errorf("expected an error due to invalid licenseConcluded, got %v", err)
	}

	// TestCase 6: everything valid:
	parser, _ = parserFromBodyContent(`
		<spdx:Snippet rdf:about="#SPDXRef-Snippet">
			<spdx:snippetFromFile>
				<spdx:File rdf:about="#SPDXRef-File" />
			</spdx:snippetFromFile>
			<spdx:range>
				<j.0:StartEndPointer>
					<j.0:startPointer>
						<j.0:LineCharPointer>
							<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
							<j.0:lineNumber>420</j.0:lineNumber>
						</j.0:LineCharPointer>
					</j.0:startPointer>
					<j.0:endPointer>
						<j.0:LineCharPointer>
							<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
							<j.0:lineNumber>310</j.0:lineNumber>
						</j.0:LineCharPointer>
					</j.0:endPointer>
				</j.0:StartEndPointer>
			</spdx:range>
			<spdx:licenseInfoInSnippet rdf:resource="http://spdx.org/rdf/terms#noassertion"/>
			<spdx:name>snippet test</spdx:name>
			<spdx:copyrightText>test</spdx:copyrightText>
			<spdx:licenseComments>comments</spdx:licenseComments>
			<rdfs:comment>comments</rdfs:comment>
			<spdx:licenseConcluded rdf:resource="http://spdx.org/rdf/terms#noassertion"/>
		</spdx:Snippet>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, err = parser.getSnippetInformationFromNode2_3(node)
	if err != nil {
		t.Fatalf("error parsing a valid example: %v", err)
	}
}

func Test_setSnippetID(t *testing.T) {
	// TestCase 1: invalid input (empty)
	err := setSnippetID("", &spdx.Snippet{})
	if err == nil {
		t.Errorf("should've raised an error for empty input")
	}

	// TestCase 2: valid input
	si := &spdx.Snippet{}
	err = setSnippetID("http://spdx.org/spdxdocs/spdx-example#SPDXRef-Snippet", si)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if si.SnippetSPDXIdentifier != "Snippet" {
		t.Errorf("expected: %s, found: %s", "Snippet", si.SnippetSPDXIdentifier)
	}
}

func Test_rdfParser2_3_parseRangeReference(t *testing.T) {
	var err error
	var node *gordfParser.Node
	var parser *rdfParser2_3
	var si *spdx.Snippet

	// TestCase 1: ResourceLiteral node without a new file shouldn't raise any error.
	si = &spdx.Snippet{}
	parser, _ = parserFromBodyContent(``)
	node = &gordfParser.Node{
		NodeType: gordfParser.RESOURCELITERAL,
		ID:       "http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource",
	}
	err = parser.parseRangeReference(node, si)
	if err != nil {
		t.Errorf("error parsing a valid node: %v", err)
	}

	// TestCase 2: invalid file in the reference should raise an error
	si = &spdx.Snippet{}
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#DoapSource">
			<spdx:fileName> test file </spdx:fileName>
		</spdx:File>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseRangeReference(node, si)
	if err == nil {
		t.Errorf("expected an error due to invalid file in the range reference, got %v", err)
	}

	// TestCase 3: A valid reference must set the file to the files map of the parser.
	si = &spdx.Snippet{}
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource">
			<spdx:fileName> test file </spdx:fileName>
		</spdx:File>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseRangeReference(node, si)
	if err != nil {
		t.Errorf("error parsing a valid input: %v", err)
	}
	if len(parser.files) != 1 {
		t.Errorf("expected parser.files to have 1 file, found %d", len(parser.files))
	}
}

func Test_rdfParser2_3_getPointerFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var node *gordfParser.Node
	var si *spdx.Snippet
	var err error
	var rt RangeType
	var number int

	// TestCase 1: invalid number in the offset field must raise an error.
	parser, _ = parserFromBodyContent(`
		<j.0:startPointer>
			<j.0:LineCharPointer>
				<j.0:reference rdf:resource="#SPDXRef-DoapSource"/>
				<j.0:offset>3-10</j.0:offset>
			</j.0:LineCharPointer>
		</j.0:startPointer>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, _, err = parser.getPointerFromNode(node, si)
	if err == nil {
		t.Errorf("should've raised an error parsing invalid offset, got %v", err)
	}

	// TestCase 2: invalid number in the lineNumber field must raise an error.
	parser, _ = parserFromBodyContent(`
		<j.0:ByteOffsetPointer>
			<j.0:reference rdf:resource="#SPDXRef-DoapSource"/>
			<j.0:offset>3-10</j.0:offset>
		</j.0:ByteOffsetPointer>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, _, err = parser.getPointerFromNode(node, si)
	if err == nil {
		t.Errorf("should've raised an error parsing invalid offset, got %v", err)
	}

	// TestCase 3: invalid predicate in the pointer field
	parser, _ = parserFromBodyContent(`
		<j.0:ByteOffsetPointer>
			<spdx:invalidTag />
			<j.0:reference rdf:resource="#SPDXRef-DoapSource"/>
			<j.0:lineNumber>3-10</j.0:lineNumber>
		</j.0:ByteOffsetPointer>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, _, err = parser.getPointerFromNode(node, si)
	if err == nil {
		t.Errorf("should've raised an error parsing invalid predicate, got %v", err)
	}

	// TestCase 4: No range type defined must also raise an error
	parser, _ = parserFromBodyContent(`
		<j.0:ByteOffsetPointer>
			<j.0:reference rdf:resource="#SPDXRef-DoapSource"/>
		</j.0:ByteOffsetPointer>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	_, _, err = parser.getPointerFromNode(node, si)
	if err == nil {
		t.Errorf("should've raised an error parsing invalid rangeType, got %v", err)
	}

	// TestCase 5: valid example
	parser, _ = parserFromBodyContent(`
		<j.0:ByteOffsetPointer>
			<j.0:reference rdf:resource="#SPDXRef-DoapSource"/>
			<j.0:offset>310</j.0:offset>
		</j.0:ByteOffsetPointer>
	`)
	node = parser.gordfParserObj.Triples[0].Subject
	rt, number, err = parser.getPointerFromNode(node, si)
	if err != nil {
		t.Fatalf("unexpected error parsing a valid node: %v", err)
	}
	if rt != BYTE_RANGE {
		t.Errorf("expected: %s, got: %s", BYTE_RANGE, rt)
	}
	if number != 310 {
		t.Errorf("expected: %d, got: %d", 310, number)
	}
}

func Test_rdfParser2_3_setSnippetRangeFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var err error
	var si *spdx.Snippet
	var node *gordfParser.Node

	// TestCase 1: range with less one pointer less must raise an error
	//			   (end-pointer missing in the range)
	parser, _ = parserFromBodyContent(`
            <j.0:StartEndPointer>
                <j.0:startPointer>
                    <j.0:LineCharPointer>
                        <j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
                        <j.0:offset>310</j.0:offset>
                    </j.0:LineCharPointer>
                </j.0:startPointer>
            </j.0:StartEndPointer>
        
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to missing end pointer, got %v", err)
	}

	// TestCase 2: triples with 0 or more than one type-triple
	parser, _ = parserFromBodyContent(`
		
            <j.0:StartEndPointer>
                <j.0:endPointer>
                    <j.0:ByteOffsetPointer>
                        <j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
                        <j.0:offset>420</j.0:offset>
                    </j.0:ByteOffsetPointer>
                </j.0:endPointer>
                <j.0:startPointer>
                    <j.0:ByteOffsetPointer>
                        <j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
                        <j.0:offset>310</j.0:offset>
                    </j.0:ByteOffsetPointer>
                </j.0:startPointer>
            </j.0:StartEndPointer>
        
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	dummyTriple := parser.gordfParserObj.Triples[0]
	// resetting the node to be associated with 3 triples which will have
	// rdf:type triple either thrice or 0 times.
	parser.nodeStringToTriples[node.String()] = []*gordfParser.Triple{
		dummyTriple, dummyTriple, dummyTriple,
	}
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to invalid rdf:type triples, got %v", err)
	}

	// TestCase 3: triples with 0 startPointer
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:endPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>420</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:endPointer>
			<j.0:endPointer>
				<j.0:LineCharPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>310</j.0:offset>
				</j.0:LineCharPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer>
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to missing start pointer, got %v", err)
	}

	// TestCase 4: triples with 0 endPointer
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:endPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>420</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:endPointer>
			<j.0:endPointer>
				<j.0:LineCharPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>310</j.0:offset>
				</j.0:LineCharPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer>
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to missing end pointer, got %v", err)
	}

	// TestCase 5: error parsing start pointer must be propagated to the range
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:startPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>42.0</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:startPointer>
			<j.0:endPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>310</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer>
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to invalid start pointer, got %v", err)
	}

	// TestCase 6: error parsing end pointer must be propagated to the range
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:startPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>420</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:startPointer>
			<j.0:endPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>31+0</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer> 
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to invalid end pointer, got %v", err)
	}

	// TestCase 7: mismatching start and end pointer must also raise an error.
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:startPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>420</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:startPointer>
			<j.0:endPointer>
				<j.0:LineCharPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:lineNumber>310</j.0:lineNumber>
				</j.0:LineCharPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer>
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err == nil {
		t.Errorf("expected an error due to mismatching start and end pointers, got %v", err)
	}

	// TestCase 8: everything valid(byte_range):
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:startPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>420</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:startPointer>
			<j.0:endPointer>
				<j.0:ByteOffsetPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:offset>310</j.0:offset>
				</j.0:ByteOffsetPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer>
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// TestCase 9: everything valid(line_range):
	parser, _ = parserFromBodyContent(`
		<j.0:StartEndPointer>
			<j.0:startPointer>
				<j.0:LineCharPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:lineNumber>420</j.0:lineNumber>
				</j.0:LineCharPointer>
			</j.0:startPointer>
			<j.0:endPointer>
				<j.0:LineCharPointer>
					<j.0:reference rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-DoapSource"/>
					<j.0:lineNumber>310</j.0:lineNumber>
				</j.0:LineCharPointer>
			</j.0:endPointer>
		</j.0:StartEndPointer>
	`)
	si = &spdx.Snippet{}
	node = parser.gordfParserObj.Triples[0].Subject
	err = parser.setSnippetRangeFromNode(node, si)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_rdfParser2_3_setSnippetToFileWithID(t *testing.T) {
	var parser *rdfParser2_3
	var fileId common.ElementID
	var si *spdx.Snippet
	var file *spdx.File
	var err error

	// TestCase 1: file id which is not associated with any file must raise an error.
	parser, _ = parserFromBodyContent("")
	si = &spdx.Snippet{}
	err = parser.setSnippetToFileWithID(si, fileId)
	if err == nil {
		t.Errorf("expected an error saying undefined file")
	}

	// TestCase 2: file exists, but snippet of the file doesn't ( it mustn't raise any error )
	fileId = common.ElementID("File1")
	file = &spdx.File{
		FileSPDXIdentifier: fileId,
	}
	parser.files[fileId] = file
	file.Snippets = nil // nil snippets
	err = parser.setSnippetToFileWithID(si, fileId)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(file.Snippets) != 1 {
		t.Errorf("expected file to have 1 snippet, got %d", len(file.Snippets))
	}
}
