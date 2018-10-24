// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Example for: *builder*, *tvsaver*

// This example demonstrates building an 'empty' SPDX document in memory that
// corresponds to a given directory's contents, including all files with their
// hashes and the package's verification code, and saving the document to disk.

package main

import (
	"fmt"
	"os"

	"github.com/swinslow/spdx-go/v0/builder"
	"github.com/swinslow/spdx-go/v0/tvsaver"
)

func main() {

	// check that we've received the right number of arguments
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("Usage: %v <package-name> <package-root-dir> <spdx-file-out>\n", args[0])
		fmt.Printf("  Build a SPDX 2.1 document with one package called <package-name>;\n")
		fmt.Printf("  create files with hashes corresponding to the files in <package-root-dir>;\n")
		fmt.Printf("  and save it out as a tag-value file to <spdx-file-out>.\n")
		return
	}

	// get the command-line arguments
	packageName := args[1]
	packageRootDir := args[2]
	fileOut := args[3]

	// to use the SPDX builder package, the first step is to define a
	// builder.Config2_1 struct. this config data can be reused, in case you
	// are building SPDX documents for several directories in sequence.
	config := &builder.Config2_1{

		// NamespacePrefix is a prefix that will be used to populate the
		// mandatory DocumentNamespace field in the Creation Info section.
		// Because it needs to be unique, the value that will be filled in
		// for the document will have the package name and verification code
		// appended to this prefix.
		NamespacePrefix: "https://example.com/whatever/testdata-",

		// CreatorType will be used for the first part of the Creator field
		// in the Creation Info section. Per the SPDX spec, it can be
		// "Person", "Organization" or "Tool".
		CreatorType: "Person",

		// Creator will be used for the second part of the Creator field in
		// the Creation Info section.
		Creator: "Jane Doe",

		// note that builder will also add the following, in addition to the
		// Creator defined above:
		// Creator: Tool: github.com/swinslow/spdx-go/v0/builder

		// Finally, you can define one or more paths that should be ignored
		// when walking through the directory. This is intended to omit files
		// that are located within the package's directory, but which should
		// be omitted from the SPDX document.
		PathsIgnored: []string{

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
	}

	// now, when we actually ask builder to walk through a directory and
	// build an SPDX document, we need to give it three things:
	//   - what to name the package; and
	//   - where the directory is located on disk; and
	//   - the config object we just defined.
	doc, err := builder.Build2_1(packageName, packageRootDir, config)
	if err != nil {
		fmt.Printf("Error while building document: %v\n", err)
		return
	}

	// if we got here, the document has been created.
	// all license info is marked as NOASSERTION, but file hashes and
	// the package verification code have been filled in appropriately.
	fmt.Printf("Successfully created document for package %s\n", packageName)

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
