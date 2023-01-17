// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"testing"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func Test_setCreator(t *testing.T) {
	// TestCase 1: invalid creator (empty)
	input := ""
	err := setCreator(input, &spdx.CreationInfo{})
	if err == nil {
		t.Errorf("shoud've raised an error due to invalid input")
	}

	// TestCase 2: invalid entity type
	input = "Company: some company"
	err = setCreator(input, &spdx.CreationInfo{})
	if err == nil {
		t.Errorf("shoud've raised an error due to unknown entity type")
	}

	// TestCase 3: valid input
	input = "Person: Jane Doe"
	ci := &spdx.CreationInfo{}
	err = setCreator(input, ci)
	if err != nil {
		t.Errorf("error parsing a valid input: %v", err)
	}
	if len(ci.Creators) != 1 {
		t.Errorf("creationInfo should've had 1 creatorPersons, found %d", len(ci.Creators))
	}
	expectedPerson := "Jane Doe"
	if ci.Creators[0].Creator != expectedPerson {
		t.Errorf("expected %s, found %s", expectedPerson, ci.Creators[0])
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
	err := parser.parseCreationInfoFromNode(&spdx.CreationInfo{}, ciNode)
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
	err = parser.parseCreationInfoFromNode(&spdx.CreationInfo{}, ciNode)
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
	ci := &spdx.CreationInfo{}
	err = parser.parseCreationInfoFromNode(ci, ciNode)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ci.LicenseListVersion != "2.6" {
		t.Errorf(`expected %s, found %s`, "2.6", ci.LicenseListVersion)
	}
	n := len(ci.Creators)
	if n != 1 {
		t.Errorf("expected 1 creatorPersons, found %d", n)
	}
	if ci.Creators[0].Creator != "fossy" {
		t.Errorf("expected %s, found %s", "fossy", ci.Creators[0].Creator)
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
