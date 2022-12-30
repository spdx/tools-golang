// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v3

import (
	"testing"
)

// ===== Utility function tests =====
func TestTextifyWrapsStringWithNewline(t *testing.T) {
	s := `this text has
a newline in it`
	want := `<text>this text has
a newline in it</text>`

	got := textify(s)

	if want != got {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestTextifyDoesNotWrapsStringWithNoNewline(t *testing.T) {
	s := `this text has no newline in it`
	want := s

	got := textify(s)

	if want != got {
		t.Errorf("expected %s, got %s", want, got)
	}
}
