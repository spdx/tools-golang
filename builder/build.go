// Package builder is used to create tools-golang data structures for a given
// directory path's contents, with hashes, etc. filled in and with empty
// license data.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package builder

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

// Config is a collection of configuration settings for builder.
// A few mandatory fields are set here
// so that they can be repeatedly reused in multiple calls to Build.
type Config struct {
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

	// PathsIgnored lists certain paths to be omitted from the built document.
	// Each string should be a path, relative to the package's dirRoot,
	// to a specific file or (for all files in a directory) ending in a slash.
	// Prefix the string with "**" to omit all instances of that file /
	// directory, regardless of where it is in the file tree.
	PathsIgnored []string

	// TestValues is used to pass fixed values for testing purposes
	// only, and should be set to nil for production use. It is only
	// exported so that it will be accessible within builder.
	TestValues map[string]string
}

// Build creates an SPDX Document, returning that document or
// error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - dirRoot: path to directory to be analyzed
//   - config: Config object
func Build(packageName string, dirRoot string, config *Config) (*spdx.Document, error) {
	// build Package section first -- will include Files and make the
	// package verification code available
	pkg, err := BuildPackageSection(packageName, dirRoot, config.PathsIgnored)
	if err != nil {
		return nil, err
	}

	ci, err := BuildCreationInfoSection(config.CreatorType, config.Creator, config.TestValues)
	if err != nil {
		return nil, err
	}

	rln, err := BuildRelationshipSection(packageName)
	if err != nil {
		return nil, err
	}

	var packageVerificationCode common.PackageVerificationCode

	if pkg.PackageVerificationCode != nil {
		packageVerificationCode = *pkg.PackageVerificationCode
	}

	doc := &spdx.Document{
		SPDXVersion:       spdx.Version,
		DataLicense:       spdx.DataLicense,
		SPDXIdentifier:    common.ElementID("DOCUMENT"),
		DocumentName:      packageName,
		DocumentNamespace: fmt.Sprintf("%s%s-%s", config.NamespacePrefix, packageName, packageVerificationCode),
		CreationInfo:      ci,
		Packages:          []*spdx.Package{pkg},
		Relationships:     []*spdx.Relationship{rln},
	}

	return doc, nil
}
