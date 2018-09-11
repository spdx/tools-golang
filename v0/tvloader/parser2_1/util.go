// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"fmt"
	"strings"
)

// used to extract key / value from embedded substrings
// returns subkey, subvalue, nil if no error, or "", "", error otherwise
func extractSubs(value string) (string, string, error) {
	// parse the value to see if it's a valid subvalue format
	sp := strings.SplitN(value, ":", 2)
	if len(sp) == 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (no colon found)", value)
	}

	subkey := strings.TrimSpace(sp[0])
	subvalue := strings.TrimSpace(sp[1])

	return subkey, subvalue, nil
}
