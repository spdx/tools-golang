// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
// copied from tvloader/parser2v2/types.go
package reader

import (
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

type rdfParser2_2 struct {
	// fields associated with gordf project which
	// will be required by rdfloader
	gordfParserObj      *gordfParser.Parser
	nodeStringToTriples map[string][]*gordfParser.Triple

	// document into which data is being parsed
	doc *v2_2.Document

	// map of packages and files.
	files            map[common.ElementID]*v2_2.File
	assocWithPackage map[common.ElementID]bool

	// mapping of nodeStrings to parsed object to save double computation.
	cache map[string]*nodeState
}

type Color int

const (
	GREY  Color = iota // represents that the node is being visited
	WHITE              // unvisited node
	BLACK              // visited node
)

type nodeState struct {
	// object will be pointer to the parsed or element being parsed.
	object interface{}
	// color of a state represents if the node is visited/unvisited/being-visited.
	Color Color
}

type AnyLicenseInfo interface {
	// ToLicenseString returns the representation of license about how it will
	// be stored in the tools-golang data model
	ToLicenseString() string
}

type SimpleLicensingInfo struct {
	AnyLicenseInfo
	comment   string
	licenseID string
	name      string
	seeAlso   []string
	example   string
}

type ExtractedLicensingInfo struct {
	SimpleLicensingInfo
	extractedText string
}

type OrLaterOperator struct {
	AnyLicenseInfo
	member SimpleLicensingInfo
}

type ConjunctiveLicenseSet struct {
	AnyLicenseInfo
	members []AnyLicenseInfo
}

type DisjunctiveLicenseSet struct {
	AnyLicenseInfo
	members []AnyLicenseInfo
}

type License struct {
	SimpleLicensingInfo
	isOsiApproved                 bool
	licenseText                   string
	standardLicenseHeader         string
	standardLicenseTemplate       string
	standardLicenseHeaderTemplate string
	isDeprecatedLicenseID         bool
	isFsfLibre                    bool
}

type ListedLicense struct {
	License
}

type LicenseException struct {
	licenseExceptionId   string
	licenseExceptionText string
	seeAlso              string // must be a valid uri
	name                 string
	example              string
	comment              string
}

type WithExceptionOperator struct {
	AnyLicenseInfo
	member           SimpleLicensingInfo
	licenseException LicenseException
}

// custom LicenseType to provide support for licences of
// type Noassertion, None and customLicenses
type SpecialLicense struct {
	AnyLicenseInfo
	value SpecialLicenseValue
}

type SpecialLicenseValue string

const (
	NONE        SpecialLicenseValue = "NONE"
	NOASSERTION SpecialLicenseValue = "NOASSERTION"
)

type RangeType string

const (
	BYTE_RANGE RangeType = "byteRange"
	LINE_RANGE RangeType = "lineRange"
)
