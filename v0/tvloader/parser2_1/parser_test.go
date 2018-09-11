// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import (
	"testing"
)

func TestParser2_1InitCreatesResetStatus(t *testing.T) {
	parser := tvParser2_1{}
	if parser.st != psStart2_1 {
		t.Errorf("parser did not begin in start state")
	}
	if parser.doc != nil {
		t.Errorf("parser did not begin with nil document")
	}
}

func TestParser2_1HasDocumentAfterCallToParseFirstTag(t *testing.T) {
	parser := tvParser2_1{}
	err := parser.parsePair2_1("SPDXVersion", "SPDX-2.1")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.doc == nil {
		t.Errorf("doc is still nil after parsing first pair")
	}
}

// ===== Parser start state change tests =====
func TestParser2_1StartMovesToCreationInfoStateAfterParsingFirstTag(t *testing.T) {
	parser := tvParser2_1{}
	err := parser.parsePair2_1("SPDXVersion", "b")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psCreationInfo2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psCreationInfo2_1)
	}
}
