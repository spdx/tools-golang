// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"testing"
)

// ===== CreationInfo section builder tests =====
func TestBuilder2_2CanBuildCreationInfoSection(t *testing.T) {
	creatorType := "Organization"
	creator := "Jane Doe LLC"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection2_2(creatorType, creator, testValues)
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

func TestBuilder2_2CanBuildCreationInfoSectionWithCreatorPerson(t *testing.T) {
	creatorType := "Person"
	creator := "John Doe"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection2_2(creatorType, creator, testValues)
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

func TestBuilder2_2CanBuildCreationInfoSectionWithCreatorTool(t *testing.T) {
	creatorType := "Tool"
	creator := "some-other-tool-2.1"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection2_2(creatorType, creator, testValues)
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
	if ci.Creators[1].Creator != "some-other-tool-2.1" {
		t.Errorf("expected %s, got %s", "some-other-tool-2.1", ci.Creators[1])
	}
}

func TestBuilder2_2CanBuildCreationInfoSectionWithInvalidPerson(t *testing.T) {
	creatorType := "Whatever"
	creator := "John Doe"
	testValues := make(map[string]string)
	testValues["Created"] = "2018-10-20T16:48:00Z"

	ci, err := BuildCreationInfoSection2_2(creatorType, creator, testValues)
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
