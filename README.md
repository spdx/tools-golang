[![Build Status](https://travis-ci.org/swinslow/spdx-go.svg?branch=master)](https://travis-ci.org/swinslow/spdx-go)
[![Coverage Status](https://coveralls.io/repos/github/swinslow/spdx-go/badge.svg)](https://coveralls.io/github/swinslow/spdx-go)

SPDX-License-Identifier: CC-BY-4.0

# spdx-go

spdx-go is a collection of Go packages intended to make it easier for Go
programs to work with [SPDX®](https://spdx.org/) files.

This software is in an early state, and its API may change significantly (hence the "v0/" directory).

## What it does

spdx-go currently works with files conformant to version 2.1 of the SPDX
specification, available at: https://spdx.org/specifications

spdx-go provides the following packages:

* *v0/spdx* - in-memory data model for the sections of an SPDX document
* *v0/tvloader* - tag-value file loader
* *v0/tvsaver* - tag-value file saver
* *v0/builder* - builds "empty" SPDX document (with hashes) for directory contents
* *v0/idsearcher* - searches for [SPDX short-form IDs](https://spdx.org/ids/) and builds SPDX document
* *v0/reporter* - generates basic license count report from SPDX document
* *v0/utils* - various utility functions that support the other spdx-go packages

Examples for how to use these packages can be found in the `examples/` directory.

## What it doesn't do

spdx-go doesn't currently do any of the following:

* work with files under any version of the SPDX spec *other than* v2.1
* work with RDF files
* convert between RDF and tag-value files, or between different versions
* enable applications to interact with SPDX files without needing to care
  (too much) about the particular SPDX file version

As a long-term goal, I am hoping to enable each of these. Code contributions
are welcome!

## Requirements

At present, spdx-go does not require anything outside the Go standard library.

## Licenses

As indicated in `LICENSE-code.txt`, spdx-go **source code files** are provided
and may be used, at your option, under *either*:
* Apache License, version 2.0 (**Apache-2.0**), **OR**
* GNU General Public License, version 2.0 or later (**GPL-2.0-or-later**).

As indicated in `LICENSE-docs.txt`, spdx-go **documentation files** are
provided and may be used under the Creative Commons Attribution
4.0 International license (**CC-BY-4.0**).

This `README.md` file is documentation, hence the CC-BY-4.0 license ID at
the top of the file.
