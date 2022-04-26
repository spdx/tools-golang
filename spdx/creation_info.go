// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Creator is a wrapper around the Creator SPDX field. The SPDX field contains two values, which requires special
// handling in order to marshal/unmarshal it to/from Go data types.
type Creator struct {
	Creator string
	// CreatorType should be one of "Person", "Organization", or "Tool"
	CreatorType string
}

// UnmarshalJSON takes an annotator in the typical one-line format and parses it into a Creator struct.
// This function is also used when unmarshalling YAML
func (c *Creator) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, "\"")
	fields := strings.SplitN(str, ": ", 2)

	if len(fields) != 2 {
		return fmt.Errorf("failed to parse Creator '%s'", str)
	}

	c.CreatorType = fields[0]
	c.Creator = fields[1]

	return nil
}

// MarshalJSON converts the receiver into a slice of bytes representing a Creator in string form.
// This function is also used with marshalling to YAML
func (c Creator) MarshalJSON() ([]byte, error) {
	if c.Creator != "" {
		return json.Marshal(fmt.Sprintf("%s: %s", c.CreatorType, c.Creator))
	}

	return []byte{}, nil
}

// CreationInfo2_1 is a Document Creation Information section of an
// SPDX Document for version 2.1 of the spec.
type CreationInfo2_1 struct {
	// 2.7: License List Version
	// Cardinality: optional, one
	LicenseListVersion string `json:"licenseListVersion"`

	// 2.8: Creators: may have multiple keys for Person, Organization
	//      and/or Tool
	// Cardinality: mandatory, one or many
	Creators []Creator `json:"creators"`

	// 2.9: Created: data format YYYY-MM-DDThh:mm:ssZ
	// Cardinality: mandatory, one
	Created string `json:"created"`

	// 2.10: Creator Comment
	// Cardinality: optional, one
	CreatorComment string `json:"comment"`
}

// CreationInfo2_2 is a Document Creation Information section of an
// SPDX Document for version 2.2 of the spec.
type CreationInfo2_2 struct {
	// 2.7: License List Version
	// Cardinality: optional, one
	LicenseListVersion string `json:"licenseListVersion"`

	// 2.8: Creators: may have multiple keys for Person, Organization
	//      and/or Tool
	// Cardinality: mandatory, one or many
	Creators []Creator `json:"creators"`

	// 2.9: Created: data format YYYY-MM-DDThh:mm:ssZ
	// Cardinality: mandatory, one
	Created string `json:"created"`

	// 2.10: Creator Comment
	// Cardinality: optional, one
	CreatorComment string `json:"comment"`
}
