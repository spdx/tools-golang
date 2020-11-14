// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"github.com/spdx/tools-golang/spdx"
	"testing"
)

func Test_setCreator(t *testing.T) {
	// TestCase 1: invalid creator (empty)
	input := ""
	err := setCreator(input, &spdx.CreationInfo2_2{})
	if err == nil {
		t.Errorf("shoud've raised an error due to invalid input")
	}

	// TestCase 2: invalid entity type
	input = "Company: some company"
	err = setCreator(input, &spdx.CreationInfo2_2{})
	if err == nil {
		t.Errorf("shoud've raised an error due to unknown entity type")
	}

	// TestCase 3: valid input
	input = "Person: Jane Doe"
	ci := &spdx.CreationInfo2_2{}
	err = setCreator(input, ci)
	if err != nil {
		t.Errorf("error parsing a valid input: %v", err)
	}
	if len(ci.CreatorPersons) != 1 {
		t.Errorf("creationInfo should've had 1 creatorPersons, found %d", len(ci.CreatorPersons))
	}
	expectedPerson := "Jane Doe"
	if ci.CreatorPersons[0] != expectedPerson {
		t.Errorf("expected %s, found %s", expectedPerson, ci.CreatorPersons[0])
	}
}

func Test_rdfParser2_2_parseCreationInfoFromNode(t *testing.T) {
	// TestCase 1: invalid creator must raise an error
	parser, _ := parserFromBodyContent(`
		<spdx:CreationInfo>
			<spdx:licenseListVersion>2.6</spdx:licenseListVersion>
			<spdx:creator>Person Unknown</spdx:creator>
			<spdx:created>2018-08-24T19:55:34Z</spdx:created>
		</spdx:CreationInfo>
	`)
	ciNode := parser.gordfParserObj.Triples[0].Subject
	err := parser.parseCreationInfoFromNode(&spdx.CreationInfo2_2{}, ciNode)
	if err == nil {
		t.Errorf("invalid creator must raise an error")
	}

	// TestCase 2: unknown predicate must also raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:CreationInfo>
			<spdx:licenseListVersion>2.6</spdx:licenseListVersion>
			<spdx:creator>Person: fossy (y)</spdx:creator>
			<spdx:creator>Organization: </spdx:creator>
			<spdx:creator>Tool: spdx2</spdx:creator>
			<spdx:created>2018-08-24T19:55:34Z</spdx:created>
			<spdx:unknownPredicate />
		</spdx:CreationInfo>
	`)
	ciNode = parser.gordfParserObj.Triples[0].Subject
	err = parser.parseCreationInfoFromNode(&spdx.CreationInfo2_2{}, ciNode)
	if err == nil {
		t.Errorf("unknown predicate must raise an error")
	}

	// TestCase 2: unknown predicate must also raise an error
	parser, _ = parserFromBodyContent(`
		<spdx:CreationInfo>
			<spdx:licenseListVersion>2.6</spdx:licenseListVersion>
			<spdx:creator>Person: fossy</spdx:creator>
			<spdx:created>2018-08-24T19:55:34Z</spdx:created>
			<rdfs:comment>comment</rdfs:comment>
		</spdx:CreationInfo>
	`)
	ciNode = parser.gordfParserObj.Triples[0].Subject
	ci := &spdx.CreationInfo2_2{}
	err = parser.parseCreationInfoFromNode(ci, ciNode)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ci.LicenseListVersion != "2.6" {
		t.Errorf(`expected %s, found %s`, "2.6", ci.LicenseListVersion)
	}
	n := len(ci.CreatorPersons)
	if n != 1 {
		t.Errorf("expected 1 creatorPersons, found %d", n)
	}
	if ci.CreatorPersons[0] != "fossy" {
		t.Errorf("expected %s, found %s", "fossy", ci.CreatorPersons[0])
	}
	expectedCreated := "2018-08-24T19:55:34Z"
	if ci.Created != expectedCreated {
		t.Errorf("expected %s, found %s", expectedCreated, ci.Created)
	}
	expectedComment := "comment"
	if ci.CreatorComment != expectedComment {
		t.Errorf("expected %s, found %s", expectedComment, ci.CreatorComment)
	}
}
