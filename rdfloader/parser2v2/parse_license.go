// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/RishabhBhatnagar/gordf/rdfloader/parser"
	"github.com/RishabhBhatnagar/gordf/rdfwriter"
	"strings"
)

func (parser *rdfParser2_2) getAnyLicenseFromNode(node *gordfParser.Node) (AnyLicenseInfo, error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	if len(associatedTriples) == 0 {
		// just a license uri string was found.
		return parser.getSpecialLicenseFromNode(node)
	}

	// we have some attributes associated with the license node.
	nodeType, err := parser.getNodeTypeFromTriples(associatedTriples, node)
	if err != nil {
		return nil, fmt.Errorf("error parsing license triple: %v", err)
	}
	switch nodeType {
	case SPDX_DISJUNCTIVE_LICENSE_SET:
		return parser.getDisjunctiveLicenseSetFromNode(node)
	case SPDX_CONJUNCTIVE_LICENSE_SET:
		return parser.getConjunctiveLicenseSetFromNode(node)
	case SPDX_EXTRACTED_LICENSING_INFO:
		return parser.getExtractedLicensingInfoFromNode(node)
	case SPDX_LISTED_LICENSE, SPDX_LICENSE:
		return parser.getLicenseFromNode(node)
	case SPDX_WITH_EXCEPTION_OPERATOR:
		return parser.getWithExceptionOperatorFromNode(node)
	case SPDX_OR_LATER_OPERATOR:
		return parser.getOrLaterOperatorFromNode(node)
	}
	return nil, fmt.Errorf("Unknown subTag (%s) found while parsing AnyLicense", nodeType)
}

func (parser *rdfParser2_2) getLicenseExceptionFromNode(node *gordfParser.Node) (exception LicenseException, err error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	for _, triple := range associatedTriples {
		value := triple.Object.ID
		switch triple.Predicate.ID {
		case SPDX_LICENSE_EXCEPTION_ID:
			exception.licenseExceptionId = value
		case SPDX_LICENSE_EXCEPTION_TEXT:
			exception.licenseExceptionText = value
		case RDFS_SEE_ALSO:
			if !isUriValid(value) {
				return exception, fmt.Errorf("invalid uri (%s) for seeAlso attribute of LicenseException", value)
			}
			exception.seeAlso = value
		case SPDX_NAME:
			exception.name = value
		case SPDX_EXAMPLE:
			exception.example = value
		default:
			return exception, fmt.Errorf("invalid predicate(%s) for LicenseException", triple.Predicate)
		}
	}
	return exception, nil
}

func (parser *rdfParser2_2) getSimpleLicensingInfoFromNode(node *gordfParser.Node) (SimpleLicensingInfo, error) {
	simpleLicensingTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	return parser.getSimpleLicensingInfoFromTriples(simpleLicensingTriples)
}

func (parser *rdfParser2_2) getWithExceptionOperatorFromNode(node *gordfParser.Node) (operator WithExceptionOperator, err error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	for _, triple := range associatedTriples {
		switch triple.Predicate.ID {
		case SPDX_MEMBER:
			member, err := parser.getSimpleLicensingInfoFromNode(triple.Object)
			if err != nil {
				return operator, fmt.Errorf("error parsing member of a WithExceptionOperator: %v", err)
			}
			operator.license = member
		case SPDX_LICENSE_EXCEPTION:
			operator.licenseException, err = parser.getLicenseExceptionFromNode(triple.Object)
		default:
			return operator, fmt.Errorf("unknown predicate (%s) for a WithExceptionOperator", triple.Predicate.ID)
		}
	}
	return operator, nil
}

func (parser *rdfParser2_2) getOrLaterOperatorFromNode(node *gordfParser.Node) (operator OrLaterOperator, err error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	n := len(associatedTriples)
	if n != 2 {
		return operator, fmt.Errorf("orLaterOperator must be associated with exactly one tag. found %v triples", n-1)
	}
	for _, triple := range associatedTriples {
		operator.license, err = parser.getSimpleLicensingInfoFromNode(triple.Object)
		if err != nil {
			return operator, fmt.Errorf("error parsing simpleLicensingInfo of OrLaterOperator: %v", err)
		}
	}
	return operator, nil
}

func (parser *rdfParser2_2) getSpecialLicenseFromNode(node *gordfParser.Node) (lic SpecialLicense, err error) {
	uri := strings.TrimSpace(node.ID)
	switch uri {
	case SPDX_NONE_CAPS, SPDX_NONE_SMALL:
		return SpecialLicense{
			value: NONE,
		}, nil
	case SPDX_NOASSERTION_SMALL, SPDX_NOASSERTION_CAPS:
		return SpecialLicense{
			value: NOASSERTION,
		}, nil
	}

	// the license is neither NONE nor NOASSERTION
	// checking if the license is among the standardLicenses
	licenseAbbreviation := getLastPartOfURI(uri)
	for _, stdLicense := range AllStandardLicenseIDS() {
		if licenseAbbreviation == stdLicense {
			return SpecialLicense{
				value: SpecialLicenseValue(stdLicense),
			}, nil
		}
	}
	return lic, fmt.Errorf("found a custom license uri (%s) without any associated fields", uri)
}

func (parser *rdfParser2_2) getDisjunctiveLicenseSetFromNode(node *gordfParser.Node) (DisjunctiveLicenseSet, error) {
	licenseSet := DisjunctiveLicenseSet{
		members: []AnyLicenseInfo{},
	}
	for _, triple := range parser.nodeToTriples(node) {
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_MEMBER:
			member, err := parser.getAnyLicenseFromNode(triple.Object)
			if err != nil {
				return licenseSet, fmt.Errorf("error parsing disjunctive license set: %v", err)
			}
			licenseSet.members = append(licenseSet.members, member)
		}
	}
	return licenseSet, nil
}

func (parser *rdfParser2_2) getConjunctiveLicenseSetFromNode(node *gordfParser.Node) (ConjunctiveLicenseSet, error) {
	licenseSet := ConjunctiveLicenseSet{
		members: []AnyLicenseInfo{},
	}
	for _, triple := range parser.nodeToTriples(node) {
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_MEMBER:
			member, err := parser.getAnyLicenseFromNode(triple.Object)
			if err != nil {
				return licenseSet, fmt.Errorf("error parsing conjunctive license set: %v", err)
			}
			licenseSet.members = append(licenseSet.members, member)
		}
	}
	return licenseSet, nil
}

func (parser *rdfParser2_2) getSimpleLicensingInfoFromTriples(triples []*gordfParser.Triple) (lic SimpleLicensingInfo, err error) {
	for _, triple := range triples {
		switch triple.Predicate.ID {
		case RDFS_COMMENT:
			lic.comment = triple.Object.ID
		case SPDX_LICENSE_ID:
			lic.licenseID = triple.Object.ID
		case SPDX_NAME:
			lic.name = triple.Object.ID
		case RDFS_SEE_ALSO:
			lic.seeAlso = append(lic.seeAlso, triple.Object.ID)
		case SPDX_EXAMPLE:
			lic.example = triple.Object.ID
		case RDF_TYPE:
			continue
		default:
			return lic, fmt.Errorf("unknown predicate(%s) for simple licensing info", triple.Predicate)
		}
	}
	return lic, nil
}

func (parser *rdfParser2_2) getLicenseFromNode(node *gordfParser.Node) (lic License, err error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	var restTriples []*gordfParser.Triple
	for _, triple := range associatedTriples {
		value := triple.Object.ID
		switch triple.Predicate.ID {
		case SPDX_IS_OSI_APPROVED:
			lic.isOsiApproved, err = boolFromString(value)
			if err != nil {
				return lic, fmt.Errorf("error parsing isOsiApproved attribute of a License: %v", err)
			}
		case SPDX_LICENSE_TEXT:
			lic.licenseText = value
		case SPDX_STANDARD_LICENSE_HEADER:
			lic.standardLicenseHeader = value
		case SPDX_STANDARD_LICENSE_TEMPLATE:
			lic.standardLicenseTemplate = value
		case SPDX_STANDARD_LICENSE_HEADER_TEMPLATE:
			lic.standardLicenseHeaderTemplate = value
		case RDFS_SEE_ALSO:
			if !isUriValid(value) {
				return lic, fmt.Errorf("%s is not a valid uri for seeAlso attribute of a License", value)
			}
			lic.seeAlso = value

		case SPDX_IS_DEPRECATED_LICENSE_ID:
			lic.isDeprecatedLicenseID, err = boolFromString(value)
			if err != nil {
				return lic, fmt.Errorf("error parsing isDeprecatedLicenseId attribute of a License: %v", err)
			}
		case SPDX_IS_FSF_LIBRE:
			lic.isFsfLibre, err = boolFromString(value)
			if err != nil {
				return lic, fmt.Errorf("error parsing isFsfLibre attribute of a License: %v", err)
			}
		default:
			restTriples = append(restTriples, triple)
		}
	}
	lic.SimpleLicensingInfo, err = parser.getSimpleLicensingInfoFromTriples(restTriples)
	if err != nil {
		return lic, fmt.Errorf("error setting simple licensing information of a License: %s", err)
	}
	return lic, nil
}

/* util methods for licenses and checksums below:*/

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
	for _, checksumTriple := range parser.nodeToTriples(checksumNode) {
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

func mapLicensesToStrings(licences []AnyLicenseInfo) []string {
	res := make([]string, len(licences), len(licences))
	for i, lic := range licences {
		res[i] = lic.ToLicenseString()
	}
	return res
}

func (lic ConjunctiveLicenseSet) ToLicenseString() string {
	return strings.Join(mapLicensesToStrings(lic.members), " AND ")
}

func (lic DisjunctiveLicenseSet) ToLicenseString() string {
	return strings.Join(mapLicensesToStrings(lic.members), " OR ")
}

/****** Type Functions ******/
func (lic ExtractedLicensingInfo) ToLicenseString() string {
	return lic.licenseID
}

func (operator OrLaterOperator) ToLicenseString() string {
	return operator.license.ToLicenseString()
}

func (lic License) ToLicenseString() string {
	return lic.licenseID
}

func (lic ListedLicense) ToLicenseString() string {
	return lic.licenseID
}

func (lic WithExceptionOperator) ToLicenseString() string {
	return lic.license.ToLicenseString()
}

func (lic SpecialLicense) ToLicenseString() string {
	return string(lic.value)
}
