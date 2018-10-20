// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v1

import (
	"fmt"
	"time"

	"github.com/swinslow/spdx-go/v0/builder"
	"github.com/swinslow/spdx-go/v0/spdx"
)

// BuildCreationInfoSection2_1 creates an SPDX Package (version 2.1), returning that
// package or error if any is encountered. Arguments:
//   - config: Config object
//   - packageName: name of package / directory
//   - code: verification code from Package
func BuildCreationInfoSection2_1(config *builder.Config2_1, packageName string, code string) (*spdx.CreationInfo2_1, error) {
	if config == nil {
		return nil, fmt.Errorf("got nil config")
	}

	// build creator slices
	cPersons := []string{}
	cOrganizations := []string{}
	cTools := []string{}
	// add builder as a tool
	cTools = append(cTools, "github.com/swinslow/spdx-go/v0/builder")

	switch config.CreatorType {
	case "Person":
		cPersons = append(cPersons, config.Creator)
	case "Organization":
		cOrganizations = append(cOrganizations, config.Creator)
	case "Tool":
		cTools = append(cTools, config.Creator)
	default:
		cPersons = append(cPersons, config.Creator)
	}

	// use test Created time if passing test values
	created := time.Now().Format("2006-01-02T15:04:05Z")
	if testVal := config.TestValues["Created"]; testVal != "" {
		created = testVal
	}

	ci := &spdx.CreationInfo2_1{
		SPDXVersion:          "SPDX-2.1",
		DataLicense:          "CC0-1.0",
		SPDXIdentifier:       "SPDXRef-DOCUMENT",
		DocumentName:         packageName,
		DocumentNamespace:    fmt.Sprintf("%s%s-%s", config.NamespacePrefix, packageName, code),
		CreatorPersons:       cPersons,
		CreatorOrganizations: cOrganizations,
		CreatorTools:         cTools,
		Created:              created,
	}
	return ci, nil
}
