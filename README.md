[![Build Status](https://travis-ci.org/swinslow/spdx-go.svg?branch=master)](https://travis-ci.org/swinslow/spdx-go)

SPDX-License-Identifier: CC-BY-4.0

# spdx-go

spdx-go is a collection of Go packages intended to make it easier for Go
programs to work with [SPDX®](https://spdx.org/) files.

This software is in an extremely pre-alpha state (hence the "v0/" directory).

## What it does

spdx-go currently works with files conformant to version 2.1 of the SPDX
specification, available at: https://spdx.org/specifications

spdx-go provides the following packages:

* *v0/spdx* - in-memory data model for the sections of an SPDX document
* *v0/tvloader* - tag-value file loader

## What it doesn't do

spdx-go doesn't currently do any of the following:

* work with files under any version of the SPDX spec *other than* v2.1
* work with RDF files
* output RDF or tag-value files
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
