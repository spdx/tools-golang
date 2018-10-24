// Package idsearcher is used to search for short-form IDs in files
// within a directory, and to build an SPDX Document containing those
// license findings.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package idsearcher

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/swinslow/spdx-go/v0/builder"
	"github.com/swinslow/spdx-go/v0/spdx"
	"github.com/swinslow/spdx-go/v0/utils"
)

// Config is a collection of configuration settings for docbuilder
// (for version 2.1 SPDX Documents). A few mandatory fields are set here
// so that they can be repeatedly reused in multiple calls to Build2_1.
type Config struct {
	// NamespacePrefix should be a URI representing a prefix for the
	// namespace with which the SPDX Document will be associated.
	// It will be used in the DocumentNamespace field in the CreationInfo
	// section, followed by the per-Document package name and a random UUID.
	NamespacePrefix string

	// BuilderPathsIgnored lists certain paths to be omitted from the built
	// document. Each string should be a path, relative to the package's
	// dirRoot, to a specific file or (for all files in a directory) ending
	// in a slash. Prefix the string with "**" to omit all instances of that
	// file / directory, regardless of where it is in the file tree.
	BuilderPathsIgnored []string

	// SearcherPathsIgnored lists certain paths that should not be searched
	// by idsearcher, even if those paths have Files present. It uses the
	// same format as BuilderPathsIgnored.
	SearcherPathsIgnored []string
}

// BuildIDsDocument creates an SPDX Document (version 2.1) and searches for
// short-form IDs in each file, filling in license fields as appropriate. It
// returns that document or error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - dirRoot: path to directory to be analyzed
//   - namespacePrefix: URI representing a prefix for the
//     namespace with which the SPDX Document will be associated
func BuildIDsDocument(packageName string, dirRoot string, idconfig *Config) (*spdx.Document2_1, error) {
	// first, build the Document using builder
	bconfig := &builder.Config2_1{
		NamespacePrefix: idconfig.NamespacePrefix,
		CreatorType:     "Tool",
		Creator:         "github.com/swinslow/spdx-go/v0/idsearcher",
		PathsIgnored:    idconfig.BuilderPathsIgnored,
	}
	doc, err := builder.Build2_1(packageName, dirRoot, bconfig)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("builder returned nil Document")
	}
	if doc.Packages == nil {
		return nil, fmt.Errorf("builder returned nil Package")
	}
	if len(doc.Packages) != 1 {
		return nil, fmt.Errorf("builder returned %d Packages", len(doc.Packages))
	}

	// now, walk through each file and find its licenses (if any)
	pkg := doc.Packages[0]
	if pkg.Files == nil {
		return nil, fmt.Errorf("builder returned nil Files in Package")
	}
	licsForPackage := map[string]int{}
	for _, f := range pkg.Files {
		// start by initializing / clearing values
		f.LicenseInfoInFile = []string{"NOASSERTION"}
		f.LicenseConcluded = "NOASSERTION"

		// check whether the searcher should ignore this file
		if utils.ShouldIgnore(f.FileName, idconfig.SearcherPathsIgnored) {
			continue
		}

		fPath := filepath.Join(dirRoot, f.FileName)
		// FIXME this is not preferable -- ignoring error
		ids, _ := searchFileIDs(fPath)
		// FIXME for now, proceed onwards with whatever IDs we obtained.
		// FIXME instead of ignoring the error, should probably either log it,
		// FIXME and/or enable the caller to configure what should happen.

		// separate out for this file's licenses
		licsForFile := map[string]int{}
		licsParens := []string{}
		for _, lid := range ids {
			// get individual elements and add for file and package
			licElements := getIndividualLicenses(lid)
			for _, elt := range licElements {
				licsForFile[elt] = 1
				licsForPackage[elt] = 1
			}
			// parenthesize if needed and add to slice for joining
			licsParens = append(licsParens, makeElement(lid))
		}

		// OK -- now we can fill in the file's details, or NOASSERTION if none
		if len(licsForFile) > 0 {
			f.LicenseInfoInFile = []string{}
			for lic := range licsForFile {
				f.LicenseInfoInFile = append(f.LicenseInfoInFile, lic)
			}
			sort.Strings(f.LicenseInfoInFile)
			// avoid adding parens and joining for single-ID items
			if len(licsParens) == 1 {
				f.LicenseConcluded = ids[0]
			} else {
				f.LicenseConcluded = strings.Join(licsParens, " AND ")
			}
		}
	}

	// and finally, we can fill in the package's details
	if len(licsForPackage) == 0 {
		pkg.PackageLicenseInfoFromFiles = []string{"NOASSERTION"}
	} else {
		pkg.PackageLicenseInfoFromFiles = []string{}
		for lic := range licsForPackage {
			pkg.PackageLicenseInfoFromFiles = append(pkg.PackageLicenseInfoFromFiles, lic)
		}
		sort.Strings(pkg.PackageLicenseInfoFromFiles)
	}

	return doc, nil
}

// ===== Utility functions =====
func searchFileIDs(filePath string) ([]string, error) {
	idsMap := map[string]int{}
	ids := []string{}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// build the scan string this way, so that we can run idsearcher on itself
	// without picking up these lines as IDs...
	tag1 := "SPDX-License-"
	tag2 := "Identifier:"
	tag := tag1 + tag2

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), tag) {
			strs := strings.SplitAfterN(scanner.Text(), tag, 2)
			// stop before trailing */ if it is present
			lidToExtract := strs[1]
			lidToExtract = strings.Split(lidToExtract, "*/")[0]
			lid := strings.TrimSpace(lidToExtract)
			lid = stripTrash(lid)
			idsMap[lid] = 1
		}
	}

	// FIXME for now, ignore scanner errors because we want to return whatever
	// FIXME IDs were in fact found. should probably be changed to either
	// FIXME log the error, and/or be configurable for what should happen.
	// if err = scanner.Err(); err != nil {
	// 	return nil, err
	// }

	// now, convert map to string
	for lid := range idsMap {
		ids = append(ids, lid)
	}

	// and sort it
	sort.Strings(ids)

	return ids, nil
}

func stripTrash(lid string) string {
	re := regexp.MustCompile(`[^\w\s\d.\-\+()]+`)
	return re.ReplaceAllString(lid, "")
}

func makeElement(lic string) string {
	if strings.Contains(lic, " AND ") || strings.Contains(lic, " OR ") {
		return fmt.Sprintf("(%s)", lic)
	}

	return lic
}

func getIndividualLicenses(lic string) []string {
	// replace parens and '+' with spaces
	lic = strings.Replace(lic, "(", " ", -1)
	lic = strings.Replace(lic, ")", " ", -1)
	lic = strings.Replace(lic, "+", " ", -1)

	// now, split by spaces, trim, and add to slice
	licElements := strings.Split(lic, " ")
	lics := []string{}
	for _, elt := range licElements {
		elt := strings.TrimSpace(elt)
		// don't add if empty or if case-insensitive operator
		if elt == "" || strings.EqualFold(elt, "AND") ||
			strings.EqualFold(elt, "OR") || strings.EqualFold(elt, "WITH") {
			continue
		}

		lics = append(lics, elt)
	}

	// sort before returning
	sort.Strings(lics)
	return lics
}
