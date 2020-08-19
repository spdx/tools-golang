// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/RishabhBhatnagar/gordf/rdfwriter"
	"strings"
)

// either single tag or a compound license with member combination of single tags.
// todo: allow all types of licenses.
func (parser *rdfParser2_2) getLicenseFromTriple(triple *gordfParser.Triple) (licenseConcluded string, err error) {
	licenseShortIdentifier := getLicenseStringFromURI(triple.Object.ID)
	// return if license is None|Noassertion
	if licenseShortIdentifier == "NONE" || licenseShortIdentifier == "NOASSERTION" {
		return licenseShortIdentifier, nil
	}

	// return if the license tag is not associated with any other triples.
	if len(parser.nodeToTriples[triple.Object.String()]) == 0 {
		return licenseShortIdentifier, nil
	}

	// no need to parse standard licenses as they have constant fields.
	// return if the license is among the standard licenses.
	for _, stdLicenseId := range AllStandardLicenseIDS() {
		if stdLicenseId == licenseShortIdentifier {
			return licenseShortIdentifier, nil
		}
	}

	// apart from the license being in the uri form, this function allows
	// license to be a collection of licenses joined by a single operator
	// (either conjunction or disjunction)

	typeTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &triple.Object.ID, &RDF_TYPE, nil)
	if len(typeTriples) == 0 {
		return "", fmt.Errorf("node(%v) not associated with a type triple", triple.Object)
	}
	if len(typeTriples) > 1 {
		return "", fmt.Errorf("node is associated with more than one type triple")
	}
	switch typeTriples[0].Object.ID {
	case SPDX_DISJUNCTIVE_LICENSE_SET, SPDX_CONJUNCTIVE_LICENSE_SET:

	case SPDX_EXTRACTED_LICENSING_INFO:
		err = parser.parseOtherLicenseFromNode(triple.Object)
		if err != nil {
			return "", err
		}
		othLic := parser.doc.OtherLicenses[len(parser.doc.OtherLicenses)-1]
		return othLic.LicenseIdentifier, nil
	default:
		return "", fmt.Errorf("not implemented error: cannot parse %s", typeTriples[0].Object)
	}
	return parser.getLicenseFromLicenseSetNode(triple.Object)
}

// Given the license URI, returns the name of the license defined
// in the last part of the uri.
// This function is susceptible to false-positives.
func getLicenseStringFromURI(uri string) string {
	licenseEnd := strings.TrimSpace(getLastPartOfURI(uri))
	lower := strings.ToLower(licenseEnd)
	if lower == "none" || lower == "noassertion" {
		return strings.ToUpper(licenseEnd)
	}
	return licenseEnd
}

// returns the checksum algorithm and it's value
// In the newer versions, these two strings will be bound to a single checksum struct
// whose pointer will be returned.
func (parser *rdfParser2_2) getChecksumFromNode(checksumNode *gordfParser.Node) (algorithm string, value string, err error) {
	var checksumValue, checksumAlgorithm string
	for _, checksumTriple := range parser.nodeToTriples[checksumNode.String()] {
		switch checksumTriple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_CHECKSUM_VALUE:
			// cardinality: exactly 1
			checksumValue = strings.TrimSpace(checksumTriple.Object.ID)
		case SPDX_ALGORITHM:
			// cardinality: exactly 1
			checksumAlgorithm, err = parser.getAlgorithmFromURI(checksumTriple.Object.ID)
			if err != nil {
				return
			}
		}
	}
	return checksumAlgorithm, checksumValue, nil
}

func (parser *rdfParser2_2) getAlgorithmFromURI(algorithmURI string) (checksumAlgorithm string, err error) {
	fragment := getLastPartOfURI(algorithmURI)
	if !strings.HasPrefix(fragment, "checksumAlgorithm") {
		return "", fmt.Errorf("checksum algorithm uri must begin with checksumAlgorithm. found %s", fragment)
	}
	algorithm := strings.TrimPrefix(fragment, "checksumAlgorithm_")
	algorithm = strings.ToLower(strings.TrimSpace(algorithm))
	switch algorithm {
	case "md2", "md4", "md5", "md6":
		checksumAlgorithm = strings.ToUpper(algorithm)
	case "sha1", "sha224", "sha256", "sha384", "sha512":
		checksumAlgorithm = strings.ToUpper(algorithm)
	default:
		return "", fmt.Errorf("unknown checksum algorithm %s", algorithm)
	}
	return
}

func (parser *rdfParser2_2) getLicenseFromLicenseSetNode(node *gordfParser.Node) (s string, err error) {
	typeLicenseSet := "undefined"
	var licenseSets []string
	for _, lst := range parser.nodeToTriples[node.String()] {
		switch lst.Predicate.ID {
		case RDF_TYPE:
			_, typeLicenseSet, err = ExtractSubs(lst.Object.ID, "#")
			if err != nil {
				return
			}
		case SPDX_MEMBER:
			licenseSets = append(licenseSets, getLicenseStringFromURI(lst.Object.ID))
		default:
			return "", fmt.Errorf("undefined predicate %s while parsing license set", lst.Predicate.ID)
		}
	}
	switch typeLicenseSet {
	case "DisjunctiveLicenseSet":
		return strings.Join(licenseSets, " OR "), nil
	case "ConjunctiveLicenseSet":
		return strings.Join(licenseSets, " AND "), nil
	default:
		return "", fmt.Errorf("unknown licenseSet type %s", typeLicenseSet)
	}
}
