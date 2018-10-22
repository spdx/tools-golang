SPDX-License-Identifier: CC-BY-4.0

The short-form ID searcher in `package idsearcher` makes the following assumptions:

- For PackageLicenseInfoFromFiles (in Package) and LicenseInfoInFile (in File),
  an exception should be treated as a separate "license". For example, in the
  expression `GPL-2.0-only WITH Classpath-exception-2.0`, each of `GPL-2.0-only`
  and `Classpath-exception-2.0` will be listed separately.
