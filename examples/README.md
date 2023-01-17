SPDX-License-Identifier: CC-BY-4.0

# tools-golang Examples

The `examples/` directory contains examples for how to use the various
tools-golang sub-packages. Sample commands below should be run from
within the example's subdirectory.

## 1-load/

*tagvalue*, *spdx*

This example demonstrates loading an SPDX tag-value file from disk into memory,
and printing some of its contents to standard output.
#### Run project: *go run example_load.go ../sample-docs/tv/hello.spdx*

## 2-load-save/

*tagvalue*, *tagvalue*

This example demonstrates loading an SPDX tag-value file from disk into memory,
and re-saving it to a different file on disk.
#### Run project: *go run example_load_save.go ../sample-docs/tv/hello.spdx test.spdx*

## 3-build/

*builder*, *tagvalue*

This example demonstrates building an 'empty' SPDX document in memory that
corresponds to a given directory's contents, including all files with their
hashes and the package's verification code, and saving the document to disk.
#### Run project: *go run example_build.go project2 ../../testdata/project2 test.spdx*

## 4-search/

*idsearcher*, *tagvalue*

This example demonstrates building an SPDX document for a directory's contents
(implicitly using *builder*); searching through that directory for [SPDX
short-form IDs](https://spdx.org/ids/); filling those IDs into the document's
Package and File license fields; and saving the resulting document to disk.
#### Run project: *go run example_search.go project2 ../../testdata/project2/folder test.spdx*

## 5-report/

*reporter*, *tagvalue*

This example demonstrates loading an SPDX tag-value file from disk into memory,
generating a basic report listing counts of the concluded licenses for its
files, and printing the report to standard output.
#### Run project: *go run example_report.go ../sample-docs/tv/hello.spdx*

## 6-licensediff

*licensediff*, *tagvalue*

This example demonstrates loading two SPDX tag-value files from disk into
memory, and generating a diff of the concluded licenses for Files in Packages
with matching IDs in each document.

This is generally only useful when run with two SPDX documents that describe
licenses for subsequent versions of the same set of files, AND if they have
the same identifier in both documents.
#### Run project: *go run example_licensediff.go ../sample-docs/tv/hello.spdx ../sample-docs/tv/hello-modified.spdx*

## 7-rdfloader

*rdfloader*, *spdx*

This example demonstrates loading an SPDX rdf file from disk into memory 
and then printing the corresponding spdx struct for the document.
#### Run project: *go run exampleRDFLoader.go ../sample-docs/rdf/SPDXRdfExample-v2.2.spdx.rdf*

## 8-jsontotv

*json*, *tagvalue*

This example demonstrates loading an SPDX json from disk into memory
and then re-saving it to a different file on disk in tag-value format.
#### Run project: *go run examplejsontotv.go ../sample-docs/json/SPDXJSONExample-v2.2.spdx.json example.spdx*

## 9-tvtojson

*json*, *tagvalue*

This example demonstrates loading an SPDX tag-value from disk into memory
and then re-saving it to a different file on disk in json format.
#### Run project: *go run exampletvtojson.go ../sample-docs/tv/hello.spdx example.json*

## 10-jsonloader

*json*

This example demonstrates loading an SPDX json from disk into memory
and then logging some of the attributes to the console.
#### Run project: *go run example_json_loader.go ../sample-docs/json/SPDXJSONExample-v2.2.spdx.json*

## 11-yamltotv

*yaml* *tagvalue*

This example demonstrates loading an SPDX yaml from disk into memory
and then re-saving it to a different file on disk in tag-value format.
#### Run project: *go run exampleyamltotv.go ../sample-docs/yaml/SPDXYAMLExample-2.2.spdx.yaml test.spdx*

## 12-tvtoyaml

*yaml* *tagvalue*

This example demonstrates loading an SPDX tag-value from disk into memory
and then re-saving it to a different file on disk in yaml format.
#### Run project: *go run exampletvtoyaml.go ../sample-docs/tv/hello.spdx example.yaml*

## 13-yamlloader

*yaml*

This example demonstrates loading an SPDX yaml from disk into memory
and then logging some of the attributes to the console.
#### Run project: *go run exampleYAMLLoader.go ../sample-docs/yaml/SPDXYAMLExample-2.2.spdx.yaml*
