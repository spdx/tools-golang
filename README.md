[![Build Status](https://github.com/spdx/tools-golang/workflows/build/badge.svg)](https://github.com/spdx/tools-golang/actions)
[![Coverage Status](https://coveralls.io/repos/github/spdx/tools-golang/badge.svg)](https://coveralls.io/github/spdx/tools-golang)
[![GitHub release](https://img.shields.io/github/release/spdx/tools-golang.svg)](https://github.com/spdx/tools-golang/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/spdx/tools-golang.svg)](https://pkg.go.dev/github.com/spdx/tools-golang)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5710/badge)](https://bestpractices.coreinfrastructure.org/projects/5710)

# SPDX tools-golang

`tools-golang` is a collection of Go packages intended to make it easier for
Go programs to work with [SPDX®](https://spdx.dev/) files.

## What it does

tools-golang currently works with files conformant to versions 2.1, 2.2 and 2.3
of the SPDX specification, available at: https://spdx.dev/specifications

tools-golang provides the following packages:

* *spdx* - in-memory data model for the sections of an SPDX document
* *tagvalue* - tag-value document reader and writer
* *rdf* - RDF document reader
* *json* - JSON document reader and writer
* *yaml* - YAML document reader and writer
* *builder* - builds "empty" SPDX document (with hashes) for directory contents
* *idsearcher* - searches for [SPDX short-form IDs](https://spdx.org/ids/) and builds an SPDX document
* *licensediff* - compares concluded licenses between files in two packages
* *reporter* - generates basic license count report from an SPDX document
* *spdxlib* - various utility functions for manipulating SPDX documents in memory
* *utils* - various utility functions that support the other tools-golang packages

Examples for how to use these packages can be found in the `examples/`
directory.

## What it doesn't do

`tools-golang` doesn't currently support files under any version of the SPDX spec prior to v2.1

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

As indicated in `LICENSE.code`, tools-golang **source code files** are
provided and may be used, at your option, under *either*:
* Apache License, version 2.0 (**Apache-2.0**), **OR**
* GNU General Public License, version 2.0 or later (**GPL-2.0-or-later**).

As indicated in `LICENSE.docs`, tools-golang **documentation files** are
provided and may be used under the Creative Commons Attribution
4.0 International license (**CC-BY-4.0**).

This `README.md` file is documentation:

`SPDX-License-Identifier: CC-BY-4.0`

## Security

For security policy and reporting security issues, please refer to [SECURITY.md](SECURITY.md)
