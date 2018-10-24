// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *idsearcher*, *tvsaver*

// This example demonstrates building an SPDX document for a directory's
// contents (implicitly using *builder*); searching through that directory for
// [SPDX short-form IDs](https://spdx.org/ids/); filling those IDs into the
// document's Package and File license fields; and saving the resulting document
// to disk.

package main

import (
	"fmt"
	"os"

	"github.com/swinslow/spdx-go/v0/idsearcher"

	"github.com/swinslow/spdx-go/v0/tvsaver"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("Usage: %v <package-name> <package-root-dir> <spdx-file-out>\n", args[0])
		fmt.Printf("  Build a SPDX 2.1 document with one package called <package-name>;\n")
		fmt.Printf("  create files with hashes corresponding to the files in <package-root-dir>;\n")
		fmt.Printf("  search for SPDX short-form IDs, and use them to fill in license data\n")
		fmt.Printf("  where possible; and save it out as a tag-value file to <spdx-file-out>.\n")
		return
	}

	// get the command-line arguments
	packageName := args[1]
	packageRootDir := args[2]
	fileOut := args[3]

	// to use the SPDX idsearcher package, the first step is to define a
	// idsearcher.Config2_1 struct. this config data can be reused, in case you
	// are building SPDX documents for several directories in sequence.
	config := &idsearcher.Config{

		// NamespacePrefix is a prefix that will be used to populate the
		// mandatory DocumentNamespace field in the Creation Info section.
		// Because it needs to be unique, the value that will be filled in
		// for the document will have the package name and verification code
		// appended to this prefix.
		NamespacePrefix: "https://example.com/whatever/testdata-",

		// CreatorType and Creator, from builder.Config2_1, are not needed for
		// idsearcher.Config. Because it is automated and doesn't assume
		// further review, the following two Creator fields are filled in:
		// Creator: Tool: github.com/swinslow/spdx-go/v0/builder
		// Creator: Tool: github.com/swinslow/spdx-go/v0/idsearcher

		// You can define one or more paths that should be ignored
		// when walking through the directory. This is intended to omit files
		// that are located within the package's directory, but which should
		// be omitted from the SPDX document.
		// This is directly passed through to builder, and uses the same
		// format as shown in examples/3-build/example_build.go.
		BuilderPathsIgnored: []string{

			// ignore all files in the .git/ directory at the package root
			"/.git/",

			// ignore all files in all __pycache__/ directories, anywhere
			// within the package directory tree
			"**/__pycache__/",

			// ignore the file with this specific path relative to the
			// package root
			"/.ignorefile",

			// or ignore all files with this filename, anywhere within the
			// package directory tree
			"**/.DS_Store",
		},

		// Finally, SearcherPathsIgnored lists certain paths that should not be
		// searched by idsearcher, even if those paths have Files present (and
		// had files filled in by builder). This is useful, for instance, if
		// your project has some directories or files with
		// "SPDX-License-Identifier:" tags, but for one reason or another you
		// want to exclude those files' tags from being picked up by the
		// searcher.
		// SearcherPathsIgnored uses the same format as BuilderPathsIgnored.
		SearcherPathsIgnored: []string{

			// Example for the Linux kernel: ignore the documentation file
			// which explains how to use SPDX short-form IDs (and therefore
			// has a bunch of "SPDX-License-Identifier:" tags that we wouldn't
			// want to pick up).
			"/Documentation/process/license-rules.rst",

			// Similar example for the Linux kernel: ignore all files in the
			// /LICENSES/ directory.
			"/LICENSES/",
		},
	}

	// now, when we actually ask idsearcher to walk through a directory and
	// build an SPDX document, we need to give it three things:
	//   - what to name the package; and
	//   - where the directory is located on disk; and
	//   - the config object we just defined.
	// these are the same arguments needed for builder, and in fact they get
	// passed through to builder (with the relevant data from the config
	// object extracted behind the scenes).
	doc, err := idsearcher.BuildIDsDocument(packageName, packageRootDir, config)
	if err != nil {
		fmt.Printf("Error while building document: %v\n", err)
		return
	}

	// if we got here, the document has been created.
	// all file hashes and the package verification code have been filled in
	// appropriately by builder.
	// And, all files with "SPDX-License-Identifier:" tags have had their
	// licenses extracted into LicenseInfoInFile and LicenseConcluded for
	// each file by idsearcher. The PackageLicenseInfoFromFiles field will
	// also be filled in with all license identifiers.
	fmt.Printf("Successfully created document and searched for IDs for package %s\n", packageName)

	// NOTE that BuildIDsDocument does NOT do any validation of the license
	// identifiers, to confirm that they are e.g. on the SPDX License List
	// or in other appropriate format (e.g., LicenseRef-...)

	// we can now save it to disk, using tvsaver.

	// create a new file for writing
	w, err := os.Create(fileOut)
	if err != nil {
		fmt.Printf("Error while opening %v for writing: %v\n", fileOut, err)
		return
	}
	defer w.Close()

	err = tvsaver.Save2_1(doc, w)
	if err != nil {
		fmt.Printf("Error while saving %v: %v", fileOut, err)
		return
	}

	fmt.Printf("Successfully saved %v\n", fileOut)
}
