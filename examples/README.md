SPDX-License-Identifier: CC-BY-4.0

# spdx-go Examples

The `examples/` directory contains examples for how to use the various spdx-go sub-packages.

## 1-load/

*tvloader*, *spdx*

This example demonstrates loading an SPDX tag-value file from disk into memory,
and printing some of its contents to standard output.

## 2-load-save/

*tvloader*, *tvsaver*

This example demonstrates loading an SPDX tag-value file from disk into memory,
and re-saving it to a different file on disk.

## 3-build/

*builder*, *tvsaver*

This example demonstrates building an 'empty' SPDX document in memory that
corresponds to a given directory's contents, including all files with their
hashes and the package's verification code, and saving the document to disk.

## 4-search/

*idsearcher*, *tvsaver*

This example demonstrates building an SPDX document for a directory's contents
(implicitly using *builder*); searching through that directory for [SPDX
short-form IDs](https://spdx.org/ids/); filling those IDs into the document's
Package and File license fields; and saving the resulting document to disk.

## 5-report/

*reporter*, *tvloader*

This example demonstrates loading an SPDX tag-value file from disk into memory,
generating a basic report listing counts of the concluded licenses for its
files, and printing the report to standard output.
