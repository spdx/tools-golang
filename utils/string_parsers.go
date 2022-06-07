package utils

import (
	"fmt"
	"strings"
)

// ExtractSubsWithSeparator used to extract key / value from embedded substrings that use an arbitrary separator.
// returns subkey, subvalue, nil if no error, or "", "", error otherwise
func ExtractSubsWithSeparator(value string, sep string) (string, string, error) {
	// parse the value to see if it's a valid subvalue format
	sp := strings.SplitN(value, sep, 2)
	if len(sp) == 1 {
		return "", "", fmt.Errorf("invalid subvalue format for %s (no %s found)", value, sep)
	}

	subkey := strings.TrimSpace(sp[0])
	subvalue := strings.TrimSpace(sp[1])

	return subkey, subvalue, nil
}

// ExtractSubs used to extract key / value from embedded substrings that use a colon ":" as the separator.
// returns subkey, subvalue, nil if no error, or "", "", error otherwise
func ExtractSubs(value string) (string, string, error) {
	return ExtractSubsWithSeparator(value, ":")
}
