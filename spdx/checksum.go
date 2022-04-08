// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

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
