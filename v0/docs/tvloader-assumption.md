SPDX-License-Identifier: CC-BY-4.0

The tag-value file loader in `package tvloader` makes the following assumptions:

* The Document Creation Info section will always come first, and be completed first.
Although the spec may not make this explicit, it appears that this is the intended
format. Unless it comes first, the parser will not be able to confirm what version
of the SPDX spec is being used. And, "SPDXID:" tags are used for not just the
Document Creation Info section but also for others (e.g. Packages, Files).
