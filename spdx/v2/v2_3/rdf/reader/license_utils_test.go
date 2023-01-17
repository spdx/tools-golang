// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"reflect"
	"testing"
)

func Test_getLicenseStringFromURI(t *testing.T) {
	// TestCase 1: NONE license
	input := SPDX_NONE_CAPS
	output := getLicenseStringFromURI(input)
	expectedOutput := "NONE"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", expectedOutput, output)
	}

	// TestCase 2: NOASSERTION license
	input = SPDX_NOASSERTION_SMALL
	output = getLicenseStringFromURI(input)
	expectedOutput = "NOASSERTION"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", expectedOutput, output)
	}

	// TestCase 3: Other license
	input = NS_SPDX + "LicenseRef-1"
	output = getLicenseStringFromURI(input)
	expectedOutput = "LicenseRef-1"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", expectedOutput, output)
	}
}

func Test_rdfParser2_3_getChecksumFromNode(t *testing.T) {
	var parser *rdfParser2_3
	var err error
	// TestCase 1: invalid checksum algorithm
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha999"/>
		</spdx:Checksum>
	`)
	checksumNode := parser.gordfParserObj.Triples[0].Subject
	_, _, err = parser.getChecksumFromNode(checksumNode)
	if err == nil {
		t.Errorf("expected an error saying invalid checksum algorithm")
	}

	// TestCase 2: invalid predicate
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1"/>
			<spdx:invalidPredicate />
		</spdx:Checksum>
	`)
	checksumNode = parser.gordfParserObj.Triples[0].Subject
	_, _, err = parser.getChecksumFromNode(checksumNode)
	if err == nil {
		t.Errorf("expected an error saying invalid predicate")
	}

	// TestCase 3: valid input
	parser, _ = parserFromBodyContent(`
		<spdx:Checksum>
			<spdx:checksumValue>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</spdx:checksumValue>
			<spdx:algorithm rdf:resource="http://spdx.org/rdf/terms#checksumAlgorithm_sha1"/>
		</spdx:Checksum>
	`)
	checksumNode = parser.gordfParserObj.Triples[0].Subject
	algorithm, value, err := parser.getChecksumFromNode(checksumNode)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if algorithm != "SHA1" {
		t.Errorf("expected checksum algorithm to be sha1, found %s", algorithm)
	}
	expectedValue := "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"
	if value != expectedValue {
		t.Errorf("expected checksumValue to be %s, found %s", expectedValue, value)
	}
}

func Test_rdfParser2_3_getAlgorithmFromURI(t *testing.T) {
	var algorithmURI string
	var err error

	// TestCase 1: checksumAlgorithm uri doesn't start with checksumAlgorithm_
	algorithmURI = NS_SPDX + "sha1"
	_, err = getAlgorithmFromURI(algorithmURI)
	if err == nil {
		t.Errorf("should've raised an error for algorithmURI that doesn't start with checksumAlgorithm_")
	}

	// TestCase 2: unknown checksum algorithm
	algorithmURI = NS_SPDX + "checksumAlgorithm_sha999"
	_, err = getAlgorithmFromURI(algorithmURI)
	if err == nil {
		t.Errorf("should've raised an error for invalid algorithm")
	}

	// TestCase 3: valid input
	algorithmURI = NS_SPDX + "checksumAlgorithm_sha256"
	algorithm, err := getAlgorithmFromURI(algorithmURI)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if algorithm != "SHA256" {
		t.Errorf("expected: SHA256, found: %s", algorithm)
	}
}

func Test_mapLicensesToStrings(t *testing.T) {
	// nothing much to test here.
	// just a dummy dry run.
	licenses := []AnyLicenseInfo{
		SpecialLicense{
			value: NONE,
		},
		SpecialLicense{
			value: NOASSERTION,
		},
	}
	licenseStrings := mapLicensesToStrings(licenses)
	expectedLicenseStrings := []string{"NONE", "NOASSERTION"}
	if !reflect.DeepEqual(licenseStrings, expectedLicenseStrings) {
		t.Errorf("expected: %+v\nfound %+v", expectedLicenseStrings, licenseStrings)
	}
}

func TestConjunctiveLicenseSet_ToLicenseString(t *testing.T) {
	var lic ConjunctiveLicenseSet
	var output, expectedOutput string

	// TestCase 1: no license in the set
	lic = ConjunctiveLicenseSet{
		members: nil,
	}
	output = lic.ToLicenseString()
	expectedOutput = ""
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}

	// TestCase 2: single license in the set
	lic = ConjunctiveLicenseSet{
		members: []AnyLicenseInfo{
			SpecialLicense{value: NOASSERTION},
		},
	}
	output = lic.ToLicenseString()
	expectedOutput = "NOASSERTION"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}

	// TestCase 3: more than one license in the set.
	lic = ConjunctiveLicenseSet{
		members: []AnyLicenseInfo{
			SpecialLicense{value: NOASSERTION},
			SpecialLicense{value: NONE},
		},
	}
	output = lic.ToLicenseString()
	expectedOutput = "NOASSERTION AND NONE"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}

	// TestCase 4: nested conjunctive license.
	lic = ConjunctiveLicenseSet{
		members: []AnyLicenseInfo{
			SpecialLicense{value: NOASSERTION},
			ConjunctiveLicenseSet{
				members: []AnyLicenseInfo{
					SpecialLicense{value: "LicenseRef-1"},
					SpecialLicense{value: NONE},
				},
			},
		},
	}
	output = lic.ToLicenseString()
	expectedOutput = "NOASSERTION AND LicenseRef-1 AND NONE"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}
}

func TestDisjunctiveLicenseSet_ToLicenseString(t *testing.T) {
	var lic DisjunctiveLicenseSet
	var output, expectedOutput string

	// TestCase 1: no license in the set
	lic = DisjunctiveLicenseSet{
		members: nil,
	}
	output = lic.ToLicenseString()
	expectedOutput = ""
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}

	// TestCase 2: single license in the set
	lic = DisjunctiveLicenseSet{
		members: []AnyLicenseInfo{
			SpecialLicense{value: NOASSERTION},
		},
	}
	output = lic.ToLicenseString()
	expectedOutput = "NOASSERTION"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}

	// TestCase 3: more than one license in the set.
	lic = DisjunctiveLicenseSet{
		members: []AnyLicenseInfo{
			SpecialLicense{value: NOASSERTION},
			SpecialLicense{value: NONE},
		},
	}
	output = lic.ToLicenseString()
	expectedOutput = "NOASSERTION OR NONE"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}

	// TestCase 4: nested conjunctive license.
	lic = DisjunctiveLicenseSet{
		members: []AnyLicenseInfo{
			SpecialLicense{value: NOASSERTION},
			DisjunctiveLicenseSet{
				members: []AnyLicenseInfo{
					SpecialLicense{value: "LicenseRef-1"},
					SpecialLicense{value: NONE},
				},
			},
		},
	}
	output = lic.ToLicenseString()
	expectedOutput = "NOASSERTION OR LicenseRef-1 OR NONE"
	if output != expectedOutput {
		t.Errorf("expected: %s, found %s", output, expectedOutput)
	}
}

func TestExtractedLicensingInfo_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	extractedLicense := ExtractedLicensingInfo{
		SimpleLicensingInfo: SimpleLicensingInfo{
			licenseID: "license",
		},
		extractedText: "extracted Text",
	}
	expectedOutput := "license"
	output := extractedLicense.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}

func TestOrLaterOperator_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	orLater := OrLaterOperator{
		member: SimpleLicensingInfo{
			licenseID: "license",
		},
	}
	expectedOutput := "license"
	output := orLater.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}

func TestLicense_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	license := License{
		SimpleLicensingInfo: SimpleLicensingInfo{
			licenseID: "license",
		},
	}
	expectedOutput := "license"
	output := license.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}

func TestListedLicense_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	ll := ListedLicense{License{
		SimpleLicensingInfo: SimpleLicensingInfo{
			licenseID: "license",
		},
	},
	}
	expectedOutput := "license"
	output := ll.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}

func TestWithExceptionOperator_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	withException := WithExceptionOperator{
		member: SimpleLicensingInfo{
			licenseID: "license",
		},
		licenseException: LicenseException{},
	}
	expectedOutput := "license"
	output := withException.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}

func TestSpecialLicense_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	specialLicense := SpecialLicense{
		value: "license",
	}
	expectedOutput := "license"
	output := specialLicense.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}

func TestSimpleLicensingInfo_ToLicenseString(t *testing.T) {
	// nothing to test (just a dry run)
	sli := SimpleLicensingInfo{
		licenseID: "license",
	}
	expectedOutput := "license"
	output := sli.ToLicenseString()
	if output != expectedOutput {
		t.Errorf("expected: %s, found: %s", expectedOutput, output)
	}
}
