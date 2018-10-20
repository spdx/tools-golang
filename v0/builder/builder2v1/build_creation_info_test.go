// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"fmt"
	"testing"

	"github.com/swinslow/spdx-go/v0/builder"
)

// ===== CreationInfo section builder tests =====
func TestBuilder2_1CanBuildCreationInfoSection(t *testing.T) {
	config := &builder.Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-whatever-",
		CreatorType:     "Organization",
		Creator:         "Jane Doe LLC",
		TestValues:      make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-20T16:48:00Z"
	packageName := "project1"
	verificationCode := "TESTCODE"

	ci, err := buildCreationInfoSection2_1(config, packageName, verificationCode)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if ci.SPDXVersion != "SPDX-2.1" {
		t.Errorf("expected %s, got %s", "SPDX-2.1", ci.SPDXVersion)
	}
	if ci.DataLicense != "CC0-1.0" {
		t.Errorf("expected %s, got %s", "CC0-1.0", ci.DataLicense)
	}
	if ci.SPDXIdentifier != "SPDXRef-DOCUMENT" {
		t.Errorf("expected %s, got %s", "SPDXRef-DOCUMENT", ci.SPDXIdentifier)
	}
	if ci.DocumentName != "project1" {
		t.Errorf("expected %s, got %s", "project1", ci.DocumentName)
	}
	wantNamespace := fmt.Sprintf("https://github.com/swinslow/spdx-docs/spdx-go/testdata-whatever-project1-%s", verificationCode)
	if ci.DocumentNamespace != wantNamespace {
		t.Errorf("expected %s, got %s", wantNamespace, ci.DocumentNamespace)
	}
	if len(ci.CreatorPersons) != 0 {
		t.Fatalf("expected %d, got %d", 0, len(ci.CreatorPersons))
	}
	if len(ci.CreatorOrganizations) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(ci.CreatorOrganizations))
	}
	if ci.CreatorOrganizations[0] != "Jane Doe LLC" {
		t.Errorf("expected %s, got %s", "Jane Doe LLC", ci.CreatorOrganizations[0])
	}
	if len(ci.CreatorTools) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(ci.CreatorTools))
	}
	if ci.CreatorTools[0] != "github.com/swinslow/spdx-go/v0/builder" {
		t.Errorf("expected %s, got %s", "github.com/swinslow/spdx-go/v0/builder", ci.CreatorTools[0])
	}
	if ci.Created != "2018-10-20T16:48:00Z" {
		t.Errorf("expected %s, got %s", "2018-10-20T16:48:00Z", ci.Created)
	}
}

func TestBuilder2_1BuildCreationInfoSectionFailsWithNilConfig(t *testing.T) {
	_, err := buildCreationInfoSection2_1(nil, "hi", "code")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestBuilder2_1CanBuildCreationInfoSectionWithCreatorPerson(t *testing.T) {
	config := &builder.Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-whatever-",
		CreatorType:     "Person",
		Creator:         "John Doe",
		TestValues:      make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-20T16:48:00Z"
	packageName := "project1"
	verificationCode := "TESTCODE"

	ci, err := buildCreationInfoSection2_1(config, packageName, verificationCode)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.CreatorPersons) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(ci.CreatorPersons))
	}
	if ci.CreatorPersons[0] != "John Doe" {
		t.Errorf("expected %s, got %s", "John Doe", ci.CreatorPersons[0])
	}
	if len(ci.CreatorOrganizations) != 0 {
		t.Fatalf("expected %d, got %d", 0, len(ci.CreatorOrganizations))
	}
	if len(ci.CreatorTools) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(ci.CreatorTools))
	}
	if ci.CreatorTools[0] != "github.com/swinslow/spdx-go/v0/builder" {
		t.Errorf("expected %s, got %s", "github.com/swinslow/spdx-go/v0/builder", ci.CreatorTools[0])
	}
}

func TestBuilder2_1CanBuildCreationInfoSectionWithCreatorTool(t *testing.T) {
	config := &builder.Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-whatever-",
		CreatorType:     "Tool",
		Creator:         "some-other-tool-2.1",
		TestValues:      make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-20T16:48:00Z"
	packageName := "project1"
	verificationCode := "TESTCODE"

	ci, err := buildCreationInfoSection2_1(config, packageName, verificationCode)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.CreatorPersons) != 0 {
		t.Fatalf("expected %d, got %d", 0, len(ci.CreatorPersons))
	}
	if len(ci.CreatorOrganizations) != 0 {
		t.Fatalf("expected %d, got %d", 0, len(ci.CreatorOrganizations))
	}
	if len(ci.CreatorTools) != 2 {
		t.Fatalf("expected %d, got %d", 2, len(ci.CreatorTools))
	}
	if ci.CreatorTools[0] != "github.com/swinslow/spdx-go/v0/builder" {
		t.Errorf("expected %s, got %s", "github.com/swinslow/spdx-go/v0/builder", ci.CreatorTools[0])
	}
	if ci.CreatorTools[1] != "some-other-tool-2.1" {
		t.Errorf("expected %s, got %s", "some-other-tool-2.1", ci.CreatorTools[1])
	}
}

func TestBuilder2_1CanBuildCreationInfoSectionWithInvalidPerson(t *testing.T) {
	config := &builder.Config2_1{
		NamespacePrefix: "https://github.com/swinslow/spdx-docs/spdx-go/testdata-whatever-",
		CreatorType:     "Whatever",
		Creator:         "John Doe",
		TestValues:      make(map[string]string),
	}
	config.TestValues["Created"] = "2018-10-20T16:48:00Z"
	packageName := "project1"
	verificationCode := "TESTCODE"

	ci, err := buildCreationInfoSection2_1(config, packageName, verificationCode)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if ci == nil {
		t.Fatalf("expected non-nil CreationInfo, got nil")
	}
	if len(ci.CreatorPersons) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(ci.CreatorPersons))
	}
	if ci.CreatorPersons[0] != "John Doe" {
		t.Errorf("expected %s, got %s", "John Doe", ci.CreatorPersons[0])
	}
	if len(ci.CreatorOrganizations) != 0 {
		t.Fatalf("expected %d, got %d", 0, len(ci.CreatorOrganizations))
	}
	if len(ci.CreatorTools) != 1 {
		t.Fatalf("expected %d, got %d", 1, len(ci.CreatorTools))
	}
	if ci.CreatorTools[0] != "github.com/swinslow/spdx-go/v0/builder" {
		t.Errorf("expected %s, got %s", "github.com/swinslow/spdx-go/v0/builder", ci.CreatorTools[0])
	}
}
