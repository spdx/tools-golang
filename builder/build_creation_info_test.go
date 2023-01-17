// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"testing"
)

// ===== CreationInfo section builder tests =====
func TestBuilderCanBuildCreationInfoSection(t *testing.T) {
	creatorType := "Organization"
	creator := "Jane Doe LLC"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection(creatorType, creator, testValues)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.Creators) != 2 {
		t.Fatalf("expected %d, got %d", 2, len(ci.Creators))
	}
	if ci.Creators[1].Creator != "Jane Doe LLC" {
		t.Errorf("expected %s, got %s", "Jane Doe LLC", ci.Creators[0].Creator)
	}
	if ci.Creators[0].Creator != "github.com/spdx/tools-golang/builder" {
		t.Errorf("expected %s, got %s", "github.com/spdx/tools-golang/builder", ci.Creators[1].Creator)
	}
	if ci.Created != "2018-10-20T16:48:00Z" {
		t.Errorf("expected %s, got %s", "2018-10-20T16:48:00Z", ci.Created)
	}
}

func TestBuilderCanBuildCreationInfoSectionWithCreatorPerson(t *testing.T) {
	creatorType := "Person"
	creator := "John Doe"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection(creatorType, creator, testValues)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.Creators) != 2 {
		t.Fatalf("expected %d, got %d", 2, len(ci.Creators))
	}
	if ci.Creators[1].Creator != "John Doe" {
		t.Errorf("expected %s, got %s", "John Doe", ci.Creators[0].Creator)
	}
	if ci.Creators[0].Creator != "github.com/spdx/tools-golang/builder" {
		t.Errorf("expected %s, got %s", "github.com/spdx/tools-golang/builder", ci.Creators[1].Creator)
	}
}

func TestBuilderCanBuildCreationInfoSectionWithCreatorTool(t *testing.T) {
	creatorType := "Tool"
	creator := "some-other-tool"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection(creatorType, creator, testValues)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.Creators) != 2 {
		t.Fatalf("expected %d, got %d", 2, len(ci.Creators))
	}
	if ci.Creators[0].Creator != "github.com/spdx/tools-golang/builder" {
		t.Errorf("expected %s, got %s", "github.com/spdx/tools-golang/builder", ci.Creators[0])
	}
	if ci.Creators[1].Creator != "some-other-tool" {
		t.Errorf("expected %s, got %s", "some-other-tool", ci.Creators[1])
	}
}

func TestBuilderCanBuildCreationInfoSectionWithInvalidPerson(t *testing.T) {
	creatorType := "Whatever"
	creator := "John Doe"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection(creatorType, creator, testValues)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.Creators) != 2 {
		t.Fatalf("expected %d, got %d", 2, len(ci.Creators))
	}
	if ci.Creators[1].Creator != "John Doe" {
		t.Errorf("expected %s, got %s", "John Doe", ci.Creators[1])
	}
	if ci.Creators[0].Creator != "github.com/spdx/tools-golang/builder" {
		t.Errorf("expected %s, got %s", "github.com/spdx/tools-golang/builder", ci.Creators[0])
	}
}
