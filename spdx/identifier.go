// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdx

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ElementID represents the identifier string portion of an SPDX element
// identifier. DocElementID should be used for any attributes which can
// contain identifiers defined in a different SPDX document.
// ElementIDs should NOT contain the mandatory 'SPDXRef-' portion.
type ElementID string

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (e ElementID) Validate() error {
	if e == "" {
		return fmt.Errorf("invalid ElementID, must not be blank")
	}

	return nil
}

func (e ElementID) String() string {
	return fmt.Sprintf("SPDXRef-%s", string(e))
}

// FromString parses an SPDX Identifier string into an ElementID.
// These strings take the form: "SPDXRef-some-identifier"
func (e *ElementID) FromString(idStr string) error {
	idFields := strings.SplitN(idStr, "SPDXRef-", 2)
	switch len(idFields) {
	case 2:
		// "SPDXRef-" prefix was present
		*e = ElementID(idFields[1])
	case 1:
		// prefix was not present
		*e = ElementID(idFields[0])
	}

	return nil
}

// UnmarshalJSON takes a SPDX Identifier string parses it into an ElementID.
// This function is also used when unmarshalling YAML
func (e *ElementID) UnmarshalJSON(data []byte) error {
	// SPDX identifier will simply be a string
	idStr := string(data)
	idStr = strings.Trim(idStr, "\"")

	return e.FromString(idStr)
}

// MarshalJSON converts the receiver into a slice of bytes representing an ElementID in string form.
// This function is also used when marshalling to YAML
func (e ElementID) MarshalJSON() ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(e.String())
}

// DocElementID represents an SPDX element identifier that could be defined
// in a different SPDX document, and therefore could have a "DocumentRef-"
// portion, such as Relationships and Annotations.
// ElementID is used for attributes in which a "DocumentRef-" portion cannot
// appear, such as a Package or File definition (since it is necessarily
// being defined in the present document).
// DocumentRefID will be an empty string for elements defined in the
// present document.
// DocElementIDs should NOT contain the mandatory 'DocumentRef-' or
// 'SPDXRef-' portions.
// SpecialID is used ONLY if the DocElementID matches a defined set of
// permitted special values for a particular field, e.g. "NONE" or
// "NOASSERTION" for the right-hand side of Relationships. If SpecialID
// is set, DocumentRefID and ElementRefID should be empty (and vice versa).
type DocElementID struct {
	DocumentRefID string
	ElementRefID  ElementID
	SpecialID     string
}

// Validate verifies that all the required fields are present.
// Returns an error if the object is invalid.
func (d DocElementID) Validate() error {
	if d.DocumentRefID == "" && d.ElementRefID.Validate() != nil && d.SpecialID == "" {
		return fmt.Errorf("invalid DocElementID, missing fields. %+v", d)
	}

	return nil
}

// FromString parses an SPDX Identifier string into a DocElementID struct.
// These strings take one of the following forms:
//  - "DocumentRef-other-document:SPDXRef-some-identifier"
//  - "SPDXRef-some-identifier"
//  - "NOASSERTION" or "NONE"
func (d *DocElementID) FromString(idStr string) error {
	// handle special cases
	if idStr == "NONE" || idStr == "NOASSERTION" {
		d.SpecialID = idStr
		return nil
	}

	var idFields []string
	// handle DocumentRef- if present
	if strings.HasPrefix(idStr, "DocumentRef-") {
		// strip out the "DocumentRef-" so we can get the value
		idFields = strings.SplitN(idStr, "DocumentRef-", 2)
		idStr = idFields[1]

		// an SPDXRef can appear after a DocumentRef, separated by a colon
		idFields = strings.SplitN(idStr, ":", 2)
		d.DocumentRefID = idFields[0]

		if len(idFields) == 2 {
			idStr = idFields[1]
		} else {
			return nil
		}
	}

	// handle SPDXRef-
	err := d.ElementRefID.FromString(idStr)
	if err != nil {
		return err
	}

	return nil
}

// MarshalString converts the receiver into a string representing a DocElementID.
// This is used when writing a spreadsheet SPDX file, for example.
func (d DocElementID) String() string {
	if d.DocumentRefID != "" && d.ElementRefID != "" {
		return fmt.Sprintf("DocumentRef-%s:%s", d.DocumentRefID, d.ElementRefID)
	} else if d.DocumentRefID != "" {
		return fmt.Sprintf("DocumentRef-%s", d.DocumentRefID)
	} else if d.ElementRefID != "" {
		return d.ElementRefID.String()
	} else if d.SpecialID != "" {
		return d.SpecialID
	}

	return ""
}

// UnmarshalJSON takes a SPDX Identifier string parses it into a DocElementID struct.
// This function is also used when unmarshalling YAML
func (d *DocElementID) UnmarshalJSON(data []byte) error {
	// SPDX identifier will simply be a string
	idStr := string(data)
	idStr = strings.Trim(idStr, "\"")

	return d.FromString(idStr)
}

// MarshalJSON converts the receiver into a slice of bytes representing a DocElementID in string form.
// This function is also used when marshalling to YAML
func (d DocElementID) MarshalJSON() ([]byte, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(d.String())
}

// MakeDocElementID takes strings for the DocumentRef- and SPDXRef- identifiers (these prefixes will be stripped if present),
// and returns a DocElementID.
// An empty string should be used for the DocumentRef- portion if it is referring to the present document.
func MakeDocElementID(docRef string, eltRef string) DocElementID {
	docRef = strings.Replace(docRef, "DocumentRef-", "", 1)
	eltRef = strings.Replace(eltRef, "SPDXRef-", "", 1)

	return DocElementID{
		DocumentRefID: docRef,
		ElementRefID:  ElementID(eltRef),
	}
}

// MakeDocElementSpecial takes a "special" string (e.g. "NONE" or
// "NOASSERTION" for the right side of a Relationship), nd returns
// a DocElementID with it in the SpecialID field. Other fields will
// be empty.
func MakeDocElementSpecial(specialID string) DocElementID {
	return DocElementID{SpecialID: specialID}
}
