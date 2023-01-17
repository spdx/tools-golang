// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"errors"
	"fmt"
	"strings"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/gordf/rdfwriter"
)

// AnyLicense is a baseClass for all the licenses
// All the types of licenses is a sub-type of AnyLicense,
// either directly or indirectly.
// This function acts as a mux for all the licenses. Based on the input, it
// decides which type of license it is and passes control to that type of
// license parser to parse the given input.
func (parser *rdfParser2_2) getAnyLicenseFromNode(node *gordfParser.Node) (AnyLicenseInfo, error) {

	currState := parser.cache[node.ID]
	if currState == nil {
		// there is no entry about the state of current package node.
		// this is the first time we're seeing this node.
		parser.cache[node.ID] = &nodeState{
			object: nil, // not storing the object as we won't retrieve it later.
			Color:  WHITE,
		}
	} else if currState.Color == GREY {
		// we have already started parsing this license node.
		// We have a cyclic dependency!
		return nil, errors.New("Couldn't parse license: found a cyclic dependency on " + node.ID)
	}

	// setting color of the state to grey to indicate that we've started to
	// parse this node once.
	parser.cache[node.ID].Color = GREY

	// setting state color to black when we're done parsing this node.
	defer func() { parser.cache[node.ID].Color = BLACK }()

	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	if len(associatedTriples) == 0 {
		// just a license uri string was found.
		return parser.getSpecialLicenseFromNode(node)
	}

	// we have some attributes associated with the license node.
	nodeType, err := getNodeTypeFromTriples(associatedTriples, node)
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
	case SPDX_SIMPLE_LICENSING_INFO:
		return parser.getSimpleLicensingInfoFromNode(node)
	}
	return nil, fmt.Errorf("Unknown subTag (%s) found while parsing AnyLicense", nodeType)
}

func (parser *rdfParser2_2) getLicenseExceptionFromNode(node *gordfParser.Node) (exception LicenseException, err error) {
	associatedTriples := rdfwriter.FilterTriples(parser.gordfParserObj.Triples, &node.ID, nil, nil)
	for _, triple := range associatedTriples {
		value := triple.Object.ID
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
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
		case RDFS_COMMENT:
			exception.comment = value
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
	var memberFound bool
	for _, triple := range associatedTriples {
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_MEMBER:
			if memberFound {
				return operator,
					fmt.Errorf("more than one member found in the WithExceptionOperator (expected only 1)")
			}
			memberFound = true
			member, err := parser.getSimpleLicensingInfoFromNode(triple.Object)
			if err != nil {
				return operator, fmt.Errorf("error parsing member of a WithExceptionOperator: %v", err)
			}
			operator.member = member
		case SPDX_LICENSE_EXCEPTION:
			operator.licenseException, err = parser.getLicenseExceptionFromNode(triple.Object)
			if err != nil {
				return operator, fmt.Errorf("error parsing licenseException of WithExceptionOperator: %v", err)
			}
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
		switch triple.Predicate.ID {
		case RDF_TYPE:
			continue
		case SPDX_MEMBER:
			operator.member, err = parser.getSimpleLicensingInfoFromNode(triple.Object)
			if err != nil {
				return operator, fmt.Errorf("error parsing simpleLicensingInfo of OrLaterOperator: %v", err)
			}
		default:
			return operator, fmt.Errorf("unknown predicate %s", triple.Predicate.ID)
		}
	}
	return operator, nil
}

// SpecialLicense is a type of license which is not defined in any of the
// spdx documents, it is a type of license defined for the sake of brevity.
// It can be [NONE|NOASSERTION|LicenseRef-<string>]
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
		default:
			return licenseSet, fmt.Errorf("unknown subTag for ConjunctiveLicenseSet: %s", triple.Predicate.ID)
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
			if !isUriValid(triple.Object.ID) {
				return lic, fmt.Errorf("%s is not a valid uri for seeAlso attribute of a License", triple.Object.ID)
			}
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
