// Package builder is used to create spdx-go data structures for a given
// directory path's contents, with hashes, etc. filled in and with empty
// license data.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package builder

import (
	"github.com/swinslow/spdx-go/v0/builder/builder2v1"
	"github.com/swinslow/spdx-go/v0/spdx"
)

// Config2_1 is a collection of configuration settings for docbuilder
// (for version 2.1 SPDX Documents). A few mandatory fields are set here
// so that they can be repeatedly reused in multiple calls to Build2_1.
type Config2_1 struct {
	// NamespacePrefix should be a URI representing a prefix for the
	// namespace with which the SPDX Document will be associated.
	// It will be used in the DocumentNamespace field in the CreationInfo
	// section, followed by the per-Document package name and a random UUID.
	NamespacePrefix string

	// CreatorType should be one of "Person", "Organization" or "Tool".
	// If not one of those strings, it will be interpreted as "Person".
	CreatorType string

	// Creator will be filled in for the given CreatorType.
	Creator string

	// TestValues is used to pass fixed values for testing purposes
	// only, and should be set to nil for production use. It is only
	// exported so that it will be accessible within docbuilder2v1.
	TestValues map[string]string
}

// Build2_1 creates an SPDX Document (version 2.1), returning that document or
// error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - dirRoot: path to directory to be analyzed
//   - config: Config object
func Build2_1(packageName string, dirRoot string, config *Config2_1) (*spdx.Document2_1, error) {
	// build Package section first -- will include Files and make the
	// package verification code available
	pkg, err := builder2v1.BuildPackageSection2_1(packageName, dirRoot)
	if err != nil {
		return nil, err
	}

	ci, err := builder2v1.BuildCreationInfoSection2_1(packageName, pkg.PackageVerificationCode, config.NamespacePrefix, config.CreatorType, config.Creator, config.TestValues)
	if err != nil {
		return nil, err
	}

	rln, err := builder2v1.BuildRelationshipSection2_1(packageName)
	if err != nil {
		return nil, err
	}

	doc := &spdx.Document2_1{
		CreationInfo:  ci,
		Packages:      []*spdx.Package2_1{pkg},
		Relationships: []*spdx.Relationship2_1{rln},
	}

	return doc, nil
}
