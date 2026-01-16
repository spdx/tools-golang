---
SPDX-License-Identifier: CC-BY-4.0
---

# Contributing

All contributions must include a "Signed-off-by" line in the commit message.

This indicates that the contribution is made pursuant to the [Developer Certificate of Origin (DCO)](https://developercertificate.org/), a copy of which is included below.

## Test coverage

Since this library is intended to be relied upon by other tools to work with SPDX data, we are aiming to ensure that it is and remains well-tested.

PRs with new code should include corresponding test files, and should continue to pass existing tests. Unit tests for `foo.go` should be placed in `foo_test.go`. Test data files and folders should be placed in the top-level `testdata/` folder.

To run the test suite, from the top-level directory run: `go test ./...`

## License information

New **code files** should include a [short-form SPDX ID](https://spdx.org/ids) at the top, indicating the project license for code, which is Apache-2.0 OR GPL-2.0-or-later. This should look like the following:

```go
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
```

New **documentation files** should include a [short-form SPDX ID](https://spdx.org/ids) at the top, indicating the project license for documentation, which is CC-BY-4.0. This should look like the following:

```text
SPDX-License-Identifier: CC-BY-4.0
```

## Developer Certificate of Origin (DCO)

```text
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```
