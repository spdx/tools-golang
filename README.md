[![Build Status](https://github.com/spdx/tools-golang/workflows/build/badge.svg)](https://github.com/spdx/tools-golang/actions)
[![Coverage Status](https://coveralls.io/repos/github/spdx/tools-golang/badge.svg)](https://coveralls.io/github/spdx/tools-golang)
[![GitHub release](https://img.shields.io/github/release/spdx/tools-golang.svg)](https://github.com/spdx/tools-golang/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/spdx/tools-golang.svg)](https://pkg.go.dev/github.com/spdx/tools-golang)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5710/badge)](https://bestpractices.coreinfrastructure.org/projects/5710)

# SPDX tools-golang

tools-golang is a collection of Go packages intended to make it easier for
Go programs to work with [SPDX®](https://spdx.dev/) files.

## Recent news

2022-04-03: **v0.3.0**: added support for saving SPDX JSON files as well as
other improvements and bugfixes. See [RELEASE-NOTES.md](./RELEASE-NOTES.md)
for full details.

## What it does

tools-golang currently works with files conformant to versions 2.1 and 2.2
of the SPDX specification, available at: https://spdx.dev/specifications

tools-golang provides the following packages:

* *spdx* - in-memory data model for the sections of an SPDX document
* *tvloader* - tag-value document loader
* *tvsaver* - tag-value document saver
* *rdfloader* - RDF document loader
* *jsonloader* - JSON document loader
* *jsonsaver* - JSON document saver
* *builder* - builds "empty" SPDX document (with hashes) for directory contents
* *idsearcher* - searches for [SPDX short-form IDs](https://spdx.org/ids/) and builds SPDX document
* *licensediff* - compares concluded licenses between files in two packages
* *reporter* - generates basic license count report from SPDX document
* *spdxlib* - various utility functions for manipulating SPDX documents in memory
* *utils* - various utility functions that support the other tools-golang packages

Examples for how to use these packages can be found in the `examples/`
directory.

## What it doesn't do

tools-golang doesn't currently do any of the following:

* work with files under any version of the SPDX spec prior to v2.1
* convert between different versions of SPDX documents (e.g., from 2.1 to 2.2)
* enable applications to interact with SPDX files without needing to care
  (too much) about the particular SPDX file version

We are working towards adding functionality for all of these. Code contributions
are welcome!

## Requirements

tools-golang uses https://github.com/spdx/gordf to manage RDF input and output.

Other than that, tools-golang does not require anything outside the Go standard
library.

## Documentation

SPDX tools-golang documentation is available on the pkg.go.dev website at https://pkg.go.dev/github.com/spdx/tools-golang.

## Contributors

Thank you to all of the contributors to spdx/tools-golang. A full list can be
found in the GitHub repo and in [the release notes](RELEASE-NOTES.md).

In particular, thank you to the following for major contributions:

JSON parsing and saving support was added by @specter25 as part of his Google
Summer of Code 2021 project.

RDF parsing support was added by @RishabhBhatnagar as part of his Google Summer
of Code 2020 project.

## Licenses

As indicated in `LICENSE-code`, tools-golang **source code files** are
provided and may be used, at your option, under *either*:
* Apache License, version 2.0 (**Apache-2.0**), **OR**
* GNU General Public License, version 2.0 or later (**GPL-2.0-or-later**).

As indicated in `LICENSE-docs`, tools-golang **documentation files** are
provided and may be used under the Creative Commons Attribution
4.0 International license (**CC-BY-4.0**).

This `README.md` file is documentation:

`SPDX-License-Identifier: CC-BY-4.0`
