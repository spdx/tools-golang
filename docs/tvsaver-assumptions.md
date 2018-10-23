SPDX-License-Identifier: CC-BY-4.0

The tag-value file saver in `package tvsaver` makes the following assumptions:

Document Creation Info
----------------------
* Mandatory fields will be treated the same way as optional fields; if they are
  set to the empty string, they will be omitted. Thus, an invalid Creation Info
  section (e.g. one that doesn't include a correct SPDXVersion field) will
  result in outputting an invalid Creation Info section.

Relationship
------------
* Same comment as above re: optional fields, for RelationshipComment.
