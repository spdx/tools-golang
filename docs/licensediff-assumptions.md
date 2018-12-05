SPDX-License-Identifier: CC-BY-4.0

The License diff tool in `package licensediff` makes the following assumptions:

- In any single Package, a given filename will only appear once. This may or may
  not be required by the SPDX spec, but it's kind of implicit in being able to
  create a diff indexed by filename.