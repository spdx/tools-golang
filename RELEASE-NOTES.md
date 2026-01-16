---
SPDX-License-Identifier: CC-BY-4.0
---

# Release Notes for spdx/tools-golang

## 0.5.7

0.5.7 released on: 2026-01-12

### What's Changed

* Fail parsing if the required prefix isn't present for IDs by @awyeth in #275

### New Contributors

* @awyeth made their first contribution in #275

### Contributors

* @awyeth

## 0.5.6

0.5.6 released on: 2025-12-23

### What's Changed

* Change the colon and space separator to just a colon by @warpkwd in #254
* build(deps): Bump sigs.k8s.io/yaml to 1.6.0 by @dependabot[bot] in #259
* Resolve two issues related to ExternalDocumentRef by @KAWAHARA-souta in #269
* fix: Fix prefixDocumentId() to use correct prefix by @KAWAHARA-souta in #272

### New Contributors

* @warpkwd made their first contribution in #254
* @KAWAHARA-souta made their first contribution in #269

### Contributors

* @warpkwd
* @dependabot
* @KAWAHARA-souta

## 0.5.5

0.5.5 released on: 2024-06-25

### What's Changed

* fix: properly normalize Windows paths, add windows test runner by @kzantow in #242
* fix: panic if JSON relationship array contains null by @kzantow in #239
* chore: provide a clearer error when using an invalid Originator by @LaurentGoderre in #246

### New Contributors

* @LaurentGoderre made their first contribution in #246

### Contributors

* @LaurentGoderre
* @kzantow

## 0.5.4

0.5.4 released on: 2024-04-17

### What's Changed

* Stop escaping HTML by @kzantow in #224
* Don't create empty `ExcludedFiles` array by @DmitriyLewen in #230
* Add external reference category `OTHER` by @mcombuechen in #229
* Remove empty packageVerificationCode in 2.2 JSON by @kzantow in #223

### New Contributors

* @mcombuechen made their first contribution in #229

### Contributors

* @mcombuechen
* @kzantow
* @DmitriyLewen

## 0.5.3

0.5.3 released on: 2023-07-27

### What's Changed

* Fix: SPDX decode error originator by @spiffcs in #221

### New Contributors

* @spiffcs made their first contribution in #221

### Contributors

* @spiffcs

## 0.5.2

0.5.2 released on: 2023-06-06

### What's Changed

* Fix duplicate shorthand relationships for opposite case by @lumjjb in #220

### Contributors

* @lumjjb

## 0.5.1

0.5.1 released on: 2023-05-26

### What's Changed

* Add ability to specify JSON output options by @DmitriyLewen in #213
* Fix some optional params: `copyrightText`, `licenseListVersion`, `packageVerificationCode` by @lumjjb in #215
* Properly output and read the filesAnalyzed field in JSON/YAML by @kzantow in #210
* Ensure no duplicates in relationships when shortcut fields are used. by @lumjjb in #218

### New Contributors

* @testwill made their first contribution in #212
* @DmitriyLewen made their first contribution in #213

### Contributors

* @kzantow
* @lumjjb
* @testwill
* @DmitriyLewen

## 0.5.0

0.5.0 released on: 2023-04-03

-rc1 released on: 2023-01-20

This is the first release which includes a significant refactoring of this library and includes the ability to convert between SPDX document versions (2.1 - 2.3).

**NOTE:** This version has a major refactoring how to use the library. This is now much more streamlined. Prior to this change, it was required to import things like `spdx/v2_2` and directly reference those version files. This refactoring moves usage to have a "common model", which ends up being the *latest* SPDX version, available at the same package across releases: `github.com/spdx/tools-golang/spdx`. This means when upgrading versions of tools-golang, you can always get the latest version supported by the library and support reading older versions due to the automatic conversions that the reading functions provide.

To get an idea of what is involved (it really isn't a lot of work), you can have a look at the Syft PR that upgraded to use the new interfaces: <https://github.com/anchore/syft/pull/1503>

After upgrading to this usage pattern, subsequent updates of the tools-golang library will only require changes to your code if the latest model changes (for example, when 3.0 is implemented -- but your older 2.x files will *still* work fine to read in and *export*).

## What's new

* Refactor: maintain the latest SPDX model and provide conversions from previous by @kzantow in #172
* Added more const for external reference to external.go by @neilnaveen in #188

### Bug fixes

* Fixed Bug For DocumentComment by @neilnaveen in #185 and #187
* Improve SPDX document validation by @neilnaveen in #200
* Read shortcut fields: documentDescribes and hasFiles by @kzantow in #201
* JSON reading/writing sets appropriate PACKAGE-MANAGER enum based on version by @lumjjb in #204

### New Contributors

* @jspeed-meyers made their first contribution in #181
* @neilnaveen made their first contribution in #185

### Contributors

* @kzantow
* @lumjjb
* @neilnaveen
* @jspeed-meyers

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
