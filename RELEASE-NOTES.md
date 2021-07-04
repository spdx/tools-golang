SPDX-License-Identifier: CC-BY-4.0

# Release Notes for spdx/tools-golang

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
