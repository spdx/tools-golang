// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"fmt"
	"time"

	"github.com/spdx/tools-golang/spdx"
)

// BuildCreationInfoSection2_2 creates an SPDX Package (version 2.2), returning that
// package or error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - code: verification code from Package
//   - namespacePrefix: prefix for DocumentNamespace (packageName and code will be added)
//   - creatorType: one of Person, Organization or Tool
//   - creator: creator string
//   - testValues: for testing only; call with nil when using in production
func BuildCreationInfoSection2_2(packageName string, code string, namespacePrefix string, creatorType string, creator string, testValues map[string]string) (*spdx.CreationInfo2_2, error) {
	// build creator slices
	cPersons := []string{}
	cOrganizations := []string{}
	cTools := []string{}
	// add builder as a tool
	cTools = append(cTools, "github.com/spdx/tools-golang/builder")

	switch creatorType {
	case "Person":
		cPersons = append(cPersons, creator)
	case "Organization":
		cOrganizations = append(cOrganizations, creator)
	case "Tool":
		cTools = append(cTools, creator)
	default:
		cPersons = append(cPersons, creator)
	}

	// use test Created time if passing test values
	location, _ := time.LoadLocation("UTC")
	locationTime := time.Now().In(location)
	created := locationTime.Format("2006-01-02T15:04:05Z")
	if testVal := testValues["Created"]; testVal != "" {
		created = testVal
	}

	ci := &spdx.CreationInfo2_2{
		SPDXVersion:          "SPDX-2.2",
		DataLicense:          "CC0-1.0",
		SPDXIdentifier:       spdx.ElementID("DOCUMENT"),
		DocumentName:         packageName,
		DocumentNamespace:    fmt.Sprintf("%s%s-%s", namespacePrefix, packageName, code),
		CreatorPersons:       cPersons,
		CreatorOrganizations: cOrganizations,
		CreatorTools:         cTools,
		Created:              created,
	}
	return ci, nil
}
