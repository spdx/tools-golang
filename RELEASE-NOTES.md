SPDX-License-Identifier: CC-BY-4.0

# Release Notes for spdx/tools-golang

## 0.4.0

0.4.0 released on 2022-01-11

### New Features and Enhancements
* SPDX v2.3 support #164
* YAML support #134
* Add reference types enumerables to SPDX pkg definition #162 #163
* Expand hash algorithm support to include all valid SPDX 2.2 and 2.3 algorithms #173

### Bug fixes
* JSON encoding and decoding not properly handling SPDXRef- prefixes #170

### Documentation and Cleanup
* Overhaul structs, refactor JSON parser and saver #133 
* YAML documentation and JSON documentation fixes #141 
* Convert SPDX structs to versioned pkgs #146
* Ensure consistency between JSON struct tags across different SPDX versions #174
* Add Security.md for handling of security issues #154
* Update build workflow to go 1.18 #148

### Contributors
* @ianling 
* @CatalinStratu
* @lumjjb 
* @pxp928 
* @kzantow 
* @puerco 
* @jedevc 

## 0.3.0

0.3.0 released on: 2022-04-03

-rc1 released on: 2022-03-27

### New Features and Enhancements
* Add support for saving SPDX JSON: #92, #94, #97, #98, #104, #106, #113
* Begin OpenSSF Best Practices process and add initial badge: #111
  * also enabled branch protection for main branch

### Bug fixes
* tvsaver: Fix incorrect tag for Snippet IDs: #95
* GitHub Actions: Fix incorrect branch for code coverage: #112
* builder: Fix file paths to be relative rather than absolute: #114
* builder: Add missing mandatory field LicenseInfoInFile: #119

### Documentation and Cleanup
* Fix link to release notes: #91
* Language fixes for JSON documentation: #108
* Add badges and links for releases and documentation: #109
* Update documentation for release: #121, #122
* Fixes for examples and sample run commands: #123, #125, #126, #127

### Contributors
* @CatalinStratu
* @specter25
* @swinslow

## 0.2.0

Released on: 2021-07-04

### New Features and Enhancements
* Add support for parsing SPDX JSON: #72, #75, #83, #84, #87
  * bug fixes in interim versions: #77, #78, #79, #80, #81, #82
* Improve handling of multiple hash checksum types: #41, #49, #60
* Enable filtering relationships by various relationship types: #71, #74
* Improve package license visibility: #65, #66
* Rename primary branch to 'main': #69
* Add release notes and push release: #85, #90

### Bug fixes
* Fix multiline (`<text>`) wrapping for various fields: #31, #53, #58, #89, #76
* Fix special SPDX IDs in right-hand side of Relationships: #59, #63, #68
* Throw error when parsing tag-value elements without SPDX IDs: #26, #64
* Fix missing colon in 'excludes' for Package Verification Code when saving tag-value documents: #86, #88
* Fix incorrect license statement: #70

### Contributors
* @autarch
* @bisakhmondal
* @ianling
* @matthewkmayer
* @RishabhBhatnagar
* @specter25
* @swinslow

## 0.1.0

Released on: 2021-03-20

### Contributors
* @abhishekspeer
* @goneall
* @RishabhBhatnagar
* @rtgdk
* @swinslow
