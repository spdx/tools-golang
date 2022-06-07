// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"fmt"
	"github.com/spdx/tools-golang/utils"
)

// ChecksumAlgorithm represents the algorithm used to generate the file checksum in the Checksum struct.
type ChecksumAlgorithm string

// The checksum algorithms mentioned in the spdxv2.2.0 https://spdx.github.io/spdx-spec/4-file-information/#44-file-checksum
const (
	SHA224 ChecksumAlgorithm = "SHA224"
	SHA1   ChecksumAlgorithm = "SHA1"
	SHA256 ChecksumAlgorithm = "SHA256"
	SHA384 ChecksumAlgorithm = "SHA384"
	SHA512 ChecksumAlgorithm = "SHA512"
	MD2    ChecksumAlgorithm = "MD2"
	MD4    ChecksumAlgorithm = "MD4"
	MD5    ChecksumAlgorithm = "MD5"
	MD6    ChecksumAlgorithm = "MD6"
)

// Checksum provides a unique identifier to match analysis information on each specific file in a package.
// The Algorithm field describes the ChecksumAlgorithm used and the Value represents the file checksum
type Checksum struct {
	Algorithm ChecksumAlgorithm `json:"algorithm"`
	Value     string            `json:"checksumValue"`
}

// FromString parses a Checksum string into a spdx.Checksum.
// These strings take the following form:
// SHA1: d6a770ba38583ed4bb4525bd96e50461655d2759
func (c *Checksum) FromString(value string) error {
	algorithm, value, err := utils.ExtractSubs(value)
	if err != nil {
		return err
	}

	c.Algorithm = ChecksumAlgorithm(algorithm)
	c.Value = value

	return nil
}

// String converts the Checksum to its string form.
// e.g. "SHA1: d6a770ba38583ed4bb4525bd96e50461655d2759"
func (c Checksum) String() string {
	return fmt.Sprintf("%s: %s", c.Algorithm, c.Value)
}

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (c Checksum) Validate() error {
	if c.Algorithm == "" || c.Value == "" {
		return fmt.Errorf("invalid checksum, missing field(s). %+v", c)
	}

	return nil
}
