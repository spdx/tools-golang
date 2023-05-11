// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"
	"strings"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

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
func (parser *rdfParser2_3) getChecksumFromNode(checksumNode *gordfParser.Node) (algorithm common.ChecksumAlgorithm, value string, err error) {
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
			checksumAlgorithm, err = getAlgorithmFromURI(checksumTriple.Object.ID)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("unknown predicate '%s' while parsing checksum node", checksumTriple.Predicate.ID)
			return
		}
	}
	return common.ChecksumAlgorithm(checksumAlgorithm), checksumValue, nil
}

func getAlgorithmFromURI(algorithmURI string) (checksumAlgorithm string, err error) {
	fragment := getLastPartOfURI(algorithmURI)
	if !strings.HasPrefix(fragment, "checksumAlgorithm_") {
		return "", fmt.Errorf("checksum algorithm uri must begin with checksumAlgorithm_. found %s", fragment)
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

// from a list of licenses, it returns a
// list of string representation of those licenses.
func mapLicensesToStrings(licences []AnyLicenseInfo) []string {
	res := make([]string, len(licences))
	for i, lic := range licences {
		res[i] = lic.ToLicenseString()
	}
	return res
}

/****** Type Functions ******/

// TODO: should probably add brackets while linearizing a nested license.
func (lic ConjunctiveLicenseSet) ToLicenseString() string {
	return strings.Join(mapLicensesToStrings(lic.members), " AND ")
}

// TODO: should probably add brackets while linearizing a nested license.
func (lic DisjunctiveLicenseSet) ToLicenseString() string {
	return strings.Join(mapLicensesToStrings(lic.members), " OR ")
}

func (lic ExtractedLicensingInfo) ToLicenseString() string {
	return lic.licenseID
}

func (operator OrLaterOperator) ToLicenseString() string {
	return operator.member.ToLicenseString()
}

func (lic License) ToLicenseString() string {
	return lic.licenseID
}

func (lic ListedLicense) ToLicenseString() string {
	return lic.licenseID
}

func (lic WithExceptionOperator) ToLicenseString() string {
	return lic.member.ToLicenseString()
}

func (lic SpecialLicense) ToLicenseString() string {
	return string(lic.value)
}

func (lic SimpleLicensingInfo) ToLicenseString() string {
	return lic.licenseID
}
