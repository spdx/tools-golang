// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package idsearcher

import (
	"testing"
)

// ===== Builder top-level Document test =====
func TestCanFindShortFormIDWhenPresent(t *testing.T) {
	filePath := "../testdata/project2/has-id.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("expected len 1, got %d", len(ids))
	}

	if ids[0] != "Apache-2.0 OR GPL-2.0-or-later" {
		t.Errorf("expected %v, got %v", "Apache-2.0 OR GPL-2.0-or-later", ids[0])
	}
}

func TestCanFindMultipleShortFormIDsWhenPresent(t *testing.T) {
	filePath := "../testdata/project2/has-multiple-ids.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("expected len 3, got %d", len(ids))
	}

	if ids[0] != "(MIT AND BSD-3-Clause) OR ISC" {
		t.Errorf("expected %v, got %v", "(MIT AND BSD-3-Clause) OR ISC", ids[0])
	}
	if ids[1] != "BSD-2-Clause" {
		t.Errorf("expected %v, got %v", "BSD-2-Clause", ids[1])
	}
	if ids[2] != "CC0-1.0" {
		t.Errorf("expected %v, got %v", "CC0-1.0", ids[2])
	}
}

func TestCannotFindShortFormIDWhenAbsent(t *testing.T) {
	filePath := "../testdata/project2/no-id.txt"

	ids, err := searchFileIDs(filePath)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(ids) != 0 {
		t.Fatalf("expected len 0, got %d", len(ids))
	}
}

func TestSearchFileIDsFailsWithInvalidFilePath(t *testing.T) {
	filePath := "./oops/nm/invalid"

	_, err := searchFileIDs(filePath)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
