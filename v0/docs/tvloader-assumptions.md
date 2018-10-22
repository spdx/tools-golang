SPDX-License-Identifier: CC-BY-4.0

The tag-value file loader in `package tvloader` makes the following assumptions:

Document Creation Info
----------------------
* The Document Creation Info section will always come first, and be completed
  first. Although the spec may not make this explicit, it appears that this is
  the intended format. Unless it comes first, the parser will not be able to
  confirm what version of the SPDX spec is being used. And, "SPDXID:" tags are
  used for not just the Document Creation Info section but also for others (e.g.
  Packages, Files).

Relationship
------------
* Relationship sections will begin with the "Relationship" tag.

Annotation
----------
* Annotation sections will begin with the "Annotator" tag.

Other License Info
------------------
* Other License sections will begin with the "LicenseID" tag.

* Any Other License section, if present, will come later than the Document
  Creation Info section and after any Package, File and Snippet sections.

* Additionally, any Other License section will not have Relationship or
  Annotation sections.

* An implication of the preceding two points is that the Other License sections
  should be the final sections in a v2.1 tag-value document, unless it is
  followed by a (deprecated as of v2.1) Review Information section.

Review
------
* Review sections will begin with the "Reviewer" tag.

* Any Review section, if present, will come later than the Document Creation
  Info section and after any Package, File, Snippet, and Other License sections.

* Additionally, any Review section will not have Relationship or
  Annotation sections.

* An implication of the preceding two points is that the Review sections
  should be the final sections in a v2.1 tag-value document.
