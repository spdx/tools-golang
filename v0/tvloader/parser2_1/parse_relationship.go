// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package tvloader

import (
	"fmt"
	"strings"
)

func (parser *tvParser2_1) parsePairForRelationship2_1(tag string, value string) error {
	if parser.rln == nil {
		return fmt.Errorf("no relationship struct created in parser rln pointer")
	}

	if tag == "Relationship" {
		// parse the value to see if it's a valid relationship format
		sp := strings.SplitN(value, " ", -1)

		// filter out any purely-whitespace items
		var rp []string
		for _, v := range sp {
			v = strings.TrimSpace(v)
			if v != "" {
				rp = append(rp, v)
			}
		}

		if len(rp) != 3 {
			return fmt.Errorf("invalid relationship format for %s", value)
		}

		parser.rln.RefA = strings.TrimSpace(rp[0])
		parser.rln.Relationship = strings.TrimSpace(rp[1])
		parser.rln.RefB = strings.TrimSpace(rp[2])
		return nil
	}

	if tag == "RelationshipComment" {
		parser.rln.RelationshipComment = value
		return nil
	}

	return fmt.Errorf("received unknown tag %v in Relationship section", tag)
}
