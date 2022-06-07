// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"encoding/json"
	"fmt"
	"github.com/spdx/tools-golang/utils"
	"strings"
)

// Creator is a wrapper around the Creator SPDX field. The SPDX field contains two values, which requires special
// handling in order to marshal/unmarshal it to/from Go data types.
type Creator struct {
	Creator string
	// CreatorType should be one of "Person", "Organization", or "Tool"
	CreatorType string
}

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (c Creator) Validate() error {
	if c.CreatorType == "" || c.Creator == "" {
		return fmt.Errorf("invalid Creator, missing fields. %+v", c)
	}

	return nil
}

// FromString takes a Creator in the typical one-line format and parses it into a Creator struct.
func (c *Creator) FromString(str string) error {
	creatorType, creator, err := utils.ExtractSubs(str)
	if err != nil {
		return err
	}

	c.CreatorType = creatorType
	c.Creator = creator

	return nil
}

// String converts the Creator into a string.
func (c Creator) String() string {
	return fmt.Sprintf("%s: %s", c.CreatorType, c.Creator)
}

// UnmarshalJSON takes a Creator in the typical one-line string format and parses it into a Creator struct.
func (c *Creator) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, "\"")

	return c.FromString(str)
}

// MarshalJSON converts the receiver into a slice of bytes representing a Creator in string form.
// This function is also used with marshalling to YAML
func (c Creator) MarshalJSON() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(c.String())
}

// CreationInfo2_1 is a Document Creation Information section of an
// SPDX Document for version 2.1 of the spec.
type CreationInfo2_1 struct {
	// 2.7: License List Version
	// Cardinality: optional, one
	LicenseListVersion string `json:"licenseListVersion,omitempty"`

	// 2.8: Creators: may have multiple keys for Person, Organization
	//      and/or Tool
	// Cardinality: mandatory, one or many
	Creators []Creator `json:"creators"`

	// 2.9: Created: data format YYYY-MM-DDThh:mm:ssZ
	// Cardinality: mandatory, one
	Created string `json:"created"`

	// 2.10: Creator Comment
	// Cardinality: optional, one
	CreatorComment string `json:"comment,omitempty"`
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
