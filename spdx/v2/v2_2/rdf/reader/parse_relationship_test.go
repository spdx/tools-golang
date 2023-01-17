// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"reflect"
	"testing"

	"github.com/spdx/gordf/rdfwriter"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func Test_getReferenceFromURI(t *testing.T) {
	// TestCase 1: noassertion uri
	ref, err := getReferenceFromURI(SPDX_NOASSERTION_CAPS)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ref.DocumentRefID != "" {
		t.Errorf("reference's documentRefID should've been empty, found %s", ref.DocumentRefID)
	}
	if ref.ElementRefID != "NOASSERTION" {
		t.Errorf("mismatching elementRefID. Found %s, expected %s", ref.ElementRefID, "NOASSERTION")
	}

	// TestCase 2: NONE uri
	ref, err = getReferenceFromURI(SPDX_NONE_CAPS)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ref.DocumentRefID != "" {
		t.Errorf("reference's documentRefID should've been empty, found %s", ref.DocumentRefID)
	}
	if ref.ElementRefID != "NONE" {
		t.Errorf("mismatching elementRefID. Found %s, expected %s", ref.ElementRefID, "NONE")
	}

	// TestCase 3: Valid URI
	ref, err = getReferenceFromURI(NS_SPDX + "SPDXRef-item1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ref.DocumentRefID != "" {
		t.Errorf("reference's documentRefID should've been empty, found %s", ref.DocumentRefID)
	}
	if ref.ElementRefID != "item1" {
		t.Errorf("mismatching elementRefID. Found %s, expected %s", ref.ElementRefID, "item1")
	}

	// TestCase 3: Invalid URI
	_, err = getReferenceFromURI(NS_SPDX + "item1")
	if err == nil {
		t.Errorf("should've raised an error for invalid input")
	}
}

func Test_getRelationshipTypeFromURI(t *testing.T) {
	// TestCase 1: valid relationshipType
	relnType := "expandedFromArchive"
	op, err := getRelationshipTypeFromURI(NS_SPDX + "relationshipType_" + relnType)
	if err != nil {
		t.Errorf("error getting relationship type from a valid input")
	}
	if op != relnType {
		t.Errorf("expected %s, found %s", relnType, op)
	}

	// TestCase2: invalid relationshipType
	relnType = "invalidRelationship"
	_, err = getRelationshipTypeFromURI(NS_SPDX + "relationshipType_" + relnType)
	if err == nil {
		t.Errorf("should've raised an error for an invalid input(%s)", relnType)
	}
}

func Test_rdfParser2_2_parseRelatedElementFromTriple(t *testing.T) {
	// TestCase 1: Package as a related element
	parser, _ := parserFromBodyContent(`
		<spdx:Relationship>
			<spdx:relatedSpdxElement>
				<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
			</spdx:relatedSpdxElement>
		</spdx:Relationship>
	`)
	reln := &v2_2.Relationship{}
	triple := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_PACKAGE)[0]
	err := parser.parseRelatedElementFromTriple(reln, triple)
	if err != nil {
		t.Errorf("error parsing a valid example")
	}
	expectedRefA := common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "",
	}
	if !reflect.DeepEqual(expectedRefA, reln.RefA) {
		t.Errorf("expected %+v, found %+v", expectedRefA, reln.RefA)
	}
	expectedRefB := common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "Saxon",
	}
	if !reflect.DeepEqual(expectedRefB, reln.RefB) {
		t.Errorf("expected %+v, found %+v", expectedRefB, reln.RefB)
	}

	// TestCase 3: invalid package as a relatedElement
	parser, _ = parserFromBodyContent(`
		<spdx:Relationship>
			<spdx:relatedSpdxElement>
				<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#Saxon"/>
			</spdx:relatedSpdxElement>
		</spdx:Relationship>
	`)
	reln = &v2_2.Relationship{}
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_PACKAGE)[0]
	err = parser.parseRelatedElementFromTriple(reln, triple)
	if err == nil {
		t.Errorf("expected an error due to invalid Package id, got %v", err)
	}

	// TestCase 4: valid File as a related element
	parser, _ = parserFromBodyContent(`
		<spdx:Relationship>
			<spdx:relatedSpdxElement>
				<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
			</spdx:relatedSpdxElement>
		</spdx:Relationship>
	`)
	reln = &v2_2.Relationship{}
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0]
	err = parser.parseRelatedElementFromTriple(reln, triple)
	if err != nil {
		t.Errorf("error parsing a valid example")
	}
	expectedRefA = common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "",
	}
	if !reflect.DeepEqual(expectedRefA, reln.RefA) {
		t.Errorf("expected %+v, found %+v", expectedRefA, reln.RefA)
	}
	expectedRefB = common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "Saxon",
	}
	if !reflect.DeepEqual(expectedRefB, reln.RefB) {
		t.Errorf("expected %+v, found %+v", expectedRefB, reln.RefB)
	}

	// TestCase 5: invalid File as a relatedElement
	parser, _ = parserFromBodyContent(`
		<spdx:Relationship>
			<spdx:relatedSpdxElement>
				<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#Saxon"/>
			</spdx:relatedSpdxElement>
		</spdx:Relationship>
	`)
	reln = &v2_2.Relationship{}
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_FILE)[0]
	err = parser.parseRelatedElementFromTriple(reln, triple)
	if err == nil {
		t.Errorf("expected an error while parsing an invalid File, got %v", err)
	}

	// TestCase 6: valid SpdxElement as a related element
	parser, _ = parserFromBodyContent(`
		<spdx:Relationship>		
			<spdx:relatedSpdxElement>
				<spdx:SpdxElement rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File"/>
			</spdx:relatedSpdxElement>
		</spdx:Relationship>
	`)
	reln = &v2_2.Relationship{}
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_SPDX_ELEMENT)[0]
	err = parser.parseRelatedElementFromTriple(reln, triple)
	if err != nil {
		t.Errorf("error parsing a valid example")
	}
	expectedRefA = common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "",
	}
	if !reflect.DeepEqual(expectedRefA, reln.RefA) {
		t.Errorf("expected %+v, found %+v", expectedRefA, reln.RefA)
	}
	expectedRefB = common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "File",
	}
	if !reflect.DeepEqual(expectedRefB, reln.RefB) {
		t.Errorf("expected %+v, found %+v", expectedRefB, reln.RefB)
	}

	// TestCase 7: invalid SpdxElement as a related element
	parser, _ = parserFromBodyContent(`
		<spdx:Relationship>		
			<spdx:relatedSpdxElement>
				<spdx:SpdxElement rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-:File"/>
			</spdx:relatedSpdxElement>
		</spdx:Relationship>
	`)
	reln = &v2_2.Relationship{}
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &RDF_TYPE, &SPDX_SPDX_ELEMENT)[0]
	err = parser.parseRelatedElementFromTriple(reln, triple)
	if err == nil {
		t.Errorf("expected an error due to invalid documentId for SpdxElement, got %v", err)
	}
}

func Test_rdfParser2_2_parseRelationship(t *testing.T) {
	// TestCase 1: invalid RefA
	parser, _ := parserFromBodyContent(`
		<spdx:File>
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err := parser.parseRelationship(triple)
	if err == nil {
		t.Errorf("should've raised an error due to invalid RefA")
	}

	// TestCase 3: invalid RefB
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#Saxon"/>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err == nil {
		t.Errorf("should've raised an error due to invalid RefB")
	}

	// TestCase 3: more than one typeTriple for relatedElement
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
					</spdx:relatedSpdxElement>
					<spdx:relatedSpdxElement>
						<spdx:File/>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err == nil {
		t.Errorf("should've raised an error due to more than one type triples")
	}

	// TestCase 4: undefined relatedSpdxElement
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relatedSpdxElement>
						<spdx:Unknown rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err == nil {
		t.Errorf("should've raised an error due to unknown relatedElement, got %v", err)
	}

	// TestCase 6: relatedElement associated with more than one type
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err == nil {
		t.Errorf("expected an error due to invalid relatedElement, got %v", err)
	}

	// TestCase 5: unknown predicate inside a relationship
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relatedSpdxElement>
						<spdx:Unknown rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
					</spdx:relatedSpdxElement>
					<spdx:unknownPredicate/>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err == nil {
		t.Errorf("should've raised an error due to unknown predicate in a relationship")
	}

	// TestCase 8: Recursive relationships mustn't raise any error:
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship rdf:about="#SPDXRef-reln">
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_describes"/>
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon">
							<spdx:relationship>
								<spdx:Relationship rdf:about="#SPDXRef-reln">
									<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_describes"/>
									<spdx:relatedSpdxElement rdf:resource="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File"/>
								</spdx:Relationship>
							</spdx:relationship>
						</spdx:Package>
					</spdx:relatedSpdxElement>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err != nil {
		t.Errorf("error parsing a valid example")
	}

	// TestCase 7: completely valid example:
	parser, _ = parserFromBodyContent(`
		<spdx:File rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-File">
			<spdx:relationship>
				<spdx:Relationship>
					<spdx:relationshipType rdf:resource="http://spdx.org/rdf/terms#relationshipType_describes"/>
					<spdx:relatedSpdxElement>
						<spdx:Package rdf:about="http://spdx.org/spdxdocs/spdx-example-444504E0-4F89-41D3-9A0C-0305E82C3301#SPDXRef-Saxon"/>
					</spdx:relatedSpdxElement>
					<rdfs:comment>comment</rdfs:comment>
				</spdx:Relationship>
			</spdx:relationship>
		</spdx:File>
	`)
	triple = rdfwriter.FilterTriples(parser.gordfParserObj.Triples, nil, &SPDX_RELATIONSHIP, nil)[0]
	err = parser.parseRelationship(triple)
	if err != nil {
		t.Errorf("unexpected error parsing a valid relationship: %v", err)
	}
	// validating parsed attributes
	if len(parser.doc.Relationships) != 1 {
		t.Errorf("after parsing a valid relationship, doc should've had 1 relationship, found %d", len(parser.doc.Relationships))
	}
	reln := parser.doc.Relationships[0]
	expectedRelnType := "describes"
	if reln.Relationship != expectedRelnType {
		t.Errorf("expected %s, found %s", expectedRelnType, reln.Relationship)
	}
	expectedRefA := common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "File",
	}
	if !reflect.DeepEqual(expectedRefA, reln.RefA) {
		t.Errorf("expected %+v, found %+v", expectedRefA, reln.RefA)
	}
	expectedRefB := common.DocElementID{
		DocumentRefID: "",
		ElementRefID:  "Saxon",
	}
	if !reflect.DeepEqual(expectedRefB, reln.RefB) {
		t.Errorf("expected %+v, found %+v", expectedRefB, reln.RefB)
	}
	expectedComment := "comment"
	if reln.RelationshipComment != expectedComment {
		t.Errorf("expected %v, found %v", expectedComment, reln.RelationshipComment)
	}
}
