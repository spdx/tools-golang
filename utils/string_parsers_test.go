// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package utils

import (
	"testing"
)

// ===== Helper function tests =====

func TestCanExtractSubvalues(t *testing.T) {
	subkey, subvalue, err := ExtractSubs("SHA1: abc123")
	if err != nil {
		t.Errorf("got error when calling utils.ExtractSubs: %v", err)
	}
	if subkey != "SHA1" {
		t.Errorf("got %v for subkey", subkey)
	}
	if subvalue != "abc123" {
		t.Errorf("got %v for subvalue", subvalue)
	}
}

func TestReturnsErrorForInvalidSubvalueFormat(t *testing.T) {
	_, _, err := ExtractSubs("blah")
	if err == nil {
		t.Errorf("expected error when calling utils.ExtractSubs for invalid format (0 colons), got nil")
	}
}
