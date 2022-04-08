// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
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
func BuildCreationInfoSection2_2(creatorType string, creator string, testValues map[string]string) (*spdx.CreationInfo2_2, error) {
	// build creator slices
	creators := []spdx.Creator{
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

	ci := &spdx.CreationInfo2_2{
		Creators: creators,
		Created:  created,
	}
	return ci, nil
}
