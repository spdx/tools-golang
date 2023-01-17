// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder

import (
	"time"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// BuildCreationInfoSection creates an SPDX Package, returning that
// package or error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - code: verification code from Package
//   - namespacePrefix: prefix for DocumentNamespace (packageName and code will be added)
//   - creatorType: one of Person, Organization or Tool
//   - creator: creator string
//   - testValues: for testing only; call with nil when using in production
func BuildCreationInfoSection(creatorType string, creator string, testValues map[string]string) (*spdx.CreationInfo, error) {
	// build creator slices
	creators := []common.Creator{
		// add builder as a tool
		{
			Creator:     "github.com/spdx/tools-golang/builder",
			CreatorType: "Tool",
		},
		{
			Creator:     creator,
			CreatorType: creatorType,
		},
	}

	// use test Created time if passing test values
	location, _ := time.LoadLocation("UTC")
	locationTime := time.Now().In(location)
	created := locationTime.Format("2006-01-02T15:04:05Z")
	if testVal := testValues["Created"]; testVal != "" {
		created = testVal
	}

	ci := &spdx.CreationInfo{
		Creators: creators,
		Created:  created,
	}
	return ci, nil
}
