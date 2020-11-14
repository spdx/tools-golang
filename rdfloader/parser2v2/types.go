// copied from tvloader/parser2v2/types.go
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v2

import (
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
)

type rdfParser2_2 struct {
	// fields associated with gordf project which
	// will be required by rdfloader
	gordfParserObj      *gordfParser.Parser
	nodeStringToTriples map[string][]*gordfParser.Triple

	// document into which data is being parsed
	doc *spdx.Document2_2

	// map of packages and files.
	files            map[spdx.ElementID]*spdx.File2_2
	assocWithPackage map[spdx.ElementID]bool
	packages         map[spdx.ElementID]*spdx.Package2_2

	// mapping of nodeStrings to parsed object to save double computation.
	cache map[string]interface{}
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
